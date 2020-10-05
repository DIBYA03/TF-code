package pretoken

import (
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/wiseco/core-platform/services/csp/cspuser"
	"github.com/wiseco/core-platform/services/csp/services"
	"github.com/wiseco/go-lib/jwt"
)

func handleAuthentication(event events.CognitoEventUserPoolsPreTokenGen) (events.CognitoEventUserPoolsPreTokenGen, error) {
	cognitoID, ok := event.Request.UserAttributes["sub"]
	if !ok {
		return event, errors.New("invalid csp agent id")
	}

	// Get csp agent id
	agentID, err := cspuser.NewUserService(services.NewSourceRequest()).GetIdByCognitoID(cognitoID)
	if err != nil {
		log.Println(err)
		return event, errors.New("invalid csp agent id")
	} else if agentID.IsZero() {
		return event, errors.New("invalid csp agent id")
	}

	event.Response.ClaimsOverrideDetails.ClaimsToAddOrOverride = map[string]string{
		jwt.KeyAgentID:   agentID.String(),
		jwt.KeyTokenType: string(jwt.TokenTypeSupportCognito),
	}

	return event, nil
}

// HandleCases generates before JWT tokens are sent
func HandleCases(event events.CognitoEventUserPoolsPreTokenGen) (events.CognitoEventUserPoolsPreTokenGen, error) {

	eventType := event.TriggerSource

	switch eventType {
	case "TokenGeneration_HostedAuth":
		log.Println("TokenGeneration_HostedAuth")
		return handleAuthentication(event)
	case "TokenGeneration_Authentication":
		// Called after authentication is completed
		log.Println("TokenGeneration_Authentication")
		return handleAuthentication(event)
	case "TokenGeneration_NewPasswordChallenge":
		log.Println("TokenGeneration_NewPasswordChallenge")
	case "TokenGeneration_AuthenticateDevice":
		log.Println("TokenGeneration_AuthenticateDevice")
		return handleAuthentication(event)
	case "TokenGeneration_RefreshTokens":
		// Called on refresh token
		log.Println("TokenGeneration_RefreshTokens")
		return handleAuthentication(event)
	default:
		fmt.Println(eventType)
	}

	return event, nil
}

// HandleCognitoPreTokenRequest is used on pretoken
func HandleCognitoPreTokenRequest(event events.CognitoEventUserPoolsPreTokenGen) (events.CognitoEventUserPoolsPreTokenGen, error) {
	event, err := HandleCases(event)
	return event, err
}
