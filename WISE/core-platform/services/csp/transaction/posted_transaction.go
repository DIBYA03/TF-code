package transaction

import (
	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/core-platform/transaction"
	data "github.com/wiseco/core-platform/transaction"
)

type PostedTransactionService interface {
	// Fetch business transactions
	ListAllPostedTransactions(params map[string]interface{}) ([]transaction.BusinessPostedTransaction, error)
	GetPostedTransactionByID(id shared.PostedTransactionID, businessID shared.BusinessID) (*transaction.BusinessPostedTransaction, error)

	ExportPostedTransaction(params map[string]interface{}) (*transaction.CSVTransaction, error)
}

type postedTransactionservice struct {
	*sqlx.DB
	sourceReq services.SourceRequest
}

func NewPostedTransactionService(r services.SourceRequest) PostedTransactionService {
	return postedTransactionservice{data.DBWrite, r}
}

func (s postedTransactionservice) ListAllPostedTransactions(params map[string]interface{}) ([]transaction.BusinessPostedTransaction, error) {
	var busID shared.BusinessID
	val, ok := params["businessId"].(shared.BusinessID)
	if ok {
		busID = val
	}

	return transaction.NewBusinessService().ListAllInternal(params, busID, "")
}

func (s postedTransactionservice) GetPostedTransactionByID(ID shared.PostedTransactionID, businessID shared.BusinessID) (*transaction.BusinessPostedTransaction, error) {
	return transaction.NewBusinessService().GetByIDInternal(ID)
}

func (s postedTransactionservice) ExportPostedTransaction(params map[string]interface{}) (*transaction.CSVTransaction, error) {
	return transaction.NewBusinessService().ExportInternal(params)
}
