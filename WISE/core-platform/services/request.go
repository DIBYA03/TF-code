/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package services

import (
	"log"
	"time"

	"github.com/google/uuid"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/partner/service"
	"github.com/wiseco/core-platform/shared"
)

type SourceRequest struct {
	JWT        string        `json:"jwt"`
	EnvID      string        `json:"envId"`
	RequestID  string        `json:"requestId"`
	GatewayID  string        `json:"gatewayId"`
	APIKey     string        `json:"apiKey"`
	SourceIP   string        `json:"sourceIP"`
	UserAgent  string        `json:"userAgent"`
	UserID     shared.UserID `json:"userId"`
	PoolID     string        `json:"poolId"`
	StartedAt  time.Time     `json:"startedAt"`
	AcceptLang string        `json:"acceptLang"`
}

func NewSourceRequest() SourceRequest {
	return SourceRequest{
		RequestID: uuid.New().String(),
		StartedAt: time.Now(),
		SourceIP:  shared.GetOutboundIP().String(),
	}
}

func (s *SourceRequest) PartnerBankRequest() partnerbank.APIRequest {
	return partnerbank.APIRequest{
		SourceRequestID: s.RequestID,
		Partner:         partnerbank.ProviderNameBBVA,
		StartedAt:       time.Now(),
		SourceIP:        s.SourceIP,
		UserAgent:       s.UserAgent,
	}
}

//CSPPartnerBankRequest use to explicity pass the user id within the request
func (s *SourceRequest) CSPPartnerBankRequest(userID shared.UserID) partnerbank.APIRequest {
	return partnerbank.APIRequest{
		SourceRequestID: s.RequestID,
		Partner:         partnerbank.ProviderNameBBVA,
		StartedAt:       time.Now(),
		SourceIP:        s.SourceIP,
		UserAgent:       s.UserAgent,
	}
}

func (s *SourceRequest) PartnerServiceRequest() service.APIRequest {
	return service.APIRequest{
		SourceRequestID: s.RequestID,
		Partner:         service.ProviderNameSendGrid,
		StartedAt:       time.Now(),
		SourceIP:        s.SourceIP,
		UserAgent:       s.UserAgent,
	}
}

//NoToCSPServiceRequest Notification to CSP source request
func NoToCSPServiceRequest(userID shared.UserID) SourceRequest {
	id := uuid.New()
	r := SourceRequest{
		UserID:    userID,
		RequestID: id.String(),
		StartedAt: time.Now(),
		SourceIP:  shared.GetOutboundIP().String(),
	}
	log.Printf("CSP Request object %v", r)
	return r

}
