/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package api

import (
	"time"

	"github.com/wiseco/core-platform/identity"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/shared"
)

type APIRequest struct {
	ResourcePath                    string              `json:"resourcePath"` // The url path for the caller
	HTTPMethod                      string              `json:"httpMethod"`
	Headers                         map[string]string   `json:"headers"`
	MultiValueHeaders               map[string][]string `json:"multiValueHeaders"`
	QueryStringParameters           map[string]string   `json:"queryStringParameters"`
	MultiValueQueryStringParameters map[string][]string `json:"multiValueQueryStringParameters"`
	PathParameters                  map[string]string   `json:"pathParameters"`
	EnvVariables                    map[string]string   `json:"envVariables"`
	Body                            string              `json:"body"`
	IsBase64Encoded                 bool                `json:"isBase64Encoded,omitempty"`
	EnvID                           string              `json:"envId"`
	RequestID                       string              `json:"requestId"`
	GatewayID                       string              `json:"gatewayId"`
	APIKey                          string              `json:"apiKey"`
	SourceIP                        string              `json:"sourceIP"`
	UserAgent                       string              `json:"userAgent"`
	UserID                          shared.UserID       `json:"userId"`
	BusinessID                      *shared.BusinessID  `json:"businessId"`
	PoolID                          string              `json:"poolId"`
	StartedAt                       time.Time           `json:"startedAt"`
	CognitoID                       string
}

func (req *APIRequest) Duration() int64 {
	duration := time.Since(req.StartedAt) / time.Millisecond
	return int64(duration)
}

func (req *APIRequest) SourceRequest() services.SourceRequest {
	lang := req.Headers["Accept-Language"]
	if lang == "" {
		lang = "en-US"
	}
	return services.SourceRequest{
		JWT:        req.Headers["Authorization"],
		EnvID:      req.EnvID,
		RequestID:  req.RequestID,
		GatewayID:  req.GatewayID,
		APIKey:     req.APIKey,
		SourceIP:   req.SourceIP,
		UserID:     req.UserID,
		PoolID:     req.PoolID,
		StartedAt:  req.StartedAt,
		AcceptLang: lang,
	}
}

func (req *APIRequest) IdentitySourceRequest() identity.SourceRequest {
	lang := req.Headers["Accept-Language"]
	if lang == "" {
		lang = "en-US"
	}
	return identity.SourceRequest{
		JWT:        req.Headers["Authorization"],
		EnvID:      req.EnvID,
		RequestID:  req.RequestID,
		GatewayID:  req.GatewayID,
		APIKey:     req.APIKey,
		SourceIP:   req.SourceIP,
		UserID:     req.UserID,
		PoolID:     req.PoolID,
		StartedAt:  req.StartedAt,
		AcceptLang: lang,
	}
}
