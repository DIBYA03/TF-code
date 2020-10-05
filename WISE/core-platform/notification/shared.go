package notification

import (
	"log"
	"regexp"
	"strings"
)

func GetMerchantName(bankDescription string) string {

	result := ""
	regEx := "[0-9]{4}(.+?)[[:space:]]{2,}"

	r, err := regexp.Compile(regEx)
	if err != nil {
		log.Println(err)
		return result
	}

	str := r.FindString(bankDescription)
	if len(str) == 0 {
		return result
	}

	result = strings.TrimSpace(string(str[4:]))

	return result
}

func (t TransactionNotification) GetMerchantName() string {
	if t.CardTransaction == nil {
		return ""
	}

	merchantName := strings.TrimSpace(t.CardTransaction.MerchantName)
	streetAddress := strings.TrimSpace(t.CardTransaction.MerchantStreetAddress)

	if len(merchantName) > 0 {
		return merchantName
	}

	if len(streetAddress) > 0 {
		if t.BankTransactionDesc != nil && len(*t.BankTransactionDesc) > 0 {
			if strings.Contains(*t.BankTransactionDesc, streetAddress) {
				return streetAddress
			}
		} else {
			return streetAddress
		}
	}

	if t.BankTransactionDesc != nil {
		return GetMerchantName(*t.BankTransactionDesc)
	}

	return ""
}

func isCheckDeposit(bankDescription string) bool {
	regEx := "[0-9]{4}[[:space:]]BRANCH DEPOSIT(.+)"

	r, err := regexp.Compile(regEx)
	if err != nil {
		log.Println(err)
		return false
	}

	str := r.FindString(bankDescription)
	if len(str) == 0 {
		return false
	}

	return true
}

func GetFeeType(bankDescription *string) string {
	result := ""
	regEx := "[0-9]{4}(.+)"

	if bankDescription == nil {
		return result
	}

	r, err := regexp.Compile(regEx)
	if err != nil {
		log.Println(err)
		return result
	}

	str := r.FindString(*bankDescription)
	if len(str) == 0 {
		return result
	}

	result = strings.TrimSpace(string(str[4:]))

	return result
}
