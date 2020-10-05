/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for all c related services
package user

import (
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	_ "github.com/wiseco/core-platform/partner/bank/bbva"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/services/document"
	emailService "github.com/wiseco/core-platform/services/email"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/go-lib/grpc"
	grpcEmail "github.com/wiseco/protobuf/golang/verification/email"
)

type consumerDatastore struct {
	sourceReq services.SourceRequest
	*sqlx.DB
}

type ConsumerService interface {
	// Fetch operations
	GetByID(shared.ConsumerID) (*Consumer, error)

	// Create c
	Create(ConsumerCreate) (*shared.ConsumerID, error)

	// Create c from auth data
	CreateFromAuth(ConsumerAuthCreate) (*shared.ConsumerID, error)

	// Update c
	Update(ConsumerUpdate) (*Consumer, error)

	// Deactivate c by id
	Deactivate(shared.ConsumerID) error

	// Consumer verification (KYC)
	Submit(shared.ConsumerID) (*Consumer, error)
	StartVerification(shared.ConsumerID, bool) (*ConsumerKYCResponse, error)
	GetVerification(shared.ConsumerID) (*ConsumerKYCResponse, error)
}

func NewConsumerService(r services.SourceRequest) ConsumerService {
	return &consumerDatastore{r, data.DBWrite}
}

// NewConsumerServiceWithout returns an c service without a source request
func NewConsumerServiceWithout() ConsumerService {
	return &consumerDatastore{services.SourceRequest{}, data.DBWrite}
}

func (db *consumerDatastore) GetByID(id shared.ConsumerID) (*Consumer, error) {
	return db.getByID(id)
}

func (db *consumerDatastore) getByID(id shared.ConsumerID) (*Consumer, error) {
	u := Consumer{}

	err := db.Get(&u, "SELECT * FROM consumer WHERE id = $1", id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &u, nil
}

func (db *consumerDatastore) Create(u ConsumerCreate) (*shared.ConsumerID, error) {
	if u.LegalAddress != nil {
		addr, err := services.ValidateAddress(u.LegalAddress, services.AddressTypeLegal)
		if err != nil {
			return nil, err
		}

		u.LegalAddress = addr
	}

	if u.MailingAddress != nil {
		addr, err := services.ValidateAddress(u.MailingAddress, services.AddressTypeMailing)
		if err != nil {
			return nil, err
		}

		u.MailingAddress = addr
	}

	if u.WorkAddress != nil {
		addr, err := services.ValidateAddress(u.WorkAddress, services.AddressTypeWork)
		if err != nil {
			return nil, err
		}

		u.WorkAddress = addr
	}

	if u.TaxID != nil {
		taxID, err := services.ValidateTaxID(u.TaxID, u.TaxIDType)
		if err != nil {
			return nil, err
		}

		u.TaxID = taxID
	} else {
		u.TaxIDType = nil
	}

	if u.Email != nil {
		esrvc := emailService.NewEmailService(db.sourceReq)

		emailType := emailService.TypeConsumer

		ea := emailService.EmailAddress(*u.Email)

		if u.IsBusinessMember {
			emailType = emailService.TypeBusinessMember
		} else {
			isAvailable, err := esrvc.IsAvailable(ea, emailService.TypeConsumer)

			if err != nil {
				return nil, err
			}

			if !isAvailable {
				return nil, errors.New("Email is not available")
			}
		}

		//This path is hit via POST business/business_id/member
		//As of right now these consumers cannot login, so for now lets set email status to active
		//TODO multi user businesses: an email should go out inviting user to create a password(this would include a verification token as well)
		ec := emailService.EmailCreate{
			EmailAddress: emailService.EmailAddress(*u.Email),
			Status:       emailService.StatusActive,
			Type:         emailType,
		}

		e, err := esrvc.Create(&ec)

		if err != nil {
			return nil, err
		}

		if !u.IsBusinessMember {
			sn, err := grpc.GetConnectionStringForService(grpc.ServiceNameVerification)
			if err != nil {
				return nil, err
			}

			client, err := grpc.NewInsecureClient(sn)
			if err != nil {
				return nil, err
			}

			defer client.CloseAndCancel()

			c := grpcEmail.NewEmailServiceClient(client.GetConn())

			req := grpcEmail.SendVerificationLinkRequest{
				CoreEmailId:  string(e.ID),
				EmailAddress: string(e.EmailAddress),
				FirstName:    u.FirstName,
				FullName:     u.FirstName + " " + u.LastName,
			}

			_, err = c.SendVerificationLink(client.GetContext(), &req)
			if err != nil {
				return nil, err
			}
		}

		u.EmailID = e.ID
	}

	// Activity type is ignored on create until consumer bank is ready
	sql := `
		INSERT INTO consumer(
			first_name, last_name, email, phone, date_of_birth, tax_id, tax_id_type, legal_address,
			mailing_address, work_address, residency, citizenship_countries, occupation,
			income_type, email_id
		)
		VALUES(
			:first_name, :last_name, :email, :phone, :date_of_birth, :tax_id, :tax_id_type,
			:legal_address, :mailing_address, :work_address, :residency, :citizenship_countries,
			:occupation, :income_type, :email_id
		)
		RETURNING id`

	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	var id shared.ConsumerID
	err = stmt.Get(&id, &u)
	if err != nil {
		log.Println("Create member failed")
		return nil, err
	}

	return &id, nil
}

//CreateFromAuth creates a minimal c from cognito
func (db *consumerDatastore) CreateFromAuth(u ConsumerAuthCreate) (*shared.ConsumerID, error) {
	// Insert statement
	sql := `
		INSERT INTO consumer(phone)
		VALUES (:phone)
		RETURNING id`

	// Execute
	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	// Return id
	var id shared.ConsumerID
	err = stmt.Get(&id, &u)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (db *consumerDatastore) Update(u ConsumerUpdate) (*Consumer, error) {
	m, err := db.getByID(u.ID)
	if err != nil {
		return nil, err
	}

	var columns []string

	if u.Phone != nil {
		columns = append(columns, "phone = :phone")
	}

	if u.Email != nil {
		var e *emailService.Email

		ea := emailService.EmailAddress(*u.Email)

		columns = append(columns, "email = :email")

		esrvc := emailService.NewEmailService(db.sourceReq)

		isAvailable, err := esrvc.IsAvailable(ea, emailService.TypeConsumer)
		if err != nil {
			return nil, err
		}

		if m.EmailID != shared.EmailID("") {
			e, err = esrvc.GetByID(m.EmailID)

			if err != nil {
				return nil, err
			}
		}

		if isAvailable {
			sendVerification := false

			if m.EmailID != shared.EmailID("") {
				err = esrvc.Deactivate(m.EmailID)

				if err != nil {
					return nil, err
				}
			} else {
				sendVerification = true
			}

			ec := emailService.EmailCreate{
				EmailAddress: emailService.EmailAddress(*u.Email),
				Status:       emailService.StatusActive,
				Type:         emailService.TypeConsumer,
			}

			e, err = esrvc.Create(&ec)

			if err != nil {
				return nil, err
			}

			u.EmailID = e.ID

			columns = append(columns, "email_id = :email_id")

			//Has this consumer had an email associated with it before? if not lets send a verification email
			if sendVerification {
				sn, err := grpc.GetConnectionStringForService(grpc.ServiceNameVerification)
				if err != nil {
					return nil, err
				}

				client, err := grpc.NewInsecureClient(sn)
				if err != nil {
					return nil, err
				}

				defer client.CloseAndCancel()

				c := grpcEmail.NewEmailServiceClient(client.GetConn())

				req := grpcEmail.SendVerificationLinkRequest{
					CoreEmailId:  string(e.ID),
					EmailAddress: string(e.EmailAddress),
					FirstName:    u.GetFirstName(),
					FullName:     u.GetFullName(),
				}

				_, err = c.SendVerificationLink(client.GetContext(), &req)
				if err != nil {
					return nil, err
				}
			}
		} else if m.EmailID == shared.EmailID("") || (e != nil && e.EmailAddress != ea) {
			//we want to error if the email is not avialable and the user is trying to enter a net new email
			return nil, errors.New("Email is not available")
		}
	}

	if u.Occupation != nil {
		columns = append(columns, "occupation = :occupation")
	}

	if u.IncomeType != nil {
		columns = append(columns, "income_type = :income_type")
	}

	if u.ActivityType != nil {
		columns = append(columns, "activity_type = :activity_type")
	}

	// Once KYC is approved fields below should be locked
	if m.KYCStatus != services.KYCStatusApproved {
		if u.FirstName != nil {
			columns = append(columns, "first_name = :first_name")
		}

		if u.MiddleName != nil {
			columns = append(columns, "middle_name = :middle_name")
		}

		if u.LastName != nil {
			columns = append(columns, "last_name = :last_name")
		}

		if u.DateOfBirth != nil {
			columns = append(columns, "date_of_birth = :date_of_birth")
		}

		if u.TaxID != nil {
			taxID, err := services.ValidateTaxID(u.TaxID, u.TaxIDType)
			if err != nil {
				return nil, err
			}

			u.TaxID = taxID
			columns = append(columns, "tax_id = :tax_id")
			columns = append(columns, "tax_id_type = :tax_id_type")
		}

		if u.Residency != nil {
			columns = append(columns, "residency = :residency")
		}

		if u.CitizenshipCountries != nil {
			columns = append(columns, "citizenship_countries = :citizenship_countries")
		}

		if u.LegalAddress != nil {
			addr, err := services.ValidateAddress(u.LegalAddress, services.AddressTypeLegal)
			if err != nil {
				return nil, err
			}

			u.LegalAddress = addr
			columns = append(columns, "legal_address = :legal_address")
		}

		if u.MailingAddress != nil {
			addr, err := services.ValidateAddress(u.MailingAddress, services.AddressTypeMailing)
			if err != nil {
				return nil, err
			}

			u.MailingAddress = addr
			columns = append(columns, "mailing_address = :mailing_address")
		}

		if u.WorkAddress != nil {
			addr, err := services.ValidateAddress(u.WorkAddress, services.AddressTypeWork)
			if err != nil {
				return nil, err
			}

			u.WorkAddress = addr
			columns = append(columns, "work_address = :work_address")
		}
	}

	// No changes requested - return c
	if len(columns) == 0 {
		return db.getByID(u.ID)
	}

	_, err = db.NamedExec(fmt.Sprintf("UPDATE consumer SET %s WHERE id = '%s'", strings.Join(columns, ", "), u.ID), u)

	if err != nil {
		log.Println(err)
		return nil, errors.Cause(err)
	}

	switch m.KYCStatus {
	case services.KYCStatusNotStarted, services.KYCStatusSubmitted:
		break
	default:
		// TODO: Update async via SQS (CSP?)
		// Update address and contacts in BBVA
		bank, err := partnerbank.GetConsumerBank(partnerbank.ProviderNameBBVA)
		if err != nil {
			return nil, err
		}

		consumerUpdate := partnerbank.UpdateConsumerRequest{
			ConsumerID: partnerbank.ConsumerID(u.ID),
		}
		isUpdateRequired := false

		srv := bank.ConsumerEntityService(db.sourceReq.PartnerBankRequest())
		if u.Email != nil {
			err := srv.UpdateContact(
				partnerbank.ConsumerID(u.ID),
				partnerbank.ConsumerPropertyTypeContactEmail,
				*u.Email,
			)
			if err != nil {
				// Log error
				log.Println(err)
			}
		}

		if u.Phone != nil {
			err := srv.UpdateContact(
				partnerbank.ConsumerID(u.ID),
				partnerbank.ConsumerPropertyTypeContactPhone,
				*u.Phone,
			)
			if err != nil {
				// Log error
				log.Println(err)
			}
		}

		if u.LegalAddress != nil {
			err := srv.UpdateAddress(
				partnerbank.ConsumerID(u.ID),
				partnerbank.ConsumerPropertyTypeAddressLegal,
				partnerbank.AddressRequest{
					Type:    partnerbank.AddressRequestTypeLegal,
					Line1:   u.LegalAddress.StreetAddress,
					Line2:   u.LegalAddress.AddressLine2,
					City:    u.LegalAddress.City,
					State:   u.LegalAddress.State,
					ZipCode: u.LegalAddress.PostalCode,
					Country: partnerbank.Country(u.LegalAddress.Country),
				},
			)
			if err != nil {
				// Log error
				log.Println(err)
			}
		}

		if u.WorkAddress != nil {
			err := srv.UpdateAddress(
				partnerbank.ConsumerID(u.ID),
				partnerbank.ConsumerPropertyTypeAddressWork,
				partnerbank.AddressRequest{
					Type:    partnerbank.AddressRequestTypeWork,
					Line1:   u.WorkAddress.StreetAddress,
					Line2:   u.WorkAddress.AddressLine2,
					City:    u.WorkAddress.City,
					State:   u.WorkAddress.State,
					ZipCode: u.WorkAddress.PostalCode,
					Country: partnerbank.Country(u.WorkAddress.Country),
				},
			)
			if err != nil {
				// Log error
				log.Println(err)
			}
		}

		if u.MailingAddress != nil {
			err := srv.UpdateAddress(
				partnerbank.ConsumerID(u.ID),
				partnerbank.ConsumerPropertyTypeAddressMailing,
				partnerbank.AddressRequest{
					Type:    partnerbank.AddressRequestTypeMailing,
					Line1:   u.MailingAddress.StreetAddress,
					Line2:   u.MailingAddress.AddressLine2,
					City:    u.MailingAddress.City,
					State:   u.MailingAddress.State,
					ZipCode: u.MailingAddress.PostalCode,
					Country: partnerbank.Country(u.MailingAddress.Country),
				},
			)
			if err != nil {
				// Log error
				log.Println(err)
			}
		}

		if m.KYCStatus == services.KYCStatusReview {
			if u.FirstName != nil {
				isUpdateRequired = true
				consumerUpdate.FirstName = *u.FirstName
			}

			if u.MiddleName != nil {
				isUpdateRequired = true
				consumerUpdate.MiddleName = *u.MiddleName
			}

			if u.LastName != nil {
				isUpdateRequired = true
				consumerUpdate.LastName = *u.LastName
			}

			if u.TaxIDType != nil && u.TaxID != nil {
				isUpdateRequired = true
				consumerUpdate.TaxIDType = partnerbank.ConsumerTaxIDType(*u.TaxIDType)
				taxID := *u.TaxID
				consumerUpdate.TaxID = string(taxID)
			}

			if u.DateOfBirth != nil {
				isUpdateRequired = true
				consumerUpdate.DateOfBirth = u.DateOfBirth.Time()
			}
		}

		if u.CitizenshipCountries != nil && len(*u.CitizenshipCountries) > 0 {
			isUpdateRequired = true
			country := *u.CitizenshipCountries
			consumerUpdate.CitizenshipCountry = partnerbank.Country(country[0])
		}

		if u.Residency != nil {
			isUpdateRequired = true
			consumerUpdate.Residency = partnerbank.ConsumerResidency{
				Country: partnerbank.Country(u.Residency.Country),
				Status:  partnerbank.ConsumerResidencyStatus(u.Residency.Status),
			}
		}

		if isUpdateRequired {
			_, err = srv.Update(consumerUpdate)
			if err != nil {
				// Log error
				log.Println(err)
			}
		}
	}

	return db.getByID(u.ID)
}

func (db *consumerDatastore) updateVerification(u ConsumerVerificationUpdate) (*Consumer, error) {
	_, err := db.NamedExec(
		fmt.Sprintf("UPDATE consumer SET kyc_status = :kyc_status WHERE id = '%s'", u.ID), u,
	)
	if err != nil {
		log.Println(err)
		return nil, errors.Cause(err)
	}

	return db.getByID(u.ID)
}

func (db *consumerDatastore) Deactivate(id shared.ConsumerID) error {
	_, err := db.Exec(fmt.Sprintf("UPDATE consumer SET deactivated = CURRENT_TIMESTAMP WHERE id = '%s'", id))
	return err
}

func (db *consumerDatastore) Submit(id shared.ConsumerID) (*Consumer, error) {
	c, err := db.getByID(id)
	if err != nil {
		return nil, err
	}

	if c.Deactivated != nil {
		return nil, errors.New("Consumer has already been deactivated")
	}

	switch c.KYCStatus {
	case services.KYCStatusSubmitted:
		return nil, errors.New("Consumer has already been submitted")
	case services.KYCStatusReview:
		return nil, errors.New("Consumer has already in review")
	case services.KYCStatusApproved:
		return nil, errors.New("Consumer has already been approved")
	case services.KYCStatusDeclined:
		return nil, errors.New("Consumer has already been declined")
	}

	update := ConsumerVerificationUpdate{
		ID:        c.ID,
		KYCStatus: services.KYCStatusSubmitted,
	}

	consumer, err := db.updateVerification(update)

	// Update in c status in database
	return consumer, StartReview(consumer)
}

func (db *consumerDatastore) StartVerification(id shared.ConsumerID, internal bool) (*ConsumerKYCResponse, error) {
	c, err := db.getByID(id)
	if err != nil {
		return nil, ConsumerKYCError{
			RawError:  err,
			ErrorType: services.KYCErrorTypeOther,
		}
	}

	switch c.KYCStatus {
	case services.KYCStatusNotStarted:
		return nil, ConsumerKYCError{
			RawError:   errors.New("Consumer has not been submitted"),
			ErrorType:  services.KYCErrorTypeOther,
			ConsumerID: &c.ID,
		}
	case services.KYCStatusApproved:
		return nil, ConsumerKYCError{
			RawError:   errors.New("Consumer has already been approved"),
			ErrorType:  services.KYCErrorTypeOther,
			ConsumerID: &c.ID,
		}
	case services.KYCStatusDeclined:
		return nil, ConsumerKYCError{
			RawError:   errors.New("Consumer has already been declined"),
			ErrorType:  services.KYCErrorTypeOther,
			ConsumerID: &c.ID,
		}
	case services.KYCStatusReview:
		return db.continueVerification(c)
	}

	if c.Deactivated != nil {
		return nil, ConsumerKYCError{
			ErrorType:  services.KYCErrorTypeDeactivated,
			ConsumerID: &c.ID,
		}
	} else if c.IsRestricted {
		return nil, ConsumerKYCError{
			ErrorType:  services.KYCErrorTypeRestricted,
			ConsumerID: &c.ID,
		}
	}

	// Check params
	kycErrors := []string{}
	if len(c.FirstName) < 2 {
		kycErrors = append(kycErrors, services.KYCParamErrorFirstName)
	}

	if len(c.LastName) < 2 {
		kycErrors = append(kycErrors, services.KYCParamErrorLastName)
	}

	if c.Email == nil || len(*c.Email) < 6 {
		kycErrors = append(kycErrors, services.KYCParamErrorEmail)
	}

	if c.Phone == nil || len(*c.Phone) < 10 {
		kycErrors = append(kycErrors, services.KYCParamErrorPhone)
	}

	if c.DateOfBirth == nil {
		kycErrors = append(kycErrors, services.KYCParamErrorDateOfBirth)
	}

	if c.TaxID == nil {
		kycErrors = append(kycErrors, services.KYCParamErrorTaxID)
	}

	if c.TaxIDType == nil {
		kycErrors = append(kycErrors, services.KYCParamErrorTaxIDType)
	}

	if c.LegalAddress == nil {
		kycErrors = append(kycErrors, services.KYCParamErrorLegalAddress)
	}

	if c.Residency == nil {
		kycErrors = append(kycErrors, services.KYCParamErrorResidency)
	}

	if len(c.CitizenshipCountries) == 0 {
		kycErrors = append(kycErrors, services.KYCParamErrorCitizenship)
	}

	if c.Occupation == nil {
		kycErrors = append(kycErrors, services.KYCParamErrorOccupation)
	}

	if len(c.IncomeType) == 0 {
		kycErrors = append(kycErrors, services.KYCParamErrorIncomeType)
	}

	if len(c.ActivityType) == 0 {
		kycErrors = append(kycErrors, services.KYCParamErrorActivityType)
	}

	if len(kycErrors) > 0 {
		return nil, ConsumerKYCError{
			ErrorType:  services.KYCErrorTypeParam,
			Values:     kycErrors,
			ConsumerID: &c.ID,
		}
	}

	bank, err := partnerbank.GetConsumerBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		return nil, ConsumerKYCError{
			RawError:   err,
			ErrorType:  services.KYCErrorTypeOther,
			ConsumerID: &c.ID,
		}
	}

	r := db.sourceReq.PartnerBankRequest()

	creq := partnerbank.CreateConsumerRequest{
		ConsumerID:  partnerbank.ConsumerID(c.ID),
		FirstName:   c.FirstName,
		MiddleName:  c.MiddleName,
		LastName:    c.LastName,
		TaxID:       c.TaxID.String(),
		TaxIDType:   partnerbank.ConsumerTaxIDType(*c.TaxIDType),
		DateOfBirth: c.DateOfBirth.Time(),
		Phone:       *c.Phone,
		Email:       *c.Email,
		Residency: partnerbank.ConsumerResidency{
			Country: partnerbank.Country(c.Residency.Country),
			Status:  partnerbank.ConsumerResidencyStatus(c.Residency.Status),
		},
		CitizenshipCountry: partnerbank.Country(c.CitizenshipCountries[0]),
		Occupation:         partnerbank.ConsumerOccupation(*c.Occupation),
		Income:             c.IncomeType.ToPartnerBankIncome(),
		ExpectedActivity:   c.ActivityType.ToPartnerBankActivity(),
		LegalAddress:       c.LegalAddress.ToPartnerBankAddress(services.AddressTypeLegal),
	}

	if c.MailingAddress != nil {
		mailingAddress := c.MailingAddress.ToPartnerBankAddress(services.AddressTypeMailing)
		creq.MailingAddress = &mailingAddress
	}

	// In case of resident alien, consumer identification document needs to be set
	isUSCitizen := false
	for _, country := range c.CitizenshipCountries {
		if country == string(partnerbank.CountryUS) {
			isUSCitizen = true
		}
	}
	var docs = []document.ConsumerDocument{}
	if internal {
		docs, err = document.NewConsumerDocumentService(db.sourceReq).ListInternal(id, 0, 10)
	} else {
		docs, err = document.NewConsumerDocumentService(db.sourceReq).List(id, 0, 10)
	}

	if err != nil {
		return nil, ConsumerKYCError{
			RawError:   err,
			ErrorType:  services.KYCErrorTypeOther,
			ConsumerID: &c.ID,
		}
	}

	if !isUSCitizen && len(docs) == 0 {
		return nil, ConsumerKYCError{
			RawError:   errors.New("No document found"),
			ErrorType:  services.KYCErrorTypeOther,
			ConsumerID: &c.ID,
		}
	}

	var idReq []*partnerbank.ConsumerIdentificationRequest
	for _, doc := range docs {
		if doc.DocType == nil {
			continue
		}

		docType := document.ConsumerIdentityDocument(*doc.DocType)
		switch docType {
		case
			document.ConsumerIdentityDocumentDriversLicense,
			document.ConsumerIdentityDocumentStateID:
			// For non-US Citizens skip DL and State ID - Passport required
			if !isUSCitizen {
				break
			}

			fallthrough
		case
			document.ConsumerIdentityDocumentPassport,
			document.ConsumerIdentityDocumentAlienRegistrationCard,
			document.ConsumerIdentityDocumentUSAVisaH1B,
			document.ConsumerIdentityDocumentUSAVisaH1C,
			document.ConsumerIdentityDocumentUSAVisaH2A,
			document.ConsumerIdentityDocumentUSAVisaH2B,
			document.ConsumerIdentityDocumentUSAVisaH3,
			document.ConsumerIdentityDocumentUSAVisaL1A,
			document.ConsumerIdentityDocumentUSAVisaL1B,
			document.ConsumerIdentityDocumentUSAVisaO1,
			document.ConsumerIdentityDocumentUSAVisaE1,
			document.ConsumerIdentityDocumentUSAVisaE3,
			document.ConsumerIdentityDocumentUSAVisaI,
			document.ConsumerIdentityDocumentUSAVisaP,
			document.ConsumerIdentityDocumentUSAVisaTN,
			document.ConsumerIdentityDocumentUSAVisaTD,
			document.ConsumerIdentityDocumentUSAVisaR1:
			// Create identity request from document
			idRequest, kycErr := createConsumerIdentityRequest(doc)
			kycErrors = append(kycErrors, kycErr...)
			if idRequest != nil {
				idReq = append(idReq, idRequest)
			}
		}
	}

	creq.Identification = idReq

	resp, err := bank.ConsumerEntityService(r).Create(creq)
	if err != nil {
		return nil, ConsumerKYCError{
			RawError:   err,
			ErrorType:  services.KYCErrorTypeOther,
			ConsumerID: &c.ID,
		}
	}

	// Update in c status in database
	c, err = db.updateVerification(
		ConsumerVerificationUpdate{
			ID:        c.ID,
			KYCStatus: services.PartnerKYCStatusFromMap[resp.KYC.Status],
		},
	)
	if err != nil {
		return nil, ConsumerKYCError{
			RawError:   err,
			ErrorType:  services.KYCErrorTypeOther,
			ConsumerID: &c.ID,
		}
	}

	// Return KYC response
	var items []string
	for _, item := range resp.KYC.IDVerifyRequired {
		items = append(items, string(item))
	}

	return &ConsumerKYCResponse{
		Status:      services.KYCStatus(resp.KYC.Status),
		ReviewItems: items,
		ConsumerID:  &c.ID,
	}, nil
}

func createConsumerIdentityRequest(doc document.ConsumerDocument) (*partnerbank.ConsumerIdentificationRequest, []string) {
	kycErrors := []string{}

	if doc.DocType == nil {
		kycErrors = append(kycErrors, services.KYCParamErrorDocType)
	}

	if doc.IssuingCountry == nil {
		kycErrors = append(kycErrors, services.KYCParamErrorIssuingCountry)
	}

	if doc.IssuedDate == nil {
		kycErrors = append(kycErrors, services.KYCParamErrorIssuedDate)
	}

	if doc.ExpirationDate == nil {
		kycErrors = append(kycErrors, services.KYCParamErrorExpirationDate)
	}

	if doc.Number == nil {
		kycErrors = append(kycErrors, services.KYCParamErrorDocNumber)
	}

	if len(kycErrors) > 0 {
		return nil, kycErrors
	}

	idRequest := partnerbank.ConsumerIdentificationRequest{
		DocumentType:   document.ConsumerDocTypeToBankDocType[document.ConsumerIdentityDocument(*doc.DocType)],
		IssueCountry:   partnerbank.Country(*doc.IssuingCountry),
		IssueDate:      doc.IssuedDate.Time(),
		ExpirationDate: doc.ExpirationDate.Time(),
		Number:         doc.Number.String(),
	}

	if doc.IssuingState != nil {
		idRequest.IssueState = *doc.IssuingState
	}

	return &idRequest, kycErrors
}

func (db *consumerDatastore) GetVerification(id shared.ConsumerID) (*ConsumerKYCResponse, error) {
	c, err := db.getByID(id)
	if err != nil {
		return nil, ConsumerKYCError{
			RawError:  err,
			ErrorType: services.KYCErrorTypeOther,
		}
	}

	switch c.KYCStatus {
	case services.KYCStatusNotStarted, services.KYCStatusSubmitted:
		return nil, ConsumerKYCError{
			RawError:   errors.New("Consumer has not entered review"),
			ErrorType:  services.KYCErrorTypeOther,
			ConsumerID: &c.ID,
		}
	}

	bank, err := partnerbank.GetConsumerBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		return nil, ConsumerKYCError{
			RawError:   err,
			ErrorType:  services.KYCErrorTypeOther,
			ConsumerID: &c.ID,
		}
	}

	resp, err := bank.ConsumerEntityService(db.sourceReq.PartnerBankRequest()).Status(partnerbank.ConsumerID(id))
	if err != nil {
		return nil, ConsumerKYCError{
			RawError:   err,
			ErrorType:  services.KYCErrorTypeOther,
			ConsumerID: &c.ID,
		}
	}

	// Update in c status in database
	c, err = db.updateVerification(
		ConsumerVerificationUpdate{
			ID:        c.ID,
			KYCStatus: services.PartnerKYCStatusFromMap[resp.KYC.Status],
		},
	)
	if err != nil {
		return nil, ConsumerKYCError{
			RawError:   err,
			ErrorType:  services.KYCErrorTypeOther,
			ConsumerID: &c.ID,
		}
	}

	// Return KYC response
	var items []string
	for _, item := range resp.KYC.IDVerifyRequired {
		items = append(items, string(item))
	}

	return &ConsumerKYCResponse{
		Status:      services.KYCStatus(resp.KYC.Status),
		ReviewItems: items,
		ConsumerID:  &c.ID,
	}, nil
}

func (db *consumerDatastore) continueVerification(c *Consumer) (*ConsumerKYCResponse, error) {
	// Check params
	kycErrors := []string{}
	if len(c.FirstName) < 2 {
		kycErrors = append(kycErrors, services.KYCParamErrorFirstName)
	}

	if len(c.LastName) < 2 {
		kycErrors = append(kycErrors, services.KYCParamErrorLastName)
	}

	if c.Email == nil || len(*c.Email) < 6 {
		kycErrors = append(kycErrors, services.KYCParamErrorEmail)
	}

	if c.Phone == nil || len(*c.Phone) < 10 {
		kycErrors = append(kycErrors, services.KYCParamErrorPhone)
	}

	if c.DateOfBirth == nil {
		kycErrors = append(kycErrors, services.KYCParamErrorDateOfBirth)
	}

	if c.TaxID == nil {
		kycErrors = append(kycErrors, services.KYCParamErrorTaxID)
	}

	if c.TaxIDType == nil {
		kycErrors = append(kycErrors, services.KYCParamErrorTaxIDType)
	}

	if c.Residency == nil {
		kycErrors = append(kycErrors, services.KYCParamErrorResidency)
	}

	if len(c.CitizenshipCountries) == 0 {
		kycErrors = append(kycErrors, services.KYCParamErrorCitizenship)
	}

	if len(kycErrors) > 0 {
		return nil, ConsumerKYCError{
			ErrorType:  services.KYCErrorTypeParam,
			Values:     kycErrors,
			ConsumerID: &c.ID,
		}
	}

	bank, err := partnerbank.GetConsumerBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		return nil, ConsumerKYCError{
			RawError:   err,
			ErrorType:  services.KYCErrorTypeOther,
			ConsumerID: &c.ID,
		}
	}

	r := db.sourceReq.PartnerBankRequest()

	ureq := partnerbank.UpdateConsumerRequest{
		ConsumerID:  partnerbank.ConsumerID(c.ID),
		FirstName:   c.FirstName,
		MiddleName:  c.MiddleName,
		LastName:    c.LastName,
		TaxID:       c.TaxID.String(),
		TaxIDType:   partnerbank.ConsumerTaxIDType(*c.TaxIDType),
		DateOfBirth: c.DateOfBirth.Time(),
		Residency: partnerbank.ConsumerResidency{
			Country: partnerbank.Country(c.Residency.Country),
			Status:  partnerbank.ConsumerResidencyStatus(c.Residency.Status),
		},
		CitizenshipCountry: partnerbank.Country(c.CitizenshipCountries[0]),
	}

	resp, err := bank.ConsumerEntityService(r).Update(ureq)
	if err != nil {
		return nil, ConsumerKYCError{
			RawError:   err,
			ErrorType:  services.KYCErrorTypeOther,
			ConsumerID: &c.ID,
		}
	}

	// Update in c status in database
	c, err = db.updateVerification(
		ConsumerVerificationUpdate{
			ID:        c.ID,
			KYCStatus: services.PartnerKYCStatusFromMap[resp.KYC.Status],
		},
	)
	if err != nil {
		return nil, ConsumerKYCError{
			RawError:   err,
			ErrorType:  services.KYCErrorTypeOther,
			ConsumerID: &c.ID,
		}
	}

	// Return KYC response
	var items []string
	for _, item := range resp.KYC.IDVerifyRequired {
		items = append(items, string(item))
	}

	return &ConsumerKYCResponse{
		Status:      services.KYCStatus(resp.KYC.Status),
		ReviewItems: items,
		ConsumerID:  &c.ID,
	}, nil
}
