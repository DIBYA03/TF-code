package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

var oldAccountAWSProfile = "master-us-west-2-saml-roles-deployment"
var oldAccountAWSRegion = "us-west-2"
var newAccountAWSProfile = "dev-us-west-2-saml-roles-deployment"
var newAccountAWSRegion = "us-west-2"
var overwriteParams = true

func main() {
	// Specify profile for config and region for requests
	oldAccountSess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String(oldAccountAWSRegion)},
		Profile: oldAccountAWSProfile,
	}))

	newAccountSess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String(newAccountAWSRegion)},
		Profile: newAccountAWSProfile,
	}))

	oldEnvironment, newEnvironment, kmsKey := getUserInput()

	parameterResults := oldSSMParameters(oldAccountSess, oldEnvironment, "")
	createNewParameters(newAccountSess, parameterResults.Parameters, oldEnvironment, newEnvironment, kmsKey)
	for {
		if parameterResults.NextToken != nil {
			parameterResults = oldSSMParameters(oldAccountSess, oldEnvironment, *parameterResults.NextToken)
			createNewParameters(newAccountSess, parameterResults.Parameters, oldEnvironment, newEnvironment, kmsKey)
		} else {
			break
		}
	}
}

func createNewParameters(sess *session.Session, params []*ssm.Parameter, oldEnvironment string, newEnvironment string, kmsKey string) {
	svc := ssm.New(sess)

	for _, param := range params {
		newParamName := strings.Replace(*param.Name, oldEnvironment, newEnvironment, 1)

		fmt.Printf("%s -> %s\n", *param.Name, newParamName)

		input := &ssm.PutParameterInput{
			KeyId:     aws.String(kmsKey),
			Name:      aws.String(newParamName),
			Overwrite: aws.Bool(overwriteParams),
			Type:      param.Type,
			Value:     param.Value,
		}

		paramResults, err := svc.PutParameter(input)
		if err != nil {
			log.Println(err)
			continue
		}

		fmt.Println(paramResults)

	}

}
func oldSSMParameters(sess *session.Session, environment string, nextToken string) *ssm.GetParametersByPathOutput {
	svc := ssm.New(sess)

	var input *ssm.GetParametersByPathInput
	if nextToken == "" {
		input = &ssm.GetParametersByPathInput{
			Path:           aws.String(fmt.Sprintf("/%s", environment)),
			Recursive:      aws.Bool(true),
			WithDecryption: aws.Bool(true),
		}
	} else {
		input = &ssm.GetParametersByPathInput{
			NextToken:      aws.String(nextToken),
			Path:           aws.String(fmt.Sprintf("/%s", environment)),
			Recursive:      aws.Bool(true),
			WithDecryption: aws.Bool(true),
		}
	}

	paramResults, err := svc.GetParametersByPath(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			panic(aerr)
		} else {
			panic(err.Error())
		}
	}

	return paramResults
}

func getUserInput() (string, string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Old Environment: ")
	oldEnvironment, _ := reader.ReadString('\n')
	oldEnvironment = strings.Replace(oldEnvironment, "\n", "", -1)

	fmt.Print("New Environment: ")
	newEnvironment, _ := reader.ReadString('\n')
	newEnvironment = strings.Replace(newEnvironment, "\n", "", -1)

	fmt.Print("New KMS Key: ")
	kmsKey, _ := reader.ReadString('\n')
	kmsKey = strings.Replace(kmsKey, "\n", "", -1)

	// Cheat way to make it look purty
	fmt.Print("\n")

	return oldEnvironment, newEnvironment, kmsKey
}
