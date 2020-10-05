package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/wiseco/go-lib/grpc"
	grpcBankTxn "github.com/wiseco/protobuf/golang/transaction/bank"
	grpcMonitor "github.com/wiseco/protobuf/golang/transaction/monitor"
)

func main() {
	var monitorClient grpcMonitor.BankTransactionMonitorServiceClient
	var txnClient grpcBankTxn.TransactionServiceClient

	var dayStart, dayEnd, dayStartLocal, dayEndLocal time.Time

	if os.Getenv("USE_TRANSACTION_SERVICE") != "true" {
		return
	}

	// Get time zone and determine start/end
	tz := os.Getenv("BATCH_TZ")
	if tz == "" {
		panic(errors.New("local timezone missing"))
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		panic(err)
	}

	utcLoc, err := time.LoadLocation("UTC")
	if err != nil {
		panic(err)
	}

	nowUTC := time.Now().UTC()
	nowLocal := nowUTC.In(loc)

	dayEndLocal = time.Date(nowLocal.Year(), nowLocal.Month(), nowLocal.Day(), 0, 0, 0, 0, loc)
	dayStartLocal = dayEndLocal.AddDate(0, 0, -1)

	dayStart = dayStartLocal.In(utcLoc)
	dayEnd = dayEndLocal.In(utcLoc)
	if dayEnd.After(nowUTC) {
		panic(fmt.Errorf("Error: day end (%v) is after current time (%v)", dayEnd, nowUTC))
	}

	dayStart = dayEnd.AddDate(0, 0, -1)
	dayStartLocal = dayStart.In(loc)

	sn, err := grpc.GetConnectionStringForService(grpc.ServiceNameTransaction)
	if err != nil {
		panic(err)
	}

	client, err := grpc.NewInsecureClient(sn)
	if err != nil {
		panic(err)
	}

	defer client.CloseAndCancel()
	monitorClient = grpcMonitor.NewBankTransactionMonitorServiceClient(client.GetConn())
	txnClient = grpcBankTxn.NewTransactionServiceClient(client.GetConn())

	// Send updates
	sendConsumerUpdates(monitorClient, dayStart, dayEnd)
	sendBusinessUpdates(monitorClient, dayStart, dayEnd)
	sendAccountUpdates(monitorClient, dayStart, dayEnd)
	sendDebitCardUpdates(monitorClient, dayStart, dayEnd)
	sendTransactionUpdates(monitorClient, txnClient, dayStart, dayEnd)
}
