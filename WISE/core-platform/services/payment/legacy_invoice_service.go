package payment

import (
	"errors"
	"fmt"

	"github.com/wiseco/core-platform/services/invoice"

	"github.com/wiseco/core-platform/shared"

	"github.com/wiseco/go-lib/id"
	grpcRoot "github.com/wiseco/protobuf/golang"
	grpcSvc "github.com/wiseco/protobuf/golang/invoice"
	grpcItem "github.com/wiseco/protobuf/golang/invoice/item"
)

type LegacyInvoiceService struct{}

func newLegacyInvoiceService() *LegacyInvoiceService {
	return &LegacyInvoiceService{}
}

func (lis *LegacyInvoiceService) CreateInvoice(request Request) (*invoice.Invoice, error) {
	if request.Notes == nil {
		return nil, errors.New("notes is required for creating invoice")
	}
	invSvc, err := invoice.NewInvoiceService()
	if err != nil {
		return nil, err
	}

	lineItemReq := &grpcItem.Item{
		Name:        *request.Notes,
		Price:       fmt.Sprintf("%f", request.Amount),
		IsTaxable:   grpcRoot.Boolean_B_FALSE,
		BusinessId:  request.BusinessID.ToPrefixString(),
		IsAvailable: grpcRoot.Boolean_B_TRUE,
	}

	item, err := invSvc.CreateItem(lineItemReq)
	if err != nil {
		return nil, err
	}

	invoiceLineItem := &grpcSvc.LineItem{
		ItemId:    item.Id,
		Quantity:  1,
		IsTaxable: grpcRoot.Boolean_B_FALSE,
	}
	var lineItems []*grpcSvc.LineItem
	lineItems = append(lineItems, invoiceLineItem)

	invoiceReq := &grpcSvc.InvoiceRequest{
		ActorId:         request.CreatedUserID.ToPrefixString(),
		BusinessId:      request.BusinessID.ToPrefixString(),
		ContactId:       fmt.Sprintf("%s%s", id.IDPrefixContact, *request.ContactID),
		LineItems:       lineItems,
		ShowBankAccount: grpcRoot.Boolean_B_TRUE,
		NotifyEmail:     grpcRoot.Boolean_B_TRUE,
	}
	if request.IPAddress != nil {
		invoiceReq.RequestIp = *request.IPAddress
	}
	if request.RequestSourceID != nil {
		invoiceReq.RequestSourceId = *request.RequestSourceID
	}
	if request.RequestSource != nil {
		invoiceReq.RequestSource = string(*request.RequestSource)
	}
	if request.Notes != nil {
		invoiceReq.Notes = *request.Notes
	}
	cardTransfer, bankTransfer := getRequestType(*request.RequestType)
	invoiceReq.PaymentTypes = &grpcSvc.InvoicePaymentType{
		AllowBankTransfer: bankTransfer,
		AllowCards:        cardTransfer,
	}

	invSvc, err = invoice.NewInvoiceService()
	if err != nil {
		return nil, err
	}
	invModel, err := invSvc.CreateInvoice(invoiceReq)
	if err != nil {
		return nil, err
	}
	return invModel, nil
}

func getRequestType(reqType PaymentRequestType) (cardTransfer grpcRoot.Boolean, bankTransfer grpcRoot.Boolean) {
	trueValue := grpcRoot.Boolean_B_TRUE
	falseValue := grpcRoot.Boolean_B_FALSE

	switch reqType {
	case PaymentRequestTypeInvoiceCard:
		return trueValue, falseValue
	case PaymentRequestTypeInvoiceBank:
		return falseValue, trueValue
	case PaymentRequestTypeInvoiceCardAndBank:
		return trueValue, trueValue
	case PaymentRequestTypeInvoiceNone:
		return falseValue, falseValue
	default:
		return falseValue, falseValue
	}
}

func (lis *LegacyInvoiceService) GetInvoice(invoiceID shared.PaymentRequestID) (*Request, error) {
	invSvc, err := invoice.NewInvoiceService()
	if err != nil {
		return nil, err
	}
	invoiceGot, err := invSvc.GetInvoiceByID(convertPaymentRequestIDToInvoiceID(invoiceID))
	return transformToLegacyRequest(invoiceGot)
}

func (lis *LegacyInvoiceService) UpdateInvoiceStatus(invoiceID shared.PaymentRequestID, status PaymentRequestStatus) error {
	invSvc, err := invoice.NewInvoiceService()
	if err != nil {
		return err
	}
	var action grpcSvc.InvoiceAction
	invoiceIDs := []string{
		convertPaymentRequestIDToInvoiceID(invoiceID).String(),
	}
	if status == PaymentRequestStatusComplete {
		action = grpcSvc.InvoiceAction_IA_MARK_AS_PAID
	} else if status == PaymentRequestStatusCanceled {
		action = grpcSvc.InvoiceAction_IA_MARK_AS_CANCELLED
	} else {
		action = grpcSvc.InvoiceAction_IA_UNSPECIFIED
	}
	err = invSvc.ExecuteInvoiceAction(invoiceIDs, action)
	return err
}

func (lis *LegacyInvoiceService) SendReminder(invoiceIDs []PaymentRequestResend) error {
	invSvc, err := invoice.NewInvoiceService()
	if err != nil {
		return err
	}
	invoiceIDsParsed := []string{}
	for _, rID := range invoiceIDs {
		invoiceIDsParsed = append(invoiceIDsParsed, convertPaymentRequestIDToInvoiceID(rID.RequestID).String())
	}
	return invSvc.ExecuteInvoiceAction(invoiceIDsParsed, grpcSvc.InvoiceAction_IA_SEND_REMINDER)
}

func (lis *LegacyInvoiceService) GetManyInvoices(businessID shared.BusinessID, limit int, offset int, status string, contactId string) ([]Request, error) {
	invSvc, err := invoice.NewInvoiceService()
	if err != nil {
		return nil, err
	}
	var cntID id.ContactID
	if contactId != "" {
		contactIDParsed, err := id.ParseContactID(contactId)
		if err != nil {
			contactIDParsed, err := id.ParseContactID(fmt.Sprintf("%s%s", id.IDPrefixContact, contactId))
			if err == nil {
				cntID = contactIDParsed
			}
		} else {
			cntID = contactIDParsed
		}
	}
	busID, _ := id.ParseBusinessID(businessID.ToPrefixString())
	invoices, err := invSvc.GetInvoices(busID, cntID, limit, offset, status)
	if err != nil {
		return nil, err
	}
	payments := []Request{}
	for _, invoice := range invoices {
		paymentObj, err := transformToLegacyRequest(invoice)
		if err != nil {
			return nil, err
		}
		payments = append(payments, *paymentObj)
	}
	return payments, nil
}

func transformToLegacyRequest(model *invoice.Invoice) (*Request, error) {
	request := Request{}

	userID, err := shared.ParseUserID(model.UserID.String())
	if err != nil {
		return &request, err
	}
	request.CreatedUserID = userID

	busID, err := shared.ParseBusinessID(model.BusinessID.String())
	if err != nil {
		return &request, err
	}
	request.BusinessID = busID

	contactUUID := model.ContactID.UUIDString()
	request.ContactID = &contactUUID

	reqID, err := shared.ParsePaymentRequestID(fmt.Sprintf("%s%s", shared.PaymentRequestPrefix, model.InvoiceID.UUIDString()))
	if err != nil {
		return &request, err
	}
	request.ID = reqID

	if amount, ok := model.Amount.Float64(); ok {
		request.Amount = amount
	} else {
		return &request, errors.New("unable to convert amount")
	}
	request.Currency = CurrencyUSD

	customMessage := model.Notes
	request.Notes = &customMessage

	status := getStatusFromGrpc(model.Status)
	request.Status = &status

	request.RequestType = getRequestTypeFromGrpc(model)

	ipAddress := model.IPAddress
	request.IPAddress = &ipAddress

	sourceId := model.RequestSourceID
	request.RequestSourceID = &sourceId

	source := model.RequestSource
	shopifySource := RequestSourceShopify
	if source != "" && source == string(RequestSourceShopify) {
		request.RequestSource = &shopifySource
	}

	createdTime, err := model.GetCreatedTime()
	if err == nil {
		request.Created = createdTime
	}

	return &request, nil

}

func getStatusFromGrpc(status grpcSvc.InvoiceRequestStatus) PaymentRequestStatus {
	switch status {
	case grpcSvc.InvoiceRequestStatus_IRT_OPEN:
		return PaymentRequestStatusPending
	case grpcSvc.InvoiceRequestStatus_IRT_PAID:
		return PaymentRequestStatusComplete
	case grpcSvc.InvoiceRequestStatus_IRT_CANCELLED:
		return PaymentRequestStatusCanceled
	case grpcSvc.InvoiceRequestStatus_IRT_PROCESSING_PAYMENT:
		return PaymentRequestStatusInProcess
	case grpcSvc.InvoiceRequestStatus_IRT_FAILED:
		return PaymentRequestStatusFailed
	default:
		return PaymentRequestStatusPending
	}
}

func getRequestTypeFromGrpc(inv *invoice.Invoice) *PaymentRequestType {
	card := PaymentRequestTypeInvoiceCard
	bank := PaymentRequestTypeInvoiceBank
	cardAndBank := PaymentRequestTypeInvoiceCardAndBank
	none := PaymentRequestTypeInvoiceNone

	if inv.AllowCard && inv.AllowBankTransfer {
		return &cardAndBank
	} else if inv.AllowBankTransfer {
		return &bank
	} else if inv.AllowCard {
		return &card
	} else {
		return &none
	}
}

func convertPaymentRequestIDToInvoiceID(pmtID shared.PaymentRequestID) id.InvoiceID {
	invId, _ := id.ParseInvoiceID(fmt.Sprintf("%s%s", id.IDPrefixInvoice, pmtID.ToUUIDString()))
	return invId
}
