/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for all internal banking related items such as accounts and cards
package banking

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type CardType string

const (
	// Debit card type
	CardTypeDebit = CardType("debit")

	// Credit card type
	CardTypeCredit = CardType("credit")

	// Prepaid debit card type
	CardTypePrepaid = CardType("prepaid")

	// Single use debit card type
	CardTypeSingleUse = CardType("singleUse")
)

type CardStatus string

const (
	// Card is inactive
	CardStatusInactive = CardStatus("inactive")

	// Card is active
	CardStatusActive = CardStatus("active")

	// Card is blocked by bank
	CardStatusBlocked = CardStatus("blocked")

	// Card is unblocked by bank
	CardStatusUnblocked = CardStatus("unblocked")

	// Card is reissued by bank
	CardStatusReissued = CardStatus("reissued")

	// Card is reissued by bank
	CardStatusLimitChanged = CardStatus("limitChanged")

	// Card is locked by user
	CardStatusLocked = CardStatus("locked")

	// Card is cancelled
	CardStatusCanceled = CardStatus("canceled")

	// Card is being embossed
	CardStatusEmbossing = CardStatus("embossing")

	// Card is being delivered
	CardStatusDelivery = CardStatus("delivery")

	// Card is being delivered
	CardStatusDelivered = CardStatus("delivered")

	// Card is being shipped
	CardStatusShipped = CardStatus("shipped")
)

type CardBlockID string

const (
	// Card lost
	CardBlockIDLost = CardBlockID("lost")

	// Card stolen
	CardBlockIDStolen = CardBlockID("stolen")

	// Card blocked due to dispute
	CardBlockIDDispute = CardBlockID("dispute")

	// Card blocked due to fraud
	CardBlockIDFraud = CardBlockID("fraud")

	// Card blocked internally
	CardBlockIDInternal = CardBlockID("internal")

	// Card locked
	CardBlockIDLocked = CardBlockID("locked")

	// Card blocked due to charge off
	CardBlockIDChargeOff = CardBlockID("chargeOff")
)

type CardBusinessName string

const (
	// Card is inactive
	CardBusinessNameLegal = CardBusinessName("legal")

	// Card is active
	CardBusinessNameDBA = CardBusinessName("dba")
)

type CardBlockArray []BankCardBlock

// SQL value marshaller
func (b CardBlockArray) Value() (driver.Value, error) {
	return json.Marshal(b)
}

// SQL scan unmarshaller
func (b *CardBlockArray) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type convertible to []byte")
	}

	var out CardBlockArray
	err := json.Unmarshal(source, &out)
	if err != nil {
		return err
	}

	*b = out
	return nil
}

func (ba CardBlockArray) ToArray() []BankCardBlock {
	return []BankCardBlock(ba)
}
