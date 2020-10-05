package service

import (
	"time"
)

type APIRequest struct {
	SourceRequestID string       `json:"sourceRequestId"`
	Partner         ProviderName `json:"partner"`
	UserID          EntityID     `json:"userId"`
	StartedAt       time.Time    `json:"startedAt"`
	SourceIP        string       `json:"sourceIP"`
	UserAgent       string       `json:"userAgent"`
	Message         string       `json:"message"`
	Elapsed         int64        `json:"elapsed"`
}
