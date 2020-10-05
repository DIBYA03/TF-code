/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package analytics

//Response  ..
type Response struct {
	Accounts                   int64    `json:"accounts" db:"accounts"`
	Cards                      int64    `json:"cards" db:"cards"`
	Contacts                   int64    `json:"contacts" db:"contacts"`
	MoneySent                  *float64 `json:"moneySent" db:"money_sent"`
	Deposits                   *float64 `json:"deposits" db:"deposits"`
	MoneyRequested             *float64 `json:"moneyRequested" db:"money_requested"`
	BusinessesInMemberReview   int64    `json:"businessesInMemberReview" db:"businesses_in_member_review"`
	BusinessesInDocumentReview int64    `json:"businessesInDocumentReview" db:"businesses_in_doc_review"`
	BusinessesInBankReview     int64    `json:"businessesInBankReview" db:"businesses_in_bank_review"`
}
