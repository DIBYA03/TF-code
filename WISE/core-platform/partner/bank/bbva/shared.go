package bbva

import (
	"regexp"
	"strings"
)

func truncateTo(input string, c int) string {
	if len(input) < c {
		return input
	}

	return input[:c]
}

func stripSpaces(input string) string {
	// Remove all spaces
	sp, err := regexp.Compile("[\\s]+")
	if err != nil {
		return input
	}

	return sp.ReplaceAllString(strings.TrimSpace(input), "")
}

func stripExtraSpaces(input string) string {
	// Shrink multiple spaces to one
	sp, err := regexp.Compile("[\\s]+")
	if err != nil {
		return input
	}

	return sp.ReplaceAllString(strings.TrimSpace(input), " ")
}

func stripBusinessName(input string) string {
	re, err := regexp.Compile("[^a-zA-Z0-9\\.\\'\\-\\&\\,\\s]+")
	if err != nil {
		return strings.ToUpper(input)
	}

	return strings.ToUpper(stripExtraSpaces(re.ReplaceAllString(input, "")))
}

func stripConsumerName(input string) string {
	re, err := regexp.Compile("[^a-zA-Z\\'\\-\\s]+")
	if err != nil {
		return strings.ToUpper(input)
	}

	return strings.ToUpper(stripExtraSpaces(re.ReplaceAllString(input, "")))
}

func stripCardName(input string) string {
	re, err := regexp.Compile("[^a-zA-Z0-9\\-\\s]+")
	if err != nil {
		return strings.ToUpper(input)
	}

	return strings.ToUpper(stripExtraSpaces(re.ReplaceAllString(input, "")))
}

func stripAccountName(input string) string {
	re, err := regexp.Compile("[^a-zA-Z0-9\\-\\s]+")
	if err != nil {
		return strings.ToUpper(input)
	}

	return strings.ToUpper(stripExtraSpaces(re.ReplaceAllString(input, "")))
}

func stripAddressPart(input string) string {
	re, err := regexp.Compile("[^a-zA-Z0-9\\'\\-\\s]+")
	if err != nil {
		return strings.ToUpper(input)
	}

	return strings.ToUpper(stripExtraSpaces(re.ReplaceAllString(input, "")))
}

func stripTaxID(input string) string {
	re, err := regexp.Compile("[^0-9\\s]+")
	if err != nil {
		return input
	}

	return stripExtraSpaces(re.ReplaceAllString(input, ""))
}
