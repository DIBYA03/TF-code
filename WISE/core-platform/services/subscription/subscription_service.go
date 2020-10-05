package subscription

import (
	"errors"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/wiseco/core-platform/auth"
	"github.com/wiseco/core-platform/partner/service/segment"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/services/data"
	user "github.com/wiseco/core-platform/services/user"
	"github.com/wiseco/core-platform/shared"
)

type subscriptionDatastore struct {
	sourceReq services.SourceRequest
	*sqlx.DB
}

type SubscriptionService interface {
	Update(SubscriptionUpdate) (*business.Business, error)
}

func NewSubscriptionService(r services.SourceRequest) SubscriptionService {
	return &subscriptionDatastore{r, data.DBWrite}
}

func (db *subscriptionDatastore) Update(u SubscriptionUpdate) (*business.Business, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckUserAccess(u.UserID)
	if err != nil {
		return nil, err
	}

	_, ok := services.SubscriptionStatusValidator[u.SubscriptionStatus]
	if !ok {
		return nil, errors.New("invalid subscription status")
	}

	// Update subscription details in business
	decisionDate := time.Now()

	bsu := business.BusinessSubscriptionUpdate{
		ID:                       u.BusinessID,
		SubscriptionStatus:       &u.SubscriptionStatus,
		SubscriptionDecisionDate: &decisionDate,
	}

	if u.SubscriptionStatus == services.SubscriptionStatusActive {
		// Set subscription start dateto Aug 1. API will be deprecated by July end
		subDate := time.Date(2020, time.August, 1, 0, 0, 0, 0, time.UTC)
		subStartDate := shared.Date(subDate)
		bsu.SubscriptionStartDate = &subStartDate
	}

	b, err := business.NewBusinessService(db.sourceReq).UpdateSubscription(bsu)
	if err != nil {
		return nil, err
	}

	// Find user subscritpion status
	bs, err := business.NewBusinessService(db.sourceReq).List(0, 10, u.UserID)
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
		usr, err := user.NewUserService(db.sourceReq).UpdateSubscription(u.UserID, userSubscriptionStatus)
		if err != nil {
			return nil, err
		}

		// push to intercom
		segment.NewSegmentService().PushToAnalyticsQueue(u.UserID, segment.CategoryConsumer, segment.ActionSubscription, usr)
	}

	log.Println("Subscription status: ", u.SubscriptionStatus, u.UserID, u.BusinessID)

	return b, nil
}
