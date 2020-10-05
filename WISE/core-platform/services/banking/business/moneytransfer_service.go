/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

// Package for all internal banking related items such as accounts and cards
package business

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
	"github.com/wiseco/core-platform/auth"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/go-lib/id"
	l "github.com/wiseco/go-lib/log"
	sg "github.com/wiseco/go-lib/sendgrid"
	"github.com/wiseco/go-lib/slack"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking"
	bsrv "github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/services/data"
)

type moneyTransferDatastore struct {
	sourceReq services.SourceRequest
	*sqlx.DB
}

type MoneyTransferService interface {
	// Get all money transfers by business ID
	GetByBusinessID(offset, limit int, businessID shared.BusinessID, userID shared.UserID) ([]MoneyTransfer, error)

	// Get money transfer by bank id
	GetByBankID(bankName banking.BankName, bankTransferID string, businessID shared.BusinessID) (*MoneyTransfer, error)
	GetByBankIDInternal(bankName banking.BankName, bankTransferID string) (*MoneyTransfer, error)

	GetAccountByTransferID(bankName, bankTransferID string, businessID shared.BusinessID) (*LinkedBankAccount, error)

	GetByIDInternal(string, shared.BusinessID) (*MoneyTransfer, error)

	GetByIDOnlyInternal(string) (*MoneyTransfer, error)

	//Transfer funds between business' internal linked accounts
	Transfer(*TransferInitiate) (*MoneyTransfer, error)

	// UpdateStatus
	UpdateStatus(businessID shared.BusinessID, moneyTransferID, status string) error

	// UpdatePostedTransaction
	UpdateDebitPostedTransaction(businessID shared.BusinessID, moneyTransferID, postedTransactionID string) error
	UpdateCreditPostedTransaction(moneyTransferID, postedTransactionID string) error

	OnMoneyTransfer(*MoneyTransfer, *string, *string, banking.PartnerName)
}

func NewMoneyTransferService(r services.SourceRequest) MoneyTransferService {
	return &moneyTransferDatastore{r, data.DBWrite}
}

func NewMoneyTransferServiceWithout() MoneyTransferService {
	return &moneyTransferDatastore{services.SourceRequest{}, data.DBWrite}
}

func (db *moneyTransferDatastore) GetByBusinessID(offset int, limit int, businessID shared.BusinessID, userID shared.UserID) ([]MoneyTransfer, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	rows := []MoneyTransfer{}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bts, err := NewBankingTransferService()
		if err != nil {
			return rows, err
		}

		return bts.GetByBusinessID(businessID, offset, limit)
	}

	err = db.Select(&rows, "SELECT * FROM business_money_transfer WHERE business_id = $1 AND created_user_id = $2", businessID, userID)
	if err != nil && err != sql.ErrNoRows {

		return nil, err
	}

	return rows, err
}

func (db *moneyTransferDatastore) GetByBankID(bankName banking.BankName, bankTransferID string, businessID shared.BusinessID) (*MoneyTransfer, error) {
	// Check access
	err := auth.NewAuthService(db.sourceReq).CheckBusinessAccess(businessID)
	if err != nil {
		return nil, err
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bts, err := NewBankingTransferService()
		if err != nil {
			return nil, err
		}

		return bts.GetByBankID(businessID, bankTransferID)
	}

	var mt MoneyTransfer
	err = db.Get(
		&mt,
		"SELECT * FROM business_money_transfer WHERE bank_transfer_id = $1 AND bank_name = $2 AND business_id = $3",
		bankTransferID,
		bankName,
		businessID,
	)
	if err != nil {
		return nil, err
	}

	return &mt, nil
}

func (db *moneyTransferDatastore) GetByBankIDInternal(bankName banking.BankName, bankTransferID string) (*MoneyTransfer, error) {
	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bts, err := NewBankingTransferService()
		if err != nil {
			return nil, err
		}

		bID := shared.BusinessID("")

		return bts.GetByBankID(bID, bankTransferID)
	}

	var mt MoneyTransfer
	err := db.Get(
		&mt,
		"SELECT * FROM business_money_transfer WHERE bank_transfer_id = $1 AND bank_name = $2",
		bankTransferID,
		bankName,
	)
	if err != nil {
		return nil, err
	}
	return &mt, nil
}

func (db *moneyTransferDatastore) GetByIDInternal(id string, businessID shared.BusinessID) (*MoneyTransfer, error) {
	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bts, err := NewBankingTransferService()
		if err != nil {
			return nil, err
		}

		return bts.GetByIDInternal(businessID, id)
	}

	var mt MoneyTransfer
	err := db.Get(
		&mt,
		"SELECT * FROM business_money_transfer WHERE id = $1 AND business_id = $2",
		id,
		businessID,
	)
	if err != nil {
		return nil, err
	}

	return &mt, nil
}

func (db *moneyTransferDatastore) GetByIDOnlyInternal(id string) (*MoneyTransfer, error) {
	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bts, err := NewBankingTransferService()
		if err != nil {
			return nil, err
		}

		return bts.GetByIDOnlyInternal(id)
	}

	var mt MoneyTransfer
	err := db.Get(
		&mt,
		"SELECT * FROM business_money_transfer WHERE id = $1",
		id,
	)
	if err != nil {
		return nil, err
	}

	return &mt, nil
}

func (db *moneyTransferDatastore) GetAccountByTransferID(bankName, bankTransferID string, businessID shared.BusinessID) (*LinkedBankAccount, error) {
	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		blas, err := NewBankingLinkedAccountService()
		if err != nil {
			return nil, err
		}

		return blas.GetAccountByTransferID(businessID, bankTransferID)
	}

	var la LinkedBankAccount
	err := db.Get(
		&la, `
		SELECT business_linked_bank_account.*
		FROM business_linked_bank_account 
		JOIN business_money_transfer ON business_linked_bank_account.id = business_money_transfer.source_account_id 
		WHERE
			business_money_transfer.id = $1 AND
			business_money_transfer.bank_name = $2 AND
			business_money_transfer.business_id = $3`,
		bankTransferID,
		bankName,
		businessID,
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &la, nil
}

func (db *moneyTransferDatastore) Transfer(transfer *TransferInitiate) (*MoneyTransfer, error) {
	if transfer.SourceAccountId == "" {
		return nil, errors.New("source account id is required")
	}

	if transfer.DestAccountId == "" {
		return nil, errors.New("destination account id is required")
	}

	if transfer.Amount == 0 {
		return nil, errors.New("amount is required")
	}

	if transfer.SourceType == "" {
		return nil, errors.New("source type is required")
	}

	if transfer.SourceType == banking.TransferTypeCard && transfer.CVVCode == nil {
		return nil, errors.New("CVV is required for debit pull")
	}

	if transfer.DestType == "" {
		return nil, errors.New("destination type is required")
	}

	if transfer.DestType != banking.TransferTypeAccount {
		return nil, errors.New("only account type is supported for destination")
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bts, err := NewBankingTransferService()
		if err != nil {
			return nil, err
		}

		var dbaID *string
		var sbaID *string

		las, err := NewBankingLinkedAccountService()
		if err != nil {
			return nil, err
		}

		var sla *LinkedBankAccount
		var dla *LinkedBankAccount

		if strings.HasPrefix(transfer.DestAccountId, id.IDPrefixLinkedBankAccount.String()) {
			las, err := NewBankingLinkedAccountService()
			if err != nil {
				return nil, err
			}

			dla, err = las.GetById(transfer.DestAccountId)
			if err != nil {
				return nil, err
			}

			dbaID = dla.BusinessBankAccountId
		} else if strings.HasPrefix(transfer.DestAccountId, id.IDPrefixBankAccount.String()) {
			da, err := NewBankAccountService(db.sourceReq).GetByIDInternal(transfer.DestAccountId)
			if err != nil {
				return nil, err
			}

			la, err := las.GetByAccountIDInternal(da.Id)
			if err != nil {
				return nil, err
			}

			dbaID = la.BusinessBankAccountId
		} else if transfer.DestType == banking.TransferTypeAccount {
			//This is an edge case from pushing the banking service, clients that have not logged out will have ids with no prefixes
			las, err := NewBankingLinkedAccountService()
			if err != nil {
				return nil, err
			}

			dla, err = las.GetById(transfer.DestAccountId)
			if err != nil {
				return nil, err
			}

			dbaID = dla.BusinessBankAccountId
		}

		las, err = NewBankingLinkedAccountService()
		if err != nil {
			return nil, err
		}

		if strings.HasPrefix(transfer.SourceAccountId, id.IDPrefixLinkedBankAccount.String()) {
			sla, err = las.GetById(transfer.SourceAccountId)
			if err != nil {
				return nil, err
			}

			sbaID = sla.BusinessBankAccountId
		} else if strings.HasPrefix(transfer.SourceAccountId, id.IDPrefixBankAccount.String()) {
			sa, err := NewBankAccountService(db.sourceReq).GetByIDInternal(transfer.SourceAccountId)
			if err != nil {
				return nil, err
			}

			sla, err = las.GetByAccountIDInternal(sa.Id)
			if err != nil {
				return nil, err
			}

			sbaID = sla.BusinessBankAccountId
		} else if transfer.SourceType == banking.TransferTypeAccount {
			//This is an edge case from pushing the banking service, clients that have not logged out will have ids with no prefixes
			sla, err = las.GetById(transfer.SourceAccountId)
			if err != nil {
				return nil, err
			}

			sbaID = sla.BusinessBankAccountId
		}

		sut := UsageTypeNone

		if transfer.SourceType == banking.TransferTypeAccount && transfer.DestType == banking.TransferTypeAccount {
			sut = *sla.UsageType
		}

		mt, err := bts.Transfer(transfer, sut, db.sourceReq)
		if err != nil {
			return nil, err
		}

		if sbaID != nil {
			id, err := id.ParseBankAccountID(*sbaID)
			if err != nil {
				return nil, err
			}

			ids := id.UUIDString()

			sbaID = &ids
		}

		if dbaID != nil {
			id, err := id.ParseBankAccountID(*dbaID)
			if err != nil {
				return nil, err
			}

			ids := id.UUIDString()

			dbaID = &ids
		}

		// create pending transaction
		db.OnMoneyTransfer(mt, sbaID, dbaID, banking.PartnerNamePlaid)

		// Notify on slack and send email
		err = db.notifyACHTransfer(sla, dla, transfer)
		if err != nil {
			log.Println("Error notifying ACH transfer", err)
		}

		return mt, nil
	}

	sourceId, source, err := db.getDestinationIDAndService(transfer)
	if err != nil {
		return nil, err
	}

	var sa *LinkedBankAccount
	var sourceBankAccountID *string
	switch transfer.SourceType {
	case banking.TransferTypeAccount:
		sa = (source.(*LinkedBankAccount))
		sourceBankAccountID = sa.BusinessBankAccountId
	default:
		sourceBankAccountID = nil
	}

	// Search destination account in linked accounts
	da, err := NewLinkedAccountService(db.sourceReq).GetById(transfer.DestAccountId, transfer.BusinessID)
	if err != nil {
		return nil, err
	}

	err = checkACHMaxLimit(sa, da, transfer)
	if err != nil {
		return nil, err
	}

	request := partnerbank.MoneyTransferRequest{
		Amount:          transfer.Amount,
		Currency:        partnerbank.Currency(transfer.Currency),
		SourceAccountID: partnerbank.MoneyTransferAccountBankID(sourceId),
		DestAccountID:   partnerbank.MoneyTransferAccountBankID(da.RegisteredAccountId),
	}

	// Register card with partner bank
	bank, err := partnerbank.GetBusinessBank(partnerbank.ProviderNameBBVA)
	if err != nil {
		return nil, err
	}

	srv, err := bank.MoneyTransferService(db.sourceReq.PartnerBankRequest(), partnerbank.BusinessID(transfer.BusinessID))
	if err != nil {
		return nil, err
	}

	resp, err := srv.Submit(&request)
	if err != nil {
		return nil, err
	}

	t := transformTransferResponse(resp)
	t.BusinessID = transfer.BusinessID
	t.CreatedUserID = transfer.CreatedUserID
	t.SourceAccountId = transfer.SourceAccountId
	t.SourceType = transfer.SourceType
	t.DestAccountId = transfer.DestAccountId
	t.DestType = transfer.DestType
	t.MonthlyInterestID = transfer.MonthlyInterestID
	t.Notes = transfer.Notes
	t.MoneyRequestID = transfer.MoneyRequestID
	t.ContactId = transfer.ContactId

	// Default/mandatory fields
	columns := []string{
		"created_user_id", "business_id", "bank_name", "bank_transfer_id", "source_account_id", "source_type",
		"dest_account_id", "dest_type", "amount", "currency", "notes", "status", "account_monthly_interest_id",
		"money_request_id", "contact_id",
	}
	// Default/mandatory values
	values := []string{
		":created_user_id", ":business_id", ":bank_name", ":bank_transfer_id", ":source_account_id", ":source_type",
		":dest_account_id", ":dest_type", ":amount", ":currency", ":notes", ":status", ":account_monthly_interest_id",
		":money_request_id", ":contact_id",
	}

	sql := fmt.Sprintf("INSERT INTO business_money_transfer(%s) VALUES(%s) RETURNING *", strings.Join(columns, ", "), strings.Join(values, ", "))

	stmt, err := db.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	l := &MoneyTransfer{}

	err = stmt.Get(l, &t)
	if err != nil {
		return nil, err
	}

	// create pending transaction
	db.OnMoneyTransfer(l, sourceBankAccountID, da.BusinessBankAccountId, banking.PartnerNamePlaid)

	// Notify on slack and send email
	err = db.notifyACHTransfer(sa, da, transfer)
	if err != nil {
		log.Println("Error notifying ACH transfer", err)
	}

	return l, nil
}

func checkACHMaxLimit(sa *LinkedBankAccount, da *LinkedBankAccount, transfer *TransferInitiate) error {
	if transfer.ContactId != nil {
		return nil
	}

	// Handle linked cards
	if sa == nil {
		return nil
	}

	if sa.Source() == LinkedAccountSourcePlaid {
		maxAmount, err := strconv.ParseFloat(os.Getenv("ACH_MAX_ALLOWED"), 64)
		if err != nil {
			log.Println(err)
			return err
		}

		if transfer.Amount > maxAmount {
			e := fmt.Sprintf("Transfers from external accounts are limited to $%s per transaction.", strconv.FormatFloat(maxAmount, 'f', 2, 64))
			return errors.New(e)
		}
	}

	if da.Source() == LinkedAccountSourcePlaid {
		maxAmount, err := strconv.ParseFloat(os.Getenv("ACH_MAX_ALLOWED"), 64)
		if err != nil {
			log.Println(err)
			return err
		}

		if transfer.Amount > maxAmount {
			e := fmt.Sprintf("Transfers to external accounts are limited to $%s per transaction.", strconv.FormatFloat(maxAmount, 'f', 2, 64))
			return errors.New(e)
		}
	}

	return nil
}

// Transform partner bank layer response
func transformTransferResponse(response *partnerbank.MoneyTransferResponse) MoneyTransfer {

	t := MoneyTransfer{}

	t.BankName = "bbva"
	t.BankTransferId = string(response.TransferID)
	t.Amount = response.Amount
	t.Currency = banking.Currency(response.Currency)
	t.Status = string(response.Status)

	println("transfer id is ", response.TransferID, response.Status)

	return t
}

//Update money transfer status
func (db moneyTransferDatastore) UpdateStatus(businessID shared.BusinessID, moneyTransferID, status string) error {
	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bts, err := NewBankingTransferService()
		if err != nil {
			return err
		}

		return bts.UpdateStatus(businessID, moneyTransferID, status)
	}

	_, err := db.Exec(`
		UPDATE business_money_transfer
		SET status = $1 WHERE bank_transfer_id = $2 AND business_id = $3`,
		status,
		moneyTransferID,
		businessID,
	)

	return err
}

func (db moneyTransferDatastore) UpdateDebitPostedTransaction(businessID shared.BusinessID, moneyTransferID, postedTransactionID string) error {
	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bts, err := NewBankingTransferService()
		if err != nil {
			return err
		}

		return bts.UpdateDebitPostedTransaction(businessID, moneyTransferID, postedTransactionID)
	}

	_, err := db.Exec(`
        UPDATE business_money_transfer
        SET posted_debit_transaction_id = $1 WHERE id = $2 AND business_id = $3`,
		postedTransactionID,
		moneyTransferID,
		businessID,
	)

	return err
}

func (db moneyTransferDatastore) UpdateCreditPostedTransaction(moneyTransferID, postedTransactionID string) error {
	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bts, err := NewBankingTransferService()
		if err != nil {
			return err
		}

		return bts.UpdateCreditPostedTransaction(moneyTransferID, postedTransactionID)
	}

	_, err := db.Exec(`
        UPDATE business_money_transfer
        SET posted_credit_transaction_id = $1 WHERE id = $2`,
		postedTransactionID,
		moneyTransferID,
	)

	return err
}

func (db *moneyTransferDatastore) OnMoneyTransfer(m *MoneyTransfer, sourceBankAccountID *string, destBankAccountID *string, partnerName banking.PartnerName) {
	if m.Status == banking.MoneyTransferStatusInProcess {
		// Check both source and destination
		var accountID string
		var codeType string
		if destBankAccountID != nil {
			codeType = CreditInProcess
			accountID = *destBankAccountID
		} else if sourceBankAccountID != nil {
			codeType = DebitInProcess
			accountID = *sourceBankAccountID
			m.Amount = -m.Amount
		}

		txnDate, err := time.Parse("2006-01-02T15:04:05", m.Created.Format("2006-01-02T15:04:05"))
		if err != nil {
			log.Println("Error parsing money transfer date ", err)
		}

		mt := PendingTransferNotification{
			BankName:        banking.BankNameBBVA,
			TransactionType: PendingMoneyTransferType,
			BankAccountID:   accountID,
			Status:          m.Status,
			Amount:          m.Amount,
			ParterName:      partnerName,
			Currency:        banking.CurrencyUSD,
			MoneyTransferID: &m.Id,
			TransactionDate: txnDate,
			ContactID:       m.ContactId,
			CodeType:        codeType,
			MoneyRequestID:  m.MoneyRequestID,
			Notes:           m.Notes,
		}

		body, err := json.Marshal(mt)
		if err != nil {
			log.Printf("Error marshalling data into %v error :%v ", mt, err)
			return
		}

		n := Notification{
			ID:         uuid.New().String(),
			EntityID:   string(m.BusinessID),
			EntityType: EntityTypeBusiness,
			BankName:   banking.BankNameBBVA,
			Type:       PendingMoneyTransfer,
			Version:    NotificationVersion,
			Created:    time.Now(),
			Data:       body,
		}

		// push notification
		err = pushToNotificationQueue(n)
		if err != nil {
			log.Println("Error pushing pending money transfer to notification queue ", err)
		} else {
			log.Println("Pushed pending money transfer to notification queue")
		}
	}
}

func pushToNotificationQueue(n Notification) error {

	out, err := json.Marshal(n)
	if err != nil {
		log.Println(err)
		return err
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("SQS_BANKING_REGION"))},
	)

	svc := sqs.New(sess)

	// URL to our queue
	qURL := os.Getenv("SQS_BANKING_URL")

	_, err = svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(5),
		MessageBody:  aws.String(string(out)),
		QueueUrl:     &qURL,
	})

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (db *moneyTransferDatastore) getDestinationIDAndService(transfer *TransferInitiate) (string, interface{}, error) {

	if transfer.SourceType == banking.TransferTypeAccount {
		// Search destination account in linked accounts
		sa, err := NewLinkedAccountService(db.sourceReq).GetById(transfer.SourceAccountId, transfer.BusinessID)
		if err != nil {
			log.Println("unable to get dest account", transfer.SourceAccountId, transfer.BusinessID, err)
			return "", nil, err
		}

		// Check for source account usage type - must be primary or clearing
		if sa.UsageType == nil {
			return "", nil, errors.New("source account type must be primary or clearing")
		}

		// Check if account has sufficient funds to transfer
		if sa.BusinessBankAccountId != nil {
			balance, err := NewBankAccountService(db.sourceReq).GetBalanceByID(*sa.BusinessBankAccountId, transfer.BusinessID)
			if err != nil {
				log.Println("error fetching account balance", err)
				return "", nil, err
			}

			if balance.PostedBalance < transfer.Amount {
				err := errors.New("Insufficient funds to initiate move money")
				log.Println(err)
				return "", nil, err
			}
		}

		return sa.RegisteredAccountId, sa, nil

	} else if transfer.SourceType == banking.TransferTypeCard {
		// Search destination account in linked cards
		sc, err := NewLinkedCardService(db.sourceReq).GetByID(transfer.SourceAccountId, transfer.BusinessID)
		if err != nil {
			return "", nil, err
		}

		return sc.RegisteredCardId, sc, nil

	}

	return "", nil, fmt.Errorf("Unable to get destination for type: %s", transfer.SourceAccountId)
}

func (db *moneyTransferDatastore) sendSlackMessage(message string) error {
	m := slack.Message{
		Text: message,
	}

	slackUrl := os.Getenv("ACH_NOTIFICATION_SLACK_URL")
	if slackUrl == "" {
		return errors.New("ACH_NOTIFICATION_SLACK_URL variable is missing")
	}

	return slack.NewSlackService(l.NewLogger()).PostToChannel(slackUrl, m)
}

func (db *moneyTransferDatastore) sendEmail(subject, body string) error {
	req := sg.EmailRequest{
		SenderEmail:   os.Getenv("WISE_SUPPORT_EMAIL"),
		SenderName:    os.Getenv("WISE_SUPPORT_NAME"),
		ReceiverEmail: os.Getenv("WISE_FRAUDOPS_INVOICE_EMAIL"),
		ReceiverName:  os.Getenv("WISE_FRAUDOPS_INVOICE_NAME"),
		Subject:       subject,
		Body:          body,
	}

	_, err := sg.NewSendGrid().SendEmail(req)

	return err
}

func (db *moneyTransferDatastore) notifyACHTransfer(sa *LinkedBankAccount, da *LinkedBankAccount, transfer *TransferInitiate) error {
	// Notify only for transfers between business' own accounts
	if transfer.ContactId != nil {
		return nil
	}

	if da == nil {
		return nil
	}

	// Handle linked cards
	if sa == nil {
		return nil
	}

	// Get business
	b, err := bsrv.NewBusinessService(db.sourceReq).GetByIdInternal(transfer.BusinessID)
	if err != nil {
		return err
	}

	var message string
	var subject string

	bName := b.Name()
	amt := shared.FormatFloatAmount(transfer.Amount)

	if da.Source() == LinkedAccountSourcePlaid {
		extAccount, err := NewExternalAccountService(db.sourceReq).GetByAccountNumberInternal(string(da.AccountNumber),
			da.RoutingNumber, transfer.BusinessID)
		if err != nil && err != sql.ErrNoRows {
			return err
		}

		if da.BankName != nil {
			message = fmt.Sprintf("%s sent $%s to %s\n", bName, amt, string(*da.BankName))
		} else {
			message = fmt.Sprintf("%s sent $%s\n", bName, amt)
		}

		if extAccount != nil {
			result, err := NewExternalAccountService(db.sourceReq).GetVerificationByAccountID(extAccount.ID, transfer.BusinessID)
			if err != nil && err != sql.ErrNoRows {
				return err
			}

			if result != nil {
				message = message + fmt.Sprintf("Verification result: %s", result.VerificationResult)
			}
		}

	} else if sa.Source() == LinkedAccountSourcePlaid {
		extAccount, err := NewExternalAccountService(db.sourceReq).GetByAccountNumberInternal(string(sa.AccountNumber),
			sa.RoutingNumber, transfer.BusinessID)
		if err != nil && err != sql.ErrNoRows {
			return err
		}

		if sa.BankName != nil {
			message = fmt.Sprintf("%s received $%s from %s\n", bName, amt, string(*sa.BankName))
		} else {
			message = fmt.Sprintf("%s received $%s\n", bName, amt)
		}

		subject = message

		if extAccount != nil {
			result, err := NewExternalAccountService(db.sourceReq).GetVerificationByAccountID(extAccount.ID, transfer.BusinessID)
			if err != nil && err != sql.ErrNoRows {
				return err
			}

			if result != nil {
				message = message + fmt.Sprintf("Verification result: %s", result.VerificationResult)
			}
		}
	}

	if message != "" {
		// send email
		err = db.sendEmail(subject, message)
		if err != nil {
			return err
		}

		// send slack message
		err = db.sendSlackMessage(message)
		if err != nil {
			return err
		}
	}

	return nil
}
