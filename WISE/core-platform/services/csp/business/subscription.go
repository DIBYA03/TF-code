package business

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/partner/service/segment"
	"github.com/wiseco/core-platform/services"
	bussrv "github.com/wiseco/core-platform/services/business"
	consrv "github.com/wiseco/core-platform/services/csp/consumer"
	"github.com/wiseco/core-platform/services/csp/cspuser"
	cspusrsrv "github.com/wiseco/core-platform/services/csp/cspuser"
	"github.com/wiseco/core-platform/services/csp/data"
	cspsrv "github.com/wiseco/core-platform/services/csp/services"
	usrsrv "github.com/wiseco/core-platform/services/user"
	"github.com/wiseco/core-platform/shared"
)

// SubscriptionService ..
type SubscriptionService interface {
	Update(SubscriptionUpdate) (*Subscription, error)
	GetByBusinessID(businessID shared.BusinessID) (*Subscription, error)
}

type subscriptionService struct {
	rdb *sqlx.DB
	wdb *sqlx.DB
	sr  cspsrv.SourceRequest
}

// NewSubscriptionService ..
func NewSubscriptionService(sr cspsrv.SourceRequest) SubscriptionService {
	return subscriptionService{wdb: data.DBWrite, rdb: data.DBRead, sr: sr}
}

func (srv subscriptionService) Update(u SubscriptionUpdate) (*Subscription, error) {
	if u.SubscriptionStatus == nil {
		return nil, errors.New("subscription status is required")
	}

	_, ok := services.SubscriptionStatusValidator[*u.SubscriptionStatus]
	if !ok {
		return nil, errors.New("invalid subscription status")
	}

	if *u.SubscriptionStatus == services.SubscriptionStatusActive && u.SubscriptionStartDate == nil {
		return nil, errors.New("subscription start date is required")
	}

	decisionDate := time.Now()

	// Update subscription details in business
	bsu := bussrv.BusinessSubscriptionUpdate{
		ID:                       u.BusinessID,
		SubscriptionStatus:       u.SubscriptionStatus,
		SubscriptionStartDate:    u.SubscriptionStartDate,
		SubscriptionDecisionDate: &decisionDate,
	}

	sr := services.NewSourceRequest()

	b, err := New(sr).ByID(u.BusinessID)
	if err != nil {
		return nil, err
	}

	sr.UserID = b.OwnerID

	_, err = bussrv.NewBusinessService(sr).UpdateSubscription(bsu)
	if err != nil {
		return nil, err
	}

	// Find user subscritpion status
	bs, err := bussrv.NewBusinessService(sr).List(0, 10, b.OwnerID)
	if err != nil {
		return nil, err
	}

	var userSubscriptionStatus services.SubscriptionStatus
	for _, b := range bs {
		if b.SubscriptionStatus == nil {
			continue
		}

		// If atleast one business is active set status to active
		if *b.SubscriptionStatus == services.SubscriptionStatusActive {
			userSubscriptionStatus = services.SubscriptionStatusActive
			break
		} else if *b.SubscriptionStatus == services.SubscriptionStatusUnpaid {
			userSubscriptionStatus = services.SubscriptionStatusUnpaid
		}
	}

	// Update subscription details in user
	if userSubscriptionStatus != "" {
		usr, err := usrsrv.NewUserService(sr).UpdateSubscription(b.OwnerID, userSubscriptionStatus)
		if err != nil {
			return nil, err
		}

		// push to intercom
		segment.NewSegmentService().PushToAnalyticsQueue(usr.ID, segment.CategoryConsumer, segment.ActionSubscription, usr)
	}

	// Update CSP operator ID
	if srv.sr.CognitoID != "" {
		agentID, err := cspuser.NewUserService(srv.sr).ByCognitoID(srv.sr.CognitoID)
		if err != nil {
			return nil, err
		}

		err = NewCSPService().UpdateSubscribedAgentID(u.BusinessID, agentID)
		if err != nil {
			return nil, err
		}
	}

	return srv.GetByBusinessID(u.BusinessID)
}

func (srv subscriptionService) GetByBusinessID(businessID shared.BusinessID) (*Subscription, error) {
	s := Subscription{
		BusinessID: businessID,
	}

	cb, err := NewCSPService().ByBusinessID(businessID)
	if err != nil {
		return nil, err
	}

	// Get CSP operator name
	if cb.SubscribedAgentID != nil {
		u, err := cspusrsrv.NewUserService(srv.sr).GetByIdInternal(*cb.SubscribedAgentID)
		if err != nil {
			return nil, err
		}

		name := u.FirstName
		if u.MiddleName != "" {
			name = name + " " + u.MiddleName
		}
		name = name + " " + u.LastName
		s.SubscribedAgentName = &name
	}

	// Get business subscription status
	b, err := New(services.NewSourceRequest()).ByID(businessID)
	if err != nil {
		return nil, err
	}
	s.BusinessSubscriptionStatus = b.SubscriptionStatus
	s.SubscriptionDecisionDate = b.SubscriptionDecisionDate
	s.SubscriptionStartDate = b.SubscriptionStartDate

	// Get user subscription status
	u, err := consrv.New().ByUserID(string(b.OwnerID))
	if err != nil {
		return nil, err
	}
	s.UserSubscriptionStatus = u.SubscriptionStatus
	s.UserID = shared.UserID(u.ID)

	return &s, nil
}
