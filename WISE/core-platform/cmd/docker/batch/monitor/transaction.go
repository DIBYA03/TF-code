package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	grpcTypes "github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking/business"
	bsrv "github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/services/user"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/core-platform/transaction"
	"github.com/wiseco/go-lib/id"
	"github.com/wiseco/go-lib/num"
	grpcTxn "github.com/wiseco/protobuf/golang/transaction"
	grpcBankTxn "github.com/wiseco/protobuf/golang/transaction/bank"
	grpcMonitor "github.com/wiseco/protobuf/golang/transaction/monitor"
)

func processTransaction(monitorClient grpcMonitor.BankTransactionMonitorServiceClient, txnID shared.PostedTransactionID) error {
	txn, err := transaction.NewBusinessService().GetByIDInternal(txnID)
	if err != nil {
		log.Println("fetch transaction failure")
		log.Println(err, txnID)
		return err
	}

	txnUUID, err := uuid.Parse(string(txnID))
	if err != nil {
		log.Println("Error parsing transaction", err)
		return err
	}

	sharedTXNID := id.BankTransactionID(txnUUID)

	created, err := grpcTypes.TimestampProto(txn.Created)
	if err != nil {
		log.Println(err, sharedTXNID.String())
		return err
	}

	modified := created

	txnDate, err := grpcTypes.TimestampProto(txn.TransactionDate)
	if err != nil {
		log.Println("txn date error")
		log.Println(err, sharedTXNID.String())
		return err
	}

	bID, err := id.ParseBusinessID("bus-" + string(txn.BusinessID))
	if err != nil {
		log.Println("business id", string(txn.BusinessID))
		log.Println(err, sharedTXNID.String())
		return err
	}

	b, err := bsrv.NewBusinessServiceWithout().GetByIdInternal(txn.BusinessID)
	if err != nil {
		log.Println(err, sharedTXNID.String())
		return nil
	}

	if b.KYCStatus != services.KYCStatusApproved {
		log.Printf("KYC Status not approved: %s %s", b.KYCStatus, sharedTXNID.String())
		return nil
	}

	var acc *business.BankAccount
	if txn.AccountID != nil {
		acc, err = business.NewAccountService().GetByIDInternal(*txn.AccountID)
		if err != nil {
			log.Println(err, sharedTXNID.String())
			return err
		}
	} else {
		log.Println("account id missing", sharedTXNID.String())
		return nil
	}

	if acc.UsageType != business.UsageTypePrimary {
		log.Println("account usage type must be primary", sharedTXNID.String())
		return nil
	}

	var conID id.ContactID
	if txn.ContactID != nil {
		conID, err = id.ParseContactID("cnt-" + *txn.ContactID)
		if err != nil {
			log.Println("contact id")
			log.Println(err, sharedTXNID.String())
			return err
		}
	}

	var btID id.BankTransferID
	if txn.MoneyTransferID != nil {
		btID, err = id.ParseBankTransferID("btr-" + *txn.MoneyTransferID)
		if err != nil {
			log.Println(err, sharedTXNID.String())
			return err
		}
	}

	if !btID.IsZero() {
		_, err = business.NewMoneyTransferServiceWithout().GetByIDInternal(btID.UUIDString(), shared.BusinessID(bID.UUIDString()))
		if err != nil {
			return fmt.Errorf("Error getting bank transfer: %v", err)
		}
	}

	var prID id.PaymentRequestID
	if txn.MoneyRequestID != nil {
		prID, err = id.ParsePaymentRequestID("pmr-" + string(*txn.MoneyRequestID))
		if err != nil {
			log.Println(err, sharedTXNID.String())
			return err
		}
	}

	category, ok := transaction.TransactionTypeToCategoryProto[txn.TransactionType]
	if !ok {
		log.Println("Invalid transaction category type", sharedTXNID.String())
		return nil
	}

	txnType := grpcTxn.BankTransactionType_BTT_UNSPECIFIED
	if txn.TransactionSubtype != nil {
		txnType, ok = transaction.TransactionSubtypeToTypeProto[*txn.TransactionSubtype]
		if !ok {
			log.Println("invalid transaction subtype", sharedTXNID.String())
			return nil
		}

		switch txnType {
		case grpcTxn.BankTransactionType_BTT_UNSPECIFIED:
			if acc.UsageType == business.UsageTypeClearing {
				if txn.Amount.IsNegative() {
					txnType = grpcTxn.BankTransactionType_BTT_CLEARING_DEBIT
				} else {
					txnType = grpcTxn.BankTransactionType_BTT_CLEARING_CREDIT
				}
			}
		case grpcTxn.BankTransactionType_BTT_INTRABANK_ONLINE_CREDIT:
			if category == grpcTxn.BankTransactionCategory_BTC_ACH {
				txnType = grpcTxn.BankTransactionType_BTT_ACH_ONLINE_CREDIT
			}
		}
	}

	var srcAccID id.BankAccountID
	var destAccID id.BankAccountID
	if txn.AccountID != nil {
		accUUID, err := uuid.Parse(*txn.AccountID)
		if err != nil {
			log.Println(err, sharedTXNID.String())
			return err
		}

		if txn.Amount.IsNegative() {
			srcAccID = id.BankAccountID(accUUID)
		} else {
			destAccID = id.BankAccountID(accUUID)
		}
	} else {
		log.Println("account id missing", sharedTXNID.String())
		return nil
	}

	var dbcID id.DebitCardID
	var dbcConID id.ConsumerID
	status := grpcTxn.BankTransactionStatus_BTS_NONCARD_POSTED
	if txn.CardID != nil {
		dbcUUID, err := uuid.Parse(*txn.CardID)
		if err != nil {
			log.Println(err, sharedTXNID.String())
			return err
		}

		dbcID = id.DebitCardID(dbcUUID)
		status = grpcTxn.BankTransactionStatus_BTS_CARD_POSTED

		if !dbcID.IsZero() {
			bc, err := business.NewCardServiceWithout().GetByIDInternal(dbcID.UUIDString())
			if err != nil {
				log.Println(err, sharedTXNID.String())
				return err
			}

			usr, err := user.NewUserServiceWithout().GetByIdInternal(bc.CardholderID)
			if err != nil {
				log.Println(err, sharedTXNID.String())
				return err
			}

			dbcConUUID, err := uuid.Parse(string(usr.ConsumerID))
			if err != nil {
				log.Println(err, sharedTXNID.String())
				return err
			}

			dbcConID = id.ConsumerID(dbcConUUID)
		}
	}

	creq := &grpcMonitor.BankTransactionRequest{
		Id:                     sharedTXNID.String(),
		BusinessId:             bID.String(),
		SourceAccountId:        srcAccID.String(),
		DestAccountId:          destAccID.String(),
		DebitCardId:            dbcID.String(),
		DebitCardHolderId:      dbcConID.String(),
		Amount:                 txn.Amount.FormatCurrency(),
		Currency:               string(txn.Currency),
		Status:                 status,
		Category:               category,
		Type:                   txnType,
		TransactionDate:        txnDate,
		BankTransferId:         btID.String(),
		PaymentRequestId:       prID.String(),
		ContactId:              conID.String(),
		PartnerTransactionId:   shared.StringValue(txn.BankTransactionID),
		PartnerTransactionDesc: shared.StringValue(txn.BankTransactionDesc),
		TransactionDesc:        shared.StringValue(txn.TransactionTitle),
		Created:                created,
		Modified:               modified,
	}

	resp, err := monitorClient.AddUpdateTransaction(context.Background(), creq)
	if err != nil {
		log.Println(err, sharedTXNID.String())
		return err
	}

	log.Println("Success: ", resp.Id)
	return nil
}

func processServiceTransaction(monitorClient grpcMonitor.BankTransactionMonitorServiceClient, txn *grpcBankTxn.Transaction) error {
	// Parse business
	busID, err := id.ParseBusinessID(txn.BusinessId)
	if err != nil {
		log.Println(err, txn.Id)
		return err
	}

	b, err := bsrv.NewBusinessServiceWithout().GetByIdInternal(shared.BusinessID(busID.UUIDString()))
	if err != nil {
		log.Println(err, txn.Id)
		return err
	}

	// Parse amount
	amount, err := num.ParseDecimal(txn.Amount)
	if err != nil {
		log.Println(err, txn.Id)
		return err
	}

	// Parse account id
	accID, err := id.ParseBankAccountID(txn.AccountId)
	if err != nil {
		log.Println(err, txn.Id)
		return err
	}

	var srcAccID, destAccID string
	if !accID.IsZero() {
		if amount.IsNegative() {
			srcAccID = txn.AccountId
		} else {
			destAccID = txn.AccountId
		}
	} else {
		log.Println("account id missing", txn.Id)
		return nil
	}

	// Parse denit card
	dbcID, err := id.ParseDebitCardID(txn.DebitCardId)
	if err != nil {
		log.Println(err, txn.Id)
		return err
	}

	var dbcConID string
	if !dbcID.IsZero() {
		bc, err := business.NewCardServiceWithout().GetByIDInternal(dbcID.UUIDString())
		if err != nil {
			log.Println(err, txn.Id)
			return err
		}

		usr, err := user.NewUserServiceWithout().GetByIdInternal(bc.CardholderID)
		if err != nil {
			log.Println(err, txn.Id)
			return err
		}

		dbcConUUID, err := uuid.Parse(string(usr.ConsumerID))
		if err != nil {
			log.Println(err, txn.Id)
			return err
		}

		dbcConID = id.ConsumerID(dbcConUUID).String()
	}

	creq := &grpcMonitor.BankTransactionRequest{
		Id:                     txn.Id,
		BusinessId:             txn.BusinessId,
		SourceAccountId:        srcAccID,
		DestAccountId:          destAccID,
		DebitCardId:            txn.DebitCardId,
		DebitCardHolderId:      dbcConID,
		Amount:                 txn.Amount,
		Currency:               txn.Currency,
		Status:                 txn.Status,
		Category:               txn.Category,
		Type:                   txn.Type,
		TransactionDate:        txn.TransactionDate,
		BankTransferId:         txn.BankTransferId,
		PaymentRequestId:       txn.PaymentRequestId,
		ContactId:              txn.ContactId,
		PartnerTransactionId:   txn.PartnerTransactionId,
		PartnerTransactionDesc: txn.PartnerTransactionDesc,
		TransactionDesc:        transaction.BankTransactionDisplayTitle(txn, b),
		Created:                txn.Created,
		Modified:               txn.Modified,
	}

	resp, err := monitorClient.AddUpdateTransaction(context.Background(), creq)
	if err != nil {
		log.Println(err, txn.Id)
		return err
	}

	log.Println("Success: ", resp.Id)
	return nil
}

func sendTransactionUpdates(monitorClient grpcMonitor.BankTransactionMonitorServiceClient, txnClient grpcBankTxn.TransactionServiceClient, dayStart, dayEnd time.Time) {
	// Process in groups of 5
	offset := 0
	limit := 5
	for {
		if os.Getenv("USE_TRANSACTION_SERVICE") == "true" {
			start, err := grpcTypes.TimestampProto(dayStart)
			if err != nil {
				panic(err)
			}

			end, err := grpcTypes.TimestampProto(dayEnd)
			if err != nil {
				panic(err)
			}

			dateRange := &grpcBankTxn.DateRange{
				Filter: grpcBankTxn.DateRangeFilter_DRF_START_END,
				Start:  start,
				End:    end,
			}

			req := &grpcBankTxn.TransactionsRequest{
				Offset:    int32(offset),
				Limit:     int32(limit),
				DateRange: dateRange,
			}

			resp, err := txnClient.GetMany(context.Background(), req)
			if err != nil {
				panic(err)
			} else if len(resp.Results) == 0 {
				log.Println("No more transactions", dayStart, dayEnd)
				break
			}

			wg := sync.WaitGroup{}
			wg.Add(len(resp.Results))
			for _, txn := range resp.Results {
				go func(txn *grpcBankTxn.Transaction) {
					defer wg.Done()
					_ = processServiceTransaction(monitorClient, txn)
				}(txn)
			}

			wg.Wait()
		} else {
			var txnIDs []shared.PostedTransactionID
			err := transaction.DBWrite.Select(
				&txnIDs,
				`
			SELECT id FROM business_transaction
			WHERE
				created >= $1 AND created < $2
			ORDER BY created ASC OFFSET $3 LIMIT $4`,
				dayStart,
				dayEnd,
				offset,
				limit,
			)
			if err != nil {
				panic(err)
			} else if len(txnIDs) == 0 {
				log.Println("No more transactions", dayStart, dayEnd)
				break
			}

			wg := sync.WaitGroup{}
			wg.Add(len(txnIDs))
			for _, txnID := range txnIDs {
				go func(id shared.PostedTransactionID) {
					defer wg.Done()
					_ = processTransaction(monitorClient, id)
				}(txnID)
			}

			wg.Wait()
		}

		offset += 5
	}
}
