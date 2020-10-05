package auth

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/services/csp/data"
	"github.com/wiseco/core-platform/services/csp/services"
)

// ServiceReq for source request object
type ServiceReq struct {
	srcReq services.SourceRequest
	*sqlx.DB
}

// Service is the interface for csp authorization
type Service interface {
	CheckUserAccess(string) error
}

// NewService creates a new service request
func NewService(r services.SourceRequest) Service {
	return &ServiceReq{r, data.DBWrite}
}

// CheckUserAccess verifies access to users in csp_user
func (a *ServiceReq) CheckUserAccess(userID string) error {
	var id string
	err := a.Get(&id, "SELECT id FROM csp_user WHERE id = $1 AND active = true", userID)
	if err != nil {
		return err
	}

	if id != a.srcReq.UserId {
		return errors.New("unauthorized")
	}

	return nil
}
