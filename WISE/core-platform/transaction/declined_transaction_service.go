package transaction

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/services"
	bsrv "github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/go-lib/grpc"
	"github.com/wiseco/go-lib/id"
	"github.com/wiseco/go-lib/num"
	"github.com/wiseco/protobuf/golang"
	grpcTxn "github.com/wiseco/protobuf/golang/transaction"
	grpcBankTxn "github.com/wiseco/protobuf/golang/transaction/bank"
)

type DeclinedTransactionService interface {
	// Fetch business transactions
	ListAllInternal(params map[string]interface{}, businessID shared.BusinessID) ([]BusinessPendingTransaction, error)
	GetByIDInternal(ID shared.PendingTransactionID, businessID shared.BusinessID) (*BusinessPendingTransaction, error)

	// Export csv
	ExportInternal(businessID *shared.BusinessID, startDate, endDate string) (*CSVTransaction, error)
}

type declinedTxnStore struct {
	cardService    PendingCardTransactionService
	accountService AccountService
	*sqlx.DB
}

// New returns a new transaction service
func NewDeclinedTransactionService() DeclinedTransactionService {
	return &declinedTxnStore{NewPendingCardService(), NewAccountService(), DBWrite}
}

func (store declinedTxnStore) ListAllInternal(params map[string]interface{}, businessID shared.BusinessID) ([]BusinessPendingTransaction, error) {
	list := []BusinessPendingTransaction{}

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

		if params["minAmount"].(string) != "" && params["maxAmount"].(string) != "" {
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
		} else if params["minAmount"].(string) != "" {
			req.AmountRange.Filter = grpcBankTxn.AmountRangeFilter_ARF_MIN
			req.AmountRange.AmountMin = params["minAmount"].(string)
		} else if params["maxAmount"].(string) != "" {
			req.AmountRange.Filter = grpcBankTxn.AmountRangeFilter_ARF_MAX
			req.AmountRange.AmountMax = params["maxAmount"].(string)
		}

		if params["type"] != "" {
			codeType, ok := params["type"].(string)
			if ok {
				switch TransactionCodeType(codeType) {
				case TransactionCodeTypeDebitInProcess:
					req.AmountRange.Type = grpcBankTxn.AmountRangeType_ART_DEBIT
					req.StatusFilter = []grpcTxn.BankTransactionStatus{
						grpcTxn.BankTransactionStatus_BTS_HOLD_SET,
						grpcTxn.BankTransactionStatus_BTS_CARD_AUTHORIZED,
						grpcTxn.BankTransactionStatus_BTS_VALIDATION,
						grpcTxn.BankTransactionStatus_BTS_REVIEW,
						grpcTxn.BankTransactionStatus_BTS_BANK_PROCESSING,
						grpcTxn.BankTransactionStatus_BTS_BANK_TRANSFER_ERROR,
					}
				case TransactionCodeTypeCreditInProcess:
					req.AmountRange.Type = grpcBankTxn.AmountRangeType_ART_CREDIT
					req.StatusFilter = []grpcTxn.BankTransactionStatus{
						grpcTxn.BankTransactionStatus_BTS_HOLD_SET,
						grpcTxn.BankTransactionStatus_BTS_CARD_AUTHORIZED,
						grpcTxn.BankTransactionStatus_BTS_VALIDATION,
						grpcTxn.BankTransactionStatus_BTS_REVIEW,
						grpcTxn.BankTransactionStatus_BTS_BANK_PROCESSING,
						grpcTxn.BankTransactionStatus_BTS_BANK_TRANSFER_ERROR,
					}
				case TransactionCodeTypeAuthApproved:
					req.StatusFilter = []grpcTxn.BankTransactionStatus{grpcTxn.BankTransactionStatus_BTS_CARD_AUTHORIZED}
				case TransactionCodeTypeHoldApproved:
					req.StatusFilter = []grpcTxn.BankTransactionStatus{grpcTxn.BankTransactionStatus_BTS_HOLD_SET}
				case TransactionCodeTypeValidation:
					req.StatusFilter = []grpcTxn.BankTransactionStatus{grpcTxn.BankTransactionStatus_BTS_VALIDATION}
				case TransactionCodeTypeReview:
					req.StatusFilter = []grpcTxn.BankTransactionStatus{grpcTxn.BankTransactionStatus_BTS_REVIEW}
				case TransactionCodeTypeTransferError:
					req.StatusFilter = []grpcTxn.BankTransactionStatus{grpcTxn.BankTransactionStatus_BTS_BANK_TRANSFER_ERROR}
				case TransactionCodeTypeDeclined:
					req.StatusFilter = []grpcTxn.BankTransactionStatus{
						grpcTxn.BankTransactionStatus_BTS_CARD_AUTH_DECLINED,
						grpcTxn.BankTransactionStatus_BTS_SYSTEM_DECLINED,
						grpcTxn.BankTransactionStatus_BTS_AGENT_DECLINED,
						grpcTxn.BankTransactionStatus_BTS_BANK_DECLINED,
					}
				case TransactionCodeTypeCanceled:
					req.StatusFilter = []grpcTxn.BankTransactionStatus{
						grpcTxn.BankTransactionStatus_BTS_AGENT_CANCELED,
						grpcTxn.BankTransactionStatus_BTS_CUSTOMER_CANCELED,
						grpcTxn.BankTransactionStatus_BTS_BANK_CANCELED,
					}
				}
			}
		} else {
			req.StatusFilter = []grpcTxn.BankTransactionStatus{
				grpcTxn.BankTransactionStatus_BTS_CARD_AUTH_DECLINED,
				grpcTxn.BankTransactionStatus_BTS_SYSTEM_DECLINED,
				grpcTxn.BankTransactionStatus_BTS_AGENT_DECLINED,
				grpcTxn.BankTransactionStatus_BTS_BANK_DECLINED,
			}
		}

		req.Offset = int32(params["offset"].(int))
		req.Limit = int32(params["limit"].(int))

		req.SortRequests = []*grpcBankTxn.SortRequest{
			&grpcBankTxn.SortRequest{
				Name:      grpcBankTxn.SortFieldName_SFN_TRANSACTION_DATE,
				Direction: golang.SortDirection_SD_DESCENDING,
			},
		}

		if params["text"] != "" {
			req.SearchTerms = strings.TrimSpace(params["text"].(string))
		}

		req.TypeFilter = []grpcTxn.BankTransactionType{}
		subtype, ok := params["subtype"].(string)
		if ok && params["subtype"] != "" {
			st := TransactionSubtype(subtype)
			typeProto, ok := TransactionSubtypeToTypeProto[st]

			stp := []grpcTxn.BankTransactionType{
				typeProto,
			}
			if !ok {
				return nil, fmt.Errorf("Invalid subtype")
			}
			req.TypeFilter = stp
		}

		resp, err := tc.GetMany(context.Background(), req)
		if err != nil {
			return list, err
		}

		for _, gtxn := range resp.Results {
			t, err := BusinessPendingTransactionFromProto(gtxn, bus)
			if err != nil {
				return list, err
			}
			list = append(list, *t)
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
		transaction_status, partner_name, source_notes, source_notes AS notes, money_request_id, created, transaction_title, transaction_subtype`

		var err error
		if businessID != "" {
			query := `SELECT ` + columns + ` FROM business_pending_transaction WHERE amount != 0 AND business_id = $1` +
				dateFilter + txnTypeFilter + contactFilter + amtFilter + txtFilter +
				` ORDER BY transaction_date DESC LIMIT $2 OFFSET $3`

			err = store.Select(&list, query, businessID, params["limit"].(int), params["offset"].(int))
		} else {
			query := `SELECT ` + columns + ` FROM business_pending_transaction WHERE amount != 0` +
				dateFilter + txnTypeFilter + contactFilter + amtFilter + txtFilter +
				` ORDER BY transaction_date DESC LIMIT $1 OFFSET $2`

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

	return list, nil
}

// Fetches all transactions without offset and limit restriction
func (store declinedTxnStore) listAllForExport(businessID *shared.BusinessID, startDate, endDate string, offset, limit int) ([]BusinessPendingTransaction, error) {
	list := []BusinessPendingTransaction{}

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
			grpcTxn.BankTransactionStatus_BTS_CARD_AUTH_DECLINED,
			grpcTxn.BankTransactionStatus_BTS_SYSTEM_DECLINED,
			grpcTxn.BankTransactionStatus_BTS_AGENT_DECLINED,
			grpcTxn.BankTransactionStatus_BTS_BANK_DECLINED,
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

			t, err := BusinessPendingTransactionFromFullProto(gtxnFull, bus)
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

		query := `SELECT * FROM business_pending_transaction WHERE ` + filterText + ` ORDER BY transaction_date DESC OFFSET $1 LIMIT $1`

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

func (store declinedTxnStore) GetByIDInternal(ptID shared.PendingTransactionID, businessID shared.BusinessID) (*BusinessPendingTransaction, error) {
	t := &BusinessPendingTransaction{}

	if os.Getenv("USE_TRANSACTION_SERVICE") == "true" {
		bus, err := bsrv.NewBusinessServiceWithout().GetByIdInternal(businessID)
		if err != nil {
			return t, err
		}

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

		busID, _ := id.ParseBusinessID(string(businessID))
		if !busID.IsZero() {
			if resp.Transaction.BusinessId != busID.String() {
				return t, errors.New("invalid business id")
			}
		}

		t, err = BusinessPendingTransactionFromFullProto(resp, bus)
		if err != nil {
			return t, err
		}
	} else {
		query := `
        SELECT
            business_pending_transaction.*,
			business_pending_transaction.bank_transaction_desc AS money_transfer_desc,
			business_pending_transaction.source_notes AS notes,
			business_card_pending_transaction.id "business_card_pending_transaction.id", 
			business_card_pending_transaction.auth_amount "business_card_pending_transaction.auth_amount", 
			business_card_pending_transaction.transaction_type "business_card_pending_transaction.transaction_type",
			business_card_pending_transaction.local_amount "business_card_pending_transaction.local_amount",
			business_card_pending_transaction.local_currency "business_card_pending_transaction.local_currency", 
			business_card_pending_transaction.local_date "business_card_pending_transaction.local_date",
			business_card_pending_transaction.billing_currency "business_card_pending_transaction.billing_currency",
			business_card_pending_transaction.merchant_category_code "business_card_pending_transaction.merchant_category_code",
			business_card_pending_transaction.merchant_name "business_card_pending_transaction.merchant_name",
			business_card_pending_transaction.merchant_street_address "business_card_pending_transaction.merchant_street_address",
			business_card_pending_transaction.merchant_city "business_card_pending_transaction.merchant_city", 
			business_card_pending_transaction.merchant_state "business_card_pending_transaction.merchant_state", 
			business_card_pending_transaction.merchant_country "business_card_pending_transaction.merchant_country", 			 
			business_hold_pending_transaction.id "business_hold_pending_transaction.id",
			business_hold_pending_transaction.amount "business_hold_pending_transaction.amount",
			business_hold_pending_transaction.hold_number "business_hold_pending_transaction.hold_number", 
			business_hold_pending_transaction.transaction_date "business_hold_pending_transaction.transaction_date",
			business_hold_pending_transaction.expiry_date "business_hold_pending_transaction.expiry_date"
        FROM business_pending_transaction
        LEFT JOIN business_card_pending_transaction ON business_pending_transaction.id = business_card_pending_transaction.transaction_id
        LEFT JOIN business_hold_pending_transaction ON business_pending_transaction.id = business_hold_pending_transaction.transaction_id
        WHERE business_pending_transaction.id = $1 AND business_pending_transaction.business_id = $2`

		err := store.Get(t, query, ptID, businessID)
		if err != nil {
			log.Println("Error fetching pending transaction ", err, ptID, businessID)
			return nil, err
		}
	}

	// Money Transfer
	// TODO: USE_BANKING_SERVICE
	if t.MoneyTransferID != nil {
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

	if t.MoneyRequestID != nil {
		query := `
		SELECT
		business_money_request_payment.id "business_money_request_payment.id", 
		business_money_request_payment.payment_date "business_money_request_payment.payment_date", 
		business_money_request_payment.card_brand "business_money_request_payment.card_brand", 
		business_money_request_payment.card_number "business_money_request_payment.card_number"
		FROM business_money_request_payment
		WHERE business_money_request_payment.request_id = $1 or business_money_request_payment.invoice_id = $2`

		e := data.DBRead.Get(t, query, *t.MoneyRequestID, *t.MoneyRequestID)
		if e != nil {
			log.Println("Error fetching money request invoice and receipt", e)
		}
	}

	if t.BankTransactionDesc != nil {
		t.OriginAccount = GetOriginAccount(TransactionTypeACH, shared.StringValue(t.BankTransactionDesc))
		t.DestinationAccount = GetDestinationAccount(TransactionTypeACH, shared.StringValue(t.BankTransactionDesc))
	}

	switch t.TransactionType {
	case TransactionTypeTransfer:
		handleDeclinedTransferTransaction(t)
	case TransactionTypeOther: // Card authorization approved
		handleDeclinedCardTransaction(t)
	default:
		log.Println("Invalid transaction type", t.TransactionType)
	}

	return t, nil
}

func handleDeclinedTransferTransaction(t *BusinessPendingTransaction) {
	switch t.CodeType {
	case TransactionCodeTypeCreditInProcess:
		s := TransactionSource{}

		if t.MoneyRequestID != nil {
			if t.Payment != nil {
				s.BankName = t.Payment.CardBrand
				s.AccountNumber = t.Payment.CardLast4
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

		s.AccountNumber = GetOriginAccount(TransactionTypeACH, shared.StringValue(t.MoneyTransferDesc))
		s.AccountHolderName = GetOriginAccountHolder(TransactionTypeACH, shared.StringValue(t.MoneyTransferDesc))
		t.Source = &s

		return

	case TransactionCodeTypeDebitInProcess:
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

			if t.Contact.BankName != nil && len(*t.Contact.BankName) > 0 {
				d.BankName = t.Contact.BankName
			}
		}

		if d.AccountNumber == nil {
			d.AccountNumber = GetDestinationAccount(TransactionTypeACH, shared.StringValue(t.MoneyTransferDesc))
		}

		if d.AccountHolderName == nil {
			d.AccountHolderName = GetDestinationAccountHolder(TransactionTypeACH, shared.StringValue(t.MoneyTransferDesc))
		}

		t.Destination = &d
	default:
		log.Println("Invalid transaction code type ", t.TransactionType)
		return

	}
}

func handleDeclinedCardTransaction(t *BusinessPendingTransaction) {
	switch t.CodeType {
	case TransactionCodeTypeAuthApproved:
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
	default:
		log.Println("Invalid transaction code type ", t.CodeType)

	}
}

func (store declinedTxnStore) ExportInternal(businessID *shared.BusinessID, startDate, endDate string) (*CSVTransaction, error) {

	var b bytes.Buffer
	w := csv.NewWriter(&b)

	header := []string{
		"S No.", "Transaction Date", "Transaction Remarks", "Transaction Type", "Amount(USD)", "Authorization Amount",
		"Merchant Name", "Merchant Street Address", "Merchant City", "Merchant State", "Merchant Country",
	}

	err := w.Write(header)
	if err != nil {
		log.Fatalln("error writing record to csv:", err)
		return nil, err
	}

	offset := 0
	limit := 20
	for {
		t, err := store.listAllForExport(businessID, startDate, endDate, offset, limit)
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
