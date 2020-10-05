package invoice

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	grpcTypes "github.com/golang/protobuf/ptypes"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/go-lib/grpc"
	"github.com/wiseco/go-lib/id"
	"github.com/wiseco/go-lib/num"
	grpcRoot "github.com/wiseco/protobuf/golang"
	grpcSvc "github.com/wiseco/protobuf/golang/invoice"
	grpcItem "github.com/wiseco/protobuf/golang/invoice/item"
	"google.golang.org/protobuf/types/known/emptypb"
)

type invoiceService struct {
	conn                 grpc.Client
	invoiceServiceClient grpcSvc.InvoiceServiceClient
	itemServiceClient    grpcItem.ItemServiceClient
}

func NewInvoiceService() (*invoiceService, error) {
	var bts *invoiceService

	bsn, err := grpc.GetConnectionStringForService(grpc.ServiceNameInvoice)
	if err != nil {
		return bts, err
	}

	conn, err := grpc.NewInsecureClient(bsn)
	if err != nil {
		return bts, err
	}

	irc := grpcSvc.NewInvoiceServiceClient(conn.GetConn())
	is := grpcItem.NewItemServiceClient(conn.GetConn())

	return &invoiceService{
		conn:                 conn,
		invoiceServiceClient: irc,
		itemServiceClient:    is,
	}, nil
}

func (is invoiceService) GetInvoiceByID(invoiceID id.InvoiceID) (*Invoice, error) {
	var inv *Invoice

	defer is.conn.CloseAndCancel()

	r := &grpcSvc.InvoiceIdGetRequest{
		InvoiceId: invoiceID.String(),
	}
	invoice, err := is.invoiceServiceClient.GetInvoice(context.Background(), r)
	if err != nil {
		return inv, err
	}
	if invoice == nil {
		return nil, errors.New("no invoice found")
	}

	return transformFullProtoToInvoice(invoice)
}

func (is invoiceService) GetInvoiceByRequestSourceID(requestSourceID id.ShopifyOrderID, businessID shared.BusinessID) (*Invoice, error) {
	var inv *Invoice

	defer is.conn.CloseAndCancel()

	r := &grpcSvc.InvoiceGetRequest{
		RequestSourceId: requestSourceID.UUIDString(),
		BusinessId:      businessID.ToPrefixString(),
		Limit:           1,
		Offset:          0,
	}
	invoices, err := is.invoiceServiceClient.GetManyInvoices(context.Background(), r)
	if err != nil {
		return inv, err
	}
	if invoices == nil || len(invoices.Invoices) == 0 {
		return nil, errors.New("no invoice found")
	}

	return transformProtoToInvoice(invoices.Invoices[0])
}

func (is invoiceService) GetInvoices(businessID id.BusinessID, contactID id.ContactID,
	limit, offset int, status string) ([]*Invoice, error) {
	var response []*Invoice
	defer is.conn.CloseAndCancel()

	r := &grpcSvc.InvoiceGetRequest{
		BusinessId: businessID.String(),
		Limit:      int32(limit),
		Offset:     int32(offset),
	}
	if !contactID.IsZero() {
		r.ContactId = contactID.String()
	}
	if status != "" {
		r.Status = []grpcSvc.InvoiceRequestStatus{
			getInvoiceRequestStatus(status),
		}
	}
	invoices, err := is.invoiceServiceClient.GetManyInvoices(context.Background(), r)
	if err != nil {
		return response, err
	}
	if invoices == nil || len(invoices.Invoices) == 0 {
		return nil, errors.New("no invoice found")
	}

	return transformProtoToInvoices(invoices)
}
func transformProtoToInvoice(request *grpcSvc.InvoiceRequest) (*Invoice, error) {
	businessIdParsed, err := id.ParseBusinessID(request.BusinessId)
	if err != nil {
		return nil, err
	}
	contactIdParsed, err := id.ParseContactID(request.ContactId)
	if err != nil {
		return nil, err
	}
	invoiceIdParsed, err := id.ParseInvoiceID(request.Id)
	if err != nil {
		return nil, err
	}
	userIdParsed, err := id.ParseUserID(request.ActorId)
	if err != nil {
		return nil, err
	}
	amountParsed, err := num.ParseDecimal(request.Amount)
	if err != nil {
		return nil, err
	}

	inv := &Invoice{
		BusinessID:      businessIdParsed,
		ContactID:       contactIdParsed,
		Amount:          amountParsed,
		InvoiceID:       invoiceIdParsed,
		Notes:           request.Notes,
		Title:           request.Title,
		ShowBankAccount: request.ShowBankAccount.ToBool(),
		UserID:          userIdParsed,
		Number:          request.Number,
		Status:          request.Status,
		IPAddress:       request.RequestIp,
		RequestSource:   request.RequestSource,
		RequestSourceID: request.RequestSourceId,
		Created:         request.Created,
	}

	if request.PaymentTypes != nil {
		inv.AllowCard = request.PaymentTypes.AllowCards.ToBool()
		inv.AllowBankTransfer = request.PaymentTypes.AllowBankTransfer.ToBool()
	}

	return inv, nil
}

func transformFullProtoToInvoice(request *grpcSvc.FullInvoice) (*Invoice, error) {
	businessIdParsed, err := id.ParseBusinessID(request.BusinessId)
	if err != nil {
		return nil, err
	}
	contactIdParsed, err := id.ParseContactID(request.ContactId)
	if err != nil {
		return nil, err
	}
	invoiceIdParsed, err := id.ParseInvoiceID(request.Id)
	if err != nil {
		return nil, err
	}
	userIdParsed, err := id.ParseUserID(request.ActorId)
	if err != nil {
		return nil, err
	}
	amountParsed, err := num.ParseDecimal(request.TotalAmount)
	if err != nil {
		return nil, err
	}

	inv := &Invoice{
		BusinessID:      businessIdParsed,
		ContactID:       contactIdParsed,
		Amount:          amountParsed,
		InvoiceID:       invoiceIdParsed,
		Notes:           request.Notes,
		Title:           request.Title,
		ShowBankAccount: request.ShowBankAccount,
		UserID:          userIdParsed,
		Number:          request.Number,
		Status:          request.Status,
		IPAddress:       request.RequestIp,
		RequestSource:   request.RequestSource,
		RequestSourceID: request.RequestSourceId,
		InvoiceViewLink: request.InvoiceLink,
		BusinessLogo:    request.BusinessDetail.Logo,
	}

	if request.PaymentTypes != nil {
		inv.AllowCard = request.PaymentTypes.AllowCards.ToBool()
	}
	if request.AccountDetails != nil {
		inv.AccountNumber = request.AccountDetails.AccountNumber
		inv.RoutingNumber = request.AccountDetails.RoutingNumber
	}
	inv.Created = request.Created
	return inv, nil
}

func transformProtoToInvoices(request *grpcSvc.InvoiceResponse) ([]*Invoice, error) {
	var responses []*Invoice
	for _, invoice := range request.Invoices {
		inv, err := transformProtoToInvoice(invoice)
		if err != nil {
			return responses, err
		}
		responses = append(responses, inv)
	}
	return responses, nil
}

func (is invoiceService) GetGrpcInvoiceByID(invoiceID id.InvoiceID) (*grpcSvc.FullInvoice, error) {
	var inv *grpcSvc.FullInvoice

	defer is.conn.CloseAndCancel()

	r := &grpcSvc.InvoiceIdGetRequest{
		InvoiceId: invoiceID.String(),
	}

	inv, err := is.invoiceServiceClient.GetInvoice(context.Background(), r)
	if err != nil {
		return inv, err
	}
	return inv, nil
}

func (is invoiceService) UpdateInvoiceStatus(invoiceID id.InvoiceID, businessID id.BusinessID, status grpcSvc.InvoiceRequestStatus) error {
	defer is.conn.CloseAndCancel()
	r := &grpcSvc.InvoiceRequest{
		Id:         invoiceID.String(),
		BusinessId: businessID.String(),
		Status:     status,
	}

	_, err := is.invoiceServiceClient.UpsertInvoiceRequest(context.Background(), r)
	if err != nil {
		return err
	}
	return nil
}

func (is invoiceService) ExecuteInvoiceAction(invoiceIDs []string, action grpcSvc.InvoiceAction) error {
	defer is.conn.CloseAndCancel()
	r := &grpcSvc.InvoiceActionRequest{
		InvoiceIds: invoiceIDs,
		Action:     action,
	}

	_, err := is.invoiceServiceClient.ExecuteAction(context.Background(), r)
	if err != nil {
		return err
	}
	return nil
}

// Follow create invoice and item is being used by shopify order creation only
func (is invoiceService) CreateInvoice(request *grpcSvc.InvoiceRequest) (*Invoice, error) {
	defer is.conn.CloseAndCancel()
	ir, err := is.invoiceServiceClient.UpsertInvoiceRequest(context.Background(), request)
	if err != nil {
		return nil, err
	}
	log.Println(fmt.Sprintf("new invoice id - %v", ir.Id))
	return transformProtoToInvoice(ir)
}

func (is invoiceService) CreateItem(request *grpcItem.Item) (*grpcItem.Item, error) {
	defer is.conn.CloseAndCancel()
	ir, err := is.itemServiceClient.UpsertItem(context.Background(), request)
	if err != nil {
		return nil, err
	}
	log.Println(fmt.Sprintf("new item id - %v", ir.Id))
	return ir, nil
}

func (is invoiceService) GetInvoiceIDFromPaymentRequestID(paymentRequestID *shared.PaymentRequestID) (*Invoice, error) {
	var invId id.InvoiceID
	if paymentRequestID != nil {
		pmReqIdStr := *paymentRequestID
		pmReqId, err := id.ParsePaymentRequestID(pmReqIdStr.ToPrefixString())
		if err != nil {
			return nil, err
		}
		invId, err = id.ParseInvoiceID(fmt.Sprintf("%s%s", id.IDPrefixInvoice, pmReqId.UUIDString()))
		if err != nil {
			return nil, err
		}
		return is.GetInvoiceByID(invId)
	}
	return nil, nil
}

func (is invoiceService) GetInvoiceCounts() (map[shared.BusinessID]int32, error) {
	defer is.conn.CloseAndCancel()
	empty := &emptypb.Empty{}
	result := map[shared.BusinessID]int32{}

	resp, err := is.invoiceServiceClient.GetInvoiceCountForBusinesses(context.Background(), empty)
	if err != nil {
		log.Println(err)
		return result, err
	}

	for _, count := range resp.Counts {
		busId, err := shared.ParseBusinessID(count.BusinessId)
		if err != nil {
			log.Println(err)
			return result, err
		}
		result[busId] = count.Count
	}
	return result, nil
}

func (is invoiceService) GetInvoiceAmounts() (num.Decimal, error) {
	defer is.conn.CloseAndCancel()
	empty := &emptypb.Empty{}
	result := num.NewZero()

	resp, err := is.invoiceServiceClient.GetInvoiceAmountForBusinesses(context.Background(), empty)
	if err != nil {
		log.Println(err)
		return result, err
	}

	for _, amount := range resp.InvoiceAmounts {
		amountParsed, err := num.ParseDecimal(amount.TotalAmount)
		if err != nil {
			log.Println(err)
			return num.NewZero(), err
		}
		result = result.Add(amountParsed)
	}
	return result, nil
}

func (is invoiceService) GetInvoiceAmountsWithFilter(businessID shared.BusinessID, startDate time.Time, endDate time.Time) (*InvoiceAmount, error) {
	defer is.conn.CloseAndCancel()
	result := &InvoiceAmount{
		BusinessID: businessID,
	}

	startDateParsed, err := grpcTypes.TimestampProto(startDate)
	if err != nil {
		return result, err
	}
	endDateParsed, err := grpcTypes.TimestampProto(endDate)
	if err != nil {
		return result, err
	}

	req := &grpcSvc.InvoiceAmountRequest{
		BusinessId: businessID.ToPrefixString(),
		DateRangeFilter: &grpcRoot.DateRange{
			Start:  startDateParsed,
			End:    endDateParsed,
			Filter: grpcRoot.DateRangeFilter_DRF_START_END,
		},
	}

	resp, err := is.invoiceServiceClient.GetInvoiceAmount(context.Background(), req)
	if err != nil {
		log.Println(err)
		return result, err
	}

	if resp.TotalPaid != "" {
		paidAmountParsed, err := shared.ParseDecimal(resp.TotalPaid)
		if err != nil {
			log.Println(err)
			return result, err
		}
		result.TotalPaid = paidAmountParsed
	}

	if resp.TotalRequested != "" {
		requestedAmountParsed, err := shared.ParseDecimal(resp.TotalRequested)
		if err != nil {
			log.Println(err)
			return result, err
		}
		result.TotalRequest = requestedAmountParsed
	}

	return result, nil
}

func getInvoiceRequestStatus(status string) grpcSvc.InvoiceRequestStatus {
	switch status {
	case "complete":
		return grpcSvc.InvoiceRequestStatus_IRT_PAID
	case "pending":
		return grpcSvc.InvoiceRequestStatus_IRT_OPEN
	default:
		return grpcSvc.InvoiceRequestStatus_IRT_UNSPECIFIED
	}
}
