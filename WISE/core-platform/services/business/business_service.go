/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for all business related services
package business

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/mail"
	"strings"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/wiseco/core-platform/auth"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	_ "github.com/wiseco/core-platform/partner/bank/bbva"
	"github.com/wiseco/core-platform/partner/service/segment"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/services/document"
	"github.com/wiseco/core-platform/services/email"
	"github.com/wiseco/core-platform/shared"
)

type businessDatastore struct {
	sourceReq services.SourceRequest
	*sqlx.DB
}

type BusinessService interface {
	// Fetch
	Count() (int, error)
	List(offset int, limit int, owner shared.UserID) ([]Business, error)
	GetById(shared.BusinessID) (*Business, error)
	GetByIdInternal(shared.BusinessID) (*Business, error)

	// Create/Modify
	Create(create *BusinessCreate) (*Business, error)
	Update(BusinessUpdate) (*Business, error)
	Deactivate(shared.BusinessID) error

	// Update subscription
	UpdateSubscription(BusinessSubscriptionUpdate) (*Business, error)

	// Submit
	Submit(shared.BusinessID) (*Business, error)
}

func NewBusinessServiceWithout() BusinessService {
	return &businessDatastore{services.NewSourceRequest(), data.DBWrite}
}

func NewBusinessService(r services.SourceRequest) BusinessService {
	return &businessDatastore{r, data.DBWrite}
}

//New returns a BusinessService without a source request
//Use this service carefully.
func New() BusinessService {
	return &businessDatastore{services.SourceRequest{}, data.DBWrite}
}

func (db *businessDatastore) Count() (int, error) {
	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM business WHERE deactivated IS NULL AND owner_id = $1", db.sourceReq.UserID).Scan(&count)

	if err != nil {
		return 0, errors.Cause(err)
	}

	return count, err
}

func (db *businessDatastore) List(offset int, limit int, ownerID shared.UserID) ([]Business, error) {
	if ownerID != db.sourceReq.UserID {
		return nil, errors.New("unauthorized")
	}

	rows := []Business{}
	err := db.Select(&rows, "SELECT * FROM business WHERE deactivated IS NULL AND owner_id = $1", ownerID)
	switch {
	case err == sql.ErrNoRows:
		return rows, nil
	case err != nil:
		return nil, errors.Cause(err)
	}

	// Set ach enabled
	for i, _ := range rows {
		rows[i].ACHPullEnabled = true
	}

	return rows, err
}

func (db *businessDatastore) Create(c *BusinessCreate) (*Business, error) {

	if c.OperationType != nil && *c.OperationType != OperationTypeLocal {
		return nil, errors.New("Only businesses with local operations are accepted at this time")
	}

	// Default/mandatory fields
	columns := []string{
		"owner_id", "legal_name", "employer_number", "dba", "entity_type", "industry_type", "tax_id",
		"tax_id_type", "origin_country", "origin_state", "origin_date", "purpose", "operation_type", "phone",
		"activity_type", "is_restricted", "website", "online_info",
	}
	// Default/mandatory values
	values := []string{
		":owner_id", ":legal_name", ":employer_number", ":dba", ":entity_type", ":industry_type",
		":tax_id", ":tax_id_type", ":origin_country", ":origin_state", ":origin_date", ":purpose", ":operation_type",
		":phone", ":activity_type", ":is_restricted", ":website", ":online_info",
	}

	// Email validation
	if c.Email != nil {
		a, err := mail.ParseAddress(*c.Email)
		if err != nil {
			return nil, errors.New("Invalid email address")
		}

		c.Email = &a.Address

		columns = append(columns, "email")
		values = append(values, ":email")

		ec := email.EmailCreate{
			EmailAddress: email.EmailAddress(a.Address),
			Status:       email.StatusActive,
			Type:         email.TypeBusiness,
		}

		e, err := email.NewEmailService(db.sourceReq).Create(&ec)

		if err != nil {
			return nil, err
		}

		c.EmailID = e.ID

		columns = append(columns, "email_id")
		values = append(values, ":email_id")
	}

	if c.Website != nil {
		columns = append(columns, "website")
		values = append(values, ":website")
	}

	if c.OnlineInfo != nil {
		columns = append(columns, "online_info")
		values = append(values, ":online_info")
	}

	if c.TaxID != nil {
		taxID, err := services.ValidateTaxID(c.TaxID, c.TaxIDType)
		if err != nil {
			return nil, err
		}

		c.TaxID = taxID
	} else {
		c.TaxIDType = nil
	}

	// Employer number
	c.EmployerNumber = generateEmployerNumber()

	if c.LegalAddress != nil {
		columns = append(columns, "legal_address")
		values = append(values, ":legal_address")
		c.LegalAddress.Type = services.AddressTypeLegal
	}

	if c.HeadquarterAddress != nil {
		columns = append(columns, "headquarter_address")
		values = append(values, ":headquarter_address")
		c.HeadquarterAddress.Type = services.AddressTypeHeadquarter
	}

	if c.MailingAddress != nil {
		columns = append(columns, "mailing_address")
		values = append(values, ":mailing_address")
		c.MailingAddress.Type = services.AddressTypeMailing
	}

	sql := fmt.Sprintf(
		"INSERT INTO business(%s) VALUES(%s) RETURNING *",
		strings.Join(columns, ", "),
		strings.Join(values, ", "),
	)

	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	b := &Business{}
	err = stmt.Get(b, &c)
	if err != nil {
		return nil, err
	}

	// Send to segment
	segment.NewSegmentService().PushToAnalyticsQueue(b.OwnerID, segment.CategoryBusiness, segment.ActionCreate, b)

	return b, nil
}

func (db *businessDatastore) GetById(id shared.BusinessID) (*Business, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(id)
	if err != nil {
		return nil, err
	}

	b := Business{}
	err = db.Get(&b, "SELECT * FROM business WHERE id = $1", id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Set ach enabled
	b.ACHPullEnabled = true
	return &b, err
}

func (db *businessDatastore) GetByIdInternal(id shared.BusinessID) (*Business, error) {
	b := Business{}
	err := db.Get(&b, "SELECT * FROM business WHERE id = $1", id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Set ach enabled
	b.ACHPullEnabled = true
	return &b, err
}

func (db *businessDatastore) Update(u BusinessUpdate) (*Business, error) {
	b, err := db.GetById(u.ID)
	if err != nil {
		return nil, err
	}

	var columns []string

	if u.ActivityType != nil {
		columns = append(columns, "activity_type = :activity_type")
	}

	if u.IndustryType != nil {
		columns = append(columns, "industry_type = :industry_type")
	}

	if u.OriginCountry != nil {
		columns = append(columns, "origin_country = :origin_country")
	}

	if u.OriginState != nil {
		columns = append(columns, "origin_state = :origin_state")
	}

	if u.OriginDate != nil {
		columns = append(columns, "origin_date = :origin_date")
	}

	if u.Purpose != nil {
		columns = append(columns, "purpose = :purpose")
	}

	if u.OperationType != nil {
		columns = append(columns, "operation_type = :operation_type")
	}

	// Email validation
	if u.Email != nil {
		a, err := mail.ParseAddress(*u.Email)
		if err != nil {
			return nil, errors.New("Invalid email address")
		}

		u.Email = &a.Address
		columns = append(columns, "email = :email")

		esrvc := email.NewEmailService(db.sourceReq)
		if b.EmailID != shared.EmailID("") {
			err = esrvc.Deactivate(b.EmailID)

			if err != nil {
				return nil, err
			}
		}

		ec := email.EmailCreate{
			EmailAddress: email.EmailAddress(a.Address),
			Status:       email.StatusActive,
			Type:         email.TypeBusiness,
		}

		e, err := esrvc.Create(&ec)

		if err != nil {
			return nil, err
		}

		u.EmailID = e.ID

		columns = append(columns, "email_id = :email_id")
	}

	if u.Phone != nil {
		columns = append(columns, "phone = :phone")
	}

	if u.IsRestricted != nil {
		columns = append(columns, "is_restricted = :is_restricted")
	}

	if u.Website != nil {
		columns = append(columns, "website = :website")
	}

	if u.OnlineInfo != nil {
		columns = append(columns, "online_info = :online_info")
	}

	if b.KYCStatus != services.KYCStatusApproved {
		if u.LegalName != nil {
			columns = append(columns, "legal_name = :legal_name")
		}

		if u.DBA != nil {
			columns = append(columns, "dba = :dba")
		}

		if u.EntityType != nil {
			columns = append(columns, "entity_type = :entity_type")
		}

		if u.TaxID != nil {
			if u.TaxIDType != nil {
				columns = append(columns, "tax_id = :tax_id")
				columns = append(columns, "tax_id_type = :tax_id_type")
			}
		}

		if u.LegalAddress != nil {
			columns = append(columns, "legal_address = :legal_address")
			u.LegalAddress.Type = services.AddressTypeLegal
		}

		if u.HeadquarterAddress != nil {
			columns = append(columns, "headquarter_address = :headquarter_address")
			u.HeadquarterAddress.Type = services.AddressTypeHeadquarter
		}

		if u.MailingAddress != nil {
			columns = append(columns, "mailing_address = :mailing_address")
			u.MailingAddress.Type = services.AddressTypeMailing
		}
	}

	// Update if changes requested
	if len(columns) > 0 {
		// Update business
		_, err = db.NamedExec(
			fmt.Sprintf(
				"UPDATE business SET %s WHERE id = '%s'",
				strings.Join(columns, ", "),
				u.ID,
			), u,
		)
		if err != nil {
			return nil, errors.Cause(err)
		}
	}

	bus, err := db.GetById(u.ID)
	if err != nil {
		return nil, err
	}

	// Update bank data
	if b.KYCStatus == services.KYCStatusReview {
		// Update address and contacts in BBVA
		bank, err := partnerbank.GetBusinessBank(partnerbank.ProviderNameBBVA)
		if err != nil {
			return nil, err
		}

		srv := bank.BusinessEntityService(db.sourceReq.PartnerBankRequest())

		ureq := partnerbank.UpdateBusinessRequest{
			BusinessID: partnerbank.BusinessID(u.ID),
			LegalName:  shared.StringValue(bus.LegalName),
			DBA:        bus.DBA,
		}

		if bus.EntityType != nil {
			ureq.EntityType = partnerbank.BusinessEntity(*bus.EntityType)
		}

		if bus.IndustryType != nil {
			ureq.IndustryType = partnerbank.BusinessIndustry(*bus.IndustryType)
		}

		if bus.Purpose != nil {
			ureq.Purpose = *bus.Purpose
		}

		if bus.OperationType != nil {
			ureq.OperationType = partnerbank.BusinessOperationType(*bus.OperationType)
		}

		if bus.TaxIDType != nil && bus.TaxID != nil {
			ureq.TaxIDType = partnerbank.BusinessTaxIDType(*bus.TaxIDType)
			taxID := *bus.TaxID
			ureq.TaxID = string(taxID)
		}

		if bus.OriginCountry != nil {
			ureq.OriginCountry = partnerbank.Country(*bus.OriginCountry)
		}

		if bus.OriginState != nil {
			ureq.OriginState = *bus.OriginState
		}

		if bus.OriginDate != nil {
			ureq.OriginDate = bus.OriginDate.Time()
		}

		// Get entity formation doc
		doc, err := db.checkEntityFormationDoc(bus)
		if err != nil {
			return nil, err
		}

		ureq.EntityFormation = &partnerbank.EntityFormationRequest{
			IssueDate:      doc.IssuedDate.Time(),
			ExpirationDate: doc.ExpirationDate.Time(),
		}

		if doc.DocType != nil {
			ureq.EntityFormation.DocumentType = partnerbank.BusinessIdentityDocument(*doc.DocType)
		}

		if doc.Number != nil {
			ureq.EntityFormation.Number = string(*doc.Number)
		}

		// Update business on bank side - log errors
		var wait sync.WaitGroup
		wait.Add(1)
		go func() {
			defer wait.Done()
			_, err = srv.Update(ureq)
		}()

		// Update only on change
		if u.Email != nil {
			wait.Add(1)
			go func() {
				defer wait.Done()
				err := srv.UpdateContact(
					partnerbank.BusinessID(u.ID),
					partnerbank.BusinessPropertyTypeContactEmail,
					*u.Email,
				)
				if err != nil {
					// Log errors
					log.Println(err)
				}
			}()
		}

		if u.Phone != nil {
			wait.Add(1)
			go func() {
				defer wait.Done()

				err := srv.UpdateContact(
					partnerbank.BusinessID(u.ID),
					partnerbank.BusinessPropertyTypeContactPhone,
					*u.Phone,
				)
				if err != nil {
					// Log errors
					log.Println(err)
				}
			}()
		}

		if u.LegalAddress != nil {
			wait.Add(1)
			go func() {
				defer wait.Done()
				err := srv.UpdateAddress(
					partnerbank.BusinessID(u.ID),
					partnerbank.BusinessPropertyTypeAddressLegal,
					u.LegalAddress.ToPartnerBankAddress(services.AddressTypeLegal),
				)
				if err != nil {
					// Log errors
					log.Println(err)
				}
			}()
		}

		if u.HeadquarterAddress != nil {
			wait.Add(1)
			go func() {
				defer wait.Done()
				err := srv.UpdateAddress(
					partnerbank.BusinessID(u.ID),
					partnerbank.BusinessPropertyTypeAddressHeadquarter,
					u.HeadquarterAddress.ToPartnerBankAddress(services.AddressTypeHeadquarter),
				)
				if err != nil {
					// Log errors
					log.Println(err)
				}
			}()
		}

		if u.MailingAddress != nil {
			wait.Add(1)
			go func() {
				defer wait.Done()
				err := srv.UpdateAddress(
					partnerbank.BusinessID(u.ID),
					partnerbank.BusinessPropertyTypeAddressMailing,
					u.MailingAddress.ToPartnerBankAddress(services.AddressTypeMailing),
				)
				if err != nil {
					// Log errors
					log.Println(err)
				}
			}()
		}

		// Wait for completion
		wait.Wait()
	}

	// Send to segment
	segment.NewSegmentService().PushToAnalyticsQueue(b.OwnerID, segment.CategoryBusiness, segment.ActionUpdate, b)
	return bus, err
}

func (db *businessDatastore) UpdateSubscription(u BusinessSubscriptionUpdate) (*Business, error) {
	var columns []string

	if u.SubscriptionDecisionDate != nil {
		columns = append(columns, "subscription_decision_date = :subscription_decision_date")
	}

	if u.SubscriptionStatus != nil {
		columns = append(columns, "subscription_status = :subscription_status")
	}

	if u.SubscriptionStartDate != nil {
		columns = append(columns, "subscription_start_date = :subscription_start_date")
	}

	// Update if changes requested
	if len(columns) > 0 {
		// Update business
		_, err := db.NamedExec(
			fmt.Sprintf(
				"UPDATE business SET %s WHERE id = '%s'",
				strings.Join(columns, ", "),
				u.ID,
			), u,
		)
		if err != nil {
			return nil, errors.Cause(err)
		}
	}
	return db.GetById(u.ID)
}

func (db *businessDatastore) updateVerification(u BusinessVerificationUpdate) (*Business, error) {
	_, err := db.NamedExec(
		fmt.Sprintf("UPDATE business SET kyc_status = :kyc_status WHERE id = '%s'", u.ID), u,
	)
	if err != nil {
		return nil, errors.Cause(err)
	}

	return db.GetById(u.ID)
}

func (db *businessDatastore) Deactivate(id shared.BusinessID) error {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(id)
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE business SET deactivated = CURRENT_TIMESTAMP WHERE id = $1", id)
	return err
}

func (db *businessDatastore) Submit(id shared.BusinessID) (*Business, error) {
	// Update verification
	b, err := db.GetById(id)
	if err != nil {
		return nil, err
	}

	if b.Deactivated != nil {
		return nil, errors.New("Business has already been deactivated")
	}

	switch b.KYCStatus {
	case services.KYCStatusSubmitted:
		return nil, errors.New("Business has already been submitted")
	case services.KYCStatusReview:
		return nil, errors.New("Business has already in review")
	case services.KYCStatusApproved:
		return nil, errors.New("Business has already been approved")
	case services.KYCStatusDeclined:
		return nil, errors.New("Business has already been declined")
	}

	// Update in c status in database
	b, err = db.updateVerification(
		BusinessVerificationUpdate{
			ID:        b.ID,
			KYCStatus: services.KYCStatusSubmitted,
		},
	)

	// Start CSP process
	return b, StartReview(b)
}

func (db *businessDatastore) checkEntityFormationDoc(b *Business) (*document.BusinessDocument, error) {
	return document.NewBusinessDocumentService(db.sourceReq).GetByID(*b.FormationDocumentID, b.ID)
}

func generateEmployerNumber() string {
	var prefix string
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const numset = "0123456789"

	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()),
	)

	// Use two random alpha values
	prefix = string(charset[seededRand.Intn(len(charset))]) + string(charset[seededRand.Intn(len(charset))])

	// Generate six digit value
	b := make([]byte, 6)
	for i := range b {
		b[i] = numset[seededRand.Intn(len(numset))]
	}

	return prefix + string(b)
}
