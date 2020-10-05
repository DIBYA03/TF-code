package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"sync"
	"unicode"

	"github.com/google/uuid"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/go-lib/grpc"
	"github.com/wiseco/go-lib/id"

	grpcTypes "github.com/golang/protobuf/ptypes"
	"github.com/wiseco/core-platform/notification"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking/business"
	"github.com/wiseco/core-platform/services/contact"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/services/payment"
	"github.com/wiseco/core-platform/services/user"
	"github.com/wiseco/core-platform/transaction"
	grpcBanking "github.com/wiseco/protobuf/golang/banking"
	grpcTxn "github.com/wiseco/protobuf/golang/transaction"
	grpcBankTxn "github.com/wiseco/protobuf/golang/transaction/bank"
)

type TransactionDetails struct {
	ContactFirstName    *string              `db:"business_contact.first_name"`
	ContactLastName     *string              `db:"business_contact.last_name"`
	ContactBusinessName *string              `db:"business_contact.business_name"`
	ContactType         *contact.ContactType `db:"business_contact.contact_type"`
	ContactName         *string
	MonthlyInterestID   *string                     `db:"business_money_transfer.account_monthly_interest_id"`
	Notes               *string                     `db:"business_money_transfer.notes"`
	MoneyRequestType    *payment.PaymentRequestType `db:"business_money_request.request_type"`
	PaymentLocation     *services.Address           `db:"business_money_request_payment.purchase_address"`
	InterestStartDate   shared.Date                 `db:"interest_start_date"`
	InterestEndDate     shared.Date                 `db:"interest_end_date"`
}

type TransactionUpdate struct {
	CounterpartyName       string
	TransactionTitle       string `db:"transaction_title"`
	TransactionDescription string `db:"transaction_desc"`
	Status                 transaction.TransactionStatus
	TransactionSubType     transaction.TransactionSubtype `db:"transaction_subtype"`
}

type BusinessDetails struct {
	MaskedCardNumber *string              `db:"business_bank_card.card_number_masked"`
	LegalName        string               `db:"business.legal_name"`
	DBA              services.StringArray `db:"business.dba"`
	BusinessName     string
}

func isOnlinePayment(posEntryMode string) bool {
	m := notification.POSEntryMode(posEntryMode[0:2])
	switch m {
	case notification.POSEntryModeKeyEntry, notification.POSEntryModeManualECommerce, notification.POSEntryModeStoredCheckout:
		return true
	default:
		return false
	}
}

func main() {
	sn, err := grpc.GetConnectionStringForService(grpc.ServiceNameTransaction)
	if err != nil {
		panic(err)
	}

	log.Println(sn)

	cl, err := grpc.NewInsecureClient(sn)
	if err != nil {
		panic(err)
	}

	defer cl.CloseAndCancel()
	txnClient := grpcBankTxn.NewTransactionServiceClient(cl.GetConn())

	// Read all transactions
	offset := 0
	for {
		transactions, err := getAllTransactions(offset, 10)
		if err != nil {
			log.Println("error fetching notifications", err)
			return
		} else if len(transactions) == 0 {
			log.Println("no more notifications:", offset, 10)
			return
		}

		var wg sync.WaitGroup
		wg.Add(len(transactions))
		for _, t := range transactions {
			go func(t transaction.BusinessPostedTransaction) {
				defer wg.Done()

				// Query core DB to get business details
				businessDetail, err := getBusinessDetails(t.BusinessID, t.CardID)
				if err != nil {
					log.Println(err, t.ID, t.BusinessID)
					return
				}

				// process transaction
				processTransaction(txnClient, t, businessDetail)
			}(t)
		}

		wg.Wait()
		offset += len(transactions)
	}
}

func getAllTransactions(offset, limit int) ([]transaction.BusinessPostedTransaction, error) {
	transactions := []transaction.BusinessPostedTransaction{}

	transactionDate := "2019-6-01"

	query := `SELECT
	business_transaction.*,
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
	WHERE business_transaction.transaction_date > $1
	ORDER BY business_transaction.transaction_date ASC
	OFFSET $2
	LIMIT $3`

	err := transaction.DBWrite.Select(&transactions, query, transactionDate, offset, limit)
	if err != nil && err == sql.ErrNoRows {
		return transactions, nil
	}

	if err != nil {
		log.Println("getAllTransactions:", err)
		return nil, err
	}

	return transactions, nil
}

func getTransactionDetailsByMoneyTransferID(moneyTransferID string, moneyRequestID *shared.PaymentRequestID) (*TransactionDetails, error) {
	t := TransactionDetails{}

	var query string
	if moneyRequestID != nil {
		query = `
		SELECT
		business_money_transfer.account_monthly_interest_id "business_money_transfer.account_monthly_interest_id",
		business_money_transfer.notes "business_money_transfer.notes",
		business_contact.first_name "business_contact.first_name",
		business_contact.last_name "business_contact.last_name",
		business_contact.business_name "business_contact.business_name",
		business_contact.contact_type "business_contact.contact_type",
		business_money_request.request_type "business_money_request.request_type",
		business_money_request_payment.purchase_address "business_money_request_payment.purchase_address"
		FROM business_money_transfer
		JOIN business_money_request ON business_money_transfer.money_request_id = business_money_request.id
		LEFT JOIN business_money_request_payment ON business_money_request_payment.request_id = business_money_request.id
		LEFT JOIN business_contact ON business_money_request.contact_id = business_contact.id
		WHERE business_money_transfer.id = $1`
	} else {
		query = `
		SELECT
		business_money_transfer.account_monthly_interest_id "business_money_transfer.account_monthly_interest_id",
		business_money_transfer.notes "business_money_transfer.notes",
		business_contact.first_name "business_contact.first_name",
		business_contact.last_name "business_contact.last_name",
		business_contact.business_name "business_contact.business_name",
		business_contact.contact_type "business_contact.contact_type",
		business_money_request.request_type "business_money_request.request_type"
		FROM business_money_transfer
		LEFT JOIN business_contact ON business_money_transfer.contact_id = business_contact.id
		LEFT JOIN business_money_request ON business_money_transfer.money_request_id = business_money_request.id
		WHERE business_money_transfer.id = $1`
	}

	err := data.DBRead.Get(&t, query, moneyTransferID)
	if err != nil {
		log.Println("Error fetching money transfer details", err, moneyTransferID)
		return nil, err
	}

	if t.ContactType != nil {
		switch *t.ContactType {
		case contact.ContactTypePerson:
			name := *t.ContactFirstName + " " + *t.ContactLastName
			t.ContactName = &name
		case contact.ContactTypeBusiness:
			name := *t.ContactBusinessName
			t.ContactName = &name
		default:
			log.Println("Unknown contact type", *t.ContactType)
			return nil, errors.New("unknown contact type")
		}
	}

	return &t, nil
}

func getBusinessDetails(businessID shared.BusinessID, cardID *string) (*BusinessDetails, error) {
	bt := BusinessDetails{}
	var err error

	if cardID != nil {
		err = data.DBRead.Get(
			&bt,
			`SELECT business.legal_name "business.legal_name", business.dba "business.dba", 
				business_bank_card.card_number_masked "business_bank_card.card_number_masked"
				FROM business 
				LEFT JOIN business_bank_card ON business_bank_card.business_id = business.id
				WHERE business_bank_card.business_id = $1 AND business_bank_card.id = $2`,
			businessID, *cardID,
		)
	} else {
		err = data.DBRead.Get(
			&bt,
			`SELECT business.legal_name "business.legal_name", business.dba "business.dba"
				FROM business 
				WHERE business.id = $1`,
			businessID,
		)
	}

	if err != nil {
		log.Println("Error retrieving business details ", err, businessID)
		return nil, err
	}

	bt.BusinessName = shared.GetBusinessName(&bt.LegalName, bt.DBA)

	return &bt, nil

}

func processTransaction(txnClient grpcBankTxn.TransactionServiceClient, t transaction.BusinessPostedTransaction, businessDetail *BusinessDetails) {
	var transactionDetail *TransactionDetails
	var err error
	u := &TransactionUpdate{}

	// Account Transaction ID
	var atxnID id.BankTransactionID
	if t.ID != "" {
		atxnUUID, err := uuid.Parse(string(t.ID))
		if err != nil {
			log.Printf("Error parsing account transaction id: %v %s", err, t.ID)
			return
		}

		atxnID = id.BankTransactionID(atxnUUID)
	}

	// Business ID
	var bID id.BusinessID
	if t.BusinessID != "" {
		bUUID, err := uuid.Parse(string(t.BusinessID))
		if err != nil {
			log.Printf("Error parsing business id: %v %s", err, t.ID)
			return
		}

		bID = id.BusinessID(bUUID)
	}

	// Bank Account ID
	var accID id.BankAccountID
	if t.AccountID != nil && len(*t.AccountID) > 0 {
		accUUID, err := uuid.Parse(shared.StringValue(t.AccountID))
		if err != nil {
			log.Printf("Error parsing bank account id: %v %s", err, t.ID)
			return
		}

		accID = id.BankAccountID(accUUID)
	}

	// Account must always exist
	acc, err := business.NewAccountService().GetByIDInternal(accID.UUIDString())
	if err != nil {
		log.Printf("Error getting bank account id: %v %s", err, t.ID)
		return
	}

	if t.MoneyTransferID != nil {
		transactionDetail, err = getTransactionDetailsByMoneyTransferID(*t.MoneyTransferID, t.MoneyRequestID)
		if err != nil {
			log.Println("processTransaction:", err, t.ID, t.MoneyRequestID)
			return
		}
	}

	if acc.UsageType == business.UsageTypeClearing && transactionDetail != nil {
		transactionDetail.ContactFirstName = nil
		transactionDetail.ContactLastName = nil
		transactionDetail.ContactBusinessName = nil
		transactionDetail.ContactType = nil
		transactionDetail.ContactName = nil
	}

	switch t.TransactionType {
	case transaction.TransactionTypeTransfer:
		u, err = processTransferTransactions(t, transactionDetail, businessDetail)
		if err != nil {
			log.Println("Error processing transfer transaction", err, t.ID)
			return
		}
	case transaction.TransactionTypeDeposit:
		amtFloat, _ := t.Amount.Float64()
		if amtFloat == 0 {
			u, err = processAccountOriginationTransaction(t, transactionDetail, businessDetail)
		} else {
			u, err = processDepositCreditTransaction(t, transactionDetail, businessDetail)
		}

		if err != nil {
			log.Println("Error processing deposit transaction", err, t.ID)
			return
		}
	case transaction.TransactionTypeACH:
		u, err = processACHTransaction(t, transactionDetail, businessDetail)
		if err != nil {
			log.Println("Error processing ach transaction", err, t.ID)
			return
		}
	case transaction.TransactionTypeATM:
		u, err = processATMTransaction(t, transactionDetail, businessDetail)
		if err != nil {
			log.Println("Error processing atm transaction", err, t.ID)
			return
		}
	case transaction.TransactionTypePurchase:
		if transactionDetail == nil {
			u, err = processCardPurchaseDebitTransaction(t, transactionDetail, businessDetail)
		} else {
			u, err = processInstantPayDebitTransaction(t, transactionDetail, businessDetail)
		}

		if err != nil {
			log.Println("Error processing card purchase transaction", err, t.ID)
			return
		}
	case transaction.TransactionTypeRefund:
		u, err = processMerchantRefundDebitTransaction(t, transactionDetail, businessDetail)
		if err != nil {
			log.Println("Error processing refund transaction", err, t.ID)
			return
		}
	case transaction.TransactionTypeFee:
		u = &TransactionUpdate{
			TransactionTitle:       shared.StringValue(t.TransactionTitle),
			TransactionDescription: shared.StringValue(t.BankTransactionDesc),
			TransactionSubType:     transaction.TransactionSubtypeFeeDebit,
		}
		break
	case transaction.TransactionTypeVisaCredit:
		u, err = processVisaCreditTransaction(t, transactionDetail, businessDetail)
		if err != nil {
			log.Println("Error processing debit pull transaction", err, t.ID)
		}
	case transaction.TransactionTypeOtherCredit:
		u = &TransactionUpdate{
			TransactionTitle:       shared.StringValue(t.TransactionTitle),
			TransactionDescription: shared.StringValue(t.BankTransactionDesc),
			TransactionSubType:     transaction.TransactionSubtypeOtherCredit,
		}
	case transaction.TransactionTypeOtherDebit:
		u = &TransactionUpdate{
			TransactionTitle:       shared.StringValue(t.TransactionTitle),
			TransactionDescription: shared.StringValue(t.BankTransactionDesc),
			TransactionSubType:     transaction.TransactionSubtypeOtherDebit,
		}
	case transaction.TransactionTypeOther:
		u = &TransactionUpdate{
			TransactionTitle:       shared.StringValue(t.TransactionTitle),
			TransactionDescription: shared.StringValue(t.BankTransactionDesc),
			TransactionSubType:     transaction.TransactionSubtypeUnspecified,
		}
	default:
		log.Println("Unhandled transaction type", t.TransactionType, t.ID)
		return
	}

	// Debit Card ID
	var dbcID id.DebitCardID
	if t.CardID != nil && len(*t.CardID) > 0 {
		dbcUUID, err := uuid.Parse(shared.StringValue(t.CardID))
		if err != nil {
			log.Printf("Error parsing debit card id: %v %s", err, t.ID)
			return
		}

		dbcID = id.DebitCardID(dbcUUID)
	}

	// Money Transfer ID
	var btID id.BankTransferID
	if t.MoneyTransferID != nil && len(*t.MoneyTransferID) > 0 {
		btUUID, err := uuid.Parse(shared.StringValue(t.MoneyTransferID))
		if err != nil {
			log.Printf("Error parsing bank transfer id: %v %s", err, t.ID)
			return
		}
		btID = id.BankTransferID(btUUID)
	}

	var mt *business.MoneyTransfer
	if !btID.IsZero() {
		mt, err = business.NewMoneyTransferServiceWithout().GetByIDOnlyInternal(btID.UUIDString())
		if err != nil {
			log.Printf("Error getting bank transfer: %v", err, t.ID)
			return
		}
	}

	// Payment Request ID
	var prID id.PaymentRequestID
	if t.MoneyRequestID != nil && len(*t.MoneyRequestID) > 0 {
		prUUID, err := uuid.Parse(string(*t.MoneyRequestID))
		if err != nil {
			log.Printf("Error parsing payment request id: %v", err, t.ID)
			return
		}

		prID = id.PaymentRequestID(prUUID)
	}

	// Contact
	var contactID id.ContactID
	if t.ContactID != nil && len(*t.ContactID) > 0 {
		contactUUID, err := uuid.Parse(shared.StringValue(t.ContactID))
		if err != nil {
			log.Printf("Error parsing contact id: %v", err, t.ID)
			return
		}

		contactID = id.ContactID(contactUUID)
	}

	status, ok := transaction.TransactionStatusToProto[u.Status]
	if !ok || status == grpcTxn.BankTransactionStatus_BTS_UNSPECIFIED {
		status, ok = transaction.TransactionStatusToProto[t.Status]
		if !ok || status == grpcTxn.BankTransactionStatus_BTS_UNSPECIFIED {
			if dbcID.IsZero() {
				status = grpcTxn.BankTransactionStatus_BTS_NONCARD_POSTED
			} else {
				status = grpcTxn.BankTransactionStatus_BTS_CARD_POSTED
			}
		}
	}

	category, ok := transaction.TransactionTypeToCategoryProto[t.TransactionType]
	if !ok || category == grpcTxn.BankTransactionCategory_BTC_UNSPECIFIED {
		log.Println("invalid or unspecified transaction category type", t.ID)
		return
	}

	subtype := transaction.TransactionSubtypeUnspecified
	txnType, _ := transaction.TransactionSubtypeToTypeProto[transaction.TransactionSubtypeUnspecified]
	if t.TransactionSubtype != nil {
		subtype = *t.TransactionSubtype

		// Try transaction subtype
		txnType, ok = transaction.TransactionSubtypeToTypeProto[*t.TransactionSubtype]
		if !ok || txnType == grpcTxn.BankTransactionType_BTT_UNSPECIFIED {
			subtype = u.TransactionSubType

			// Try update subtype
			txnType, ok = transaction.TransactionSubtypeToTypeProto[u.TransactionSubType]
			if !ok || txnType == grpcTxn.BankTransactionType_BTT_UNSPECIFIED {
				log.Println("invalid or unspecified transaction subtype:", u.TransactionSubType, atxnID.String(), t.CodeType)
				return
			}
		}
	} else {
		subtype = u.TransactionSubType
		txnType, ok = transaction.TransactionSubtypeToTypeProto[u.TransactionSubType]
		if !ok || txnType == grpcTxn.BankTransactionType_BTT_UNSPECIFIED {
			log.Println("invalid or unspecified transaction subtype:", u.TransactionSubType, atxnID.String(), t.CodeType)
			return
		}
	}

	counterpartyType := grpcTxn.BankTransactionCounterpartyType_BTCT_UNSPECIFIED
	var interestDate shared.Date
	switch txnType {
	case grpcTxn.BankTransactionType_BTT_UNSPECIFIED:
		if acc.UsageType == business.UsageTypeClearing {
			if mt == nil {
				log.Println("clearing account transaction missing money transfer", t.ID)
				return
			}
			if acc.Id == os.Getenv("WISE_CLEARING_ACCOUNT_ID") {
				if mt.MonthlyInterestID != nil {
					txnType = grpcTxn.BankTransactionType_BTT_INTEREST_DEBIT
				} else {
					txnType = grpcTxn.BankTransactionType_BTT_CARD_PAYMENT_DEBIT
				}
			} else if acc.Id == os.Getenv("WISE_PROMO_CLEARING_ACCOUNT_ID") {
				txnType = grpcTxn.BankTransactionType_BTT_PROMO_DEBIT
			}
		}
	case grpcTxn.BankTransactionType_BTT_INTRABANK_ONLINE_CREDIT:
		if category == grpcTxn.BankTransactionCategory_BTC_ACH {
			txnType = grpcTxn.BankTransactionType_BTT_ACH_ONLINE_CREDIT
		}
	case grpcTxn.BankTransactionType_BTT_INTEREST_CREDIT:
		if transactionDetail != nil && !transactionDetail.InterestStartDate.IsZero() {
			interestDate = transactionDetail.InterestStartDate
		}
	case grpcTxn.BankTransactionType_BTT_ACH_CREDIT:
		if t.TransactionSubtype != nil && *t.TransactionSubtype == transaction.TransactionSubtypeACHTransferShopifyCredit {
			counterpartyType = grpcTxn.BankTransactionCounterpartyType_BTCT_SHOPPING_SHOPIFY_PAYOUT
		}
	case grpcTxn.BankTransactionType_BTT_ACH_DEBIT:
		if t.TransactionSubtype != nil && *t.TransactionSubtype == transaction.TransactionSubtypeACHTransferShopifyDebit {
			counterpartyType = grpcTxn.BankTransactionCounterpartyType_BTCT_SHOPPING_SHOPIFY_REFUND
		}
	}

	txnDate, err := grpcTypes.TimestampProto(t.TransactionDate)
	if err != nil {
		log.Println(err, t.ID)
		return
	}

	var title = shared.StringValue(t.TransactionTitle)
	if title == "" {
		title = u.TransactionTitle
	}

	desc := t.TransactionDesc
	if desc == "" {
		desc = u.TransactionDescription
	}

	created, err := grpcTypes.TimestampProto(t.Created)
	if err != nil {
		log.Println(err, t.ID)
		return
	}

	// Transactions are first so generate event and thread id - back patch banking
	evID, err := id.NewEventID()
	if err != nil {
		log.Println(err, t.ID)
		return
	}

	evThreadID, err := id.NewEventThreadID()
	if err != nil {
		log.Println(err, t.ID)
		return
	}

	notes := shared.StringValue(t.Notes)
	if notes == "" && transactionDetail != nil {
		notes = shared.StringValue(transactionDetail.Notes)
	}

	req := &grpcBankTxn.UpsertTransactionRequest{
		Id:                     atxnID.String(),
		BusinessId:             bID.String(),
		AccountId:              accID.String(),
		DebitCardId:            dbcID.String(),
		BankTransferId:         btID.String(),
		PaymentRequestId:       prID.String(),
		ContactId:              contactID.String(),
		PartnerName:            grpcBanking.PartnerName_PN_BBVA,
		PartnerTransactionId:   shared.StringValue(t.BankTransactionID),
		PartnerTransactionDesc: strings.TrimSpace(shared.StringValue(t.BankTransactionDesc)),
		EventId:                evID.String(),
		EventThreadId:          evThreadID.String(),
		Status:                 status,
		Category:               category,
		Type:                   txnType,
		Amount:                 t.Amount.FormatCurrency(),
		Currency:               string(t.Currency),
		TransactionDate:        txnDate,
		Counterparty:           strings.TrimSpace(u.CounterpartyName),
		CounterpartyType:       counterpartyType,
		InterestDate:           interestDate.String(),
		Created:                created,
		LegacyType:             string(t.TransactionType),
		LegacyCodeType:         string(t.CodeType),
		LegacySubtype:          string(subtype),
		LegacyTitle:            title,
		LegacyDescription:      desc,
		LegacyNotes:            notes,
	}

	if t.CardTransaction != nil {
		card := &transaction.BusinessCardTransaction{}
		query := "SELECT * FROM business_card_transaction where id = $1 LIMIT 1"
		err := transaction.DBWrite.Get(card, query, t.CardTransaction.ID)
		if err == nil {
			network, ok := transaction.CardNetworkToProto[strings.ToUpper(card.TransactionNetwork)]
			if !ok {
				log.Println("invalid card network", card.TransactionNetwork, t.ID)
				return
			}

			usr, err := user.NewUserServiceWithout().GetByIdInternal(card.CardHolderID)
			if err != nil {
				log.Println(err)
				return
			}

			cardHolderUUID, _ := id.ParseUUID(string(usr.ConsumerID))
			authDate, _ := grpcTypes.TimestampProto(card.AuthDate)
			localDate, _ := grpcTypes.TimestampProto(card.LocalDate)
			created, _ := grpcTypes.TimestampProto(card.Created)
			req.CardRequest = &grpcBankTxn.UpsertCardTransactionRequest{
				CardHolderId:          id.ConsumerID(cardHolderUUID).String(),
				NetworkTransactionId:  card.CardTransactionID,
				Network:               network,
				AuthAmount:            card.AuthAmount.FormatCurrency(),
				AuthDate:              authDate,
				AuthResponseCode:      card.AuthResponseCode,
				AuthNumber:            card.AuthNumber,
				CardTransactionType:   string(card.TransactionType),
				LocalAmount:           card.LocalAmount.FormatCurrency(),
				LocalCurrency:         card.LocalCurrency,
				LocalDate:             localDate,
				BillingCurrency:       card.BillingCurrency,
				PosEntryMode:          card.POSEntryMode,
				PosConditionCode:      card.POSConditionCode,
				AcquirerBin:           card.AcquirerBIN,
				MerchantName:          card.MerchantName,
				MerchantCategoryCode:  card.MerchantCategoryCode,
				AcceptorId:            card.MerchantID,
				AcceptorTerminal:      card.MerchantTerminal,
				AcceptorStreetAddress: card.MerchantStreetAddress,
				AcceptorCity:          card.MerchantCity,
				AcceptorState:         card.MerchantState,
				AcceptorCountry:       card.MerchantCountry,
				Created:               created,
			}

			if card.POSEntryMode != "" && isOnlinePayment(card.POSEntryMode) {
				subtype = transaction.TransactionSubtypeCardPurchaseDebitOnline
				req.Type = grpcTxn.BankTransactionType_BTT_CARD_PURCHASE_ONLINE_DEBIT
			}
		} else if err != sql.ErrNoRows {
			log.Println(err, t.ID)
			return
		}
	}

	if t.HoldTransaction != nil {
		hold := &transaction.BusinessHoldTransaction{}
		query := "SELECT * FROM business_hold_transaction where id = $1 LIMIT 1"
		err := transaction.DBWrite.Get(hold, query, t.HoldTransaction.ID)
		if err == nil {
			txnDate, _ := grpcTypes.TimestampProto(hold.Date)
			expiryDate, _ := grpcTypes.TimestampProto(hold.ExpiryDate)
			created, _ := grpcTypes.TimestampProto(hold.Created)
			req.HoldRequest = &grpcBankTxn.UpsertHoldTransactionRequest{
				HoldNumber:      hold.Number,
				Amount:          hold.Amount.FormatCurrency(),
				TransactionDate: txnDate,
				ExpiryDate:      expiryDate,
				Created:         created,
			}
		} else if err != sql.ErrNoRows {
			log.Println(err, t.ID)
			return
		}
	}

	// Create transaction in new service
	resp, err := txnClient.Upsert(context.Background(), req)
	if err != nil {
		log.Println("Transaction service error:", err, t.ID)
	} else {
		log.Println("Transaction service success:", resp.Transaction.Id)
	}
}

func processTransferTransactions(t transaction.BusinessPostedTransaction,
	transactionDetail *TransactionDetails, businessDetail *BusinessDetails) (*TransactionUpdate, error) {

	if t.CodeType == transaction.TransactionCodeTypeCreditPosted {
		if t.MoneyRequestID != nil {
			return processMoneyRequestCreditTransaction(t, transactionDetail, businessDetail)
		} else if transactionDetail != nil && transactionDetail.MonthlyInterestID != nil {
			return processInterestCreditTransaction(t, transactionDetail, businessDetail)
		} else {
			return processTransferCreditTransaction(t, transactionDetail, businessDetail)
		}
	} else {
		return processTransferDebitTransaction(t, transactionDetail, businessDetail)
	}
}

func processMoneyRequestCreditTransaction(t transaction.BusinessPostedTransaction,
	transactionDetail *TransactionDetails, businessDetail *BusinessDetails) (*TransactionUpdate, error) {

	u := TransactionUpdate{}

	if transactionDetail == nil {
		return nil, fmt.Errorf("Unable to fetch transaction detail: %s", t.CodeType)
	}

	var contactName string
	if transactionDetail.MoneyRequestType == nil {
		// default request type is invoice card
		requestType := "Card"

		contactName = shared.StringValue(transactionDetail.ContactName)
		u.CounterpartyName = contactName
		u.TransactionTitle = fmt.Sprintf(notification.CreditViaCardTransactionTitle, contactName, businessDetail.BusinessName)
		u.TransactionDescription = fmt.Sprintf(
			notification.CreditViaCardTransactionDescription,
			contactName, businessDetail.BusinessName, requestType)
		u.TransactionSubType = transaction.TransactionSubtypeCardOnlineCredit
	} else {
		switch *transactionDetail.MoneyRequestType {
		case payment.PaymentRequestTypePOS:
			requestType := "Card Reader"

			var city string
			if transactionDetail.PaymentLocation != nil {
				city = transactionDetail.PaymentLocation.City
			}

			u.TransactionTitle = fmt.Sprintf(notification.CreditViaCardReaderTransactionTitle, businessDetail.BusinessName)
			u.TransactionDescription = fmt.Sprintf(notification.CreditViaCardReaderTransactionDescription,
				businessDetail.BusinessName, requestType, city)
			u.TransactionSubType = transaction.TransactionSubtypeCardReaderCredit
		case payment.PaymentRequestTypeInvoiceCard:
			requestType := "Card"

			contactName = shared.StringValue(transactionDetail.ContactName)
			u.CounterpartyName = contactName
			u.TransactionTitle = fmt.Sprintf(notification.CreditViaCardTransactionTitle, contactName, businessDetail.BusinessName)
			u.TransactionDescription = fmt.Sprintf(
				notification.CreditViaCardTransactionDescription,
				contactName, businessDetail.BusinessName, requestType)
			u.TransactionSubType = transaction.TransactionSubtypeCardOnlineCredit
		case payment.PaymentRequestTypeInvoiceBank:
			requestType := "Bank Transfer"

			contactName = shared.StringValue(transactionDetail.ContactName)
			u.CounterpartyName = contactName
			u.TransactionTitle = fmt.Sprintf(notification.CreditViaBankTransactionTitle, contactName, businessDetail.BusinessName)
			u.TransactionDescription = fmt.Sprintf(
				notification.CreditViaBankTransactionDescription,
				contactName, businessDetail.BusinessName, requestType)
			u.TransactionSubType = transaction.TransactionSubtypeBankOnlineCredit
		default:
			log.Println("Unknown money request type", *transactionDetail.MoneyRequestType)
			return nil, errors.New("Unknown money request type")
		}
	}

	return &u, nil
}

func processInterestCreditTransaction(t transaction.BusinessPostedTransaction,
	transactionDetail *TransactionDetails, businessDetail *BusinessDetails) (*TransactionUpdate, error) {
	u := TransactionUpdate{}

	err := transaction.DBRead.Get(
		transactionDetail,
		`SELECT
			start_date "interest_start_date",
			end_date "interest_end_date"
		FROM business_account_monthly_interest 
		WHERE id = $1`,
		*transactionDetail.MonthlyInterestID)
	if err != nil {
		log.Println("Error fetching interest details", err, shared.StringValue(transactionDetail.MonthlyInterestID))
		return nil, err
	}

	date := transactionDetail.InterestStartDate.Time().Format("Jan-06")

	amtFloat, _ := t.Amount.Float64()
	amt := shared.FormatFloatAmount(math.Abs(amtFloat))

	u.TransactionTitle = fmt.Sprintf(notification.CreditInterestTransferTransactionTitle, date)
	u.TransactionDescription = fmt.Sprintf(notification.CreditInterestTransferTransactionDescription,
		businessDetail.BusinessName, amt, date)
	u.TransactionSubType = transaction.TransactionSubtypeInterestTransferCredit

	return &u, nil
}

func processTransferCreditTransaction(t transaction.BusinessPostedTransaction,
	transactionDetail *TransactionDetails, businessDetail *BusinessDetails) (*TransactionUpdate, error) {
	u := TransactionUpdate{}

	amtFloat, _ := t.Amount.Float64()
	amt := shared.FormatFloatAmount(math.Abs(amtFloat))

	senderName := shared.StringValue(transaction.GetOriginAccountHolder(transaction.TransactionTypeTransfer, shared.StringValue(t.BankTransactionDesc)))

	u.CounterpartyName = senderName
	u.TransactionTitle = fmt.Sprintf(notification.CreditWiseTransferTransactionTitle, senderName)
	u.TransactionDescription = fmt.Sprintf(notification.CreditWiseTransferTransactionDescription,
		amt, senderName)
	u.TransactionSubType = transaction.TransactionSubtypeWiseTransferCredit

	return &u, nil
}

func processInstantPayDebitTransaction(t transaction.BusinessPostedTransaction,
	transactionDetail *TransactionDetails, businessDetail *BusinessDetails) (*TransactionUpdate, error) {
	u := TransactionUpdate{}

	contactName := shared.StringValue(transactionDetail.ContactName)
	u.CounterpartyName = contactName
	u.TransactionTitle = fmt.Sprintf(notification.DebitCardInstantPayTransactionTitle, contactName)
	u.TransactionDescription = fmt.Sprintf(
		notification.DebitCardInstantPayTransactionDescription,
		businessDetail.BusinessName, contactName)
	u.TransactionSubType = transaction.TransactionSubtypeCardPushDebit

	return &u, nil
}

func processDepositCreditTransaction(t transaction.BusinessPostedTransaction,
	transactionDetail *TransactionDetails, businessDetail *BusinessDetails) (*TransactionUpdate, error) {
	u := TransactionUpdate{}

	amtFloat, _ := t.Amount.Float64()
	amt := shared.FormatFloatAmount(math.Abs(amtFloat))

	senderName := shared.StringValue(transaction.GetOriginAccountHolder(transaction.TransactionTypeDeposit, shared.StringValue(t.BankTransactionDesc)))

	if senderName != "" {
		u.CounterpartyName = senderName
		u.TransactionTitle = fmt.Sprintf(notification.CreditDepositWireTransactionTitle, senderName)
		u.TransactionDescription = fmt.Sprintf(notification.CreditDepositWireTransactionDescription, amt, senderName)
	} else {
		u.TransactionTitle = fmt.Sprintf(notification.CreditDepositTransactionTitle)
		u.TransactionDescription = fmt.Sprintf(notification.CreditDepositTransactionDescription, amt)
	}

	u.TransactionSubType = transaction.TransactionSubtypeWireTransferCredit

	return &u, nil
}

func processAccountOriginationTransaction(t transaction.BusinessPostedTransaction,
	transactionDetail *TransactionDetails, businessDetail *BusinessDetails) (*TransactionUpdate, error) {
	u := TransactionUpdate{}

	u.TransactionSubType = transaction.TransactionSubtypeAccountOriginationCredit

	return &u, nil
}

func processACHTransaction(t transaction.BusinessPostedTransaction,
	transactionDetail *TransactionDetails, businessDetail *BusinessDetails) (*TransactionUpdate, error) {

	if t.CodeType == transaction.TransactionCodeTypeCreditPosted {
		if t.MoneyRequestID != nil {
			return processMoneyRequestCreditTransaction(t, transactionDetail, businessDetail)
		} else {
			return processACHCreditTransaction(t, transactionDetail, businessDetail)
		}
	} else {
		if t.ContactID != nil {
			return processACHExternalDebitTransaction(t, transactionDetail, businessDetail)
		} else {
			return processACHInternalDebitTransaction(t, transactionDetail, businessDetail)
		}
	}

}

func processACHExternalDebitTransaction(t transaction.BusinessPostedTransaction,
	transactionDetail *TransactionDetails, businessDetail *BusinessDetails) (*TransactionUpdate, error) {
	u := TransactionUpdate{}

	contactName := shared.StringValue(transactionDetail.ContactName)
	u.CounterpartyName = contactName
	u.TransactionTitle = fmt.Sprintf(notification.DebitACHTransactionTitle, contactName)
	u.TransactionDescription = fmt.Sprintf(notification.DebitACHTransactionDescription, businessDetail.BusinessName, contactName)
	u.TransactionSubType = transaction.TransactionSubtypeACHTransferDebit

	return &u, nil
}

func processACHInternalDebitTransaction(t transaction.BusinessPostedTransaction,
	transactionDetail *TransactionDetails, businessDetail *BusinessDetails) (*TransactionUpdate, error) {
	u := TransactionUpdate{}

	receiverName := shared.StringValue(transaction.GetDestinationAccountHolder(transaction.TransactionTypeACH, *t.BankTransactionDesc))

	u.CounterpartyName = receiverName
	u.TransactionTitle = fmt.Sprintf(notification.DebitACHTransactionTitle, receiverName)
	u.TransactionDescription = fmt.Sprintf(notification.DebitACHTransactionDescription, businessDetail.BusinessName, receiverName)
	u.TransactionSubType = transaction.TransactionSubtypeACHTransferDebit

	return &u, nil
}

func processACHCreditTransaction(t transaction.BusinessPostedTransaction,
	transactionDetail *TransactionDetails, businessDetail *BusinessDetails) (*TransactionUpdate, error) {
	u := TransactionUpdate{}

	amtFloat, _ := t.Amount.Float64()
	amt := shared.FormatFloatAmount(math.Abs(amtFloat))

	senderName := shared.StringValue(transaction.GetOriginAccountHolder(transaction.TransactionTypeACH, *t.BankTransactionDesc))

	if senderName != "" {
		u.CounterpartyName = senderName
		u.TransactionTitle = fmt.Sprintf(notification.CreditACHTransactionTitle, senderName)
		u.TransactionDescription = fmt.Sprintf(notification.CreditACHTransactionDescription, amt, senderName)
	} else {
		u.TransactionTitle = fmt.Sprintf(notification.CreditACHTransactionWithoutSenderTitle)
		u.TransactionDescription = fmt.Sprintf(notification.CreditACHTransactionWithoutSenderDescription, amt)
	}

	u.TransactionSubType = transaction.TransactionSubtypeACHTransferCredit

	return &u, nil
}

func processMerchantRefundDebitTransaction(t transaction.BusinessPostedTransaction,
	transactionDetail *TransactionDetails, businessDetail *BusinessDetails) (*TransactionUpdate, error) {
	u := TransactionUpdate{}

	amtFloat, _ := t.Amount.Float64()
	amt := shared.FormatFloatAmount(math.Abs(amtFloat))

	merchantName := strings.TrimSpace(shared.StringValue(t.CardTransaction.MerchantName))
	streetAddress := strings.TrimSpace(shared.StringValue(t.CardTransaction.MerchantStreetAddress))
	if merchantName == "" && streetAddress != "" {
		if strings.Contains(shared.StringValue(t.BankTransactionDesc), streetAddress) {
			merchantName = streetAddress
		}
	}
	if merchantName == "" {
		merchantName = notification.GetMerchantName(shared.StringValue(t.BankTransactionDesc))
	}

	if t.CardTransaction.TransactionType != nil && t.CardTransaction.TransactionType.IsRefundTypeInstantPay() {
		u.TransactionDescription = fmt.Sprintf(notification.CreditCardInstantPayTransactionDescription, amt)
		if merchantName != "" {
			u.CounterpartyName = merchantName
			u.TransactionTitle = fmt.Sprintf(notification.CreditCardInstantPayTransactionTitle, merchantName)
		} else {
			u.TransactionTitle = fmt.Sprintf(notification.CreditCardInstantPayTransactionWithoutSenderTitle)
		}

		u.TransactionSubType = transaction.TransactionSubtypeCardPullCredit
	} else {
		u.TransactionDescription = fmt.Sprintf(notification.CreditMerchantRefundTransactionDescription, amt)
		if merchantName != "" {
			u.CounterpartyName = merchantName
			u.TransactionTitle = fmt.Sprintf(notification.CreditMerchantRefundTransactionTitle, merchantName)
		} else {
			u.TransactionTitle = fmt.Sprintf(notification.CreditMerchantRefundTransactionWithoutMerchantTitle)
		}

		u.TransactionSubType = transaction.TransactionSubtypeMerchantRefundCredit
	}

	return &u, nil
}

func processVisaCreditTransaction(
	t transaction.BusinessPostedTransaction,
	transactionDetail *TransactionDetails,
	businessDetail *BusinessDetails) (*TransactionUpdate, error) {

	u := TransactionUpdate{}
	amtFloat, _ := t.Amount.Float64()
	amt := shared.FormatFloatAmount(math.Abs(amtFloat))
	merchantName := strings.TrimSpace(shared.StringValue(t.CardTransaction.MerchantName))
	streetAddress := strings.TrimSpace(shared.StringValue(t.CardTransaction.MerchantStreetAddress))
	if merchantName == "" && streetAddress != "" {
		if strings.Contains(shared.StringValue(t.BankTransactionDesc), streetAddress) {
			merchantName = streetAddress
		}
	}

	if merchantName == "" {
		merchantName = notification.GetMerchantName(shared.StringValue(t.BankTransactionDesc))
	}

	if t.CodeType != transaction.TransactionCodeTypeCreditPosted {
		return &u, fmt.Errorf("Visa credit code type must be credit posted: %s", t.CodeType)
	}

	u.TransactionDescription = fmt.Sprintf(notification.CreditCardInstantPayTransactionDescription, amt)
	if merchantName != "" {
		u.CounterpartyName = merchantName
		u.TransactionTitle = fmt.Sprintf(notification.CreditCardInstantPayTransactionTitle, merchantName)
	} else {
		u.TransactionTitle = fmt.Sprintf(notification.CreditCardInstantPayTransactionWithoutSenderTitle)
	}

	u.TransactionSubType = transaction.TransactionSubtypeCardPullCredit
	return &u, nil
}

func processATMTransaction(t transaction.BusinessPostedTransaction,
	transactionDetail *TransactionDetails, businessDetail *BusinessDetails) (*TransactionUpdate, error) {
	u := TransactionUpdate{}

	amtFloat, _ := t.Amount.Float64()
	amt := shared.FormatFloatAmount(math.Abs(amtFloat))

	var location string
	if t.CardTransaction != nil {
		if t.CardTransaction.MerchantCity != nil && isLetter(shared.StringValue(t.CardTransaction.MerchantCity)) {
			l := strings.TrimSpace(shared.StringValue(t.CardTransaction.MerchantCity))
			location = l
		}

		if t.CardTransaction.MerchantState != nil && shared.StringValue(t.CardTransaction.MerchantState) != "" {
			var l string
			if location != "" {
				location += ", "
			} else {
				// initialize
				location = l
			}

			location += shared.StringValue(t.CardTransaction.MerchantState)
		}
	}

	maskedCardNumber := shared.StringValue(businessDetail.MaskedCardNumber)
	var lastFour string
	if len(maskedCardNumber) >= 4 {
		lastFour = string(maskedCardNumber[len(maskedCardNumber)-4:])
	}

	u.TransactionTitle = fmt.Sprintf(notification.DebitCardATMTransactionTitle)
	if location != "" {
		u.TransactionDescription = fmt.Sprintf(notification.DebitCardATMTransactionDescription,
			businessDetail.BusinessName, amt, location, lastFour)
	} else {
		u.TransactionDescription = fmt.Sprintf(notification.DebitCardATMTransactionWithoutLocationDescription, businessDetail.BusinessName,
			amt, lastFour)
	}

	u.TransactionSubType = transaction.TransactionSubtypeCardATMDebit

	return &u, nil
}

func processCardPurchaseDebitTransaction(t transaction.BusinessPostedTransaction,
	transactionDetail *TransactionDetails, businessDetail *BusinessDetails) (*TransactionUpdate, error) {
	u := TransactionUpdate{}

	if businessDetail.MaskedCardNumber == nil {
		return nil, errors.New("Unable to get debit card number")
	}

	amtFloat, _ := t.Amount.Float64()
	amt := shared.FormatFloatAmount(math.Abs(amtFloat))

	var location, merchantName, streetAddress string
	if t.CardTransaction != nil {
		if t.CardTransaction.MerchantCity != nil && isLetter(shared.StringValue(t.CardTransaction.MerchantCity)) {
			l := strings.TrimSpace(shared.StringValue(t.CardTransaction.MerchantCity))
			location = l
		}

		if t.CardTransaction.MerchantState != nil && shared.StringValue(t.CardTransaction.MerchantState) != "" {
			var l string
			if location != "" {
				location += ", "
			} else {
				// initialize
				location = l
			}

			location += shared.StringValue(t.CardTransaction.MerchantState)
		}

		merchantName = strings.TrimSpace(shared.StringValue(t.CardTransaction.MerchantName))
		streetAddress = strings.TrimSpace(shared.StringValue(t.CardTransaction.MerchantStreetAddress))

		if merchantName == "" && streetAddress != "" {
			if strings.Contains(shared.StringValue(t.BankTransactionDesc), streetAddress) {
				merchantName = streetAddress
			}
		}

		if merchantName == "" {
			merchantName = notification.GetMerchantName(shared.StringValue(t.BankTransactionDesc))
		}
	}

	maskedCardNumber := shared.StringValue(businessDetail.MaskedCardNumber)
	var lastFour string
	if len(maskedCardNumber) >= 4 {
		lastFour = string(maskedCardNumber[len(maskedCardNumber)-4:])
	}

	if merchantName != "" {
		u.CounterpartyName = merchantName
		u.TransactionTitle = fmt.Sprintf(notification.DebitCardPurchaseTransactionTitle, merchantName)
		if location != "" {
			u.TransactionDescription = fmt.Sprintf(
				notification.DebitCardPurchaseTransactionDescription,
				businessDetail.BusinessName,
				merchantName,
				amt,
				location,
				lastFour,
			)
		} else {
			u.TransactionDescription = fmt.Sprintf(
				notification.DebitCardPurchaseTransactionWithoutLocationDescription,
				businessDetail.BusinessName,
				merchantName,
				amt,
				lastFour,
			)
		}
	} else {
		u.TransactionTitle = fmt.Sprintf(notification.DebitCardPurchaseTransactionWithoutMerchantTitle, amt)
		if location != "" {
			u.TransactionDescription = fmt.Sprintf(
				notification.DebitCardPurchaseTransactionWithoutMerchantDescription,
				businessDetail.BusinessName,
				amt,
				location,
				lastFour)
		} else {
			u.TransactionDescription = fmt.Sprintf(
				notification.DebitCardPurchaseTransactionGenericDescription,
				businessDetail.BusinessName,
				amt,
				lastFour)
		}
	}

	u.TransactionSubType = transaction.TransactionSubtypeCardPurchaseDebit
	return &u, nil
}

func processTransferDebitTransaction(t transaction.BusinessPostedTransaction,
	transactionDetail *TransactionDetails, businessDetail *BusinessDetails) (*TransactionUpdate, error) {
	u := TransactionUpdate{}

	var contactName string
	if transactionDetail != nil {
		contactName = shared.StringValue(transactionDetail.ContactName)
	}

	if contactName == "" {
		contactName = shared.StringValue(transaction.GetDestinationAccountHolder(transaction.TransactionTypeTransfer, shared.StringValue(t.BankTransactionDesc)))
	}

	u.CounterpartyName = contactName
	u.TransactionTitle = fmt.Sprintf(notification.DebitWiseTransferTransactionTitle, contactName)
	u.TransactionDescription = fmt.Sprintf(
		notification.DebitWiseTransferTransactionDescription,
		businessDetail.BusinessName, contactName)

	u.TransactionSubType = transaction.TransactionSubtypeWiseTransferDebit
	return &u, nil
}

func updateTransaction(ID shared.PostedTransactionID, u *TransactionUpdate) error {
	var columns []string

	if len(u.TransactionDescription) > 0 {
		columns = append(columns, "transaction_desc = :transaction_desc")
	}

	if len(u.TransactionTitle) > 0 {
		columns = append(columns, "transaction_title = :transaction_title")
	}

	if len(u.TransactionSubType) > 0 {
		columns = append(columns, "transaction_subtype = :transaction_subtype")
	}

	query := fmt.Sprintf(
		"UPDATE business_transaction SET %s WHERE id = '%s'",
		strings.Join(columns, ", "),
		ID,
	)
	_, err := transaction.DBWrite.NamedExec(
		query, u,
	)
	if err != nil {
		log.Println("Error updating transaction", err, query)
		return err
	}

	return nil
}

func isLetter(s string) bool {

	if len(s) == 0 {
		return false
	}

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
