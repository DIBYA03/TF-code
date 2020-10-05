package auth

import (
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/wiseco/core-platform/services/csp/services"
)

func NewPostConfirmSourceRequest(e events.CognitoEventUserPoolsPostConfirmation) services.SourceRequest {
	return services.SourceRequest{
		RequestId: uuid.New().String(),
		UserId:    e.Request.UserAttributes[CognitoSub],
		PoolId:    e.UserPoolID,
		StartedAt: time.Now(),
	}
}

func NewPreSignUpSourceRequest(e events.CognitoEventUserPoolsPreSignup) services.SourceRequest {
	return services.SourceRequest{
		RequestId: uuid.New().String(),
		UserId:    e.Request.UserAttributes[CognitoSub],
		PoolId:    e.UserPoolID,
		StartedAt: time.Now(),
	}
}

func NewPreAuthenticationSourceRequest(e events.CognitoEventUserPoolsPreAuthentication) services.SourceRequest {
	return services.SourceRequest{
		RequestId: uuid.New().String(),
		UserId:    e.Request.UserAttributes[CognitoSub],
		PoolId:    e.UserPoolID,
		StartedAt: time.Now(),
	}
}
