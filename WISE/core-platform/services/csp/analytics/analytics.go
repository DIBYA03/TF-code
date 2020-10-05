/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package analytics

import (
	"github.com/wiseco/core-platform/services/invoice"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	busBanking "github.com/wiseco/core-platform/services/banking/business"
	cspData "github.com/wiseco/core-platform/services/csp/data"
	"github.com/wiseco/core-platform/services/data"
)

// Service .
type Service interface {
	Metrics() (Response, error)
	TopMetrics()
}

type service struct {
	coreDB *sqlx.DB
	cspDB  *sqlx.DB
}

//NewAnlytics ..
func NewAnlytics() Service {
	return service{data.DBWrite, cspData.DBWrite}
}

func (s service) Metrics() (Response, error) {
	var ana Response
	clearingID := os.Getenv("WISE_CLEARING_BUSINESS_ID")

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bas, err := busBanking.NewBankingAccountService()
		if err != nil {
			return ana, err
		}

		ana.Accounts, err = bas.FindCount()
		if err != nil {
			return ana, err
		}
	} else {
		err := s.coreDB.Get(&ana, `SELECT (SELECT COUNT(*) FROM business_bank_account) AS accounts`)
		if err != nil {
			return ana, err
		}
	}

	if os.Getenv("USE_INVOICE_SERVICE") == "true" {
		err := s.coreDB.Get(&ana, `SELECT (SELECT COUNT(*) FROM business_bank_card) AS cards,
	(SELECT COUNT(*) FROM business_contact) AS contacts,
	(SELECT SUM(amount) FROM business_money_request where request_type = 'pos') AS money_requested,
	(SELECT SUM(posted_balance) FROM business_bank_account) AS deposits,
	(SELECT SUM(amount)	 FROM business_money_transfer WHERE business_id != $1) AS money_sent`, clearingID)

		if err != nil {
			return ana, err
		}

		invSvc, err := invoice.NewInvoiceService()
		if err != nil {
			log.Println(err)
			return ana, err
		}
		amountResp, err := invSvc.GetInvoiceAmounts()
		if err != nil {
			log.Println(err)
			return ana, err
		}
		if amountResp.IsPositive() {
			amntFlt, ok := amountResp.Float64()
			if ok {
				totalMoneyRequested := *ana.MoneyRequested + amntFlt
				ana.MoneyRequested = &totalMoneyRequested
			}
		}

	} else {
		err := s.coreDB.Get(&ana, `SELECT (SELECT COUNT(*) FROM business_bank_card) AS cards,
	(SELECT COUNT(*) FROM business_contact) AS contacts,
	(SELECT SUM(amount) FROM business_money_request) AS money_requested,
	(SELECT SUM(posted_balance) FROM business_bank_account) AS deposits,
	(SELECT SUM(amount)	 FROM business_money_transfer WHERE business_id != $1) AS money_sent`, clearingID)

		if err != nil {
			return ana, err
		}
	}

	err := s.cspDB.Get(&ana, `SELECT (SELECT COUNT(*) FROM business WHERE review_status = 'memberReview') AS businesses_in_member_review,
	(SELECT COUNT(*) FROM business WHERE review_status = 'docReview') AS businesses_in_doc_review,
	(SELECT COUNT(*) FROM business WHERE review_status = 'bankReview') AS businesses_in_bank_review`)

	return ana, err
}

func (s service) TopMetrics() {
	//TODO
}
