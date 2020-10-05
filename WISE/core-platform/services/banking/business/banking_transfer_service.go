package business

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"github.com/wiseco/core-platform/partner/bank/bbva"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/go-lib/grpc"
	"github.com/wiseco/go-lib/id"
	"github.com/wiseco/go-lib/num"
	grpcRoot "github.com/wiseco/protobuf/golang"
	grpcBanking "github.com/wiseco/protobuf/golang/banking"
	"github.com/wiseco/protobuf/golang/banking/transfer"
	grpcBankTransfer "github.com/wiseco/protobuf/golang/banking/transfer"
	grpcBankTxn "github.com/wiseco/protobuf/golang/transaction/bank"
	"github.com/xtgo/uuid"
)

type bankingTransferService struct {
	conn grpc.Client
	tc   grpcBankTransfer.TransferServiceClient
}

func NewBankingTransferService() (*bankingTransferService, error) {
	var bts *bankingTransferService

	tc, conn, err := getNewClientAndConn()
	if err != nil {
		return bts, err
	}

	return &bankingTransferService{
		conn: conn,
		tc:   tc,
	}, nil
}

func getNewClientAndConn() (grpcBankTransfer.TransferServiceClient, grpc.Client, error) {
	var tc grpcBankTransfer.TransferServiceClient
	var conn grpc.Client

	bsn, err := grpc.GetConnectionStringForService(grpc.ServiceNameBanking)
	if err != nil {
		return tc, conn, err
	}

	conn, err = grpc.NewInsecureClient(bsn)
	if err != nil {
		return tc, conn, err
	}

	tc = grpcBankTransfer.NewTransferServiceClient(conn.GetConn())

	return tc, conn, nil
}

func (bts bankingTransferService) GetByBusinessID(bID shared.BusinessID, offset, limit int) ([]MoneyTransfer, error) {
	var mts []MoneyTransfer

	defer bts.conn.CloseAndCancel()

	gmr := &grpcBankTransfer.GetManyRequest{
		BusinessId: string(bID),
		Limit:      int32(limit),
		Offset:     int32(offset),
	}

	ts, err := bts.tc.GetMany(context.Background(), gmr)
	if err != nil {
		return mts, err
	}

	for _, t := range ts.Transfers {
		mt, err := transformProtoTransferToMoneyTransfer(t)
		if err != nil {
			return mts, err
		}

		mts = append(mts, *mt)
	}

	return mts, nil
}

func (bts bankingTransferService) GetByBusinessAndContactID(bID shared.BusinessID, cID string, offset, limit int) ([]MoneyTransfer, error) {
	var mts []MoneyTransfer

	defer bts.conn.CloseAndCancel()

	cUUID, err := uuid.Parse(cID)
	if err != nil {
		return mts, err
	}

	gmr := &grpcBankTransfer.GetManyRequest{
		BusinessId: string(bID),
		ContactId:  id.ContactID(cUUID).String(),
		Limit:      int32(limit),
		Offset:     int32(offset),
	}

	ts, err := bts.tc.GetMany(context.Background(), gmr)
	if err != nil {
		return mts, err
	}

	for _, t := range ts.Transfers {
		mt, err := transformProtoTransferToMoneyTransfer(t)
		if err != nil {
			return mts, err
		}

		mts = append(mts, *mt)
	}

	return mts, nil
}

func (bts bankingTransferService) GetByBankID(bID shared.BusinessID, ptID string) (*MoneyTransfer, error) {
	var mt *MoneyTransfer

	defer bts.conn.CloseAndCancel()

	gr := &grpcBankTransfer.GetRequest{
		PartnerTransferId: ptID,
		PartnerName:       grpcBanking.PartnerName_PN_BBVA,
	}

	t, err := bts.tc.Get(context.Background(), gr)
	if err != nil {
		return mt, err
	}

	return transformProtoTransferToMoneyTransfer(t)
}

func (bts bankingTransferService) GetByIDInternal(bID shared.BusinessID, tID string) (*MoneyTransfer, error) {
	var mt *MoneyTransfer

	defer bts.conn.CloseAndCancel()

	s := tID
	if !strings.HasPrefix(tID, id.IDPrefixBankTransfer.String()) {
		s = fmt.Sprintf("%s%s", id.IDPrefixBankTransfer, tID)
	}

	btID, err := id.ParseBankTransferID(s)
	if err != nil {
		return mt, err
	}

	gr := &grpcBankTransfer.GetRequest{
		TransferId: btID.String(),
	}

	t, err := bts.tc.Get(context.Background(), gr)
	if err != nil {
		return mt, err
	}

	return transformProtoTransferToMoneyTransfer(t)
}

func (bts bankingTransferService) GetProtoByTransactionID(txnID string) (*grpcBankTransfer.Transfer, error) {
	t := new(grpcBankTransfer.Transfer)

	defer bts.conn.CloseAndCancel()

	tsn, err := grpc.GetConnectionStringForService(grpc.ServiceNameTransaction)
	if err != nil {
		return t, err
	}

	tc, err := grpc.NewInsecureClient(tsn)

	defer tc.CloseAndCancel()

	if err != nil {
		return t, err
	}

	txnClient := grpcBankTxn.NewTransactionServiceClient(tc.GetConn())

	s := txnID

	if !strings.HasPrefix(s, id.IDPrefixBankTransaction.String()) {
		//This should be thought about further:
		//We are passing pnt(pending transaction ids around) but the transaction service does not recognize that prefix
		if strings.HasPrefix(s, id.IDPrefixPendingTransaction.String()) {
			//Remove pnt prefix
			s = s[4:]
		}

		s = fmt.Sprintf("%s%s", id.IDPrefixBankTransaction, s)
	}

	transactionID, err := id.ParseBankTransactionID(s)
	if err != nil {
		return t, err
	}

	req := &grpcBankTxn.TransactionIDRequest{
		Id: transactionID.String(),
	}

	resp, err := txnClient.GetByID(context.Background(), req)
	if err != nil {
		return t, err
	}

	gr := &grpcBankTransfer.GetRequest{
		TransferId: resp.BankTransferId,
	}

	return bts.tc.Get(context.Background(), gr)
}

func (bts bankingTransferService) GetByIDOnlyInternal(tID string) (*MoneyTransfer, error) {
	var mt *MoneyTransfer

	defer bts.conn.CloseAndCancel()

	s := tID
	if !strings.HasPrefix(tID, id.IDPrefixBankTransfer.String()) {
		s = fmt.Sprintf("%s%s", id.IDPrefixBankTransfer, tID)
	}

	btID, err := id.ParseBankTransferID(s)
	if err != nil {
		return mt, err
	}

	gr := &grpcBankTransfer.GetRequest{
		TransferId: btID.String(),
	}

	t, err := bts.tc.Get(context.Background(), gr)
	if err != nil {
		return mt, err
	}

	return transformProtoTransferToMoneyTransfer(t)
}

func (bts bankingTransferService) GetByPaymentRequestID(prID shared.PaymentRequestID) (*MoneyTransfer, error) {
	var mt *MoneyTransfer

	defer bts.conn.CloseAndCancel()

	gr := &grpcBankTransfer.GetRequest{
		PaymentRequestId: prID.ToPrefixString(),
	}

	t, err := bts.tc.Get(context.Background(), gr)
	if err != nil {
		if strings.Contains(err.Error(), "Cannot find transfer to get") {
			return mt, sql.ErrNoRows
		}

		return mt, err
	}

	return transformProtoTransferToMoneyTransfer(t)
}

func (bts bankingTransferService) GetByAccountIDPartnerTransferStatusAndAmount(aID string, pts []grpcBankTransfer.PartnerTransferStatus, amount num.Decimal) (*MoneyTransfer, error) {
	var mt *MoneyTransfer

	defer bts.conn.CloseAndCancel()

	s := aID
	if !strings.HasPrefix(aID, id.IDPrefixBankAccount.String()) {
		s = fmt.Sprintf("%s%s", id.IDPrefixBankAccount, aID)
	}

	paID, err := id.ParseBankAccountID(s)
	if err != nil {
		return mt, err
	}

	gr := &grpcBankTransfer.GetRequest{
		AccountId:             paID.String(),
		PartnerTransferStatus: pts,
		Amount:                amount.String(),
	}

	t, err := bts.tc.Get(context.Background(), gr)
	if err != nil {
		if strings.Contains(err.Error(), "Cannot find transfer to get") {
			return mt, sql.ErrNoRows
		}

		return mt, err
	}

	return transformProtoTransferToMoneyTransfer(t)
}

func (bts bankingTransferService) UpdateStatus(bID shared.BusinessID, ptID, status string) error {
	defer bts.conn.CloseAndCancel()

	gr := &grpcBankTransfer.GetRequest{
		PartnerTransferId: ptID,
		PartnerName:       grpcBanking.PartnerName_PN_BBVA,
	}

	t, err := bts.tc.Get(context.Background(), gr)
	if err != nil {
		return err
	}

	ts := grpcBankTransfer.TransferStatus_TS_UNSPECIFIED
	pts := grpcBankTransfer.PartnerTransferStatus_TPS_UNSPECIFIED

	switch bbva.MoveMoneyStatus(status) {
	case bbva.MoveMoneyStatusInProcess:
		pts = grpcBankTransfer.PartnerTransferStatus_TPS_IN_PROCESS
		ts = grpcBankTransfer.TransferStatus_TS_BANK_PROCESSING
	case bbva.MoveMoneyStatusCancelled:
		pts = grpcBankTransfer.PartnerTransferStatus_TPS_CANCELED
		ts = grpcBankTransfer.TransferStatus_TS_USER_CANCELED
	case bbva.MoveMoneyStatusPosted:
		pts = grpcBankTransfer.PartnerTransferStatus_TPS_POSTED
		ts = grpcBankTransfer.TransferStatus_TS_POSTED
	case bbva.MoveMoneyStatusDebitSent:
		pts = grpcBankTransfer.PartnerTransferStatus_TPS_DEBIT_SENT
		ts = grpcBankTransfer.TransferStatus_TS_BANK_PROCESSING
	case bbva.MoveMoneyStatusCreditSent:
		pts = grpcBankTransfer.PartnerTransferStatus_TPS_CREDIT_SENT
		ts = grpcBankTransfer.TransferStatus_TS_BANK_PROCESSING
	case bbva.MoveMoneyStatusSettled:
		pts = grpcBankTransfer.PartnerTransferStatus_TPS_SETTLED
		ts = grpcBankTransfer.TransferStatus_TS_POSTED
	case bbva.MoveMoneyStatusDisbursed:
		pts = grpcBankTransfer.PartnerTransferStatus_TPS_DISBURSED
		ts = grpcBankTransfer.TransferStatus_TS_POSTED
	case bbva.MoveMoneyStatusCheckCleared:
		pts = grpcBankTransfer.PartnerTransferStatus_TPS_CHECK_CLEARED
		ts = grpcBankTransfer.TransferStatus_TS_POSTED
	case bbva.MoveMoneyStatusPullFailed:
		pts = grpcBankTransfer.PartnerTransferStatus_TPS_PULL_FAILED
		ts = grpcBankTransfer.TransferStatus_TS_BANK_DECLINED
	case bbva.MoveMoneyStatusPullFailedRefunded:
		pts = grpcBankTransfer.PartnerTransferStatus_TPS_PULL_FAILED_REFUNDED
		ts = grpcBankTransfer.TransferStatus_TS_BANK_REFUNDED
	case bbva.MoveMoneyStatusPullFailedUnderReview:
		pts = grpcBankTransfer.PartnerTransferStatus_TPS_PULL_FAILED_UNDER_REVIEW
		ts = grpcBankTransfer.TransferStatus_TS_BANK_IN_REVIEW
	case bbva.MoveMoneyStatusPushFailedRefunded:
		pts = grpcBankTransfer.PartnerTransferStatus_TPS_PUSH_FAILED_REFUNDED
		ts = grpcBankTransfer.TransferStatus_TS_BANK_REFUNDED
	case bbva.MoveMoneyStatusPushFailedUnderReview:
		pts = grpcBankTransfer.PartnerTransferStatus_TPS_PUSH_FAILED_UNDER_REVIEW
		ts = grpcBankTransfer.TransferStatus_TS_BANK_IN_REVIEW
	}

	ur := &grpcBankTransfer.UpdateRequest{
		TransferId:            t.Id,
		TransferStatus:        ts,
		PartnerTransferStatus: pts,
	}

	_, err = bts.tc.Update(context.Background(), ur)

	return err
}

func (bts bankingTransferService) UpdateDebitPostedTransaction(bID shared.BusinessID, tID, postedDebitTransactionID string) error {
	defer bts.conn.CloseAndCancel()

	s := tID
	if !strings.HasPrefix(tID, id.IDPrefixBankTransfer.String()) {
		s = fmt.Sprintf("%s%s", id.IDPrefixBankTransfer, tID)
	}

	btID, err := id.ParseBankTransferID(s)
	if err != nil {
		return err
	}

	ur := &grpcBankTransfer.UpdateRequest{
		TransferId:               btID.String(),
		PostedDebitTransactionId: postedDebitTransactionID,
	}

	_, err = bts.tc.Update(context.Background(), ur)

	return err
}

func (bts bankingTransferService) UpdateCreditPostedTransaction(tID, postedCreditTransactionID string) error {
	defer bts.conn.CloseAndCancel()

	s := tID
	if !strings.HasPrefix(tID, id.IDPrefixBankTransfer.String()) {
		s = fmt.Sprintf("%s%s", id.IDPrefixBankTransfer, tID)
	}

	btID, err := id.ParseBankTransferID(s)
	if err != nil {
		return err
	}

	ur := &grpcBankTransfer.UpdateRequest{
		TransferId:                btID.String(),
		PostedCreditTransactionId: postedCreditTransactionID,
	}

	_, err = bts.tc.Update(context.Background(), ur)

	return err
}

func (bts bankingTransferService) Transfer(ti *TransferInitiate, sut UsageType, sourceReq services.SourceRequest) (*MoneyTransfer, error) {
	var mt *MoneyTransfer
	var t *grpcBankTransfer.Transfer
	var err error

	defer bts.conn.CloseAndCancel()

	conID := ""
	if ti.ContactId != nil {
		cUUID, err := uuid.Parse(*ti.ContactId)
		if err != nil {
			return mt, err
		}

		conID = id.ContactID(cUUID).String()
	}

	notes := ""

	if ti.Notes != nil {
		notes = *ti.Notes
	}

	miID := ""
	if ti.MonthlyInterestID != nil {
		tmiID, err := id.ParseMonthlyInterestID(fmt.Sprintf("%s%s", id.IDPrefixMonthlyInterest, *ti.MonthlyInterestID))
		if err != nil {
			return mt, err
		}

		miID = tmiID.String()
	}

	prID := ""
	if ti.MoneyRequestID != nil {
		prID = ti.MoneyRequestID.ToPrefixString()
	}

	//TODO this is a hack
	//We should be creating and getting the payee before this
	//We need to update the api and have the front end use those endpoints
	if ti.DestType == banking.TransferTypeCheck {
		if ti.AddressID == "" {
			return mt, errors.New("Address ID is required to pay by check")
		}

		lps := NewLinkedPayeeService(sourceReq)

		payee, err := lps.GetByAddressID(ti.AddressID, ti.BusinessID)
		if err != nil {
			if err == sql.ErrNoRows {
				//TODO this is a bad hack, the client should be explicitly creating a linked_payee
				//instead of relying on addressID
				lpc := LinkedPayeeCreate{
					BusinessID: ti.BusinessID,
					ContactID:  *ti.ContactId,
					AddressID:  ti.AddressID,
				}

				payee, err = lps.Create(lpc)
				if err != nil {
					return mt, err
				}
			} else {
				return mt, err
			}
		}

		ti.DestAccountId = payee.ID
	}

	st, saID := getAccountTypeAndIDFromTransferType(ti.SourceType, ti.SourceAccountId)
	dt, daID := getAccountTypeAndIDFromTransferType(ti.DestType, ti.DestAccountId)

	er := &grpcBankTransfer.ExecuteRequest{
		BusinessId:        ti.BusinessID.ToPrefixString(),
		ContactId:         conID,
		MonthlyInterestId: miID,
		PaymentRequestId:  prID,
		SourceAccountId:   saID,
		SourceType:        st,
		DestAccountId:     daID,
		DestType:          dt,
		Amount:            fmt.Sprintf("%f", ti.Amount),
		Currency:          string(banking.CurrencyUSD),
		Notes:             notes,
		SendEmail:         ti.SendEmail,
		ActorId:           ti.BusinessID.ToPrefixString(),
		ActorType:         grpcRoot.ActorType_AT_BUSINESS,
	}

	if ti.CVVCode != nil {
		er.Cvv2CvcCode = *ti.CVVCode
	}

	tt := getTransferType(ti, sut)

	switch tt {
	case grpcBankTransfer.TransferType_TT_ACH_PUSH:
		t, err = bts.tc.ExecuteACHPush(context.Background(), er)
	case grpcBankTransfer.TransferType_TT_ACH_PULL:
		t, err = bts.tc.ExecuteACHPull(context.Background(), er)
	case grpcBankTransfer.TransferType_TT_DEBIT_PUSH:
		t, err = bts.tc.ExecuteDebitPush(context.Background(), er)
	case grpcBankTransfer.TransferType_TT_DEBIT_PULL:
		t, err = bts.tc.ExecuteDebitPull(context.Background(), er)
	case grpcBankTransfer.TransferType_TT_CHECK:
		t, err = bts.tc.ExecuteCheck(context.Background(), er)
	default:
		return mt, fmt.Errorf("Could not figure out transfer type for transfer initiate:%v", ti)
	}

	if err != nil {
		return mt, err
	}

	mt, err = transformProtoTransferToMoneyTransfer(t)
	if err != nil {
		return mt, err
	}

	if t.TransferStatus == grpcBankTransfer.TransferStatus_TS_AUTO_DECLINED {
		code := int32(1000)
		for _, f := range t.Failures {
			if f.Code > 200 {
				//The lower the code the more pressing the error message
				if f.Code < code {
					err = errors.New(f.Message)
					code = f.Code
				}
			}
		}
	}

	return mt, err
}

func (bts bankingTransferService) Approve(txnID, cspUserID, ipAddress string) error {
	defer bts.conn.CloseAndCancel()

	trp, err := bts.GetProtoByTransactionID(txnID)
	if err != nil {
		return err
	}

	if trp.TransferStatus != grpcBankTransfer.TransferStatus_TS_AGENT_IN_REVIEW {
		return errors.New("Only transfers with status agent review can be approved")
	}

	tc, conn, err := getNewClientAndConn()

	defer conn.CloseAndCancel()

	if err != nil {
		return err
	}

	for _, f := range trp.Failures {
		aReq := &transfer.CloseValidatorFailureRequest{
			ValidatorFailureId: f.Id,
			CspUserId:          cspUserID,
			IpAddress:          ipAddress,
		}

		_, err := tc.CloseValidatorFailure(context.Background(), aReq)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	s := txnID

	if !strings.HasPrefix(s, id.IDPrefixBankTransaction.String()) {
		//This should be thought about further:
		//We are passing pnt(pending transaction ids around) but the transaction service does not recognize that prefix
		if strings.HasPrefix(s, id.IDPrefixPendingTransaction.String()) {
			//Remove pnt prefix
			s = s[4:]
		}

		s = fmt.Sprintf("%s%s", id.IDPrefixBankTransaction, s)
	}

	transactionID, err := id.ParseBankTransactionID(s)
	if err != nil {
		return err
	}

	aer := &grpcBankTransfer.AgentExecuteRequest{
		TransferId:    trp.Id,
		TransactionId: transactionID.String(),
		CspUserId:     cspUserID,
	}

	_, err = tc.AgentExecute(context.Background(), aer)

	return err
}

func (bts bankingTransferService) Decline(txnID, cspUserID, ipAddress string) error {
	defer bts.conn.CloseAndCancel()

	trp, err := bts.GetProtoByTransactionID(txnID)
	if err != nil {
		return err
	}

	s := txnID

	if !strings.HasPrefix(s, id.IDPrefixBankTransaction.String()) {
		//This should be thought about further:
		//We are passing pnt(pending transaction ids around) but the transaction service does not recognize that prefix
		if strings.HasPrefix(s, id.IDPrefixPendingTransaction.String()) {
			//Remove pnt prefix
			s = s[4:]
		}

		s = fmt.Sprintf("%s%s", id.IDPrefixBankTransaction, s)
	}

	transactionID, err := id.ParseBankTransactionID(s)
	if err != nil {
		return err
	}

	adr := &grpcBankTransfer.AgentDeclineRequest{
		TransferId:    trp.Id,
		TransactionId: transactionID.String(),
		CspUserId:     cspUserID,
		IpAddress:     ipAddress,
	}

	tc, conn, err := getNewClientAndConn()

	defer conn.CloseAndCancel()

	if err != nil {
		return err
	}

	_, err = tc.AgentDecline(context.Background(), adr)

	return err
}

func transformProtoTransferToMoneyTransfer(t *grpcBankTransfer.Transfer) (*MoneyTransfer, error) {
	var miID *string
	var mrID *shared.PaymentRequestID

	mt := new(MoneyTransfer)

	if t.Id == "" {
		return mt, services.ErrorNotFound{}.New("")
	}

	mtID, err := id.ParseBankTransferID(t.Id)
	if err != nil {
		return mt, err
	}

	bID, err := shared.ParseBusinessID(t.BusinessId)
	if err != nil {
		return mt, err
	}

	pctID, err := id.ParseBankTransactionID(t.PostedCreditTransactionId)
	if err != nil {
		return mt, err
	}

	var pctUUID *string

	if !pctID.IsZero() {
		s := pctID.UUIDString()
		pctUUID = &s
	}

	pdtID, err := id.ParseBankTransactionID(t.PostedDebitTransactionId)
	if err != nil {
		return mt, err
	}

	var pdtUUID *string
	if !pdtID.IsZero() {
		s := pdtID.UUIDString()
		pdtUUID = &s
	}

	conID, err := id.ParseContactID(t.ContactId)
	if err != nil {
		return mt, err
	}

	var conUUID *string

	if !conID.IsZero() {
		s := conID.UUIDString()
		conUUID = &s
	}

	created, err := ptypes.Timestamp(t.Created)
	if err != nil {
		return mt, err
	}

	if t.MonthlyInterestId != "" {
		tempID, err := id.ParseMonthlyInterestID(t.MonthlyInterestId)
		if err != nil {
			return mt, err
		}

		if !tempID.IsZero() {
			s := tempID.UUIDString()
			miID = &s
		}
	}

	if t.PaymentRequestId != "" {
		tempID, err := id.ParsePaymentRequestID(t.PaymentRequestId)
		if err != nil {
			return mt, err
		}

		if !tempID.IsZero() {
			s := shared.PaymentRequestID(tempID.UUIDString())
			mrID = &s
		}
	}

	st := getTransferTypeFromProtoAccountType(t.SourceType)
	dt := getTransferTypeFromProtoAccountType(t.DestType)
	status := getTransferStatusFromProto(t.TransferStatus)

	amount, err := strconv.ParseFloat(t.Amount, 64)
	if err != nil {
		return mt, err
	}

	mt.Id = mtID.UUIDString()
	mt.BusinessID = bID
	mt.MonthlyInterestID = miID
	mt.PostedDebitTransactionID = pdtUUID
	mt.PostedCreditTransactionID = pctUUID
	mt.MoneyRequestID = mrID
	mt.ContactId = conUUID
	mt.BankName = "bbva"
	mt.BankTransferId = t.PartnerTransferId
	mt.SourceAccountId = t.SourceAccountId
	mt.SourceType = st
	mt.DestAccountId = t.DestAccountId
	mt.DestType = dt
	mt.Amount = amount
	mt.Currency = banking.CurrencyUSD
	mt.Notes = &t.Notes
	mt.Status = status
	mt.SendEmail = t.SendEmail
	mt.Created = created
	mt.ErrorCause = t.ErrorCause

	return mt, nil
}

func getTransferStatusFromProto(tsp grpcBankTransfer.TransferStatus) string {
	var ts string

	switch tsp {
	case grpcBankTransfer.TransferStatus_TS_BANK_PROCESSING:
		ts = banking.MoneyTransferStatusInProcess
	case grpcBankTransfer.TransferStatus_TS_POSTED:
		ts = banking.MoneyTransferStatusPosted
	case grpcBankTransfer.TransferStatus_TS_AGENT_IN_REVIEW:
		ts = banking.MoneyTransferStatusPending
	case grpcBankTransfer.TransferStatus_TS_AUTO_DECLINED:
		ts = banking.MoneyTransferStatusCanceled
	case grpcBankTransfer.TransferStatus_TS_BANK_ERROR:
		ts = banking.MoneyTransferStatusBankError
	}

	return ts
}

func getAccountTypeAndIDFromTransferType(tt banking.TransferType, ID string) (grpcBankTransfer.AccountType, string) {
	var at grpcBankTransfer.AccountType

	switch tt {
	case banking.TransferTypeAccount:
		at = grpcBankTransfer.AccountType_AT_LINKED_ACCOUNT

		if !strings.HasPrefix(ID, id.IDPrefixLinkedBankAccount.String()) && !strings.HasPrefix(ID, id.IDPrefixBankAccount.String()) {
			ID = id.IDPrefixLinkedBankAccount.String() + ID
		}
	case banking.TransferTypeCard:
		at = grpcBankTransfer.AccountType_AT_LINKED_CARD

		if !strings.HasPrefix(ID, id.IDPrefixLinkedCard.String()) {
			ID = id.IDPrefixLinkedCard.String() + ID
		}
	case banking.TransferTypeCheck:
		at = grpcBankTransfer.AccountType_AT_LINKED_PAYEE

		if !strings.HasPrefix(ID, id.IDPrefixLinkedPayee.String()) {
			ID = id.IDPrefixLinkedPayee.String() + ID
		}
	}

	return at, ID
}

func getTransferTypeFromProtoAccountType(at grpcBankTransfer.AccountType) banking.TransferType {
	var tt banking.TransferType

	switch at {
	case grpcBankTransfer.AccountType_AT_ACCOUNT_BUSINESS, grpcBankTransfer.AccountType_AT_LINKED_ACCOUNT:
		tt = banking.TransferTypeAccount
	case grpcBankTransfer.AccountType_AT_LINKED_CARD:
		tt = banking.TransferTypeCard
	case grpcBankTransfer.AccountType_AT_LINKED_PAYEE:
		tt = banking.TransferTypeCheck
	}

	return tt
}

func getTransferType(ti *TransferInitiate, sut UsageType) grpcBankTransfer.TransferType {
	tt := grpcBankTransfer.TransferType_TT_UNSPECIFIED

	if ti.SourceType == banking.TransferTypeAccount {
		switch ti.DestType {
		case banking.TransferTypeAccount:
			if sut == UsageTypePrimary || sut == UsageTypeClearing {
				tt = grpcBankTransfer.TransferType_TT_ACH_PUSH
			} else {
				tt = grpcBankTransfer.TransferType_TT_ACH_PULL
			}
		case banking.TransferTypeCard:
			tt = grpcBankTransfer.TransferType_TT_DEBIT_PUSH
		case banking.TransferTypeCheck:
			tt = grpcBankTransfer.TransferType_TT_CHECK
		}
	}

	if ti.SourceType == banking.TransferTypeCard {
		switch ti.DestType {
		case banking.TransferTypeAccount:
			tt = grpcBankTransfer.TransferType_TT_DEBIT_PULL
		}
	}

	return tt
}

func getSourceUsageType(ti *TransferInitiate) UsageType {
	var ut UsageType

	ut = UsageTypePrimary

	return ut
}
