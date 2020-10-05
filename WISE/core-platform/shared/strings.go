package shared

import (
	"regexp"
	"strings"
)

func StripSpaces(input string) string {
	// Remove all spaces
	sp, err := regexp.Compile("[\\s]+")
	if err != nil {
		return input
	}

	return sp.ReplaceAllString(strings.TrimSpace(input), "")
}

func StripExtraSpaces(input string) string {
	// Shrink multiple spaces to one
	sp, err := regexp.Compile("[\\s]+")
	if err != nil {
		return input
	}

	return sp.ReplaceAllString(strings.TrimSpace(input), " ")
}

func StripNonDigits(input string) string {
	re, err := regexp.Compile("[^0-9\\s]+")
	if err != nil {
		return input
	}

	return StripExtraSpaces(re.ReplaceAllString(input, ""))
}
