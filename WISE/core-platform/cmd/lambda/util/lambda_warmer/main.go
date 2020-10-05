package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	lambdaH "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type lambdaWarmingDetail struct {
	Name    string
	Payload string
}

var sess = session.Must(session.NewSession())

// lambdaFunctionList gets all lambdas to search
func lambdaFunctionList() ([]lambdaWarmingDetail, error) {
	log.Println("Grabbing list of lambda functions")

	environment := os.Getenv("API_ENV")
	lambdaName := lambdacontext.FunctionName

	var functionList []lambdaWarmingDetail
	ssmPath := fmt.Sprintf("/%s/%s/", environment, lambdaName)

	svc := ssm.New(sess)
	input := &ssm.GetParametersByPathInput{
		Path:      aws.String(ssmPath),
		Recursive: aws.Bool(true),
	}

	err := svc.GetParametersByPathPages(input,
		func(page *ssm.GetParametersByPathOutput, lastPage bool) bool {
			for _, param := range page.Parameters {
				fnParts := strings.Split(*param.Name, "/")
				lambdaName := fnParts[len(fnParts)-1]
				payload, err := strconv.Unquote(string(*param.Value))
				if err != nil {
					log.Println(lambdaName, err)
				}

				newLambdaDetail := lambdaWarmingDetail{
					Name:    lambdaName,
					Payload: payload,
				}

				functionList = append(functionList, newLambdaDetail)
			}

			return !lastPage
		},
	)

	return functionList, err
}

// invokeLambdas invokes a list of lambdas with ping data
func invokeLambdas(fnNames []lambdaWarmingDetail) {
	log.Println("Attempting to invoke lambdas")

	svc := lambda.New(sess)
	for _, fn := range fnNames {
		_, err := svc.Invoke(
			&lambda.InvokeInput{
				FunctionName: aws.String(fn.Name),
				Payload:      []byte(fn.Payload),
			},
		)
		if err != nil {
			log.Printf("error invoking lambda '%s': %s", fn.Name, err)
		} else {
			log.Printf("lambda invoked: %s", fn.Name)
		}
	}
}

// HandleRequest is the lambda handler
func HandleRequest(ctx context.Context) error {
	log.Print("Function Name:", lambdacontext.FunctionName)

	// Get all the lambdas to invoke
	lambdas, err := lambdaFunctionList()
	if err != nil {
		return err
	}

	// invoke all the lambdas
	invokeLambdas(lambdas)

	return nil
}

func main() {
	lambdaH.Start(HandleRequest)
}
