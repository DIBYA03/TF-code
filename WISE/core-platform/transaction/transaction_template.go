package transaction

import (
	"fmt"

	bsrv "github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/go-lib/locale"
	grpcTxn "github.com/wiseco/protobuf/golang/transaction"
	grpcBankTxn "github.com/wiseco/protobuf/golang/transaction/bank"
)

const AccountOriginationTransactionTitle = "Wise bank account created"
const InterestTransferCreditTransactionTitle = "Interest earned for %s"
const IntrabankTransferCreditTransactionTitle = "Funds from %s"
const ACHCreditTransactionTitle = "Funds from %s"
const ACHCreditTransactionTitleWithoutCounterparty = "Received Bank Transfer"
const DepositCreditTransactionTitle = "Received Wire Transfer from %s"
const DepositCreditTransactionTitleWithoutCounterparty = "Received Wire Transfer"
const DepositCreditTransaction = "Received Deposit"
const CardInstantPayCreditTransactionTitle = "Received Instant by Wise from %s"
const CardInstantPayCreditTransactionTitleWithoutCounterparty = "Received Instant by Wise"
const CardDebitPullCreditTransactionTitle = "Received Instant by Wise from %s"
const CardDebitPullCreditTransactionTitleWithoutCounterparty = "Received Instant by Wise"
const ViaCardReaderCreditTransactionTitle = "%s got paid"
const ViaCardReaderCreditTransactionTitleNoBusiness = "Paid via reader"
const ViaCardCreditTransactionTitle = "%s paid %s's Invoice"
const ViaCardCreditTransactionTitleNoBusiness = "Paid %s's Invoice"
const ViaBankCreditTransactionTitle = "%s paid %s's Invoice"
const ViaBankCreditTransactionTitleNoBusiness = "Paid %s's Invoice"
const CardPurchaseRefundCreditTransactionTitle = "Received merchant refund from %s"
const CardPurchaseRefundCreditTransactionTitleWithoutCounterparty = "Received merchant refund"
const PromoCreditTransactionTitle = "Received promo credit"
const CardPurchaseDebitTransactionTitle = "Paid %s"
const CardPurchaseDebitTransactionTitleWithoutCounterparty = "Paid $%s"
const IntrabankTransferDebitTransactionTitle = "Sent to %s"
const ACHDebitTransactionTitle = "Sent to %s"
const ACHDebitExternalTransactionTitle = "Sent a Bank Transfer"
const CardATMDebitTransactionTitle = "Withdrew at ATM"
const CheckDebitTransactionTitle = "Sent to %s"
const CardInstantPayDebitTransactionTitle = "Sent to %s"
const CardFeeDebitTransactionTitle = "Card fee"

const DebitTransactionTitle = "Debit sent to %s"
const DebitTransactionTitleWithoutCounterparty = "Debit sent"

const CreditTransactionTitle = "Credit received from %s"
const CreditTransactionTitleWithoutCounterparty = "Credit received"
const CreditTransactionTitleCheck = "Deposit received"

const ShopifyACHCreditTransactionTitle = "Payout from %s"

// Pending Transactions
const CardAuthorizationDebitTransactionTitle = "Authorized at %s"
const CardAuthorizationDebitTransactionTitleWithoutCounterparty = "Authorized $%s"
const HoldAuthorizationDebitTransactionTitle = "Hold placed for $%s"

func BankTransactionDisplayTitle(t *grpcBankTxn.Transaction, bus *bsrv.Business) string {
	if t.LegacyTitle != "" {
		return t.LegacyTitle
	}

	switch t.Type {
	case grpcTxn.BankTransactionType_BTT_ORIGINATION:
		return AccountOriginationTransactionTitle
	case grpcTxn.BankTransactionType_BTT_INTEREST_CREDIT:
		d, err := locale.ParseDate(t.InterestDate)
		if err != nil {
			return CreditTransactionTitleWithoutCounterparty
		}

		return fmt.Sprintf(InterestTransferCreditTransactionTitle, d.Format())
	case grpcTxn.BankTransactionType_BTT_INTRABANK_CREDIT:
		return fmt.Sprintf(IntrabankTransferCreditTransactionTitle, t.Counterparty)
	case grpcTxn.BankTransactionType_BTT_ACH_CREDIT:
		if t.CounterpartyType == grpcTxn.BankTransactionCounterpartyType_BTCT_SHOPPING_SHOPIFY_PAYOUT {
			return fmt.Sprintf(ShopifyACHCreditTransactionTitle, t.Counterparty)
		} else if t.Counterparty != "" {
			return fmt.Sprintf(ACHCreditTransactionTitle, t.Counterparty)
		} else {
			return ACHCreditTransactionTitleWithoutCounterparty
		}
	case grpcTxn.BankTransactionType_BTT_WIRE_CREDIT:
		if t.Counterparty != "" {
			return fmt.Sprintf(DepositCreditTransactionTitle, t.Counterparty)
		} else {
			return DepositCreditTransactionTitleWithoutCounterparty
		}
	case grpcTxn.BankTransactionType_BTT_CARD_PULL_CREDIT:
		if t.Counterparty != "" {
			return fmt.Sprintf(CardDebitPullCreditTransactionTitle, t.Counterparty)
		} else {
			return CardDebitPullCreditTransactionTitleWithoutCounterparty
		}
	case grpcTxn.BankTransactionType_BTT_CARD_PUSH_CREDIT:
		if t.Counterparty != "" {
			return fmt.Sprintf(CardInstantPayCreditTransactionTitle, t.Counterparty)
		} else {
			return CardInstantPayCreditTransactionTitleWithoutCounterparty
		}
	case grpcTxn.BankTransactionType_BTT_CARD_PAYMENT_READER_CREDIT, grpcTxn.BankTransactionType_BTT_CARD_PAYMENT_ONLINE_CREDIT:
		if t.Counterparty != "" {
			if bus != nil {
				return fmt.Sprintf(ViaCardCreditTransactionTitle, t.Counterparty, bus.Name())
			} else {
				return fmt.Sprintf(ViaCardCreditTransactionTitleNoBusiness, t.Counterparty)
			}
		} else {
			if bus != nil {
				return fmt.Sprintf(ViaCardReaderCreditTransactionTitle, bus.Name())
			} else {
				return ViaCardReaderCreditTransactionTitleNoBusiness
			}
		}
	case grpcTxn.BankTransactionType_BTT_INTRABANK_ONLINE_CREDIT, grpcTxn.BankTransactionType_BTT_ACH_ONLINE_CREDIT:
		if bus != nil {
			return fmt.Sprintf(ViaBankCreditTransactionTitle, t.Counterparty, bus.Name())
		} else {
			return fmt.Sprintf(ViaBankCreditTransactionTitleNoBusiness, t.Counterparty)
		}
	case grpcTxn.BankTransactionType_BTT_CARD_PURCHASE_REFUND_CREDIT:
		if t.Counterparty != "" {
			return fmt.Sprintf(CardPurchaseRefundCreditTransactionTitle, t.Counterparty)
		} else {
			return CardPurchaseRefundCreditTransactionTitleWithoutCounterparty
		}
	case grpcTxn.BankTransactionType_BTT_PROMO_CREDIT:
		return PromoCreditTransactionTitle
	case grpcTxn.BankTransactionType_BTT_CARD_PURCHASE_DEBIT, grpcTxn.BankTransactionType_BTT_CARD_PURCHASE_ONLINE_DEBIT:
		if t.Counterparty != "" {
			return fmt.Sprintf(CardPurchaseDebitTransactionTitle, t.Counterparty)
		} else {
			return CardPurchaseDebitTransactionTitleWithoutCounterparty
		}
	case grpcTxn.BankTransactionType_BTT_INTRABANK_DEBIT:
		return fmt.Sprintf(IntrabankTransferDebitTransactionTitle, t.Counterparty)
	case grpcTxn.BankTransactionType_BTT_ACH_DEBIT:
		if t.Counterparty != "" {
			return fmt.Sprintf(ACHDebitTransactionTitle, t.Counterparty)
		} else {
			return ACHDebitExternalTransactionTitle
		}
	case grpcTxn.BankTransactionType_BTT_CARD_ATM_DEBIT:
		return CardATMDebitTransactionTitle
	case grpcTxn.BankTransactionType_BTT_CARD_PUSH_DEBIT, grpcTxn.BankTransactionType_BTT_CARD_PULL_DEBIT:
		return fmt.Sprintf(CardInstantPayDebitTransactionTitle, t.Counterparty)
	case grpcTxn.BankTransactionType_BTT_CHECK_DEBIT:
		return fmt.Sprintf(CheckDebitTransactionTitle, t.Counterparty)
	case grpcTxn.BankTransactionType_BTT_FEE_DEBIT:
		return CardFeeDebitTransactionTitle
	case grpcTxn.BankTransactionType_BTT_HOLD_CREDIT:
		return fmt.Sprintf(HoldAuthorizationDebitTransactionTitle, t.Amount)
	case grpcTxn.BankTransactionType_BTT_DEPOSIT_CREDIT:
		return DepositCreditTransaction
	case grpcTxn.BankTransactionType_BTT_CHECK_CREDIT:
		return CreditTransactionTitleCheck
	case
		grpcTxn.BankTransactionType_BTT_CHECK_ONLINE_CREDIT,
		grpcTxn.BankTransactionType_BTT_WIRE_DEBIT,
		grpcTxn.BankTransactionType_BTT_CARD_PAYMENT_REFUND_DEBIT:
		fallthrough
	default:
		if t.Amount[:1] == "-" {
			if t.Counterparty != "" {
				return DebitTransactionTitle
			} else {
				return DebitTransactionTitleWithoutCounterparty
			}
		} else {
			if t.Counterparty != "" {
				return CreditTransactionTitle
			} else {
				return CreditTransactionTitleWithoutCounterparty
			}
		}
	}
}
