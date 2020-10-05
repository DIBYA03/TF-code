package pretoken

import (
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/wiseco/core-platform/api/auth"
	idsrv "github.com/wiseco/core-platform/identity"
	usersrv "github.com/wiseco/core-platform/services/user"
	"github.com/wiseco/go-lib/jwt"
)

func handleAuthentication(event events.CognitoEventUserPoolsPreTokenGen) (events.CognitoEventUserPoolsPreTokenGen, error) {
	providerID := idsrv.ProviderID(event.UserName)
	identity, err := idsrv.NewIdentityService(auth.NewPreTokenGenIdentitySourceRequest(event)).GetByProviderID(
		providerID, idsrv.ProviderNameCognito, idsrv.ProviderSource(event.UserPoolID),
	)
	if err != nil {
		return event, err
	}

	user, err := usersrv.NewUserService(auth.NewPreTokenGenSourceRequest(event)).GetUserWithIdentity(identity.ID)
	if err != nil {
		return event, err
	}

	event.Response.ClaimsOverrideDetails.ClaimsToAddOrOverride = map[string]string{
		jwt.KeyIdentityID: string(identity.ID),
		jwt.KeyUserID:     user.ID.ToPrefixString(),
		jwt.KeyConsumerID: user.ConsumerID.ToPrefixString(),
		jwt.KeyTokenType:  string(jwt.TokenTypeClientCognito),
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
