/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for transaction services
package transaction

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/auth"
	"github.com/wiseco/core-platform/notification/activity"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking"
	"github.com/wiseco/core-platform/services/banking/business"
	"github.com/wiseco/core-platform/shared"
)

type disputeDatastore struct {
	sourceReq services.SourceRequest
	*sqlx.DB
}

type DisputeService interface {
	// Read
	GetById(string, shared.PostedTransactionID, shared.BusinessID) (*Dispute, error)

	// Create dispute
	Create(*DisputeCreate) (*Dispute, error)

	// Cancel dispute
	Cancel(*DisputeCancel) (*Dispute, error)
}

func NewDisputeService(r services.SourceRequest) DisputeService {
	return &disputeDatastore{r, DBWrite}
}

func (db *disputeDatastore) GetById(ID string, txnID shared.PostedTransactionID, bID shared.BusinessID) (*Dispute, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(bID)
	if err != nil {
		return nil, err
	}

	dispute := Dispute{}

	err = db.Get(&dispute, "SELECT * FROM business_transaction_dispute WHERE id = $1 AND transaction_id = $2 AND business_id = $3", ID, txnID, bID)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return nil, err
	}

	return &dispute, err
}

func (db *disputeDatastore) Create(d *DisputeCreate) (*Dispute, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(d.BusinessID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Check transaction status
	t, err := NewBusinessService().GetByID(d.TransactionID, db.sourceReq.UserID, d.BusinessID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if t.Dispute.DisputeStatus != nil {
		return nil, errors.New("Transactions once disputed cannot be disputed again")
	}

	if d.Category == "" {
		return nil, errors.New("Dispute category is required")
	}

	_, ok := categoryTo[d.Category]
	if !ok {
		return nil, errors.New("Dispute category is invalid")
	}

	// Receipt id is not required for fraudulent transactions
	if d.ReceiptID == nil && d.Category != CategoryFraudulentCharge {
		return nil, errors.New("Receipt id is required")
	}

	if d.Summary == nil && d.Category != CategoryFraudulentCharge {
		return nil, errors.New("Dispute description is required")
	}

	d.DisputeNumber = shared.GetRandomString("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 8)

	// Add status
	status := DisputeStatusDisputed
	d.DisputeStatus = &status

	// Default/mandatory fields
	columns := []string{
		"dispute_number", "transaction_id", "receipt_id", "created_user_id", "business_id", "category", "summary", "dispute_status",
	}
	// Default/mandatory values
	values := []string{
		":dispute_number", ":transaction_id", ":receipt_id", ":created_user_id", ":business_id", ":category", ":summary", ":dispute_status",
	}

	sql := fmt.Sprintf("INSERT INTO business_transaction_dispute(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))

	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	dispute := &Dispute{}

	err = stmt.Get(dispute, &d)
	if err != nil {
		return nil, err
	}

	// Block card in case of fradulent charge
	if d.Category == CategoryFraudulentCharge {
		cards, err := business.NewCardService(db.sourceReq).GetByBusinessID(0, 20, d.BusinessID, d.CreatedUserID)
		if err != nil {
			log.Println("error fetching user card", err)
			return nil, err
		}

		id := getActiveCardId(cards)

		if id != "" {
			block := business.BankCardBlockCreate{
				BusinessID: d.BusinessID,
			}
			block.CardholderID = d.CreatedUserID
			block.BlockID = banking.CardBlockIDDispute
			block.CardID = id
			_, err = business.NewCardService(db.sourceReq).BlockBankCard(&block)
			if err != nil {
				log.Println("Error blocking card", err)
				return nil, err
			}
		}
	}

	// Add dispute to activity stream
	amt := t.Amount.FormatCurrency()
	db.onTxnDisputed(*dispute, amt)

	return dispute, nil
}

func getActiveCardId(cards []business.BankCard) string {

	for _, c := range cards {
		if c.CardStatus == "active" {
			return c.Id
		}
	}

	return ""

}

func (db *disputeDatastore) onTxnDisputed(d Dispute, amt string) error {

	desc, _ := categoryTo[d.Category]

	dispute := activity.Dispute{
		EntityID:      string(d.BusinessID),
		UserID:        d.CreatedUserID,
		TransactionID: string(d.TransactionID),
		Amount:        amt,
		Category:      desc,
	}

	err := activity.NewDisputeCreator().Create(dispute)
	if err != nil {
		log.Println(err)
	}

	return err
}

func (db *disputeDatastore) Cancel(d *DisputeCancel) (*Dispute, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(d.BusinessID)
	if err != nil {
		return nil, err
	}

	// Add cancel dispute to activity stream
	t, err := NewBusinessService().GetByID(d.TransactionID, db.sourceReq.UserID, d.BusinessID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if t.Dispute.DisputeStatus == nil {
		return nil, errors.New("Only disputed transactions can be cancelled")
	}

	if *t.Dispute.DisputeStatus != DisputeStatusDisputed {
		return nil, errors.New("Only transactions with disputed status can be cancelled")
	}

	_, err = db.Exec(("UPDATE business_transaction_dispute SET dispute_status = '" + string(DisputeStatusDisputedCancelled) + "' WHERE id = $1"), d.Id)

	if err != nil {
		log.Println("error updating dispute status", err)
		return nil, err
	}

	// Add dispute to activity stream
	amt := t.Amount.FormatCurrency()
	db.onDisputeCancelled(d, amt)

	return db.GetById(d.Id, d.TransactionID, d.BusinessID)

}

func (db *disputeDatastore) onDisputeCancelled(d *DisputeCancel, amt string) error {

	dispute := activity.Dispute{
		EntityID:      string(d.BusinessID),
		UserID:        d.CreatedUserID,
		TransactionID: string(d.TransactionID),
		Amount:        amt,
		Category:      "",
	}

	err := activity.NewDisputeCreator().Delete(dispute)
	if err != nil {
		log.Println(err)
	}

	return err
}
