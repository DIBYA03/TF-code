package bank

import (
	"time"

	"github.com/google/uuid"
	"github.com/wiseco/core-platform/shared"
)

type APIRequest struct {
	SourceRequestID string       `json:"sourceRequestId"`
	Partner         ProviderName `json:"partner"`
	StartedAt       time.Time    `json:"startedAt"`
	SourceIP        string       `json:"sourceIP"`
	UserAgent       string       `json:"userAgent"`
	Message         string       `json:"message"`
	Elapsed         int64        `json:"elapsed"`
}

type APIResponse struct {
	RequestID       string       `json:"requestId"`
	SourceRequestID string       `json:"sourceRequestId"`
	Partner         ProviderName `json:"partner"`
	StartedAt       time.Time    `json:"startedAt"`
	SourceIP        string       `json:"sourceIP"`
	UserAgent       string       `json:"userAgent"`
	Message         string       `json:"message"`
	Elapsed         int64        `json:"elapsed"`
}

func (r *APIRequest) Duration() int64 {
	duration := time.Since(r.StartedAt) / time.Millisecond
	return int64(duration)
}

func (r *APIRequest) New() APIRequest {
	return APIRequest{
		SourceRequestID: r.SourceRequestID,
		Partner:         r.Partner,
		StartedAt:       r.StartedAt,
		SourceIP:        r.SourceIP,
		UserAgent:       r.UserAgent,
		Message:         r.Message,
		Elapsed:         r.Elapsed,
	}
}

func NewAPIRequest() APIRequest {
	return APIRequest{
		SourceRequestID: uuid.New().String(),
		StartedAt:       time.Now(),
		SourceIP:        shared.GetOutboundIP().String(),
	}
}
