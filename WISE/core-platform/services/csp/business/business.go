/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package business

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/partner/service/segment"
	"github.com/wiseco/core-platform/services"
	busBanking "github.com/wiseco/core-platform/services/banking/business"
	bus "github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/services/csp"
	cspData "github.com/wiseco/core-platform/services/csp/data"
	"github.com/wiseco/core-platform/services/csp/document"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/go-lib/grpc"
	"github.com/wiseco/protobuf/golang"
	"github.com/wiseco/protobuf/golang/verification/alloy"
	"github.com/wiseco/protobuf/golang/verification/business"
	grpcClear "github.com/wiseco/protobuf/golang/verification/clear"
)

//Service the csp business services
type Service interface {

	//Business By ID
	ByID(shared.BusinessID) (Business, error)

	//ListAll  gets all the businesses with different statuses
	ListAll(limit, offset int) ([]*Business, error)

	UpdateID(shared.BusinessID, BusinessUpdate) (Business, error)

	ListStatus(status string, limit, offset int) ([]*Business, error)

	UpdateKYC(businessID shared.BusinessID, status csp.KYCStatus) error

	RunMiddesk(shared.BusinessID) (string, error)

	GetMiddesk(shared.BusinessID) (string, error)

	RunKYB(string) (*KYBResult, error)

	RunClearVerification(shared.BusinessID) (string, error)

	GetClearVerification(shared.BusinessID) (string, error)

	GetActiveBusinessIDsInternal() ([]string, error)
}

type service struct {
	*sqlx.DB
	sourceReq services.SourceRequest
}

//New will return a new csp business service
func New(r services.SourceRequest) Service {
	return service{data.DBWrite, r}
}

func (s service) ByID(id shared.BusinessID) (Business, error) {
	var biz Business

	err := s.Get(&biz, `
	SELECT * FROM business
	 WHERE id = $1`, id)
	if err != nil {
		return biz, err
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bas, err := busBanking.NewBankingAccountService()
		if err != nil {
			return biz, err
		}

		var account *busBanking.BankAccount
		accounts, err := bas.GetByBusinessID(id, 10, 0)
		if err == sql.ErrNoRows {
			return biz, nil
		}

		if err != nil {
			return biz, err
		}

		for _, acc := range accounts {
			if acc.UsageType == busBanking.UsageTypePrimary {
				account = &acc
				break
			}
		}

		if account != nil {
			biz.AvailableBalance = &account.AvailableBalance
			biz.PostedBalance = &account.PostedBalance
		}
	} else {
		err = s.Get(&biz, `
		SELECT posted_balance, available_balance FROM 
		business_bank_account 
		WHERE business_id = $1`, id)
		if err != nil {
			return biz, err
		}
	}

	bank, err := partnerbank.GetProxyBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		log.Printf("Error getting partner bank %v", err)
	}
	NC, err := bank.ProxyService(s.sourceReq.PartnerBankRequest()).GetBusinessBankID(partnerbank.BusinessID(id))
	if err != nil && err != sql.ErrNoRows {
		log.Printf("error getting bankID %v", err)
	} else {
		biz.BankID = NC
	}

	return biz, nil
}

func (s service) ListAll(limit, offset int) ([]*Business, error) {
	var list []*Business
	err := s.Select(&list, "SELECT * FROM business ORDER BY created DESC LIMIT $1 OFFSET $2", limit, offset)
	if err == sql.ErrNoRows {
		log.Printf("no business %v", err)
		return []*Business{}, nil
	}

	return list, err
}

func (s service) ListStatus(status string, limit, offset int) ([]*Business, error) {
	list := make([]*Business, 0)
	err := s.Select(&list, "SELECT * FROM business WHERE kyc_status = $1 ORDER BY updated DESC LIMIT $2 OFFSET $3", status, limit, offset)
	if err == sql.ErrNoRows {
		return []*Business{}, nil
	}
	return list, err
}

func (s service) UpdateID(id shared.BusinessID, updates BusinessUpdate) (Business, error) {
	var biz Business
	keys := services.SQLGenForUpdate(updates)
	q := fmt.Sprintf("UPDATE business SET %s WHERE id = '%s' RETURNING *", keys, id)
	stmt, err := s.PrepareNamed(q)
	if err != nil {
		log.Printf("error preparing stmt %v query %v", err, q)
		return biz, err
	}
	err = stmt.Get(&biz, updates)
	if err != nil {
		return biz, fmt.Errorf("error keys: %v err: %v", keys, err)
	}

	if updates.LegalName != nil || updates.EntityType != nil {
		businessName := ""
		if updates.LegalName != nil {
			businessName = *updates.LegalName
		}

		businessEntityType := ""
		if updates.EntityType != nil {
			businessEntityType = *updates.EntityType
		}
		NewCSPService().UpdateBusinessNameAndEntityType(id, businessName, businessEntityType)
	}

	// Update bank data
	if biz.KYCStatus == services.KYCStatusReview {
		// Update address and contacts in BBVA
		bank, err := partnerbank.GetBusinessBank(partnerbank.ProviderNameBBVA)
		if err != nil {
			return biz, err
		}

		srv := bank.BusinessEntityService(s.sourceReq.PartnerBankRequest())

		ureq := partnerbank.UpdateBusinessRequest{
			BusinessID: partnerbank.BusinessID(biz.ID),
			LegalName:  shared.StringValue(biz.LegalName),
			DBA:        biz.DBA,
		}

		if biz.EntityType != nil {
			ureq.EntityType = partnerbank.BusinessEntity(*biz.EntityType)
		}

		if biz.IndustryType != nil {
			ureq.IndustryType = partnerbank.BusinessIndustry(*biz.IndustryType)
		}

		if biz.Purpose != nil {
			ureq.Purpose = *biz.Purpose
		}

		if biz.OperationType != nil {
			ureq.OperationType = partnerbank.BusinessOperationType(*biz.OperationType)
		}

		if biz.TaxIDType != nil && biz.TaxID != nil {
			ureq.TaxIDType = partnerbank.BusinessTaxIDType(*biz.TaxIDType)
			taxID := *biz.TaxID
			ureq.TaxID = string(taxID)
		}

		if biz.OriginCountry != nil {
			ureq.OriginCountry = partnerbank.Country(*biz.OriginCountry)
		}

		if biz.OriginState != nil {
			ureq.OriginState = *biz.OriginState
		}

		if biz.OriginDate != nil {
			ureq.OriginDate = biz.OriginDate.Time()
		}

		// Get entity formation do
		if biz.FormationDocumentID != nil {
			doc, err := document.NewDocumentService().GetByID(biz.ID, *biz.FormationDocumentID)
			if err != nil {
				return biz, err
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
		}

		// Update business on bank side - log errors
		var wait sync.WaitGroup
		wait.Add(1)
		go func() {
			defer wait.Done()
			_, err = srv.Update(ureq)
		}()

		// Update only on change
		if biz.Email != nil {
			wait.Add(1)
			go func() {
				defer wait.Done()
				err := srv.UpdateContact(
					partnerbank.BusinessID(biz.ID),
					partnerbank.BusinessPropertyTypeContactEmail,
					*biz.Email,
				)
				if err != nil {
					// Log errors
					log.Println(err)
				}
			}()
		}

		if biz.Phone != nil {
			wait.Add(1)
			go func() {
				defer wait.Done()

				err := srv.UpdateContact(
					partnerbank.BusinessID(biz.ID),
					partnerbank.BusinessPropertyTypeContactPhone,
					*biz.Phone,
				)
				if err != nil {
					// Log errors
					log.Println(err)
				}
			}()
		}

		if biz.LegalAddress != nil {
			wait.Add(1)
			go func() {
				defer wait.Done()
				err := srv.UpdateAddress(
					partnerbank.BusinessID(biz.ID),
					partnerbank.BusinessPropertyTypeAddressLegal,
					biz.LegalAddress.ToPartnerBankAddress(services.AddressTypeLegal),
				)
				if err != nil {
					// Log errors
					log.Println(err)
				}
			}()
		}

		if biz.HeadquarterAddress != nil {
			wait.Add(1)
			go func() {
				defer wait.Done()
				err := srv.UpdateAddress(
					partnerbank.BusinessID(biz.ID),
					partnerbank.BusinessPropertyTypeAddressHeadquarter,
					biz.HeadquarterAddress.ToPartnerBankAddress(services.AddressTypeHeadquarter),
				)
				if err != nil {
					// Log errors
					log.Println(err)
				}
			}()
		}

		if biz.MailingAddress != nil {
			wait.Add(1)
			go func() {
				defer wait.Done()
				err := srv.UpdateAddress(
					partnerbank.BusinessID(biz.ID),
					partnerbank.BusinessPropertyTypeAddressMailing,
					biz.MailingAddress.ToPartnerBankAddress(services.AddressTypeMailing),
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

	segment.NewSegmentService().PushToAnalyticsQueue(biz.OwnerID, segment.CategoryBusiness, segment.ActionUpdate, biz)
	return biz, err
}

type businessVerificationPreReq struct {
	bus          Business
	client       grpc.Client
	businessName string
}

func (s service) initBusinessVerificationPreReq(id shared.BusinessID) (businessVerificationPreReq, error) {
	res := businessVerificationPreReq{}
	var b Business
	var err error
	var client grpc.Client

	b, err = s.ByID(id)
	if err != nil {
		return res, err
	}

	sn, err := grpc.GetConnectionStringForService(grpc.ServiceNameVerification)
	if err != nil {
		return res, err
	}

	client, err = grpc.NewInsecureClient(sn)
	if err != nil {
		return res, err
	}

	if b.Name() == "" {
		err = errors.New("Business name or DBA is required")
		return res, err
	}

	var businessName string
	switch *b.EntityType {
	case bus.EntityTypeSoleProprietor:
		businessName = b.Name()
	default:
		if b.LegalName != nil {
			businessName = *b.LegalName
		} else {
			businessName = b.Name()
		}
	}

	if b.Phone == nil {
		err = errors.New("Phone number is required")
		return res, err
	}
	if b.TaxID == nil {
		err = errors.New("TaxId is required")
		return res, err
	}

	res.bus = b
	res.businessName = businessName
	res.client = client

	return res, err
}

func (s service) UpdateKYC(businessID shared.BusinessID, status csp.KYCStatus) error {
	_, err := s.Exec("UPDATE business SET kyc_status = $2 WHERE id = $1", businessID, status)
	return err
}

func (s service) RunMiddesk(id shared.BusinessID) (string, error) {
	bv, err := s.initBusinessVerificationPreReq(id)
	if err != nil {
		return "", err
	}

	b := bv.bus
	defer bv.client.CloseAndCancel()

	service := business.NewBusinessServiceClient(bv.client.GetConn())
	vReq := &business.VerificationRequest{
		BusinessId: string(b.ID),
		Name:       bv.businessName,
		Tin:        *b.TaxID,
		Addresses: []*business.AddressRequest{
			&business.AddressRequest{
				Line_1:     b.LegalAddress.StreetAddress,
				Line_2:     b.LegalAddress.AddressLine2,
				Locality:   b.LegalAddress.City,
				AdminArea:  b.LegalAddress.State,
				PostalCode: b.LegalAddress.PostalCode,
			},
		},
		PhoneNumbers: []string{*b.Phone},
	}

	res, err := service.Verification(bv.client.GetContext(), vReq)
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", errors.New("Invalid response")
	}

	return string(res.RawResults), nil
}

func (s service) RunClearVerification(id shared.BusinessID) (string, error) {
	bv, err := s.initBusinessVerificationPreReq(id)
	if err != nil {
		return "", err
	}

	b := bv.bus

	defer bv.client.CloseAndCancel()

	ph := strings.Trim(*b.Phone, "+")
	service := grpcClear.NewClearServiceClient(bv.client.GetConn())
	vReq := &grpcClear.BusinessVerificationRequest{
		BusinessId: string(b.ID),
		Name:       bv.businessName,
		Phone:      ph,
		TaxId:      *b.TaxID,
		Address: &golang.AddressRequest{
			Line_1:     b.LegalAddress.StreetAddress,
			Locality:   b.LegalAddress.City,
			AdminArea:  b.LegalAddress.State,
			PostalCode: b.LegalAddress.PostalCode,
		},
	}

	res, err := service.RiskInformBusinessVerification(bv.client.GetContext(), vReq)
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", errors.New("Invalid response")
	}

	return string(res.ResultData), nil
}

func (s service) GetMiddesk(id shared.BusinessID) (string, error) {
	b, err := s.ByID(id)
	if err != nil {
		return "", err
	}

	sn, err := grpc.GetConnectionStringForService(grpc.ServiceNameVerification)
	if err != nil {
		return "", err
	}

	client, err := grpc.NewInsecureClient(sn)
	if err != nil {
		return "", err
	}

	defer client.CloseAndCancel()
	service := business.NewBusinessServiceClient(client.GetConn())
	vReq := &business.GetBusinessRequest{
		BusinessId: string(b.ID),
	}

	res, err := service.GetBusiness(client.GetContext(), vReq)
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", errors.New("Invalid response")
	}

	return string(res.RawResults), nil

}

func (s service) GetClearVerification(id shared.BusinessID) (string, error) {
	b, err := s.ByID(id)
	if err != nil {
		return "", err
	}

	sn, err := grpc.GetConnectionStringForService(grpc.ServiceNameVerification)
	if err != nil {
		return "", err
	}

	client, err := grpc.NewInsecureClient(sn)
	if err != nil {
		return "", err
	}

	defer client.CloseAndCancel()
	service := grpcClear.NewClearServiceClient(client.GetConn())

	vReq := &grpcClear.GetBusinessRequest{
		BusinessId: string(b.ID),
	}

	res, err := service.GetBusiness(client.GetContext(), vReq)
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", errors.New("Invalid response")
	}

	return string(res.ResultData), nil
}

func (s service) RunKYB(tin string) (*KYBResult, error) {
	b, err := s.GetByTIN(tin)
	if err != nil {
		return nil, err
	}

	var businessName string
	switch *b.EntityType {
	case "soleProprietor":
		businessName = shared.GetBusinessName(b.LegalName, b.DBA)
	default:
		if b.LegalName == nil {
			return nil, errors.New("Business legal name is required")
		}

		businessName = *b.LegalName
	}

	if len(businessName) == 0 {
		return nil, errors.New("Business legal name or DBA is required")
	}

	var wg sync.WaitGroup

	result := KYBResult{}

	var kybSummary *KYBSummary
	var alloyScore *float64
	var middeskScore *float64

	// Alloy KYB
	wg.Add(1)
	go func() {
		defer wg.Done()
		alloyScore, err = s.runAlloyKYB(b)
	}()

	// Middesk KYB
	wg.Add(1)
	go func() {
		defer wg.Done()
		middeskScore, kybSummary, err = s.runMiddeskKYB(b)
	}()

	wg.Wait()

	if err != nil || kybSummary == nil || alloyScore == nil || middeskScore == nil {
		return nil, errors.New("Unable to complete KYB process")
	}

	if kybSummary.Watchlist != WatchListStatusFound {
		result.WiseScore = *alloyScore + *middeskScore
	}

	switch {
	case result.WiseScore < 0:
		result.Result = KYCStatusDeclined
	case result.WiseScore <= 50:
		result.Result = KYCStatusReview
	case result.WiseScore <= 100:
		result.Result = KYCStatusApproved
	default:
		result.Result = KYCStatusUnKnown
	}

	result.KYBSummary = *kybSummary
	return &result, nil
}

func (s service) GetByTIN(tin string) (*Business, error) {
	var b Business

	err := s.Get(&b, `SELECT * FROM business WHERE tax_id = $1`, tin)
	if err != nil {
		return nil, err
	}

	return &b, nil
}

func (s service) runAlloyKYB(b *Business) (*float64, error) {
	sn, err := grpc.GetConnectionStringForService(grpc.ServiceNameVerification)
	if err != nil {
		return nil, err
	}
	client, err := grpc.NewInsecureClient(sn)
	if err != nil {
		return nil, err
	}
	defer client.CloseAndCancel()
	c := alloy.NewAlloyServiceClient(client.GetConn())

	ph := strings.Trim(*b.Phone, "+")

	vReq := &alloy.BusinessVerificationRequest{
		Name:          *b.LegalName,
		Phone:         ph,
		AddressLine_1: b.LegalAddress.StreetAddress,
		City:          b.LegalAddress.City,
		State:         b.LegalAddress.State,
		PostalCode:    b.LegalAddress.PostalCode,
		Country:       b.LegalAddress.Country,
		Tin:           *b.TaxID,
	}

	res, err := c.BusinessVerification(client.GetContext(), vReq)
	if err != nil {
		return nil, err
	}

	score := (res.Score * 100) / 2
	return &score, nil
}

func (s service) runMiddeskKYB(b *Business) (*float64, *KYBSummary, error) {
	var summary KYBSummary
	var m MiddeskResponse
	resp, err := s.GetMiddesk(b.ID)
	if err != nil {
		return nil, nil, err
	}

	err = json.Unmarshal([]byte(resp), &m)
	if err != nil {
		return nil, nil, err
	}

	var middeskScore float64

	var yrsInBusiness int
	var tinMatch int
	var registrationMatch int
	var watchListMatch int
	var addressMatch int

	if m.Formation != nil {
		layout := "2006-01-02"
		date := m.Formation.FormationDate

		formationDate, err := time.Parse(layout, date)
		if err != nil {
			summary.InBusiness = VerificationStatusVerified
		} else {
			yrsInBusiness = (time.Now().Year() - formationDate.Year()) + 1

			if yrsInBusiness > MaxYearsInBusinessWeightage {
				yrsInBusiness = MaxYearsInBusinessWeightage
			}

			summary.InBusiness = VerificationStatusVerified
		}
	} else {
		summary.InBusiness = VerificationStatusUnverified
	}

	for _, sum := range m.Summary {
		if sum.Name == "Tin" {
			switch sum.Status {
			case SummaryStatusSuccess:
				tinMatch = TinMatchWeightage
				summary.Tin = VerificationStatusVerified
			default:
				summary.Tin = VerificationStatusUnverified
			}
		}

		if sum.Name == "Sos" {
			switch sum.Status {
			case SummaryStatusSuccess:
				registrationMatch = FormationMatchWeightage
				summary.Formation = VerificationStatusVerified
			default:
				summary.Formation = VerificationStatusUnverified
			}
		}

		if sum.Name == "Address" {
			switch sum.Status {
			case SummaryStatusSuccess:
				addressMatch = AddressMatchWeightage
				summary.Address = VerificationStatusVerified
			default:
				summary.Address = VerificationStatusUnverified
			}
		}

		if sum.Name == "Watchlist" {
			switch sum.Status {
			case SummaryStatusFailure:
				summary.Watchlist = WatchListStatusFound
			case SummaryStatusSuccess:
				watchListMatch = WatchlistNotFoundWeightage
				summary.Watchlist = WatchListStatusNotFound
			default:
				summary.Watchlist = WatchListStatusUnKnown
			}
		}
	}

	middeskScore = float64(yrsInBusiness + tinMatch + registrationMatch + watchListMatch + addressMatch)
	return &middeskScore, &summary, nil
}

// For tools/Projectlock
func (s service) GetActiveBusinessIDsInternal() ([]string, error) {
	var businessIds = []string{}
	cspDb := cspData.DBRead
	err := cspDb.Select(&businessIds, "SELECT business_id FROM business WHERE review_status IN ('training','trainingComplete') ORDER BY created")
	return businessIds, err
}
