package business

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/notification/activity"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
)

type BusinessActivity interface {
	GetByID(string) (*activity.Activity, error)
	List(int, int, shared.BusinessID, shared.UserID) (*[]activity.Activity, error)
}

type store struct {
	request services.SourceRequest
	*sqlx.DB
}

//New Returns a new Business service
func New(request services.SourceRequest) BusinessActivity {
	return &store{request, data.DBWrite}
}

func (s *store) GetByID(id string) (*activity.Activity, error) {
	var a activity.Activity
	err := s.Get(&a, "SELECT * FROM business_activity WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return nil, services.ErrorNotFound{}.New(fmt.Sprintf("activity with id:%s not found", id))
	} else if err != nil {
		return nil, err
	}

	return activity.ComposeActivity(a, s.request.AcceptLang)
}

func (s *store) List(offset int, limit int, businessID shared.BusinessID, userID shared.UserID) (*[]activity.Activity, error) {

	list := []activity.Activity{}

	err := s.Select(&list, `SELECT * FROM business_activity WHERE entity_id = $1 
	UNION SELECT * FROM user_activity WHERE entity_id = $2 ORDER BY activity_date DESC LIMIT $3 OFFSET $4`, businessID, userID, limit, offset)
	if err == sql.ErrNoRows {
		return &list, nil
	} else if err != nil {
		return nil, err
	}

	return activity.ComposeActivityList(list, s.request.AcceptLang)
}
