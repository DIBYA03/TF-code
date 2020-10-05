package jwt

const (
	// Token types
	TokenTypeClientAccess   = "client_access"
	TokenTypeClientRefresh  = "client_refresh"
	TokenTypeClientCognito  = "client_cognito"
	TokenTypeSupportCognito = "support_cognito"

	// Claim Keys
	KeyTokenType  = "token_type"
	KeyIdentityID = "identity_id"
	KeyUserID     = "user_id"
	KeyConsumerID = "consumer_id"
	KeyAgentID    = "agent_id"
)
