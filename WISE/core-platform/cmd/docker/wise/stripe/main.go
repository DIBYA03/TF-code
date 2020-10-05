package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func main() {

	http.HandleFunc("/stripe/webhook", func(w http.ResponseWriter, r *http.Request) {

		responseData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		responseString := string(responseData)

		// Stripe-Signature is used to verify validity of stripe payload to prevent replay attacks
		stripeSignature := r.Header.Get("Stripe-Signature")

		userAgent := r.UserAgent()
		ipAddress := strings.Split(r.RemoteAddr, ":")[0]

		// Sends message to queue
		err = sendMessage(responseString, stripeSignature, userAgent, ipAddress)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		// Return 200 to stripe
		w.WriteHeader(http.StatusOK)

	})

	// healthcheck
	http.HandleFunc("/healthcheck.html", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	})

	containerPort := os.Getenv("CONTAINER_LISTEN_PORT")
	http.ListenAndServe(fmt.Sprintf(":%s", containerPort), nil)
}

func sendMessage(body string, stripeSignature string, userAgent string, ipAddress string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	svc := sqs.New(sess)

	// URL to our queue
	qURL := os.Getenv("MONEY_REQUEST_SQS_URL")

	_, err = svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"Stripe-Signature": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(stripeSignature),
			},

			"User-Agent": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(userAgent),
			},
			"Source-IP": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(ipAddress),
			},
		},
		MessageBody: aws.String(body),
		QueueUrl:    &qURL,
	})

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
