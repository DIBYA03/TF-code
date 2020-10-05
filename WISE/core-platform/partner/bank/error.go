package bank

import (
	"errors"
)

type ErrorCode string

type Error struct {
	RawError   error     `json:"rawError"`
	Code       ErrorCode `json:"code"`
	HTTPStatus int       `json:"httpStatus"`
}

func (e Error) Error() string {
	return e.RawError.Error()
}

const (
	ErrorCodeNotImplemented                     = ErrorCode("notImplemented")
	ErrorCodeInvalidBankProvider                = ErrorCode("invalidBankProvider")
	ErrorCodeLogRawNotification                 = ErrorCode("logRawNotification")
	ErrorCodeInvalidNotificationType            = ErrorCode("invalidNotificationCategory")
	ErrorCodeInvalidNotificationAction          = ErrorCode("invalidNotificationAction")
	ErrorCodeInvalidNotificationAttribute       = ErrorCode("invalidNotificationCategoryAttribute")
	ErrorCodeInvalidNotificationFormat          = ErrorCode("invalidNotificationFormat")
	ErrorCodeInternalDatabaseError              = ErrorCode("internalDatabaseError")
	ErrorCodeDuplicateNotification              = ErrorCode("duplicateNotification")
	ErrorCodeProcessTransactionNotification     = ErrorCode("errorProcessingTransactionNotification")
	ErrorCodeInvalidBankEntity                  = ErrorCode("invalidBankEntity")
	ErrorCodeInvalidNotificationTransactionType = ErrorCode("invalidNotificationTransactionType")
	ErrorCodeInvalidNotificationTransactionCode = ErrorCode("invalidNotificationTransactionCode")
	ErrorCodeHandleNotifications                = ErrorCode("errorHandlingNotifications")
	ErrorCodeHandleNotificationMessage          = ErrorCode("errorHandlingNotificationMessage")
	ErrorCodeInvalidCountryCode                 = ErrorCode("invalidCountryCode")
	ErrorCodeInvalidCitizenStatus               = ErrorCode("invalidCitizenStatus")
	ErrorCodeInvalidDocumentType                = ErrorCode("invalidDocumentType")
	ErrorCodeInvalidIssueDateFormat             = ErrorCode("invalidIssueDateFormat")
	ErrorCodeInvalidExpDateFormat               = ErrorCode("invalidExpDateFormat")
	ErrorCodeInvalidOccupation                  = ErrorCode("invalidOccupation")
	ErrorCodeInvalidIncome                      = ErrorCode("invalidIncome")
	ErrorCodeInvalidActivity                    = ErrorCode("invalidActivity")
	ErrorCodeInvalidCountry                     = ErrorCode("invalidCountry")
	ErrorCodeInvalidAddressType                 = ErrorCode("invalidAddressType")
	ErrorCodeInvalidKYCStatus                   = ErrorCode("invalidKYCStatus")
	ErrorCodeInvalidEntityType                  = ErrorCode("invalidEntityType")
	ErrorCodeInvalidIndustryType                = ErrorCode("invalidIndustryType")
	ErrorCodeInvalidMemberType                  = ErrorCode("invalidMemberType")
	ErrorCodeInvalidMemberTitle                 = ErrorCode("invalidMemberTitle")
	ErrorCodeInvalidBirthDate                   = ErrorCode("invalidBirthDate")
	ErrorCodeInvalidConsumerContactType         = ErrorCode("invalidConsumerContactType")
	ErrorCodeInvalidConsumerAddressType         = ErrorCode("invalidConsumerAddressType")
	ErrorCodeInvalidEntityID                    = ErrorCode("invalidEntityId")
	ErrorCodeInvalidAccountOpenDate             = ErrorCode("invalidAccountOpenDate")
	ErrorCodeInvalidAccountStatus               = ErrorCode("invalidAccountStatus")
	ErrorCodeInvalidMoneyTransferStatus         = ErrorCode("invalidMoneyTransferStatus")
	ErrorCodeInvalidMoneyTransferStatusDate     = ErrorCode("invalidMoneyTransferStatusDate")
	ErrorCodeInvalidMoneyTransferReasonCode     = ErrorCode("invalidMoneyTransferReasonCode")
	ErrorCodeInvalidCardStatus                  = ErrorCode("invalidCardStatus")
	ErrorCodeInvalidNotificationStatus          = ErrorCode("invalidNotificationStatus")
	ErrorCodeInvalidBusinessContactType         = ErrorCode("invalidBusinessContactType")
	ErrorCodeInvalidBusinessOp                  = ErrorCode("invalidBusinessOp")
	ErrorCodeInvalidTaxIDType                   = ErrorCode("invalidTaxIDType")
	ErrorCodeNoPayload                          = ErrorCode("noPayload")
)

func ErrorCodeRaw(e ErrorCode) string {
	return string(e)
}

var errorCodeMap = map[ErrorCode]string{
	ErrorCodeNotImplemented:                     "not implemented",
	ErrorCodeInvalidBankProvider:                "invalid bank provider",
	ErrorCodeLogRawNotification:                 "error logging raw notification",
	ErrorCodeInvalidNotificationType:            "invalid notification type",
	ErrorCodeInvalidNotificationAction:          "invalid notification action",
	ErrorCodeInvalidNotificationAttribute:       "invalid notification attribute",
	ErrorCodeInvalidNotificationFormat:          "invalid notification format",
	ErrorCodeProcessTransactionNotification:     "error processing transaction notification",
	ErrorCodeInvalidBankEntity:                  "invalid bank entity",
	ErrorCodeInvalidNotificationTransactionType: "invalid notification transaction type",
	ErrorCodeInvalidNotificationTransactionCode: "invalid notification transaction code",
	ErrorCodeHandleNotifications:                "error handling notifications",
	ErrorCodeHandleNotificationMessage:          "error handling notification message",
	ErrorCodeInvalidCountryCode:                 "invalid country code",
	ErrorCodeInvalidCitizenStatus:               "invalid citizenship status",
	ErrorCodeInvalidDocumentType:                "invalid document type",
	ErrorCodeInvalidIssueDateFormat:             "invalid issue date format",
	ErrorCodeInvalidExpDateFormat:               "invalid expiration date format",
	ErrorCodeInvalidOccupation:                  "invalid occupation",
	ErrorCodeInvalidIncome:                      "invalid income",
	ErrorCodeInvalidActivity:                    "invalid activity",
	ErrorCodeInvalidCountry:                     "invalid country",
	ErrorCodeInvalidAddressType:                 "invalid address type",
	ErrorCodeInvalidKYCStatus:                   "invalid KYC status",
	ErrorCodeInvalidEntityType:                  "invalid entity type",
	ErrorCodeInvalidIndustryType:                "invalid industry type",
	ErrorCodeInvalidMemberType:                  "invalid member type",
	ErrorCodeInvalidMemberTitle:                 "invalid member title",
	ErrorCodeInvalidBirthDate:                   "invalid birth date",
	ErrorCodeInvalidConsumerContactType:         "invalid consumer contact type",
	ErrorCodeInvalidConsumerAddressType:         "invalid consumer address type",
	ErrorCodeInvalidEntityID:                    "invalid entity id",
	ErrorCodeInvalidAccountOpenDate:             "invalid account open date",
	ErrorCodeInvalidAccountStatus:               "invalid account status",
	ErrorCodeInvalidMoneyTransferStatus:         "invalid money transfer status",
	ErrorCodeInvalidMoneyTransferStatusDate:     "invalid money transfer status date",
	ErrorCodeInvalidMoneyTransferReasonCode:     "invalid money transfer reason code",
	ErrorCodeInvalidCardStatus:                  "invalid card status",
	ErrorCodeInvalidNotificationStatus:          "invalid notification status",
	ErrorCodeInvalidBusinessContactType:         "invalid business contact type",
	ErrorCodeInvalidBusinessOp:                  "invalid business operation",
	ErrorCodeInvalidTaxIDType:                   "invalid tax id type",
	ErrorCodeNoPayload:                          "no payload",
}

func ErrorCodeString(e ErrorCode) string {
	return errorCodeMap[e]
}

func NewErrorFromCode(e ErrorCode) error {
	return &Error{
		RawError: errors.New(ErrorCodeString(e)),
		Code:     e,
	}
}
