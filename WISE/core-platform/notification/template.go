package notification

// Account and debit card notifications

const ApplicationApproved = `Congratulations! Your account application has been approved and your Wise checking account has been created.`
const DebitCardActivated = `%s's card ending in %s has been activated.`
const DebitCardBlocked = `%s's card ending in %s has been blocked.`
const DebitCardUnblocked = `%s's card ending in %s has been unblocked.`
const DebitCardShipped = `%s's card ending in %s has been shipped.`

// Account origination
const AccountOriginationTransactionTitle = `Wise bank account created`
const AccountOriginationTransactionDescription = `Your Wise bank account ending in %s has been created`

// Card declined and merchant name is available
const DebitCardDeclinedWithMerchant = `Your card %s was declined at %s for the amount of $%s`

// Card declined but merchant name not available
const DebitCardDeclinedWithoutMerchant = `Your card %s was declined for the amount of $%s`

// Card charged and merchant name is available
const CardDebitPostedWithMerchant = `Your card %s was charged $%s at %s`

// Card charged but merchant name is not available
const CardDebitPostedWithoutMerchant = `Your card %s was charged $%s`

// Amount credited to card - refunds, etc..
const CardCreditPosted = `Your account %s has been credited $%s`

const AccountCreditPosted = `Your account %s has been credited $%s`

const AccountDebitPosted = `You sent a transfer to %s for amount $%s`
const AccountDebitPostedWithoutContact = `You sent a transfer for amount $%s`

const AccountDebitInProcess = `You initiated a transfer to %s for amount $%s`

const AccountCreditInProcess = `You initiated a transfer from account %s for $%s`

const CreditViaCardReaderNotificationHeader = `%s got paid $%s!`
const CreditViaCardReaderNotificationBody = `%s got paid via %s in %s. $%s instantly available in %s's Wise account ending in %s.`
const CreditViaCardReaderTransactionTitle = `%s got paid`
const CreditViaCardReaderTransactionDescription = `%s got paid via %s in %s`

const CreditViaCardNotificationHeader = `%s got paid $%s!`
const CreditViaCardNotificationBody = `%s paid %s's invoice via %s. $%s instantly available in %s's Wise account ending in %s.`
const CreditViaCardTransactionTitle = `%s paid %s's Invoice`
const CreditViaCardTransactionDescription = `%s paid %s via %s`

const CreditViaBankNotificationHeader = `%s got paid $%s!`
const CreditViaBankNotificationBody = `%s paid %s's invoice via %s. $%s instantly available in %s's Wise account ending in %s.`
const CreditViaBankTransactionTitle = `%s paid %s's Invoice`
const CreditViaBankTransactionDescription = `%s paid %s via %s`

const CreditWiseTransferNotificationHeader = `%s received $%s!`
const CreditWiseTransferNotificationBody = `%s received a Wise Transfer from %s. $%s instantly available in %s's Wise account ending in %s.`
const CreditWiseTransferTransactionTitle = `Funds from %s`
const CreditWiseTransferTransactionDescription = `Received $%s Wise Transfer from %s`

const CreditACHNotificationHeader = `%s received $%s!`
const CreditACHNotificationBody = `%s received a Bank Transfer from %s. $%s instantly available in %s's Wise account ending in %s.`
const CreditACHNotificationWithoutSenderBody = `%s received a Bank Transfer. $%s instantly available in %s's Wise account ending in %s.`
const CreditACHTransactionTitle = `Funds from %s`
const CreditACHTransactionWithoutSenderTitle = `Received Bank Transfer`
const CreditACHTransactionDescription = `Received $%s Bank Transfer from %s`
const CreditACHTransactionWithoutSenderDescription = `Received $%s Bank Transfer`

const CreditDepositWireNotificationBody = `%s received a Wire Transfer from %s. $%s instantly available in %s's Wise account ending in %s.`
const CreditDepositWireTransactionTitle = `Received Wire Transfer from %s`
const CreditDepositWireTransactionDescription = `Received $%s Wire Transfer from %s`

const CreditDepositCheckNotificationBody = `Check for $%s deposited in %s's Wise account ending in %s.`
const CreditDepositCheckTransactionTitle = `Received Check Deposit`
const CreditDepositCheckTransactionDescription = `Received $%s Check Deposit`

const CreditACHShopifyNotificationHeader = `%s received $%s!`
const CreditACHShopifyNotificationBody = `%s received a Payout from %s. $%s instantly available in %s's Wise account ending in %s.`
const CreditACHShopifyTransactionTitle = `Payout from %s`
const CreditACHShopifyTransactionDescription = `Received $%s Payout from %s`

const CreditDepositNotificationHeader = `%s received $%s!`
const CreditDepositNotificationBody = `%s received a Deposit. $%s instantly available in %s's Wise account ending in %s.`
const CreditDepositTransactionTitle = `Received Deposit`
const CreditDepositTransactionDescription = `Received $%s Deposit`

const CreditInterestTransferNotificationHeader = `%s earned $%s in interest!`
const CreditInterestTransferNotificationBody = `%s's Wise account ending in %s earned $%s in interest for %s`
const CreditInterestTransferTransactionTitle = `Interest earned for %s`
const CreditInterestTransferTransactionDescription = `%s earned $%s in interest for %s`

const CreditMerchantRefundNotificationHeader = `%s received a merchant refund of $%s!`
const CreditMerchantRefundNotificationBody = `%s received a merchant refund from %s. $%s instantly available in %s's Wise account ending in %s`
const CreditMerchantRefundNotificationWithoutMerchantBody = `%s received a merchant refund. $%s instantly available in %s's Wise account ending in %s`
const CreditMerchantRefundTransactionTitle = `Received merchant refund from %s`
const CreditMerchantRefundTransactionWithoutMerchantTitle = `Received merchant refund`
const CreditMerchantRefundTransactionDescription = `Received $%s merchant refund`

const CreditCardInstantPayNotificationHeader = `%s received $%s`
const CreditCardInstantPayNotificationBody = `%s received a Instant by Wise Transfer from %s. $%s instantly available in %s's Wise account ending in %s`
const CreditCardInstantPayNotificationWithoutSenderBody = `%s received a Instant by Wise. $%s instantly available in %s's Wise account ending in %s`
const CreditCardInstantPayTransactionTitle = `Received Instant by Wise from %s`
const CreditCardInstantPayTransactionWithoutSenderTitle = `Received Instant by Wise`
const CreditCardInstantPayTransactionDescription = `Received $%s Instant by Wise`

const CreditCardDebitPullNotificationHeader = `%s received $%s`
const CreditCardDebitPullNotificationBody = `%s received $%s via Instant by Wise from a visa debit card ending in %s. Funds are now available in your account ending in %s`
const CreditCardDebitPullNotificationWithoutSenderBody = `%s received $%s via Instant by Wise. Funds are now available in your account ending in %s`
const CreditCardDebitPullTransactionTitle = `Received Instant by Wise from a visa debit card ending in %s`
const CreditCardDebitPullTransactionWithoutSenderTitle = `Received Instant by Wise`
const CreditCardDebitPullTransactionDescription = `Received $%s Instant by Wise`

const CreditCardVisaCreditNotificationHeader = `%s received $%s`
const CreditCardVisaCreditNotificationBody = `%s received $%s via visa credit from %s. Funds are now available in your account ending in %s`
const CreditCardVisaCreditNotificationWithoutSenderBody = `%s received $%s via visa credit. Funds are now available in your account ending in %s`
const CreditCardVisaCreditTransactionTitle = `Received visa credit from %s`
const CreditCardVisaCreditTransactionWithoutSenderTitle = `Received visa credit`
const CreditCardVisaCreditTransactionDescription = `Received $%s visa credit`

const DebitCardPurchaseNotificationHeader = `%s paid %s $%s`
const DebitCardPurchaseNotificationWithoutMerchantHeader = `%s paid $%s`
const DebitCardPurchaseNotificationBody = `%s paid %s $%s in %s via card ending in %s`
const DebitCardPurchaseNotificationWithoutMerchantBody = `%s paid $%s in %s via card ending in %s`
const DebitCardPurchaseNotificationWithoutLocationBody = `%s paid %s $%s via card ending in %s`
const DebitCardPurchaseNotificationWithoutCardNumberBody = `%s paid %s $%s via card`
const DebitCardPurchaseNotificationGenericBody = `%s paid $%s via card ending in %s`
const DebitCardPurchaseNotificationGenericWithoutCardNumberBody = `%s paid $%s via card`
const DebitCardPurchaseTransactionTitle = `Paid %s`
const DebitCardPurchaseTransactionWithoutMerchantTitle = `Paid $%s`
const DebitCardPurchaseTransactionDescription = `%s paid %s $%s in %s via card ending in %s`
const DebitCardPurchaseTransactionWithoutMerchantDescription = `%s paid $%s in %s via card ending in %s`
const DebitCardPurchaseTransactionWithoutLocationDescription = `%s paid %s $%s via card ending in %s`
const DebitCardPurchaseTransactionWithoutCardNumberDescription = `%s paid %s $%s via card`
const DebitCardPurchaseTransactionGenericDescription = `%s paid $%s via card ending in %s`
const DebitCardPurchaseTransactionGenericWithoutCardNumberDescription = `%s paid $%s via card`

const DebitCardATMNotificationHeader = `%s withdrew $%s`
const DebitCardATMNotificationBody = `%s withdrew $%s in %s via card ending in %s`
const DebitCardATMNotificationWithoutCardNumberBody = `%s withdrew $%s in %s via card`
const DebitCardATMNotificationWithoutLocationBody = `%s withdrew $%s via card ending in %s`
const DebitCardATMNotificationGenericBody = `%s withdrew $%s via card`
const DebitCardATMTransactionTitle = `Withdrew at ATM`
const DebitCardATMTransactionDescription = `%s withdrew $%s in %s via card ending in %s`
const DebitCardATMTransactionWithoutCardNumberDescription = `%s withdrew $%s in %s via card`
const DebitCardATMTransactionWithoutLocationDescription = `%s withdrew $%s via card ending in %s`
const DebitCardATMTransactionGenericDescription = `%s withdrew $%s via card`

const DebitFeeNotificationHeader = `%s charged $%s`
const DebitFeeNotificationBody = `%s charged a %s of $%s`
const DebitFeeNotificationGenericBody = `%s charged a fee of $%s`
const DebitFeeTransactionTitle = `Charged %s`
const DebitFeeTransactionDescription = `%s charged $%s as %s`
const DebitFeeTransactionWithoutTypeTitle = `Charged fee`
const DebitFeeTransactionWithoutTypeDescription = `%s charged $%s as fee`

const DebitWiseTransferNotificationHeader = `%s sent $%s`
const DebitWiseTransferNotificationBody = `%s sent a Wise Transfer to %s. $%s debited from %s's Wise account ending in %s`
const DebitWiseTransferNotificationWithoutContactBody = `%s sent a Wise Transfer. $%s debited from %s's Wise account ending in %s`
const DebitWiseTransferTransactionTitle = `Sent to %s`
const DebitWiseTransferTransactionWithoutContactTitle = `Sent a Wise Transfer`
const DebitWiseTransferTransactionDescription = `%s sent a Wise Transfer to %s`
const DebitWiseTransferTransactionWithoutContactDescription = `%s sent a Wise Transfer`

const DebitWiseCheckNotificationHeader = `%s paid %s $%s`
const DebitWiseCheckNotificationBody = `%s paid %s $%s via check`
const DebitWiseCheckTransactionTitle = `Sent to %s`
const DebitWiseCheckTransactionDescription = `%s sent a Check Transfer to %s`

const DebitACHNotificationHeader = `%s sent $%s`
const DebitACHNotificationBody = `%s sent a Bank Transfer to %s. $%s debited from %s's Wise account ending in %s`
const DebitACHTransactionTitle = `Sent to %s`
const DebitACHTransactionDescription = `%s sent a Bank Transfer to %s`

const DebitExternalACHNotificationHeader = `%s sent $%s`
const DebitExternalACHNotificationBody = `%s sent a Bank Transfer. $%s debited from %s's Wise account ending in %s`
const DebitExternalACHTransactionTitle = `Sent a Bank Transfer`
const DebitExternalACHTransactionDescription = `Sent a Bank Transfer`

const DebitCardInstantPayNotificationHeader = `%s sent $%s`
const DebitCardInstantPayNotificationBody = `%s sent an Instant by Wise Transfer to %s. $%s debited from %s's Wise account ending in %s`
const DebitCardInstantPayTransactionTitle = `Sent to %s`
const DebitCardInstantPayTransactionDescription = `%s sent Instant by Wise Transfer to %s`

// Pending Transactions
const DebitCardAuthorizationNotificationHeader = `%s's card ending in %s was authorized`
const DebitCardAuthorizationNotificationBody = `%s's card ending in %s was authorized at %s in %s for the amount of $%s`
const DebitCardAuthorizationNotificationWithoutMerchantBody = `%s's card %s was authorized in %s for the amount of $%s`
const DebitCardAuthorizationNotificationWithoutLocationBody = `%s's card %s was authorized at %s for the amount of $%s`
const DebitCardAuthorizationNotificationGenericBody = `%s's card ending in %s was authorized for the amount of $%s`
const DebitCardAuthorizationTransactionTitle = `Authorized at %s`
const DebitCardAuthorizationTransactionWithoutMerchantTitle = `Authorized $%s`
const DebitCardAuthorizationTransactionDescription = `%s's card ending in %s was authorized at %s in %s for the amount of $%s`
const DebitCardAuthorizationTransactionWithoutMerchantDescription = `%s's card ending in %s was authorized in %s for the amount of $%s`
const DebitCardAuthorizationTransactionWithoutLocationDescription = `%s's card ending in %s was authorized at %s for the amount of $%s`
const DebitCardAuthorizationTransactionGenericDescription = `%s's card ending in %s was authorized for the amount of $%s`

const DebitCardDeclinedNotificationHeader = `%s's card ending in %s was declined`
const DebitCardDeclinedNotificationBody = `%s's card ending in %s was declined at %s in %s for the amount of $%s`
const DebitCardDeclinedNotificationWithoutMerchantBody = `%s's card ending in %s was declined in %s for the amount of $%s`
const DebitCardDeclinedNotificationWithoutLocationBody = `%s's card ending in %s was declined at %s for the amount of $%s`
const DebitCardDeclinedNotificationGenericBody = `%s's card ending in %s was declined for the amount of $%s`

const AccountHoldNotificationHeader = `Hold placed for $%s`
const AccountHoldNotificationBody = `Hold placed for the amount of $%s on %s's account ending in %s`
const AccountHoldTransactionTitle = `Hold placed for $%s`
const AccountHoldTransactionDescription = `Hold for the amount of $%s placed on %s's account ending in %s`

const DebitCardDeclinedTransactionTitle = `Declined at %s`
const DebitCardDeclinedTransactionWithoutMerchantTitle = `Declined $%s`
const DebitCardDeclinedTransactionDescription = `%s's card ending in %s was declined at %s in %s for the amount of $%s`
const DebitCardDeclinedTransactionWithoutMerchantDescription = `%s's card ending in %s was declined in %s for the amount of $%s`
const DebitCardDeclinedTransactionWithoutLocationDescription = `%s's card ending in %s was declined at %s for the amount of $%s`
const DebitCardDeclinedTransactionGenericDescription = `%s's card ending in %s was declined for the amount of $%s`

const OtherCreditNotificationHeader = `%s received $%s!`
const OtherCreditNotificationBody = `%s received a Credit. $%s instantly available in %s's Wise account ending in %s.`
const OtherCreditTransactionTitle = `Received Credit`
const OtherCreditTransactionDescription = `Received $%s Credit`
