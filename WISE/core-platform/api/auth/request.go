package auth

import (
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/wiseco/core-platform/identity"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/shared"
)

// TODO: Parse user id and return error
func NewPostConfirmSourceRequest(e events.CognitoEventUserPoolsPostConfirmation) services.SourceRequest {
	return services.SourceRequest{
		RequestID: uuid.New().String(),
		UserID:    shared.UserID(e.Request.UserAttributes[CognitoSub]),
		PoolID:    e.UserPoolID,
		StartedAt: time.Now(),
	}
}

func NewPreTokenGenSourceRequest(e events.CognitoEventUserPoolsPreTokenGen) services.SourceRequest {
	return services.SourceRequest{
		RequestID: uuid.New().String(),
		UserID:    shared.UserID(e.Request.UserAttributes[CognitoSub]),
		PoolID:    e.UserPoolID,
		StartedAt: time.Now(),
	}
}

func NewPreSignUpSourceRequest(e events.CognitoEventUserPoolsPreSignup) services.SourceRequest {
	return services.SourceRequest{
		RequestID: uuid.New().String(),
		UserID:    shared.UserID(e.Request.UserAttributes[CognitoSub]),
		PoolID:    e.UserPoolID,
		StartedAt: time.Now(),
	}
}

func NewPostConfirmIdentitySourceRequest(e events.CognitoEventUserPoolsPostConfirmation) identity.SourceRequest {
	return identity.SourceRequest{
		RequestID: uuid.New().String(),
		UserID:    shared.UserID(e.Request.UserAttributes[CognitoSub]),
		PoolID:    e.UserPoolID,
		StartedAt: time.Now(),
	}
}

func NewPreTokenGenIdentitySourceRequest(e events.CognitoEventUserPoolsPreTokenGen) identity.SourceRequest {
	return identity.SourceRequest{
		RequestID: uuid.New().String(),
		UserID:    shared.UserID(e.Request.UserAttributes[CognitoSub]),
		PoolID:    e.UserPoolID,
		StartedAt: time.Now(),
	}
}

func NewPreSignUpIdentitySourceRequest(e events.CognitoEventUserPoolsPreSignup) identity.SourceRequest {
	return identity.SourceRequest{
		RequestID: uuid.New().String(),
		UserID:    shared.UserID(e.Request.UserAttributes[CognitoSub]),
		PoolID:    e.UserPoolID,
		StartedAt: time.Now(),
	}
}
