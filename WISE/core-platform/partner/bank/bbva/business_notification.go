package bbva

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/partner/bank/data"
)

type NotificationIDField struct {
	ID string `json:"id"`
}

type EntityFormationNotification struct {
	BusinessType NotificationIDField `json:"businessType"`
	Country      NotificationIDField `json:"country"`
	State        NotificationIDField `json:"state"`
	Date         string              `json:"date"`
}

type BusinessLegalDocumentNotification struct {
	LegalDocumentType NotificationIDField `json:"legalDocumentType"`
	Country           NotificationIDField `json:"country"`
	State             NotificationIDField `json:"state"`
	IssueDate         string              `json:"issueDate"`
	ExpirationDate    string              `json:"expirationDate"`
	DocumentNumber    string              `json:"documentNumber"`
}

type BusinessDocumentNotification struct {
	BusinessDocumentType NotificationIDField `json:"businessDocumentType"`
	DocumentNumber       string              `json:"documentNumber"`
}

type BusinessNotification struct {
	BusinessID          string                              `json:"id"`
	LegalName           string                              `json:"legalName"`
	RegistrationPurpose string                              `json:"registrationPurpose"`
	DBA                 string                              `json:"doingBusinessAs"`
	EconomicActivity    NotificationIDField                 `json:"economicActivity"`
	EntityType          NotificationIDField                 `json:"entityType"`
	Formation           *EntityFormationNotification        `json:"formation"`
	LegalDocuments      []BusinessLegalDocumentNotification `json:"legalDocuments"`
	BusinessDocuments   []BusinessDocumentNotification      `json:"businessDocuments"`
}

var partnerBusinessActionNotificationTo = map[EventType]bank.NotificationAction{
	EventTypeBusinessCreate:          bank.NotificationActionCreate,
	EventTypeBusinessUpdate:          bank.NotificationActionUpdate,
	EventTypeBusinessOwnerCreate:     bank.NotificationActionCreate,
	EventTypeBusinessOwnerDelete:     bank.NotificationActionDelete,
	EventTypeBusinessIndicatorUpdate: bank.NotificationActionUpdate,
	EventTypeBusinessContactCreate:   bank.NotificationActionCreate,
	EventTypeBusinessContactUpdate:   bank.NotificationActionUpdate,
	EventTypeBusinessContactDelete:   bank.NotificationActionDelete,
	EventTypeBusinessAddressCreate:   bank.NotificationActionCreate,
	EventTypeBusinessAddressUpdate:   bank.NotificationActionUpdate,
	EventTypeBusinessAddressDelete:   bank.NotificationActionDelete,
	EventTypeBusinessKYC:             bank.NotificationActionUpdate,
}

var partnerBusinessAttributeToNotification = map[EventType]bank.NotificationAttribute{
	EventTypeBusinessAddressCreate:   bank.NotificationAttributeAddress,
	EventTypeBusinessAddressUpdate:   bank.NotificationAttributeAddress,
	EventTypeBusinessAddressDelete:   bank.NotificationAttributeAddress,
	EventTypeBusinessIndicatorUpdate: bank.NotificationAttributeActivity,
	EventTypeBusinessKYC:             bank.NotificationAttributeKYC,
}

func (s *notificationService) processBusinessEntityNotification(n Notification) error {
	hasPayload := true
	var d BusinessNotification
	err := n.unmarshalData(&d)
	if err != nil {
		if perr, ok := err.(*bank.Error); ok && perr.Code == bank.ErrorCodeNoPayload {
			hasPayload = false
		} else {
			return &bank.Error{
				RawError: err,
				Code:     bank.ErrorCodeProcessTransactionNotification,
			}
		}
	}

	// Get customer
	customerID := d.BusinessID
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

	action, ok := partnerBusinessActionNotificationTo[n.EventType]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationAction)
	}

	entityType, ok := partnerEntityToMap[BusinessEntity(strings.ToLower(d.EntityType.ID))]
	if !ok && hasPayload {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidEntityType)
	}

	// Formation information
	var operations *bank.BusinessOperationType
	var originCountry *bank.Country
	var originState *string
	var originDate *time.Time
	if d.Formation != nil {
		ops, ok := parterOperationTypeMapTo[BusinessOperation(strings.ToLower(d.Formation.BusinessType.ID))]
		if ok {
			operations = &ops
		}

		originState = &d.Formation.State.ID

		country, ok := partnerCountryTo[Country3Alpha(d.Formation.Country.ID)]
		if ok {
			originCountry = &country
		}

		date, err := time.Parse("2006-01-02", d.Formation.Date)
		if err == nil {
			originDate = &date
		}
	}

	// Tax ID type
	var taxID *string
	var taxIDType *bank.BusinessTaxIDType
	for _, doc := range d.BusinessDocuments {
		tidType, ok := partnerTaxIDToMap[BusinessTINType(strings.ToLower(doc.BusinessDocumentType.ID))]
		if ok {
			taxIDType = &tidType
			taxID = &doc.DocumentNumber
		}
	}

	// Documents
	var documents []bank.BusinessDocumentNotification
	for _, doc := range d.LegalDocuments {
		docType, ok := partnerBusinessDocumentTo[BusinessIdentityDocument(strings.ToLower(doc.LegalDocumentType.ID))]
		if !ok {
			continue
		}

		issueState := &doc.State.ID
		var issueCountry *bank.Country
		country, ok := partnerCountryTo[Country3Alpha(doc.Country.ID)]
		if ok {
			issueCountry = &country
		}

		var issueDate *time.Time
		issDate, err := time.Parse("2006-01-02", doc.IssueDate)
		if err == nil {
			issueDate = &issDate
		}

		var expDate *time.Time
		exp, err := time.Parse("2006-01-02", doc.ExpirationDate)
		if err == nil {
			expDate = &exp
		}

		ncDoc := bank.BusinessDocumentNotification{
			DocumentType:   docType,
			Number:         doc.DocumentNumber,
			Issuer:         issueState,
			IssueDate:      issueDate,
			IssueState:     issueState,
			IssueCountry:   issueCountry,
			ExpirationDate: expDate,
		}

		documents = append(documents, ncDoc)
	}

	// Business notification
	nc := bank.BusinessNotification{
		BankID:        bank.BusinessBankID(d.BusinessID),
		LegalName:     d.LegalName,
		TaxID:         taxID,
		TaxIDType:     taxIDType,
		EntityType:    &entityType,
		OperationType: operations,
		OriginCountry: originCountry,
		OriginState:   originState,
		OriginDate:    originDate,
		Documents:     documents,
	}

	b, err := json.Marshal(nc)
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

// Business owner
type BusinessNotificationOwnerCustomerID struct {
	CustomerID string `json:"customerId"`
}

type BusinessNotificationOwner struct {
	OwnerID            string                              `json:"id"`
	OwnerCustomerID    BusinessNotificationOwnerCustomerID `json:"ownerIdentification"`
	OwnerType          NotificationIDField                 `json:"ownerType"`
	Ownership          int                                 `json:"ownerShip"`
	ProfessionPosition NotificationIDField                 `json:"professionPosition"`
}

func (s *notificationService) processBusinessOwnerNotification(n Notification) error {
	var d BusinessNotificationOwner
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

	action, ok := partnerBusinessActionNotificationTo[n.EventType]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationAction)
	}

	me, err := notificationEntityFromCustomerID(d.OwnerCustomerID.CustomerID)
	if err != nil {
		return err
	}

	isControllingManager := false
	ownerType := BusinessMemberType(strings.ToLower(d.OwnerType.ID))
	if ownerType == BusinessMemberTypeController {
		isControllingManager = true
	} else if ownerType != BusinessMemberTypeOwner {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidMemberType)
	}

	title, ok := partnerMemberTitleToMap[BusinessMemberTitle(strings.ToLower(d.ProfessionPosition.ID))]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidMemberTitle)
	}

	var member = bank.BusinessMemberNotification{
		MemberBankID:         bank.MemberBankID(me.BankID),
		ConsumerID:           bank.ConsumerID(me.EntityID),
		IsControllingManager: isControllingManager,
		Ownership:            d.Ownership,
		Title:                title,
		// TitleDesc:            titleDesc,
	}

	nc := bank.BusinessMembersNotification{
		BankID:  bank.BusinessBankID(e.BankID),
		Members: []bank.BusinessMemberNotification{member},
	}

	b, err := json.Marshal(nc)
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

type BusinessNotificationIndicator struct {
	IndicatorID string `json:"indicatorId"`
	Active      bool   `json:"active"`
}

type BusinessNotificationIndicators struct {
	Indicators []BusinessNotificationIndicator `json:"indicators"`
}

func (s *notificationService) processBusinessIndicatorNotification(n Notification) error {
	var d BusinessNotificationIndicators
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

	action, ok := partnerBusinessActionNotificationTo[n.EventType]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationAction)
	}

	var activities []bank.ExpectedActivity
	for _, indicator := range d.Indicators {
		if indicator.Active {
			actType, ok := partnerActivityTo[ExpectedActivity(strings.ToLower(indicator.IndicatorID))]
			if ok {
				activities = append(activities, actType)
			}
		}
	}

	nc := bank.BusinessAttributeNotification{
		BankID:     bank.BusinessBankID(e.BankID),
		Activities: activities,
	}

	b, err := json.Marshal(nc)
	if err != nil {
		return err
	}

	attr := bank.NotificationAttributeActivity
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

func (s *notificationService) processBusinessContactNotification(n Notification) error {
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

	action, ok := partnerBusinessActionNotificationTo[n.EventType]
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
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidBusinessContactType)
	}

	ca := bank.BusinessAttributeNotification{
		BankID:      bank.BusinessBankID(e.BankID),
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

func (s *notificationService) processBusinessAddressNotification(n Notification) error {
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

	action, ok := partnerBusinessActionNotificationTo[n.EventType]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationAction)
	}

	country, ok := partnerCountryTo[Country3Alpha(d.Location.Country.ID)]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidCountry)
	}

	addrType, ok := partnerBusinessNotificationAddressTo[NotificationAddressType(strings.ToLower(d.Type.ID))]
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

	ca := bank.BusinessAttributeNotification{
		BankID:      bank.BusinessBankID(e.BankID),
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

func (s *notificationService) processBusinessKYCChangeNotification(n Notification) error {
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

	action, ok := partnerBusinessActionNotificationTo[n.EventType]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationAction)
	}

	kycStatus, ok := kycStatusToPartnerMap[KYCStatus(strings.ToUpper(d.Status.ID))]
	if !ok {
		return bank.NewErrorFromCode(bank.ErrorCodeInvalidKYCStatus)
	}

	cn := bank.BusinessKYCNotification{
		BankID:    bank.BusinessBankID(e.BankID),
		Risk:      partnerKYCRiskTo[d.RiskScore],
		Result:    partnerKYCResultTo[d.ProcessScore],
		KYCStatus: kycStatus,
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

func (s *notificationService) processBusinessNotificationMessage(n Notification) error {

	switch n.EventType {
	case EventTypeBusinessCreate, EventTypeBusinessUpdate:
		return s.processBusinessEntityNotification(n)
	case EventTypeBusinessOwnerCreate, EventTypeBusinessOwnerDelete:
		return s.processBusinessOwnerNotification(n)
	case EventTypeBusinessIndicatorUpdate:
		return s.processBusinessIndicatorNotification(n)
	case EventTypeBusinessContactCreate, EventTypeBusinessContactUpdate, EventTypeBusinessContactDelete:
		return s.processBusinessContactNotification(n)
	case EventTypeBusinessAddressCreate, EventTypeBusinessAddressUpdate, EventTypeBusinessAddressDelete:
		return s.processBusinessAddressNotification(n)
	case EventTypeBusinessKYC:
		return s.processBusinessKYCChangeNotification(n)
	}

	return bank.NewErrorFromCode(bank.ErrorCodeInvalidNotificationFormat)
}
