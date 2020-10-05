package signature

import (
	"encoding/json"
	l "log"
	"time"

	"github.com/wiseco/core-platform/services"
	bus "github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/services/document"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/go-lib/hellosign"
	"github.com/wiseco/go-lib/log"
)

func HandleMessage(body string) error {
	var sqsMessage SQSMessage

	err := json.Unmarshal([]byte(body), &sqsMessage)
	if err != nil {
		return err
	}

	switch sqsMessage.EventType {
	case EventTypeSignatureRequestAllSigned:
		return handleSignatureRequestAllSigned(sqsMessage.SignatureRequestID, SignatureRequestStatusCompleted)
	default:
		return nil
	}
}

func handleSignatureRequestAllSigned(signatureRequestID string, status SignatureRequestStatus) error {
	// 1. Download signed doc
	lo := log.NewLogger()
	docContent, err := hellosign.NewHellosignService(lo).DownloadDocument(signatureRequestID)
	if err != nil {
		return err
	}

	signature, err := NewSignatureService(services.NewSourceRequest()).GetBySignatureRequestID(signatureRequestID)
	if err != nil {
		return err
	}

	b, err := bus.NewBusinessService(services.NewSourceRequest()).GetByIdInternal(signature.BusinessID)
	if err != nil {
		return err
	}

	// 2. Upload to s3 bucket
	store, err := document.NewAWSS3DocStorageFromContent(string(signature.BusinessID), document.BusinessPrefix, "application/pdf", docContent)
	if err != nil {
		l.Println("Error uploading control person doc to s3", err)
		return err
	}

	key, err := store.Key()
	if err != nil {
		l.Println("Error getting control person document key", err)
		return err
	}

	// 3. Create business doc
	create := document.BusinessDocumentCreate{
		BusinessID:    b.ID,
		CreatedUserID: b.OwnerID,
		ContentType:   "application/pdf",
		StorageKey:    key,
	}

	docType := "other"
	create.DocType = &docType

	docNumber := document.DocumentNumber("12341234")
	create.Number = &docNumber

	issueDate := shared.Date(time.Now().AddDate(-10, 0, 0))
	create.IssuedDate = &issueDate

	expDate := shared.Date(time.Now().AddDate(10, 0, 0))
	create.ExpirationDate = &expDate

	state := "any"
	create.IssuingState = &state

	country := "US"
	create.IssuingCountry = &country

	sr := services.NewSourceRequest()
	sr.UserID = b.OwnerID

	doc, err := document.NewBusinessDocumentService(sr).CreateInternal(create)
	if err != nil {
		l.Println("error creating business doc")
		return err
	}

	// 4. Update signature request status
	signature, err = NewSignatureService(services.NewSourceRequest()).UpdateSignatureStatus(signature.ID, doc.ID, SignatureRequestStatusCompleted)
	if err != nil {
		l.Println(err)
		return err
	}

	l.Println("Webhook processed successfully for signature request id ", signature.ID, signature.SignatureStatus)

	return nil
}
