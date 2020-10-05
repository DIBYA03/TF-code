package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/ssm"

	"golang.org/x/crypto/ssh/terminal"
)

// Used for params from SSM parameter store later
var clientID, userPoolID string

func main() {
	sess := session.New()

	username, password, environment := getUserInput()
	ssmParameters(sess, environment)

	svc := cognitoidentityprovider.New(sess)

	idpAuth := initiateAuth(svc, username, password)

	respondToAuthChallenge(svc, idpAuth)
}

func initiateAuth(svc *cognitoidentityprovider.CognitoIdentityProvider, username string, password string) *cognitoidentityprovider.AdminInitiateAuthOutput {
	input := &cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow: aws.String("ADMIN_NO_SRP_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(username),
			"PASSWORD": aws.String(password),
		},
		ClientId:   aws.String(clientID),
		UserPoolId: aws.String(userPoolID),
	}

	idpResults, err := svc.AdminInitiateAuth(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			panic(aerr)
		} else {
			panic(err.Error())
		}
	}

	return idpResults
}

func respondToAuthChallenge(svc *cognitoidentityprovider.CognitoIdentityProvider, idpAuth *cognitoidentityprovider.AdminInitiateAuthOutput) {
	mfaToken := getMFAToken()

	input := &cognitoidentityprovider.AdminRespondToAuthChallengeInput{
		ChallengeName: aws.String("SMS_MFA"),
		ChallengeResponses: map[string]*string{
			"USERNAME":     idpAuth.ChallengeParameters["USER_ID_FOR_SRP"],
			"SMS_MFA_CODE": aws.String(mfaToken),
		},
		ClientId:   aws.String(clientID),
		Session:    idpAuth.Session,
		UserPoolId: aws.String(userPoolID),
	}

	idpChallenge, err := svc.AdminRespondToAuthChallenge(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			panic(aerr)
		} else {
			panic(err.Error())
		}
	}

	fmt.Println(idpChallenge)
}

func ssmParameters(sess *session.Session, environment string) {
	svc := ssm.New(sess)

	input := &ssm.GetParametersInput{
		Names: []*string{
			aws.String(fmt.Sprintf("/%s/cognito/user_pool/id", environment)),
			aws.String(fmt.Sprintf("/%s/cognito/user_pool/client/admin/id", environment)),
		},
		WithDecryption: aws.Bool(true),
	}

	paramResults, err := svc.GetParameters(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			panic(aerr)
		} else {
			panic(err.Error())
		}
	}

	for _, param := range paramResults.Parameters {
		if *param.Name == fmt.Sprintf("/%s/cognito/user_pool/id", environment) {
			userPoolID = *param.Value
		}
		if *param.Name == fmt.Sprintf("/%s/cognito/user_pool/client/admin/id", environment) {
			clientID = *param.Value
		}
	}
}

func getUserInput() (string, string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Environment: ")
	environment, _ := reader.ReadString('\n')
	environment = strings.Replace(environment, "\n", "", -1)

	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.Replace(username, "\n", "", -1)

	fmt.Print("Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)

	// Cheat way to make it look purty
	fmt.Print("\n")

	return username, password, environment
}

func getMFAToken() string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("MFA Token: ")
	mfaToken, _ := reader.ReadString('\n')
	mfaToken = strings.Replace(mfaToken, "\n", "", -1)

	return mfaToken
}
