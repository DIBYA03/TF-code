package main

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/wiseco/core-platform/services/signature"
)

func main() {
	http.HandleFunc("/hellosign/webhook", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(0)
		if err != nil {
			log.Println("error parsing form", err)
			http.Error(w, "Internal server error", 500)
			return
		}

		responseData := r.FormValue("json")

		// validate header
		messageMAC, err := base64.StdEncoding.DecodeString(r.Header.Get("Content-MD5"))
		if err != nil {
			http.Error(w, "Error decoding header", 401)
			return
		}

		apiKey := os.Getenv("HELLOSIGN_API_KEY")
		if len(apiKey) == 0 {
			log.Println("HELLOSIGN_API_KEY is missing")
			panic("HELLOSIGN_API_KEY key missing")
		}

		mac := hmac.New(md5.New, []byte(apiKey))
		mac.Write([]byte(responseData))
		expectedMAC := []byte(hex.EncodeToString(mac.Sum(nil)))

		if !hmac.Equal(messageMAC, expectedMAC) {
			http.Error(w, "Content-MD5 header mismatch", 401)
			return
		}

		// Sends message to queue
		err = sendMessage(responseData)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Return 200 to hellosign
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello API Event Received"))
	})

	// healthcheck
	http.HandleFunc("/healthcheck.html", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	})

	containerPort := os.Getenv("CONTAINER_LISTEN_PORT")
	err := http.ListenAndServeTLS(fmt.Sprintf(":%s", containerPort), "./ssl/cert.pem", "./ssl/key.pem", nil)
	if err != nil {
		log.Fatal("ListenAndServeTLS: ", err)
	}
}

func sendMessage(eventBody string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	svc := sqs.New(sess)

	qURL := os.Getenv("SIGNATURE_SQS_URL")
	if len(qURL) == 0 {
		panic("Signature sqs url missing")
	}

	var eventMessage EventMessage
	err = json.Unmarshal([]byte(eventBody), &eventMessage)
	if err != nil {
		return err
	}

	eventType, ok := EventTypeMap[eventMessage.Event.EventType]
	if !ok {
		return errors.New("Unknown event type")
	}

	sqsMessage := signature.SQSMessage{
		SignatureRequestID: eventMessage.SignatureRequest.SignatureRequestID,
		EventType:          eventType,
	}

	body, err := json.Marshal(sqsMessage)
	if err != nil {
		return err
	}

	_, err = svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageBody:  aws.String(string(body)),
		QueueUrl:     &qURL,
	})

	if err != nil {
		return err
	}

	return nil
}
