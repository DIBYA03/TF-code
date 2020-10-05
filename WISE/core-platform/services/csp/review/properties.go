/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package review

import (
	"github.com/jmoiron/sqlx/types"
	"github.com/wiseco/core-platform/services"
	biz "github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/shared"
)

const (
	EntityTypeSoleProprietor = "soleProprietor"
)

//Response the response of business verification
type Response struct {
	Status        string                `json:"status"`
	ReviewItems   *services.StringArray `json:"reviewItems"`
	Notes         *types.JSONText       `json:"notes"`
	BusinessName  string                `json:"-"`
	BusinessOwner shared.UserID         `json:"-"`
	EntityType    *string               `json:"-"`
}

type ConsumerKYCResponse struct {
	Status      string                `json:"status"`
	ReviewItems *services.StringArray `json:"reviewItems"`
	Notes       *types.JSONText       `json:"notes"`
}

type ErrorResponse struct {
	Raw       error             `json:"-"`
	ErrorType VerificationError `json:"errorType"`
	Values    []string          `json:"values"`
	Business  *biz.Business     `json:"business"`
}

type ConsumerKYCError struct {
	Raw       error             `json:"-"`
	ErrorType VerificationError `json:"errorType"`
	Values    []string          `json:"values"`
}

func (v ConsumerKYCError) Error() string {
	return v.ErrorType.String()
}

func (v ErrorResponse) Error() string {
	return v.ErrorType.String()
}

func NewErrorResponse(raw VerificationError, values []string, b *biz.Business) ErrorResponse {
	return ErrorResponse{
		Raw:       raw,
		ErrorType: raw,
		Values:    values,
		Business:  b,
	}
}
