package transaction

import (
	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/core-platform/transaction"
	data "github.com/wiseco/core-platform/transaction"
)

type DeclinedTransactionService interface {
	// Fetch business transactions
	ListAllDeclinedTransaction(params map[string]interface{}) ([]transaction.BusinessPendingTransaction, error)
	GetDeclinedTransactionByID(id shared.PendingTransactionID, businessID shared.BusinessID) (*transaction.BusinessPendingTransaction, error)

	ExportDeclinedTransaction(businessID *shared.BusinessID, startDate, endDate string) (*transaction.CSVTransaction, error)
}

type declinedTransactionservice struct {
	*sqlx.DB
	sourceReq services.SourceRequest
}

func NewDeclinedTransactionService(r services.SourceRequest) DeclinedTransactionService {
	return declinedTransactionservice{data.DBWrite, r}
}

func (s declinedTransactionservice) ListAllDeclinedTransaction(params map[string]interface{}) ([]transaction.BusinessPendingTransaction, error) {
	var busID shared.BusinessID
	val, ok := params["businessId"].(shared.BusinessID)
	if ok {
		busID = val
	}

	return transaction.NewDeclinedTransactionService().ListAllInternal(params, busID)
}

func (s declinedTransactionservice) GetDeclinedTransactionByID(ID shared.PendingTransactionID, businessID shared.BusinessID) (*transaction.BusinessPendingTransaction, error) {
	return transaction.NewDeclinedTransactionService().GetByIDInternal(ID, businessID)
}

func (s declinedTransactionservice) ExportDeclinedTransaction(businessID *shared.BusinessID, startDate, endDate string) (*transaction.CSVTransaction, error) {
	return transaction.NewDeclinedTransactionService().ExportInternal(businessID, startDate, endDate)
}
