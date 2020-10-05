package order

import (
	"database/sql"
	"log"
	"os"
	"strconv"

	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/services/contact"
	"github.com/wiseco/core-platform/services/payment"
)

type shopifyDatastore struct {
	sourceReq services.SourceRequest
}

type ShopifyOrderService interface {
	HandleShopifyOrder(ShopifyOrderMessage) error
}

func NewShopifyOrderService(r services.SourceRequest) ShopifyOrderService {
	return &shopifyDatastore{r}
}

func (s shopifyDatastore) HandleShopifyOrder(msg ShopifyOrderMessage) error {

	switch msg.EventType {
	case EventTypeOrderCreated:
		return s.onOrderCreated(msg)
	case EventTypeOrderPaid:
		return s.onOrderPaid(msg)
	case EventTypeOrderCanceled:
		return s.onOrderCanceled(msg)
	default:
		break
	}

	return nil
}

func (s shopifyDatastore) onOrderCanceled(msg ShopifyOrderMessage) error {
	sr := services.NewSourceRequest()

	req, err := payment.NewRequestService(sr).GetByRequestSourceIDInternal(msg.ID, msg.BusinessID)
	if err != nil {
		return err
	}

	// If already canceled return
	if req.Status != nil && *req.Status == payment.PaymentRequestStatusCanceled {
		return nil
	}

	reqUpdate := payment.RequestStatusUpdate{
		ID:         req.ID,
		BusinessID: req.BusinessID,
		Status:     payment.PaymentRequestStatusCanceled,
	}

	sr.UserID = req.CreatedUserID

	req, err = payment.NewRequestService(sr).UpdateRequestStatus(&reqUpdate)
	if err != nil {
		return err
	}

	return nil
}

func (s shopifyDatastore) onOrderPaid(msg ShopifyOrderMessage) error {
	sr := services.NewSourceRequest()

	req, err := payment.NewRequestService(sr).GetByRequestSourceIDInternal(msg.ID, msg.BusinessID)
	if err != nil {
		return err
	}

	// If already paid return
	if req.Status != nil && *req.Status == payment.PaymentRequestStatusComplete {
		return nil
	}

	reqUpdate := payment.RequestStatusUpdate{
		ID:         req.ID,
		BusinessID: req.BusinessID,
		Status:     payment.PaymentRequestStatusComplete,
	}

	sr.UserID = req.CreatedUserID

	req, err = payment.NewRequestService(sr).UpdateRequestStatus(&reqUpdate)
	if err != nil {
		return err
	}

	return nil
}

func (s shopifyDatastore) onOrderCreated(msg ShopifyOrderMessage) error {
	b, err := business.NewBusinessServiceWithout().GetByIdInternal(msg.BusinessID)
	if err != nil {
		return err
	}

	// Check if contact exists
	c, err := contact.NewContactServiceWithout().GetByPhoneEmailInternal(msg.Phone, msg.Email, msg.BusinessID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	sr := services.NewSourceRequest()
	sr.UserID = b.OwnerID

	if err != nil && err == sql.ErrNoRows {
		// Create contact if not found
		create := contact.ContactCreate{
			BusinessID: msg.BusinessID,
			Type:       contact.ContactTypePerson,
			UserID:     b.OwnerID,
		}

		category := contact.ContactCategoryShopify
		create.Category = &category

		if msg.Phone != nil {
			create.PhoneNumber = *msg.Phone
		}

		if msg.Email != nil {
			create.Email = *msg.Email
		}

		if msg.FirstName == nil && msg.LastName == nil {
			if msg.Email != nil {
				create.FirstName = msg.Email
			} else if msg.Phone != nil {
				create.PhoneNumber = *msg.Phone
			}
		} else {
			if msg.FirstName != nil {
				create.FirstName = msg.FirstName
			}

			if msg.LastName != nil {
				create.LastName = msg.LastName
			}
		}

		c, err = contact.NewContactService(sr).Create(&create)
		if err != nil {
			return err
		}
	}

	amt, err := strconv.ParseFloat(msg.Amount, 64)
	if err != nil {
		return err
	}

	notes := "Invoice for shopify order #" + strconv.FormatInt(msg.OrderID, 10)
	source := payment.RequestSourceShopify
	orderID := msg.ID.UUIDString()

	req := payment.RequestInitiate{
		BusinessID:      b.ID,
		CreatedUserID:   b.OwnerID,
		Amount:          amt,
		ContactID:       &c.ID,
		Currency:        payment.CurrencyUSD,
		IPAddress:       &sr.SourceIP,
		RequestSource:   &source,
		RequestSourceID: &orderID,
		Notes:           &notes,
	}

	maxCardOnlineAmount, err := strconv.ParseFloat(os.Getenv("CARD_ONLINE_MAX_MONEY_REQUEST_ALLOWED"), 64)
	if err != nil {
		log.Println("CARD_ONLINE_MAX_MONEY_REQUEST_ALLOWED is missing")
		return err
	}
	if amt > maxCardOnlineAmount {
		req.RequestType = payment.PaymentRequestTypeInvoiceNone
	} else {
		req.RequestType = payment.PaymentRequestTypeInvoiceCardAndBank
	}

	_, err = payment.NewRequestService(sr).Request(&req)
	if err != nil {
		return err
	}

	return nil
}
