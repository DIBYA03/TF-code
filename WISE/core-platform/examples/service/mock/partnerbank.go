package mock

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"

	partnerbank "github.com/wiseco/core-platform/partner/bank"
)

func NewPartnerBankConsumerEntity(userId string) partnerbank.PartnerBankEntity {
	var now = time.Now()

	b, _ := json.Marshal(map[string]interface{}{"custom": "field"})

	return partnerbank.PartnerBankEntity{
		ID:           partnerbank.PartnerEntityID(uuid.New().String()),
		BankName:     partnerbank.ProviderNameBBVA,
		EntityType:   partnerbank.EntityTypeConsumer,
		EntityID:     partnerbank.EntityID(userId),
		BankEntityID: partnerbank.BankEntityID(uuid.New().String()),
		BankExtra:    types.JSONText(string(b)),
		Created:      now,
		Updated:      now,
	}
}

func NewPartnerBankBusinessEntity(businessId string) partnerbank.PartnerBankEntity {
	var now = time.Now()

	b, _ := json.Marshal(map[string]interface{}{"custom": "field"})

	return partnerbank.PartnerBankEntity{
		ID:           partnerbank.PartnerEntityID(uuid.New().String()),
		BankName:     partnerbank.ProviderNameBBVA,
		EntityType:   partnerbank.EntityTypeBusiness,
		EntityID:     partnerbank.EntityID(businessId),
		BankEntityID: partnerbank.BankEntityID(uuid.New().String()),
		BankExtra:    types.JSONText(string(b)),
		Created:      now,
		Updated:      now,
	}
}

/*
func NewPartnerBankAccount(accountId string) PartnerBankAccount {
	var now = time.Now()

	return PartnerBankAccount{
		ID:               uuid.New().String(),
		BankName:         "bbva",
		AccountType:      PartnerBankAccountTypeChecking,
		IsPrimary:        true,
		PartnerAccountID: "AC-666265c4-c11d-4369-95cc-a660323a588c",
		AccountID:        accountId,
		AccountNumber:    "7238723451",
		RoutingNumber:    "234234226",
		AvailableBalance: 500.00,
		PostedBalance:    500.00,
		Currency:         "USD",
		Extra:            map[string]interface{}{"custom": "field"},
		Created:          now,
		Updated:          now,
	}
}

func NewPartnerBankCard() PartnerBankCard {
	var now = time.Now()

	return PartnerBankCard{
		ID:            uuid.New().String(),
		BankName:      "bbva",
		CardType:      PartnerBankCardStatusActive,
		PartnerCardID: uuid.New().String(),
		CardID:        uuid.New().String(),
		Extra:         map[string]interface{}{"custom": "field"},
		Created:       now,
		Updated:       now,
	}
}

func NewPartnerBankMoneyTransfer(transferId string) PartnerBankMoneyTransfer {
	var now = time.Now()

	return PartnerBankMoneyTransfer{
		ID:                uuid.New().String(),
		BankName:          "bbva",
		TransferID:        transferId,
		PartnerTransferID: uuid.New().String(),
		Extra:             map[string]interface{}{"custom": "field"},
		Created:           now,
		Updated:           now,
	}
}

func NewPartnerBankCardApprovedTransaction(transactionId string) PartnerBankTransaction {
	return PartnerBankTransaction{
		ID:                   uuid.New().String(),
		BankName:             "bbva",
		TransactionID:        transactionId,
		PartnerTransactionID: uuid.New().String(),
		RawTransaction: map[string]interface{}{
			"notification_id":        "NO-d6bd5eee-446b-48cc-a9d9-32dd110f74ee",
			"notification_type":      "transactions",
			"notification_version":   "1.0.0",
			"notification_reason":    "authorization_approved",
			"notification_sent_date": "2018-06-28T19:26:10Z",
			"notification_data": map[string]interface{}{
				"transaction_id": "6750209325000000066",
				"user_id":        "CO-ff5df61b-cb4b-4ab0-9edf-97a570431b57",
				"account_id":     "AC-6b6c980b-8da4-4104-b9bf-2f9eb141b669",
				"card_id":        "DC-a3dba49c-19d3-4763-a94f-c0b27e646d5c",
				"hold_data": map[string]interface{}{
					"number":      66,
					"amount":      10.94,
					"date":        "2018-06-28T00:00:00",
					"expiry_date": "2018-07-02T00:00:00",
				},
				"card_details": map[string]interface{}{
					"visa_transaction_id":         "468179699465984",
					"authorization_amount":        10.94,
					"authorization_date":          "2018-06-28T19:25:46Z",
					"authorization_response":      "00",
					"authorization_number":        "281276",
					"card_transaction_type":       "002000",
					"local_transaction_amount":    10.94,
					"local_transaction_currency":  "USD",
					"local_transaction_date_time": "2018-06-28T13:25:46",
					"pos_entry_mode":              "050",
					"pos_condition_code":          "01000008045",
					"merchant_category_code":      "5812",
					"acquirer_bin":                "476501",
					"card_acceptor_id":            "000266015396887",
					"card_acceptor_terminal":      "08673556",
					"card_acceptor_address":       "OSHA EXPRESS",
					"card_acceptor_city":          "SAN FRANCISCO",
					"card_acceptor_state":         "CA",
					"card_acceptor_country":       "US",
				},
			},
		},
		Created: time.Now(),
	}
}

func NewPartnerBankCardDeclinedTransaction(transactionId string) PartnerBankTransaction {
	return PartnerBankTransaction{
		ID:                   uuid.New().String(),
		BankName:             "bbva",
		TransactionID:        transactionId,
		PartnerTransactionID: uuid.New().String(),
		RawTransaction: map[string]interface{}{
			"notification_id":        "NO-d6bd5eee-446b-48cc-a9d9-32dd110f74ee",
			"notification_type":      "transactions",
			"notification_version":   "1.0.0",
			"notification_reason":    "authorization_declined",
			"notification_sent_date": "2018-06-28T19:26:10Z",
			"notification_data": map[string]interface{}{
				"transaction_id": "6750209325000000066",
				"user_id":        "CO-ff5df61b-cb4b-4ab0-9edf-97a570431b57",
				"account_id":     "AC-6b6c980b-8da4-4104-b9bf-2f9eb141b669",
				"card_id":        "DC-a3dba49c-19d3-4763-a94f-c0b27e646d5c",
				"card_details": map[string]interface{}{
					"visa_transaction_id":         "468179699465984",
					"authorization_amount":        10.94,
					"authorization_date":          "2018-06-28T19:25:46Z",
					"authorization_response":      "00",
					"authorization_number":        "281276",
					"card_transaction_type":       "002000",
					"local_transaction_amount":    10.94,
					"local_transaction_currency":  "USD",
					"local_transaction_date_time": "2018-06-28T13:25:46",
					"pos_entry_mode":              "050",
					"pos_condition_code":          "01000008045",
					"merchant_category_code":      "5812",
					"acquirer_bin":                "476501",
					"card_acceptor_id":            "000266015396887",
					"card_acceptor_terminal":      "08673556",
					"card_acceptor_address":       "OSHA EXPRESS",
					"card_acceptor_city":          "SAN FRANCISCO",
					"card_acceptor_state":         "CA",
					"card_acceptor_country":       "US",
				},
			},
		},
		Created: time.Now(),
	}
}
*/
