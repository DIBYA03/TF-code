/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/
package identity

import (
	"time"

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
