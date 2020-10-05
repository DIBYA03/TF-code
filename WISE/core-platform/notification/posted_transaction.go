package notification

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/wiseco/core-platform/notification/activity"
	"github.com/wiseco/core-platform/partner/service/sendgrid"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking"
	"github.com/wiseco/core-platform/services/banking/business"
	"github.com/wiseco/core-platform/services/banking/business/contact"
	bus "github.com/wiseco/core-platform/services/business"
	con "github.com/wiseco/core-platform/services/contact"
	core "github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/services/payment"
	usr "github.com/wiseco/core-platform/services/user"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/core-platform/transaction"
	"github.com/wiseco/go-lib/num"
	grpcBankTransfer "github.com/wiseco/protobuf/golang/banking/transfer"
)

func processPostedCardTransaction(uID shared.UserID, bID shared.BusinessID, accID *string, cID *string, amt num.Decimal, notification TransactionNotification,
	txn transaction.BusinessCardTransactionCreate) (*TransactionMessage, error) {
	var location *string
	if len(txn.MerchantCity) > 0 && isLetter(txn.MerchantCity) {
		l := strings.TrimSpace(txn.MerchantCity)
		location = &l
	}

	if len(txn.MerchantState) > 0 {
		var l string
		if location != nil {
			*location = *location + ", "
		} else {
			// initialize
			location = &l
		}

		*location = *location + txn.MerchantState
	}

	merchantName := notification.GetMerchantName()

	var m *TransactionMessage
	var err error
	switch notification.TransactionType {
	case TransactionTypePurchase:
		m, err = processCardDebitPurchaseTransaction(cID, bID, merchantName, notification.CardTransaction.POSEntryMode, location, amt)
		if err != nil {
			log.Println("Error processing purchase card transaction", err)
			return nil, err
		}
	case TransactionTypeATM:
		m, err = processCardDebitATMTransaction(cID, bID, location, amt)
		if err != nil {
			log.Println("Error processing atm card transaction", err)
			return nil, err
		}
	case TransactionTypeRefund:
		m, err = processRefundTransaction(accID, bID, txn.MerchantName, amt, notification)
		if err != nil {
			log.Println("Error processing posted card transaction", err)
			return nil, err
		}

		if m.SenderReceiverName != nil {
			txn.MerchantName = *m.SenderReceiverName
		}
	case TransactionTypeFee:
		return nil, nil
	case TransactionTypeVisaCredit:
		m, err = processCardVisaCreditransaction(accID, bID, amt, notification)
		if err != nil {
			log.Println("Error processing visa credit transaction", err)
			return nil, err
		}

		if m.SenderReceiverName != nil {
			txn.MerchantName = *m.SenderReceiverName
			merchantName = *m.SenderReceiverName
		}
	default:
		e := fmt.Sprintf("Invalid card transaction type %s", notification.TransactionType)
		return nil, errors.New(e)
	}

	// Push to activity stream
	t := activity.CardTransaction{
		EntityID:     string(bID),
		UserID:       uID,
		Amount:       amt.Abs().FormatCurrency(),
		Merchant:     merchantName,
		BusinessName: m.BusinessName,
	}

	ID := onCardTransaction(t, m.ActivityType, notification.CodeType)
	m.ActivtiyID = ID

	return m, nil
}

func processCardDebitPurchaseTransaction(cID *string, bID shared.BusinessID, merchantName string,
	posEntryMode POSEntryMode, location *string, amt num.Decimal) (*TransactionMessage, error) {
	m := TransactionMessage{}

	amount := amt.Abs().FormatCurrency()

	// Get business name
	type BusinessDetails struct {
		CardNumber string               `db:"business_bank_card.card_number_masked"`
		LegalName  string               `db:"business.legal_name"`
		DBA        services.StringArray `db:"business.dba"`
	}

	bt := BusinessDetails{}

	// If card id is nil
	if cID == nil {
		err := core.DBRead.Get(
			&bt,
			`
			SELECT
				business.legal_name "business.legal_name",
				business.dba "business.dba"
			FROM business 
			WHERE id = $1`,
			bID,
		)
		if err != nil {
			log.Println("Error retrieving business details:", err)
			return nil, err
		}

		businessName := shared.GetBusinessName(&bt.LegalName, bt.DBA)

		if len(merchantName) > 0 {
			m.Counterparty = merchantName
			m.NotificationHeader = fmt.Sprintf(DebitCardPurchaseNotificationHeader, businessName, merchantName, amount)
			m.TransactionTitle = fmt.Sprintf(DebitCardPurchaseTransactionTitle, merchantName)
			m.NotificationBody = fmt.Sprintf(DebitCardPurchaseNotificationWithoutCardNumberBody, businessName, merchantName, amount)
			m.TransactionDescription = fmt.Sprintf(DebitCardPurchaseTransactionWithoutCardNumberDescription, businessName, merchantName, amount)
		} else {
			m.NotificationHeader = fmt.Sprintf(DebitCardPurchaseNotificationWithoutMerchantHeader, businessName, amount)
			m.TransactionTitle = fmt.Sprintf(DebitCardPurchaseTransactionWithoutMerchantTitle, amount)
			m.NotificationBody = fmt.Sprintf(DebitCardPurchaseNotificationGenericWithoutCardNumberBody, businessName, amount)
			m.TransactionDescription = fmt.Sprintf(DebitCardPurchaseTransactionGenericWithoutCardNumberDescription, businessName, amount)
		}

		m.ActivityType = activity.TypeCardReaderPurchaseDebit
		m.BusinessName = businessName

		return &m, nil
	}

	err := core.DBRead.Get(
		&bt,
		`
		SELECT
			business.legal_name "business.legal_name",
			business.dba "business.dba", 
			business_bank_card.card_number_masked "business_bank_card.card_number_masked"
		FROM business_bank_card 
		JOIN business ON business_bank_card.business_id = business.id
		WHERE business_bank_card.business_id = $1 AND business_bank_card.id = $2`,
		bID, *cID,
	)
	if err != nil {
		log.Println("Error retrieving business details 1:", err)
		return nil, err
	}

	businessName := shared.GetBusinessName(&bt.LegalName, bt.DBA)

	lastFour := string(bt.CardNumber[len(bt.CardNumber)-4:])

	if len(merchantName) > 0 {
		m.Counterparty = merchantName
		m.NotificationHeader = fmt.Sprintf(DebitCardPurchaseNotificationHeader, businessName, merchantName, amount)
		m.TransactionTitle = fmt.Sprintf(DebitCardPurchaseTransactionTitle, merchantName)

		if posEntryMode.isOnlinePayment() {
			m.NotificationBody = fmt.Sprintf(DebitCardPurchaseNotificationWithoutLocationBody, businessName, merchantName, amount, lastFour)
			m.TransactionDescription = fmt.Sprintf(DebitCardPurchaseTransactionWithoutLocationDescription, businessName, merchantName, amount, lastFour)

			m.ActivityType = activity.TypeCardReaderPurchaseDebitOnline
			m.BusinessName = businessName

			return &m, nil
		}

		if location != nil {
			m.NotificationBody = fmt.Sprintf(DebitCardPurchaseNotificationBody, businessName, merchantName, amount, *location, lastFour)
			m.TransactionDescription = fmt.Sprintf(DebitCardPurchaseTransactionDescription, businessName, merchantName, amount, *location, lastFour)
		} else {
			m.NotificationBody = fmt.Sprintf(DebitCardPurchaseNotificationWithoutLocationBody, businessName, merchantName, amount, lastFour)
			m.TransactionDescription = fmt.Sprintf(DebitCardPurchaseTransactionWithoutLocationDescription, businessName, merchantName, amount, lastFour)
		}
	} else {
		m.NotificationHeader = fmt.Sprintf(DebitCardPurchaseNotificationWithoutMerchantHeader, businessName, amount)
		m.TransactionTitle = fmt.Sprintf(DebitCardPurchaseTransactionWithoutMerchantTitle, amount)

		if posEntryMode.isOnlinePayment() {
			m.NotificationBody = fmt.Sprintf(DebitCardPurchaseNotificationGenericBody, businessName, amount, lastFour)
			m.TransactionDescription = fmt.Sprintf(DebitCardPurchaseTransactionGenericDescription, businessName, amount, lastFour)

			m.ActivityType = activity.TypeCardReaderPurchaseDebitOnline
			m.BusinessName = businessName

			return &m, nil
		}

		if location != nil {
			m.NotificationBody = fmt.Sprintf(DebitCardPurchaseNotificationWithoutMerchantBody, businessName, amount, *location, lastFour)
			m.TransactionDescription = fmt.Sprintf(DebitCardPurchaseTransactionWithoutMerchantDescription, businessName, amount, *location, lastFour)
		} else {
			m.NotificationBody = fmt.Sprintf(DebitCardPurchaseNotificationGenericBody, businessName, amount, lastFour)
			m.TransactionDescription = fmt.Sprintf(DebitCardPurchaseTransactionGenericDescription, businessName, amount, lastFour)
		}
	}

	m.ActivityType = activity.TypeCardReaderPurchaseDebit
	m.BusinessName = businessName

	return &m, nil
}

func processCardDebitATMTransaction(cID *string, bID shared.BusinessID, location *string, amt num.Decimal) (*TransactionMessage, error) {
	m := TransactionMessage{}

	amount := amt.Abs().FormatCurrency()

	// Get business name
	type BusinessDetails struct {
		CardNumber string               `db:"business_bank_card.card_number_masked"`
		LegalName  string               `db:"business.legal_name"`
		DBA        services.StringArray `db:"business.dba"`
	}

	bt := BusinessDetails{}

	// If card id is nil
	if cID == nil {
		err := core.DBRead.Get(
			&bt,
			`
			SELECT
				business.legal_name "business.legal_name",
				business.dba "business.dba"
			FROM business 
			WHERE id = $1`,
			bID,
		)
		if err != nil {
			log.Println("Error retrieving business details:", err)
			return nil, err
		}

		businessName := shared.GetBusinessName(&bt.LegalName, bt.DBA)
		m.NotificationHeader = fmt.Sprintf(DebitCardATMNotificationHeader, businessName, amount)
		m.TransactionTitle = fmt.Sprintf(DebitCardATMTransactionTitle)

		if location != nil {
			m.NotificationBody = fmt.Sprintf(DebitCardATMNotificationWithoutCardNumberBody, businessName, amount, *location)
			m.TransactionDescription = fmt.Sprintf(DebitCardATMTransactionWithoutCardNumberDescription, businessName, amount, *location)
		} else {
			m.NotificationBody = fmt.Sprintf(DebitCardATMNotificationGenericBody, businessName, amount)
			m.TransactionDescription = fmt.Sprintf(DebitCardATMTransactionGenericDescription, businessName, amount)
		}

		m.ActivityType = activity.TypeCardATMDebit
		m.BusinessName = businessName

		return &m, nil
	}

	err := core.DBRead.Get(
		&bt,
		`
		SELECT
			business.legal_name "business.legal_name",
			business.dba "business.dba", 
			business_bank_card.card_number_masked "business_bank_card.card_number_masked"
		FROM business_bank_card 
		JOIN business ON business_bank_card.business_id = business.id
		WHERE business_bank_card.business_id = $1 AND business_bank_card.id = $2`,
		bID, *cID,
	)
	if err != nil {
		log.Println("Error retrieving business details 2:", err)
		return nil, err
	}

	businessName := shared.GetBusinessName(&bt.LegalName, bt.DBA)

	lastFour := string(bt.CardNumber[len(bt.CardNumber)-4:])

	m.NotificationHeader = fmt.Sprintf(DebitCardATMNotificationHeader, businessName, amount)
	m.TransactionTitle = fmt.Sprintf(DebitCardATMTransactionTitle)
	if location != nil {
		m.NotificationBody = fmt.Sprintf(DebitCardATMNotificationBody, businessName, amount, *location, lastFour)
		m.TransactionDescription = fmt.Sprintf(DebitCardATMTransactionDescription, businessName, amount, *location, lastFour)
	} else {
		m.NotificationBody = fmt.Sprintf(DebitCardATMNotificationWithoutLocationBody, businessName, amount, lastFour)
		m.TransactionDescription = fmt.Sprintf(DebitCardATMTransactionWithoutLocationDescription, businessName, amount, lastFour)
	}

	m.ActivityType = activity.TypeCardATMDebit
	m.BusinessName = businessName

	return &m, nil
}

func processDebitFeeTransaction(uID shared.UserID, bID shared.BusinessID, accID *string,
	amt num.Decimal, notification TransactionNotification) (*TransactionMessage, error) {

	m := TransactionMessage{}

	// Get business name
	type BusinessDetails struct {
		AccountNumber string               `db:"business_bank_account.account_number"`
		LegalName     string               `db:"business.legal_name"`
		DBA           services.StringArray `db:"business.dba"`
	}

	bt := BusinessDetails{}

	err := core.DBRead.Get(
		&bt,
		`SELECT business.legal_name "business.legal_name", business.dba "business.dba"
		FROM business
		WHERE business.id = $1`,
		bID,
	)
	if err != nil {
		log.Println("Error retrieving business details ", err)
		return nil, err
	}

	if accID != nil && *accID != "" {
		acc, err := business.NewAccountService().GetByIDInternal(*accID)
		if err != nil {
			log.Println("Error retrieving account details ", err)
			return nil, err
		}

		bt.AccountNumber = acc.AccountNumber
	}

	amount := amt.Abs().FormatCurrency()

	businessName := shared.GetBusinessName(&bt.LegalName, bt.DBA)
	feeType := GetFeeType(notification.BankTransactionDesc)

	if len(feeType) > 0 {
		m.NotificationHeader = fmt.Sprintf(DebitFeeNotificationHeader, businessName, amount)
		m.NotificationBody = fmt.Sprintf(DebitFeeNotificationBody, businessName, feeType, amount)
		m.TransactionTitle = fmt.Sprintf(DebitFeeTransactionTitle, feeType)
		m.TransactionDescription = fmt.Sprintf(DebitFeeTransactionDescription, businessName, amount, feeType)
	} else {
		m.NotificationHeader = fmt.Sprintf(DebitFeeNotificationHeader, businessName, amount)
		m.NotificationBody = fmt.Sprintf(DebitFeeNotificationGenericBody, businessName, amount)
		m.TransactionTitle = DebitFeeTransactionWithoutTypeTitle
		m.TransactionDescription = fmt.Sprintf(DebitFeeTransactionWithoutTypeDescription, businessName, amount)
	}

	m.ActivityType = activity.TypeFeeDebit
	m.BusinessName = businessName

	fAmount, _ := amt.Abs().Float64()

	// Add to activity stream
	t := activity.AccountTransaction{
		EntityID:            string(bID),
		UserID:              uID,
		Amount:              activity.AccountTransactionAmount(fAmount),
		TransactionDate:     notification.TransactionDate,
		BusinessName:        &m.BusinessName,
		ContactName:         &feeType,
		InterestEarnedMonth: nil,
	}
	ID := onMoneyTransfer(t, m.ActivityType, notification.CodeType)
	m.ActivtiyID = ID

	return &m, nil
}

func processCardVisaCreditransaction(accID *string, bID shared.BusinessID, amt num.Decimal, notification TransactionNotification) (*TransactionMessage, error) {
	m := TransactionMessage{}

	amount := amt.Abs().FormatCurrency()

	// Get business name
	type BusinessDetails struct {
		AccountNumber string               `db:"business_bank_account.account_number"`
		LegalName     string               `db:"business.legal_name"`
		DBA           services.StringArray `db:"business.dba"`
	}

	bt := BusinessDetails{}

	err := core.DBRead.Get(
		&bt,
		`
		SELECT
			business.legal_name "business.legal_name",
			business.dba "business.dba"
		FROM business 
		WHERE id = $1`,
		bID,
	)
	if err != nil {
		log.Println("Error retrieving business details 2:", err)
		return nil, err
	}

	if accID != nil && *accID != "" {
		acc, err := business.NewAccountService().GetByIDInternal(*accID)
		if err != nil {
			log.Println("Error retrieving account details ", err)
			return nil, err
		}

		bt.AccountNumber = acc.AccountNumber
	}

	senderName := transaction.GetVisaCreditSenderName(notification.BankTransactionDesc)
	m.Counterparty = senderName

	businessName := shared.GetBusinessName(&bt.LegalName, bt.DBA)
	accountNumber := string(bt.AccountNumber[len(bt.AccountNumber)-4:])

	m.NotificationHeader = fmt.Sprintf(CreditCardVisaCreditNotificationHeader, businessName, amount)
	m.TransactionDescription = fmt.Sprintf(CreditCardVisaCreditTransactionDescription, amount)
	if len(senderName) > 0 {
		m.NotificationBody = fmt.Sprintf(CreditCardVisaCreditNotificationBody, businessName, amount, senderName, accountNumber)
		m.TransactionTitle = fmt.Sprintf(CreditCardVisaCreditTransactionTitle, senderName)
	} else {
		m.NotificationBody = fmt.Sprintf(CreditCardVisaCreditNotificationWithoutSenderBody, businessName, amount, accountNumber)
		m.TransactionTitle = CreditCardVisaCreditTransactionWithoutSenderTitle
	}

	m.ActivityType = activity.TypeCardVisaCredit
	m.BusinessName = businessName
	m.SenderReceiverName = &senderName

	return &m, nil
}

func processRefundTransaction(accID *string, bID shared.BusinessID, merchantName string, amt num.Decimal, notification TransactionNotification) (*TransactionMessage, error) {

	if notification.CardTransaction != nil && len(notification.CardTransaction.TransactionType) > 0 {
		transactionType := notification.CardTransaction.TransactionType

		if transactionType.IsRefundTypeInstantPay() {
			return processCardCreditInstantPayTransaction(accID, bID, notification.BankTransactionDesc, amt)
		} else {
			return processCardMerchantRefundTransaction(accID, bID, merchantName, amt)
		}
	}

	return processCardMerchantRefundTransaction(accID, bID, merchantName, amt)
}

func processCardCreditInstantPayTransaction(accID *string, bID shared.BusinessID, transactionDescription *string, amt num.Decimal) (*TransactionMessage, error) {
	m := TransactionMessage{}

	amount := amt.Abs().FormatCurrency()

	// Get business name
	type BusinessDetails struct {
		AccountNumber string               `db:"business_bank_account.account_number"`
		LegalName     string               `db:"business.legal_name"`
		DBA           services.StringArray `db:"business.dba"`
	}

	bt := BusinessDetails{}

	err := core.DBRead.Get(
		&bt, `
		SELECT
			legal_name "business.legal_name",
			dba "business.dba"
		FROM business
		WHERE id = $1`,
		bID,
	)
	if err != nil {
		log.Println("Error retrieving business details 3:", err)
		return nil, err
	}

	if accID != nil && *accID != "" {
		acc, err := business.NewAccountService().GetByIDInternal(*accID)
		if err != nil {
			log.Println("Error retrieving account details ", err)
			return nil, err
		}

		bt.AccountNumber = acc.AccountNumber
	}

	businessName := shared.GetBusinessName(&bt.LegalName, bt.DBA)
	accountNumber := string(bt.AccountNumber[len(bt.AccountNumber)-4:])

	senderName := transaction.GetInstantPaySenderName(transactionDescription)
	m.Counterparty = senderName

	m.NotificationHeader = fmt.Sprintf(CreditCardInstantPayNotificationHeader, businessName, amount)
	m.TransactionDescription = fmt.Sprintf(CreditCardInstantPayTransactionDescription, amount)
	if len(senderName) > 0 {
		m.NotificationBody = fmt.Sprintf(CreditCardInstantPayNotificationBody, businessName, senderName, amount, businessName, accountNumber)
		m.TransactionTitle = fmt.Sprintf(CreditCardInstantPayTransactionTitle, senderName)
	} else {
		m.NotificationBody = fmt.Sprintf(CreditCardInstantPayNotificationWithoutSenderBody, businessName, amount, businessName, accountNumber)
		m.TransactionTitle = fmt.Sprintf(CreditCardInstantPayTransactionWithoutSenderTitle)
	}

	m.ActivityType = activity.TypeCardPushCredit
	m.BusinessName = businessName
	m.SenderReceiverName = &senderName

	return &m, nil
}

func processCardMerchantRefundTransaction(accID *string, bID shared.BusinessID, merchantName string, amt num.Decimal) (*TransactionMessage, error) {
	m := TransactionMessage{}

	amount := amt.Abs().FormatCurrency()

	// Get business name
	type BusinessDetails struct {
		AccountNumber string               `db:"business_bank_account.account_number"`
		LegalName     string               `db:"business.legal_name"`
		DBA           services.StringArray `db:"business.dba"`
	}

	bt := BusinessDetails{}

	err := core.DBRead.Get(
		&bt, `
	    SELECT
    	    legal_name "business.legal_name",
        	dba "business.dba"
	    FROM business
    	WHERE id = $1`,
		bID,
	)
	if err != nil {
		log.Println("Error retrieving business details 4:", err)
		return nil, err
	}

	if accID != nil && *accID != "" {
		acc, err := business.NewAccountService().GetByIDInternal(*accID)
		if err != nil {
			log.Println("Error retrieving account details ", err)
			return nil, err
		}

		bt.AccountNumber = acc.AccountNumber
	}

	businessName := shared.GetBusinessName(&bt.LegalName, bt.DBA)

	accountNumber := string(bt.AccountNumber[len(bt.AccountNumber)-4:])

	m.Counterparty = merchantName

	m.NotificationHeader = fmt.Sprintf(CreditMerchantRefundNotificationHeader, businessName, amount)
	m.TransactionDescription = fmt.Sprintf(CreditMerchantRefundTransactionDescription, amount)
	if len(merchantName) > 0 {
		m.NotificationBody = fmt.Sprintf(CreditMerchantRefundNotificationBody, businessName, merchantName, amount, businessName, accountNumber)
		m.TransactionTitle = fmt.Sprintf(CreditMerchantRefundTransactionTitle, merchantName)
	} else {
		m.NotificationBody = fmt.Sprintf(CreditMerchantRefundNotificationWithoutMerchantBody, businessName, amount, businessName, accountNumber)
		m.TransactionTitle = fmt.Sprintf(CreditMerchantRefundTransactionWithoutMerchantTitle)
	}

	m.ActivityType = activity.TypeMerchantRefundCredit
	m.BusinessName = businessName

	return &m, nil
}

func processMoneyTransferTransaction(uID shared.UserID, bID shared.BusinessID, accID *string, contactID *string, moneyTransferID *string,
	requestID *shared.PaymentRequestID, monthlyInterestID *string, amt num.Decimal, action string, notification TransactionNotification) (*TransactionMessage, error) {

	var m *TransactionMessage
	var err error

	if notification.CodeType == transaction.TransactionCodeTypeCreditPosted {
		if requestID != nil && *requestID != "" {
			m, err = processMoneyRequestCreditTransaction(bID, accID, contactID, requestID, amt, notification.TransactionDate)
			if err != nil {
				log.Println(err)
				return nil, err
			}
		} else if monthlyInterestID != nil && *monthlyInterestID != "" {
			m, err = processTransferInterestCreditTransaction(bID, accID, monthlyInterestID, amt)
			if err != nil {
				log.Println(err)
				return nil, err
			}
		} else {
			switch notification.TransactionType {
			case TransactionTypeTransfer:
				m, err = processTransferCreditTransaction(bID, accID, amt, notification.TransactionDate, *notification.BankTransactionDesc)
				if err != nil {
					log.Println(err)
					return nil, err
				}
			case TransactionTypeDeposit:
				m, err = processDepositCreditTransaction(bID, accID, amt, notification.TransactionDate, *notification.BankTransactionDesc)
				if err != nil {
					log.Println(err)
					return nil, err
				}
			}
		}

	} else if notification.CodeType == transaction.TransactionCodeTypeDebitPosted {
		if moneyTransferID != nil && *moneyTransferID != "" {
			m, err = processTransferDebitTransaction(bID, accID, contactID, moneyTransferID, amt, notification.BankTransactionDesc)
			if err != nil {
				log.Println(err)
				return nil, err
			}

			c, err := contact.NewMoneyTransferService(services.NewSourceRequest()).GetContactByTransferId(*notification.BankMoneyTransferID, bID)
			if err != nil && err != sql.ErrNoRows {
				log.Println("Error fetching contact for debit posted transfer ", err, *notification.BankMoneyTransferID, bID)
			}

			if err == nil {
				user, err := usr.NewUserServiceWithout().GetByIdInternal(c.Contact.UserID)
				if err != nil {
					log.Println(err)
					return nil, err
				}

				if c.SendEmail {
					// Send email to sender and receiver
					sendEmail(bID, c, user, amt, m.ActivityType)
				}
			}
		}
	}

	// Check if type is transfer
	if notification.TransactionType == TransactionTypeTransfer {
		// Update money transfer status
		business.NewMoneyTransferService(services.NewSourceRequest()).UpdateStatus(bID,
			*notification.BankMoneyTransferID, action)
	}

	fAmount, _ := amt.Abs().Float64()

	// Add to activity stream
	t := activity.AccountTransaction{
		EntityID:            string(bID),
		UserID:              uID,
		Amount:              activity.AccountTransactionAmount(fAmount),
		ContactName:         m.SenderReceiverName,
		TransactionDate:     notification.TransactionDate,
		BusinessName:        &m.BusinessName,
		InterestEarnedMonth: m.InterestEarnedMonth,
	}

	ID := onMoneyTransfer(t, m.ActivityType, notification.CodeType)
	m.ActivtiyID = ID

	return m, nil
}

func processACHTransaction(uID shared.UserID, bID shared.BusinessID, accID *string, contactID *string, moneyTransferID *string, requestID *shared.PaymentRequestID,
	amt num.Decimal, action string, notification TransactionNotification) (*TransactionMessage, error) {
	var name *string

	var m *TransactionMessage
	var err error

	isSnapcheck := false

	if notification.CodeType == transaction.TransactionCodeTypeCreditPosted {

		if requestID != nil && *requestID != "" {
			m, err = processMoneyRequestCreditTransaction(bID, accID, contactID, requestID, amt, notification.TransactionDate)
			if err != nil {
				return nil, err
			}
		} else {
			//Check for a snapcheck transaction
			isSnapcheck, err = isSnapcheckTransaction(notification.BankTransactionDesc)
			if err != nil {
				return nil, err
			}

			if isSnapcheck {
				m, err = processCheckCreditTransaction(bID, accID, amt, notification.TransactionDate, *notification.BankTransactionDesc)
				if err != nil {
					return nil, err
				}
			} else {
				m, err = processACHCreditTransaction(bID, accID, moneyTransferID, amt, *notification.BankTransactionDesc, notification.TransactionDate)
				if err != nil {
					return nil, err
				}
			}
		}

		name = m.SenderReceiverName

	} else if notification.CodeType == transaction.TransactionCodeTypeDebitPosted {
		// ACH transaction here
		if contactID != nil && *contactID != "" {
			if transaction.TransactionType(notification.TransactionType) == transaction.TransactionTypePurchase {
				m, err = processPurchaseDebitTransaction(bID, accID, contactID, moneyTransferID, amt)
				if err != nil {
					return nil, err
				}
			} else {
				m, err = processACHContactDebitTransaction(bID, accID, contactID, moneyTransferID, amt)
				if err != nil {
					return nil, err
				}
			}

			c, err := contact.NewMoneyTransferService(services.NewSourceRequest()).GetContactByTransferId(*notification.BankMoneyTransferID, bID)
			if err == nil {
				user, err := usr.NewUserServiceWithout().GetByIdInternal(c.Contact.UserID)
				if err != nil {
					log.Println(err)
					return nil, err
				}

				if c.SendEmail {
					// Send email to sender and receiver
					sendEmail(bID, c, user, amt, m.ActivityType)
				}
			}

			name = m.SenderReceiverName
		} else if moneyTransferID != nil && *moneyTransferID != "" {
			m, err = processACHTransferDebitTransaction(bID, accID, moneyTransferID, amt, *notification.BankTransactionDesc)
			if err != nil {
				return nil, err
			}

			name = m.SenderReceiverName
		} else {
			m, err = processACHExternalDebitTransaction(bID, accID, amt, *notification.BankTransactionDesc)
			if err != nil {
				return nil, err
			}

			name = m.SenderReceiverName
		}

	}

	if isSnapcheck {
		bts, err := business.NewBankingTransferService()
		if err != nil {
			return nil, err
		}

		if notification.BankAccountID != nil {
			pts := []grpcBankTransfer.PartnerTransferStatus{
				grpcBankTransfer.PartnerTransferStatus_TPS_SNAPCHECK_DEPOSITED,
			}

			mt, err := bts.GetByAccountIDPartnerTransferStatusAndAmount(*accID, pts, amt)
			if err != nil {
				return nil, err
			}

			notification.BankMoneyTransferID = &mt.BankTransferId
		} else {
			return nil, errors.New("Snapcheck notification sent without BankAccountID")
		}
	}

	// Check if type is transfer
	if notification.BankMoneyTransferID != nil {
		// Update money transfer status
		business.NewMoneyTransferService(services.NewSourceRequest()).UpdateStatus(bID,
			*notification.BankMoneyTransferID, action)
	}

	fAmount, _ := amt.Abs().Float64()

	// Add to activity stream
	t := activity.AccountTransaction{
		EntityID:            string(bID),
		UserID:              uID,
		Amount:              activity.AccountTransactionAmount(fAmount),
		ContactName:         name,
		TransactionDate:     notification.TransactionDate,
		BusinessName:        &m.BusinessName,
		InterestEarnedMonth: nil,
	}

	ID := onMoneyTransfer(t, m.ActivityType, notification.CodeType)
	m.ActivtiyID = ID

	return m, nil
}

func sendEmail(bID shared.BusinessID, c *contact.ContactTransferDetails, user *usr.User, amt num.Decimal, activityType activity.Type) error {

	// Send email to recepient
	date := time.Now().Format("Jan _2, 2006")

	var name string
	switch c.Contact.Type {
	case con.ContactTypeBusiness:
		name = *c.Contact.BusinessName
	case con.ContactTypePerson:
		name = *c.Contact.FirstName + " " + *c.Contact.LastName
	default:
		log.Println("Invalid contact type ", c.Contact.Type)
		return errors.New("Invalid contact type")
	}

	_, err := sendgrid.NewSendGridServiceWithout().SendEmail(sendEmailToRecepient(c.Contact.Email, name,
		date, c.Notes, amt.Abs().FormatCurrency(), activityType))
	if err != nil {
		log.Println(err)
		return err
	}

	// Send email to sender
	_, err = sendgrid.NewSendGridServiceWithout().SendEmail(sendEmailToSender(*user.Email,
		user.FirstName+" "+user.LastName, date, c.Notes, amt.Abs().FormatCurrency(), c.AccountNumber))
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// Send email to recepient
func sendEmailToRecepient(rEmail string, rName string, date string, description *string, amount string, activityType activity.Type) sendgrid.EmailRequest {

	var desc = ""
	if description != nil {
		desc = *description
	}

	var body string
	if activityType == activity.TypeCheckDebit {
		body = fmt.Sprintf(services.CheckPaymentReceivedEmail, rName, date, amount, desc)
	} else {
		body = fmt.Sprintf(services.PaymentReceivedEmail, rName, date, amount, desc)
	}

	return sendgrid.EmailRequest{
		SenderEmail:   os.Getenv("WISE_SUPPORT_EMAIL"),
		SenderName:    os.Getenv("WISE_SUPPORT_NAME"),
		ReceiverEmail: rEmail,
		ReceiverName:  rName,
		Subject:       services.PaymentReceivedSubject,
		Body:          body,
	}
}

// Send email to recepient
func sendEmailToSender(rEmail string, rName string, date string, description *string, amount string, accountNumber *string) sendgrid.EmailRequest {

	var desc = ""
	if description != nil {
		desc = *description
	}

	body := fmt.Sprintf(services.PaymentSentEmail, rName, services.MaskLeft(*accountNumber, 4), date, rName, amount, desc)

	return sendgrid.EmailRequest{
		SenderEmail:   os.Getenv("WISE_SUPPORT_EMAIL"),
		SenderName:    os.Getenv("WISE_SUPPORT_NAME"),
		ReceiverEmail: rEmail,
		ReceiverName:  rName,
		Subject:       services.PaymentSentSubject,
		Body:          body,
	}
}

// business name, amount, method, location, account number, date
func processMoneyRequestCreditTransaction(bID shared.BusinessID, accID *string, contactID *string, requestID *shared.PaymentRequestID, amt num.Decimal, txnDate time.Time) (*TransactionMessage, error) {
	m := TransactionMessage{}

	amount := amt.Abs().FormatCurrency()

	// Get business name
	type BusinessDetails struct {
		AccountNumber      string               `db:"business_bank_account.account_number"`
		LegalName          string               `db:"business.legal_name"`
		DBA                services.StringArray `db:"business.dba"`
		ContactName        string
		MoneyRequestMethod payment.PaymentRequestType `db:"business_money_request.request_type"`
		PaymentLocation    *services.Address          `db:"business_money_request_payment.purchase_address"`
	}

	bt := BusinessDetails{}

	err := core.DBRead.Get(
		&bt, `
	    SELECT
    	    legal_name "business.legal_name",
        	dba "business.dba"
	    FROM business
    	WHERE id = $1`,
		bID,
	)
	if err != nil {
		log.Println("Error retrieving business details 10:", err)
		return nil, err
	}

	if contactID != nil && *contactID != "" {
		cont, err := con.NewContactServiceWithout().GetByIDInternal(*contactID)
		if err != nil {
			log.Println("Error retrieving business details 11:", err)
			return nil, err
		}

		bt.ContactName = cont.Name()
	}

	if requestID != nil && *requestID != "" {
		isPOS := isPOSRequest(requestID)
		if os.Getenv("USE_INVOICE_SERVICE") == "true" && !isPOS {
			err := core.DBRead.Get(
				&bt,
				`SELECT
				business_money_request_payment.purchase_address "business_money_request_payment.purchase_address"
				FROM business_money_request_payment WHERE business_money_request_payment.invoice_id = $1`, *requestID,
			)
			if err != nil {
				log.Println("Error while fetching the invoice from business details:", err)
				return nil, err
			}
			invoiceDetail, err := getInvoiceFromInvoiceId(requestID)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			if invoiceDetail.AllowCard {
				bt.MoneyRequestMethod = payment.PaymentRequestTypeInvoiceCard
			}

			log.Println("Updated money request object with invoice details from invoice service")
		} else {
			err := core.DBRead.Get(
				&bt,
				`SELECT
				business_money_request_payment.purchase_address "business_money_request_payment.purchase_address",
				business_money_request.request_type "business_money_request.request_type"
			FROM business_money_request
			JOIN business_money_request_payment ON business_money_request.id = business_money_request_payment.request_id
			WHERE business_money_request.business_id = $1 AND business_money_request_payment.request_id = $2`,
				bID, *requestID,
			)
			if err != nil {
				log.Println("Error retrieving business details 11:", err)
				return nil, err
			}
		}
	}
	if err != nil {
		log.Println("Error retrieving business details 5:", err)
		return nil, err
	}

	if accID != nil && *accID != "" {
		acc, err := business.NewAccountService().GetByIDInternal(*accID)
		if err != nil {
			log.Println("Error retrieving account details ", err)
			return nil, err
		}

		bt.AccountNumber = acc.AccountNumber
	}

	var requestType string
	switch bt.MoneyRequestMethod {
	case payment.PaymentRequestTypePOS:
		requestType = "Card Reader"

		businessName := shared.GetBusinessName(&bt.LegalName, bt.DBA)

		accountNumber := string(bt.AccountNumber[len(bt.AccountNumber)-4:])

		if bt.PaymentLocation != nil {
			m.NotificationHeader = fmt.Sprintf(CreditViaCardReaderNotificationHeader, businessName, amount)
			m.NotificationBody = fmt.Sprintf(CreditViaCardReaderNotificationBody, businessName, requestType, bt.PaymentLocation.City, amount, businessName, accountNumber)
			m.TransactionTitle = fmt.Sprintf(CreditViaCardReaderTransactionTitle, businessName)
			m.TransactionDescription = fmt.Sprintf(CreditViaCardReaderTransactionDescription, businessName, requestType, bt.PaymentLocation.City)
			m.ActivityType = activity.TypeCardReaderCredit
			m.BusinessName = businessName

			return &m, nil
		}
	case payment.PaymentRequestTypeInvoiceBank:
		requestType = "Bank Transfer"

		businessName := shared.GetBusinessName(&bt.LegalName, bt.DBA)

		accountNumber := string(bt.AccountNumber[len(bt.AccountNumber)-4:])

		m.Counterparty = bt.ContactName
		m.NotificationHeader = fmt.Sprintf(CreditViaBankNotificationHeader, businessName, amount)
		m.NotificationBody = fmt.Sprintf(CreditViaBankNotificationBody, bt.ContactName, businessName, requestType, amount, businessName, accountNumber)
		m.TransactionTitle = fmt.Sprintf(CreditViaBankTransactionTitle, bt.ContactName, businessName)
		m.TransactionDescription = fmt.Sprintf(CreditViaBankTransactionDescription, bt.ContactName, businessName, requestType)
		m.ActivityType = activity.TypeBankOnlineCredit
		m.SenderReceiverName = &bt.ContactName
		m.BusinessName = businessName

		return &m, nil
	default:
		requestType = "Card"

		businessName := shared.GetBusinessName(&bt.LegalName, bt.DBA)

		accountNumber := string(bt.AccountNumber[len(bt.AccountNumber)-4:])

		m.Counterparty = bt.ContactName
		m.NotificationHeader = fmt.Sprintf(CreditViaCardNotificationHeader, businessName, amount)
		m.NotificationBody = fmt.Sprintf(CreditViaCardNotificationBody, bt.ContactName, businessName, requestType, amount, businessName, accountNumber)
		m.TransactionTitle = fmt.Sprintf(CreditViaCardTransactionTitle, bt.ContactName, businessName)
		m.TransactionDescription = fmt.Sprintf(CreditViaCardTransactionDescription, bt.ContactName, businessName, requestType)
		m.ActivityType = activity.TypeCardOnlineCredit
		m.SenderReceiverName = &bt.ContactName
		m.BusinessName = businessName

		return &m, nil
	}

	return nil, errors.New("Error in processMoneyRequestCreditTransaction")
}

func processTransferCreditTransaction(bID shared.BusinessID, accID *string, amt num.Decimal, txnDate time.Time, transferDesc string) (*TransactionMessage, error) {
	m := TransactionMessage{}

	amount := amt.Abs().FormatCurrency()

	// Get business name
	type BusinessDetails struct {
		AccountNumber string               `db:"business_bank_account.account_number"`
		LegalName     string               `db:"business.legal_name"`
		DBA           services.StringArray `db:"business.dba"`
	}

	bt := BusinessDetails{}

	err := core.DBRead.Get(
		&bt, `
        SELECT
            legal_name "business.legal_name",
            dba "business.dba"
        FROM business
        WHERE id = $1`,
		bID,
	)
	if err != nil {
		log.Println("Error retrieving business details 6:", err)
		return nil, err
	}

	if accID != nil && *accID != "" {
		acc, err := business.NewAccountService().GetByIDInternal(*accID)
		if err != nil {
			log.Println("Error retrieving account details ", err)
			return nil, err
		}

		bt.AccountNumber = acc.AccountNumber
	}

	businessName := shared.GetBusinessName(&bt.LegalName, bt.DBA)

	accountNumber := string(bt.AccountNumber[len(bt.AccountNumber)-4:])

	senderName := transaction.GetOriginAccountHolder(transaction.TransactionTypeTransfer, transferDesc)
	m.Counterparty = shared.StringValue(senderName)

	m.NotificationHeader = fmt.Sprintf(CreditWiseTransferNotificationHeader, businessName, amount)
	m.NotificationBody = fmt.Sprintf(CreditWiseTransferNotificationBody, businessName, *senderName, amount, businessName, accountNumber)
	m.TransactionTitle = fmt.Sprintf(CreditWiseTransferTransactionTitle, *senderName)
	m.TransactionDescription = fmt.Sprintf(CreditWiseTransferTransactionDescription, amount, *senderName)

	m.SenderReceiverName = senderName
	m.ActivityType = activity.TypeWiseTransferCredit

	return &m, nil
}

func processACHCreditTransaction(bID shared.BusinessID, accID *string, moneyTransferID *string, amt num.Decimal,
	transferDesc string, transactionDate time.Time) (*TransactionMessage, error) {
	m := TransactionMessage{}

	amount := amt.Abs().FormatCurrency()

	// Get business name
	type BusinessDetails struct {
		AccountNumber string               `db:"business_bank_account.account_number"`
		LegalName     string               `db:"business.legal_name"`
		DBA           services.StringArray `db:"business.dba"`
	}

	bt := BusinessDetails{}

	err := core.DBRead.Get(
		&bt, `
        SELECT
            legal_name "business.legal_name",
            dba "business.dba"
        FROM business
        WHERE id = $1`,
		bID,
	)
	if err != nil {
		log.Println("Error retrieving business details 7:", err)
		return nil, err
	}

	if accID != nil && *accID != "" {
		acc, err := business.NewAccountService().GetByIDInternal(*accID)
		if err != nil {
			log.Println("Error retrieving account details ", err)
			return nil, err
		}

		bt.AccountNumber = acc.AccountNumber
	}

	businessName := shared.GetBusinessName(&bt.LegalName, bt.DBA)
	accountNumber := string(bt.AccountNumber[len(bt.AccountNumber)-4:])

	var senderName *string
	if moneyTransferID != nil && *moneyTransferID != "" {
		senderName = transaction.GetOriginAccountHolder(transaction.TransactionTypeACH, transferDesc)
	} else {
		senderName = transaction.GetExternalACHSource(transferDesc)
	}
	m.Counterparty = shared.StringValue(senderName)

	m.NotificationHeader = fmt.Sprintf(CreditACHNotificationHeader, businessName, amount)
	if len(*senderName) > 0 {
		m.NotificationBody = fmt.Sprintf(CreditACHNotificationBody, businessName, *senderName, amount, businessName, accountNumber)
		m.TransactionTitle = fmt.Sprintf(CreditACHTransactionTitle, *senderName)
		m.TransactionDescription = fmt.Sprintf(CreditACHTransactionDescription, amount, *senderName)
	} else {
		m.NotificationBody = fmt.Sprintf(CreditACHNotificationWithoutSenderBody, businessName, amount, businessName, accountNumber)
		m.TransactionTitle = fmt.Sprintf(CreditACHTransactionWithoutSenderTitle)
		m.TransactionDescription = fmt.Sprintf(CreditACHTransactionWithoutSenderDescription, amount)
	}

	m.SenderReceiverName = senderName
	m.ActivityType = activity.TypeACHTransferCredit
	m.BusinessName = businessName

	// check for shopify transaction
	if isShopifyTransaction(senderName) {
		m.ActivityType = activity.TypeACHTransferShopifyCredit

		m.NotificationBody = fmt.Sprintf(CreditACHShopifyNotificationBody, businessName, *senderName, amount, businessName, accountNumber)
		m.TransactionTitle = fmt.Sprintf(CreditACHShopifyTransactionTitle, *senderName)
		m.TransactionDescription = fmt.Sprintf(CreditACHShopifyTransactionDescription, amount, *senderName)
	}

	return &m, nil
}

func processDepositCreditTransaction(bID shared.BusinessID, accID *string, amt num.Decimal, txnDate time.Time, transferDesc string) (*TransactionMessage, error) {

	wireOriginAccount := transaction.GetOriginAccountHolder(transaction.TransactionTypeDeposit, transferDesc)
	if wireOriginAccount != nil && len(*wireOriginAccount) > 0 {
		return processWireCreditTransaction(bID, accID, amt, txnDate, transferDesc)
	}

	isCheckDeposit := isCheckDeposit(transferDesc)
	if isCheckDeposit {
		return processCheckCreditTransaction(bID, accID, amt, txnDate, transferDesc)
	}

	return processGenericDepositCreditTransaction(bID, accID, amt, txnDate)
}

func processWireCreditTransaction(bID shared.BusinessID, accID *string, amt num.Decimal, txnDate time.Time, transferDesc string) (*TransactionMessage, error) {
	m := TransactionMessage{}

	amount := amt.Abs().FormatCurrency()

	// Get business name
	type BusinessDetails struct {
		AccountNumber string               `db:"business_bank_account.account_number"`
		LegalName     string               `db:"business.legal_name"`
		DBA           services.StringArray `db:"business.dba"`
	}

	bt := BusinessDetails{}

	err := core.DBRead.Get(
		&bt, `
        SELECT
            legal_name "business.legal_name",
            dba "business.dba"
        FROM business
        WHERE id = $1`,
		bID,
	)
	if err != nil {
		log.Println("Error retrieving business details 8:", err)
		return nil, err
	}

	if accID != nil && *accID != "" {
		acc, err := business.NewAccountService().GetByIDInternal(*accID)
		if err != nil {
			log.Println("Error retrieving account details ", err)
			return nil, err
		}

		bt.AccountNumber = acc.AccountNumber
	}

	businessName := shared.GetBusinessName(&bt.LegalName, bt.DBA)
	accountNumber := string(bt.AccountNumber[len(bt.AccountNumber)-4:])
	senderName := transaction.GetOriginAccountHolder(transaction.TransactionTypeDeposit, transferDesc)
	m.Counterparty = shared.StringValue(senderName)

	m.NotificationHeader = fmt.Sprintf(CreditDepositNotificationHeader, businessName, amount)
	m.NotificationBody = fmt.Sprintf(CreditDepositWireNotificationBody, businessName, *senderName, amount, businessName, accountNumber)

	m.TransactionTitle = fmt.Sprintf(CreditDepositWireTransactionTitle, *senderName)
	m.TransactionDescription = fmt.Sprintf(CreditDepositWireTransactionDescription, amount, *senderName)

	m.SenderReceiverName = senderName
	m.ActivityType = activity.TypeWireTransferCredit
	m.BusinessName = businessName

	return &m, nil
}

func processCheckCreditTransaction(bID shared.BusinessID, accID *string, amt num.Decimal, txnDate time.Time, transferDesc string) (*TransactionMessage, error) {
	m := TransactionMessage{}

	amount := amt.Abs().FormatCurrency()

	// Get business name
	type BusinessDetails struct {
		AccountNumber string               `db:"business_bank_account.account_number"`
		LegalName     string               `db:"business.legal_name"`
		DBA           services.StringArray `db:"business.dba"`
	}

	bt := BusinessDetails{}

	err := core.DBRead.Get(
		&bt,
		`SELECT business.legal_name "business.legal_name", business.dba "business.dba" 
		FROM business 
		WHERE business.id = $1`,
		bID,
	)
	if err != nil {
		log.Println("Error retrieving business details ", err)
		return nil, err
	}

	if accID != nil && *accID != "" {
		acc, err := business.NewAccountService().GetByIDInternal(*accID)
		if err != nil {
			log.Println("Error retrieving account details ", err)
			return nil, err
		}

		bt.AccountNumber = acc.AccountNumber
	}

	businessName := shared.GetBusinessName(&bt.LegalName, bt.DBA)
	accountNumber := string(bt.AccountNumber[len(bt.AccountNumber)-4:])

	m.NotificationHeader = fmt.Sprintf(CreditDepositNotificationHeader, businessName, amount)
	m.NotificationBody = fmt.Sprintf(CreditDepositCheckNotificationBody, amount, businessName, accountNumber)

	m.TransactionTitle = fmt.Sprintf(CreditDepositCheckTransactionTitle)
	m.TransactionDescription = fmt.Sprintf(CreditDepositCheckTransactionDescription, amount)

	m.ActivityType = activity.TypeCheckCredit
	m.BusinessName = businessName

	return &m, nil
}

func processGenericDepositCreditTransaction(bID shared.BusinessID, accID *string, amt num.Decimal, txnDate time.Time) (*TransactionMessage, error) {
	m := TransactionMessage{}

	amount := amt.Abs().FormatCurrency()

	// Get business name
	type BusinessDetails struct {
		AccountNumber string               `db:"business_bank_account.account_number"`
		LegalName     string               `db:"business.legal_name"`
		DBA           services.StringArray `db:"business.dba"`
	}

	bt := BusinessDetails{}

	err := core.DBRead.Get(
		&bt,
		`SELECT business.legal_name "business.legal_name", business.dba "business.dba"
		FROM business 
		WHERE business.id = $1`,
		bID,
	)
	if err != nil {
		log.Println("Error retrieving business details ", err)
		return nil, err
	}

	if accID != nil && *accID != "" {
		acc, err := business.NewAccountService().GetByIDInternal(*accID)
		if err != nil {
			log.Println("Error retrieving account details ", err)
			return nil, err
		}

		bt.AccountNumber = acc.AccountNumber
	}

	businessName := shared.GetBusinessName(&bt.LegalName, bt.DBA)
	accountNumber := string(bt.AccountNumber[len(bt.AccountNumber)-4:])

	m.NotificationHeader = fmt.Sprintf(CreditDepositNotificationHeader, businessName, amount)
	m.NotificationBody = fmt.Sprintf(CreditDepositNotificationBody, businessName, amount, businessName, accountNumber)

	m.TransactionTitle = fmt.Sprintf(CreditDepositTransactionTitle)
	m.TransactionDescription = fmt.Sprintf(CreditDepositTransactionDescription, amount)

	m.ActivityType = activity.TypeDepositCredit
	m.BusinessName = businessName

	return &m, nil
}

func processTransferInterestCreditTransaction(bID shared.BusinessID, accID *string, monthlyInterestID *string, amt num.Decimal) (*TransactionMessage, error) {
	m := TransactionMessage{}

	amount := amt.Abs().FormatCurrency()

	// Get business name
	type BusinessDetails struct {
		AccountNumber string               `db:"business_bank_account.account_number"`
		LegalName     string               `db:"business.legal_name"`
		DBA           services.StringArray `db:"business.dba"`
	}

	type InterestDetails struct {
		StartDate shared.Date `db:"business_account_monthly_interest.start_date"`
		EndDate   shared.Date `db:"business_account_monthly_interest.end_date"`
	}

	bt := BusinessDetails{}
	i := InterestDetails{}

	err := core.DBRead.Get(
		&bt, `
        SELECT
            legal_name "business.legal_name",
            dba "business.dba"
        FROM business
        WHERE id = $1`,
		bID,
	)
	if err != nil {
		log.Println("Error retrieving business details 9:", err)
		return nil, err
	}

	if accID != nil && *accID != "" {
		acc, err := business.NewAccountService().GetByIDInternal(*accID)
		if err != nil {
			log.Println("Error retrieving account details ", err)
			return nil, err
		}

		bt.AccountNumber = acc.AccountNumber
	}

	err = transaction.DBRead.Get(
		&i, `
		SELECT
			business_account_monthly_interest.start_date "business_account_monthly_interest.start_date", 
			business_account_monthly_interest.end_date "business_account_monthly_interest.end_date"
		FROM business_account_monthly_interest 
		WHERE id = $1`,
		monthlyInterestID)

	businessName := shared.GetBusinessName(&bt.LegalName, bt.DBA)

	accountNumber := string(bt.AccountNumber[len(bt.AccountNumber)-4:])
	date := i.StartDate.Time().Format("Jan 2006")

	m.NotificationHeader = fmt.Sprintf(CreditInterestTransferNotificationHeader, businessName, amount)
	m.NotificationBody = fmt.Sprintf(CreditInterestTransferNotificationBody, businessName, accountNumber, amount, date)
	m.TransactionTitle = fmt.Sprintf(CreditInterestTransferTransactionTitle, date)
	m.TransactionDescription = fmt.Sprintf(CreditInterestTransferTransactionDescription, businessName, amount, date)

	m.ActivityType = activity.TypeInterestTransferCredit
	m.BusinessName = businessName
	m.InterestEarnedMonth = &date

	return &m, nil
}

func processTransferDebitTransaction(bID shared.BusinessID, accID *string, contactID *string, moneyTransferID *string,
	amt num.Decimal, transferDesc *string) (*TransactionMessage, error) {
	m := TransactionMessage{}

	amount := amt.Abs().FormatCurrency()

	// Get business name
	type BusinessDetails struct {
		AccountNumber string               `db:"business_bank_account.account_number"`
		LegalName     string               `db:"business.legal_name"`
		DBA           services.StringArray `db:"business.dba"`
		ContactName   string
		DestType      banking.TransferType `db:"business_money_transfer.dest_type"`
		Created       time.Time            `db:"business_money_transfer.created"`
	}

	bt := BusinessDetails{}

	err := core.DBRead.Get(
		&bt, `
        SELECT
            legal_name "business.legal_name",
            dba "business.dba"
        FROM business
        WHERE id = $1`,
		bID,
	)
	if err != nil {
		log.Println("Error retrieving business details 10:", err)
		return nil, err
	}

	if contactID != nil && *contactID != "" {
		cont, err := con.NewContactServiceWithout().GetByIDInternal(*contactID)
		if err != nil {
			log.Println("Error retrieving business details 11:", err)
			return nil, err
		}

		bt.ContactName = cont.Name()
	}

	log.Println("Contact:", contactID, bt.ContactName)

	if accID != nil && *accID != "" {
		acc, err := business.NewAccountService().GetByIDInternal(*accID)
		if err != nil {
			log.Println("Error retrieving account details ", err)
			return nil, err
		}

		bt.AccountNumber = acc.AccountNumber
	}

	if moneyTransferID != nil && *moneyTransferID != "" {
		mtr, err := business.NewMoneyTransferServiceWithout().GetByIDOnlyInternal(*moneyTransferID)
		if err != nil {
			log.Println("Error retrieving business details 13:", err)
			return nil, err
		}

		bt.DestType = mtr.DestType
		bt.Created = mtr.Created
	}

	businessName := shared.GetBusinessName(&bt.LegalName, bt.DBA)

	accountNumber := string(bt.AccountNumber[len(bt.AccountNumber)-4:])

	m.Counterparty = bt.ContactName
	if banking.TransferType(bt.DestType) == banking.TransferTypeCheck {
		m.NotificationHeader = fmt.Sprintf(DebitWiseCheckNotificationHeader, businessName, bt.ContactName, amount)
		m.NotificationBody = fmt.Sprintf(DebitWiseCheckNotificationBody, businessName, bt.ContactName, amount)
		m.TransactionTitle = fmt.Sprintf(DebitWiseCheckTransactionTitle, bt.ContactName)
		m.TransactionDescription = fmt.Sprintf(DebitWiseCheckTransactionDescription, businessName, bt.ContactName)

		m.ActivityType = activity.TypeCheckDebit

	} else {
		// extract contact name from description
		if bt.ContactName == "" {
			destName := transaction.GetDestinationAccountHolder(transaction.TransactionTypeTransfer, shared.StringValue(transferDesc))
			if destName != nil {
				bt.ContactName = *destName
			}
		}

		m.NotificationHeader = fmt.Sprintf(DebitWiseTransferNotificationHeader, businessName, amount)
		if bt.ContactName != "" {
			m.NotificationBody = fmt.Sprintf(DebitWiseTransferNotificationBody, businessName, bt.ContactName, amount, businessName, accountNumber)
			m.TransactionTitle = fmt.Sprintf(DebitWiseTransferTransactionTitle, bt.ContactName)
			m.TransactionDescription = fmt.Sprintf(DebitWiseTransferTransactionDescription, businessName, bt.ContactName)
		} else {
			m.NotificationBody = fmt.Sprintf(DebitWiseTransferNotificationWithoutContactBody, businessName, amount, businessName, accountNumber)
			m.TransactionTitle = DebitWiseTransferTransactionWithoutContactTitle
			m.TransactionDescription = fmt.Sprintf(DebitWiseTransferTransactionWithoutContactDescription, businessName)
		}

		m.ActivityType = activity.TypeWiseTransferDebit
	}

	m.SenderReceiverName = &bt.ContactName
	m.BusinessName = businessName
	return &m, nil
}

func processACHTransferDebitTransaction(bID shared.BusinessID, accID *string, moneyTransferID *string, amt num.Decimal, transferDesc string) (*TransactionMessage, error) {
	m := TransactionMessage{}

	amount := amt.Abs().FormatCurrency()

	// Get business name
	type BusinessDetails struct {
		AccountNumber string               `db:"business_bank_account.account_number"`
		LegalName     string               `db:"business.legal_name"`
		DBA           services.StringArray `db:"business.dba"`
	}

	bt := BusinessDetails{}

	err := core.DBRead.Get(
		&bt, `
        SELECT
            legal_name "business.legal_name",
            dba "business.dba"
        FROM business
        WHERE id = $1`,
		bID,
	)
	if err != nil {
		log.Println("Error retrieving business details 14:", err)
		return nil, err
	}

	if accID != nil && *accID != "" {
		acc, err := business.NewAccountService().GetByIDInternal(*accID)
		if err != nil {
			log.Println("Error retrieving account details ", err)
			return nil, err
		}

		bt.AccountNumber = acc.AccountNumber
	}

	businessName := shared.GetBusinessName(&bt.LegalName, bt.DBA)

	accountNumber := string(bt.AccountNumber[len(bt.AccountNumber)-4:])
	receiverName := transaction.GetDestinationAccountHolder(transaction.TransactionTypeACH, transferDesc)
	m.Counterparty = shared.StringValue(receiverName)

	m.NotificationHeader = fmt.Sprintf(DebitACHNotificationHeader, businessName, amount)
	m.NotificationBody = fmt.Sprintf(DebitACHNotificationBody, businessName, *receiverName, amount, businessName, accountNumber)
	m.TransactionTitle = fmt.Sprintf(DebitACHTransactionTitle, *receiverName)
	m.TransactionDescription = fmt.Sprintf(DebitACHTransactionDescription, businessName, *receiverName)

	m.SenderReceiverName = receiverName
	m.ActivityType = activity.TypeACHTransferDebit
	m.BusinessName = businessName

	return &m, nil
}

func processACHExternalDebitTransaction(bID shared.BusinessID, accID *string, amt num.Decimal, transferDesc string) (*TransactionMessage, error) {
	m := TransactionMessage{}

	amount := amt.Abs().FormatCurrency()

	// Get business name
	type BusinessDetails struct {
		AccountNumber string               `db:"business_bank_account.account_number"`
		LegalName     string               `db:"business.legal_name"`
		DBA           services.StringArray `db:"business.dba"`
	}

	bt := BusinessDetails{}

	err := core.DBRead.Get(
		&bt, `
        SELECT
            legal_name "business.legal_name",
            dba "business.dba"
        FROM business
        WHERE id = $1`,
		bID,
	)
	if err != nil {
		log.Println("Error retrieving business details 15:", err)
		return nil, err
	}

	if accID != nil && *accID != "" {
		acc, err := business.NewAccountService().GetByIDInternal(*accID)
		if err != nil {
			log.Println("Error retrieving account details ", err)
			return nil, err
		}

		bt.AccountNumber = acc.AccountNumber
	}

	businessName := shared.GetBusinessName(&bt.LegalName, bt.DBA)

	accountNumber := string(bt.AccountNumber[len(bt.AccountNumber)-4:])
	receiverName := transaction.GetExternalACHDestination(transferDesc)
	m.Counterparty = receiverName

	if len(receiverName) > 0 {
		m.NotificationHeader = fmt.Sprintf(DebitACHNotificationHeader, businessName, amount)
		m.NotificationBody = fmt.Sprintf(DebitACHNotificationBody, businessName, receiverName, amount, businessName, accountNumber)
		m.TransactionTitle = fmt.Sprintf(DebitACHTransactionTitle, receiverName)
		m.TransactionDescription = fmt.Sprintf(DebitACHTransactionDescription, businessName, receiverName)

		m.SenderReceiverName = &receiverName
	} else {
		m.NotificationHeader = fmt.Sprintf(DebitExternalACHNotificationHeader, businessName, amount)
		m.NotificationBody = fmt.Sprintf(DebitExternalACHNotificationBody, businessName, amount, businessName, accountNumber)
		m.TransactionTitle = fmt.Sprintf(DebitExternalACHTransactionTitle)
		m.TransactionDescription = fmt.Sprintf(DebitExternalACHTransactionDescription)
	}

	m.ActivityType = activity.TypeACHTransferDebit
	m.BusinessName = businessName

	return &m, nil
}

func processACHContactDebitTransaction(bID shared.BusinessID, accID *string, contactID *string,
	moneyTransferID *string, amt num.Decimal) (*TransactionMessage, error) {
	m := TransactionMessage{}

	amount := amt.Abs().FormatCurrency()

	// Get business name
	type BusinessDetails struct {
		AccountNumber string               `db:"business_bank_account.account_number"`
		LegalName     string               `db:"business.legal_name"`
		DBA           services.StringArray `db:"business.dba"`
		ContactName   string
		DestType      banking.TransferType `db:"dest_type"`
		Created       time.Time            `db:"business_money_transfer.created"`
	}

	bt := BusinessDetails{}

	err := core.DBRead.Get(
		&bt, `
        SELECT
            legal_name "business.legal_name",
            dba "business.dba"
        FROM business
        WHERE id = $1`,
		bID,
	)
	if err != nil {
		log.Println("Error retrieving business details 16:", err)
		return nil, err
	}

	if contactID != nil && *contactID != "" {
		cont, err := con.NewContactServiceWithout().GetByIDInternal(*contactID)
		if err != nil {
			log.Println("Error retrieving business details 11:", err)
			return nil, err
		}

		bt.ContactName = cont.Name()
	}

	if accID != nil && *accID != "" {
		acc, err := business.NewAccountService().GetByIDInternal(*accID)
		if err != nil {
			log.Println("Error retrieving account details ", err)
			return nil, err
		}

		bt.AccountNumber = acc.AccountNumber
	}

	if moneyTransferID != nil && *moneyTransferID != "" {
		mtr, err := business.NewMoneyTransferServiceWithout().GetByIDOnlyInternal(*moneyTransferID)
		if err != nil {
			log.Println("Error retrieving business details 13:", err)
			return nil, err
		}

		bt.DestType = mtr.DestType
	}

	businessName := shared.GetBusinessName(&bt.LegalName, bt.DBA)

	accountNumber := string(bt.AccountNumber[len(bt.AccountNumber)-4:])

	m.Counterparty = bt.ContactName
	m.NotificationHeader = fmt.Sprintf(DebitACHNotificationHeader, businessName, amount)
	m.NotificationBody = fmt.Sprintf(DebitACHNotificationBody, businessName, bt.ContactName, amount, businessName, accountNumber)
	m.TransactionTitle = fmt.Sprintf(DebitACHTransactionTitle, bt.ContactName)
	m.TransactionDescription = fmt.Sprintf(DebitACHTransactionDescription, businessName, bt.ContactName)

	m.SenderReceiverName = &bt.ContactName
	m.ActivityType = activity.TypeACHTransferDebit
	m.BusinessName = businessName

	return &m, nil
}

func processPurchaseDebitTransaction(bID shared.BusinessID, accID *string, contactID *string,
	moneyTransferID *string, amt num.Decimal) (*TransactionMessage, error) {
	m := TransactionMessage{}

	amount := amt.Abs().FormatCurrency()

	// Get business name
	type BusinessDetails struct {
		AccountNumber string               `db:"business_bank_account.account_number"`
		LegalName     string               `db:"business.legal_name"`
		DBA           services.StringArray `db:"business.dba"`
		ContactName   string
		DestType      banking.TransferType `db:"dest_type"`
		Created       time.Time            `db:"business_money_transfer.created"`
	}

	bt := BusinessDetails{}

	err := core.DBRead.Get(
		&bt, `
        SELECT
            legal_name "business.legal_name",
            dba "business.dba"
        FROM business
        WHERE id = $1`,
		bID,
	)
	if err != nil {
		log.Println("Error retrieving business details 20:", err)
		return nil, err
	}

	if contactID != nil && *contactID != "" {
		cont, err := con.NewContactServiceWithout().GetByIDInternal(*contactID)
		if err != nil {
			log.Println("Error retrieving business details 11:", err)
			return nil, err
		}
		bt.ContactName = cont.Name()
	}

	if accID != nil && *accID != "" {
		acc, err := business.NewAccountService().GetByIDInternal(*accID)
		if err != nil {
			log.Println("Error retrieving account details ", err)
			return nil, err
		}

		bt.AccountNumber = acc.AccountNumber
	}

	if moneyTransferID != nil && *moneyTransferID != "" {
		mtr, err := business.NewMoneyTransferServiceWithout().GetByIDOnlyInternal(*moneyTransferID)
		if err != nil {
			log.Println("Error retrieving business details 13:", err)
			return nil, err
		}

		bt.DestType = mtr.DestType
		bt.Created = mtr.Created
	}

	businessName := shared.GetBusinessName(&bt.LegalName, bt.DBA)

	accountNumber := string(bt.AccountNumber[len(bt.AccountNumber)-4:])

	m.SenderReceiverName = &bt.ContactName
	m.Counterparty = bt.ContactName

	if bt.DestType == banking.TransferTypeCard {
		m.NotificationHeader = fmt.Sprintf(DebitCardInstantPayNotificationHeader, businessName, amount)
		m.NotificationBody = fmt.Sprintf(DebitCardInstantPayNotificationBody, businessName, bt.ContactName, amount, businessName, accountNumber)
		m.TransactionTitle = fmt.Sprintf(DebitCardInstantPayTransactionTitle, bt.ContactName)
		m.TransactionDescription = fmt.Sprintf(DebitCardInstantPayTransactionDescription, businessName, bt.ContactName)

		m.ActivityType = activity.TypeCardPushDebit
	} else {
		m.NotificationHeader = fmt.Sprintf(DebitACHNotificationHeader, businessName, amount)
		m.NotificationBody = fmt.Sprintf(DebitACHNotificationBody, businessName, bt.ContactName, amount, businessName, accountNumber)
		m.TransactionTitle = fmt.Sprintf(DebitACHTransactionTitle, bt.ContactName)
		m.TransactionDescription = fmt.Sprintf(DebitACHTransactionDescription, businessName, bt.ContactName)

		m.ActivityType = activity.TypeACHTransferDebit
	}

	m.BusinessName = businessName
	return &m, nil
}

func isShopifyTransaction(counterparty *string) bool {
	if counterparty == nil || len(*counterparty) == 0 {
		return false
	}

	if strings.EqualFold(*counterparty, bus.PartnerNameShopify.String()) {
		return true
	}

	return false
}

func isSnapcheckTransaction(desc *string) (bool, error) {
	ret := false

	snapReg, err := regexp.Compile(`[0-9]+\s+CREDIT FOR SNAPCHECK\s+PAYMENT`)
	if err != nil {
		return ret, err
	}

	if desc != nil {
		found := snapReg.MatchString(*desc)
		if found {
			ret = true
		}
	}

	return ret, nil
}

func processDebitPullTransaction(uID shared.UserID, bID shared.BusinessID, accID *string, moneyTransferID *string,
	amt num.Decimal, action string, notification TransactionNotification) (*TransactionMessage, error) {
	var name *string

	var m *TransactionMessage
	var err error

	if notification.CodeType == transaction.TransactionCodeTypeCreditPosted {

		m, err = processDebitPullCreditTransaction(bID, accID, amt, *notification.BankTransactionDesc)
		if err != nil {
			return nil, err
		}

		name = m.SenderReceiverName
	}

	// Check if type is transfer
	if notification.BankMoneyTransferID != nil {
		// Update money transfer status
		business.NewMoneyTransferService(services.NewSourceRequest()).UpdateStatus(bID,
			*notification.BankMoneyTransferID, action)
	}

	fAmount, _ := amt.Abs().Float64()

	// Add to activity stream
	t := activity.AccountTransaction{
		EntityID:            string(bID),
		UserID:              uID,
		Amount:              activity.AccountTransactionAmount(fAmount),
		ContactName:         name,
		TransactionDate:     notification.TransactionDate,
		BusinessName:        &m.BusinessName,
		InterestEarnedMonth: nil,
	}

	ID := onMoneyTransfer(t, m.ActivityType, notification.CodeType)
	m.ActivtiyID = ID

	return m, nil
}

func processDebitPullCreditTransaction(bID shared.BusinessID, accID *string, amt num.Decimal, transferDesc string) (*TransactionMessage, error) {
	m := TransactionMessage{}

	amount := amt.Abs().FormatCurrency()

	// Get business name
	type BusinessDetails struct {
		AccountNumber string               `db:"business_bank_account.account_number"`
		LegalName     string               `db:"business.legal_name"`
		DBA           services.StringArray `db:"business.dba"`
	}

	bt := BusinessDetails{}

	err := core.DBRead.Get(
		&bt, `
        SELECT
            legal_name "business.legal_name",
            dba "business.dba"
        FROM business
        WHERE id = $1`,
		bID,
	)
	if err != nil {
		log.Println("Error retrieving business details 24", err)
		return nil, err
	}

	if accID != nil && *accID != "" {
		acc, err := business.NewAccountService().GetByIDInternal(*accID)
		if err != nil {
			log.Println("Error retrieving account details ", err)
			return nil, err
		}

		bt.AccountNumber = acc.AccountNumber
	}

	businessName := shared.GetBusinessName(&bt.LegalName, bt.DBA)

	accountNumber := string(bt.AccountNumber[len(bt.AccountNumber)-4:])

	cardNumber := transaction.GetOriginAccount(transaction.TransactionTypeACH, transferDesc)

	m.NotificationHeader = fmt.Sprintf(CreditCardDebitPullNotificationHeader, businessName, amount)
	if len(*cardNumber) > 0 {
		m.NotificationBody = fmt.Sprintf(CreditCardDebitPullNotificationBody, businessName, amount, *cardNumber, accountNumber)
		m.TransactionTitle = fmt.Sprintf(CreditCardDebitPullTransactionTitle, *cardNumber)
		m.TransactionDescription = fmt.Sprintf(CreditCardDebitPullTransactionDescription, amount)
	} else {
		m.NotificationBody = fmt.Sprintf(CreditCardDebitPullNotificationWithoutSenderBody, businessName, amount, accountNumber)
		m.TransactionTitle = fmt.Sprintf(CreditCardDebitPullTransactionWithoutSenderTitle)
		m.TransactionDescription = fmt.Sprintf(CreditCardDebitPullTransactionDescription, amount)
	}

	m.SenderReceiverName = cardNumber
	m.ActivityType = activity.TypeCardPullCredit
	m.BusinessName = businessName

	return &m, nil
}

func processOtherCreditTransaction(uID shared.UserID, bID shared.BusinessID, accID *string, moneyTransferID *string,
	amt num.Decimal, action string, notification TransactionNotification) (*TransactionMessage, error) {

	amount := amt.Abs().FormatCurrency()

	// Get business name
	type BusinessDetails struct {
		AccountNumber string               `db:"business_bank_account.account_number"`
		LegalName     string               `db:"business.legal_name"`
		DBA           services.StringArray `db:"business.dba"`
	}

	m := TransactionMessage{}
	bt := BusinessDetails{}

	err := core.DBRead.Get(
		&bt, `
        SELECT
            legal_name "business.legal_name",
            dba "business.dba"
        FROM business
        WHERE id = $1`,
		bID,
	)
	if err != nil {
		log.Println("Error retrieving business details 24", err)
		return nil, err
	}

	if accID != nil && *accID != "" {
		acc, err := business.NewAccountService().GetByIDInternal(*accID)
		if err != nil {
			log.Println("Error retrieving account details ", err)
			return nil, err
		}

		bt.AccountNumber = acc.AccountNumber
	}

	businessName := shared.GetBusinessName(&bt.LegalName, bt.DBA)
	accountNumber := string(bt.AccountNumber[len(bt.AccountNumber)-4:])

	m.NotificationHeader = fmt.Sprintf(OtherCreditNotificationHeader, businessName, amount)
	m.NotificationBody = fmt.Sprintf(OtherCreditNotificationBody, businessName, amount, businessName, accountNumber)
	m.TransactionTitle = OtherCreditTransactionTitle
	m.TransactionDescription = fmt.Sprintf(OtherCreditTransactionDescription, amount)

	m.ActivityType = activity.TypeOtherCredit
	m.BusinessName = businessName

	fAmount, _ := amt.Abs().Float64()

	// Add to activity stream
	t := activity.AccountTransaction{
		EntityID:            string(bID),
		UserID:              uID,
		Amount:              activity.AccountTransactionAmount(fAmount),
		TransactionDate:     notification.TransactionDate,
		BusinessName:        &m.BusinessName,
		InterestEarnedMonth: nil,
	}
	ID := onMoneyTransfer(t, m.ActivityType, notification.CodeType)
	m.ActivtiyID = ID

	return &m, nil
}

func isPOSRequest(requestID *shared.PaymentRequestID) bool {
	type Request struct {
		ID          shared.PaymentRequestID     `db:"id"`
		RequestType *payment.PaymentRequestType `db:"request_type"`
	}

	isPOSRequest := false
	query := `select request_type from business_money_request where id = $1`
	req := Request{}
	err := core.DBRead.Get(&req, query, *requestID)

	if err == nil && req.RequestType != nil && *req.RequestType == payment.PaymentRequestTypePOS {
		isPOSRequest = true
	}

	return isPOSRequest
}
