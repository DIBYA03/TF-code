package user

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/notification/activity"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/data"
)

type UserActivity interface {
	GetByID(string) (*activity.Activity, error)
	List(string) (*[]activity.Activity, error)
}

type store struct {
	request services.SourceRequest
	*sqlx.DB
}

//New Returns a new UserActivity service
func New(request services.SourceRequest) UserActivity {
	return &store{request, data.DBWrite}
}

func (s *store) GetByID(id string) (*activity.Activity, error) {
	var activity activity.Activity
	err := s.Get(&activity, "SELECT * FROM user_activity WHERE id = $1", id)
	if err != nil && err == sql.ErrNoRows {
		return nil, services.ErrorNotFound{}.New(fmt.Sprintf("activity with id:%s not found", id))
	}

	if err != nil {
		return nil, err
	}
	//TODO:
	//We have an activity lets construct the message
	return &activity, nil
}

func (s *store) List(userID string) (*[]activity.Activity, error) {
	var list []activity.Activity
	err := s.Select(&list, "SELECT * FROM user_activity WHERE entity_id = $1", userID)
	if err != nil && err == sql.ErrNoRows {
		return &list, nil
	}

	if err != nil {
		return nil, err
	}

	//TODO:
	//We have an activity list lets construct the message

	return &list, nil
}
