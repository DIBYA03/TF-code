package business

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/csp/cspuser"
	"github.com/wiseco/core-platform/services/csp/data"
	cspsrv "github.com/wiseco/core-platform/services/csp/services"
)

// NotesService ..
type NotesService interface {
	Create(NotesCreate) (Notes, error)
	Update(string, NotesUpdate) (Notes, error)
	List(businessID string, limit, offset int) ([]Notes, error)
	ByID(businessID string, ID string) (Notes, error)
	Delete(businessID string, ID string) error
}

type notesService struct {
	rdb *sqlx.DB
	wdb *sqlx.DB
	sr  cspsrv.SourceRequest
}

func (service notesService) Create(create NotesCreate) (Notes, error) {
	var note Notes
	userID, err := cspuser.NewUserService(service.sr).ByCognitoID(service.sr.CognitoID)
	if err != nil {
		return note, err
	}
	if create.Notes == "" {
		return note, errors.New("cant create empty notes")
	}
	create.UserID = userID
	keys := services.SQLGenInsertKeys(create)
	values := services.SQLGenInsertValues(create)
	q := fmt.Sprintf("INSERT INTO business_notes (%s) VALUES(%s) RETURNING *", keys, values)
	stmt, err := service.wdb.PrepareNamed(q)
	if err != nil {
		return note, err
	}
	err = stmt.Get(&note, create)
	return note, err
}

// NewNotesService ..
func NewNotesService(sr cspsrv.SourceRequest) NotesService {
	return notesService{wdb: data.DBWrite, rdb: data.DBRead, sr: sr}
}

func (service notesService) Update(ID string, updates NotesUpdate) (Notes, error) {
	var note Notes

	keys := services.SQLGenForUpdate(updates)
	q := fmt.Sprintf("UPDATE business_notes SET %s WHERE id = '%s' RETURNING *", keys, ID)
	stmt, err := service.wdb.PrepareNamed(q)
	if err != nil {
		log.Printf("error preparing stmt %v query %v", err, q)
		return note, err
	}
	err = stmt.Get(&note, updates)
	if err != nil {
		return note, fmt.Errorf("error keys: %v err: %v", keys, err)
	}
	return note, nil
}

func (service notesService) List(businessID string, limit, offset int) ([]Notes, error) {
	list := make([]Notes, 0)

	err := service.rdb.Select(&list, `
			SELECT b.*,u.first_name,u.last_name,u.picture FROM business_notes AS b
			JOIN csp_user AS u ON b.user_id = u.id
			WHERE b.business_id = $3
			ORDER BY b.created DESC LIMIT $1 OFFSET $2 `, limit, offset, businessID)

	if err == sql.ErrNoRows {
		return list, nil
	}

	return list, err
}

func (service notesService) ByID(businessID, ID string) (Notes, error) {
	var note Notes
	err := service.rdb.Get(&note, "SELECT * FROM business_notes WHERE id = $1 AND business_id = $2", ID, businessID)
	if err != nil && err == sql.ErrNoRows {
		return note, services.ErrorNotFound{}.New("")
	}
	return note, err
}

func (service notesService) Delete(businessID, noteID string) error {
	_, err := service.wdb.Exec("DELETE FROM business_notes WHERE business_id = $1 AND id = $2", businessID, noteID)
	return err
}
