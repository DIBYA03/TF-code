package bbva

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/partner/bank/data"
)

var partnerTransactionCodeTo = map[NotificationReason]bank.TransactionCode{
	NotificationReasonAuthorizationApproved: bank.TransactionCodeAuthApproved,
	NotificationReasonAuthorizationDeclined: bank.TransactionCodeAuthDeclined,
	NotificationReasonAuthorizationReversal: bank.TransactionCodeAuthReversed,
	NotificationReasonFundsHoldSet:          bank.TransactionCodeHoldApproved,
	NotificationReasonFundsHoldReleased:     bank.TransactionCodeHoldReleased,
	NotificationReasonCardDebitPosted:       bank.TransactionCodeDebitPosted,
	NotificationReasonCardCreditPosted:      bank.TransactionCodeCreditPosted,
	NotificationReasonNonCardDebitPosted:    bank.TransactionCodeDebitPosted,
	NotificationReasonNonCardCreditPosted:   bank.TransactionCodeCreditPosted,
}

var partnerTransactionNotificationActionTo = map[NotificationReason]bank.NotificationAction{
	NotificationReasonAuthorizationApproved: bank.NotificationActionAuthorize,
	NotificationReasonAuthorizationDeclined: bank.NotificationActionAuthorize,
	NotificationReasonAuthorizationReversal: bank.NotificationActionAuthorize,
	NotificationReasonFundsHoldSet:          bank.NotificationActionHold,
	NotificationReasonFundsHoldReleased:     bank.NotificationActionHold,
	NotificationReasonCardDebitPosted:       bank.NotificationActionPosted,
	NotificationReasonCardCreditPosted:      bank.NotificationActionPosted,
	NotificationReasonNonCardDebitPosted:    bank.NotificationActionPosted,
	NotificationReasonNonCardCreditPosted:   bank.NotificationActionPosted,
}

type NotificationTransactionType string

const (
	NotificationTransactionTypeACH         = NotificationTransactionType("ACH")
	NotificationTransactionTypeAdjustment  = NotificationTransactionType("ADJUSTMENT")
	NotificationTransactionTypeATM         = NotificationTransactionType("ATM")
	NotificationTransactionTypeCheck       = NotificationTransactionType("CHECK")
	NotificationTransactionTypeDeposit     = NotificationTransactionType("DEPOSIT")
	NotificationTransactionTypeFee         = NotificationTransactionType("FEE")
	NotificationTransactionTypeOtherCredit = NotificationTransactionType("OTHER_CREDIT")
	NotificationTransactionTypeOtherDebit  = NotificationTransactionType("OTHER_DEBIT")
	NotificationTransactionTypePurchase    = NotificationTransactionType("PURCHASE")
	NotificationTransactionTypeRefund      = NotificationTransactionType("REFUND")
	NotificationTransactionTypeReturn      = NotificationTransactionType("RETURN")
	NotificationTransactionTypeReversal    = NotificationTransactionType("REVERSAL")
	NotificationTransactionTypeTransfer    = NotificationTransactionType("TRANSFER")
	NotificationTransactionTypeVisaCredit  = NotificationTransactionType("VISA_CREDIT")
	NotificationTransactionTypeWithdrawal  = NotificationTransactionType("WITHDRAWAL")
	NotificationTransactionTypeOther       = NotificationTransactionType("OTHER")
)

var partnerTransactionTypeTo = map[NotificationTransactionType]bank.TransactionType{
	NotificationTransactionTypeACH:         bank.TransactionTypeACH,
	NotificationTransactionTypeAdjustment:  bank.TransactionTypeAdjustment,
	NotificationTransactionTypeATM:         bank.TransactionTypeATM,
	NotificationTransactionTypeCheck:       bank.TransactionTypeCheck,
	NotificationTransactionTypeDeposit:     bank.TransactionTypeDeposit,
	NotificationTransactionTypeFee:         bank.TransactionTypeFee,
	NotificationTransactionTypeOtherCredit: bank.TransactionTypeOtherCredit,
	NotificationTransactionTypeOtherDebit:  bank.TransactionTypeOtherDebit,
	NotificationTransactionTypePurchase:    bank.TransactionTypePurchase,
	NotificationTransactionTypeRefund:      bank.TransactionTypeRefund,
	NotificationTransactionTypeReversal:    bank.TransactionTypeReversal,
	NotificationTransactionTypeTransfer:    bank.TransactionTypeTransfer,
	NotificationTransactionTypeVisaCredit:  bank.TransactionTypeVisaCredit,
	NotificationTransactionTypeWithdrawal:  bank.TransactionTypeWithdrawal,
	NotificationTransactionTypeOther:       bank.TransactionTypeOther,
}

type TransactionNotificationData struct {
	ID              string                       `json:"transaction_id"`
	UserID          string                       `json:"user_id"`
	AccountID       *string                      `json:"account_id"`
	CardID          *string                      `json:"card_id"`
	MoveMoneyID     *string                      `json:"move_money_id"`
	Description     string                       `json:"transaction_description"`
	TransactionType NotificationTransactionType  `json:"transaction_type"`
	PostedAmount    float64                      `json:"posted_amount"`
	PostedBalance   *float64                     `json:"posted_balance"`
	PostedDate      time.Time                    `json:"posted_date"`
	Hold            *TransactionNotificationHold `json:"hold_data"`
	Card            *TransactionNotificationCard `json:"card_details"`
}

type TransactionNotificationHold struct {
	Number     int       `json:"number"`
	Amount     float64   `json:"amount"`
	Date       time.Time `json:"date"`
	ExpiryDate time.Time `json:"expiry_date"`
}

type TransactionNotificationCard struct {
	VisaTransactionID          string        `json:"visa_transaction_id"`
	AuthorizationAmount        float64       `json:"authorization_amount"`
	AuthorizationDate          time.Time     `json:"authorization_date"`
	AuthorizationResponse      string        `json:"authorization_response"`
	AuthorizationNumber        string        `json:"authorization_number"`
	CardTransactionType        string        `json:"card_transaction_type"`
	LocalTransactionAmount     float64       `json:"local_transaction_amount"`
	LocalTransactionCurrency   string        `json:"local_transaction_currency"`
	LocalTransactionDateTime   DateTimeLocal `json:"local_transaction_date_time"`
	BillingTransactionCurrency string        `json:"billing_transaction_currency"`
	POSEntryMode               string        `json:"pos_entry_mode"`
	POSConditionCode           string        `json:"pos_condition_code"`
	MerchantCategoryCode       string        `json:"merchant_category_code"`
	MerchantName               string        `json:"merchant_name"`
	AcquirerBIN                string        `json:"acquirer_bin"`
	CardAcceptorID             string        `json:"card_acceptor_id"`
	CardAcceptorTerminal       string        `json:"card_acceptor_terminal"`
	CardAcceptorAddress        string        `json:"card_acceptor_address"`
	CardAcceptorCity           string        `json:"card_acceptor_city"`
	CardAcceptorState          string        `json:"card_acceptor_state"`
	CardAcceptorCountry        string        `json:"card_acceptor_country"`
}

func processTransactionHold(d TransactionNotificationData) (*bank.HoldTransaction, error) {
	if d.Hold == nil {
		return nil, nil
	}

	return &bank.HoldTransaction{
		Number:     strconv.FormatInt(int64(d.Hold.Number), 10),
		Amount:     d.Hold.Amount,
		Date:       d.Hold.Date,
		ExpiryDate: d.Hold.ExpiryDate,
	}, nil
}

func processTransactionCard(d TransactionNotificationData) (*bank.CardTransaction, error) {
	if d.Card == nil {
		return nil, nil
	}

	return &bank.CardTransaction{
		CardTransactionID:     d.Card.VisaTransactionID,
		TransactionNetwork:    bank.CardTransactionNetworkVisa,
		AuthAmount:            d.Card.AuthorizationAmount,
		AuthDate:              d.Card.AuthorizationDate,
		AuthResponseCode:      d.Card.AuthorizationResponse,
		AuthNumber:            d.Card.AuthorizationNumber,
		TransactionType:       d.Card.CardTransactionType,
		LocalAmount:           d.Card.LocalTransactionAmount,
		LocalCurrency:         bank.Currency(strings.ToLower(d.Card.LocalTransactionCurrency)),
		BillingCurrency:       bank.Currency(strings.ToLower(d.Card.BillingTransactionCurrency)),
		LocalDate:             d.Card.LocalTransactionDateTime.Time(),
		POSEntryMode:          d.Card.POSEntryMode,
		POSConditionCode:      d.Card.POSConditionCode,
		AcquirerBIN:           d.Card.AcquirerBIN,
		MerchantID:            d.Card.CardAcceptorID,
		MerchantCategoryCode:  d.Card.MerchantCategoryCode,
		MerchantTerminal:      d.Card.CardAcceptorTerminal,
		MerchantName:          d.Card.MerchantName,
		MerchantStreetAddress: d.Card.CardAcceptorAddress,
		MerchantCity:          d.Card.CardAcceptorCity,
		MerchantState:         d.Card.CardAcceptorState,
		MerchantCountry:       d.Card.CardAcceptorCountry,
	}, nil
}

func processTransactionData(n Notification, d TransactionNotificationData) (*bank.Transaction, error) {
	// Process any hold details
	hold, err := processTransactionHold(d)
	if err != nil {
		return nil, err
	}

	// Process any card details
	card, err := processTransactionCard(d)
	if err != nil {
		return nil, err
	}

	code, ok := partnerTransactionCodeTo[n.Reason]
	if !ok {
		return nil, bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationTransactionCode)
	}

	transactionType, ok := partnerTransactionTypeTo[d.TransactionType]
	if !ok {
		switch code {
		case
			bank.TransactionCodeDebitPosted, bank.TransactionCodeCreditPosted,
			bank.TransactionCodeAuthApproved, bank.TransactionCodeAuthDeclined,
			bank.TransactionCodeAuthReversed, bank.TransactionCodeHoldApproved,
			bank.TransactionCodeHoldReleased:
			transactionType = bank.TransactionTypeOther
		default:
			return nil, bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationTransactionType)
		}
	}

	var amount float64
	if d.PostedAmount != 0 {
		amount = d.PostedAmount
	} else {
		if card != nil && card.AuthAmount != 0 {
			amount = card.AuthAmount
		} else if hold != nil {
			amount = hold.Amount
		}
	}

	var currency = bank.CurrencyUSD

	var desc string = d.Description

	// Use posted time first
	var transactionDate time.Time = d.PostedDate

	// Use hold date (utc) if zero
	if transactionDate.IsZero() && hold != nil {
		transactionDate = hold.Date
	}

	// Use card auth date (UTC) if zero
	if transactionDate.IsZero() && card != nil {
		transactionDate = card.AuthDate
	}

	// Default to sent date
	if transactionDate.IsZero() {
		transactionDate = n.SentDate
	}

	// Last resort use current time (UTC)
	if transactionDate.IsZero() {
		transactionDate = time.Now().UTC()
	}

	// Grock the MM id if omitted
	if d.MoveMoneyID == nil {
		// Match MM prefix with UUID
		re, err := regexp.Compile(`MM-[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}`)
		if err == nil {
			match := re.Find([]byte(desc))
			if match != nil {
				moveMoneyID := string(match)
				d.MoveMoneyID = &moveMoneyID
			}
		}
	}

	return &bank.Transaction{
		BankName:            bank.ProviderNameBBVA,
		BankTransactionID:   d.ID,
		TransactionType:     transactionType,
		AccountID:           d.AccountID,
		CardID:              d.CardID,
		CodeType:            code,
		Amount:              amount,
		PostedBalance:       d.PostedBalance,
		Currency:            currency,
		CardTransaction:     card,
		HoldTransaction:     hold,
		BankMoneyTransferID: d.MoveMoneyID,
		BankTransactionDesc: &desc,
		TransactionDate:     transactionDate,
	}, nil
}

func (s *notificationService) processTransactionNotificationMessage(n Notification) error {
	var d TransactionNotificationData
	err := n.unmarshalData(&d)
	if err != nil {
		return &bank.Error{
			RawError: err,
			Code:     bank.ErrorCodeProcessTransactionNotification,
		}
	}

	e, err := notificationEntityFromCustomerID(d.UserID)
	if err != nil {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidEntityID)
	}

	nt, ok := partnerNotificationTypeTo[n.Type]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationType)
	}

	action, ok := partnerTransactionNotificationActionTo[n.Reason]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationAction)
	}

	tx, err := processTransactionData(n, d)
	if err != nil {
		return err
	}

	b, err := json.Marshal(tx)
	if err != nil {
		return &bank.Error{
			RawError: err,
			Code:     bank.ErrorCodeProcessTransactionNotification,
		}
	}

	pn := data.NotificationCreate{
		EntityID:   bank.NotificationEntityID(e.EntityID),
		EntityType: bank.NotificationEntityType(e.EntityType),
		BankName:   bank.ProviderNameBBVA,
		SourceID:   bank.SourceID(n.notificationID()),
		Type:       nt,
		Action:     action,
		Version:    n.Version,
		Created:    n.SentDate,
		Data:       b,
	}

	return s.forwardNotification(pn)
}
