package transaction

import (
	"log"
	"regexp"
	"strings"
)

func GetOriginAccount(transactionType TransactionType, description string) *string {
	result := ""
	var regEx string
	switch transactionType {
	case TransactionTypeTransfer:
		fallthrough
	case TransactionTypeACH:
		regEx = "ORIG(.+)DEST"
	default:
		regEx = "ORIG(.+)DEST"
	}

	r, err := regexp.Compile(regEx)
	if err != nil {
		log.Println(err)
		return &result
	}

	str := r.FindString(description)
	if len(str) == 0 {
		result = GetAccount(description)
		return &result
	}

	words := strings.Fields(str)
	if len(words) > 1 {
		return &words[1]
	}

	return &result
}
func GetOriginAccountHolder(transactionType TransactionType, description string) *string {
	result := ""
	var regEx string
	switch transactionType {
	case TransactionTypeTransfer:
		fallthrough
	case TransactionTypeACH:
		regEx = "ORIG(.+)DEST"
	case TransactionTypeDeposit:
		regEx = "ORG(.+)"
	default:
		regEx = "ORIG(.+)DEST"
	}

	r, err := regexp.Compile(regEx)
	if err != nil {
		log.Println(err)
		return &result
	}

	str := r.FindString(description)
	if len(str) == 0 {
		result = GetSenderName(description)
		return &result
	}

	words := strings.Fields(str)
	if transactionType == TransactionTypeDeposit {
		if len(words) > 2 {
			result = (words[1] + " " + words[2])
			return &result
		}
	} else {
		if len(words) > 4 {
			result = (words[2] + " " + words[3])
			return &result
		}
	}

	return &result
}
func GetDestinationAccount(transactionType TransactionType, description string) *string {
	result := ""
	var regEx string
	switch transactionType {
	case TransactionTypeTransfer:
		fallthrough
	case TransactionTypeACH:
		regEx = "DEST(.+)"
	default:
		regEx = "DEST(.+)"
	}

	r, err := regexp.Compile(regEx)
	if err != nil {
		log.Println(err)
		return &result
	}

	str := r.FindString(description)
	if len(str) == 0 {
		result = GetAccount(description)
		return &result
	}

	words := strings.Fields(str)
	if len(words) > 1 {
		return &words[1]
	}

	return &result
}
func GetDestinationAccountHolder(transactionType TransactionType, description string) *string {
	result := ""
	var regEx string
	switch transactionType {
	case TransactionTypeTransfer:
		fallthrough
	case TransactionTypeACH:
		regEx = "DEST(.+)"
	default:
		regEx = "DEST(.+)"
	}

	r, err := regexp.Compile(regEx)
	if err != nil {
		log.Println(err)
		return &result
	}

	str := r.FindString(description)
	if len(str) == 0 {
		result = GetReceiverName(description)
		return &result
	}

	words := strings.Fields(str)
	if len(words) > 3 {
		result = (words[2] + " " + words[3])
		return &result
	}

	return &result
}
func GetExternalACHDestination(description string) string {

	result := ""
	description = StandardizeSpaces(description)
	regEx := "DEBIT FOR(.+)CO REF" //DEBIT FOR  <receiver name>     <transfer type>  CO REF
	r, err := regexp.Compile(regEx)
	if err != nil {
		log.Println(err)
		return result
	}

	str := r.FindString(description)
	if len(str) == 0 {
		return GetReceiverName(description)
	}

	words := strings.Fields(str)
	if len(words) > 5 {
		for i := 2; i < len(words)-3; i++ {
			result = result + words[i] + " "
		}
	}

	result = strings.TrimSpace(result)
	return result
}

func GetExternalACHSource(description string) *string {

	result := ""
	description = StandardizeSpaces(description)
	regEx := "CREDIT FOR(.+)CO REF" //CREDIT FOR  <receiver name>     <transfer type>  CO REF

	r, err := regexp.Compile(regEx)
	if err != nil {
		log.Println(err)
		return &result
	}

	str := r.FindString(description)
	if len(str) == 0 {
		result = GetReceiverName(description)
		return &result
	}

	words := strings.Fields(str)
	if len(words) > 5 {
		for i := 2; i < len(words)-3; i++ {
			result = result + words[i] + " "
		}
	}

	result = strings.TrimSpace(result)

	return &result
}

func StandardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
func GetInstantPaySenderName(description *string) string {
	result := ""
	if description == nil {
		return result
	}

	words := strings.Split(*description, " ")
	if len(words) > 2 {
		result = words[1] + " " + words[2]
	}

	return result
}
func GetSenderName(description string) string {
	result := ""
	description = StandardizeSpaces(description)
	regEx := "From(.+)(Account:|Card:)"
	r, err := regexp.Compile(regEx)
	if err != nil {
		log.Println(err)
		return result
	}

	str := r.FindString(description)
	if len(str) == 0 {
		return result
	}

	words := strings.Fields(str)
	if len(words) > 2 {
		for i := 1; i < len(words)-1; i++ {
			result = result + words[i] + " "
		}
	}

	result = strings.TrimSpace(result)
	return result
}
func GetReceiverName(description string) string {
	result := ""
	description = StandardizeSpaces(description)
	regEx := "To(.+)(Account:|Card:)"
	r, err := regexp.Compile(regEx)
	if err != nil {
		log.Println(err)
		return result
	}

	str := r.FindString(description)
	if len(str) == 0 {
		return result
	}

	words := strings.Fields(str)
	if len(words) > 2 {
		for i := 1; i < len(words)-1; i++ {
			result = result + words[i] + " "
		}
	}

	result = strings.TrimSpace(result)
	return result
}

func GetAccountNumber(description string) string {
	result := ""
	description = StandardizeSpaces(description)
	regEx := "(Account:|Card:)(.+)"
	r, err := regexp.Compile(regEx)
	if err != nil {
		log.Println(err)
		return result
	}

	str := r.FindString(description)
	if len(str) == 0 {
		return result
	}

	words := strings.Fields(str)
	if len(words) > 1 {
		account := words[1]
		result = account[4:]
	}

	result = strings.TrimSpace(result)
	return result
}

func GetAccount(description string) string {

	result := ""
	description = StandardizeSpaces(description)
	regEx := "(Account:|Card:)(.+)"

	r, err := regexp.Compile(regEx)
	if err != nil {
		log.Println(err)
		return result
	}

	str := r.FindString(description)
	if len(str) == 0 {
		return result
	}

	words := strings.Fields(str)
	if len(words) > 1 {
		account := words[1]
		result = account[4:]
	}

	result = strings.TrimSpace(result)

	return result
}

func GetVisaCreditSenderName(description *string) string {
	result := ""
	if description == nil {
		return result
	}

	words := strings.Split(*description, " ")
	if len(words) > 1 {
		for i := 1; i < len(words); i++ {
			result = result + words[i] + " "
		}
	}

	result = strings.TrimSpace(result)

	return result
}
