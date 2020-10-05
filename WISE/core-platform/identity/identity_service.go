/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/
package identity

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/shared"
)

type identityService struct {
	sourceReq SourceRequest
	wdb       *sqlx.DB
	rdb       *sqlx.DB
}

type IdentityService interface {
	// Fetch operations
	GetByID(shared.IdentityID) (*Identity, error)

	// Get by Provider
	GetByProviderID(ProviderID, ProviderName, ProviderSource) (*Identity, error)

	// Create c
	Create(IdentityCreate) (*shared.IdentityID, error)

	// Deactivate c by id
	Deactivate(shared.IdentityID) error

	// Delete
	Delete(shared.IdentityID) error

	// Update
	Update(IdentityUpdate) error
}

func NewIdentityService(r SourceRequest) IdentityService {
	return &identityService{r, DBWrite, DBRead}
}

// NewIdentityServiceWithout returns an c service without a source request
func NewIdentityServiceWithout() IdentityService {
	return &identityService{SourceRequest{}, DBWrite, DBRead}
}

func (s *identityService) GetByID(id shared.IdentityID) (*Identity, error) {
	return s.getByID(id)
}

func (s *identityService) getByID(id shared.IdentityID) (*Identity, error) {
	i := Identity{}

	err := s.wdb.Get(&i, "SELECT * FROM identity WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	return &i, nil
}

func (s *identityService) GetByProviderID(id ProviderID, name ProviderName, source ProviderSource) (*Identity, error) {
	i := Identity{}

	err := s.wdb.Get(&i, "SELECT * FROM identity WHERE provider_id = $1 AND provider_name = $2 AND provider_source = $3", id, name, source)
	if err != nil {
		return nil, err
	}

	return &i, nil
}

//CreateFromAuth creates a minimal c from cognito
func (s *identityService) Create(u IdentityCreate) (*shared.IdentityID, error) {
	// Insert statement
	sql := `
		INSERT INTO identity (provider_id, provider_name, provider_source, phone)
		VALUES (:provider_id, :provider_name, :provider_source, :phone)
		RETURNING id`

	// Execute
	stmt, err := s.wdb.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	// Return id
	var id shared.IdentityID
	err = stmt.Get(&id, &u)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (s *identityService) Deactivate(id shared.IdentityID) error {
	_, err := s.wdb.Exec("UPDATE identity SET deactivated = CURRENT_TIMESTAMP WHERE id = $1", id)
	return err
}

func (s *identityService) Delete(id shared.IdentityID) error {
	identity, err := s.getByID(id)
	if err != nil {
		return err
	}

	_, err = s.wdb.Exec("DELETE FROM identity WHERE id = $1", id)
	if err != nil {
		return err
	}

	// Remove from Cognito
	if identity.ProviderName == ProviderNameCognito {
		sess := session.Must(session.NewSession(&aws.Config{}))
		prov := cognitoidentityprovider.New(sess)

		del := cognitoidentityprovider.AdminDeleteUserInput{
			UserPoolId: aws.String(string(identity.ProviderSource)),
			Username:   aws.String(identity.Phone),
		}

		_, err = prov.AdminDeleteUser(&del)
	}

	return err
}

func (s *identityService) Update(u IdentityUpdate) error {
	_, err := s.wdb.Exec("UPDATE identity SET phone = $1 WHERE id = $2", u.Phone, u.ID)
	return err
}
