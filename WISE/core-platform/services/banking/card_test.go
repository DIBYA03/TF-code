package banking_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wiseco/core-platform/services/banking"
)

func TestCardGetCardNumberLastFour(t *testing.T) {
	assert := require.New(t)

	testCases := []struct {
		maskedCardNumber string
		expected         string
		err              error
	}{
		{"", "", errors.New("Unable to get last four digits of card, CardNumberMasked not long enough. length: 0")},
		{"X", "", errors.New("Unable to get last four digits of card, CardNumberMasked not long enough. length: 1")},
		{"XXXXXXXXXXX4568", "4568", nil},
	}

	for _, tt := range testCases {
		bc := &banking.BankCard{CardNumberMasked: tt.maskedCardNumber}

		actual, err := bc.GetCardNumberLastFour()

		assert.Equal(tt.expected, actual)
		assert.Equal(tt.err, err)
	}
}
