package notification

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/wiseco/core-platform/notification/activity"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking/business"
	core "github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/core-platform/transaction"
	"github.com/wiseco/go-lib/num"
)

func processNonPostedCardTransaction(uID shared.UserID, bID shared.BusinessID, amt num.Decimal, notification TransactionNotification,
	txn transaction.BusinessCardTransactionCreate) (*TransactionMessage, error) {
	var m *TransactionMessage
	var err error

	amount := amt.Abs().FormatCurrency()

	businessName, err := getBusinessName(bID)
	if err != nil {
		return nil, err
	}

	// Fetch card number
	card, err := business.NewCardService(services.NewSourceRequest()).GetByBankCardId(*notification.BankCardID, uID)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to retrieve card %s for user %s", *notification.BankCardID, uID))
	}

	cardNumber, err := card.GetCardNumberLastFour()

	if err != nil {
		return nil, err
	}

	var location *string
	if len(txn.MerchantCity) > 0 && isLetter(txn.MerchantCity) {
		l := strings.TrimSpace(txn.MerchantCity)
		location = &l
	}

	merchantName := notification.GetMerchantName()

	switch notification.CodeType {
	case transaction.TransactionCodeTypeAuthApproved:
		m, err = processCardAuthorizationTransaction(cardNumber, amount, location, merchantName, notification.CardTransaction.POSEntryMode, *businessName)
		if err != nil {
			return nil, err
		}
	case transaction.TransactionCodeTypeAuthDeclined:
		m, err = processCardDeclinedTransaction(cardNumber, amount, location, merchantName, notification.CardTransaction.POSEntryMode,
			notification.CardTransaction.AuthResponseCode, *businessName)
		if err != nil {
			return nil, err
		}
	case transaction.TransactionCodeTypeAuthReversed:
		m = &TransactionMessage{
			ActivityType: activity.TypeCardTransaction,
		}
	}

	// Push to activity stream
	t := activity.CardTransaction{
		EntityID:     string(bID),
		UserID:       uID,
		Amount:       amount,
		Number:       cardNumber,
		Merchant:     notification.GetMerchantName(),
		BusinessName: *businessName,
	}

	ID := onCardTransaction(t, m.ActivityType, notification.CodeType)
	m.ActivtiyID = ID

	return m, nil
}

func processCardAuthorizationTransaction(cardNumber string, amount string, purchaseLocation *string, merchantName string,
	posEntryMode POSEntryMode, businessName string) (*TransactionMessage, error) {
	m := TransactionMessage{}

	m.NotificationHeader = fmt.Sprintf(DebitCardAuthorizationNotificationHeader, businessName, cardNumber)

	if len(merchantName) > 0 {
		m.TransactionTitle = fmt.Sprintf(DebitCardAuthorizationTransactionTitle, merchantName)

		if posEntryMode.isOnlinePayment() {
			m.NotificationBody = fmt.Sprintf(DebitCardAuthorizationNotificationWithoutLocationBody, businessName, cardNumber, merchantName, amount)
			m.TransactionDescription = fmt.Sprintf(DebitCardAuthorizationTransactionWithoutLocationDescription, businessName, cardNumber, merchantName, amount)

			m.ActivityType = activity.TypeCardReaderPurchaseDebitOnline
			m.BusinessName = businessName

			return &m, nil
		}

		if purchaseLocation != nil {
			m.NotificationBody = fmt.Sprintf(DebitCardAuthorizationNotificationBody, businessName, cardNumber, merchantName, *purchaseLocation, amount)
			m.TransactionDescription = fmt.Sprintf(DebitCardAuthorizationTransactionDescription, businessName, cardNumber, merchantName, *purchaseLocation, amount)
		} else {
			m.NotificationBody = fmt.Sprintf(DebitCardAuthorizationNotificationWithoutLocationBody, businessName, cardNumber, merchantName, amount)
			m.TransactionDescription = fmt.Sprintf(DebitCardAuthorizationTransactionWithoutLocationDescription, businessName, cardNumber, merchantName, amount)
		}
	} else {
		m.TransactionTitle = fmt.Sprintf(DebitCardAuthorizationTransactionWithoutMerchantTitle, amount)

		if posEntryMode.isOnlinePayment() {
			m.NotificationBody = fmt.Sprintf(DebitCardAuthorizationNotificationGenericBody, businessName, cardNumber, amount)
			m.TransactionDescription = fmt.Sprintf(DebitCardAuthorizationTransactionGenericDescription, businessName, cardNumber, amount)

			m.ActivityType = activity.TypeCardReaderPurchaseDebitOnline
			m.BusinessName = businessName

			return &m, nil
		}

		if purchaseLocation != nil {
			m.NotificationBody = fmt.Sprintf(DebitCardAuthorizationNotificationWithoutMerchantBody, businessName, cardNumber, *purchaseLocation, amount)
			m.TransactionDescription = fmt.Sprintf(DebitCardAuthorizationTransactionWithoutMerchantDescription, businessName, cardNumber, *purchaseLocation, amount)
		} else {
			m.NotificationBody = fmt.Sprintf(DebitCardAuthorizationNotificationGenericBody, businessName, cardNumber, amount)
			m.TransactionDescription = fmt.Sprintf(DebitCardAuthorizationTransactionGenericDescription, businessName, cardNumber, amount)
		}
	}

	m.ActivityType = activity.TypeCardReaderPurchaseDebit
	m.BusinessName = businessName

	return &m, nil
}

func processCardDeclinedTransaction(cardNumber string, amount string, purchaseLocation *string, merchantName string,
	posEntryMode POSEntryMode, authResponseCode AuthResponseCode, businessName string) (*TransactionMessage, error) {
	m := TransactionMessage{}

	m.NotificationHeader = fmt.Sprintf(DebitCardDeclinedNotificationHeader, businessName, cardNumber)

	if len(merchantName) > 0 {
		m.TransactionTitle = fmt.Sprintf(DebitCardDeclinedTransactionTitle, merchantName)

		if posEntryMode.isOnlinePayment() {
			m.NotificationBody = fmt.Sprintf(DebitCardDeclinedNotificationWithoutLocationBody, businessName, cardNumber, merchantName, amount)
			m.TransactionDescription = fmt.Sprintf(DebitCardDeclinedTransactionWithoutLocationDescription, businessName, cardNumber, merchantName, amount)

			m.ActivityType = activity.TypeCardReaderPurchaseDebitOnline
			m.BusinessName = businessName

			return &m, nil
		}

		if purchaseLocation != nil {
			m.NotificationBody = fmt.Sprintf(DebitCardDeclinedNotificationBody, businessName, cardNumber, merchantName, *purchaseLocation, amount)
			m.TransactionDescription = fmt.Sprintf(DebitCardDeclinedTransactionDescription, businessName, cardNumber, merchantName, *purchaseLocation, amount)
		} else {
			m.NotificationBody = fmt.Sprintf(DebitCardDeclinedNotificationWithoutLocationBody, businessName, cardNumber, merchantName, amount)
			m.TransactionDescription = fmt.Sprintf(DebitCardDeclinedTransactionWithoutLocationDescription, businessName, cardNumber, merchantName, amount)
		}
	} else {
		m.TransactionTitle = fmt.Sprintf(DebitCardDeclinedTransactionWithoutMerchantTitle, amount)

		if posEntryMode.isOnlinePayment() {
			m.NotificationBody = fmt.Sprintf(DebitCardDeclinedNotificationGenericBody, businessName, cardNumber, amount)
			m.TransactionDescription = fmt.Sprintf(DebitCardDeclinedTransactionGenericDescription, businessName, cardNumber, amount)

			m.ActivityType = activity.TypeCardReaderPurchaseDebitOnline
			m.BusinessName = businessName

			return &m, nil
		}

		if purchaseLocation != nil {
			m.NotificationBody = fmt.Sprintf(DebitCardDeclinedNotificationWithoutMerchantBody, businessName, cardNumber, *purchaseLocation, amount)
			m.TransactionDescription = fmt.Sprintf(DebitCardDeclinedTransactionWithoutMerchantDescription, businessName, cardNumber, *purchaseLocation, amount)
		} else {
			m.NotificationBody = fmt.Sprintf(DebitCardDeclinedNotificationGenericBody, businessName, cardNumber, amount)
			m.TransactionDescription = fmt.Sprintf(DebitCardDeclinedTransactionGenericDescription, businessName, cardNumber, amount)
		}
	}

	declineReason, ok := AuthRespCodeToMessage[authResponseCode]
	if ok {
		m.NotificationBody = m.NotificationBody + " due to " + declineReason
	}

	m.ActivityType = activity.TypeCardReaderPurchaseDebit
	m.BusinessName = businessName

	return &m, nil
}

func processTransactionTypeOther(accountID *string, bID shared.BusinessID, uID shared.UserID, amt num.Decimal, codeType transaction.TransactionCodeType) (*TransactionMessage, error) {
	// Get business name
	type BusinessDetails struct {
		LegalName     string               `db:"business.legal_name"`
		DBA           services.StringArray `db:"business.dba"`
		AccountNumber string               `db:"business_bank_account.account_number"`
	}

	bt := BusinessDetails{}

	var err error
	err = core.DBRead.Get(
		&bt,
		`SELECT business.legal_name "business.legal_name", business.dba "business.dba"
			FROM business 
			WHERE business.id = $1`, bID,
	)
	if err != nil {
		log.Println("Error retrieving business details ", err)
		return nil, err
	}

	if accountID != nil && *accountID != "" {
		acc, err := business.NewAccountService().GetByIDInternal(*accountID)
		if err != nil {
			log.Println("Error retrieving account details ", err)
			return nil, err
		}

		bt.AccountNumber = acc.AccountNumber
	}

	businessName := shared.GetBusinessName(&bt.LegalName, bt.DBA)

	switch codeType {
	case transaction.TransactionCodeTypeHoldApproved:
		return processHoldApprovedTransaction(bt.AccountNumber, bID, uID, businessName, amt)
	case transaction.TransactionCodeTypeHoldReleased:
		return processHoldReleasedTransaction(bID, uID, businessName, amt)
	case transaction.TransactionCodeTypeAuthApproved, transaction.TransactionCodeTypeAuthReversed, transaction.TransactionCodeTypeAuthDeclined:
		return nil, nil
	default:
		msg := fmt.Sprintf("unhandled code type: %s", string(codeType))
		return nil, errors.New(msg)
	}
}

func processHoldApprovedTransaction(accountNumber string, bID shared.BusinessID, uID shared.UserID,
	businessName string, amt num.Decimal) (*TransactionMessage, error) {

	amount := amt.Abs().FormatCurrency()

	accountNumber = string(accountNumber[len(accountNumber)-4:])

	m := TransactionMessage{}

	m.ActivityType = activity.TypeHoldApproved
	m.BusinessName = businessName
	m.NotificationHeader = fmt.Sprintf(AccountHoldNotificationHeader, amount)
	m.NotificationBody = fmt.Sprintf(AccountHoldNotificationBody, amount, businessName, accountNumber)
	m.TransactionTitle = fmt.Sprintf(AccountHoldTransactionTitle, amount)
	m.TransactionDescription = fmt.Sprintf(AccountHoldTransactionDescription, amount, businessName, accountNumber)

	fAmount, _ := amt.Abs().Float64()

	a := activity.AccountTransaction{
		BusinessName:    &businessName,
		EntityID:        string(bID),
		UserID:          uID,
		Amount:          activity.AccountTransactionAmount(fAmount),
		TransactionDate: time.Now(),
	}

	_, err := activity.NewTransferCreator().HoldApproved(a)
	if err != nil {
		log.Println(err)
	}

	return &m, nil
}

func processHoldReleasedTransaction(bID shared.BusinessID, uID shared.UserID,
	businessName string, amt num.Decimal) (*TransactionMessage, error) {
	m := TransactionMessage{}

	m.ActivityType = activity.TypeHoldReleased
	m.BusinessName = businessName

	fAmount, _ := amt.Abs().Float64()

	a := activity.AccountTransaction{
		BusinessName:    &businessName,
		EntityID:        string(bID),
		UserID:          uID,
		Amount:          activity.AccountTransactionAmount(fAmount),
		TransactionDate: time.Now(),
	}

	_, err := activity.NewTransferCreator().HoldReleased(a)
	if err != nil {
		log.Println(err)
	}

	return &m, nil
}
