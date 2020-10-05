package notification

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/wiseco/core-platform/services/invoice"

	grpcTypes "github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/wiseco/core-platform/notification/activity"
	"github.com/wiseco/core-platform/notification/push"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking"
	"github.com/wiseco/core-platform/services/banking/business"
	bus "github.com/wiseco/core-platform/services/business"
	core "github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/services/payment"
	usrsrv "github.com/wiseco/core-platform/services/user"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/core-platform/transaction"
	"github.com/wiseco/go-lib/grpc"
	"github.com/wiseco/go-lib/id"
	"github.com/wiseco/go-lib/locale"
	"github.com/wiseco/go-lib/num"
	grpcBanking "github.com/wiseco/protobuf/golang/banking"
	grpcBankTransfer "github.com/wiseco/protobuf/golang/banking/transfer"
	grpcShopify "github.com/wiseco/protobuf/golang/shopping/shopify"
	grpcTxn "github.com/wiseco/protobuf/golang/transaction"
	grpcBankTxn "github.com/wiseco/protobuf/golang/transaction/bank"
)

func sendPostedToTransactionService(
	n *Notification,
	t *transaction.BusinessPostedTransaction,
	c *transaction.BusinessCardTransaction,
	h *transaction.BusinessHoldTransaction) error {

	// Bank Transaction ID
	var btxID id.BankTransactionID
	if t.ID != "" {
		btxnUUID, err := uuid.Parse(string(t.ID))
		if err != nil {
			return fmt.Errorf("Error parsing account transaction id: %v", err)
		}

		btxID = id.BankTransactionID(btxnUUID)
	}

	sn, err := grpc.GetConnectionStringForService(grpc.ServiceNameTransaction)
	if err != nil {
		return err
	}

	cl, err := grpc.NewInsecureClient(sn)
	if err != nil {
		return err
	}

	defer cl.CloseAndCancel()
	txnClient := grpcBankTxn.NewTransactionServiceClient(cl.GetConn())
	txnResp, err := txnClient.GetByID(context.Background(), &grpcBankTxn.TransactionIDRequest{Id: btxID.String()})
	if err != nil {
		log.Println(err)
		txnResp = &grpcBankTxn.Transaction{}
	}

	// Business ID
	var busID id.BusinessID
	if t.BusinessID != "" {
		bUUID, err := uuid.Parse(string(t.BusinessID))
		if err != nil {
			return fmt.Errorf("Error parsing business id: %v", err)
		}

		busID = id.BusinessID(bUUID)
	}

	// Bank Account ID
	var accID id.BankAccountID
	if t.AccountID != nil && *t.AccountID != "" {
		accUUID, err := uuid.Parse(shared.StringValue(t.AccountID))
		if err != nil {
			return fmt.Errorf("Error parsing bank account id: %v", err)
		}

		accID = id.BankAccountID(accUUID)
	}

	// Account must always exist
	acc, err := business.NewAccountService().GetByIDInternal(accID.UUIDString())
	if err != nil {
		return fmt.Errorf("Error getting bank account id: %v", err)
	}

	// Debit Card ID
	var dbcID id.DebitCardID
	if t.CardID != nil && *t.CardID != "" {
		dbcUUID, err := uuid.Parse(shared.StringValue(t.CardID))
		if err != nil {
			return fmt.Errorf("Error parsing debit card id: %v", err)
		}

		dbcID = id.DebitCardID(dbcUUID)
	}

	// Money Transfer ID
	var btID id.BankTransferID
	if t.MoneyTransferID != nil && *t.MoneyTransferID != "" {
		btUUID, err := uuid.Parse(shared.StringValue(t.MoneyTransferID))
		if err != nil {
			return fmt.Errorf("Error parsing bank transfer id: %v", err)
		}

		btID = id.BankTransferID(btUUID)
	}

	var interestDate shared.Date
	var mt *business.MoneyTransfer
	if !btID.IsZero() {
		mt, err = business.NewMoneyTransferServiceWithout().GetByIDOnlyInternal(btID.UUIDString())
		if err != nil {
			return fmt.Errorf("Error getting bank transfer: %v", err)
		}

		if mt.MonthlyInterestID != nil {
			mi := &transaction.BusinessAccountMonthlyInterest{}
			err := transaction.DBRead.Get(
				mi,
				"SELECT * FROM business_account_monthly_interest WHERE id = $1",
				shared.StringValue(mt.MonthlyInterestID),
			)
			if err != nil {
				log.Println("Error fetching interest details", err, shared.StringValue(mt.MonthlyInterestID))
				return fmt.Errorf("Error fetching interest details %v %s", err, shared.StringValue(mt.MonthlyInterestID))
			}

			interestDate = mi.StartDate
		}
	}

	var grpcDate string
	if !interestDate.IsZero() {
		d := locale.Date(interestDate)
		grpcDate = d.Format()
	}

	// Payment Request ID
	var prID id.PaymentRequestID
	if t.MoneyRequestID != nil && *t.MoneyRequestID != "" {
		prUUID, err := uuid.Parse(string(*t.MoneyRequestID))
		if err != nil {
			return fmt.Errorf("Error parsing payment request id: %v", err)
		}

		prID = id.PaymentRequestID(prUUID)
	}

	// Contact
	var contactID id.ContactID
	if t.ContactID != nil && *t.ContactID != "" {
		contactUUID, err := uuid.Parse(shared.StringValue(t.ContactID))
		if err != nil {
			return fmt.Errorf("Error parsing contact id: %v", err)
		}

		contactID = id.ContactID(contactUUID)
	}

	// Notification event
	evUUID, err := uuid.Parse(n.ID)
	if err != nil {
		return fmt.Errorf("Error parsing event id: %v", err)
	}

	evID := id.EventID(evUUID)

	// TODO: Get Thread ID from Bank transfer if none available
	var evThreadID id.EventThreadID
	if txnResp.EventThreadId != "" {
		evThreadID, err = id.ParseEventThreadID(txnResp.EventThreadId)
		if err != nil {
			return fmt.Errorf("Error in event thread id: %v", err)
		}
	} else {
		evThreadID, _ = id.NewEventThreadID()
		if err != nil {
			return fmt.Errorf("Error in event thread id: %v", err)
		}
	}

	status, ok := transaction.TransactionStatusToProto[t.Status]
	if !ok {
		return errors.New("invalid transaction status")
	}

	category, ok := transaction.TransactionTypeToCategoryProto[t.TransactionType]
	if !ok {
		return errors.New("invalid transaction category type")
	}

	counterpartyType := grpcTxn.BankTransactionCounterpartyType_BTCT_UNSPECIFIED
	txnType, _ := transaction.TransactionSubtypeToTypeProto[transaction.TransactionSubtypeUnspecified]
	legacySubtype := transaction.TransactionSubtypeUnspecified
	if t.TransactionSubtype != nil {
		legacySubtype = *t.TransactionSubtype

		txnType, ok = transaction.TransactionSubtypeToTypeProto[*t.TransactionSubtype]
		if !ok {
			return errors.New("invalid transaction subtype")
		}

		switch txnType {
		case grpcTxn.BankTransactionType_BTT_UNSPECIFIED:
			if acc.UsageType == business.UsageTypeClearing {
				if t.Amount.IsNegative() {
					txnType = grpcTxn.BankTransactionType_BTT_CLEARING_DEBIT
				} else {
					if mt != nil && mt.MonthlyInterestID != nil {
						txnType = grpcTxn.BankTransactionType_BTT_INTEREST_CREDIT
					} else {
						txnType = grpcTxn.BankTransactionType_BTT_CLEARING_CREDIT
					}
				}

				if acc.Id == os.Getenv("WISE_CLEARING_ACCOUNT_ID") {
					if mt != nil && mt.MonthlyInterestID != nil {
						counterpartyType = grpcTxn.BankTransactionCounterpartyType_BTCT_INTEREST_PAYOUT_DEBIT
					} else {
						counterpartyType = grpcTxn.BankTransactionCounterpartyType_BTCT_CARD_PAYMENT_DEBIT
					}
				} else if acc.Id == os.Getenv("WISE_PROMO_CLEARING_ACCOUNT_ID") {
					counterpartyType = grpcTxn.BankTransactionCounterpartyType_BTCT_PROMO_DEBIT
				}
			}
		case grpcTxn.BankTransactionType_BTT_INTRABANK_ONLINE_CREDIT:
			if category == grpcTxn.BankTransactionCategory_BTC_ACH {
				txnType = grpcTxn.BankTransactionType_BTT_ACH_ONLINE_CREDIT
			}
		case grpcTxn.BankTransactionType_BTT_ACH_CREDIT:
			if *t.TransactionSubtype == transaction.TransactionSubtypeACHTransferShopifyCredit {
				counterpartyType = grpcTxn.BankTransactionCounterpartyType_BTCT_SHOPPING_SHOPIFY_PAYOUT
			}
		case grpcTxn.BankTransactionType_BTT_ACH_DEBIT:
			if *t.TransactionSubtype == transaction.TransactionSubtypeACHTransferShopifyDebit {
				counterpartyType = grpcTxn.BankTransactionCounterpartyType_BTCT_SHOPPING_SHOPIFY_REFUND
			}
		}
	}

	// If unspecified default to other debit or credit
	if txnType == grpcTxn.BankTransactionType_BTT_UNSPECIFIED {
		if t.Amount.IsNegative() {
			txnType = grpcTxn.BankTransactionType_BTT_OTHER_DEBIT
		} else {
			txnType = grpcTxn.BankTransactionType_BTT_OTHER_CREDIT
		}
	}

	txnDate, err := grpcTypes.TimestampProto(t.TransactionDate)
	if err != nil {
		return err
	}

	req := &grpcBankTxn.UpsertTransactionRequest{
		Id:                     btxID.String(),
		BusinessId:             busID.String(),
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
		InterestDate:           grpcDate,
		Counterparty:           strings.TrimSpace(t.Counterparty),
		CounterpartyType:       counterpartyType,
		LegacyType:             string(t.TransactionType),
		LegacyCodeType:         string(t.CodeType),
		LegacySubtype:          string(legacySubtype),
		LegacyTitle:            shared.StringValue(t.TransactionTitle),
		LegacyDescription:      t.TransactionDesc,
		LegacyNotes:            shared.StringValue(t.SourceNotes),
	}

	if t.NotificationID != nil {
		req.NotificationId = *t.NotificationID
	}

	if c != nil {
		network, ok := transaction.CardNetworkToProto[strings.ToUpper(c.TransactionNetwork)]
		if !ok {
			return fmt.Errorf("invalid card network: %s", strings.ToUpper(c.TransactionNetwork))
		}

		usr, err := usrsrv.NewUserServiceWithout().GetByIdInternal(c.CardHolderID)
		if err != nil {
			return err
		}

		cardHolderUUID, _ := id.ParseUUID(string(usr.ConsumerID))
		authDate, _ := grpcTypes.TimestampProto(c.AuthDate)
		localDate, _ := grpcTypes.TimestampProto(c.LocalDate)
		created, _ := grpcTypes.TimestampProto(c.Created)
		req.CardRequest = &grpcBankTxn.UpsertCardTransactionRequest{
			CardHolderId:          id.ConsumerID(cardHolderUUID).String(),
			NetworkTransactionId:  c.CardTransactionID,
			Network:               network,
			AuthAmount:            c.AuthAmount.FormatCurrency(),
			AuthDate:              authDate,
			AuthResponseCode:      c.AuthResponseCode,
			AuthNumber:            c.AuthNumber,
			CardTransactionType:   string(c.TransactionType),
			LocalAmount:           c.LocalAmount.FormatCurrency(),
			LocalCurrency:         c.LocalCurrency,
			LocalDate:             localDate,
			BillingCurrency:       c.BillingCurrency,
			PosEntryMode:          c.POSEntryMode,
			PosConditionCode:      c.POSConditionCode,
			AcquirerBin:           c.AcquirerBIN,
			MerchantName:          c.MerchantName,
			MerchantCategoryCode:  c.MerchantCategoryCode,
			AcceptorId:            c.MerchantID,
			AcceptorTerminal:      c.MerchantTerminal,
			AcceptorStreetAddress: c.MerchantStreetAddress,
			AcceptorCity:          c.MerchantCity,
			AcceptorState:         c.MerchantState,
			AcceptorCountry:       c.MerchantCountry,
			Created:               created,
		}
	}

	if h != nil {
		txnDate, _ := grpcTypes.TimestampProto(h.Date)
		expiryDate, _ := grpcTypes.TimestampProto(h.ExpiryDate)
		created, _ := grpcTypes.TimestampProto(h.Created)
		req.HoldRequest = &grpcBankTxn.UpsertHoldTransactionRequest{
			HoldNumber:      h.Number,
			Amount:          h.Amount.FormatCurrency(),
			TransactionDate: txnDate,
			ExpiryDate:      expiryDate,
			Created:         created,
		}
	}

	resp, err := txnClient.Upsert(context.Background(), req)
	if err != nil {
		log.Println("Transaction service error:", err)
	} else {
		log.Println("Transaction Posted service success:", resp.Transaction.Id)
	}

	return err
}

func createPostedTransaction(
	n *Notification,
	t transaction.BusinessPostedTransactionCreate,
	c *transaction.BusinessCardTransactionCreate,
	h *transaction.BusinessHoldTransactionCreate) (*transaction.BusinessPostedTransaction, error) {

	log.Println("createPostedTransaction")

	srv := transaction.NewCardService()

	// TODO: Use a transaction or pass all three values
	trx, err := transaction.NewBusinessService().Create(t)
	if err != nil {
		return nil, fmt.Errorf("creating posted transaction error: %v", err)
	}

	trx.Counterparty = t.Counterparty
	trx.Status = t.Status

	var card *transaction.BusinessCardTransaction
	if c != nil {
		c.TransactionID = string(trx.ID)
		crd, err := srv.Create(c)
		if err == nil {
			card = &crd
		}
	}

	var hold *transaction.BusinessHoldTransaction
	if h != nil {
		h.TransactionID = string(trx.ID)
		hld, err := srv.CreateHold(h)
		if err == nil {
			hold = &hld
		}
	}

	if os.Getenv("USE_TRANSACTION_SERVICE") == "true" {
		err = sendPostedToTransactionService(n, trx, card, hold)
		if err != nil {
			log.Println("Transaction service response error:", err)
		}
	}

	return transaction.NewBusinessService().GetByIDInternal(trx.ID)
}

func sendPendingToTransactionService(
	n *Notification,
	t *transaction.BusinessPendingTransaction,
	c *transaction.BusinessCardTransaction,
	h *transaction.BusinessHoldTransaction) error {

	// Bank Transaction ID
	var btxID id.BankTransactionID
	if t.ID != "" {
		btxnUUID, err := uuid.Parse(string(t.ID))
		if err != nil {
			return fmt.Errorf("Error parsing account transaction id: %v", err)
		}

		btxID = id.BankTransactionID(btxnUUID)
	}

	sn, err := grpc.GetConnectionStringForService(grpc.ServiceNameTransaction)
	if err != nil {
		return err
	}

	cl, err := grpc.NewInsecureClient(sn)
	if err != nil {
		return err
	}

	defer cl.CloseAndCancel()
	txnClient := grpcBankTxn.NewTransactionServiceClient(cl.GetConn())
	txnResp, err := txnClient.GetByID(context.Background(), &grpcBankTxn.TransactionIDRequest{Id: btxID.String()})
	if err != nil {
		//TODO ERROR TYPES!
		if strings.Contains(err.Error(), "bank transaction not found") {
			err = nil
			txnResp = &grpcBankTxn.Transaction{}
		} else {
			return err
		}
	}

	// Business ID
	var busID id.BusinessID
	if t.BusinessID != "" {
		bUUID, err := uuid.Parse(string(t.BusinessID))
		if err != nil {
			return fmt.Errorf("Error parsing business id: %v", err)
		}
		busID = id.BusinessID(bUUID)
	}

	// Bank Account ID
	var accID id.BankAccountID
	if t.AccountID != nil && *t.AccountID != "" {
		accUUID, err := uuid.Parse(shared.StringValue(t.AccountID))
		if err != nil {
			return fmt.Errorf("Error parsing bank account id: %v", err)
		}
		accID = id.BankAccountID(accUUID)
	}

	// Debit Card ID
	var dbcID id.DebitCardID
	if t.CardID != nil && *t.CardID != "" {
		dbcUUID, err := uuid.Parse(shared.StringValue(t.CardID))
		if err != nil {
			return fmt.Errorf("Error parsing debit card id: %v", err)
		}
		dbcID = id.DebitCardID(dbcUUID)
	}

	// Money Transfer ID
	var btID id.BankTransferID
	if t.MoneyTransferID != nil && *t.MoneyTransferID != "" {
		btUUID, err := uuid.Parse(shared.StringValue(t.MoneyTransferID))
		if err != nil {
			return fmt.Errorf("Error parsing bank transfer id: %v", err)
		}
		btID = id.BankTransferID(btUUID)
	}

	// Payment Request ID
	var prID id.PaymentRequestID
	if t.MoneyRequestID != nil && *t.MoneyRequestID != "" {
		prUUID, err := uuid.Parse(string(*t.MoneyRequestID))
		if err != nil {
			return fmt.Errorf("Error parsing payment request id: %v", err)
		}
		prID = id.PaymentRequestID(prUUID)
	}

	// Contact
	var contactID id.ContactID
	if t.ContactID != nil && *t.ContactID != "" {
		contactUUID, err := uuid.Parse(shared.StringValue(t.ContactID))
		if err != nil {
			return fmt.Errorf("Error parsing contact id: %v", err)
		}
		contactID = id.ContactID(contactUUID)
	}

	// Notification event
	evUUID, err := uuid.Parse(n.ID)
	if err != nil {
		return fmt.Errorf("Error parsing event id: %v", err)
	}

	evID := id.EventID(evUUID)
	status, ok := transaction.TransactionStatusToProto[t.Status]
	if !ok {
		return errors.New("invalid transaction status")
	}

	// TODO: Get Thread ID from bank transfer if none available
	var evThreadID id.EventThreadID
	if txnResp.EventThreadId != "" {
		evThreadID, err = id.ParseEventThreadID(txnResp.EventThreadId)
		if err != nil {
			return fmt.Errorf("Error in event thread id: %v", err)
		}
	} else {
		evThreadID, _ = id.NewEventThreadID()
		if err != nil {
			return fmt.Errorf("Error in event thread id: %v", err)
		}
	}

	category, ok := transaction.TransactionTypeToCategoryProto[t.TransactionType]
	if !ok {
		return errors.New("invalid transaction category type")
	}

	legacySubtype := transaction.TransactionSubtypeUnspecified
	txnType, _ := transaction.TransactionSubtypeToTypeProto[transaction.TransactionSubtypeUnspecified]
	if t.TransactionSubtype != nil {
		legacySubtype = *t.TransactionSubtype
		txnType, ok = transaction.TransactionSubtypeToTypeProto[*t.TransactionSubtype]
		if !ok {
			return errors.New("invalid transaction subtype")
		}
	}

	// If unspecified default to other debit or credit
	if txnType == grpcTxn.BankTransactionType_BTT_UNSPECIFIED {
		if t.Amount.IsNegative() {
			txnType = grpcTxn.BankTransactionType_BTT_OTHER_DEBIT
		} else {
			txnType = grpcTxn.BankTransactionType_BTT_OTHER_CREDIT
		}
	}

	txnDate, err := grpcTypes.TimestampProto(t.TransactionDate)
	if err != nil {
		return err
	}

	req := &grpcBankTxn.UpsertTransactionRequest{
		Id:                     btxID.String(),
		BusinessId:             busID.String(),
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
		Counterparty:           strings.TrimSpace(t.Counterparty),
		LegacyType:             string(t.TransactionType),
		LegacyCodeType:         string(t.CodeType),
		LegacySubtype:          string(legacySubtype),
		LegacyTitle:            shared.StringValue(t.TransactionTitle),
		LegacyDescription:      t.TransactionDesc,
		LegacyNotes:            shared.StringValue(t.SourceNotes),
	}

	if t.NotificationID != nil {
		req.NotificationId = *t.NotificationID
	}

	if c != nil {
		network, ok := transaction.CardNetworkToProto[strings.ToUpper(c.TransactionNetwork)]
		if !ok {
			return fmt.Errorf("invalid card network: %s", strings.ToUpper(c.TransactionNetwork))
		}
		usr, err := usrsrv.NewUserServiceWithout().GetByIdInternal(c.CardHolderID)
		if err != nil {
			return err
		}
		cardHolderUUID, _ := id.ParseUUID(string(usr.ConsumerID))
		authDate, _ := grpcTypes.TimestampProto(c.AuthDate)
		localDate, _ := grpcTypes.TimestampProto(c.LocalDate)
		created, _ := grpcTypes.TimestampProto(c.Created)
		req.CardRequest = &grpcBankTxn.UpsertCardTransactionRequest{
			CardHolderId:          id.ConsumerID(cardHolderUUID).String(),
			NetworkTransactionId:  c.CardTransactionID,
			Network:               network,
			AuthAmount:            c.AuthAmount.FormatCurrency(),
			AuthDate:              authDate,
			AuthResponseCode:      c.AuthResponseCode,
			AuthNumber:            c.AuthNumber,
			CardTransactionType:   string(c.TransactionType),
			LocalAmount:           c.LocalAmount.FormatCurrency(),
			LocalCurrency:         c.LocalCurrency,
			LocalDate:             localDate,
			BillingCurrency:       c.BillingCurrency,
			PosEntryMode:          c.POSEntryMode,
			PosConditionCode:      c.POSConditionCode,
			AcquirerBin:           c.AcquirerBIN,
			MerchantName:          c.MerchantName,
			MerchantCategoryCode:  c.MerchantCategoryCode,
			AcceptorId:            c.MerchantID,
			AcceptorTerminal:      c.MerchantTerminal,
			AcceptorStreetAddress: c.MerchantStreetAddress,
			AcceptorCity:          c.MerchantCity,
			AcceptorState:         c.MerchantState,
			AcceptorCountry:       c.MerchantCountry,
			Created:               created,
		}
	}

	if h != nil {
		txnDate, _ := grpcTypes.TimestampProto(h.Date)
		expiryDate, _ := grpcTypes.TimestampProto(h.ExpiryDate)
		created, _ := grpcTypes.TimestampProto(h.Created)
		req.HoldRequest = &grpcBankTxn.UpsertHoldTransactionRequest{
			HoldNumber:      h.Number,
			Amount:          h.Amount.FormatCurrency(),
			TransactionDate: txnDate,
			ExpiryDate:      expiryDate,
			Created:         created,
		}
	}

	resp, err := txnClient.Upsert(context.Background(), req)
	if err != nil {
		log.Println("Transaction service error:", err)
	} else {
		log.Println("Transaction Pending service success:", resp.Transaction.Id)
	}

	return err
}

func createPendingTransaction(
	n *Notification,
	t transaction.BusinessPendingTransactionCreate,
	c *transaction.BusinessCardTransactionCreate,
	h *transaction.BusinessHoldTransactionCreate) (*transaction.BusinessPendingTransaction, error) {

	srv := transaction.NewPendingCardService()
	var trx *transaction.BusinessPendingTransaction
	var err error

	// TODO: Use a transaction or pass all three values
	switch t.CodeType {
	case transaction.TransactionCodeTypeAuthReversed:
		fallthrough
	case transaction.TransactionCodeTypeHoldReleased:
		trx = &transaction.BusinessPendingTransaction{
			ID:                  t.ID,
			BusinessID:          t.BusinessID,
			BankName:            t.BankName,
			BankTransactionID:   t.BankTransactionID,
			BankExtra:           t.BankExtra,
			TransactionType:     t.TransactionType,
			CardID:              t.CardID,
			CodeType:            t.CodeType,
			Amount:              t.Amount,
			Currency:            t.Currency,
			MoneyTransferID:     t.MoneyTransferID,
			SourceNotes:         t.SourceNotes,
			ContactID:           t.ContactID,
			BankTransactionDesc: t.BankTransactionDesc,
			TransactionDesc:     t.TransactionDesc,
			TransactionDate:     t.TransactionDate,
			TransactionStatus:   t.TransactionStatus,
			PartnerName:         t.PartnerName,
			TransactionTitle:    &t.TransactionTitle,
			TransactionSubtype:  &t.TransactionSubtype,
			MoneyRequestID:      t.MoneyRequestID,
		}

		if t.AccountID != nil && strings.HasPrefix(*t.AccountID, id.IDPrefixBankAccount.String()) {
			aID, err := id.ParseBankAccountID(*t.AccountID)
			if err != nil {
				return nil, err
			}

			s := aID.UUIDString()
			trx.AccountID = &s
		}

		break
	default:
		trx, err = transaction.NewPendingTransactionService().CreateTransaction(t)
		if err != nil {
			return nil, fmt.Errorf("creating pending transaction error: %v", err)
		}
	}

	trx.Status = t.Status
	trx.Counterparty = t.Counterparty

	var card *transaction.BusinessCardTransaction
	if c != nil {
		c.TransactionID = string(trx.ID)
		crd, err := srv.Create(c)
		if err == nil {
			card = &crd
		}
	}

	var hold *transaction.BusinessHoldTransaction
	if h != nil {
		h.TransactionID = string(trx.ID)
		hld, err := srv.CreateHold(h)
		if err == nil {
			hold = &hld
		}
	}

	if os.Getenv("USE_TRANSACTION_SERVICE") == "true" {
		err = sendPendingToTransactionService(n, trx, card, hold)
		if err != nil {
			log.Println("Transaction service response error:", err)
		}
	}

	return transaction.NewPendingTransactionService().GetByIDInternal(trx.ID, t.BusinessID)
}

func processTransactionNotification(n Notification) error {
	log.Println("Notification ID is ", n.ID)

	var t TransactionNotification
	err := json.Unmarshal(n.Data, &t)
	if err != nil {
		return fmt.Errorf("transaction notification error: %v", err)
	}

	l := fmt.Sprintf("Transaction Notification Code Type %s Transaction Type %s", string(t.CodeType), string(t.TransactionType))
	log.Println(l)

	// Get user ID
	usrID, err := NewNotificationService().getUserIDByEntity(n.EntityID, n.EntityType)
	if err != nil {
		log.Println("User ID not found")
		return err
	}

	// Get business ID
	busID, err := NewNotificationService().getBusinessIDByEntity(n.EntityID, n.EntityType, t.BankAccountID)
	if err != nil {
		log.Println("Business ID not found")
		return err
	}

	createActivity := true
	sendPushNotification := true
	var m *TransactionMessage
	var activityID *string
	var accountID *string
	var account *business.BankAccount
	if t.BankAccountID != nil {
		account, err = business.NewAccountService().GetByBankAccountId(*t.BankAccountID, *busID)
		if err == nil {
			accountID = &account.Id
		}

		if account != nil {
			updateAccountBalance(account)

			// Don't create activity or send push when not primary
			if account.UsageType != business.UsageTypePrimary {
				createActivity = false
				sendPushNotification = false
			}
		}
	}

	var cardID *string
	if t.BankCardID != nil {
		card, err := business.NewCardService(services.NewSourceRequest()).GetByBankCardId(*t.BankCardID, *usrID)
		if err == nil {
			cardID = &card.Id
		}
	}

	var transfer *TransferDetails
	var moneyTransferID *string
	var moneyRequestID *shared.PaymentRequestID
	var contactID *string
	var notes *string
	var monthlyInterestID *string
	if t.BankMoneyTransferID != nil && *t.BankMoneyTransferID != "" {
		mt := TransferDetails{}

		tr, err := business.NewMoneyTransferServiceWithout().GetByBankIDInternal(banking.BankName(t.BankName), *t.BankMoneyTransferID)
		if err != nil {
			log.Println("MoneyTransfer:", err)
		} else {

			if tr.Status == banking.MoneyTransferStatusBankError {
				return fmt.Errorf("Bank error: %s", tr.ErrorCause)
			}

			mt.MoneyTransferID = &tr.Id
			mt.TransferContactID = tr.ContactId
			mt.RequestID = tr.MoneyRequestID
			mt.Notes = tr.Notes
			mt.MonthlyInterestID = tr.MonthlyInterestID

			moneyTransferID = mt.MoneyTransferID

			notes = mt.Notes
			monthlyInterestID = mt.MonthlyInterestID
			transfer = &mt
		}

		if mt.TransferContactID != nil && *mt.TransferContactID != "" {
			query := `
				SELECT
					first_name "business_contact.first_name",
					last_name "business_contact.last_name",
					business_name "business_contact.business_name",
					email "business_contact.email",
					phone_number "business_contact.phone_number"
				FROM business_contact
				WHERE id = $1`
			err = core.DBRead.Get(&mt, query, *mt.TransferContactID)
			if err != nil {
				log.Println("Contact:", err)
			}
		}

		query := `
			SELECT
				owner_id "business.owner_id",
				legal_name "business.legal_name",
				dba "business.dba",
				phone "business.phone",
				email "business.email"
			FROM business
			WHERE id = $1`

		err = core.DBRead.Get(&mt, query, busID)
		if err != nil {
			log.Println("Business:", err)
		}

		if mt.RequestID != nil && *mt.RequestID != "" {
			isPOS := isPOSRequest(mt.RequestID)
			if os.Getenv("USE_INVOICE_SERVICE") == "true" && !isPOS {
				err := fillTransferDetailsFromInvoiceSvc(&mt)
				if err != nil {
					log.Println("MoneyRequest (Invoice conversion error):", err)
				} else {
					moneyRequestID = mt.RequestID
					if mt.LinkedBankAccountID != nil && *mt.LinkedBankAccountID != "" {
						la, err := business.NewLinkedAccountServiceWithout().GetByIDInternal(*mt.LinkedBankAccountID)
						if err != nil {
							log.Println("LinkedBankAccount:", err)
						} else {
							mt.BankName = la.BankName
							accNum := string(la.AccountNumber)
							mt.AccountNumber = &accNum
						}
					}
				}
			} else {
				query := `
				SELECT
					business_money_request.contact_id "business_money_request.contact_id",
					business_money_request.request_type "business_money_request.request_type",
					business_money_request.request_source "business_money_request.request_source",
					business_invoice.id "business_invoice.id",
					business_invoice.invoice_number "business_invoice.invoice_number",
					business_money_request_payment.id "business_money_request_payment.id",
					business_money_request_payment.linked_bank_account_id "business_money_request_payment.linked_bank_account_id",
					business_money_request_payment.payment_date "business_money_request_payment.payment_date"
				FROM business_money_request
				LEFT JOIN business_money_request_payment ON business_money_request_payment.request_id = business_money_request.id
				LEFT JOIN business_invoice ON business_invoice.request_id = business_money_request.id
				WHERE business_money_request.id = $1`

				err = core.DBRead.Get(&mt, query, *mt.RequestID)
				if err != nil {
					log.Println("MoneyRequest:", err)
				} else {
					moneyRequestID = mt.RequestID
					if mt.LinkedBankAccountID != nil && *mt.LinkedBankAccountID != "" {
						la, err := business.NewLinkedAccountServiceWithout().GetByIDInternal(*mt.LinkedBankAccountID)
						if err != nil {
							log.Println("LinkedBankAccount:", err)
						} else {
							mt.BankName = la.BankName
							accNum := string(la.AccountNumber)
							mt.AccountNumber = &accNum
						}
					}
				}
			}
		}

		// In case of money request, do not get clearing account contact details
		if account.UsageType != business.UsageTypeClearing {
			if mt.RequestContactID != nil && *mt.RequestContactID != "" {
				contactID = mt.RequestContactID
			} else if mt.TransferContactID != nil && *mt.TransferContactID != "" {
				contactID = mt.TransferContactID
			}
		}
	}

	// Handle older notification styles
	desc := t.MoneyTransferDesc
	if t.BankTransactionDesc != nil {
		desc = t.BankTransactionDesc
	}

	txnID := uuid.New().String()

	isSnapcheck, err := isSnapcheckTransaction(t.BankTransactionDesc)
	if err != nil {
		return err
	}

	if isSnapcheck {
		bts, err := business.NewBankingTransferService()
		if err != nil {
			return err
		}

		if accountID != nil {
			pts := []grpcBankTransfer.PartnerTransferStatus{
				grpcBankTransfer.PartnerTransferStatus_TPS_SNAPCHECK_DEPOSITED,
			}

			mt, err := bts.GetByAccountIDPartnerTransferStatusAndAmount(*accountID, pts, t.Amount)
			if err != nil {
				return err
			}

			moneyTransferID = &mt.Id
		} else {
			return errors.New("Snapcheck notification sent without BankAccountID")
		}

	}

	// Store transaction
	postedTXNCreate := transaction.BusinessPostedTransactionCreate{
		// Generate Transaction ID for every transaction
		ID:                  shared.PostedTransactionID(txnID),
		BankTransactionID:   &t.BankTransactionID,
		BusinessID:          *busID,
		BankName:            t.BankName,
		TransactionType:     transaction.TransactionType(t.TransactionType),
		AccountID:           accountID,
		CardID:              cardID,
		CodeType:            transaction.TransactionCodeType(t.CodeType),
		Amount:              t.Amount,
		Currency:            transaction.Currency(t.Currency),
		MoneyTransferID:     moneyTransferID,
		ContactID:           contactID,
		TransactionDate:     t.TransactionDate,
		BankTransactionDesc: desc,
		MoneyRequestID:      moneyRequestID,
		SourceNotes:         notes,
		NotificationID:      &n.ID,
	}

	pendingTXNCreate := transaction.BusinessPendingTransactionCreate{
		// Generate Transaction ID for every transaction
		ID:                  shared.PendingTransactionID(txnID),
		BankTransactionID:   &t.BankTransactionID,
		BusinessID:          *busID,
		BankName:            t.BankName,
		TransactionType:     transaction.TransactionType(t.TransactionType),
		AccountID:           accountID,
		CardID:              cardID,
		CodeType:            transaction.TransactionCodeType(t.CodeType),
		Amount:              t.Amount,
		Currency:            transaction.Currency(t.Currency),
		MoneyTransferID:     moneyTransferID,
		ContactID:           contactID,
		TransactionDate:     t.TransactionDate,
		BankTransactionDesc: t.BankTransactionDesc,
		MoneyRequestID:      moneyRequestID,
		SourceNotes:         notes,
		NotificationID:      &n.ID,
	}

	var cardTXN *transaction.BusinessCardTransactionCreate
	if t.CardTransaction != nil {

		authAmount, err := num.NewFromFloat(t.CardTransaction.AuthAmount)
		if err != nil {
			return err
		}

		localAmount, err := num.NewFromFloat(t.CardTransaction.LocalAmount)
		if err != nil {
			return err
		}

		c := transaction.BusinessCardTransactionCreate{
			CardHolderID:          *usrID,
			CardTransactionID:     t.CardTransaction.CardTransactionID, // Can card transaction ID be null?
			TransactionNetwork:    t.CardTransaction.TransactionNetwork,
			AuthAmount:            authAmount,
			AuthDate:              t.CardTransaction.AuthDate,
			AuthResponseCode:      string(t.CardTransaction.AuthResponseCode),
			AuthNumber:            t.CardTransaction.AuthNumber,
			TransactionType:       transaction.CardTransactionType(t.CardTransaction.TransactionType),
			LocalAmount:           localAmount,
			LocalCurrency:         t.CardTransaction.LocalCurrency,
			LocalDate:             t.CardTransaction.LocalDate,
			BillingCurrency:       t.CardTransaction.BillingCurrency,
			POSEntryMode:          string(t.CardTransaction.POSEntryMode),
			POSConditionCode:      t.CardTransaction.POSConditionCode,
			AcquirerBIN:           t.CardTransaction.AcquirerBIN,
			MerchantID:            t.CardTransaction.MerchantID,
			MerchantCategoryCode:  t.CardTransaction.MerchantCategoryCode,
			MerchantTerminal:      t.CardTransaction.MerchantTerminal,
			MerchantName:          t.CardTransaction.MerchantName,
			MerchantStreetAddress: t.CardTransaction.MerchantStreetAddress,
			MerchantCity:          t.CardTransaction.MerchantCity,
			MerchantState:         t.CardTransaction.MerchantState,
			MerchantCountry:       t.CardTransaction.MerchantCountry,
		}

		cardTXN = &c

		// Process card activity
		if createActivity {
			m, err = processCardTransaction(*usrID, *busID, postedTXNCreate.AccountID, cardID, t.Amount, t, *cardTXN)
			if err != nil {
				return err
			}
		}
	}

	var holdTXN *transaction.BusinessHoldTransactionCreate
	if t.HoldTransaction != nil {
		amount, err := num.NewFromFloat(t.HoldTransaction.Amount)
		if err != nil {
			return err
		}

		h := transaction.BusinessHoldTransactionCreate{
			Number:     t.HoldTransaction.Number,
			Amount:     amount,
			Date:       t.HoldTransaction.Date,
			ExpiryDate: t.HoldTransaction.ExpiryDate,
		}

		holdTXN = &h
	}

	if createActivity {
		switch t.TransactionType {
		case TransactionTypeTransfer:
			if t.BankMoneyTransferID != nil {
				m, err = processMoneyTransferTransaction(*usrID, *busID, postedTXNCreate.AccountID, contactID,
					moneyTransferID, postedTXNCreate.MoneyRequestID, monthlyInterestID, t.Amount, string(n.Action), t)
				if err != nil {
					return err
				}
			}
		case TransactionTypePurchase:
			if moneyTransferID == nil {
				break
			}
			fallthrough
		case TransactionTypeACH:
			m, err = processACHTransaction(*usrID, *busID, accountID, contactID, moneyTransferID, postedTXNCreate.MoneyRequestID, t.Amount, string(n.Action), t)
			if err != nil {
				return err
			}
		case TransactionTypeDeposit:
			// Account origination hack
			if t.Amount.Zero == true {
				m, err = processAccountOrigination(*usrID, *busID, t)
				if err != nil {
					log.Println("error processing account origination", err)
					return err
				}

				// Dont send push for account origination
				sendPushNotification = false
			} else {
				m, err = processMoneyTransferTransaction(*usrID, *busID, postedTXNCreate.AccountID, contactID, moneyTransferID, postedTXNCreate.MoneyRequestID,
					monthlyInterestID, t.Amount, string(n.Action), t)
				if err != nil {
					return err
				}
			}
		case TransactionTypeVisaCredit:
			// visacredit card transaction will be handled in processCardTransaction method
			if t.CardTransaction != nil {
				break
			}

			m, err = processDebitPullTransaction(*usrID, *busID, accountID, contactID, t.Amount, string(n.Action), t)
			if err != nil {
				return err
			}
		case TransactionTypeOther:
			message, err := processTransactionTypeOther(accountID, *busID, *usrID, t.Amount, t.CodeType)
			if err != nil {
				return err
			}

			if message != nil {
				m = message
			}
		case TransactionTypeFee:
			m, err = processDebitFeeTransaction(*usrID, *busID, accountID, t.Amount, t)
			if err != nil {
				return err
			}
		case TransactionTypeOtherCredit:
			m, err = processOtherCreditTransaction(*usrID, *busID, accountID, contactID, t.Amount, string(n.Action), t)
			if err != nil {
				return err
			}
		}
	}

	if m != nil && m.Counterparty != "" {
		postedTXNCreate.Counterparty = m.Counterparty
	}

	var postedTXN *transaction.BusinessPostedTransaction
	var pendingTXN *transaction.BusinessPendingTransaction

	switch t.CodeType {
	case transaction.TransactionCodeTypeDebitPosted:
		fallthrough
	case transaction.TransactionCodeTypeCreditPosted:
		if m != nil {
			postedTXNCreate.TransactionDesc = m.TransactionDescription
			postedTXNCreate.TransactionTitle = m.TransactionTitle
			postedTXNCreate.TransactionSubtype = transaction.ActivityToTransactionSubtype[m.ActivityType]
		}

		updatePostedData(*busID, postedTXNCreate, transfer, account)

		if isBankTransferMoneyRequest(transfer, t.CodeType) {
			handleBankTransferMoneyRequest(transfer, t.Amount, *busID, contactID)
		}

		foundPending := false

		// Don't send push notification if card was already authorized
		if postedTXNCreate.BankTransactionID != nil && *postedTXNCreate.BankTransactionID != "" {
			pt, err := transaction.NewPendingTransactionService().GetTransactionByBankTransactionID(*postedTXNCreate.BankTransactionID, postedTXNCreate.BusinessID)
			if err == nil && pt != nil && isPendingTransaction(pt) {
				foundPending = true

				if postedTXNCreate.CardID != nil && pt.CodeType == transaction.TransactionCodeTypeAuthApproved {
					sendPushNotification = false
				}

				log.Println("Pending w/ Partner txn ID:", pt.ID)
				postedTXNCreate.ID = shared.PostedTransactionID(pt.ID)
			}
		}

		if postedTXNCreate.MoneyTransferID != nil && *postedTXNCreate.MoneyTransferID != "" && foundPending == false {
			aID, err := id.ParseBankAccountID(shared.StringValue(postedTXNCreate.AccountID))
			if err != nil {
				log.Println(err)
			}

			if err == nil {
				pt, err := transaction.NewPendingTransactionService().GetTransactionByMoneyTransferID(
					shared.StringValue(postedTXNCreate.MoneyTransferID),
					aID,
					postedTXNCreate.BusinessID)
				if err == nil && pt != nil && isPendingTransaction(pt) {
					foundPending = true

					if postedTXNCreate.CardID != nil && pt.CodeType == transaction.TransactionCodeTypeAuthApproved {
						sendPushNotification = false
					}

					if isSnapcheck {
						log.Println("Pending w/ Snapcheck:", pt.ID)
					} else {
						log.Println("Pending w/ Transfer:", pt.ID)
					}

					postedTXNCreate.ID = shared.PostedTransactionID(pt.ID)
				} else {
					if err != nil {
						log.Println(err)
					}
				}
			}
		}

		// Delete any existing pending transaction before creating new posted transaction
		if foundPending {
			err := transaction.NewPendingTransactionService().DeleteTransaction(postedTXNCreate.BankTransactionID, postedTXNCreate.MoneyTransferID, postedTXNCreate.BusinessID)
			if err != nil {
				// Log error and fall through to continue processing notification
				// We don't delete on new txn service - just update
				log.Printf("error deleting pending transaction: %v", err)
			}
		}

		if cardTXN != nil {
			postedTXNCreate.Status = transaction.TransactionStatusCardPosted
		} else {
			postedTXNCreate.Status = transaction.TransactionStatusNonCardPosted
		}

		postedTXN, err = createPostedTransaction(&n, postedTXNCreate, cardTXN, holdTXN)
		if err != nil {
			return err
		}

		if postedTXN.TransactionSubtype != nil && *postedTXN.TransactionSubtype == transaction.TransactionSubtypeACHTransferShopifyCredit {
			err = fetchShopifyPayout(postedTXN.ID, postedTXN.BusinessID, t.Amount, postedTXN.TransactionDate)
			if err != nil {
				log.Println(err)
			}
		}

		updateActivityResourceID(string(postedTXN.ID), activityID)
	case transaction.TransactionCodeTypeAuthApproved:
		if m != nil {
			pendingTXNCreate.TransactionDesc = m.TransactionDescription
			pendingTXNCreate.TransactionTitle = m.TransactionTitle
			pendingTXNCreate.TransactionSubtype = transaction.ActivityToTransactionSubtype[m.ActivityType]
		}

		pendingTXNCreate.Status = transaction.TransactionStatusCardAuthorized
		pendingTXN, err = createPendingTransaction(&n, pendingTXNCreate, cardTXN, holdTXN)
		if err != nil {
			return err
		}

		updateActivityResourceID(string(pendingTXN.ID), activityID)
	case transaction.TransactionCodeTypeHoldApproved:
		if m != nil {
			pendingTXNCreate.TransactionDesc = m.TransactionDescription
			pendingTXNCreate.TransactionTitle = m.TransactionTitle
			pendingTXNCreate.TransactionSubtype = transaction.ActivityToTransactionSubtype[m.ActivityType]
		}

		pendingTXNCreate.Status = transaction.TransactionStatusHoldSet
		pendingTXN, err = createPendingTransaction(&n, pendingTXNCreate, cardTXN, holdTXN)
		if err != nil {
			return err
		}

		updateActivityResourceID(string(pendingTXN.ID), activityID)
	case transaction.TransactionCodeTypeAuthReversed, transaction.TransactionCodeTypeHoldReleased:
		if t.BankTransactionID != "" && busID != nil {
			pt, err := transaction.NewPendingTransactionService().GetTransactionByBankTransactionID(*pendingTXNCreate.BankTransactionID, *busID)
			if err == nil && pt != nil {
				log.Println("Pending w/ Partner txn ID:", pt.ID)
				pendingTXNCreate.ID = shared.PendingTransactionID(pt.ID)
			}

			if t.CodeType == transaction.TransactionCodeTypeAuthReversed {
				pendingTXNCreate.Status = transaction.TransactionStatusCardAuthReversed
			} else if t.CodeType == transaction.TransactionCodeTypeHoldReleased {
				pendingTXNCreate.Status = transaction.TransactionStatusHoldReleased
			}

			if pendingTXNCreate.ID != "" {
				if pt.TransactionSubtype != nil {
					pendingTXNCreate.TransactionSubtype = *pt.TransactionSubtype
				}

				pendingTXN, err = createPendingTransaction(&n, pendingTXNCreate, cardTXN, holdTXN)
				if err != nil {
					return err
				}

				err = transaction.NewPendingTransactionService().DeleteTransaction(&t.BankTransactionID, nil, *busID)
				if err != nil {
					return fmt.Errorf("error deleting pending transaction: %v", err)
				}
			}
		}

		sendPushNotification = false
	case transaction.TransactionCodeTypeAuthDeclined:
		if m != nil {
			pendingTXNCreate.TransactionDesc = m.TransactionDescription
			pendingTXNCreate.TransactionTitle = m.TransactionTitle
			pendingTXNCreate.TransactionSubtype = transaction.ActivityToTransactionSubtype[m.ActivityType]
		}

		pendingTXNCreate.Status = transaction.TransactionStatusCardAuthDeclined
		pendingTXN, err = createPendingTransaction(&n, pendingTXNCreate, cardTXN, holdTXN)
		if err != nil {
			return err
		}

		updateActivityResourceID(string(pendingTXN.ID), activityID)
	default:
		sendPushNotification = false
	}

	// Send push notification
	if m != nil && len(m.NotificationBody) > 0 && sendPushNotification {
		var ID string
		var isPosted bool
		if postedTXN != nil {
			ID = postedTXN.ID.ToPrefixString()
			isPosted = true
		} else if pendingTXN != nil {
			ID = pendingTXN.ID.ToPrefixString()
			isPosted = false
		}

		n := push.Notification{
			UserID:     usrID,
			Provider:   push.TempTextProvider{PushTitle: m.NotificationHeader, PushBody: m.NotificationBody},
			BusinessID: busID,
		}

		// Do not send transaction ID in case of declined transaction
		if t.CodeType != transaction.TransactionCodeTypeAuthDeclined {
			n.TransactionID = &ID
		}

		push.Notify(n, isPosted)
	}

	return nil
}

func isPendingTransaction(pt *transaction.BusinessPendingTransaction) bool {
	if pt == nil {
		return false
	}

	log.Println("Pending txn status is ", pt.Status)

	switch pt.Status {
	case transaction.TransactionStatusHoldSet:
		fallthrough
	case transaction.TransactionStatusCardAuthorized:
		fallthrough
	case transaction.TransactionStatusValidation:
		fallthrough
	case transaction.TransactionStatusReview:
		fallthrough
	case transaction.TransactionStatusBankProcessing:
		return true
	default:
		return false
	}
}

func updateActivityResourceID(txnID string, activityID *string) {
	if activityID != nil {
		log.Println("Activity ID - Update Resource ID ", *activityID)

		sql := `UPDATE business_activity SET resource_id = $1 WHERE id = $2`
		_, err := core.DBWrite.Exec(sql, txnID, *activityID)
		if err != nil {
			log.Println(err)
		}
	}
}

func updatePostedData(busID shared.BusinessID, txn transaction.BusinessPostedTransactionCreate, transfer *TransferDetails, account *business.BankAccount) {
	// Update posted transaction id
	if transfer != nil {
		switch txn.CodeType {
		case transaction.TransactionCodeTypeDebitPosted:
			err := business.NewMoneyTransferService(services.NewSourceRequest()).UpdateDebitPostedTransaction(busID, *transfer.MoneyTransferID, string(txn.ID))
			if err != nil {
				log.Println("Failed to update money transfer debit posted transaction: ", *transfer.MoneyTransferID, err)
			}

		case transaction.TransactionCodeTypeCreditPosted:
			err := business.NewMoneyTransferService(services.NewSourceRequest()).UpdateCreditPostedTransaction(*transfer.MoneyTransferID, string(txn.ID))
			if err != nil {
				log.Println("Failed to update money transfer credit posted transaction: ", *transfer.MoneyTransferID, err)
			}
		}

	}

}

func updateAccountBalance(account *business.BankAccount) {
	// Update account balance from partner bank data
	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bas, err := business.NewBankingAccountService()
		if err != nil {
			log.Println("updateAccountBalance Error", err)
			return
		}

		_, err = bas.GetBalanceByID(account.Id, false)

		return
	}

	providerName, ok := banking.ToPartnerBankName[account.BankName]
	if !ok {
		log.Printf("Partner bank %s does not exist", account.BankName)
		return
	}

	bank, err := partnerbank.GetBusinessBank(providerName)
	if err != nil {
		log.Println(err)
		return
	}

	srv, err := bank.BankAccountService(partnerbank.NewAPIRequest(), partnerbank.BusinessID(account.BusinessID))
	if err != nil {
		log.Println(err)
		return
	}

	resp, err := srv.Get(partnerbank.AccountBankID(account.BankAccountId))
	if err != nil {
		log.Println(err)
		return
	}

	sql := `
			UPDATE business_bank_account
			SET available_balance = $1, posted_balance = $2
			WHERE id = $3`
	_, err = core.DBWrite.Exec(sql, resp.AvailableBalance, resp.PostedBalance, account.Id)
	if err != nil {
		log.Println(err)
	}
}

// Adds card transaction to activity stream and construct body for push notification
func processCardTransaction(usrID shared.UserID, busID shared.BusinessID, accountID *string, cID *string, amt num.Decimal, notification TransactionNotification,
	txn transaction.BusinessCardTransactionCreate) (*TransactionMessage, error) {
	log.Println("Processing card transaction ", notification.CodeType)
	switch notification.CodeType {
	case transaction.TransactionCodeTypeAuthDeclined:
		fallthrough
	case transaction.TransactionCodeTypeAuthApproved, transaction.TransactionCodeTypeAuthReversed:
		return processNonPostedCardTransaction(usrID, busID, amt, notification, txn)
	case transaction.TransactionCodeTypeDebitPosted, transaction.TransactionCodeTypeCreditPosted:
		return processPostedCardTransaction(usrID, busID, accountID, cID, amt, notification, txn)
	default:
		return nil, fmt.Errorf("unhandled card transaction code type: %s", notification.CodeType)
	}
}

// TODO - better logic here
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

func processAccountOrigination(usrID shared.UserID, busID shared.BusinessID, notification TransactionNotification) (*TransactionMessage, error) {
	account, err := business.NewAccountService().GetByBankAccountId(*notification.BankAccountID, busID)
	if err != nil {
		return nil, fmt.Errorf("failed to process account origination: %v", err)
	}

	amt := num.NewZero()

	fAmount, _ := amt.Abs().Float64()

	t := activity.AccountTransaction{
		EntityID:        string(busID),
		UserID:          usrID,
		Amount:          activity.AccountTransactionAmount(fAmount),
		TransactionDate: notification.TransactionDate,
		Origin:          services.MaskLeft(account.AccountNumber, 4),
	}

	return onAccountOriginated(t)

}

func onCardTransaction(t activity.CardTransaction, activityType activity.Type, codeType transaction.TransactionCodeType) *string {

	switch codeType {
	case transaction.TransactionCodeTypeDebitPosted:
		ID, err := activity.NewCardTransationCreator().PostedDebit(t, activityType)
		if err != nil {
			log.Println(err)
		}

		return ID
	case transaction.TransactionCodeTypeAuthApproved:
		ID, err := activity.NewCardTransationCreator().Authorized(t)
		if err != nil {
			log.Println(err)
		}

		return ID
	case transaction.TransactionCodeTypeAuthReversed:
		ID, err := activity.NewCardTransationCreator().AuthReversed(t)
		if err != nil {
			log.Println(err)
		}

		return ID
	case transaction.TransactionCodeTypeHoldReleased:
		ID, err := activity.NewCardTransationCreator().HoldExpired(t)
		if err != nil {
			log.Println(err)
		}

		return ID
	case transaction.TransactionCodeTypeAuthDeclined:
		ID, err := activity.NewCardTransationCreator().Declined(t)
		if err != nil {
			log.Println(err)
		}

		return ID
	case transaction.TransactionCodeTypeCreditPosted:
		ID, err := activity.NewCardTransationCreator().PostedCredit(t, activityType)
		if err != nil {
			log.Println(err)
		}

		return ID
	}

	return nil
}

func onMoneyTransfer(t activity.AccountTransaction, activityType activity.Type, codeType transaction.TransactionCodeType) *string {
	var ID *string

	switch codeType {
	case transaction.TransactionCodeTypeCreditPosted:
		ID, _ = activity.NewTransferCreator().CreditPosted(t, activityType)
	case transaction.TransactionCodeTypeDebitPosted:
		ID, _ = activity.NewTransferCreator().DebitPosted(t, activityType)
	}

	return ID
}

func onAccountOriginated(t activity.AccountTransaction) (*TransactionMessage, error) {

	ID, err := activity.NewTransferCreator().AccountOriginated(t)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	accountNumber := string(t.Origin[len(t.Origin)-4:])

	m := TransactionMessage{
		TransactionDescription: fmt.Sprintf(AccountOriginationTransactionDescription, accountNumber),
		TransactionTitle:       AccountOriginationTransactionTitle,
		ActivtiyID:             ID,
		ActivityType:           activity.TypeAccountOrigination,
	}

	return &m, nil

}

func handleBankTransferMoneyRequest(transfer *TransferDetails, amount num.Decimal, businessID shared.BusinessID, contactID *string) error {
	// Generate receipt
	rNumber := shared.GetRandomString("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 8) + "-" +
		shared.GetRandomString("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 5)

	amt, ok := amount.Float64()
	if !ok {
		log.Println("Invalid transaction amount ", amount.FormatCurrency())
		return errors.New("Invalid transaction amount")
	}

	businessName := shared.GetBusinessName(transfer.BusinessLegalName, transfer.BusinessDBA)

	accountNumber := shared.StringValue(transfer.AccountNumber)
	accountNumber = string(accountNumber[len(accountNumber)-4:])

	r := payment.ReceiptGenerate{
		ContactFirstName:    transfer.ContactFirstName,
		ContactLastName:     transfer.ContactLastName,
		ContactBusinessName: transfer.ContactBusinessName,
		ContactEmail:        transfer.ContactEmail,
		Amount:              amt,
		Notes:               transfer.Notes,
		BusinessID:          businessID,
		UserID:              transfer.UserID,
		BusinessName:        businessName,
		ContactID:           contactID,
		BusinessPhone:       transfer.BusinessPhone,
		RequestID:           transfer.RequestID,
		ReceiptNumber:       rNumber,
		InvoiceID:           transfer.InvoiceID,
		InvoiceNumber:       transfer.InvoiceNumber,
		PaymentDate:         transfer.PaymentDate,
		PaymentBrand:        transfer.BankName,
		PaymentNumber:       &accountNumber,
	}
	content, _, err := payment.NewPaymentService(services.NewSourceRequest()).GenerateReceipt(r)
	if err != nil {
		log.Println("error generating receipt")
		//return err
	}

	paidDate := (*transfer.PaymentDate).Format("Jan _2, 2006")

	// Send receipt to consumer
	receiptRequest := payment.ReceiptRequest{
		RequestSource:       transfer.RequestSource,
		ContactFirstName:    transfer.ContactFirstName,
		ContactLastName:     transfer.ContactLastName,
		ContactBusinessName: transfer.ContactBusinessName,
		ContactEmail:        transfer.ContactEmail,
		ContactPhone:        transfer.ContactPhone,
		BusinessName:        businessName,
		BusinessEmail:       transfer.BusinessEmail,
		Notes:               transfer.Notes,
		Amount:              amt,
		ReceiptNumber:       rNumber,
		PaymentDate:         paidDate,
		Content:             content,
	}

	err = payment.NewPaymentService(services.NewSourceRequest()).SendReceiptToCustomer(receiptRequest, "", "")
	if err != nil {
		log.Println("error sending receipt to customer", err)
		return err
	}

	// Send receipt to business
	err = payment.NewPaymentService(services.NewSourceRequest()).SendReceiptToBusiness(receiptRequest, "")
	if err != nil {
		log.Println("error sending receipt to business", err)
		return err
	}

	// Update request status
	reqUpdate := payment.RequestUpdate{
		ID:     *transfer.RequestID,
		Status: payment.PaymentRequestStatusComplete,
	}
	err = payment.NewPaymentService(services.NewSourceRequest()).UpdateRequestStatus(&reqUpdate)
	if err != nil {
		log.Println("Error updating request status", err)
		return err
	}

	// Update payment status
	paymentUpdate := payment.Payment{
		ID:     *transfer.PaymentID,
		Status: payment.PaymentStatusSucceeded,
	}
	err = payment.NewPaymentService(services.NewSourceRequest()).UpdatePaymentStatus(&paymentUpdate)
	if err != nil {
		log.Println("Error updating payment status", err)
		return err
	}

	return nil
}

func isBankTransferMoneyRequest(transfer *TransferDetails, codeType transaction.TransactionCodeType) bool {
	if transfer == nil {
		return false
	}

	if transfer.RequestType == nil {
		return false
	}

	if *transfer.RequestType != string(payment.PaymentRequestTypeInvoiceBank) {
		return false
	}

	if codeType != transaction.TransactionCodeTypeCreditPosted {
		return false
	}

	return true
}

func fetchShopifyPayout(txnID shared.PostedTransactionID, businessID shared.BusinessID, amount num.Decimal, transactionDate time.Time) error {
	bus, err := bus.NewPartnerService(services.NewSourceRequest()).GetShopifyBusinessByID(businessID)
	if err != nil {
		log.Println(err)
		return err
	}

	if bus.InstallStatus == grpcShopify.ShopifyInstallStatus_SIS_INACTIVE {
		log.Println("Shopify app has been uninstalled")
		return nil
	}

	amt := amount.Abs().FormatCurrency()

	sn, err := grpc.GetConnectionStringForService(grpc.ServiceNameShopping)
	if err != nil {
		return err
	}

	client, err := grpc.NewInsecureClient(sn)
	if err != nil {
		return err
	}

	defer client.CloseAndCancel()
	shopifyServiceClient := grpcShopify.NewShopifyBusinessServiceClient(client.GetConn())

	layout := "2006-01-02"
	maxDate := transactionDate.Format(layout)
	minDate := transactionDate.AddDate(0, -1, 0).Format(layout)

	req := grpcShopify.ShopifyPayoutRequest{
		BusinessId:    businessID.ToPrefixString(),
		Amount:        amt,
		PayoutStatus:  grpcShopify.ShopifyPayoutStatus_SPS_PAID,
		MinDate:       minDate,
		MaxDate:       maxDate,
		TransactionId: txnID.ToPrefixString(),
	}

	payout, err := shopifyServiceClient.AddShopifyPayout(client.GetContext(), &req)
	if err != nil {
		return err
	}

	log.Println("Added shopify payout: ", txnID, payout.Id, payout.BusinessId)

	return nil
}
func fillTransferDetailsFromInvoiceSvc(transferDetails *TransferDetails) error {
	invoiceDetail, err := getInvoiceFromInvoiceId(transferDetails.RequestID)
	if err != nil {
		return err
	}
	td := TransferDetails{}
	query := `
				SELECT
					business_money_request_payment.id "business_money_request_payment.id",
					business_money_request_payment.linked_bank_account_id "business_money_request_payment.linked_bank_account_id",
					business_money_request_payment.payment_date "business_money_request_payment.payment_date"
				FROM business_money_request_payment 
				WHERE business_money_request_payment.invoice_id = $1`

	err = core.DBRead.Get(&td, query, invoiceDetail.InvoiceID.UUIDString())
	if err != nil {
		return err
	}
	// fill the retried values into the transfer details
	transferDetails.PaymentID = td.PaymentID
	transferDetails.LinkedBankAccountID = td.LinkedBankAccountID
	transferDetails.PaymentDate = td.PaymentDate

	contactIdUUIDStr := invoiceDetail.ContactID.UUIDString()
	transferDetails.RequestContactID = &contactIdUUIDStr
	reqestTypeCard := string(payment.PaymentRequestTypeInvoiceCard)
	if invoiceDetail.AllowCard {
		transferDetails.RequestType = &reqestTypeCard
	}
	invNumber := fmt.Sprintf("%05d", invoiceDetail.Number)
	transferDetails.InvoiceNumber = &invNumber
	invIdStr := invoiceDetail.InvoiceID.UUIDString()
	transferDetails.InvoiceID = &invIdStr
	return nil
}

func getInvoiceFromInvoiceId(paymentRequestId *shared.PaymentRequestID) (*invoice.Invoice, error) {
	invSvc, err := invoice.NewInvoiceService()
	if err != nil {
		return nil, err
	}
	resp, err := invSvc.GetInvoiceIDFromPaymentRequestID(paymentRequestId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
