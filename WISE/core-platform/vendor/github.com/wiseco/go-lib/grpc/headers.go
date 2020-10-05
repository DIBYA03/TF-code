package grpc

import (
	"context"

	"github.com/wiseco/go-lib/http"
	"github.com/wiseco/go-lib/id"
)

type Header interface {
	Get(key string) string

	// Request id
	GetRequestID() string

	// Auth
	GetAuthorization() string
	GetSecretKey() string

	// User auth
	GetAuthClientKey() id.ClientKey
	GetAuthUserID() id.UserID
	GetAuthConsumerID() id.ConsumerID

	// Agent auth
	GetAgentName() string
	GetAgentEmail() string
	GetAgentID() string

	// IP headers
	GetWsRealIP() string
	GetForwardedFor() string
	GetForwardedHost() string

	// General request headers
	GetWsIdempotencyKey() string
	GetWsBusinessID() id.BusinessID
	GetWsConsumerID() id.ConsumerID
}

type header struct {
	ctx context.Context
}

func NewHeader(ctx context.Context) Header {
	return &header{ctx}
}

func (h *header) Get(key string) string {
	return HeaderValueFromIncoming(h.ctx, key)
}

func (h *header) GetRequestID() string {
	return h.Get(http.HeaderRequestID)
}

func (h *header) GetAuthorization() string {
	return h.Get(http.HeaderAuthorization)
}

func (h *header) GetSecretKey() string {
	return h.Get(http.HeaderSecretKey)
}

func (h *header) GetAuthClientKey() id.ClientKey {
	key, _ := id.ParseClientKey(HeaderValueFromIncoming(h.ctx, http.HeaderAuthClientKey))
	return key
}

func (h *header) GetAuthUserID() id.UserID {
	userID, _ := id.ParseUserID(HeaderValueFromIncoming(h.ctx, http.HeaderAuthUserID))
	return userID
}

func (h *header) GetAuthConsumerID() id.ConsumerID {
	conID, _ := id.ParseConsumerID(HeaderValueFromIncoming(h.ctx, http.HeaderAuthConsumerID))
	return conID
}

func (h *header) GetAgentName() string {
	return h.Get(http.HeaderAgentName)
}

func (h *header) GetAgentEmail() string {
	return h.Get(http.HeaderAgentEmail)
}

func (h *header) GetAgentID() string {
	return h.Get(http.HeaderAgentID)
}

func (h *header) GetWsRealIP() string {
	return h.Get(http.HeaderRealIP)
}

func (h *header) GetForwardedFor() string {
	return h.Get(http.HeaderForwardedFor)
}

func (h *header) GetForwardedHost() string {
	return h.Get(http.HeaderForwardedHost)
}

func (h *header) GetWsIdempotencyKey() string {
	return h.Get(http.HeaderIdempotency)
}

func (h *header) GetWsBusinessID() id.BusinessID {
	wsBusID, _ := id.ParseBusinessID(HeaderValueFromIncoming(h.ctx, http.HeaderBusinessID))
	return wsBusID
}

func (h *header) GetWsConsumerID() id.ConsumerID {
	wsConID, _ := id.ParseConsumerID(HeaderValueFromIncoming(h.ctx, http.HeaderConsumerID))
	return wsConID
}
