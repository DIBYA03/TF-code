/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/wiseco/core-platform/shared"
)

type SourceRequest struct {
	JWT        string    `json:"jwt"`
	EnvId      string    `json:"envId"`
	RequestId  string    `json:"requestId"`
	APIId      string    `json:"apiId"`
	APIKey     string    `json:"apiKey"`
	SourceIP   string    `json:"sourceIP"`
	UserAgent  string    `json:"userAgent"`
	UserId     string    `json:"userId"`
	PoolId     string    `json:"poolId"`
	StartedAt  time.Time `json:"startedAt"`
	AcceptLang string    `json:"acceptLang"`
	CognitoID  string
}

func NewSourceRequest() SourceRequest {
	return SourceRequest{
		RequestId: uuid.New().String(),
		StartedAt: time.Now(),
		SourceIP:  shared.GetOutboundIP().String(),
	}
}

func NewSRRequest(cognitoID string) SourceRequest {
	return SourceRequest{
		RequestId: uuid.New().String(),
		StartedAt: time.Now(),
		SourceIP:  shared.GetOutboundIP().String(),
		CognitoID: cognitoID,
	}
}
