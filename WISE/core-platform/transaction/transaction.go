/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package transaction

import (
	"strings"

	"github.com/wiseco/core-platform/notification/activity"
	grpcBanking "github.com/wiseco/protobuf/golang/banking"
	grpcTxn "github.com/wiseco/protobuf/golang/transaction"
)

type TransactionStatus string

const (
	TransactionStatusUnspecified = TransactionStatus("")

	// Hold
	TransactionStatusHoldSet      = TransactionStatus("holdSet")
	TransactionStatusHoldReleased = TransactionStatus("holdReleased")

	// Cards
	TransactionStatusCardAuthorized   = TransactionStatus("cardAuthorized")
	TransactionStatusCardAuthDeclined = TransactionStatus("cardAuthDeclined")
	TransactionStatusCardAuthReversed = TransactionStatus("cardAuthReversed")
	TransactionStatusCardPosted       = TransactionStatus("cardPosted")

	// Accounts
	TransactionStatusValidation       = TransactionStatus("validation")
	TransactionStatusReview           = TransactionStatus("review")
	TransactionStatusSystemDeclined   = TransactionStatus("systemDeclined")
	TransactionStatusAgentDeclined    = TransactionStatus("agentDeclined")
	TransactionStatusBankDeclined     = TransactionStatus("bankDeclined")
	TransactionStatusAgentCanceled    = TransactionStatus("agentCanceled")
	TransactionStatusCustomerCanceled = TransactionStatus("customerCanceled")
	TransactionStatusBankCanceled     = TransactionStatus("bankCanceled")
	TransactionStatusBankProcessing   = TransactionStatus("bankProcessing")
	TransactionStatusNonCardPosted    = TransactionStatus("nonCardPosted")
	TransactionStatusTransferError    = TransactionStatus("transferError")
)

var TransactionStatusToProto = map[TransactionStatus]grpcTxn.BankTransactionStatus{
	TransactionStatusUnspecified:      grpcTxn.BankTransactionStatus_BTS_UNSPECIFIED,
	TransactionStatusHoldSet:          grpcTxn.BankTransactionStatus_BTS_HOLD_SET,
	TransactionStatusHoldReleased:     grpcTxn.BankTransactionStatus_BTS_HOLD_RELEASED,
	TransactionStatusCardAuthorized:   grpcTxn.BankTransactionStatus_BTS_CARD_AUTHORIZED,
	TransactionStatusCardAuthDeclined: grpcTxn.BankTransactionStatus_BTS_CARD_AUTH_DECLINED,
	TransactionStatusCardAuthReversed: grpcTxn.BankTransactionStatus_BTS_CARD_AUTH_REVERSED,
	TransactionStatusCardPosted:       grpcTxn.BankTransactionStatus_BTS_CARD_POSTED,
	TransactionStatusValidation:       grpcTxn.BankTransactionStatus_BTS_VALIDATION,
	TransactionStatusReview:           grpcTxn.BankTransactionStatus_BTS_REVIEW,
	TransactionStatusSystemDeclined:   grpcTxn.BankTransactionStatus_BTS_SYSTEM_DECLINED,
	TransactionStatusAgentDeclined:    grpcTxn.BankTransactionStatus_BTS_AGENT_DECLINED,
	TransactionStatusBankDeclined:     grpcTxn.BankTransactionStatus_BTS_BANK_DECLINED,
	TransactionStatusAgentCanceled:    grpcTxn.BankTransactionStatus_BTS_AGENT_CANCELED,
	TransactionStatusCustomerCanceled: grpcTxn.BankTransactionStatus_BTS_CUSTOMER_CANCELED,
	TransactionStatusBankCanceled:     grpcTxn.BankTransactionStatus_BTS_BANK_CANCELED,
	TransactionStatusBankProcessing:   grpcTxn.BankTransactionStatus_BTS_BANK_PROCESSING,
	TransactionStatusNonCardPosted:    grpcTxn.BankTransactionStatus_BTS_NONCARD_POSTED,
	TransactionStatusTransferError:    grpcTxn.BankTransactionStatus_BTS_BANK_TRANSFER_ERROR,
}

var TransactionStatusFromProto = map[grpcTxn.BankTransactionStatus]TransactionStatus{
	grpcTxn.BankTransactionStatus_BTS_UNSPECIFIED:         TransactionStatusUnspecified,
	grpcTxn.BankTransactionStatus_BTS_HOLD_SET:            TransactionStatusHoldSet,
	grpcTxn.BankTransactionStatus_BTS_HOLD_RELEASED:       TransactionStatusHoldReleased,
	grpcTxn.BankTransactionStatus_BTS_CARD_AUTHORIZED:     TransactionStatusCardAuthorized,
	grpcTxn.BankTransactionStatus_BTS_CARD_AUTH_DECLINED:  TransactionStatusCardAuthDeclined,
	grpcTxn.BankTransactionStatus_BTS_CARD_AUTH_REVERSED:  TransactionStatusCardAuthReversed,
	grpcTxn.BankTransactionStatus_BTS_CARD_POSTED:         TransactionStatusCardPosted,
	grpcTxn.BankTransactionStatus_BTS_VALIDATION:          TransactionStatusValidation,
	grpcTxn.BankTransactionStatus_BTS_REVIEW:              TransactionStatusReview,
	grpcTxn.BankTransactionStatus_BTS_SYSTEM_DECLINED:     TransactionStatusSystemDeclined,
	grpcTxn.BankTransactionStatus_BTS_AGENT_DECLINED:      TransactionStatusAgentDeclined,
	grpcTxn.BankTransactionStatus_BTS_BANK_DECLINED:       TransactionStatusBankDeclined,
	grpcTxn.BankTransactionStatus_BTS_AGENT_CANCELED:      TransactionStatusAgentCanceled,
	grpcTxn.BankTransactionStatus_BTS_CUSTOMER_CANCELED:   TransactionStatusCustomerCanceled,
	grpcTxn.BankTransactionStatus_BTS_BANK_CANCELED:       TransactionStatusBankCanceled,
	grpcTxn.BankTransactionStatus_BTS_BANK_PROCESSING:     TransactionStatusBankProcessing,
	grpcTxn.BankTransactionStatus_BTS_NONCARD_POSTED:      TransactionStatusNonCardPosted,
	grpcTxn.BankTransactionStatus_BTS_BANK_TRANSFER_ERROR: TransactionStatusTransferError,
}

const (
	// Bank account
	TransactionSourceTypeAccount = "account"

	// Debit or credit card
	TransactionSourceTypeCard = "card"
)

type TransactionType string

const (
	TransactionTypeUnspecified = TransactionType("")

	// ACH
	TransactionTypeACH = TransactionType("ach")

	// Adjustment
	TransactionTypeAdjustment = TransactionType("adjustment")

	// ATM
	TransactionTypeATM = TransactionType("atm")

	// Check
	TransactionTypeCheck = TransactionType("check")

	// Deposit
	TransactionTypeDeposit = TransactionType("deposit")

	// Fee
	TransactionTypeFee = TransactionType("fee")

	// Other Credit
	TransactionTypeOtherCredit = TransactionType("otherCredit")

	// Other Debit
	TransactionTypeOtherDebit = TransactionType("otherDebit")

	// Purchase transaction
	TransactionTypePurchase = TransactionType("purchase")

	// Refund
	TransactionTypeRefund = TransactionType("refund")

	// Return
	TransactionTypeReturn = TransactionType("return")

	// Reversal
	TransactionTypeReversal = TransactionType("reversal")

	// Transfer
	TransactionTypeTransfer = TransactionType("transfer")

	// Visa credit
	TransactionTypeVisaCredit = TransactionType("visaCredit")

	// Withdrawal
	TransactionTypeWithdrawal = TransactionType("withdrawal")

	// Other
	TransactionTypeOther = TransactionType("other")
)

var TransactionTypeToCategoryProto = map[TransactionType]grpcTxn.BankTransactionCategory{
	TransactionTypeUnspecified: grpcTxn.BankTransactionCategory_BTC_UNSPECIFIED,
	TransactionTypeACH:         grpcTxn.BankTransactionCategory_BTC_ACH,
	TransactionTypeAdjustment:  grpcTxn.BankTransactionCategory_BTC_ADJUSTMENT,
	TransactionTypeATM:         grpcTxn.BankTransactionCategory_BTC_ATM,
	TransactionTypeCheck:       grpcTxn.BankTransactionCategory_BTC_CHECK,
	TransactionTypeDeposit:     grpcTxn.BankTransactionCategory_BTC_DEPOSIT,
	TransactionTypeFee:         grpcTxn.BankTransactionCategory_BTC_FEE,
	TransactionTypeOtherCredit: grpcTxn.BankTransactionCategory_BTC_OTHER_CREDIT,
	TransactionTypeOtherDebit:  grpcTxn.BankTransactionCategory_BTC_OTHER_DEBIT,
	TransactionTypePurchase:    grpcTxn.BankTransactionCategory_BTC_PURCHASE,
	TransactionTypeRefund:      grpcTxn.BankTransactionCategory_BTC_REFUND,
	TransactionTypeReturn:      grpcTxn.BankTransactionCategory_BTC_RETURN,
	TransactionTypeReversal:    grpcTxn.BankTransactionCategory_BTC_REVERSAL,
	TransactionTypeTransfer:    grpcTxn.BankTransactionCategory_BTC_INTRABANK,
	TransactionTypeVisaCredit:  grpcTxn.BankTransactionCategory_BTC_VISA_CREDIT,
	TransactionTypeWithdrawal:  grpcTxn.BankTransactionCategory_BTC_WITHDRAWAL,
	TransactionTypeOther:       grpcTxn.BankTransactionCategory_BTC_OTHER,
}

var TransactionTypeFromCategoryProto = map[grpcTxn.BankTransactionCategory]TransactionType{
	grpcTxn.BankTransactionCategory_BTC_UNSPECIFIED:  TransactionTypeUnspecified,
	grpcTxn.BankTransactionCategory_BTC_ACH:          TransactionTypeACH,
	grpcTxn.BankTransactionCategory_BTC_ADJUSTMENT:   TransactionTypeAdjustment,
	grpcTxn.BankTransactionCategory_BTC_ATM:          TransactionTypeATM,
	grpcTxn.BankTransactionCategory_BTC_CHECK:        TransactionTypeCheck,
	grpcTxn.BankTransactionCategory_BTC_DEPOSIT:      TransactionTypeDeposit,
	grpcTxn.BankTransactionCategory_BTC_FEE:          TransactionTypeFee,
	grpcTxn.BankTransactionCategory_BTC_OTHER_CREDIT: TransactionTypeOtherCredit,
	grpcTxn.BankTransactionCategory_BTC_OTHER_DEBIT:  TransactionTypeOtherDebit,
	grpcTxn.BankTransactionCategory_BTC_PURCHASE:     TransactionTypePurchase,
	grpcTxn.BankTransactionCategory_BTC_REFUND:       TransactionTypeRefund,
	grpcTxn.BankTransactionCategory_BTC_RETURN:       TransactionTypeReturn,
	grpcTxn.BankTransactionCategory_BTC_REVERSAL:     TransactionTypeReversal,
	grpcTxn.BankTransactionCategory_BTC_INTRABANK:    TransactionTypeTransfer,
	grpcTxn.BankTransactionCategory_BTC_VISA_CREDIT:  TransactionTypeVisaCredit,
	grpcTxn.BankTransactionCategory_BTC_WITHDRAWAL:   TransactionTypeWithdrawal,
	grpcTxn.BankTransactionCategory_BTC_OTHER:        TransactionTypeOther,
}

type TransactionCodeType string

const (
	TransactionCodeTypeStatusChange   = TransactionCodeType("statusChange")
	TransactionCodeTypeAuthApproved   = TransactionCodeType("authApproved")
	TransactionCodeTypeAuthDeclined   = TransactionCodeType("authDeclined")
	TransactionCodeTypeAuthReversed   = TransactionCodeType("authReversed")
	TransactionCodeTypeHoldApproved   = TransactionCodeType("holdApproved")
	TransactionCodeTypeHoldReleased   = TransactionCodeType("holdReleased")
	TransactionCodeTypeDebitPosted    = TransactionCodeType("debitPosted")
	TransactionCodeTypeCreditPosted   = TransactionCodeType("creditPosted")
	TransactionCodeTypeTransferChange = TransactionCodeType("transferChange")

	//Custom types
	TransactionCodeTypeCreditInProcess = TransactionCodeType("creditInProcess")
	TransactionCodeTypeDebitInProcess  = TransactionCodeType("debitInProcess")

	// New Mappings
	TransactionCodeTypeValidation = TransactionCodeType("validation")
	TransactionCodeTypeReview     = TransactionCodeType("review")

	// Aggregate code types
	TransactionCodeTypeTransferError = TransactionCodeType("transferError")
	TransactionCodeTypeDeclined      = TransactionCodeType("declined")
	TransactionCodeTypeCanceled      = TransactionCodeType("canceled")
)

const (
	// Wise Network (Intrabank)
	TransactionNetworkTypeWise = "wise"

	// ACH network
	TransactionNetworkTypeACH = "ach"

	// Visa network
	TransactionNetworkTypeVisa = "visa"

	// Mastercard network
	TransactionNetworkTypeMC = "mc"

	// Plus network
	TransactionNetworkTypePlus = "plus"
)

var CardNetworkToProto = map[string]grpcBanking.DebitCardNetwork{
	"VISA": grpcBanking.DebitCardNetwork_DCN_VISA,
}

type CardTransactionType string

const (
	CardTransactionTypeMoneyTransfer = CardTransactionType("26")
)

func (t CardTransactionType) IsRefundTypeInstantPay() bool {
	return strings.HasPrefix(string(t), string(CardTransactionTypeMoneyTransfer))
}

type TransactionSubtype string

const (
	TransactionSubtypeUnspecified              = TransactionSubtype("")
	TransactionSubtypeAccountOriginationCredit = TransactionSubtype("accountOriginationCredit")
	TransactionSubtypeCardReaderCredit         = TransactionSubtype("cardReaderCredit")
	TransactionSubtypeCardOnlineCredit         = TransactionSubtype("cardOnlineCredit")
	TransactionSubtypeBankOnlineCredit         = TransactionSubtype("bankOnlineCredit")
	TransactionSubtypeInterestTransferCredit   = TransactionSubtype("interestTransferCredit")
	TransactionSubtypeWiseTransferCredit       = TransactionSubtype("wiseTransferCredit")
	TransactionSubtypeACHTransferCredit        = TransactionSubtype("achTransferCredit")
	TransactionSubtypeACHTransferShopifyCredit = TransactionSubtype("achTransferShopifyCredit")
	TransactionSubtypeWireTransferCredit       = TransactionSubtype("wireTransferCredit")
	TransactionSubtypeCheckCredit              = TransactionSubtype("checkCredit")
	TransactionSubtypeDepositCredit            = TransactionSubtype("depositCredit")
	TransactionSubtypeMerchantRefundCredit     = TransactionSubtype("merchantRefundCredit")
	TransactionSubtypeCardPullCredit           = TransactionSubtype("cardPullCredit")
	TransactionSubtypeOtherCredit              = TransactionSubtype("otherCredit")
	TransactionSubtypeCardPullDebit            = TransactionSubtype("cardPullDebit")
	TransactionSubtypeCardPurchaseDebit        = TransactionSubtype("cardPurchaseDebit")
	TransactionSubtypeCardPurchaseDebitOnline  = TransactionSubtype("cardPurchaseDebitOnline")
	TransactionSubtypeWiseTransferDebit        = TransactionSubtype("wiseTransferDebit")
	TransactionSubtypeACHTransferDebit         = TransactionSubtype("achTransferDebit")
	TransactionSubtypeACHTransferShopifyDebit  = TransactionSubtype("achTransferShopifyDebit")
	TransactionSubtypeWireTransferDebit        = TransactionSubtype("wireTransferDebit")
	TransactionSubtypeCardATMDebit             = TransactionSubtype("cardATMDebit")
	TransactionSubtypeCardPushDebit            = TransactionSubtype("cardPushDebit")
	TransactionSubtypeCardPushCredit           = TransactionSubtype("cardPushCredit")
	TransactionSubtypeCheckDebit               = TransactionSubtype("checkDebit")
	TransactionSubtypeFeeDebit                 = TransactionSubtype("feeDebit")
	TransactionSubtypeHoldApproved             = TransactionSubtype("holdApproved")
	TransactionSubtypeHoldReleased             = TransactionSubtype("holdReleased")
	TransactionSubtypeOtherDebit               = TransactionSubtype("otherDebit")
)

var ActivityToTransactionSubtype = map[activity.Type]TransactionSubtype{
	activity.TypeAccountOrigination:            TransactionSubtypeAccountOriginationCredit,
	activity.TypeCardReaderCredit:              TransactionSubtypeCardReaderCredit,
	activity.TypeCardOnlineCredit:              TransactionSubtypeCardOnlineCredit,
	activity.TypeBankOnlineCredit:              TransactionSubtypeBankOnlineCredit,
	activity.TypeInterestTransferCredit:        TransactionSubtypeInterestTransferCredit,
	activity.TypeWiseTransferCredit:            TransactionSubtypeWiseTransferCredit,
	activity.TypeACHTransferCredit:             TransactionSubtypeACHTransferCredit,
	activity.TypeWireTransferCredit:            TransactionSubtypeWireTransferCredit,
	activity.TypeCheckCredit:                   TransactionSubtypeCheckCredit,
	activity.TypeDepositCredit:                 TransactionSubtypeDepositCredit,
	activity.TypeMerchantRefundCredit:          TransactionSubtypeMerchantRefundCredit,
	activity.TypeCardPullCredit:                TransactionSubtypeCardPullCredit,
	activity.TypeCardPullDebit:                 TransactionSubtypeCardPullDebit,
	activity.TypeCardReaderPurchaseDebit:       TransactionSubtypeCardPurchaseDebit,
	activity.TypeCardReaderPurchaseDebitOnline: TransactionSubtypeCardPurchaseDebitOnline,
	activity.TypeWiseTransferDebit:             TransactionSubtypeWiseTransferDebit,
	activity.TypeACHTransferDebit:              TransactionSubtypeACHTransferDebit,
	activity.TypeCardATMDebit:                  TransactionSubtypeCardATMDebit,
	activity.TypeCardPushDebit:                 TransactionSubtypeCardPushDebit,
	activity.TypeCardPushCredit:                TransactionSubtypeCardPushCredit,
	activity.TypeCheckDebit:                    TransactionSubtypeCheckDebit,
	activity.TypeACHTransferShopifyCredit:      TransactionSubtypeACHTransferShopifyCredit,
	activity.TypeHoldApproved:                  TransactionSubtypeHoldApproved,
	activity.TypeHoldReleased:                  TransactionSubtypeHoldReleased,
	activity.TypeFeeDebit:                      TransactionSubtypeFeeDebit,
	activity.TypeOtherCredit:                   TransactionSubtypeOtherCredit,
	activity.TypeCardVisaCredit:                TransactionSubtypeCardPullCredit,
}

var TransactionSubtypeToTypeProto = map[TransactionSubtype]grpcTxn.BankTransactionType{
	TransactionSubtypeUnspecified: grpcTxn.BankTransactionType_BTT_UNSPECIFIED,

	// Origination
	TransactionSubtypeAccountOriginationCredit: grpcTxn.BankTransactionType_BTT_ORIGINATION,

	// Credits
	TransactionSubtypeInterestTransferCredit:   grpcTxn.BankTransactionType_BTT_INTEREST_CREDIT,
	TransactionSubtypeWiseTransferCredit:       grpcTxn.BankTransactionType_BTT_INTRABANK_CREDIT,
	TransactionSubtypeACHTransferCredit:        grpcTxn.BankTransactionType_BTT_ACH_CREDIT,
	TransactionSubtypeWireTransferCredit:       grpcTxn.BankTransactionType_BTT_WIRE_CREDIT,
	TransactionSubtypeCardPullCredit:           grpcTxn.BankTransactionType_BTT_CARD_PULL_CREDIT,
	TransactionSubtypeCardPushCredit:           grpcTxn.BankTransactionType_BTT_CARD_PUSH_CREDIT,
	TransactionSubtypeCardReaderCredit:         grpcTxn.BankTransactionType_BTT_CARD_PAYMENT_READER_CREDIT,
	TransactionSubtypeCardOnlineCredit:         grpcTxn.BankTransactionType_BTT_CARD_PAYMENT_ONLINE_CREDIT,
	TransactionSubtypeBankOnlineCredit:         grpcTxn.BankTransactionType_BTT_INTRABANK_ONLINE_CREDIT,
	TransactionSubtypeCheckCredit:              grpcTxn.BankTransactionType_BTT_CHECK_CREDIT,
	TransactionSubtypeMerchantRefundCredit:     grpcTxn.BankTransactionType_BTT_CARD_PURCHASE_REFUND_CREDIT,
	TransactionSubtypeACHTransferShopifyCredit: grpcTxn.BankTransactionType_BTT_ACH_CREDIT,
	TransactionSubtypeOtherCredit:              grpcTxn.BankTransactionType_BTT_OTHER_CREDIT,
	TransactionSubtypeDepositCredit:            grpcTxn.BankTransactionType_BTT_DEPOSIT_CREDIT,
	TransactionSubtypeHoldApproved:             grpcTxn.BankTransactionType_BTT_HOLD_CREDIT,

	// Debits
	TransactionSubtypeCardPurchaseDebit:       grpcTxn.BankTransactionType_BTT_CARD_PURCHASE_DEBIT,
	TransactionSubtypeCardPurchaseDebitOnline: grpcTxn.BankTransactionType_BTT_CARD_PURCHASE_ONLINE_DEBIT,
	TransactionSubtypeWiseTransferDebit:       grpcTxn.BankTransactionType_BTT_INTRABANK_DEBIT,
	TransactionSubtypeACHTransferDebit:        grpcTxn.BankTransactionType_BTT_ACH_DEBIT,
	TransactionSubtypeWireTransferDebit:       grpcTxn.BankTransactionType_BTT_WIRE_DEBIT,
	TransactionSubtypeCardATMDebit:            grpcTxn.BankTransactionType_BTT_CARD_ATM_DEBIT,
	TransactionSubtypeCardPushDebit:           grpcTxn.BankTransactionType_BTT_CARD_PUSH_DEBIT,
	TransactionSubtypeCardPullDebit:           grpcTxn.BankTransactionType_BTT_CARD_PULL_DEBIT,
	TransactionSubtypeCheckDebit:              grpcTxn.BankTransactionType_BTT_CHECK_DEBIT,
	TransactionSubtypeFeeDebit:                grpcTxn.BankTransactionType_BTT_FEE_DEBIT,
	TransactionSubtypeACHTransferShopifyDebit: grpcTxn.BankTransactionType_BTT_ACH_DEBIT,
	TransactionSubtypeOtherDebit:              grpcTxn.BankTransactionType_BTT_OTHER_DEBIT,
}

var TransactionSubtypeFromTypeProto = map[grpcTxn.BankTransactionType]TransactionSubtype{
	grpcTxn.BankTransactionType_BTT_UNSPECIFIED: TransactionSubtypeUnspecified,

	// Origination
	grpcTxn.BankTransactionType_BTT_ORIGINATION: TransactionSubtypeAccountOriginationCredit,

	// Credits
	grpcTxn.BankTransactionType_BTT_INTEREST_CREDIT:             TransactionSubtypeInterestTransferCredit,
	grpcTxn.BankTransactionType_BTT_INTRABANK_CREDIT:            TransactionSubtypeWiseTransferCredit,
	grpcTxn.BankTransactionType_BTT_ACH_CREDIT:                  TransactionSubtypeACHTransferCredit,
	grpcTxn.BankTransactionType_BTT_WIRE_CREDIT:                 TransactionSubtypeWireTransferCredit,
	grpcTxn.BankTransactionType_BTT_CARD_PULL_CREDIT:            TransactionSubtypeCardPullCredit,
	grpcTxn.BankTransactionType_BTT_CARD_PUSH_CREDIT:            TransactionSubtypeCardPushCredit,
	grpcTxn.BankTransactionType_BTT_CARD_PAYMENT_READER_CREDIT:  TransactionSubtypeCardReaderCredit,
	grpcTxn.BankTransactionType_BTT_CARD_PAYMENT_ONLINE_CREDIT:  TransactionSubtypeCardOnlineCredit,
	grpcTxn.BankTransactionType_BTT_INTRABANK_ONLINE_CREDIT:     TransactionSubtypeBankOnlineCredit,
	grpcTxn.BankTransactionType_BTT_ACH_ONLINE_CREDIT:           TransactionSubtypeBankOnlineCredit,
	grpcTxn.BankTransactionType_BTT_CHECK_CREDIT:                TransactionSubtypeCheckCredit,
	grpcTxn.BankTransactionType_BTT_CHECK_ONLINE_CREDIT:         TransactionSubtypeCheckCredit,
	grpcTxn.BankTransactionType_BTT_CARD_PURCHASE_REFUND_CREDIT: TransactionSubtypeMerchantRefundCredit,
	grpcTxn.BankTransactionType_BTT_PROMO_CREDIT:                TransactionSubtypeWiseTransferCredit,
	grpcTxn.BankTransactionType_BTT_OTHER_CREDIT:                TransactionSubtypeOtherCredit,

	// Debits
	grpcTxn.BankTransactionType_BTT_CARD_PURCHASE_DEBIT:        TransactionSubtypeCardPurchaseDebit,
	grpcTxn.BankTransactionType_BTT_CARD_PURCHASE_ONLINE_DEBIT: TransactionSubtypeCardPurchaseDebitOnline,
	grpcTxn.BankTransactionType_BTT_INTRABANK_DEBIT:            TransactionSubtypeWiseTransferDebit,
	grpcTxn.BankTransactionType_BTT_ACH_DEBIT:                  TransactionSubtypeACHTransferDebit,
	grpcTxn.BankTransactionType_BTT_WIRE_DEBIT:                 TransactionSubtypeWireTransferDebit,
	grpcTxn.BankTransactionType_BTT_CARD_ATM_DEBIT:             TransactionSubtypeCardATMDebit,
	grpcTxn.BankTransactionType_BTT_CARD_PUSH_DEBIT:            TransactionSubtypeCardPushDebit,
	grpcTxn.BankTransactionType_BTT_CARD_PULL_DEBIT:            TransactionSubtypeCardPullDebit,
	grpcTxn.BankTransactionType_BTT_CHECK_DEBIT:                TransactionSubtypeCheckDebit,
	grpcTxn.BankTransactionType_BTT_FEE_DEBIT:                  TransactionSubtypeFeeDebit,
	grpcTxn.BankTransactionType_BTT_OTHER_DEBIT:                TransactionSubtypeOtherDebit,
}
