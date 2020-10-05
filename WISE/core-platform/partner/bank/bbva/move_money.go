package bbva

import (
	"errors"
	"time"

	partnerbank "github.com/wiseco/core-platform/partner/bank"
)

// Metadata fields must all be a string or typed string
type MoveMoneyMetadata struct {
	Currency Currency `json:"currency"`
}

// https://bbvaopenplatform.com/docs/guides%7Capicontent%7C08-move-money?code=674527&token=5c7df9a7e8288600018c9108

// MoveMoneyRequest
//
// origin_account - alphanumeric account registration identifier (RA-) of the funded source account.
// destination_account - alphanumeric account registration identifier (RA-) of the destination account.
//
// Each key-value pair can be no longer than 200 characters total.
type MoveMoneyRequest struct {
	OriginAccount      string            `json:"origin_account"`
	DestinationAccount string            `json:"destination_account"`
	Amount             float64           `json:"amount"`
	Metadata           MoveMoneyMetadata `json:"metadata"`
}

// MoveMoneyResponse
//
// About 404 Not Found errors
// When registering the destination bank account for a particular move money transaction, the request
// header OP-User-Id must specify the user ID of the funded source account. This enables the source
// account to operate on the destination account.
type MoveMoneyResponse struct {
	MoveMoneyID          string          `json:"move_money_id"`
	TransactionStatus    MoveMoneyStatus `json:"transaction_status"`
	DateLastStatusChange time.Time       `json:"date_last_status_change"`
}

type MoveMoneyStatus string

const (
	MoveMoneyStatusInProcess             = MoveMoneyStatus("in_process")
	MoveMoneyStatusCancelled             = MoveMoneyStatus("canceled")
	MoveMoneyStatusPosted                = MoveMoneyStatus("posted")
	MoveMoneyStatusDebitSent             = MoveMoneyStatus("debit_sent")
	MoveMoneyStatusCreditSent            = MoveMoneyStatus("credit_sent")
	MoveMoneyStatusDisbursed             = MoveMoneyStatus("disbursed")
	MoveMoneyStatusResolved              = MoveMoneyStatus("resolved")
	MoveMoneyStatusPullFailed            = MoveMoneyStatus("pull_failed")
	MoveMoneyStatusPullFailedRefunded    = MoveMoneyStatus("pull_failed_refunded")
	MoveMoneyStatusPullFailedUnderReview = MoveMoneyStatus("pull_failed_under_review")
	MoveMoneyStatusPushFailedRefunded    = MoveMoneyStatus("push_failed_refunded")
	MoveMoneyStatusPushFailedUnderReview = MoveMoneyStatus("push_failed_under_review")
	MoveMoneyStatusSettled               = MoveMoneyStatus("settled")
	MoveMoneyStatusCheckDisbursed        = MoveMoneyStatus("disbursed_check")
	MoveMoneyStatusCheckCleared          = MoveMoneyStatus("check_cleared")
)

var partnerMoneyTransferStatusFrom = map[partnerbank.MoneyTransferStatus]MoveMoneyStatus{
	partnerbank.MoneyTransferStatusInProcess:      MoveMoneyStatusInProcess,
	partnerbank.MoneyTransferStatusCanceled:       MoveMoneyStatusCancelled,
	partnerbank.MoneyTransferStatusPosted:         MoveMoneyStatusPosted,
	partnerbank.MoneyTransferStatusDebitSent:      MoveMoneyStatusDebitSent,
	partnerbank.MoneyTransferStatusCreditSent:     MoveMoneyStatusCreditSent,
	partnerbank.MoneyTransferStatusReviewResolved: MoveMoneyStatusResolved,
	partnerbank.MoneyTransferStatusPullFailed:     MoveMoneyStatusPullFailed,
	partnerbank.MoneyTransferStatusPullRefunded:   MoveMoneyStatusPullFailedRefunded,
	partnerbank.MoneyTransferStatusPullReview:     MoveMoneyStatusPullFailedUnderReview,
	partnerbank.MoneyTransferStatusPushRefunded:   MoveMoneyStatusPushFailedRefunded,
	partnerbank.MoneyTransferStatusPushReview:     MoveMoneyStatusPushFailedUnderReview,
	partnerbank.MoneyTransferStatusSettled:        MoveMoneyStatusSettled,
	partnerbank.MoneyTransferStatusDisbursed:      MoveMoneyStatusDisbursed,
	partnerbank.MoneyTransferStatusCheckDisbursed: MoveMoneyStatusCheckDisbursed,
	partnerbank.MoneyTransferStatusCheckCleared:   MoveMoneyStatusCheckCleared,
}

var partnerMoneyTransferStatusTo = map[MoveMoneyStatus]partnerbank.MoneyTransferStatus{
	MoveMoneyStatusInProcess:             partnerbank.MoneyTransferStatusInProcess,
	MoveMoneyStatusCancelled:             partnerbank.MoneyTransferStatusCanceled,
	MoveMoneyStatusPosted:                partnerbank.MoneyTransferStatusPosted,
	MoveMoneyStatusDebitSent:             partnerbank.MoneyTransferStatusDebitSent,
	MoveMoneyStatusCreditSent:            partnerbank.MoneyTransferStatusCreditSent,
	MoveMoneyStatusResolved:              partnerbank.MoneyTransferStatusReviewResolved,
	MoveMoneyStatusPullFailed:            partnerbank.MoneyTransferStatusPullFailed,
	MoveMoneyStatusPullFailedRefunded:    partnerbank.MoneyTransferStatusPullRefunded,
	MoveMoneyStatusPullFailedUnderReview: partnerbank.MoneyTransferStatusPullReview,
	MoveMoneyStatusPushFailedRefunded:    partnerbank.MoneyTransferStatusPushRefunded,
	MoveMoneyStatusPushFailedUnderReview: partnerbank.MoneyTransferStatusPushReview,
	MoveMoneyStatusSettled:               partnerbank.MoneyTransferStatusSettled,
	MoveMoneyStatusDisbursed:             partnerbank.MoneyTransferStatusDisbursed,
	MoveMoneyStatusCheckDisbursed:        partnerbank.MoneyTransferStatusCheckDisbursed,
	MoveMoneyStatusCheckCleared:          partnerbank.MoneyTransferStatusCheckCleared,
}

type GetMoveMoneyCorrected struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Data        string `json:"data"`
}

type GetMoveMoneyResponse struct {
	MoveMoneyID          string                 `json:"move_money_id"`
	OriginAccount        string                 `json:"origin_account"`
	DestinationAccount   string                 `json:"destination_account"`
	Amount               float64                `json:"amount"`
	TransactionStatus    MoveMoneyStatus        `json:"transaction_status"`
	Created              time.Time              `json:"date_created"`
	DateLastStatusChange time.Time              `json:"date_last_status_change"`
	OriginReturnedCode   *string                `json:"if_origin_returned_code"`
	OriginReturnedDesc   *string                `json:"if_origin_returned_description"`
	DestReturnedCode     *string                `json:"if_destination_returned_code"`
	DestReturnedDesc     *string                `json:"if_destination_returned_description"`
	OriginCorrected      *GetMoveMoneyCorrected `json:"if_origin_corrected"`
	DestCorrected        *GetMoveMoneyCorrected `json:"if_destination_corrected"`
	Metadata             MoveMoneyMetadata      `json:"metadata"`
}

func (m *GetMoveMoneyResponse) partnerMoneyTransferResponseTo() (*partnerbank.MoneyTransferResponse, error) {

	currency, ok := partnerCurrencyTo[m.Metadata.Currency]
	if !ok {
		return nil, errors.New("invalid currency")
	}

	status, ok := partnerMoneyTransferStatusTo[m.TransactionStatus]
	if !ok {
		return nil, errors.New("invalid transfer status")
	}

	mResp := &partnerbank.MoneyTransferResponse{
		TransferID:           partnerbank.MoneyTransferBankID(m.MoveMoneyID),
		SourceAccountID:      partnerbank.MoneyTransferAccountBankID(m.OriginAccount),
		DestAccountID:        partnerbank.MoneyTransferAccountBankID(m.DestinationAccount),
		Amount:               m.Amount,
		Currency:             currency,
		Created:              m.Created,
		Status:               status,
		LastStatusChange:     m.DateLastStatusChange,
		SourceReturnCodeType: partnerbank.ReturnCodeTypeEmpty,
		DestReturnCodeType:   partnerbank.ReturnCodeTypeEmpty,
		SourceChangeCodeType: partnerbank.ChangeCodeTypeEmpty,
		DestChangeCodeType:   partnerbank.ChangeCodeTypeEmpty,
	}

	if m.OriginReturnedCode != nil {
		mResp.SourceReturnCodeType = partnerbank.ReturnCodeTypeACH
		retCode := partnerbank.ReturnCode(*m.OriginReturnedCode)
		mResp.SourceReturnCode = &retCode
		mResp.SourceReturnDesc = m.OriginReturnedDesc
	}

	if m.DestReturnedCode != nil {
		mResp.DestReturnCodeType = partnerbank.ReturnCodeTypeACH
		retCode := partnerbank.ReturnCode(*m.DestReturnedCode)
		mResp.DestReturnCode = &retCode
		mResp.DestReturnDesc = m.DestReturnedDesc
	}

	if m.OriginCorrected != nil {
		mResp.SourceChangeCodeType = partnerbank.ChangeCodeTypeACH
		retCode := partnerbank.ChangeCode(m.OriginCorrected.Code)
		mResp.SourceChangeCode = &retCode
		mResp.SourceChangeDesc = &m.OriginCorrected.Description
		mResp.SourceChangeData = &m.OriginCorrected.Data
	}

	if m.DestCorrected != nil {
		mResp.DestChangeCodeType = partnerbank.ChangeCodeTypeACH
		retCode := partnerbank.ChangeCode(m.DestCorrected.Code)
		mResp.DestChangeCode = &retCode
		mResp.DestChangeDesc = &m.DestCorrected.Description
		mResp.DestChangeData = &m.DestCorrected.Data
	}

	return mResp, nil
}
