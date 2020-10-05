package bbva

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/partner/bank/data"
)

type ConsumerPersonalInfoNotification struct {
	SSN                string `json:"ssn"`
	FirstName          string `json:"firstName"`
	MiddleName         string `json:"middleName"`
	LastName           string `json:"lastName"`
	BirthDate          string `json:"dob"`
	CitizenStatus      string `json:"citizenshipStatus"`
	CitizenshipCountry string `json:"citizenshipCountry"`
}

type ConsumerIdentityNotificationDocumentType string

const (
	ConsumerIdentityNotificationDocumentTypeDriversLicense        = ConsumerIdentityNotificationDocumentType("DL")
	ConsumerIdentityNotificationDocumentTypePassport              = ConsumerIdentityNotificationDocumentType("PASSPORT")
	ConsumerIdentityNotificationDocumentTypeSocialSecurityCard    = ConsumerIdentityNotificationDocumentType("SSN")
	ConsumerIdentityNotificationDocumentTypeStateID               = ConsumerIdentityNotificationDocumentType("STATE_ID")
	ConsumerIdentityNotificationDocumentTypeWorkPermit            = ConsumerIdentityNotificationDocumentType("WORK_PERMIT")
	ConsumerIdentityNotificationDocumentTypeAlienRegistrationCard = ConsumerIdentityNotificationDocumentType("ARC")
)

var partnerConsumerIdentityDocumentTypeTo = map[ConsumerIdentityNotificationDocumentType]bank.ConsumerIdentityDocument{
	ConsumerIdentityNotificationDocumentTypeDriversLicense:        bank.ConsumerIdentityDocumentDriversLicense,
	ConsumerIdentityNotificationDocumentTypePassport:              bank.ConsumerIdentityDocumentPassport,
	ConsumerIdentityNotificationDocumentTypeSocialSecurityCard:    bank.ConsumerIdentityDocumentSocialSecurityCard,
	ConsumerIdentityNotificationDocumentTypeStateID:               bank.ConsumerIdentityDocumentStateID,
	ConsumerIdentityNotificationDocumentTypeWorkPermit:            bank.ConsumerIdentityDocumentWorkPermit,
	ConsumerIdentityNotificationDocumentTypeAlienRegistrationCard: bank.ConsumerIdentityDocumentAlienRegistrationCard,
}

type ConsumerIdentityDocumentNotification struct {
	DocumentType   NotificationIDField `json:"documentType"`
	State          NotificationIDField `json:"state"`
	Country        NotificationIDField `json:"country"`
	IssueDate      *string             `json:"issueDate"`
	ExpirationDate *string             `json:"expirationDate"`
	DocumentNumber string              `json:"documentNumber"`
}

type ConsumerEntityNotification struct {
	CustomerID           string                                 `json:"customerId"`
	FirstName            string                                 `json:"firstName"`
	MiddleName           string                                 `json:"middleName"`
	LastName             string                                 `json:"lastName"`
	BirthData            ConsumerBirthDataNotification          `json:"dob"`
	Nationalities        []NotificationIDField                  `json:"nationalities"`
	SocialSecurityNumber string                                 `json:"socialSecurityNumber"`
	Residence            ConsumerResidenceNotification          `json:"residence"`
	PEP                  NotificationIDField                    `json:"politicallyExposedPerson"`
	IdentityDocuments    []ConsumerIdentityDocumentNotification `json:"identityDocuments"`
}

type ConsumerBirthDataNotification struct {
	BirthDate string `json:"dob"`
}

type ConsumerResidenceNotification struct {
	ResidenceType NotificationIDField `json:"residenceType"`
}

func (s *notificationService) processConsumerNotificationMessage(n Notification) error {
	switch n.EventType {
	case EventTypeConsumerCreate, EventTypeConsumerUpdate:
		return s.processConsumerEntityNotification(n)
	case EventTypeConsumerContactCreate, EventTypeConsumerContactUpdate, EventTypeConsumerContactDelete:
		return s.processConsumerContactNotification(n)
	case EventTypeConsumerAddressCreate, EventTypeConsumerAddressUpdate, EventTypeConsumerAddressDelete:
		return s.processConsumerAddressNotification(n)
	case EventTypeConsumerOccupationUpdate:
		return s.processConsumerOccupationNotification(n)
	case EventTypeConsumerKYC:
		return s.processConsumerKYCChangeNotification(n)
	}

	return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationFormat)
}

var partnerConsumerActionNotificationTo = map[EventType]bank.NotificationAction{
	EventTypeConsumerCreate:           bank.NotificationActionCreate,
	EventTypeConsumerUpdate:           bank.NotificationActionUpdate,
	EventTypeConsumerContactCreate:    bank.NotificationActionCreate,
	EventTypeConsumerContactUpdate:    bank.NotificationActionUpdate,
	EventTypeConsumerContactDelete:    bank.NotificationActionDelete,
	EventTypeConsumerAddressCreate:    bank.NotificationActionCreate,
	EventTypeConsumerAddressUpdate:    bank.NotificationActionUpdate,
	EventTypeConsumerAddressDelete:    bank.NotificationActionDelete,
	EventTypeConsumerOccupationUpdate: bank.NotificationActionUpdate,
	EventTypeConsumerKYC:              bank.NotificationActionUpdate,
}

var partnerConsumerAttributeToNotification = map[EventType]bank.NotificationAttribute{
	EventTypeConsumerAddressCreate:    bank.NotificationAttributeAddress,
	EventTypeConsumerAddressUpdate:    bank.NotificationAttributeAddress,
	EventTypeConsumerAddressDelete:    bank.NotificationAttributeAddress,
	EventTypeConsumerOccupationUpdate: bank.NotificationAttributeOccupation,
	EventTypeConsumerKYC:              bank.NotificationAttributeKYC,
}

type NotificationContact struct {
	ContactID      string              `json:"contactDetailId"`
	Value          string              `json:"contact"`
	Type           NotificationIDField `json:"contactType"`
	IsPreferential bool                `json:"isPreferential"`
}

func (s *notificationService) processConsumerEntityNotification(n Notification) error {
	var d ConsumerEntityNotification
	err := n.unmarshalData(&d)
	if err != nil {
		// Default to regular consumer entity variant
		return s.processConsumerEntityNotification(n)
	}

	// Get customer
	customerID := d.CustomerID
	if customerID == "" {
		customerID, err = n.customerID()
		if err != nil {
			return err
		}
	}

	e, err := notificationEntityFromCustomerID(customerID)
	if err != nil {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidEntityID)
	}

	nt, ok := partnerNotificationTypeTo[n.Type]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationType)
	}

	action, ok := partnerConsumerActionNotificationTo[n.EventType]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationAction)
	}

	var dob *time.Time
	date, err := time.Parse("2006-01-02", d.BirthData.BirthDate)
	if err != nil {
		bank.NewErrorFromCode(bank.ErrorCodeInvalidBirthDate)
	}

	dob = &date

	var country bank.Country
	for _, nat := range d.Nationalities {
		cid, ok := partnerCountryTo[Country3Alpha(nat.ID)]
		if ok {
			country = cid
		}
	}

	status := USResidencyStatus(strings.ToLower(d.Residence.ResidenceType.ID))
	residency, err := partnerResidencyTo(status)
	if err != nil {
		return err
	}

	var middleName *string
	if d.MiddleName != "" {
		middleName = &d.MiddleName
	}

	cn := bank.ConsumerNotification{
		BankID:             bank.ConsumerBankID(e.BankID),
		TaxID:              d.SocialSecurityNumber,
		TaxIDType:          bank.ConsumerTaxIDTypeSSN,
		FirstName:          d.FirstName,
		MiddleName:         middleName,
		LastName:           d.LastName,
		DateOfBirth:        dob,
		Residency:          &residency,
		CitizenshipCountry: &country,
	}

	var docs []bank.ConsumerDocument
	for _, doc := range d.IdentityDocuments {
		if doc.DocumentType.ID == "" {
			continue
		}

		docType, ok := partnerConsumerIdentityDocumentTypeTo[ConsumerIdentityNotificationDocumentType(doc.DocumentType.ID)]
		if !ok {
			return bank.NewErrorFromCode(bank.ErrorCodeInvalidDocumentType)
		}

		var issueDate *time.Time
		if doc.IssueDate != nil {
			date, err = time.Parse("2006-01-02", *doc.IssueDate)
			if err != nil {
				return bank.NewErrorFromCode(bank.ErrorCodeInvalidIssueDateFormat)
			}

			issueDate = &date
		}

		var expDate *time.Time
		if doc.ExpirationDate != nil {
			date, err = time.Parse("2006-01-02", *doc.ExpirationDate)
			if err != nil {
				return bank.NewErrorFromCode(bank.ErrorCodeInvalidExpDateFormat)
			}

			expDate = &date
		}

		var state *string
		if doc.State.ID != "" {
			state = &doc.State.ID
		}

		var country *bank.Country
		if doc.Country.ID != "" {
			c := bank.Country(doc.Country.ID)
			country = &c
		}

		d := bank.ConsumerDocument{
			DocumentType:   docType,
			Number:         doc.DocumentNumber,
			IssueDate:      issueDate,
			ExpirationDate: expDate,
			IssueState:     state,
			IssueCountry:   country,
		}

		docs = append(docs, d)
	}

	cn.Documents = docs

	b, err := json.Marshal(cn)
	if err != nil {
		return err
	}

	pn := data.NotificationCreate{
		EntityID:   bank.NotificationEntityID(e.EntityID),
		EntityType: bank.NotificationEntityType(e.EntityType),
		BankName:   bank.ProviderNameBBVA,
		SourceID:   bank.SourceID(n.notificationID()),
		Type:       nt,
		Action:     action,
		Version:    n.Version,
		Created:    n.SentDate,
		Data:       b,
	}

	return s.forwardNotification(pn)
}

func (s *notificationService) processConsumerContactNotification(n Notification) error {
	var d NotificationContact
	err := n.unmarshalData(&d)
	if err != nil {
		return &bank.Error{
			RawError: err,
			Code:     bank.ErrorCodeProcessTransactionNotification,
		}
	}

	// Get customer
	customerID, err := n.customerID()
	if err != nil {
		return err
	}

	e, err := notificationEntityFromCustomerID(customerID)
	if err != nil {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidEntityID)
	}

	nt, ok := partnerNotificationTypeTo[n.Type]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationType)
	}

	action, ok := partnerConsumerActionNotificationTo[n.EventType]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationAction)
	}

	var attr bank.NotificationAttribute
	var email *string
	var phone *string
	switch ContactType(strings.ToLower(d.Type.ID)) {
	case ContactTypeEmail:
		email = &d.Value
		attr = bank.NotificationAttributeEmail
	case ContactTypePhone, ContactTypeMobileNumber:
		phone = &d.Value
		attr = bank.NotificationAttributePhone
	default:
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidConsumerContactType)
	}

	ca := bank.ConsumerAttributeNotification{
		BankID:      bank.ConsumerBankID(e.BankID),
		AttributeID: d.ContactID,
		Phone:       phone,
		Email:       email,
	}

	b, err := json.Marshal(ca)
	if err != nil {
		return err
	}

	pn := data.NotificationCreate{
		EntityID:   bank.NotificationEntityID(e.EntityID),
		EntityType: bank.NotificationEntityType(e.EntityType),
		BankName:   bank.ProviderNameBBVA,
		SourceID:   bank.SourceID(n.notificationID()),
		Type:       nt,
		Attribute:  &attr,
		Action:     action,
		Version:    n.Version,
		Created:    n.SentDate,
		Data:       b,
	}

	return s.forwardNotification(pn)
}

type NotificationAddressLocation struct {
	StreetName string
	City       string
	State      NotificationIDField
	Country    NotificationIDField
	ZipCode    string
}

type NotificationAddress struct {
	AddressID      string                      `json:"addressId"`
	Value          string                      `json:"contact"`
	Type           NotificationIDField         `json:"addressType"`
	Location       NotificationAddressLocation `json:"location"`
	IsPreferential bool                        `json:"isPreferential"`
}

type NotificationAddressType string

const (
	NotificationAddressTypeLegal        = NotificationAddressType("legal")
	NotificationAddressTypePostal       = NotificationAddressType("postal")
	NotificationAddressTypeWork         = NotificationAddressType("work")
	NotificationAddressTypeMailing      = NotificationAddressType("mailing")
	NotificationAddressTypeHeadquarter  = NotificationAddressType("headquarter")
	NotificationAddressTypeHeadquarters = NotificationAddressType("headquarters")
	NotificationAddressTypeBilling      = NotificationAddressType("billing")
	NotificationAddressTypeAlternative  = NotificationAddressType("alternative")
)

var partnerConsumerNotificationAddressTo = map[NotificationAddressType]bank.AddressRequestType{
	NotificationAddressTypeLegal:       bank.AddressRequestTypeLegal,
	NotificationAddressTypePostal:      bank.AddressRequestTypeMailing,
	NotificationAddressTypeWork:        bank.AddressRequestTypeWork,
	NotificationAddressTypeAlternative: bank.AddressRequestTypeOther,
}

var partnerBusinessNotificationAddressTo = map[NotificationAddressType]bank.AddressRequestType{
	NotificationAddressTypeLegal:        bank.AddressRequestTypeLegal,
	NotificationAddressTypeMailing:      bank.AddressRequestTypeMailing,
	NotificationAddressTypePostal:       bank.AddressRequestTypeMailing,
	NotificationAddressTypeHeadquarter:  bank.AddressRequestTypeHeadquarter,
	NotificationAddressTypeHeadquarters: bank.AddressRequestTypeHeadquarter,
	NotificationAddressTypeAlternative:  bank.AddressRequestTypeOther,
}

func (s *notificationService) processConsumerAddressNotification(n Notification) error {
	var d NotificationAddress
	err := n.unmarshalData(&d)
	if err != nil {
		return &bank.Error{
			RawError: err,
			Code:     bank.ErrorCodeProcessTransactionNotification,
		}
	}

	customerID, err := n.customerID()
	if err != nil {
		return err
	}

	e, err := notificationEntityFromCustomerID(customerID)
	if err != nil {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidEntityID)
	}

	nt, ok := partnerNotificationTypeTo[n.Type]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationType)
	}

	action, ok := partnerConsumerActionNotificationTo[n.EventType]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationAction)
	}

	country, ok := partnerCountryTo[Country3Alpha(d.Location.Country.ID)]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidCountry)
	}

	addrType, ok := partnerConsumerNotificationAddressTo[NotificationAddressType(strings.ToLower(d.Type.ID))]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidAddressType)
	}

	addr := bank.AddressResponse{
		Type:    addrType,
		Line1:   d.Location.StreetName,
		City:    d.Location.City,
		State:   d.Location.State.ID,
		ZipCode: d.Location.ZipCode,
		Country: country,
	}

	ca := bank.ConsumerAttributeNotification{
		BankID:      bank.ConsumerBankID(e.BankID),
		AttributeID: d.AddressID,
		Address:     &addr,
	}

	b, err := json.Marshal(ca)
	if err != nil {
		return err
	}

	attr := bank.NotificationAttributeAddress
	pn := data.NotificationCreate{
		EntityID:   bank.NotificationEntityID(e.EntityID),
		EntityType: bank.NotificationEntityType(e.EntityType),
		BankName:   bank.ProviderNameBBVA,
		SourceID:   bank.SourceID(n.notificationID()),
		Type:       nt,
		Attribute:  &attr,
		Action:     action,
		Version:    n.Version,
		Created:    n.SentDate,
		Data:       b,
	}

	return s.forwardNotification(pn)
}

type NotificationConsumerProperty struct {
	Occupation string `json:"occupationCode"`
}

func (s *notificationService) processConsumerOccupationNotification(n Notification) error {
	var d NotificationConsumerProperty
	err := n.unmarshalData(&d)
	if err != nil {
		return &bank.Error{
			RawError: err,
			Code:     bank.ErrorCodeProcessTransactionNotification,
		}
	}

	customerID, err := n.customerID()
	if err != nil {
		return err
	}

	e, err := notificationEntityFromCustomerID(customerID)
	if err != nil {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidEntityID)
	}

	nt, ok := partnerNotificationTypeTo[n.Type]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationType)
	}

	action, ok := partnerConsumerActionNotificationTo[n.EventType]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationAction)
	}

	occ, ok := occupationToPartnerMap[ConsumerOccupation(strings.ToLower(d.Occupation))]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidOccupation)
	}

	ca := bank.ConsumerAttributeNotification{
		BankID:     bank.ConsumerBankID(e.BankID),
		Occupation: &occ,
	}

	b, err := json.Marshal(ca)
	if err != nil {
		return err
	}

	attr := bank.NotificationAttributeOccupation
	pn := data.NotificationCreate{
		EntityID:   bank.NotificationEntityID(e.EntityID),
		EntityType: bank.NotificationEntityType(e.EntityType),
		BankName:   bank.ProviderNameBBVA,
		SourceID:   bank.SourceID(n.notificationID()),
		Type:       nt,
		Attribute:  &attr,
		Action:     action,
		Version:    n.Version,
		Created:    n.SentDate,
		Data:       b,
	}

	return s.forwardNotification(pn)
}

type NotificationDigitalFootprint struct {
	Link string `json:"link"`
}

type NotificationDecisionReason struct {
	Source string `json:"type"`
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

type ProcessScore string

const (
	ProcessScorePass = ProcessScore("PASS")
	ProcessScoreFail = ProcessScore("FAIL")
)

var partnerKYCResultTo = map[ProcessScore]bank.KYCResult{
	ProcessScorePass: bank.KYCResultPass,
	ProcessScoreFail: bank.KYCResultFail,
}

type RiskScore string

const (
	RiskScoreHigh   = RiskScore("HIGH")
	RiskScoreMedium = RiskScore("MEDIUM")
	RiskScoreLow    = RiskScore("LOW")
)

var partnerKYCRiskTo = map[RiskScore]bank.KYCRisk{
	RiskScoreHigh:   bank.KYCRiskHigh,
	RiskScoreMedium: bank.KYCRiskMedium,
	RiskScoreLow:    bank.KYCRiskLow,
}

type NotificationKYCStatusChange struct {
	Status           NotificationIDField            `json:"status"`
	ProcessScore     ProcessScore                   `json:"processScore"`
	RiskScore        RiskScore                      `json:"riskScore"`
	DigitalFootprint []NotificationDigitalFootprint `json:"digitalFootprint"`
	DecisionReasons  []NotificationDecisionReason   `json:"decisionReasons"`
}

func (s *notificationService) processConsumerKYCChangeNotification(n Notification) error {
	var d NotificationKYCStatusChange
	err := n.unmarshalData(&d)
	if err != nil {
		return &bank.Error{
			RawError: err,
			Code:     bank.ErrorCodeProcessTransactionNotification,
		}
	}

	customerID, err := n.customerID()
	if err != nil {
		return err
	}

	e, err := notificationEntityFromCustomerID(customerID)
	if err != nil {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidEntityID)
	}

	nt, ok := partnerNotificationTypeTo[n.Type]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationType)
	}

	action, ok := partnerConsumerActionNotificationTo[n.EventType]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationAction)
	}

	kycStatus, ok := kycStatusToPartnerMap[KYCStatus(strings.ToUpper(d.Status.ID))]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidKYCStatus)
	}

	var notes []bank.ConsumerKYCNotesNotification
	for _, dec := range d.DecisionReasons {
		note := bank.ConsumerKYCNotesNotification{
			Source: dec.Source,
			Desc:   dec.Detail,
		}

		notes = append(notes, note)
	}

	cn := bank.ConsumerKYCNotification{
		BankID:    bank.ConsumerBankID(e.BankID),
		Risk:      partnerKYCRiskTo[d.RiskScore],
		Result:    partnerKYCResultTo[d.ProcessScore],
		KYCStatus: kycStatus,
		KYCNotes:  notes,
	}

	b, err := json.Marshal(cn)
	if err != nil {
		return err
	}

	attr := bank.NotificationAttributeKYC
	pn := data.NotificationCreate{
		EntityID:   bank.NotificationEntityID(e.EntityID),
		EntityType: bank.NotificationEntityType(e.EntityType),
		BankName:   bank.ProviderNameBBVA,
		SourceID:   bank.SourceID(n.notificationID()),
		Type:       nt,
		Attribute:  &attr,
		Action:     action,
		Version:    n.Version,
		Created:    n.SentDate,
		Data:       b,
	}

	return s.forwardNotification(pn)
}
