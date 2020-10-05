package business

import (
	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
)

type accountClosureDataStore struct {
	sourceRequest services.SourceRequest
	rdb           *sqlx.DB
	wdb           *sqlx.DB
}

//AccountClosureService ...
type AccountClosureService interface {
	Create(create AccountClosureCreate) (*AccountClosureItem, error)
	GetByBusinessID(ID shared.BusinessID) (*AccountClosureItem, error)
}

//NewAccountClosureService ...
func NewAccountClosureService(r services.SourceRequest) AccountClosureService {
	return &accountClosureDataStore{r, data.DBRead, data.DBWrite}
}

func (s accountClosureDataStore) Create(create AccountClosureCreate) (*AccountClosureItem, error) {
	_, err := s.GetByBusinessID(create.BusinessID)
	if err == nil {
		return nil, services.ErrorNotFound{}.New("A request for account closure is already in progress")
	}
	_, err = s.wdb.Exec(`INSERT INTO account_closure_request(business_id, reason, description) 
	VALUES($1, $2, $3)`, create.BusinessID, create.Reason, create.Description)

	req, err := s.GetByBusinessID(create.BusinessID)
	return req, err
}

func (s accountClosureDataStore) GetByBusinessID(ID shared.BusinessID) (*AccountClosureItem, error) {
	var item = AccountClosureItem{}
	err := s.rdb.Get(&item, "SELECT * FROM account_closure_request WHERE business_id = $1 AND status IN ('pending', 'approved')", ID)
	return &item, err
}
