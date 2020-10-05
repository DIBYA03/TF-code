package business

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
)

type cardReaderService struct {
	rdb       *sqlx.DB
	wdb       *sqlx.DB
	sourceReq services.SourceRequest
}

// CardReaderService ..
type CardReaderService interface {
	List(businessID shared.BusinessID) ([]CardReader, error)
	ByID(ID shared.CardReaderID, businessID shared.BusinessID) (*CardReader, error)
	Create(CardReaderCreate) (*CardReader, error)
	Update(id shared.CardReaderID, updates CardReaderUpdate) (*CardReader, error)
	Deactivate(id shared.CardReaderID, businessID shared.BusinessID) error
}

// NewCardReaderService ..
func NewCardReaderService(r services.SourceRequest) CardReaderService {
	return cardReaderService{wdb: data.DBWrite, rdb: data.DBRead, sourceReq: r}
}

func (service cardReaderService) List(businessID shared.BusinessID) ([]CardReader, error) {
	list := make([]CardReader, 0)
	err := service.rdb.Select(&list, "SELECT * FROM card_reader WHERE deactivated IS NULL AND business_id = $1", businessID)
	if err != nil && err == sql.ErrNoRows {
		return list, nil
	}
	return list, err
}

func (service cardReaderService) ByID(id shared.CardReaderID, businessID shared.BusinessID) (*CardReader, error) {

	cardReader := CardReader{}
	err := service.rdb.Get(&cardReader, "SELECT * FROM card_reader WHERE id = $1 AND business_id = $2", id, businessID)
	if err != nil && err == sql.ErrNoRows {
		return nil, services.ErrorNotFound{}.New("")
	}

	return &cardReader, err
}

func (service cardReaderService) Create(c CardReaderCreate) (*CardReader, error) {
	if err := validateCardReader(c); err != nil {
		return nil, err
	}
	// Default/mandatory fields
	columns := []string{
		"business_id", "alias", "device_type", "serial_number",
	}
	// Default/mandatory values
	values := []string{
		":business_id", ":alias", ":device_type", ":serial_number",
	}

	sql := fmt.Sprintf("INSERT INTO card_reader(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))

	stmt, err := service.wdb.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	cardReader := &CardReader{}

	err = stmt.Get(cardReader, &c)
	if err != nil {
		return nil, err
	}

	return cardReader, nil
}

func validateCardReader(c CardReaderCreate) error {
	if c.SerialNumber == "" {
		return errors.New("Serial number is required")
	}
	return nil
}

func (service cardReaderService) Update(id shared.CardReaderID, updates CardReaderUpdate) (*CardReader, error) {
	var reader CardReader

	keys := services.SQLGenForUpdate(updates)
	q := fmt.Sprintf("UPDATE card_reader SET %s WHERE id = '%s' RETURNING *", keys, id)
	stmt, err := service.wdb.PrepareNamed(q)
	if err != nil {
		log.Printf("error preparing stmt %v query %v", err, q)
		return &reader, err
	}
	err = stmt.Get(&reader, updates)
	if err != nil {
		return &reader, fmt.Errorf("error keys: %v err: %v", keys, err)
	}
	return &reader, nil
}

func (service cardReaderService) Deactivate(id shared.CardReaderID, businessID shared.BusinessID) error {
	_, err := service.wdb.Exec("UPDATE card_reader SET deactivated = CURRENT_TIMESTAMP WHERE id = $1 AND business_id = $2", id, businessID)
	return err
}
