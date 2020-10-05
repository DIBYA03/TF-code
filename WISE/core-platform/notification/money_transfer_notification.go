package notification

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/wiseco/core-platform/notification/activity"
	"github.com/wiseco/core-platform/services"
	busbanking "github.com/wiseco/core-platform/services/banking/business"
	core "github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/core-platform/transaction"
	"github.com/wiseco/go-lib/num"
)

func processPendingMoneyTransferNotification(n Notification) error {
	var mt PendingTransferNotification

	err := json.Unmarshal(n.Data, &mt)
	if err != nil {
		return fmt.Errorf("pending money transfer notification error: %v", err)
	}

	if mt.MoneyTransferID == nil {
		return fmt.Errorf("money transfer ID is null")
	}

	bID, err := shared.ParseBusinessID(n.EntityID)
	if err != nil {
		return fmt.Errorf("error parsing business id: %s", n.EntityID)
	}

	// Check status of money transfer
	m, err := busbanking.NewMoneyTransferService(services.NewSourceRequest()).GetByIDInternal(*mt.MoneyTransferID, bID)
	if err != nil {
		return err
	}

	if m.Status == string(ActionPosted) {
		return nil
	}

	// Create pending transaction
	amount, err := num.NewDecimalFin(mt.Amount)
	if err != nil {
		return err
	}

	t := transaction.BusinessPendingTransactionCreate{
		// Generate Transaction ID for every transaction
		ID:                shared.PendingTransactionID(uuid.New().String()),
		BusinessID:        bID,
		BankName:          n.BankName,
		TransactionType:   transaction.TransactionType(mt.TransactionType),
		AccountID:         &mt.BankAccountID,
		Amount:            num.Decimal{V: amount},
		Currency:          transaction.Currency(mt.Currency),
		MoneyTransferID:   mt.MoneyTransferID,
		ContactID:         mt.ContactID,
		TransactionDate:   mt.TransactionDate,
		PartnerName:       mt.PartnerName,
		TransactionStatus: &mt.Status,
		CodeType:          mt.CodeType,
		MoneyRequestID:    mt.MoneyRequestID,
		SourceNotes:       mt.Notes,
	}

	switch t.CodeType {
	case transaction.TransactionCodeTypeCreditInProcess:
		return processCreditInProcessTransfer(&n, mt, t)
	case transaction.TransactionCodeTypeDebitInProcess:
		return processDebitInProcessTransfer(&n, mt, t)
	default:
		log.Println(fmt.Sprintf("invalid pending transaction code type %s", t.CodeType))
		return nil
	}

}

func processDebitInProcessTransfer(n *Notification, mt PendingTransferNotification, pendingTransaction transaction.BusinessPendingTransactionCreate) error {
	la, err := busbanking.NewMoneyTransferService(services.NewSourceRequest()).GetAccountByTransferID(pendingTransaction.BankName, *pendingTransaction.MoneyTransferID, pendingTransaction.BusinessID)
	if err != nil {
		log.Println("Error fetching linked account ", err)
		return err
	}

	businessName, err := getBusinessName(pendingTransaction.BusinessID)
	if err != nil {
		return err
	}

	txnTitle := fmt.Sprintf(DebitACHTransactionTitle, la.AccountHolderName)
	txnDescription := fmt.Sprintf(DebitACHTransactionDescription, *businessName, la.AccountHolderName)
	moneyTransferDesc := "DEST: " + string(la.AccountNumber[len(la.AccountNumber)-4:]) + " " + la.AccountHolderName

	pendingTransaction.TransactionSubtype = transaction.ActivityToTransactionSubtype[activity.TypeACHTransferDebit]
	pendingTransaction.TransactionTitle = txnTitle
	pendingTransaction.TransactionDesc = txnDescription
	pendingTransaction.BankTransactionDesc = &moneyTransferDesc
	pendingTransaction.Status = transaction.TransactionStatusBankProcessing
	pendingTransaction.Counterparty = la.AccountHolderName

	_, err = createPendingTransaction(n, pendingTransaction, nil, nil)
	if err != nil {
		return err
	}

	fAmount, _ := pendingTransaction.Amount.Abs().Float64()

	// Generate activity stream
	//name := c.Contact.Name()
	a := activity.AccountTransaction{
		BusinessName:    businessName,
		EntityID:        string(pendingTransaction.BusinessID),
		UserID:          la.UserID,
		ContactName:     &la.AccountHolderName,
		Amount:          activity.AccountTransactionAmount(fAmount),
		TransactionDate: time.Now(),
	}

	_, err = activity.NewTransferCreator().DebitInProcess(a)
	if err != nil {
		log.Println("Error creating debit in process activity stream", err)
		return nil
	}

	return nil

}

func processCreditInProcessTransfer(n *Notification, mt PendingTransferNotification, pendingTransaction transaction.BusinessPendingTransactionCreate) error {
	la, err := busbanking.NewMoneyTransferService(services.NewSourceRequest()).GetAccountByTransferID(pendingTransaction.BankName, *pendingTransaction.MoneyTransferID, pendingTransaction.BusinessID)
	if err != nil {
		log.Println("Error fetching linked account ", err)
		return err
	}

	businessName, err := getBusinessName(pendingTransaction.BusinessID)
	if err != nil {
		return err
	}

	amt := shared.FormatFloatAmount(math.Abs(mt.Amount))

	moneyTransferDesc := "ORIG: " + string(la.AccountNumber[len(la.AccountNumber)-4:]) + " " + la.AccountHolderName + " DEST: "

	var txnTitle string
	var txnDescription string
	if mt.MoneyRequestID != nil {
		txnTitle = fmt.Sprintf(CreditViaBankTransactionTitle, la.AccountHolderName, *businessName)
		txnDescription = fmt.Sprintf(CreditViaBankTransactionDescription, la.AccountHolderName, *businessName, "Bank Transfer")
		pendingTransaction.TransactionSubtype = transaction.ActivityToTransactionSubtype[activity.TypeBankOnlineCredit]
		pendingTransaction.Counterparty = la.AccountHolderName
	} else {
		txnTitle = fmt.Sprintf(CreditACHTransactionTitle, la.AccountHolderName)
		txnDescription = fmt.Sprintf(CreditACHTransactionDescription, amt, la.AccountHolderName)
		pendingTransaction.TransactionSubtype = transaction.ActivityToTransactionSubtype[activity.TypeACHTransferCredit]
		pendingTransaction.Counterparty = la.AccountHolderName
	}

	pendingTransaction.TransactionTitle = txnTitle
	pendingTransaction.TransactionDesc = txnDescription
	pendingTransaction.BankTransactionDesc = &moneyTransferDesc
	pendingTransaction.Status = transaction.TransactionStatusBankProcessing

	_, err = createPendingTransaction(n, pendingTransaction, nil, nil)
	if err != nil {
		return err
	}

	fAmount, _ := pendingTransaction.Amount.Abs().Float64()

	// Generate activity stream
	a := activity.AccountTransaction{
		BusinessName:    businessName,
		EntityID:        string(pendingTransaction.BusinessID),
		UserID:          la.UserID,
		Amount:          activity.AccountTransactionAmount(fAmount),
		Origin:          la.AccountHolderName,
		TransactionDate: time.Now(),
	}

	_, err = activity.NewTransferCreator().CreditInProcess(a)
	if err != nil {
		log.Println("Error creating credit in process activity stream", err)
		return nil
	}

	return nil
}

func processMoneyTransferNotification(n Notification) error {

	switch n.Action {
	case ActionUpdate:
		return processUpdateMoneyTransferNotification(n)
	default:
		log.Println(fmt.Errorf("invalid money transfer notification action: %s", n.Action))
		return nil
	}
}

func processUpdateMoneyTransferNotification(n Notification) error {
	var m MoneyTransferStatusNotification
	err := json.Unmarshal(n.Data, &m)

	if err != nil {
		return fmt.Errorf("money transfer notification error: %v", err)
	}

	// Get business ID
	bID, err := NewNotificationService().getBusinessIDByEntity(n.EntityID, n.EntityType, nil)
	if err != nil {
		log.Println("Business ID not found")
		return err
	}

	var moneyTransferID *string
	var mt busbanking.MoneyTransfer

	// Update money transfer status
	err = busbanking.NewMoneyTransferService(services.NewSourceRequest()).UpdateStatus(*bID, m.MoneyTransferID, m.Status)
	if err != nil {
		log.Println("Error updating money transfer status ", err)
	}

	if os.Getenv("USE_BANKING_SERVICE") == "true" {
		bts, err := busbanking.NewBankingTransferService()
		if err != nil {
			return err
		}

		transfer, err := bts.GetByBankID(*bID, m.MoneyTransferID)
		if err != nil {
			return err
		}

		if transfer != nil {
			mt = *transfer
		}
	} else {
		err = core.DBRead.Get(
			&mt,
			`SELECT * FROM business_money_transfer WHERE bank_transfer_id = $1`,
			m.MoneyTransferID,
		)
	}

	if err == nil {
		moneyTransferID = &mt.Id

		// Update pending transaction
		u := transaction.BusinessPendingTransactionUpdate{
			BusinessID:      *bID,
			Status:          m.Status,
			MoneyTransferID: *moneyTransferID,
		}
		err = transaction.NewPendingTransactionService().UpdateMoneyTransferStatus(u)
		if err != nil {
			log.Println("Error updating pending money transfer status ", err)
		}
	}

	return err
}

func getBusinessName(bID shared.BusinessID) (*string, error) {
	type BusinessName struct {
		LegalName string               `db:"business.legal_name"`
		DBA       services.StringArray `db:"business.dba"`
	}

	bt := BusinessName{}
	err := core.DBRead.Get(
		&bt,
		`SELECT business.legal_name "business.legal_name", business.dba "business.dba"
		FROM business
		WHERE business.id = $1`,
		bID,
	)
	if err != nil {
		log.Println("Error retrieving business details ", err)
		return nil, err
	}

	businessName := shared.GetBusinessName(&bt.LegalName, bt.DBA)

	return &businessName, nil

}
