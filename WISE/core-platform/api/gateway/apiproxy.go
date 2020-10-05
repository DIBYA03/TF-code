/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package gateway

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/wiseco/core-platform/api"
	idsrv "github.com/wiseco/core-platform/identity"
	"github.com/wiseco/core-platform/services"
	usersrv "github.com/wiseco/core-platform/services/user"
	"github.com/wiseco/core-platform/shared"
)

const WiseBusinessID = "X-Wise-Business-ID"

const AuthorizerClaims = "claims"

const (
	AuthorizerClaimKeyAUD           = "aud"
	AuthorizerClaimKeyAuthTime      = "auth_time"
	AuthorizerClaimKeyUsername      = "cognito:username"
	AuthorizerClaimKeyEventId       = "event_id"
	AuthorizerClaimKeyExp           = "exp"
	AuthorizerClaimKeyIAT           = "iat"
	AuthorizerClaimKeyISS           = "iss"
	AuthorizerClaimKeyPhone         = "phone_number"
	AuthorizerClaimKeyPhoneVerified = "phone_number_verified"
	AuthorizerClaimKeySub           = "sub"
	AuthorizerClaimKeyTokenUse      = "token_use"
	AuthorizerClaimKeyIdentityID    = "identity_id"
	AuthorizerClaimKeyUserID        = "user_id"
)

func NewAPIRequest(req events.APIGatewayProxyRequest) api.APIRequest {
	var c map[string]interface{} = req.RequestContext.Authorizer[AuthorizerClaims].(map[string]interface{})
	cognitoID, ok := c[AuthorizerClaimKeySub].(string)
	if !ok {
		cognitoID = c[AuthorizerClaimKeyUsername].(string)
	}

	poolID := os.Getenv("COGNITO_USER_POOL_ID")
	if poolID == "" {
		poolURL, ok := c[AuthorizerClaimKeyISS].(string)
		if ok {
			urlParts := strings.Split(poolURL, "/")
			if len(urlParts) > 0 {
				poolID = urlParts[len(urlParts)-1]
			}
		}
	}

	var userID shared.UserID
	var err error

	// Check user id key
	userIDStr, ok := c[AuthorizerClaimKeyUserID].(string)
	if ok {
		userID, err = shared.ParseUserID(userIDStr)
	} else {
		// Find using cognito id
		ir := idsrv.SourceRequest{
			RequestID: uuid.New().String(),
			StartedAt: time.Now(),
		}

		identity, err := idsrv.NewIdentityService(ir).GetByProviderID(
			idsrv.ProviderID(cognitoID), idsrv.ProviderNameCognito, idsrv.ProviderSource(poolID),
		)
		if err == nil {
			sr := services.SourceRequest{
				RequestID: uuid.New().String(),
				StartedAt: time.Now(),
			}

			id, err := usersrv.NewUserService(sr).GetUserIDWithIdentity(identity.ID)
			if err == nil {
				userID = *id
			}
		}
	}

	// Log error and continue with user id as empty string
	// Downstream auth access checks will fail which is ok
	if err != nil {
		log.Println(err)
	}

	r := api.APIRequest{
		ResourcePath:                    req.Resource,
		HTTPMethod:                      req.HTTPMethod,
		Headers:                         req.Headers,
		MultiValueHeaders:               req.MultiValueHeaders,
		QueryStringParameters:           req.QueryStringParameters,
		MultiValueQueryStringParameters: req.MultiValueQueryStringParameters,
		PathParameters:                  req.PathParameters,
		EnvVariables:                    req.StageVariables,
		Body:                            req.Body,
		IsBase64Encoded:                 req.IsBase64Encoded,
		EnvID:                           req.RequestContext.Stage,
		RequestID:                       req.RequestContext.RequestID,
		GatewayID:                       req.RequestContext.APIID,
		APIKey:                          req.RequestContext.Identity.APIKey,
		SourceIP:                        req.RequestContext.Identity.SourceIP,
		UserAgent:                       req.RequestContext.Identity.UserAgent,
		UserID:                          userID,
		PoolID:                          poolID,
		StartedAt:                       time.Now(),
	}

	businessID, err := shared.ParseBusinessID(r.SingleHeaderValue(WiseBusinessID))
	if err == nil {
		r.BusinessID = &businessID
	}

	return r
}

//NewCSPAPIRequest CSP Temp request
func NewCSPAPIRequest(req events.APIGatewayProxyRequest) api.APIRequest {
	var c map[string]interface{} = req.RequestContext.Authorizer[AuthorizerClaims].(map[string]interface{})
	cognitoID, ok := c[AuthorizerClaimKeySub].(string)
	if !ok {
		cognitoID = c[AuthorizerClaimKeyUsername].(string)
	}

	poolURL, ok := c[AuthorizerClaimKeyISS].(string)
	var poolID string
	if ok {
		urlParts := strings.Split(poolURL, "/")
		if len(urlParts) > 0 {
			poolID = urlParts[len(urlParts)-1]
		}
	}

	return api.APIRequest{
		ResourcePath:                    req.Resource,
		HTTPMethod:                      req.HTTPMethod,
		Headers:                         req.Headers,
		MultiValueHeaders:               req.MultiValueHeaders,
		QueryStringParameters:           req.QueryStringParameters,
		MultiValueQueryStringParameters: req.MultiValueQueryStringParameters,
		PathParameters:                  req.PathParameters,
		EnvVariables:                    req.StageVariables,
		Body:                            req.Body,
		IsBase64Encoded:                 req.IsBase64Encoded,
		EnvID:                           req.RequestContext.Stage,
		RequestID:                       req.RequestContext.RequestID,
		GatewayID:                       req.RequestContext.APIID,
		APIKey:                          req.RequestContext.Identity.APIKey,
		SourceIP:                        req.RequestContext.Identity.SourceIP,
		UserAgent:                       req.RequestContext.Identity.UserAgent,
		CognitoID:                       cognitoID,
		PoolID:                          poolID,
		StartedAt:                       time.Now(),
	}
}

func ProxyResponse(resp api.APIResponse) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode:        resp.StatusCode,
		Headers:           resp.Headers,
		MultiValueHeaders: resp.MultiValueHeaders,
		Body:              resp.Body,
		IsBase64Encoded:   resp.IsBase64Encoded,
	}
}
