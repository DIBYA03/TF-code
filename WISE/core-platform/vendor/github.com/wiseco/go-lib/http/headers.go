package http

// Use lowercase keys for GRPC header compat
const (
	// Server generated request header
	HeaderRequestID = "ws-request-id"

	// Auth headers
	HeaderAuthorization = "authorization"
	HeaderSecretKey     = "ws-secret-key"

	// Client auth headers
	HeaderAuthClientKey  = "wsa-client-key"
	HeaderAuthUserID     = "wsa-user-id"
	HeaderAuthConsumerID = "wsa-consumer-id"

	// Support
	HeaderAgentName  = "wsa-agent-name"
	HeaderAgentEmail = "wsa-agent-email"
	HeaderAgentID    = "wsa-agent-id"

	// IP headers
	HeaderForwardedFor  = "x-forwarded-for"
	HeaderForwardedHost = "x-forwarded-host"
	HeaderRealIP        = "ws-real-ip"

	// Endpoint request headers
	HeaderIdempotency = "ws-idempotency-key"
	HeaderBusinessID  = "ws-business-id"
	HeaderConsumerID  = "ws-consumer-id"

	// Legacy headers
	HeaderLegacyBusinessID = "x-business-id"

	// CORS
	HeaderCorsAllowOrigin      = "access-control-allow-origin"
	HeaderCorsAllowCredentials = "access-control-allow-credentials"
)
