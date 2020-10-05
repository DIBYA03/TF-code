package bank

import "time"

type MoneyTransferAccountType string

const (
	MoneyTransferAccountTypeBank = MoneyTransferAccountType("bank")
	MoneyTransferAccountTypeCard = MoneyTransferAccountType("card")
)

type MoneyTransferAccountBankID string

func (id MoneyTransferAccountBankID) String() string {
	return string(id)
}

type MoneyTransferRequest struct {
	SourceAccountID MoneyTransferAccountBankID `json:"sourceAccountId"`
	DestAccountID   MoneyTransferAccountBankID `json:"destAccountId"`
	Amount          float64                    `json:"amount"`
	Currency        Currency                   `json:"currency,omitempty"`
}

type MoneyTransferBankID string

func (id MoneyTransferBankID) String() string {
	return string(id)
}

type ReturnCode string

func (c ReturnCode) String() string {
	return string(c)
}

type ReturnCodeType string

const (
	// Empty return code
	ReturnCodeTypeEmpty = ReturnCodeType("")

	// ACH return code R01 - R85
	ReturnCodeTypeACH = ReturnCodeType("ach")
)

type ChangeCode string

func (c ChangeCode) String() string {
	return string(c)
}

type ChangeCodeType string

const (
	// Empty change code
	ChangeCodeTypeEmpty = ChangeCodeType("")

	// ACH change code C01 - C69
	ChangeCodeTypeACH = ChangeCodeType("ach")
)

type MoneyTransferStatus string

const (
	// Transfer is pending - can be cancelled
	MoneyTransferStatusPending = MoneyTransferStatus("pending")

	// Transfer submitted and is in process
	MoneyTransferStatusInProcess = MoneyTransferStatus("inProcess")

	// Transfer has been canceled
	MoneyTransferStatusCanceled = MoneyTransferStatus("canceled")

	// Transfer Posted (intrabank)
	MoneyTransferStatusPosted = MoneyTransferStatus("posted")

	// Transfer settled (extenbal or push to debit)
	MoneyTransferStatusSettled = MoneyTransferStatus("settled")

	// Debit sent to origin bank
	MoneyTransferStatusDebitSent = MoneyTransferStatus("debitSent")

	// Credit sent to destination bank
	MoneyTransferStatusCreditSent = MoneyTransferStatus("creditSent")

	// Review issue resolved
	MoneyTransferStatusReviewResolved = MoneyTransferStatus("reviewResolved")

	// Pull failed
	MoneyTransferStatusPullFailed = MoneyTransferStatus("pullFailed")

	// Pull refunded
	MoneyTransferStatusPullRefunded = MoneyTransferStatus("pullRefunded")

	// Pull transfer under review
	MoneyTransferStatusPullReview = MoneyTransferStatus("pullReview")

	// Push refunded
	MoneyTransferStatusPushRefunded = MoneyTransferStatus("pushRefunded")

	// Push under review
	MoneyTransferStatusPushReview = MoneyTransferStatus("pushReview")

	// Disbursed
	MoneyTransferStatusDisbursed = MoneyTransferStatus("disbursed")

	// Check disbursed
	MoneyTransferStatusCheckDisbursed = MoneyTransferStatus("checkDisbursed")

	// Check cleared
	MoneyTransferStatusCheckCleared = MoneyTransferStatus("checkCleared")
)

type MoneyTransferResponse struct {
	TransferID           MoneyTransferBankID        `json:"transferId"`
	SourceAccountID      MoneyTransferAccountBankID `json:"sourceAccountId"`
	DestAccountID        MoneyTransferAccountBankID `json:"destAccountId"`
	Amount               float64                    `json:"amount"`
	Currency             Currency                   `json:"currency,omitempty"`
	Created              time.Time                  `json:"created"`
	Status               MoneyTransferStatus        `json:"transferStatus"`
	LastStatusChange     time.Time                  `json:"lastStatusChange"`
	SourceReturnCodeType ReturnCodeType             `json:"sourceReturnCodeType"`
	SourceReturnCode     *ReturnCode                `json:"sourceReturnCode"`
	SourceReturnDesc     *string                    `json:"sourceReturnDesc"`
	DestReturnCodeType   ReturnCodeType             `json:"destReturnCodeType"`
	DestReturnCode       *ReturnCode                `json:"destReturnCode"`
	DestReturnDesc       *string                    `json:"destReturnDesc"`
	SourceChangeCodeType ChangeCodeType             `json:"sourceChangeCodeType "`
	SourceChangeCode     *ChangeCode                `json:"sourceChangeCode"`
	SourceChangeDesc     *string                    `json:"sourceChangeDesc"`
	SourceChangeData     *string                    `json:"sourceChangeData"`
	DestChangeCodeType   ChangeCodeType             `json:"destChangeCodeType "`
	DestChangeCode       *ChangeCode                `json:"destChangeCode"`
	DestChangeDesc       *string                    `json:"destChangeDes"`
	DestChangeData       *string                    `json:"destChangeData"`
}
