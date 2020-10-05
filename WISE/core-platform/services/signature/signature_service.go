package signature

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/wiseco/go-lib/hellosign"
	"github.com/wiseco/go-lib/log"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
)

type signatureDatastore struct {
	sourceReq services.SourceRequest
	*sqlx.DB
}

type SignatureService interface {
	Create(shared.BusinessID, SignatureRequestTemplate) (*SignatureRequest, error)

	GetByBusinessID(shared.BusinessID, SignatureRequestTemplate) (*SignatureRequest, error)
	GetBySignatureRequestID(string) (*SignatureRequest, error)

	UpdateSignatureStatus(shared.SignatureRequestID, shared.BusinessDocumentID, SignatureRequestStatus) (*SignatureRequest, error)
}

func NewSignatureService(r services.SourceRequest) SignatureService {
	return &signatureDatastore{r, data.DBWrite}
}

func (s *signatureDatastore) Create(businessID shared.BusinessID, templateType SignatureRequestTemplate) (*SignatureRequest, error) {
	l := log.NewLogger()
	_, ok := SignatureRequestTemplateTo[templateType]
	if !ok {
		return nil, errors.New("Invalid template type")
	}

	req, err := s.getSignatureRequest(businessID)
	if err != nil {
		l.Error(err.Error())
		return nil, err
	}

	signatureRequest := constructSignatureRequest(req, templateType)

	resp, err := hellosign.NewHellosignService(l).CreateSignatureRequest(signatureRequest)
	if err != nil {
		return nil, err
	}

	// Store in DB
	signatureCreate := SignatureRequestCreate{
		TemplateType:       templateType,
		TemplateProvider:   SignatureRequestProviderHellosign,
		BusinessID:         businessID,
		SignatureRequestID: resp.SignatureRequest.SignatureRequestID,
		SignatureStatus:    SignatureRequestStatusPending,
	}

	if len(resp.SignatureRequest.Signatures) > 0 {
		signatureCreate.SignatureID = resp.SignatureRequest.Signatures[0].SignatureID
	}

	// Default/mandatory fields
	columns := []string{
		"business_id", "template_type", "template_provider", "signature_request_id", "signature_id", "signature_status",
	}
	// Default/mandatory values
	values := []string{
		":business_id", ":template_type", ":template_provider", ":signature_request_id", ":signature_id", ":signature_status",
	}

	sql := fmt.Sprintf(
		"INSERT INTO signature_request(%s) VALUES(%s) RETURNING *",
		strings.Join(columns, ", "),
		strings.Join(values, ", "),
	)

	stmt, err := s.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	sign := &SignatureRequest{}
	err = stmt.Get(sign, &signatureCreate)
	if err != nil {
		return nil, err
	}

	sign.SignURL = resp.SignatureRequest.SigningURL

	return sign, nil
}

func (s *signatureDatastore) GetByBusinessID(businessID shared.BusinessID, templateType SignatureRequestTemplate) (*SignatureRequest, error) {
	lo := log.NewLogger()

	// If no entry return nil
	signature, err := s.getByBusinessID(businessID, templateType)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	// If document is completed do not get sign url
	if signature.SignatureStatus == SignatureRequestStatusCompleted {
		return signature, nil
	}

	// If entry, check business and member modified id > created id then recreate it else return new sign url
	var businessModified time.Time
	err = s.Get(&businessModified, "SELECT modified FROM business WHERE id = $1", businessID)
	if err != nil {
		return nil, err
	}

	var businessMemberModified time.Time
	err = s.Get(&businessMemberModified, "SELECT modified FROM business_member WHERE business_id = $1", businessID)
	if err != nil {
		return nil, err
	}

	if businessModified.After(signature.Modified) || businessMemberModified.After(signature.Modified) {
		// Cancel existing signature request
		err := hellosign.NewHellosignService(lo).CancelSignatureRequest(signature.SignatureRequestID)
		if err != nil {
			return nil, err
		}

		// Create new request
		req, err := s.getSignatureRequest(businessID)
		if err != nil {
			return nil, err
		}

		signatureRequest := constructSignatureRequest(req, templateType)

		resp, err := hellosign.NewHellosignService(lo).CreateSignatureRequest(signatureRequest)
		if err != nil {
			return nil, err
		}

		// Update DB
		_, err = s.Exec(`UPDATE signature_request SET signature_request_id = $1, signature_id = $2 WHERE business_id = $3`,
			resp.SignatureRequest.SignatureRequestID, resp.SignatureRequest.Signatures[0].SignatureID, businessID)
		if err != nil {
			return nil, err
		}

		signature, err = s.getByBusinessID(businessID, templateType)
		if err != nil {
			return nil, err
		}

	}

	// get new signed url
	url, err := hellosign.NewHellosignService(lo).GetEmbeddedSignURL(signature.SignatureID)
	if err != nil {
		return nil, err
	}

	lo.Info("Sign url " + url.Embedded.SignURL + strconv.FormatInt(url.Embedded.ExpiresAt, 10))

	signature.SignURL = url.Embedded.SignURL

	return signature, nil
}

func (s *signatureDatastore) getByBusinessID(businessID shared.BusinessID, templateType SignatureRequestTemplate) (*SignatureRequest, error) {
	signature := SignatureRequest{}
	err := s.Get(&signature, "SELECT * FROM signature_request WHERE business_id = $1 AND template_type = $2", businessID, templateType)
	if err != nil {
		return nil, err
	}

	return &signature, nil
}

func (s *signatureDatastore) getByID(ID shared.SignatureRequestID) (*SignatureRequest, error) {
	signature := SignatureRequest{}
	err := s.Get(&signature, "SELECT * FROM signature_request WHERE id = $1", ID)
	if err != nil {
		return nil, err
	}

	return &signature, nil
}

func (s *signatureDatastore) getSignatureRequest(businessID shared.BusinessID) (*SignatureRequestJoin, error) {
	req := SignatureRequestJoin{}
	err := s.Get(&req,
		`SELECT legal_name,dba,entity_type,business.legal_address,business.tax_id "business.tax_id", consumer.email "email",
		business.tax_id_type "business.tax_id_type", first_name, middle_name, last_name, title_type, title_other, ownership 
		FROM business 
		JOIN business_member ON business.id = business_member.business_id 
		JOIN consumer ON business_member.consumer_id = consumer.id
		WHERE business.id = $1 AND business_member.is_controlling_manager = true`, businessID)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func (s *signatureDatastore) GetBySignatureRequestID(signatureRequestID string) (*SignatureRequest, error) {
	signature := SignatureRequest{}
	err := s.Get(&signature, "SELECT * FROM signature_request WHERE signature_request_id = $1", signatureRequestID)
	if err != nil {
		return nil, err
	}

	return &signature, nil
}

func (s *signatureDatastore) UpdateSignatureStatus(signatureRequestID shared.SignatureRequestID, documentID shared.BusinessDocumentID, status SignatureRequestStatus) (*SignatureRequest, error) {
	// Update DB
	_, err := s.Exec(`UPDATE signature_request SET signature_status = $1, document_id = $2 WHERE id = $3`,
		status, documentID, signatureRequestID)
	if err != nil {
		return nil, err
	}

	return s.getByID(signatureRequestID)
}

func constructSignatureRequest(r *SignatureRequestJoin, templateType SignatureRequestTemplate) hellosign.SignatureRequest {

	controlPersonName := r.FirstName
	if r.MiddleName != nil && len(*r.MiddleName) > 0 {
		controlPersonName = controlPersonName + " " + *r.MiddleName
	}
	controlPersonName = controlPersonName + " " + r.LastName

	if r.EntityType == business.EntityTypeSoleProprietor {
		r.LegalName = &controlPersonName
	}

	request := hellosign.SignatureRequest{
		Email:        r.EmailAddress,
		Name:         controlPersonName,
		TemplateName: SignatureRequestTemplateTo[templateType],
	}

	var dba string
	if r.DBA != nil && len(r.DBA) > 0 {
		if len(r.DBA[0]) > 0 {
			dba = r.DBA[0]
		}
	}

	if r.TitleOther != nil && len(*r.TitleOther) > 0 {
		r.TitleType = *r.TitleOther
	}

	fields := []hellosign.CustomField{}

	if r.LegalName != nil {
		legalName := hellosign.CustomField{
			Name:  LegalName,
			Value: *r.LegalName,
		}
		fields = append(fields, legalName)
	}

	legalAddress := r.LegalAddress.StreetAddress
	if len(r.LegalAddress.AddressLine2) > 0 {
		legalAddress = legalAddress + ", " + r.LegalAddress.AddressLine2
	}
	streetAddress := hellosign.CustomField{
		Name:  StreetAddress,
		Value: legalAddress,
	}
	fields = append(fields, streetAddress)

	cityField := hellosign.CustomField{
		Name:  City,
		Value: r.LegalAddress.City,
	}
	fields = append(fields, cityField)

	stateField := hellosign.CustomField{
		Name:  State,
		Value: r.LegalAddress.State,
	}
	fields = append(fields, stateField)

	opStateField := hellosign.CustomField{
		Name:  OperatingState,
		Value: r.LegalAddress.State,
	}
	fields = append(fields, opStateField)

	postalCodeField := hellosign.CustomField{
		Name:  PostalCode,
		Value: r.LegalAddress.PostalCode,
	}
	fields = append(fields, postalCodeField)

	nameField := hellosign.CustomField{
		Name:  ControlPersonName,
		Value: controlPersonName,
	}
	fields = append(fields, nameField)

	title, ok := controlPersonTitle[r.TitleType]
	if !ok {
		title = r.TitleType
	}

	titleField := hellosign.CustomField{
		Name:  ControlPersonTitle,
		Value: title,
	}
	fields = append(fields, titleField)

	tinField := hellosign.CustomField{
		Name:  EIN,
		Value: r.TaxID.String(),
	}
	fields = append(fields, tinField)

	llpField := hellosign.CustomField{
		Name:  LimitedLiabilityPartnership,
		Value: "false",
	}
	singleLLCField := hellosign.CustomField{
		Name:  SingleLimitedLiabilityCompany,
		Value: "false",
	}
	multiLLCField := hellosign.CustomField{
		Name:  MultiLimitedLiabilityCompany,
		Value: "false",
	}
	corpField := hellosign.CustomField{
		Name:  Corporation,
		Value: "false",
	}
	solePropField := hellosign.CustomField{
		Name:  SoleProp,
		Value: "false",
	}
	dbaField := hellosign.CustomField{
		Name:  DBA,
		Value: "false",
	}

	switch r.EntityType {
	case BusinessEntityLimitedLiabilityPartnership:
		llpField = hellosign.CustomField{
			Name:  LimitedLiabilityPartnership,
			Value: "true",
		}
	case BusinessEntityLimitedLiabilityCompany:
		if r.Ownership == 100 {
			singleLLCField = hellosign.CustomField{
				Name:  SingleLimitedLiabilityCompany,
				Value: "true",
			}
		} else {
			multiLLCField = hellosign.CustomField{
				Name:  MultiLimitedLiabilityCompany,
				Value: "true",
			}
		}
	case BusinessEntityUnlistedCorporation:
		corpField = hellosign.CustomField{
			Name:  Corporation,
			Value: "true",
		}
	case BusinessEntitySoleProprietor:
		if len(dba) > 0 {
			dbaField = hellosign.CustomField{
				Name:  DBA,
				Value: "true",
			}
		} else {
			solePropField = hellosign.CustomField{
				Name:  SoleProp,
				Value: "true",
			}
		}
	}

	fields = append(fields, llpField)
	fields = append(fields, singleLLCField)
	fields = append(fields, multiLLCField)
	fields = append(fields, corpField)
	fields = append(fields, solePropField)
	fields = append(fields, dbaField)

	t := time.Now()

	dayField := hellosign.CustomField{
		Name:  Day,
		Value: strconv.Itoa(t.Day()),
	}
	fields = append(fields, dayField)

	monthField := hellosign.CustomField{
		Name:  Month,
		Value: t.Month().String()[:3],
	}
	fields = append(fields, monthField)

	yearField := hellosign.CustomField{
		Name:  Year,
		Value: strconv.Itoa(t.Year())[2:],
	}
	fields = append(fields, yearField)

	request.CustomField = fields

	return request
}
