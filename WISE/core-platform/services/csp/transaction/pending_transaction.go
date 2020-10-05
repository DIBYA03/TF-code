package transaction

import (
	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/core-platform/transaction"
	data "github.com/wiseco/core-platform/transaction"
)

type PendingTransactionService interface {
	// Fetch business transactions
	ListAllPendingTransaction(params map[string]interface{}) ([]transaction.BusinessPendingTransaction, error)
	GetPendingTransactionByID(id shared.PendingTransactionID, businessID shared.BusinessID) (*transaction.BusinessPendingTransaction, error)

	ExportPendingTransaction(params map[string]interface{}) (*transaction.CSVTransaction, error)
}

type pendingTransactionservice struct {
	*sqlx.DB
	sourceReq services.SourceRequest
}

func NewPendingTransactionService(r services.SourceRequest) PendingTransactionService {
	return pendingTransactionservice{data.DBWrite, r}
}

func (s pendingTransactionservice) ListAllPendingTransaction(params map[string]interface{}) ([]transaction.BusinessPendingTransaction, error) {
	var busID shared.BusinessID
	val, ok := params["businessId"].(shared.BusinessID)
	if ok {
		busID = val
	}

	return transaction.NewPendingTransactionService().ListAllInternal(params, busID)
}

func (s pendingTransactionservice) GetPendingTransactionByID(ID shared.PendingTransactionID, businessID shared.BusinessID) (*transaction.BusinessPendingTransaction, error) {
	return transaction.NewPendingTransactionService().GetByIDInternal(ID, businessID)
}

func (s pendingTransactionservice) ExportPendingTransaction(params map[string]interface{}) (*transaction.CSVTransaction, error) {
	return transaction.NewPendingTransactionService().ExportInternal(params)
}
