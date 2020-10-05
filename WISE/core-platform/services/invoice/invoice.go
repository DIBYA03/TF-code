package invoice

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/go-lib/id"
	"github.com/wiseco/go-lib/num"
	grpcSvc "github.com/wiseco/protobuf/golang/invoice"
)

type Invoice struct {
	InvoiceID         id.InvoiceID
	BusinessID        id.BusinessID
	ContactID         id.ContactID
	Amount            num.Decimal
	Notes             string
	Title             string
	Number            int64
	AllowCard         bool
	AllowBankTransfer bool
	ShowBankAccount   bool
	UserID            id.UserID
	Created           *timestamp.Timestamp
	Status            grpcSvc.InvoiceRequestStatus
	AccountNumber     string
	RoutingNumber     string
	IPAddress         string
	RequestSource     string
	RequestSourceID   string
	InvoiceViewLink   string
	BusinessLogo      string
}

type InvoiceAmount struct {
	BusinessID   shared.BusinessID
	TotalRequest shared.Decimal
	TotalPaid    shared.Decimal
}

func (i *Invoice) GetCreatedTime() (time.Time, error) {
	return ptypes.Timestamp(i.Created)
}
