package consumer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	partnerbank "github.com/wiseco/core-platform/partner/bank"
	_ "github.com/wiseco/core-platform/partner/bank/bbva"
	"github.com/wiseco/core-platform/services"
	coreDB "github.com/wiseco/core-platform/services/data"
	usrv "github.com/wiseco/core-platform/services/user"
	"github.com/wiseco/core-platform/shared"
	goLibClear "github.com/wiseco/go-lib/clear"
	"github.com/wiseco/go-lib/grpc"
	"github.com/wiseco/protobuf/golang"
	"github.com/wiseco/protobuf/golang/verification/alloy"
	"github.com/wiseco/protobuf/golang/verification/clear"
	"github.com/wiseco/protobuf/golang/verification/phone"
)

// Service  user/consumer services
type Service interface {
	UserList() ([]User, error)
	ByFilter(params map[string]interface{}) ([]User, error)
	ByID(id shared.ConsumerID) (User, error)
	ByUserID(id string) (User, error)
	UpdateID(id shared.ConsumerID, u *Update) (User, error)
	ByPhone(phone string) (User, error)
	UpdateKYC(id shared.ConsumerID, status string) (*usrv.Consumer, error)
	RunClearKYC(consumerID shared.ConsumerID) (string, error)
	GetClearKYC(consumerID shared.ConsumerID) (string, error)
	RunAlloyKYC(consumerID shared.ConsumerID) (string, error)
	GetAlloyKYC(consumerID shared.ConsumerID) (string, error)
	PhoneVerification(consumerID shared.ConsumerID) (string, error)
}

type service struct {
	sourceReq services.SourceRequest
}

// New a new user/consumer service
func New() Service {
	return service{}
}

// NewWithSource returns a service with a source request
func NewWithSource(sourceReq services.SourceRequest) Service {
	return service{sourceReq}
}

// List Lists all the users
func (s service) UserList() ([]User, error) {
	var list []User
	err := coreDB.DBRead.Select(&list, `
						SELECT
						wise_user.id, wise_user.consumer_id, wise_user.identity_id,
						wise_user.partner_id, wise_user.email, wise_user.email_verified,
						wise_user.phone, wise_user.phone_verified, wise_user.notification, wise_user.deactivated,
						wise_user.created, wise_user.modified, consumer.first_name, consumer.middle_name,
						consumer.last_name, consumer.date_of_birth, consumer.tax_id, consumer.tax_id_type,
						consumer.kyc_status, consumer.legal_address, consumer.mailing_address,
						consumer.work_address, consumer.residency, consumer.citizenship_countries,
						consumer.occupation, consumer.income_type, consumer.activity_type, consumer.is_restricted
						FROM
						wise_user
						JOIN
						consumer ON wise_user.consumer_id = consumer.id`)
	if err != nil && err == sql.ErrNoRows {
		return list, nil
	}

	return list, err
}

func (s service) ByFilter(params map[string]interface{}) ([]User, error) {
	clause := ""
	if params["firstName"] != nil && params["firstName"] != "" {
		fname := strings.ToLower(params["firstName"].(string))
		clause = "LOWER(c.first_name) LIKE '" + fname + "%'"

	} else if params["phone"] != nil && params["phone"] != "" {
		phone := params["phone"].(string)
		clause = "usr.phone = '" + phone + "'"

	} else if params["userId"] != nil && params["userId"] != "" {
		ID := params["userId"].(string)
		uID, err := shared.ParseUserID(ID)
		if err == nil {
			ID = string(uID)
		}
		clause = "usr.id = '" + ID + "'"
	} else if params["coId"] != nil && params["coId"] != "" {
		coID := params["coId"].(string)
		bank, err := partnerbank.GetProxyBank(partnerbank.ProviderNameBBVA)
		if err != nil {
			log.Printf("Error getting partner bank %v", err)
			return nil, err
		}

		sr := services.NewSourceRequest()
		coIDItem := partnerbank.ConsumerBankID(coID)
		consumerID, err := bank.ProxyService(sr.PartnerBankRequest()).GetConsumerID(coIDItem)
		clause = "usr.consumer_id = '" + string(*consumerID) + "'"
	}

	columns := `usr.id, usr.consumer_id, usr.identity_id, usr.partner_id, usr.email, usr.email_verified, 
				usr.phone, usr.phone_verified, usr.notification, usr.deactivated, usr.created, usr.modified, 
				c.first_name, c.middle_name, c.last_name, c.date_of_birth, c.tax_id, c.tax_id_type, c.kyc_status, c.legal_address, c.mailing_address,
				c.work_address, c.residency, c.citizenship_countries, c.occupation, c.income_type, c.activity_type, c.is_restricted`

	q := fmt.Sprintf("SELECT %v FROM wise_user usr JOIN consumer c ON usr.consumer_id = c.id WHERE %v ORDER BY c.first_name asc", columns, clause)

	var list []User
	err := coreDB.DBRead.Select(&list, q)

	if err != nil && err == sql.ErrNoRows {
		return list, nil
	}
	return list, err

}

// ByID Gets a user with a consumer id
func (s service) ByID(ConsumerID shared.ConsumerID) (User, error) {
	var usr User
	err := coreDB.DBRead.Get(&usr, `
						SELECT
						wise_user.id, wise_user.consumer_id, wise_user.identity_id,
						wise_user.partner_id, wise_user.email, wise_user.email_verified,
						wise_user.phone, wise_user.phone_verified, wise_user.notification, wise_user.deactivated,
						wise_user.created, wise_user.modified, consumer.first_name, consumer.middle_name,
						consumer.last_name, consumer.date_of_birth, consumer.tax_id, consumer.tax_id_type,
						consumer.kyc_status, consumer.legal_address, consumer.mailing_address,
						consumer.work_address, consumer.residency, consumer.citizenship_countries,
						consumer.occupation, consumer.income_type, consumer.activity_type, consumer.is_restricted
						FROM
						consumer
						JOIN
						wise_user ON consumer.id = wise_user.consumer_id
						WHERE consumer.id = $1`, ConsumerID)

	if err != nil && err == sql.ErrNoRows {
		return usr, services.ErrorNotFound{}.New("")
	}
	bank, err := partnerbank.GetProxyBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		log.Printf("Error getting partner bank %v", err)
	}
	CO, err := bank.ProxyService(s.sourceReq.PartnerBankRequest()).GetConsumerBankID(partnerbank.ConsumerID(ConsumerID))
	if err != nil && err != sql.ErrNoRows {
		log.Printf("error getting bankID %v", err)
	} else {
		usr.BankID = CO
	}
	return usr, err
}

func (s service) ByUserID(id string) (User, error) {
	var usr User
	err := coreDB.DBRead.Get(
		&usr, `
		SELECT
            wise_user.id, wise_user.consumer_id, wise_user.identity_id,
            wise_user.partner_id, wise_user.email, wise_user.email_verified, 
            wise_user.phone, wise_user.phone_verified, wise_user.notification, wise_user.deactivated,  wise_user.subscription_status,
            wise_user.created, wise_user.modified, consumer.first_name, consumer.middle_name,
            consumer.last_name, consumer.date_of_birth, consumer.tax_id, consumer.tax_id_type,
            consumer.kyc_status, consumer.legal_address, consumer.mailing_address,
            consumer.work_address, consumer.residency, consumer.citizenship_countries,
            consumer.occupation, consumer.income_type, consumer.activity_type, consumer.is_restricted
        FROM
            wise_user
        JOIN
            consumer ON wise_user.consumer_id = consumer.id
        WHERE
            wise_user.id = $1`,
		id,
	)
	if err != nil && err == sql.ErrNoRows {
		return usr, services.ErrorNotFound{}.New("")
	}
	return usr, err
}

func (s service) ByPhone(phone string) (User, error) {
	var usr User
	err := coreDB.DBRead.Get(
		&usr, `
		SELECT
            wise_user.id, wise_user.consumer_id, wise_user.identity_id,
            wise_user.partner_id, wise_user.email, wise_user.email_verified, 
            wise_user.phone, wise_user.phone_verified, wise_user.notification, wise_user.deactivated,
            wise_user.created, wise_user.modified, consumer.first_name, consumer.middle_name,
            consumer.last_name, consumer.date_of_birth, consumer.tax_id, consumer.tax_id_type,
            consumer.kyc_status, consumer.legal_address, consumer.mailing_address,
            consumer.work_address, consumer.residency, consumer.citizenship_countries,
            consumer.occupation, consumer.income_type, consumer.activity_type, consumer.is_restricted
        FROM
            wise_user
        JOIN
            consumer ON wise_user.consumer_id = consumer.id
        WHERE
            wise_user.phone = $1`,
		phone,
	)
	if err != nil && err == sql.ErrNoRows {
		return usr, services.ErrorNotFound{}.New("")
	}

	bank, err := partnerbank.GetProxyBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		log.Printf("Error getting partner bank %v", err)
	}
	CO, err := bank.ProxyService(s.sourceReq.PartnerBankRequest()).GetConsumerBankID(partnerbank.ConsumerID(usr.ConsumerID))
	if err != nil && err != sql.ErrNoRows {
		log.Printf("error getting bankID %v", err)
	} else {
		usr.BankID = CO
	}

	return usr, err
}

func (s service) UpdateID(id shared.ConsumerID, updates *Update) (User, error) {
	var kycStatus string
	var u User
	err := coreDB.DBRead.Get(&kycStatus, "SELECT kyc_status FROM consumer WHERE id = $1", id)
	if err != nil {
		log.Printf("Error getting consumer status %v", err)
		return User{}, err
	}

	if kycStatus == "approved" {
		dropValues(updates)
	}
	keys := services.SQLGenForUpdate(*updates)
	q := fmt.Sprintf("UPDATE consumer SET %s WHERE id = '%s' RETURNING *", keys, id)
	stmt, err := coreDB.DBWrite.PrepareNamed(q)
	if err != nil {
		log.Printf("error preparing stmt %v query %v", err, q)
		return u, err
	}
	err = stmt.Get(&u, updates)

	if updates.Email != nil {
		updateUserEmail(id, *updates.Email)
	}
	if kycStatus != "notStarted" {
		if err := s.updateConsumerOnBBVA(id, updates); err != nil {
			log.Printf("error updating consumer on BBVA")
		}
	}

	return u, err
}

func updateUserEmail(consumerID shared.ConsumerID, email string) error {
	var err error
	if email != "" {
		_, err = coreDB.DBWrite.Exec("UPDATE wise_user SET email = $1 WHERE consumer_id = $2", email, consumerID)
	}
	if err != nil {
		log.Printf("error update user email %v", err)
	}
	return err
}

// If KYC status is approved lets drop these values for update
func dropValues(u *Update) *Update {
	u.FirstName = nil
	u.MiddleName = nil
	u.LastName = nil
	u.DateOfBirth = nil
	u.TaxID = nil
	u.TaxIDType = nil
	u.Residency = nil
	u.CitizenshipCountries = nil
	return u
}

func validateAddress(u *Update) error {
	if u.LegalAddress != nil {
		addr, err := services.ValidateAddress(u.LegalAddress, services.AddressTypeLegal)
		if err != nil {
			return err
		}
		u.LegalAddress = addr
	}
	if u.MailingAddress != nil {
		addr, err := services.ValidateAddress(u.MailingAddress, services.AddressTypeMailing)
		if err != nil {
			return err
		}
		u.MailingAddress = addr
	}
	if u.WorkAddress != nil {
		addr, err := services.ValidateAddress(u.WorkAddress, services.AddressTypeWork)
		if err != nil {
			return err
		}
		u.WorkAddress = addr
	}
	return nil
}

func (s service) updateConsumerOnBBVA(id shared.ConsumerID, u *Update) error {
	bank, err := partnerbank.GetConsumerBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		return err
	}

	srv := bank.ConsumerEntityService(s.sourceReq.PartnerBankRequest())
	if u.Email != nil {
		err := srv.UpdateContact(
			partnerbank.ConsumerID(id),
			partnerbank.ConsumerPropertyTypeContactEmail,
			*u.Email,
		)
		if err != nil {
			return err
		}
	}

	if u.Phone != nil {
		err := srv.UpdateContact(
			partnerbank.ConsumerID(id),
			partnerbank.ConsumerPropertyTypeContactPhone,
			*u.Phone,
		)
		if err != nil {
			return err
		}
	}

	if u.LegalAddress != nil {
		err := srv.UpdateAddress(
			partnerbank.ConsumerID(id),
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
			return err
		}
	}

	if u.WorkAddress != nil {
		err := srv.UpdateAddress(
			partnerbank.ConsumerID(id),
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
			return err
		}
	}

	if u.MailingAddress != nil {
		err := srv.UpdateAddress(
			partnerbank.ConsumerID(id),
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
			return err
		}
	}

	return nil
}

// Update KYC Status
func (s service) UpdateKYC(id shared.ConsumerID, status string) (*usrv.Consumer, error) {
	var c usrv.Consumer
	err := coreDB.DBWrite.Get(&c, "UPDATE consumer SET kyc_status = $1 WHERE id = $2 RETURNING *", status, id)
	if err != nil {
		log.Printf("error updating consumer status %v", err)
		return nil, err
	}

	return &c, nil
}

// AIRSTREAM
func (s service) RunClearKYC(consumerID shared.ConsumerID) (string, error) {
	log.Println("In RunClearKYC")
	sn, err := grpc.GetConnectionStringForService(grpc.ServiceNameVerification)
	if err != nil {
		return "", err
	}
	client, err := grpc.NewInsecureClient(sn)
	if err != nil {
		return "", err
	}
	defer client.CloseAndCancel()
	c := clear.NewClearServiceClient(client.GetConn())

	consumer, err := usrv.NewConsumerServiceWithout().GetByID(consumerID)
	if err != nil {
		return "", err
	}
	ph := strings.Trim(*consumer.Phone, "+")
	if consumer.TaxID == nil {
		return "", errors.New("Tax id for member cannot be null")
	}
	vReq := &clear.ConsumerVerificationRequest{
		ConsumerId:  string(consumerID),
		FirstName:   consumer.FirstName,
		LastName:    consumer.LastName,
		Email:       *consumer.Email,
		DateOfBirth: consumer.DateOfBirth.String(),
		Address: &golang.AddressRequest{
			Line_1:     consumer.LegalAddress.StreetAddress,
			Locality:   consumer.LegalAddress.City,
			AdminArea:  consumer.LegalAddress.State,
			PostalCode: consumer.LegalAddress.PostalCode,
		},
		TaxId: string(*consumer.TaxID),
		Phone: ph,
	}

	res, err := c.RiskInformConsumerVerification(context.Background(), vReq)
	if err != nil {
		log.Println("Error occured while calling consumerVerification", err)
		return "", err
	}
	log.Println("Successfully completed clear consumerVerification")
	return string(res.ResultData), nil
}

func (s service) GetClearKYC(consumerID shared.ConsumerID) (string, error) {

	sn, err := grpc.GetConnectionStringForService(grpc.ServiceNameVerification)
	if err != nil {
		return "", err
	}
	client, err := grpc.NewInsecureClient(sn)
	if err != nil {
		return "", err
	}
	defer client.CloseAndCancel()
	c := clear.NewClearServiceClient(client.GetConn())

	gReq := &clear.GetConsumerRequest{
		ConsumerId:   string(consumerID),
		ClearKycType: string(goLibClear.RiskInformPerson),
	}
	res, err := c.GetConsumer(client.GetContext(), gReq)
	if err != nil {
		return "", err
	}

	return string(res.ResultData), nil
}

func (s service) RunAlloyKYC(consumerID shared.ConsumerID) (string, error) {
	sn, err := grpc.GetConnectionStringForService(grpc.ServiceNameVerification)
	if err != nil {
		return "", err
	}
	client, err := grpc.NewInsecureClient(sn)
	if err != nil {
		return "", err
	}
	defer client.CloseAndCancel()
	c := alloy.NewAlloyServiceClient(client.GetConn())

	consumer, err := usrv.NewConsumerServiceWithout().GetByID(consumerID)
	if err != nil {
		return "", err
	}
	ph := strings.Trim(*consumer.Phone, "+")
	if consumer.TaxID == nil {
		return "", errors.New("Tax id for member cannot be null")
	}
	vReq := &alloy.ConsumerVerificationRequest{
		ConsumerId:    string(consumer.ID),
		FirstName:     consumer.FirstName,
		LastName:      consumer.LastName,
		Email:         *consumer.Email,
		Dob:           consumer.DateOfBirth.String(),
		AddressLine_1: consumer.LegalAddress.StreetAddress,
		City:          consumer.LegalAddress.City,
		State:         consumer.LegalAddress.State,
		PostalCode:    consumer.LegalAddress.PostalCode,
		Country:       consumer.LegalAddress.Country,
		Ssn:           string(*consumer.TaxID),
		Phone:         ph,
	}

	res, err := c.ConsumerVerification(context.Background(), vReq)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(res.ResultData), nil
}

func (s service) GetAlloyKYC(consumerID shared.ConsumerID) (string, error) {

	sn, err := grpc.GetConnectionStringForService(grpc.ServiceNameVerification)
	if err != nil {
		return "", err
	}
	client, err := grpc.NewInsecureClient(sn)
	if err != nil {
		return "", err
	}
	defer client.CloseAndCancel()
	c := alloy.NewAlloyServiceClient(client.GetConn())

	gReq := &alloy.GetConsumerRequest{
		ConsumerId: string(consumerID),
	}

	res, err := c.GetConsumer(client.GetContext(), gReq)
	if err != nil {
		return "", err
	}

	return string(res.ResultData), nil
}

func (s service) PhoneVerification(conusmerID shared.ConsumerID) (string, error) {
	sn, err := grpc.GetConnectionStringForService(grpc.ServiceNameVerification)
	if err != nil {
		return "", err
	}
	client, err := grpc.NewInsecureClient(sn)
	if err != nil {
		return "", err
	}
	defer client.CloseAndCancel()
	c := phone.NewPhoneServiceClient(client.GetConn())

	consumer, err := usrv.NewConsumerServiceWithout().GetByID(conusmerID)
	if err != nil {
		return "", err
	}

	vReq := &phone.VerificationRequest{
		ConsumerId:  string(conusmerID),
		PhoneNumber: *consumer.Phone,
	}
	res, err := c.Verification(context.Background(), vReq)

	if err != nil {
		return "", err
	}
	return res.Raw, nil
}
