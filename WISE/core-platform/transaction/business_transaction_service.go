package transaction

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/wiseco/core-platform/services/invoice"
	"github.com/wiseco/core-platform/services/payment"

	grpcTypes "github.com/golang/protobuf/ptypes"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/wiseco/core-platform/auth"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking"
	busBanking "github.com/wiseco/core-platform/services/banking/business"
	bsrv "github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/services/contact"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/go-lib/grpc"
	"github.com/wiseco/go-lib/id"
	"github.com/wiseco/go-lib/num"
	"github.com/wiseco/protobuf/golang"
	"github.com/wiseco/protobuf/golang/shopping/shopify"
	grpcTxn "github.com/wiseco/protobuf/golang/transaction"
	grpcBankTxn "github.com/wiseco/protobuf/golang/transaction/bank"
)

// Service is the transaction service
type BusinessService interface {
	// Create transaction
	Create(BusinessPostedTransactionCreate) (*BusinessPostedTransaction, error)

	// Log transaction
	Log(*BusinessPostedTransactionCreate, *BusinessCardTransactionCreate, *BusinessHoldTransactionCreate) error

	// Fetch business transactions
	ListAllInternal(params map[string]interface{}, businessID shared.BusinessID, accountID string) ([]BusinessPostedTransaction, error)
	ListAll(params map[string]interface{}, userID shared.UserID, businessID shared.BusinessID, accountID string) ([]BusinessPostedTransaction, error)
	GetByIDInternal(shared.PostedTransactionID) (*BusinessPostedTransaction, error)
	GetByID(shared.PostedTransactionID, shared.UserID, shared.BusinessID) (*BusinessPostedTransaction, error)

	// Export csv
	ExportInternal(params map[string]interface{}) (*CSVTransaction, error)
	Export(userID shared.UserID, businessID shared.BusinessID, startDate, endDate string) (*CSVTransaction, error)

	Update(BusinessPostedTransactionUpdate, shared.UserID) (*BusinessPostedTransaction, error)
}

type transactionStore struct {
	cardService    CardService
	accountService AccountService
	*sqlx.DB
}

// New returns a new transaction service
func NewBusinessService() BusinessService {
	return &transactionStore{NewCardService(), NewAccountService(), DBWrite}
}

func getTransactionServiceClient() (grpcBankTxn.TransactionServiceClient, grpc.Client, error) {
	if os.Getenv("USE_TRANSACTION_SERVICE") == "true" {
		sn, err := grpc.GetConnectionStringForService(grpc.ServiceNameTransaction)
		if err != nil {
			return nil, nil, err
		}

		cl, err := grpc.NewInsecureClient(sn)
		if err != nil {
			log.Println("NewInsecureClient:", err)
			return nil, nil, err
		}

		return grpcBankTxn.NewTransactionServiceClient(cl.GetConn()), cl, nil
	}

	return nil, nil, nil
}

func (store transactionStore) ListAllInternal(params map[string]interface{}, businessID shared.BusinessID, accountID string) ([]BusinessPostedTransaction, error) {

	list := []BusinessPostedTransaction{}

	if os.Getenv("USE_TRANSACTION_SERVICE") == "true" {
		tc, cl, err := getTransactionServiceClient()
		if err != nil {
			return list, err
		}

		defer cl.CloseAndCancel()
		req := &grpcBankTxn.TransactionsRequest{}

		var bus *bsrv.Business
		if string(businessID) != "" {
			bus, err = bsrv.NewBusinessServiceWithout().GetByIdInternal(businessID)
			if err != nil {
				return list, err
			}

			busUUID, err := id.ParseUUID(string(bus.ID))
			if err != nil {
				return list, err
			}

			req.BusinessId = id.BusinessID(busUUID).String()
		}

		if accountID != "" {
			var bacID id.BankAccountID
			if strings.HasPrefix(accountID, id.IDPrefixBankAccount.String()) {
				bacID, err = id.ParseBankAccountID(accountID)
				if err != nil {
					return list, err
				}

			} else {
				accUUID, err := id.ParseUUID(accountID)
				if err != nil {
					return list, err
				}

				bacID = id.BankAccountID(accUUID)
			}
			req.AccountId = bacID.String()
		}

		contactID, ok := params["contactId"].(string)
		if ok {
			contactUUID, err := id.ParseUUID(contactID)
			if err != nil {
				return list, err
			}

			req.ContactId = id.ContactID(contactUUID).String()
		}

		req.DateRange = &grpcBankTxn.DateRange{
			Filter: grpcBankTxn.DateRangeFilter_DRF_UNSPECIFIED,
		}

		if params["startDate"] != "" && params["endDate"] != "" {
			req.DateRange.Filter = grpcBankTxn.DateRangeFilter_DRF_START_END

			startDate, err := grpc.ParseTimestampProto(params["startDate"].(string))
			if err != nil {
				startDate, err = grpc.ParseDateProto(params["startDate"].(string))
				if err != nil {
					return list, err
				}
			}

			req.DateRange.Start = startDate
			endDate, err := grpc.ParseTimestampProto(params["endDate"].(string))
			if err != nil {
				endDate, err = grpc.ParseDateProto(params["endDate"].(string))
				if err != nil {
					return list, err
				}
			}

			req.DateRange.End = endDate
		}

		req.AmountRange = &grpcBankTxn.AmountRange{
			Filter: grpcBankTxn.AmountRangeFilter_ARF_UNSPECIFIED,
		}

		if params["minAmount"] != nil && params["minAmount"] != "" && params["maxAmount"] != nil && params["maxAmount"] != "" {
			req.AmountRange.Filter = grpcBankTxn.AmountRangeFilter_ARF_MIN_MAX

			// Parse amounts to check high/low value (web app workaround)
			min, err := num.ParseDecimal(params["minAmount"].(string))
			if err != nil {
				return list, err
			}

			_ = min.V.Abs(min.V)

			max, err := num.ParseDecimal(params["maxAmount"].(string))
			if err != nil {
				return list, err
			}

			_ = max.V.Abs(max.V)

			// Swap entries based on higher/lower abs number
			if max.V.Cmp(min.V) < 0 {
				req.AmountRange.AmountMin = max.FormatCurrency()
				req.AmountRange.AmountMax = min.FormatCurrency()
			} else {
				req.AmountRange.AmountMin = min.FormatCurrency()
				req.AmountRange.AmountMax = max.FormatCurrency()
			}
		} else if params["minAmount"] != nil && params["minAmount"] != "" {
			req.AmountRange.Filter = grpcBankTxn.AmountRangeFilter_ARF_MIN
			req.AmountRange.AmountMin = params["minAmount"].(string)
		} else if params["minAmount"] != nil && params["maxAmount"] != "" {
			req.AmountRange.Filter = grpcBankTxn.AmountRangeFilter_ARF_MAX
			req.AmountRange.AmountMax = params["maxAmount"].(string)
		}

		if params["type"] != "" {
			codeType, ok := params["type"].(string)
			if ok {
				// Convert to range filter
				switch TransactionCodeType(codeType) {
				case TransactionCodeTypeDebitPosted:
					req.AmountRange.Type = grpcBankTxn.AmountRangeType_ART_DEBIT
				case TransactionCodeTypeCreditPosted:
					req.AmountRange.Type = grpcBankTxn.AmountRangeType_ART_CREDIT
				}
			}
		}

		req.StatusFilter = []grpcTxn.BankTransactionStatus{
			grpcTxn.BankTransactionStatus_BTS_CARD_POSTED,
			grpcTxn.BankTransactionStatus_BTS_NONCARD_POSTED,
		}

		req.Offset = int32(params["offset"].(int))
		req.Limit = int32(params["limit"].(int))

		req.SortRequests = []*grpcBankTxn.SortRequest{
			&grpcBankTxn.SortRequest{
				Name:      grpcBankTxn.SortFieldName_SFN_TRANSACTION_DATE,
				Direction: golang.SortDirection_SD_DESCENDING,
			},
		}

		if params["text"] != nil && params["text"] != "" {
			req.SearchTerms = strings.TrimSpace(params["text"].(string))
		}

		req.TypeFilter = []grpcTxn.BankTransactionType{}
		subtype, ok := params["subtype"].(string)
		if ok && params["subtype"] != "" {
			st := TransactionSubtype(subtype)
			if st == TransactionSubtypeACHTransferShopifyCredit {
				req.CounterpartyTypeFilter = []grpcTxn.BankTransactionCounterpartyType{grpcTxn.BankTransactionCounterpartyType_BTCT_SHOPPING_SHOPIFY_PAYOUT}
			} else if st == TransactionSubtypeACHTransferShopifyDebit {
				req.CounterpartyTypeFilter = []grpcTxn.BankTransactionCounterpartyType{grpcTxn.BankTransactionCounterpartyType_BTCT_SHOPPING_SHOPIFY_REFUND}
			} else {
				typeProto, ok := TransactionSubtypeToTypeProto[st]
				stp := []grpcTxn.BankTransactionType{
					typeProto,
				}
				if !ok {
					return nil, fmt.Errorf("Invalid subtype")
				}
				req.TypeFilter = stp
			}
		}
		resp, err := tc.GetMany(context.Background(), req)
		if err != nil {
			return list, err
		}

		for _, gtxn := range resp.Results {
			pt, err := BusinessPostedTransactionFromProto(gtxn, bus)
			if err != nil {
				return list, err
			}

			list = append(list, *pt)
		}
	} else {
		var dateFilter string = ""
		if params["startDate"] != "" && params["endDate"] != "" {
			dateFilter = " AND transaction_date >=  '" + params["startDate"].(string) + "' AND transaction_date <= '" + params["endDate"].(string) + "'"
		}

		var txnTypeFilter string = ""
		if params["type"] != "" {
			txnTypeFilter = " AND code_type = '" + params["type"].(string) + "'"
		}

		var contactFilter string = ""
		if params["contactId"] != "" {
			contactFilter = " AND contact_id = '" + params["contactId"].(string) + "'"
		}

		var amtFilter string = ""
		if params["minAmount"] != "" {
			amtFilter = " AND ABS(amount) >= " + params["minAmount"].(string)
		}

		if params["maxAmount"] != "" {
			amtFilter = amtFilter + " AND ABS(amount) <= " + params["maxAmount"].(string)
		}

		var txtFilter string = ""
		if params["text"] != "" {
			txtFilter = " AND (bank_transaction_desc ILIKE  '%" + params["text"].(string) + "%' OR transaction_desc ILIKE '%" + params["text"].(string) + "%')"
		}

		columns := `
			id, business_id, bank_name, bank_transaction_id, bank_extra, transaction_type, account_id,
			card_id, code_type, amount, currency, money_transfer_id, contact_id, bank_transaction_desc,
			bank_transaction_desc AS money_transfer_desc, transaction_desc, transaction_date,
			source_notes, source_notes AS notes, transaction_title, transaction_subtype, created`

		var err error
		if businessID != "" && accountID != "" {
			query := `SELECT ` + columns + ` FROM business_transaction WHERE business_id = $1 AND account_id = $2` + dateFilter + txnTypeFilter + contactFilter + amtFilter + txtFilter +
				` ORDER BY transaction_date DESC LIMIT $3 OFFSET $4`

			err = store.Select(&list, query, businessID, accountID, params["limit"].(int), params["offset"].(int))
		} else if businessID != "" {
			query := `SELECT ` + columns + ` FROM business_transaction WHERE business_id = $1` + dateFilter + txnTypeFilter + contactFilter + amtFilter + txtFilter +
				` ORDER BY transaction_date DESC LIMIT $2 OFFSET $3`

			err = store.Select(&list, query, businessID, params["limit"].(int), params["offset"].(int))
		} else if accountID != "" {
			query := `SELECT ` + columns + ` FROM business_transaction WHERE account_id = $1` + dateFilter + txnTypeFilter + contactFilter + amtFilter + txtFilter +
				` ORDER BY transaction_date DESC LIMIT $2 OFFSET $3`

			err = store.Select(&list, query, accountID, params["limit"].(int), params["offset"].(int))
		} else {
			var query string
			if dateFilter != "" || txnTypeFilter != "" || contactFilter != "" || amtFilter != "" || txtFilter != "" {
				query = `SELECT ` + columns + ` FROM business_transaction WHERE 1=1` + dateFilter + txnTypeFilter + contactFilter + amtFilter + txtFilter +
					` ORDER BY transaction_date DESC LIMIT $1 OFFSET $2`
			} else {
				query = `SELECT ` + columns + ` FROM business_transaction ORDER BY transaction_date DESC LIMIT $1 OFFSET $2`
			}

			err = store.Select(&list, query, params["limit"].(int), params["offset"].(int))
		}

		if err != nil && err == sql.ErrNoRows {
			log.Println(err)
			return list, nil
		}

		if err != nil {
			log.Println(err)
			return nil, err
		}
	}

	if os.Getenv("USE_INVOICE_SERVICE") == "true" {
		list = store.updateWalletType(list)
	}
	return list, nil
}

// updates the walletType in BusinessPostedTransaction.Source.WalletType if moneyRequestID is not NULL
func (store transactionStore) updateWalletType(transactions []BusinessPostedTransaction) []BusinessPostedTransaction {
	// collect all moneyRequestID
	requestIDs := []shared.PaymentRequestID{}
	for _, txn := range transactions {
		if txn.MoneyRequestID != nil {
			requestIDs = append(requestIDs, *txn.MoneyRequestID)
		}
	}

	// If no requestIDs found, then nothing to update, return
	if len(requestIDs) == 0 {
		return transactions
	}

	// Get the payments info from request IDs
	sourceReq := services.NewSourceRequest()
	payments, err := payment.NewPaymentService(sourceReq).GetPayments(requestIDs)
	if err != nil {
		log.Println(err.Error())
		return transactions
	}

	// no payment info found for those transactions
	if len(payments) == 0 {
		return transactions
	}

	// Make map of invoice_id and walletType { invoiceID:walletType }
	// we are using key as string, to simplify mapping of id.InvoiceID to shared.PaymentRequestID
	walletMap := map[string]string{}
	for _, pm := range payments {
		if pm.WalletType != nil {
			walletMap[pm.InvoiceID.UUIDString()] = *pm.WalletType
		}
	}

	// update the walletType in transactions
	for idx, txn := range transactions {
		if txn.MoneyRequestID != nil {
			if val, ok := walletMap[txn.MoneyRequestID.ToUUIDString()]; ok {
				if txn.Source != nil {
					txn.Source.WalletType = &val
				} else {
					transactions[idx].Source = &TransactionSource{
						WalletType: &val,
					}
				}
			}
		}
	}

	return transactions
}

func (store transactionStore) ListAll(params map[string]interface{}, userID shared.UserID, businessID shared.BusinessID, accountID string) ([]BusinessPostedTransaction, error) {
	sourceReq := services.NewSourceRequest()
	sourceReq.UserID = userID
	err := auth.NewAuthService(sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	resp := []BusinessPostedTransaction{}

	if string(businessID) == "" {
		return resp, errors.New("business id required")
	}

	return store.ListAllInternal(params, businessID, accountID)
}

// Fetches all transactions without offset and limit restriction
func (store transactionStore) listAllForExport(businessID *shared.BusinessID, startDate, endDate string, offset, limit int) ([]BusinessPostedTransaction, error) {
	list := []BusinessPostedTransaction{}

	if os.Getenv("USE_TRANSACTION_SERVICE") == "true" {
		tc, cl, err := getTransactionServiceClient()
		if err != nil {
			return list, err
		}

		defer cl.CloseAndCancel()
		req := &grpcBankTxn.TransactionsRequest{}

		var bus *bsrv.Business
		if businessID != nil {
			busUUID, err := id.ParseUUID(string(*businessID))
			if err != nil {
				return list, err
			}

			busID := id.BusinessID(busUUID)
			if busID.IsZero() {
				return list, errors.New("invalid business id")
			}

			req.BusinessId = busID.String()

			bus, err = bsrv.NewBusinessServiceWithout().GetByIdInternal(*businessID)
			if err != nil {
				return list, err
			}
		}

		start, err := grpc.ParseTimestampProto(startDate)
		if err != nil {
			start, err = grpc.ParseDateProto(startDate)
			if err != nil {
				return list, err
			}
		}

		end, err := grpc.ParseTimestampProto(endDate)
		if err != nil {
			end, err = grpc.ParseDateProto(endDate)
			if err != nil {
				return list, err
			}
		}

		req.DateRange = &grpcBankTxn.DateRange{
			Filter: grpcBankTxn.DateRangeFilter_DRF_START_END,
			Start:  start,
			End:    end,
		}

		req.StatusFilter = []grpcTxn.BankTransactionStatus{
			grpcTxn.BankTransactionStatus_BTS_CARD_POSTED,
			grpcTxn.BankTransactionStatus_BTS_NONCARD_POSTED,
		}

		req.Offset = int32(offset)
		req.Limit = int32(limit)

		req.SortRequests = []*grpcBankTxn.SortRequest{
			&grpcBankTxn.SortRequest{
				Name:      grpcBankTxn.SortFieldName_SFN_TRANSACTION_DATE,
				Direction: golang.SortDirection_SD_DESCENDING,
			},
		}

		resp, err := tc.GetMany(context.Background(), req)
		if err != nil {
			return list, err
		}

		for _, gtxn := range resp.Results {
			// TODO: Improve performance with full txn fetch in service
			fullReq := &grpcBankTxn.TransactionIDRequest{
				Id: gtxn.Id,
			}

			gtxnFull, err := tc.GetFullByID(context.Background(), fullReq)
			if err != nil {
				return list, err
			}

			t, err := BusinessPostedTransactionFromFullProto(gtxnFull, bus)
			if err != nil {
				return list, err
			}

			list = append(list, *t)
		}
	} else {
		var filterText string
		if businessID != nil {
			filterText = "business_id = '" + string(*businessID) + "' AND "
		}

		filterText = filterText + " transaction_date >=  '" + startDate + "' AND transaction_date <= '" + endDate + "'"

		query := `SELECT * FROM business_transaction WHERE ` + filterText + ` ORDER BY transaction_date DESC OFFSET $1 LIMIT $2`

		err := store.Select(&list, query, offset, limit)
		if err != nil && err == sql.ErrNoRows {
			log.Println(err)
			return list, nil
		}

		if err != nil {
			log.Println(err)
			return nil, err
		}
	}

	return list, nil
}

func (store transactionStore) Create(t BusinessPostedTransactionCreate) (*BusinessPostedTransaction, error) {
	var trx BusinessPostedTransaction

	if t.AccountID != nil && strings.HasPrefix(*t.AccountID, string(id.IDPrefixBankAccount)) {
		aID, err := id.ParseBankAccountID(*t.AccountID)
		if err != nil {
			return &trx, err
		}

		s := aID.UUIDString()
		t.AccountID = &s
	}

	keys := shared.SQLGenInsertKeys(t)
	values := shared.SQLGenInsertValues(t)
	query := fmt.Sprintf("INSERT INTO business_transaction(%s) VALUES(%s) RETURNING *", keys, values)

	stmt, err := store.PrepareNamed(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = stmt.Get(&trx, t)
	if err != nil {
		log.Println(err)
	}

	return &trx, err
}

func (store transactionStore) GetByIDInternal(ptID shared.PostedTransactionID) (*BusinessPostedTransaction, error) {
	t := &BusinessPostedTransaction{}

	if os.Getenv("USE_TRANSACTION_SERVICE") == "true" {
		tc, cl, err := getTransactionServiceClient()
		if err != nil {
			return t, err
		}

		defer cl.CloseAndCancel()
		txnUUID, err := id.ParseUUID(string(ptID))
		if err != nil {
			return t, err
		}

		req := &grpcBankTxn.TransactionIDRequest{
			Id: id.BankTransactionID(txnUUID).String(),
		}

		resp, err := tc.GetFullByID(context.Background(), req)
		if err != nil {
			return t, err
		}

		busID, err := id.ParseBusinessID(resp.Transaction.BusinessId)
		if err != nil {
			return t, err
		}

		bus, err := bsrv.NewBusinessServiceWithout().GetByIdInternal(shared.BusinessID(busID.UUIDString()))
		if err != nil {
			return t, err
		}

		t, err = BusinessPostedTransactionFromFullProto(resp, bus)
		if err != nil {
			log.Println(err)
			return t, err
		}
	} else {
		query := `
        SELECT
            business_transaction.*,
            business_transaction.bank_transaction_desc AS money_transfer_desc,
            business_transaction.source_notes AS notes,
            business_card_transaction.id "business_card_transaction.id",
            business_card_transaction.auth_amount "business_card_transaction.auth_amount",
            business_card_transaction.transaction_type "business_card_transaction.transaction_type",
            business_card_transaction.local_amount "business_card_transaction.local_amount",
            business_card_transaction.local_currency "business_card_transaction.local_currency",
            business_card_transaction.local_date "business_card_transaction.local_date",
            business_card_transaction.billing_currency "business_card_transaction.billing_currency",
			business_card_transaction.merchant_category_code "business_card_transaction.merchant_category_code",
            business_card_transaction.merchant_name "business_card_transaction.merchant_name",
            business_card_transaction.merchant_street_address "business_card_transaction.merchant_street_address",
            business_card_transaction.merchant_city "business_card_transaction.merchant_city",
            business_card_transaction.merchant_state "business_card_transaction.merchant_state",
            business_card_transaction.merchant_country "business_card_transaction.merchant_country",
            business_hold_transaction.id "business_hold_transaction.id",
            business_hold_transaction.amount "business_hold_transaction.amount",
            business_hold_transaction.hold_number "business_hold_transaction.hold_number",
            business_hold_transaction.transaction_date "business_hold_transaction.transaction_date",
            business_hold_transaction.expiry_date "business_hold_transaction.expiry_date"
        FROM business_transaction
        LEFT JOIN business_card_transaction ON business_transaction.id = business_card_transaction.transaction_id
        LEFT JOIN business_hold_transaction ON business_transaction.id = business_hold_transaction.transaction_id
        WHERE business_transaction.id = $1`

		err := store.Get(t, query, ptID)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}

	// Transaction Dispute
	query := `
		SELECT
			id "business_transaction_dispute.id",
            receipt_id "business_transaction_dispute.receipt_id",
            category "business_transaction_dispute.category",
            summary "business_transaction_dispute.summary",
			dispute_status "business_transaction_dispute.dispute_status",
			dispute_number "business_transaction_dispute.dispute_number",
			created "business_transaction_dispute.created"
		FROM business_transaction_dispute
		WHERE transaction_id = $1`

	err := store.Get(t, query, ptID)
	if err == sql.ErrNoRows {
		t.Dispute = nil
	} else if err != nil {
		log.Println("Dispute:", err)
	}

	// Transaction Notes
	query = `
        SELECT transaction_notes "business_transaction_annotation.transaction_notes"
        FROM business_transaction_annotation
        WHERE transaction_id = $1`

	err = store.Get(t, query, ptID)
	if err == sql.ErrNoRows {
		t.TransactionNotes = nil
	} else if err != nil {
		log.Println("Notes:", err)
	}

	// Transaction Attachment
	query = `
	    SELECT
    	    id "business_transaction_attachment.id",
        	deleted "business_transaction_attachment.deleted"
	    FROM business_transaction_attachment
    	WHERE transaction_id = $1 AND deleted IS NULL`

	err = store.Get(t, query, ptID)
	if err == sql.ErrNoRows {
		t.AttachmentID = nil
		t.AttachmentDeleted = nil
	} else if err != nil {
		log.Println("Attachment:", err)
	}

	// Money Transfer
	if t.MoneyTransferID != nil {
		if os.Getenv("USE_BANKING_SERVICE") == "true" {
			bts, err := busBanking.NewBankingTransferService()
			if err != nil {
				log.Println("MoneyTransfer:", err)
				return nil, err
			}

			mt, err := bts.GetByIDOnlyInternal(*t.MoneyTransferID)
			if err != nil {
				log.Println("MoneyTransfer:", err)
			}

			t.MoneyTransfer = new(MoneyTransfer)

			if mt != nil {
				t.MoneyTransfer.SourceAccountID = mt.SourceAccountId
				t.MoneyTransfer.SourceType = mt.SourceType
				t.MoneyTransfer.DestAccountID = mt.DestAccountId
				t.MoneyTransfer.DestType = mt.DestType
				t.MoneyTransfer.Amount = &mt.Amount

				currency := Currency(mt.Currency)
				t.MoneyTransfer.Currency = &currency

				t.MoneyTransfer.Notes = mt.Notes
				t.SourceNotes = mt.Notes
				t.MoneyTransfer.Status = &mt.Status
				t.MoneyTransfer.Created = &mt.Created
			}
		} else {
			query := `
	        SELECT
			    source_account_id "business_money_transfer.source_account_id",
				source_type "business_money_transfer.source_type",
				dest_account_id "business_money_transfer.dest_account_id",
				dest_type "business_money_transfer.dest_type",
	    	    amount "business_money_transfer.amount",
    	    	currency "business_money_transfer.currency",
        		notes "business_money_transfer.notes",
				notes "source_notes",
        		status "business_money_transfer.status",
        		created "business_money_transfer.created"
	        FROM business_money_transfer
			WHERE id = $1`

			err := data.DBRead.Get(t, query, *t.MoneyTransferID)
			if err != nil {
				log.Println("MoneyTransfer:", err)
			}
		}
	}

	// Contact
	if t.ContactID != nil {
		query := `
            SELECT
                id "business_contact.id",
                contact_category "business_contact.contact_category",
                contact_type "business_contact.contact_type",
                engagement "business_contact.engagement",
                job_title "business_contact.job_title",
                business_name "business_contact.business_name",
                first_name "business_contact.first_name",
                last_name "business_contact.last_name",
                phone_number "business_contact.phone_number",
                email "business_contact.email",
                mailing_address "business_contact.mailing_address"
            FROM business_contact
            WHERE id = $1`

		err := data.DBRead.Get(t, query, *t.ContactID)
		if err != nil {
			log.Println("Contact:", err)
		}
	}

	switch t.CodeType {
	case TransactionCodeTypeCreditPosted:
		if t.MoneyTransfer != nil {
			if os.Getenv("USE_BANKING_SERVICE") == "true" {
				blas, err := busBanking.NewBankingLinkedAccountService()
				if err != nil {
					return nil, err
				}

				la, err := blas.GetById(t.MoneyTransfer.SourceAccountID)
				if err != nil {
					log.Println("LinkedAccount:", err)
				}

				if la != nil {
					t.Contact = new(Contact)

					accNumber := contact.AccountNumber(la.AccountNumber)
					t.Contact.AccountNumber = &accNumber

					t.Contact.RoutingNumber = &la.RoutingNumber
				}
			} else {
				query := `
	        	SELECT
    	    		account_number "business_contact.account_number",
        			routing_number "business_contact.routing_number"
				FROM business_linked_bank_account
				WHERE id = $1`

				err := data.DBRead.Get(t, query, t.MoneyTransfer.SourceAccountID)
				if err != nil {
					log.Println("LinkedAccount:", err)
				}
			}
		}
	case TransactionCodeTypeDebitPosted:
		if t.MoneyTransfer != nil {
			switch t.MoneyTransfer.DestType {
			case banking.TransferTypeAccount:
				if os.Getenv("USE_BANKING_SERVICE") == "true" {
					blas, err := busBanking.NewBankingLinkedAccountService()
					if err != nil {
						return nil, err
					}

					la, err := blas.GetById(t.MoneyTransfer.DestAccountID)
					if err != nil {
						log.Println("LinkedAccount:", err)
					}

					if la != nil {
						t.Contact = new(Contact)

						accNumber := contact.AccountNumber(la.AccountNumber)
						t.Contact.AccountNumber = &accNumber

						t.Contact.RoutingNumber = &la.RoutingNumber
					}
				} else {
					query := `
					SELECT
		    		    account_number "business_contact.account_number",
				        routing_number "business_contact.routing_number",
        				bank_name "business_contact.bank_name"
					FROM business_linked_bank_account
					WHERE id = $1`

					err := data.DBRead.Get(t, query, t.MoneyTransfer.DestAccountID)
					if err != nil {
						log.Println("LinkedAccount:", err)
					}
				}
			case banking.TransferTypeCard:
				if os.Getenv("USE_BANKING_SERVICE") == "true" {
					blcs, err := busBanking.NewBankingLinkedCardService()
					if err != nil {
						return nil, err
					}

					lc, err := blcs.GetByID(t.MoneyTransfer.DestAccountID)
					if err != nil {
						log.Println("LinkedCard:", err)
					}

					if lc != nil {
						t.Contact = new(Contact)

						t.Contact.CardNumber = &lc.CardNumberMasked
						t.Contact.CardBrand = &lc.CardBrand
					}
				} else {
					query := `
					SELECT
						business_linked_card.card_number_masked "business_contact.card_number_masked",
						business_linked_card.card_brand "business_contact.card_brand"
					FROM business_linked_card
					WHERE id = $1`

					err := data.DBRead.Get(t, query, t.MoneyTransfer.DestAccountID)
					if err != nil {
						log.Println("LinkedCard:", err)
					}
				}
			}
		}
	}

	if t.MoneyRequestID != nil {
		// Check request type
		isPOSRequest := false
		tmpPostedTxn := &BusinessPostedTransaction{}
		posCheckQuery := `SELECT business_money_request.amount "business_money_request.amount" FROM business_money_request WHERE business_money_request.id = $1 and business_money_request.request_type = 'pos';`
		err = data.DBRead.Get(tmpPostedTxn, posCheckQuery, *t.MoneyRequestID)
		if err != nil && err == sql.ErrNoRows {
			isPOSRequest = false
		} else if tmpPostedTxn.MoneyRequest != nil {
			isPOSRequest = true
		}

		if os.Getenv("USE_INVOICE_SERVICE") == "true" && !isPOSRequest {
			invSvc, err := invoice.NewInvoiceService()
			if err != nil {
				log.Println("Error fetching money request invoice and receipt from invoice service", err)
			} else {
				resp, err := invSvc.GetInvoiceIDFromPaymentRequestID(t.MoneyRequestID)
				if err != nil {
					log.Println(err)
				} else {
					query := `
	        	SELECT
		        business_receipt.id "business_receipt.id",
				business_receipt.receipt_number "business_receipt.receipt_number",
        		business_money_request_payment.id "business_money_request_payment.id",
				business_money_request_payment.receipt_id "business_money_request_payment.receipt_id",
		        business_money_request_payment.payment_date "business_money_request_payment.payment_date",
				business_money_request_payment.purchase_address "business_money_request_payment.purchase_address",
        		business_money_request_payment.card_brand "business_money_request_payment.card_brand",
				business_money_request_payment.card_number "business_money_request_payment.card_number",
				business_money_request_payment.fee_amount "business_money_request_payment.fee_amount",
				business_money_request_payment.wallet_type "business_money_request_payment.wallet_type"
			FROM business_money_request_payment
			LEFT JOIN business_receipt ON business_money_request_payment.invoice_id = business_receipt.invoice_id_v2
    	    WHERE business_money_request_payment.invoice_id = $1`

					err = data.DBRead.Get(t, query, *t.MoneyRequestID)
					amount, ok := resp.Amount.Float64()
					if !ok {
						log.Println("unable to convert to float")
					} else {
						t.MoneyRequest = &payment.RequestMini{
							Amount: amount,
						}
						invoiceID := resp.InvoiceID.UUIDString()
						invNumber := fmt.Sprintf("%05d", resp.Number)
						t.Invoice = &payment.InvoiceMini{
							InvoiceNumber: &invNumber,
							Id:            &invoiceID,
							ViewLink:      &resp.InvoiceViewLink,
						}
						if t.Receipt != nil {
							encodedInvoiceID := base64.RawURLEncoding.EncodeToString([]byte(resp.InvoiceID.String()))
							viewLink := fmt.Sprintf("%s/invoice-receipt?token=%s",
								os.Getenv("PAYMENTS_URL"), encodedInvoiceID)
							t.Receipt.ViewLink = &viewLink
						}

					}
					log.Println(err)
				}
			}

			if err != nil {
				log.Println("MoneyRequest (Invoice conversion error):", err)
			}
		} else {
			query := `
	        SELECT
    		    business_invoice.id "business_invoice.id",
				business_invoice.invoice_number "business_invoice.invoice_number",
		        business_receipt.id "business_receipt.id",
				business_receipt.receipt_number "business_receipt.receipt_number",
        		business_money_request_payment.id "business_money_request_payment.id",
				business_money_request_payment.receipt_id "business_money_request_payment.receipt_id",
		        business_money_request_payment.payment_date "business_money_request_payment.payment_date",
				business_money_request_payment.purchase_address "business_money_request_payment.purchase_address",
        		business_money_request_payment.card_brand "business_money_request_payment.card_brand",
				business_money_request_payment.card_number "business_money_request_payment.card_number",
				business_money_request_payment.fee_amount "business_money_request_payment.fee_amount",
				business_money_request_payment.wallet_type "business_money_request_payment.wallet_type",
				business_money_request.id "business_money_request.id",
				business_money_request.amount "business_money_request.amount"
			FROM business_money_request
			LEFT JOIN business_receipt ON business_money_request.id = business_receipt.request_id
    	    LEFT JOIN business_invoice ON business_money_request.id = business_invoice.request_id
			LEFT JOIN business_money_request_payment ON business_money_request.id = business_money_request_payment.request_id
    	    WHERE business_money_request.id = $1 `

			err := data.DBRead.Get(t, query, *t.MoneyRequestID)
			if err != nil {
				log.Println("Error fetching money request invoice and receipt", err)
			}
		}
	}

	if t.MoneyTransferDesc != nil {
		t.OriginAccount = GetOriginAccount(TransactionTypeTransfer, shared.StringValue(t.MoneyTransferDesc))
		t.DestinationAccount = GetDestinationAccount(TransactionTypeTransfer, shared.StringValue(t.MoneyTransferDesc))
	}

	switch t.TransactionType {
	case TransactionTypeTransfer:
		handleTransferTransaction(t)
	case TransactionTypeACH:
		handleACHTransaction(t)
	case TransactionTypeDeposit:
		handleDepositTransaction(t)
	case TransactionTypeRefund:
		if t.TransactionSubtype != nil && *t.TransactionSubtype == TransactionSubtypeCardPushCredit {
			handlePushToCreditTransaction(t)
		} else {
			handleCardTransaction(t)
		}
	case TransactionTypePurchase:
		if t.MoneyTransferID != nil {
			handlePushToDebitTransaction(t)
		} else {
			handleCardTransaction(t)
		}
	case TransactionTypeATM:
		handleCardTransaction(t)
	}

	return t, nil
}

func (store transactionStore) GetByID(id shared.PostedTransactionID, userID shared.UserID, businessID shared.BusinessID) (*BusinessPostedTransaction, error) {
	sourceReq := services.NewSourceRequest()
	sourceReq.UserID = userID
	err := auth.NewAuthService(sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	t, err := store.GetByIDInternal(id)
	if err != nil {
		return t, err
	}

	if t.BusinessID != businessID {
		return t, errors.New("unauthorized")
	}

	return t, nil
}

func handlePushToDebitTransaction(t *BusinessPostedTransaction) {
	switch t.CodeType {
	case TransactionCodeTypeDebitPosted:
		d := TransactionDestination{}

		if t.Contact != nil {
			if t.Contact.BusinessName != nil && len(*t.Contact.BusinessName) > 0 {
				d.AccountHolderName = t.Contact.BusinessName
			} else if t.Contact.FirstName != nil && t.Contact.LastName != nil {
				name := *t.Contact.FirstName + " " + *t.Contact.LastName
				d.AccountHolderName = &name
			}

			if t.Contact.CardNumber != nil && len(*t.Contact.CardNumber) > 0 {
				cardNumber := services.MaskLeft(string(*t.Contact.CardNumber), 4)
				cardNumber = string(cardNumber[len(cardNumber)-4:])
				d.AccountNumber = &cardNumber
			}

			if t.Contact.CardBrand != nil && len(*t.Contact.CardBrand) > 0 {
				cardBrand := strings.Title(strings.ToLower(*t.Contact.CardBrand))
				d.BankName = &cardBrand
			}

			t.Destination = &d
			return
		}

		d.AccountNumber = GetDestinationAccount(TransactionTypeTransfer, shared.StringValue(t.BankTransactionDesc))
		d.AccountHolderName = GetDestinationAccountHolder(TransactionTypeTransfer, shared.StringValue(t.BankTransactionDesc))
		t.Destination = &d
	default:
		log.Println("Unhandled transaction code type", t.CodeType)

	}
}

func handleTransferTransaction(t *BusinessPostedTransaction) {
	switch t.CodeType {
	case TransactionCodeTypeCreditPosted:
		s := TransactionSource{}

		if t.MoneyRequestID != nil {
			if t.Payment != nil {
				s.CardBrand = t.Payment.CardBrand
				s.CardLast4 = t.Payment.CardLast4
				s.WalletType = t.Payment.WalletType
				s.PurchaseAddress = t.Payment.PurchaseAddress
			}

			if t.Contact != nil {
				if t.Contact.BusinessName != nil && len(*t.Contact.BusinessName) > 0 {
					s.CardHolderName = t.Contact.BusinessName
				} else if t.Contact.FirstName != nil && t.Contact.LastName != nil {
					name := *t.Contact.FirstName + " " + *t.Contact.LastName
					s.CardHolderName = &name
				}
			}

			t.Source = &s

			// Attach payment summary
			ps := PaymentSummary{
				InvoiceAmount: t.MoneyRequest.Amount,
				FeeAmount:     t.Payment.FeeAmount,
				NetAmount:     t.Amount,
			}
			t.PaymentSummary = &ps

			return
		}

		s.AccountNumber = GetOriginAccount(TransactionTypeTransfer, shared.StringValue(t.BankTransactionDesc))
		s.AccountHolderName = GetOriginAccountHolder(TransactionTypeTransfer, shared.StringValue(t.BankTransactionDesc))
		t.Source = &s

		return

	case TransactionCodeTypeDebitPosted:
		d := TransactionDestination{}

		if t.Contact != nil {
			if t.Contact.BusinessName != nil && len(*t.Contact.BusinessName) > 0 {
				d.AccountHolderName = t.Contact.BusinessName
			} else if t.Contact.FirstName != nil && t.Contact.LastName != nil {
				name := *t.Contact.FirstName + " " + *t.Contact.LastName
				d.AccountHolderName = &name
			}

			if t.Contact.AccountNumber != nil && len(*t.Contact.AccountNumber) > 0 {
				accNumber := services.MaskLeft(string(*t.Contact.AccountNumber), 4)
				accNumber = string(accNumber[len(accNumber)-4:])
				d.AccountNumber = &accNumber
			}

			t.Destination = &d
			return
		}

		d.AccountNumber = GetDestinationAccount(TransactionTypeTransfer, shared.StringValue(t.BankTransactionDesc))
		d.AccountHolderName = GetDestinationAccountHolder(TransactionTypeTransfer, shared.StringValue(t.BankTransactionDesc))
		t.Destination = &d

	}
}

func handleACHTransaction(t *BusinessPostedTransaction) {
	switch t.CodeType {
	case TransactionCodeTypeCreditPosted:
		s := TransactionSource{}

		if t.MoneyRequestID != nil {
			if t.Payment != nil {
				s.BankName = t.Payment.CardBrand
				s.AccountNumber = t.Payment.CardLast4
				s.WalletType = t.Payment.WalletType
			}

			if t.Contact != nil {
				if t.Contact.BusinessName != nil && len(*t.Contact.BusinessName) > 0 {
					s.AccountHolderName = t.Contact.BusinessName
				} else if t.Contact.FirstName != nil && t.Contact.LastName != nil {
					name := *t.Contact.FirstName + " " + *t.Contact.LastName
					s.AccountHolderName = &name
				}
			}

			t.Source = &s
			return
		}

		s.AccountNumber = GetOriginAccount(TransactionTypeACH, shared.StringValue(t.BankTransactionDesc))
		s.AccountHolderName = GetOriginAccountHolder(TransactionTypeACH, shared.StringValue(t.BankTransactionDesc))
		t.Source = &s

		if t.TransactionSubtype != nil && *t.TransactionSubtype == TransactionSubtypeACHTransferShopifyCredit {
			getShopifyPayoutDetails(t)
		}

		return
	case TransactionCodeTypeDebitPosted:
		d := TransactionDestination{}

		if t.Contact != nil && t.Contact.Id != nil && len(*t.Contact.Id) > 0 {
			if t.Contact.BusinessName != nil && len(*t.Contact.BusinessName) > 0 {
				d.AccountHolderName = t.Contact.BusinessName
			} else if t.Contact.FirstName != nil && t.Contact.LastName != nil {
				name := *t.Contact.FirstName + " " + *t.Contact.LastName
				d.AccountHolderName = &name
			}

			if t.Contact.AccountNumber != nil && len(*t.Contact.AccountNumber) > 0 {
				accNumber := services.MaskLeft(string(*t.Contact.AccountNumber), 4)
				accNumber = string(accNumber[len(accNumber)-4:])
				d.AccountNumber = &accNumber
			}

			if t.Contact.BankName != nil && len(*t.Contact.BankName) > 0 {
				d.BankName = t.Contact.BankName
			}
		}

		if d.AccountNumber == nil {
			d.AccountNumber = GetDestinationAccount(TransactionTypeACH, shared.StringValue(t.BankTransactionDesc))
		}

		if d.AccountHolderName == nil {
			if t.MoneyTransferID != nil {
				d.AccountHolderName = GetDestinationAccountHolder(TransactionTypeACH, shared.StringValue(t.BankTransactionDesc))
			} else {
				destName := GetExternalACHDestination(shared.StringValue(t.BankTransactionDesc))
				if len(destName) > 0 {
					d.AccountHolderName = &destName
				}
			}
		}

		t.Destination = &d
	}
}

func handleDepositTransaction(t *BusinessPostedTransaction) {
	switch t.CodeType {
	case TransactionCodeTypeCreditPosted:
		s := TransactionSource{}

		s.AccountHolderName = GetOriginAccountHolder(TransactionTypeDeposit, shared.StringValue(t.BankTransactionDesc))
		t.Source = &s

		return
	}
}

func handleCardTransaction(t *BusinessPostedTransaction) {
	switch t.CodeType {
	case TransactionCodeTypeDebitPosted:
		if t.CardTransaction != nil {
			s := TransactionSource{}
			d := TransactionDestination{}
			a := services.Address{}
			// purchase location
			if t.CardTransaction.MerchantStreetAddress != nil && isLetter(*t.CardTransaction.MerchantStreetAddress) {
				a.StreetAddress = strings.TrimSpace(*t.CardTransaction.MerchantStreetAddress)
			}

			if t.CardTransaction.MerchantCity != nil && isLetter(*t.CardTransaction.MerchantCity) {
				a.City = strings.TrimSpace(*t.CardTransaction.MerchantCity)
			}

			if t.CardTransaction.MerchantState != nil && isLetter(*t.CardTransaction.MerchantState) {
				a.State = strings.TrimSpace(*t.CardTransaction.MerchantState)
			}

			if t.CardTransaction.MerchantCountry != nil && isLetter(*t.CardTransaction.MerchantCountry) {
				a.Country = strings.TrimSpace(*t.CardTransaction.MerchantCountry)
			}

			s.PurchaseAddress = &a
			t.Source = &s

			if t.CardTransaction.MerchantName != nil && len(*t.CardTransaction.MerchantName) > 0 {
				name := strings.TrimSpace(*t.CardTransaction.MerchantName)
				d.AccountHolderName = &name
			}

			t.Destination = &d
		}
	case TransactionCodeTypeCreditPosted:
		if t.CardTransaction != nil {

			if t.CardTransaction.MerchantName != nil && len(*t.CardTransaction.MerchantName) > 0 {
				s := TransactionSource{}

				name := strings.TrimSpace(*t.CardTransaction.MerchantName)
				s.AccountHolderName = &name

				t.Source = &s
			}
		}
	}
}

func handlePushToCreditTransaction(t *BusinessPostedTransaction) {
	s := TransactionSource{}

	name := GetInstantPaySenderName(t.BankTransactionDesc)
	s.AccountHolderName = &name

	t.Source = &s
}

func isLetter(s string) bool {
	a := strings.Split(strings.TrimSpace(s), " ")

	for _, s := range a {
		for _, r := range s {
			if !unicode.IsLetter(r) {
				return false
			}
		}
	}

	return true
}

func (store transactionStore) ExportInternal(params map[string]interface{}) (*CSVTransaction, error) {

	var b bytes.Buffer
	w := csv.NewWriter(&b)

	header := []string{
		"S No.", "Transaction Date", "Transaction Remarks", "Transaction Type", "Amount(USD)", "Authorization Amount",
		"Merchant Name", "Merchant Street Address", "Merchant City", "Merchant State", "Merchant Country", "Notes",
	}

	err := w.Write(header)
	if err != nil {
		log.Fatalln("error writing record to csv:", err)
		return nil, err
	}

	offset := 0
	limit := 20
	params["limit"] = limit

	for {
		// Using 'ListAllInternal()' instead of 'listAllForExport()'
		// Fetch all transactions between given dates
		// t, err := store.listAllForExport(businessID, startDate, endDate, offset, limit)

		var busID shared.BusinessID
		val, ok := params["businessId"].(shared.BusinessID)
		if ok {
			busID = val
		}
		params["offset"] = offset

		t, err := store.ListAllInternal(params, busID, "")

		if err != nil {
			return nil, err
		}

		if len(t) == 0 {
			break
		}

		for i, txn := range t {

			txnDate := txn.TransactionDate.Format("01/02/06")
			txnType := txn.CodeType

			if txnType == "debitPosted" {
				txnType = "debit"
			} else {
				txnType = "credit"
			}

			txnAmount := txn.Amount.FormatCurrency()

			var authAmt, merchantName, merchantStAddress, merchantCity, merchantState, merchantCountry string
			if txn.CardTransaction != nil {
				amtVal, _ := txn.CardTransaction.AuthAmount.Float64()
				authAmt = shared.FormatFloatAmount(amtVal)
				merchantName = shared.StringValue(txn.CardTransaction.MerchantName)
				merchantStAddress = shared.StringValue(txn.CardTransaction.MerchantStreetAddress)
				merchantCity = shared.StringValue(txn.CardTransaction.MerchantCity)
				merchantState = shared.StringValue(txn.CardTransaction.MerchantState)
				merchantCountry = shared.StringValue(txn.CardTransaction.MerchantCountry)
			}

			var txnDesc string
			if len(txn.TransactionDesc) > 0 {
				txnDesc = txn.TransactionDesc
			} else {
				txnDesc = shared.StringValue(txn.BankTransactionDesc)
			}

			data := []string{
				strconv.Itoa(i + 1),
				txnDate,
				txnDesc,
				string(txnType),
				txnAmount,
				authAmt,
				merchantName,
				merchantStAddress,
				merchantCity,
				merchantState,
				merchantCountry,
			}

			if txn.SourceNotes != nil {
				data = append(data, *txn.SourceNotes)

			}

			if err := w.Write(data); err != nil {
				log.Fatalln("error writing record to csv:", err)
				return nil, err
			}
		}

		offset += limit
	}

	w.Flush()

	csv := CSVTransaction{
		Data: b.String(),
	}

	return &csv, nil
}

func (store transactionStore) Export(userID shared.UserID, businessID shared.BusinessID, startDate, endDate string) (*CSVTransaction, error) {
	sourceReq := services.NewSourceRequest()
	sourceReq.UserID = userID
	err := auth.NewAuthService(sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}
	params := make(map[string]interface{})
	params["businessId"] = businessID
	params["startDate"] = startDate
	params["endDate"] = endDate
	return store.ExportInternal(params)
}

func DeleteBusinessTransactions(businessIDs []shared.BusinessID) error {
	tx := DBWrite.MustBegin()

	// Delete transaction receipts
	query, args, err := sqlx.In(`
        DELETE FROM business_transaction_attachment
        WHERE transaction_id IN (select id FROM business_transaction WHERE business_id IN (?))`,
		businessIDs,
	)
	if err != nil {
		btr := DBWrite.Rebind(query)
		tx.MustExec(btr, args...)
	}

	// Delete transaction disputes
	query, args, err = sqlx.In(`
        DELETE FROM business_transaction_dispute
        WHERE transaction_id IN (select id FROM business_transaction WHERE business_id IN (?))`,
		businessIDs,
	)
	if err != nil {
		btd := DBWrite.Rebind(query)
		tx.MustExec(btd, args...)
	}

	// Delete card transactions
	query, args, err = sqlx.In(`
		DELETE FROM business_card_transaction
		WHERE transaction_id IN (select id FROM business_transaction WHERE business_id IN (?))`,
		businessIDs,
	)
	if err != nil {
		bct := DBWrite.Rebind(query)
		tx.MustExec(bct, args...)
	}

	// Delete hold transactions
	query, args, err = sqlx.In(`
		DELETE FROM business_hold_transaction
		WHERE transaction_id IN (select id FROM business_transaction WHERE business_id IN (?))`,
		businessIDs,
	)
	if err != nil {
		bcht := DBWrite.Rebind(query)
		tx.MustExec(bcht, args...)
	}

	// Delete transactions
	query, args, err = sqlx.In("DELETE FROM business_transaction WHERE business_id IN (?)", businessIDs)
	if err != nil {
		bt := DBWrite.Rebind(query)
		tx.MustExec(bt, args...)
	}

	return tx.Commit()
}

func (store *transactionStore) Update(u BusinessPostedTransactionUpdate, userID shared.UserID) (*BusinessPostedTransaction, error) {
	_, err := NewBusinessService().GetByID(shared.PostedTransactionID(u.TransactionID), userID, u.BusinessID)
	if err != nil {
		return nil, err
	}

	t := BusinessPostedTransactionAnnotation{}
	err = store.Get(&t, "SELECT * FROM business_transaction_annotation WHERE transaction_id = $1", u.TransactionID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err != nil {
		return store.createTransactionAnnotation(u)
	} else {
		return store.updateTransactionAnnotation(u)
	}
}

func (store *transactionStore) updateTransactionAnnotation(u BusinessPostedTransactionUpdate) (*BusinessPostedTransaction, error) {
	var columns []string

	if u.TransactionNotes == nil {
		return nil, errors.New("Notes cannot be empty")
	}

	columns = append(columns, "transaction_notes = :transaction_notes")

	_, err := store.NamedExec(
		fmt.Sprintf(
			"UPDATE business_transaction_annotation SET %s WHERE transaction_id = '%s'",
			strings.Join(columns, ", "),
			u.TransactionID,
		), u,
	)
	if err != nil {
		return nil, errors.Cause(err)
	}

	return NewBusinessService().GetByIDInternal(shared.PostedTransactionID(u.TransactionID))
}

func (store *transactionStore) createTransactionAnnotation(u BusinessPostedTransactionUpdate) (*BusinessPostedTransaction, error) {
	if u.TransactionNotes == nil {
		return nil, errors.New("Notes cannot be empty")
	}

	// Default/mandatory fields
	columns := []string{
		"transaction_id", "transaction_notes",
	}
	// Default/mandatory values
	values := []string{
		":transaction_id", ":transaction_notes",
	}

	sql := fmt.Sprintf("INSERT INTO business_transaction_annotation(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))

	stmt, err := store.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	t := &BusinessPostedTransactionAnnotation{}

	err = stmt.Get(t, &u)
	if err != nil {
		return nil, err
	}

	return NewBusinessService().GetByIDInternal(shared.PostedTransactionID(u.TransactionID))
}

func getShopifyPayoutDetails(t *BusinessPostedTransaction) {
	sn, err := grpc.GetConnectionStringForService(grpc.ServiceNameShopping)
	if err != nil {
		log.Println(err)
		return
	}

	client, err := grpc.NewInsecureClient(sn)
	if err != nil {
		log.Println(err)
		return
	}

	defer client.CloseAndCancel()
	shopifyServiceClient := shopify.NewShopifyBusinessServiceClient(client.GetConn())

	req := shopify.ShopifyTransactionRequest{
		TransactionId: t.ID.ToPrefixString(),
	}

	payout, err := shopifyServiceClient.GetShopifyPayoutByTransactionID(client.GetContext(), &req)
	if err != nil {
		log.Println(err)
		return
	}

	amount, err := num.ParseDecimal(payout.Amount)
	if err != nil {
		log.Println(err)
		return
	}

	adjustmentsFeeAmount, err := num.ParseDecimal(payout.AdjustmentsFeeAmount)
	if err != nil {
		log.Println(err)
		return
	}

	adjustmentsGrossAmount, err := num.ParseDecimal(payout.AdjustmentsGrossAmount)
	if err != nil {
		log.Println(err)
		return
	}

	chargesFeeAmount, err := num.ParseDecimal(payout.ChargesFeeAmount)
	if err != nil {
		log.Println(err)
		return
	}

	chargesGrossAmount, err := num.ParseDecimal(payout.ChargesGrossAmount)
	if err != nil {
		log.Println(err)
		return
	}

	refundsFeeAmount, err := num.ParseDecimal(payout.RefundsFeeAmount)
	if err != nil {
		log.Println(err)
		return
	}

	refundsGrossAmount, err := num.ParseDecimal(payout.RefundsGrossAmount)
	if err != nil {
		log.Println(err)
		return
	}

	reservedFundsFeeAmount, err := num.ParseDecimal(payout.ReservedFundsFeeAmount)
	if err != nil {
		log.Println(err)
		return
	}

	reservedFundsGrossAmount, err := num.ParseDecimal(payout.ReservedFundsGrossAmount)
	if err != nil {
		log.Println(err)
		return
	}

	retriedPayoutsFeeAmount, err := num.ParseDecimal(payout.RetriedPayoutsFeeAmount)
	if err != nil {
		log.Println(err)
		return
	}

	retriedPayoutsGrossAmount, err := num.ParseDecimal(payout.RetriedPayoutsGrossAmount)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("payout date is ", payout.PayoutDate)

	payoutDate, err := grpcTypes.Timestamp(payout.PayoutDate)
	if err != nil {
		log.Println(err)
		return
	}

	shopifyPayout := ShopifyPayout{
		PayoutID:                  payout.PayoutId,
		PayoutStatus:              payout.PayoutStatus,
		Currency:                  payout.Currency,
		Amount:                    amount,
		PayoutDate:                payoutDate,
		AdjustmentsFeeAmount:      adjustmentsFeeAmount,
		AdjustmentsGrossAmount:    adjustmentsGrossAmount,
		ChargesFeeAmount:          chargesFeeAmount,
		ChargesGrossAmount:        chargesGrossAmount,
		RefundsFeeAmount:          refundsFeeAmount,
		RefundsGrossAmount:        refundsGrossAmount,
		ReservedFundsFeeAmount:    reservedFundsFeeAmount,
		ReservedFundsGrossAmount:  reservedFundsGrossAmount,
		RetriedPayoutsFeeAmount:   retriedPayoutsFeeAmount,
		RetriedPayoutsGrossAmount: retriedPayoutsGrossAmount,
	}

	t.ShopifyPayout = &shopifyPayout
}
