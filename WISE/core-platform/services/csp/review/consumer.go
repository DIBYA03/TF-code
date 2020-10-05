package review

import (
	"errors"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/csp/consumer"
	coreDB "github.com/wiseco/core-platform/services/data"
	usrv "github.com/wiseco/core-platform/services/user"
	"github.com/wiseco/core-platform/shared"
)

//ConsumerVerification ..
type ConsumerVerification interface {

	// User verification (KYC)
	StartVerification(shared.ConsumerID) (*ConsumerKYCResponse, error)

	//Continue Verification
	GetVerification(shared.ConsumerID) (*ConsumerKYCResponse, error)

	//Continue verification
	ContinueVerification(*usrv.Consumer) (*ConsumerKYCResponse, error)
}

type consumerVerification struct {
	sourceReq services.SourceRequest
	*sqlx.DB
}

func (v consumerVerification) GetUser(id string) (usrv.User, error) {
	var user usrv.User
	err := v.Get(&user, "SELECT * from wise_user WHERE id = $1", id)
	return user, err
}

//NewConsumerVerfication ..
func NewConsumerVerfication(s services.SourceRequest) ConsumerVerification {
	return consumerVerification{s, coreDB.DBWrite}
}

func (v consumerVerification) StartVerification(consumerID shared.ConsumerID) (*ConsumerKYCResponse, error) {
	resp, err := v.startConsumerVerification(consumerID)
	if err != nil {
		return nil, err
	}

	return &ConsumerKYCResponse{
		Status:      resp.Status,
		ReviewItems: resp.ReviewItems,
	}, nil
}

func getConsumer(id shared.ConsumerID) (*usrv.Consumer, error) {
	var consumer usrv.Consumer
	if err := coreDB.DBWrite.Get(&consumer, "SELECT * FROM consumer WHERE id = $1", id); err != nil {
		log.Printf("error getting consumer by id %v", err)
		return nil, err
	}
	return &consumer, nil
}

func getConsumerUserID(id string) shared.ConsumerID {
	var consumerID shared.ConsumerID
	err := coreDB.DBWrite.Get(&consumerID, "SELECT consumer_id FROM wise_user WHERE id = $1", id)
	if err != nil {
		log.Printf("error getting consumer by user id %v", err)
	}
	return consumerID
}

func (v consumerVerification) startConsumerVerification(id shared.ConsumerID) (*ConsumerKYCResponse, error) {
	c, err := getConsumer(id)
	if err != nil {
		return nil, ConsumerKYCError{
			Raw:       err,
			ErrorType: services.KYCErrorTypeOther,
		}
	}

	if err := checkConsumerStatus(c); err != nil {
		return nil, err
	}

	if c.KYCStatus == services.KYCStatusReview {
		return v.ContinueVerification(c)
	}
	//Check if status is not started
	if c.KYCStatus != services.KYCStatusNotStarted {
		return nil, ConsumerKYCError{
			ErrorType: services.KYCErrorTypeInProgress,
		}
	}

	if len(checkConsumerCommons(c)) > 0 {
		return nil, ConsumerKYCError{
			ErrorType: services.KYCErrorTypeParam,
			Values:    checkConsumerCommons(c),
		}
	}

	bank, err := partnerbank.GetConsumerBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		return nil, ConsumerKYCError{
			Raw:       err,
			ErrorType: services.KYCErrorTypeOther,
		}
	}

	r := v.sourceReq.PartnerBankRequest()

	creq := v.buildConsumerRequest(c)

	if c.MailingAddress != nil {
		mailingAddress := c.MailingAddress.ToPartnerBankAddress(services.AddressTypeMailing)
		creq.MailingAddress = &mailingAddress
	}

	resp, err := bank.ConsumerEntityService(r).Create(creq)
	if err != nil {
		return nil, ConsumerKYCError{
			Raw:       err,
			ErrorType: services.KYCErrorTypeOther,
		}
	}
	// Update in consumer status
	c, err = consumer.New().UpdateKYC(c.ID, resp.KYC.Status.String())
	if err != nil {
		return nil, ConsumerKYCError{
			Raw:       err,
			ErrorType: services.KYCErrorTypeOther,
		}
	}

	// Return KYC response
	var items services.StringArray
	for _, item := range resp.KYC.IDVerifyRequired {
		items = append(items, string(item))
	}

	submitted := time.Now()
	cspcs := consumer.NewCSPService()
	cspConsumer, err := cspcs.ByConsumerID(string(c.ID))
	if cspConsumer == nil {
		cspcs.CSPConsumerCreate(consumer.CSPConsumerCreate{
			ConsumerName: &c.FirstName,
			ConsumerID:   c.ID,
			IDVs:         &items,
			Submitted:    &submitted,
			Status:       resp.KYC.Status.String(),
		})
	}

	return &ConsumerKYCResponse{
		Status:      string(resp.KYC.Status),
		ReviewItems: &items,
	}, nil
}

func (v consumerVerification) buildConsumerRequest(c *usrv.Consumer) partnerbank.CreateConsumerRequest {
	return partnerbank.CreateConsumerRequest{
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

}

func (v consumerVerification) GetVerification(id shared.ConsumerID) (*ConsumerKYCResponse, error) {
	c, err := getConsumer(id)
	if err != nil {
		log.Printf("Error getting consumer %v", err)
		return nil, ConsumerKYCError{
			Raw:       err,
			ErrorType: services.KYCErrorTypeOther,
		}
	}

	bank, err := partnerbank.GetConsumerBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		log.Printf("Error getting partner bank %v", err)
		return nil, ConsumerKYCError{
			Raw:       err,
			ErrorType: services.KYCErrorTypeOther,
		}
	}

	resp, err := bank.ConsumerEntityService(v.sourceReq.PartnerBankRequest()).Status(partnerbank.ConsumerID(id))
	if err != nil {
		log.Printf("Error getting consumer status %v", err)
		return nil, ConsumerKYCError{
			Raw:       err,
			ErrorType: services.KYCErrorTypeOther,
		}
	}

	c, err = consumer.New().UpdateKYC(c.ID, resp.KYC.Status.String())
	if err != nil {
		return nil, ConsumerKYCError{
			Raw:       err,
			ErrorType: services.KYCErrorTypeOther,
		}
	}

	var items services.StringArray
	for _, item := range resp.KYC.IDVerifyRequired {
		items = append(items, string(item))
	}

	st := resp.KYC.Status.String()
	consumer.NewCSPService().CSPConsumerUpdateByConsumerID(c.ID, consumer.CSPConsumerUpdate{
		Status: &st,
		IDVs:   &items,
	})

	return &ConsumerKYCResponse{
		Status:      string(resp.KYC.Status),
		ReviewItems: &items,
	}, nil
}

func (v consumerVerification) ContinueVerification(c *usrv.Consumer) (*ConsumerKYCResponse, error) {
	if len(checkConsumerCommons(c)) > 0 {
		return nil, ConsumerKYCError{
			ErrorType: services.KYCErrorTypeParam,
			Values:    checkConsumerCommons(c),
		}
	}

	bank, err := partnerbank.GetConsumerBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		return nil, ConsumerKYCError{
			Raw:       err,
			ErrorType: services.KYCErrorTypeOther,
		}
	}

	r := v.sourceReq.PartnerBankRequest()

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
			Raw:       err,
			ErrorType: services.KYCErrorTypeOther,
		}
	}

	// Update core consumer status
	c, err = consumer.New().UpdateKYC(c.ID, resp.KYC.Status.String())
	if err != nil {
		return nil, ConsumerKYCError{
			Raw:       err,
			ErrorType: services.KYCErrorTypeOther,
		}
	}

	//Update csp consumer status
	consumer.NewCSPService().UpdateStatus(c.ID, resp.KYC.Status.String())
	// Return KYC response
	var items services.StringArray
	for _, item := range resp.KYC.IDVerifyRequired {
		items = append(items, string(item))
	}

	return &ConsumerKYCResponse{
		Status:      string(resp.KYC.Status),
		ReviewItems: &items,
	}, nil
}

//Check status
func checkConsumerStatus(c *usrv.Consumer) error {
	if c.Deactivated != nil {
		return ConsumerKYCError{
			ErrorType: services.KYCErrorTypeDeactivated,
		}
	}
	if c.IsRestricted {
		return ConsumerKYCError{
			ErrorType: services.KYCErrorTypeRestricted,
		}
	}
	switch c.KYCStatus {
	case services.KYCStatusApproved:
		return ConsumerKYCError{
			Raw:       errors.New("Consumer has already been approved"),
			ErrorType: services.KYCErrorTypeOther,
		}
	case services.KYCStatusDeclined:
		return ConsumerKYCError{
			Raw:       errors.New("Consumer has already been declined"),
			ErrorType: services.KYCErrorTypeOther,
		}
	}
	return nil

}

func consumerIDVS(consumerID shared.ConsumerID) []partnerbank.IDVerify {
	c, err := consumer.NewCSPService().CSPConsumerByConsumerID(consumerID)
	var idvs []partnerbank.IDVerify
	if err != nil {
		return []partnerbank.IDVerify{}
	}
	if c.IDVs != nil {
		for _, idv := range *c.IDVs {
			idvs = append(idvs, reviewToIDV(idv))
		}
	}
	return idvs
}

func reviewToIDV(reviewIDV string) partnerbank.IDVerify {
	switch reviewIDV {
	case partnerbank.IDVerifyAddress.String():
		return partnerbank.IDVerifyAddress
	case partnerbank.IDVerifyDateOfBirth.String():
		return partnerbank.IDVerifyDateOfBirth
	case partnerbank.IDVerifyFullName.String():
		return partnerbank.IDVerifyFullName
	case partnerbank.IDVerifyTaxId.String():
		return partnerbank.IDVerifyTaxId
	case partnerbank.IDVerifyMismatch.String():
		return partnerbank.IDVerifyMismatch
	case partnerbank.IDVerifyOFAC.String():
		return partnerbank.IDVerifyOFAC
	case partnerbank.IDVerifyOther.String():
		return partnerbank.IDVerifyOther
	case partnerbank.IDVerifyFormationDoc.String():
		return partnerbank.IDVerifyFormationDoc
	case partnerbank.IDVerifyPrimaryDoc.String():
		return partnerbank.IDVerifyPrimaryDoc
	case partnerbank.IDVerifySecondaryDoc.String():
		return partnerbank.IDVerifySecondaryDoc
	}
	return partnerbank.IDVerifyOther
}
