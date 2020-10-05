package activity

import (
	"bytes"
	"fmt"
	"text/template"
)

//TemplateName type
type TemplateName string
type Language string

type Template struct {
}

const (
	cardCreateTempl                       = TemplateName("cardCreate")
	cardStatusUpdateTempl                 = TemplateName("cardStatusUpdate")
	cardStatusUpdateBusinessNameTempl     = TemplateName("cardStatusUpdateBusinessName")
	cardAuthorizeTempl                    = TemplateName("cardAuthorize")
	cardAuthorizeGenericTempl             = TemplateName("cardAuthorizeGeneric")
	cardAuthReversalTempl                 = TemplateName("authReversal")
	cardAuthReversalGenericTempl          = TemplateName("authReversalGeneric")
	cardHoldExpireGenericTempl            = TemplateName("cardExpireGeneric")
	cardHoldApproveTempl                  = TemplateName("cardHoldApprove")
	cardHoldApproveGenericTempl           = TemplateName("cardHoldApproveGeneric")
	cardAuthorizeBusinessNameTempl        = TemplateName("cardAuthorizeBusinessName")
	cardAuthorizeBusinessNameGenericTempl = TemplateName("cardAuthorizeBusinessNameGeneric")
	cardDeclineTempl                      = TemplateName("cardDecline")
	cardDeclineGenericTempl               = TemplateName("cardDeclineGenericTempl")
	cardDeclineBusinessNameTempl          = TemplateName("cardDeclineBusinessName")
	cardDeclineBusinessNameGenericTempl   = TemplateName("cardDeclineBusinessNameGenericTempl")
	cardPostedDebitTempl                  = TemplateName("cardPostedDebit")
	cardPostedCreditTempl                 = TemplateName("cardPostedCredit")
	cardPostedCreditMerchantRefundTempl   = TemplateName("cardPostedCreditMerchantRefund")
	cardPostedDebitGenericTempl           = TemplateName("cardPostedDebitGeneric") // Without merchant name
	cardPostedDebitCardReaderTempl        = TemplateName("cardPostedDebitCardReader")
	cardPostedDebitCardReaderGenericTempl = TemplateName("cardPostedDebitCardReaderGeneric")
	cardPostedDebitCardATMTempl           = TemplateName("cardPostedDebitCardATM")
	cardPushDebitCreditTempl              = TemplateName("cardPushDebitCredit")
	cardPushDebitCreditGenericTempl       = TemplateName("cardPushDebitCreditGeneric")
	cardVisaCreditTempl                   = TemplateName("cardVisaCredit")
	cardVisaCreditGenericTempl            = TemplateName("cardVisaCreditGeneric")
)

const (
	accountOriginatedTempl         = TemplateName("accountOriginated")
	accountPostedDebitTempl        = TemplateName("accountPostedDebit")
	accountPostedDebitGenericTempl = TemplateName("accountPostedDebitGeneric")
	accountPostedCreditTempl       = TemplateName("accountPostedCredit")
	accountInProcessDebitTempl     = TemplateName("accountInProcessDebit")
	accountInProcessCreditTempl    = TemplateName("accountInProcessCredit")

	accountCardReaderCreditTempl   = TemplateName("accountRequestCardReader")
	accountCardCreditTempl         = TemplateName("accountRequestCard")
	accountBankCreditTempl         = TemplateName("accountRequestBank")
	accountWiseTransferCreditTempl = TemplateName("accountWiseTransferCredit")
	accountInterestCreditTempl     = TemplateName("accountInterestCredit")
	accountOtherCreditTempl        = TemplateName("accountOtherCredit")

	accountACHTransferShopifyCreditTempl  = TemplateName("accountACHTransferShopifyCredit")
	accountACHTransferCreditTempl         = TemplateName("accountACHTransferCredit")
	accountACHTransferCreditGenericTempl  = TemplateName("accountACHTransferCreditGeneric")
	accountWireTransferCreditTempl        = TemplateName("accountWireTransferCredit")
	accountWireTransferCreditGenericTempl = TemplateName("accountWireTransferCreditGeneric")
	accountCheckCreditTempl               = TemplateName("accountCheckCredit")
	accountDepositCreditTempl             = TemplateName("accountDepositCredit")
	accountDebitPullCreditTempl           = TemplateName("accountDebitPullCredit")
	accountDebitPullCreditGenericTempl    = TemplateName("accountDebitPullCreditGeneric")

	accountWiseTransferDebitTempl       = TemplateName("accountWiseTransferDebit")
	accountACHTransferDebitTempl        = TemplateName("accountACHTransferDebit")
	accountACHTransferDebitGenericTempl = TemplateName("accountACHTransferDebitGenericTempl")
	accountPushDebitDebitTempl          = TemplateName("accountPushDebitDebit")
	accountFeeDebitTempl                = TemplateName("accountFeeDebit")
	accountFeeDebitGenericTempl         = TemplateName("accountFeeDebitGeneric")

	accountHoldApprovedTempl = TemplateName("accountHoldApproved")
	accountHoldReleasedTempl = TemplateName("accountHoldReleased")

	accountCheckDebitTempl = TemplateName("accountCheckDebitTempl")
)

const (
	contactUpdateTempl = TemplateName("contactUpdate")
	contactCreateTempl = TemplateName("contactCreate")
	contactDeleteTempl = TemplateName("contactDelete")
)

const (
	disputeCreateTempl = TemplateName("disputeCreate")
	disputeDeleteTempl = TemplateName("disputeDelete")
)

const (
	consumerEmailUpdateTempl   = TemplateName("consumerEmailUpdate")
	consumerPhoneUpdateTempl   = TemplateName("consumerPhoneUpdate")
	consumerAddressUpdateTempl = TemplateName("consumerAddressUpdate")
)

const (
	businessEmailUpdateTempl   = TemplateName("businessEmailUpdate")
	businessPhoneUpdateTempl   = TemplateName("businessPhoneUpdate")
	businessAddressUpdateTempl = TemplateName("businessAddressUpdate")
)

const (
	LangEnglish = Language("en-US")
	LangSpanish = Language("es-SV")
)

var langs = map[Language]Language{
	LangEnglish: LangEnglish,
	LangSpanish: LangSpanish,
}

func (l Language) IsValid() bool {
	return langs[l] != ""
}

var templs = map[TemplateName]TemplateName{
	cardAuthorizeTempl:                    cardAuthorizeTempl,
	cardAuthorizeGenericTempl:             cardAuthorizeGenericTempl,
	cardDeclineTempl:                      cardDeclineTempl,
	cardDeclineGenericTempl:               cardDeclineGenericTempl,
	cardPostedDebitTempl:                  cardPostedDebitTempl,
	cardPostedCreditTempl:                 cardPostedCreditTempl,
	cardPostedDebitGenericTempl:           cardPostedDebitGenericTempl,
	accountOriginatedTempl:                accountOriginatedTempl,
	accountPostedDebitTempl:               accountPostedDebitTempl,
	accountPostedDebitGenericTempl:        accountPostedDebitGenericTempl,
	accountInProcessDebitTempl:            accountInProcessDebitTempl,
	accountPostedCreditTempl:              accountPostedCreditTempl,
	accountInProcessCreditTempl:           accountInProcessCreditTempl,
	contactUpdateTempl:                    contactUpdateTempl,
	contactCreateTempl:                    contactCreateTempl,
	contactDeleteTempl:                    contactDeleteTempl,
	disputeCreateTempl:                    disputeCreateTempl,
	disputeDeleteTempl:                    disputeDeleteTempl,
	consumerEmailUpdateTempl:              consumerEmailUpdateTempl,
	consumerPhoneUpdateTempl:              consumerPhoneUpdateTempl,
	consumerAddressUpdateTempl:            consumerAddressUpdateTempl,
	businessEmailUpdateTempl:              businessEmailUpdateTempl,
	businessPhoneUpdateTempl:              businessPhoneUpdateTempl,
	businessAddressUpdateTempl:            businessAddressUpdateTempl,
	cardCreateTempl:                       cardCreateTempl,
	cardStatusUpdateTempl:                 cardStatusUpdateTempl,
	cardStatusUpdateBusinessNameTempl:     cardStatusUpdateBusinessNameTempl,
	accountCardReaderCreditTempl:          accountCardReaderCreditTempl,
	accountCardCreditTempl:                accountCardCreditTempl,
	accountBankCreditTempl:                accountBankCreditTempl,
	accountWiseTransferCreditTempl:        accountWiseTransferCreditTempl,
	accountACHTransferCreditTempl:         accountACHTransferCreditTempl,
	accountACHTransferShopifyCreditTempl:  accountACHTransferShopifyCreditTempl,
	accountACHTransferCreditGenericTempl:  accountACHTransferCreditGenericTempl,
	accountWireTransferCreditTempl:        accountWireTransferCreditTempl,
	accountWireTransferCreditGenericTempl: accountWireTransferCreditGenericTempl,
	accountCheckCreditTempl:               accountCheckCreditTempl,
	accountDepositCreditTempl:             accountDepositCreditTempl,
	accountDebitPullCreditTempl:           accountDebitPullCreditTempl,
	accountDebitPullCreditGenericTempl:    accountDebitPullCreditGenericTempl,
	accountInterestCreditTempl:            accountInterestCreditTempl,
	accountOtherCreditTempl:               accountOtherCreditTempl,
	cardPostedCreditMerchantRefundTempl:   cardPostedCreditMerchantRefundTempl,
	cardPostedDebitCardReaderTempl:        cardPostedDebitCardReaderTempl,
	cardPostedDebitCardReaderGenericTempl: cardPostedDebitCardReaderGenericTempl,
	accountWiseTransferDebitTempl:         accountWiseTransferDebitTempl,
	accountACHTransferDebitTempl:          accountACHTransferDebitTempl,
	accountACHTransferDebitGenericTempl:   accountACHTransferDebitGenericTempl,
	accountPushDebitDebitTempl:            accountPushDebitDebitTempl,
	accountFeeDebitTempl:                  accountFeeDebitTempl,
	accountFeeDebitGenericTempl:           accountFeeDebitGenericTempl,
	cardPostedDebitCardATMTempl:           cardPostedDebitCardATMTempl,
	cardAuthorizeBusinessNameTempl:        cardAuthorizeBusinessNameTempl,
	cardAuthorizeBusinessNameGenericTempl: cardAuthorizeBusinessNameGenericTempl,
	cardAuthReversalTempl:                 cardAuthReversalTempl,
	cardAuthReversalGenericTempl:          cardAuthReversalGenericTempl,
	cardHoldApproveTempl:                  cardHoldApproveTempl,
	cardHoldApproveGenericTempl:           cardHoldApproveGenericTempl,
	cardDeclineBusinessNameTempl:          cardDeclineBusinessNameTempl,
	cardDeclineBusinessNameGenericTempl:   cardDeclineBusinessNameGenericTempl,
	cardPushDebitCreditTempl:              cardPushDebitCreditTempl,
	cardPushDebitCreditGenericTempl:       cardPushDebitCreditGenericTempl,
	cardVisaCreditTempl:                   cardVisaCreditTempl,
	cardVisaCreditGenericTempl:            cardVisaCreditGenericTempl,
	accountCheckDebitTempl:                accountCheckDebitTempl,
	accountHoldReleasedTempl:              accountHoldReleasedTempl,
	accountHoldApprovedTempl:              accountHoldApprovedTempl,
}

func (templ TemplateName) IsValid() bool {
	return templs[templ] != ""
}

var enUStransactionTempl = map[TemplateName]string{
	// Cards
	cardCreateTempl:                       "A new debit card is being produced for you.",
	cardPostedDebitTempl:                  "Your card ending in {{.Number}} was charged ${{.Amount}} at {{.Merchant}}",
	cardPostedDebitGenericTempl:           "Your card ending in {{.Number}} was charged ${{.Amount}}",
	cardDeclineTempl:                      "Your card ending in {{.Number}} was declined at {{.Merchant}} for the amount of ${{.Amount}}",
	cardDeclineGenericTempl:               "Your card ending in {{.Number}} was declined for the amount of ${{.Amount}}",
	cardStatusUpdateTempl:                 "Your card ending in {{.Number}} has been {{.Status}}",
	cardStatusUpdateBusinessNameTempl:     "{{.BusinessName}}'s card ending in {{.Number}} has been {{.Status}}",
	cardDeclineBusinessNameTempl:          "{{.BusinessName}}'s card {{.Number}} was declined at {{.Merchant}} for the amount of ${{.Amount}}",
	cardDeclineBusinessNameGenericTempl:   "{{.BusinessName}}'s card {{.Number}} was declined for the amount of ${{.Amount}}",
	cardPostedCreditTempl:                 "Your account ending in {{.AccountNumber}} has been credited ${{.Amount}}",
	cardAuthorizeTempl:                    "Your card ending in {{.Number}} was authorized at {{.Merchant}} for the amount of ${{.Amount}}",
	cardAuthorizeGenericTempl:             "Your card ending in {{.Number}} was authorized for the amount of ${{.Amount}}",
	cardPostedDebitCardReaderTempl:        "{{.BusinessName}} paid {{.Merchant}} ${{.Amount}}",
	cardPostedDebitCardReaderGenericTempl: "{{.BusinessName}} paid ${{.Amount}}",
	cardPostedDebitCardATMTempl:           "{{.BusinessName}} withdrew ${{.Amount}}",
	cardPostedCreditMerchantRefundTempl:   "Received ${{.Amount}} merchant refund",
	cardAuthorizeBusinessNameTempl:        "{{.BusinessName}}'s card ending in {{.Number}} was authorized at {{.Merchant}} for ${{.Amount}}",
	cardAuthorizeBusinessNameGenericTempl: "{{.BusinessName}}'s card ending in {{.Number}} was authorized for ${{.Amount}}",
	cardAuthReversalTempl:                 "Authorization on {{.BusinessName}}'s card ending in {{.Number}} for ${{.Amount}} at {{.Merchant}} was reversed",
	cardAuthReversalGenericTempl:          "Authorization on {{.BusinessName}}'s card ending in {{.Number}} for ${{.Amount}} was reversed",
	cardPushDebitCreditTempl:              "Received ${{.Amount}} Instant by Wise from {{.Merchant}}",
	cardPushDebitCreditGenericTempl:       "Received ${{.Amount}} Instant by Wise",
	cardVisaCreditTempl:                   "Received ${{.Amount}} visa credit from {{.Merchant}}",
	cardVisaCreditGenericTempl:            "Received ${{.Amount}} visa credit",

	// Accounts
	accountOriginatedTempl:                "Your account {{.Origin}} has been created",
	accountPostedDebitTempl:               "You sent a transfer to {{.ContactName}} for ${{.Amount}}",
	accountPostedDebitGenericTempl:        "You sent a transfer for ${{.Amount}}",
	accountInProcessDebitTempl:            "{{.BusinessName}} initiated a transfer of ${{.Amount}} to {{.ContactName}}",
	accountPostedCreditTempl:              "Your account has been credited ${{.Amount}}",
	accountInProcessCreditTempl:           "{{.Origin}} initiated a transfer of ${{.Amount}} to {{.BusinessName}}",
	accountCardReaderCreditTempl:          "{{.BusinessName}} got paid ${{.Amount}} via Card Reader",
	accountCardCreditTempl:                "{{.ContactName}} paid ${{.Amount}} invoice via Card",
	accountBankCreditTempl:                "{{.ContactName}} paid ${{.Amount}} invoice via Bank Transfer",
	accountWiseTransferCreditTempl:        "Received ${{.Amount}} Wise Transfer from {{.ContactName}}",
	accountACHTransferShopifyCreditTempl:  "Received ${{.Amount}} Payout from {{.ContactName}}",
	accountACHTransferCreditTempl:         "Received ${{.Amount}} Bank Transfer from {{.ContactName}}",
	accountACHTransferCreditGenericTempl:  "Received ${{.Amount}} Bank Transfer",
	accountWireTransferCreditTempl:        "Received ${{.Amount}} Wire Transfer from {{.ContactName}}",
	accountWireTransferCreditGenericTempl: "Received ${{.Amount}} Wire Transfer",
	accountCheckCreditTempl:               "Received ${{.Amount}} Check Deposit",
	accountDepositCreditTempl:             "Received ${{.Amount}} Bank Deposit",
	accountDebitPullCreditTempl:           "Received ${{.Amount}} Instant by Wise from a visa debit card ending in {{.ContactName}}",
	accountDebitPullCreditGenericTempl:    "Received ${{.Amount}} Instant by Wise",
	accountInterestCreditTempl:            "{{.BusinessName}} earned ${{.Amount}} in interest for {{.InterestEarnedMonth}}",
	accountOtherCreditTempl:               "Received ${{.Amount}} Credit",
	accountWiseTransferDebitTempl:         "{{.BusinessName}} sent ${{.Amount}} to {{.ContactName}}",
	accountACHTransferDebitTempl:          "{{.BusinessName}} sent ${{.Amount}} to {{.ContactName}}",
	accountACHTransferDebitGenericTempl:   "{{.BusinessName}} sent ${{.Amount}}",
	accountPushDebitDebitTempl:            "{{.BusinessName}} sent ${{.Amount}} to {{.ContactName}} via Instant by Wise",
	accountFeeDebitTempl:                  "{{.BusinessName}} paid {{.ContactName}} of ${{.Amount}}",
	accountFeeDebitGenericTempl:           "{{.BusinessName}} paid fee of ${{.Amount}}",
	accountCheckDebitTempl:                "{{.BusinessName}} paid {{.ContactName}} ${{.Amount}} via check",
	accountHoldReleasedTempl:              "Hold released on {{.BusinessName}}'s account for ${{.Amount}}",
	accountHoldApprovedTempl:              "Hold placed on {{.BusinessName}}'s account for ${{.Amount}}",

	// Contacts
	contactUpdateTempl: "You updated {{.Name}}'s details",
	contactCreateTempl: "You added {{.Name}} as a new contact",
	contactDeleteTempl: "You deleted {{.Name}}",

	// Dispute
	disputeCreateTempl: "You have disputed ${{.Amount}} transaction as {{.Category}}",
	disputeDeleteTempl: "You have cancelled ${{.Amount}} dispute",

	// Consumer
	consumerEmailUpdateTempl:   "Your email has been updated to {{.Email}}",
	consumerPhoneUpdateTempl:   "Your phone number has been updated to {{.Phone}}",
	consumerAddressUpdateTempl: "Your {{.Type}} address has been updated to {{.Line1}}, {{.City}}, {{.State}} {{.ZipCode}}",

	// Business
	businessEmailUpdateTempl:   "Your business email has been updated to {{.Email}}",
	businessPhoneUpdateTempl:   "Your business phone number has been updated to {{.Phone}}",
	businessAddressUpdateTempl: "Your business {{.Type}} address has been updated to {{.Line1}}, {{.City}}, {{.State}} {{.ZipCode}}",
}

var esSVtransactionTempl = map[TemplateName]string{
	cardAuthorizeTempl:       "{{.Amount}} autorizado en {{.Merchant}}",
	cardDeclineTempl:         "Tu tarjeta fue genada en {{.Merchant}} por la cantidad de {{.Amount}}",
	cardPostedDebitTempl:     "Tu tarjeta fue cargado {{.Amount}} en {{.Merchant}}",
	cardPostedCreditTempl:    "Tu tarjeta fue creditada {{.Amount}}",
	accountPostedDebitTempl:  "Tu cuenta fue cargada {{.Amount}} en {{.Merchant}}",
	accountPostedCreditTempl: "Tu cuenta fue creditada {{.Amount}}",
	contactUpdateTempl:       "Modificaste la informacion de {{.Name}}",
	contactCreateTempl:       "Creaste {{.Name}} como un {{.Category}}",
	contactDeleteTempl:       "Borraste a {{.Name}}",
}

func (t Template) NewName(v string) TemplateName {
	return TemplateName(v)
}

//NewWithLang returns a template for the specified language
func (t Template) NewWithLang(templ TemplateName, lang Language, v interface{}) (string, error) {
	println("new with lang", templ, lang)

	if !templ.IsValid() {
		return "", fmt.Errorf("template %v is not a valid template", v)
	}

	switch lang {
	case LangEnglish:
		return t.parse(templ, v, enUStransactionTempl)
	case LangSpanish:
		return t.parse(templ, v, esSVtransactionTempl)
	default:
		return t.parse(templ, v, enUStransactionTempl)
	}
}

func (t Template) parse(name TemplateName, v interface{}, m map[TemplateName]string) (string, error) {
	println("parsing...")

	temp, ok := m[name]
	if !ok {
		return "", fmt.Errorf("Unable to parse %v into %v", name, v)
	}

	templ, err := template.New(string(name)).Parse(temp)
	buf := new(bytes.Buffer)
	err = templ.Execute(buf, v)

	println("parsed string is ", buf.String())

	return buf.String(), err
}
