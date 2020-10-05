package business

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/go-lib/grpc"
	"github.com/wiseco/go-lib/id"
	grpcRoot "github.com/wiseco/protobuf/golang"
	grpcBanking "github.com/wiseco/protobuf/golang/banking"
	grpcBankLinkedAccount "github.com/wiseco/protobuf/golang/banking/linked_account"
	"github.com/xtgo/uuid"
)

type bankingLinkedAccountService struct {
	conn grpc.Client
	lac  grpcBankLinkedAccount.LinkedAccountServiceClient
}

func NewBankingLinkedAccountService() (*bankingLinkedAccountService, error) {
	var bts *bankingLinkedAccountService

	bsn, err := grpc.GetConnectionStringForService(grpc.ServiceNameBanking)
	if err != nil {
		return bts, err
	}

	conn, err := grpc.NewInsecureClient(bsn)
	if err != nil {
		return bts, err
	}

	lac := grpcBankLinkedAccount.NewLinkedAccountServiceClient(conn.GetConn())

	return &bankingLinkedAccountService{
		conn: conn,
		lac:  lac,
	}, nil
}

func (blas bankingLinkedAccountService) LinkBankAccount(la *LinkedBankAccount, cID shared.ConsumerID) (*LinkedBankAccount, error) {
	var lba *LinkedBankAccount

	defer blas.conn.CloseAndCancel()

	conID := ""
	if la.ContactId != nil {
		cUUID, err := uuid.Parse(*la.ContactId)
		if err != nil {
			return lba, err
		}

		conID = id.ContactID(cUUID).String()
	}

	aID := ""
	if la.BusinessBankAccountId != nil {
		baID, err := id.ParseBankAccountID(fmt.Sprintf("%s%s", id.IDPrefixBankAccount, *la.BusinessBankAccountId))
		if err != nil {
			return lba, err
		}

		aID = baID.String()
	}

	bn := ""
	if la.BankName != nil {
		bn = *la.BankName
	}

	eln := grpcBankLinkedAccount.ExternalLinkName_ELN_S_UNSPECIFIED
	elrID := ""
	elID := ""

	if la.SourceId != nil {
		elrID = *la.SourceId
		elID = *la.SourceAccountId
		eln = grpcBankLinkedAccount.ExternalLinkName_ELN_PLAID
	}

	alias := ""
	if la.Alias != nil {
		alias = *la.Alias
	}

	accountName := ""
	if la.AccountName != nil {
		accountName = *la.AccountName
	}

	ut := getProtoUsageTypeFromPermission(la.Permission)
	lat := grpcBankLinkedAccount.LinkedAccountType_LAT_BUSINESS_CHECKING
	last := getProtoSubTypeFromUsageType(la.UsageType)

	rr := &grpcBankLinkedAccount.RegisterRequest{
		ConsumerId:            cID.ToPrefixString(),
		BusinessId:            la.BusinessID.ToPrefixString(),
		ContactId:             conID,
		AccountId:             aID,
		PartnerName:           grpcBanking.PartnerName_PN_BBVA,
		IssuerAccountNumber:   la.AccountNumber.String(),
		IssuerRoutingNumber:   la.RoutingNumber,
		IssuerBankName:        bn,
		IssuerAccountName:     accountName,
		ExternalLinkId:        elID,
		ExternalLinkRequestId: elrID,
		ExternalLinkName:      eln,
		HolderName:            la.AccountHolderName,
		Alias:                 alias,
		Currency:              string(banking.CurrencyUSD),
		UsageType:             ut,
		LinkedAccountType:     lat,
		LinkedAccountSubtype:  last,
	}

	lbap, err := blas.lac.Register(context.Background(), rr)
	if err != nil {
		return lba, err
	}

	return transformProtoLinkedAccountToLinkedBankAccount(lbap)
}

func (blas bankingLinkedAccountService) GetAccountByTransferID(bID shared.BusinessID, ptID string) (*LinkedBankAccount, error) {
	var lba *LinkedBankAccount

	defer blas.conn.CloseAndCancel()

	s := ptID
	if !strings.HasPrefix(s, id.IDPrefixBankTransfer.String()) {
		s = fmt.Sprintf("%s%s", id.IDPrefixBankTransfer, s)
	}

	tID, err := id.ParseBankTransferID(s)
	if err != nil {
		return lba, err
	}

	r := &grpcBankLinkedAccount.GetSourceAccountByTransferIDRequest{
		TransferId: tID.String(),
	}

	lbap, err := blas.lac.GetSourceAccountByTransferID(context.Background(), r)
	if err != nil {
		return lba, err
	}

	return transformProtoLinkedAccountToLinkedBankAccount(lbap)
}

func (blas bankingLinkedAccountService) GetByAccountIDInternal(aID string) (*LinkedBankAccount, error) {
	var lba *LinkedBankAccount

	defer blas.conn.CloseAndCancel()

	s := aID
	if !strings.HasPrefix(aID, id.IDPrefixBankAccount.String()) {
		s = fmt.Sprintf("%s%s", id.IDPrefixBankAccount, aID)
	}

	baID, err := id.ParseBankAccountID(s)
	if err != nil {
		return lba, err
	}

	gr := &grpcBankLinkedAccount.GetRequest{
		AccountId: baID.String(),
	}

	lbap, err := blas.lac.Get(context.Background(), gr)
	if err != nil {
		return lba, err
	}

	return transformProtoLinkedAccountToLinkedBankAccount(lbap)
}

func (blas bankingLinkedAccountService) GetByAccountNumber(bID shared.BusinessID, aN, rN string) (*LinkedBankAccount, error) {
	lba := new(LinkedBankAccount)

	defer blas.conn.CloseAndCancel()

	gr := &grpcBankLinkedAccount.GetRequest{
		BusinessId:          bID.ToPrefixString(),
		IssuerAccountNumber: aN,
		IssuerRoutingNumber: rN,
	}

	lbap, err := blas.lac.Get(context.Background(), gr)
	//TODO remove this in favor of new error proto
	if err != nil {
		if strings.Contains(err.Error(), "Cannot find linked account to Get") {
			return lba, sql.ErrNoRows
		} else {
			return lba, err
		}
	}

	return transformProtoLinkedAccountToLinkedBankAccount(lbap)
}

func (blas bankingLinkedAccountService) ListWithContact(bID shared.BusinessID, cID shared.ContactID, stfs []grpcBanking.LinkedSubtype, limit, offset int) ([]LinkedBankAccount, error) {
	var lbas []LinkedBankAccount

	defer blas.conn.CloseAndCancel()

	s := string(cID)
	if !strings.HasPrefix(s, id.IDPrefixContact.String()) {
		s = fmt.Sprintf("%s%s", id.IDPrefixContact, s)
	}

	conID, err := id.ParseContactID(s)
	if err != nil {
		return lbas, err
	}

	sf := []grpcBankLinkedAccount.LinkedAccountStatus{
		grpcBankLinkedAccount.LinkedAccountStatus_LAS_ACTIVE,
	}

	gmr := &grpcBankLinkedAccount.GetManyRequest{
		BusinessId:    bID.ToPrefixString(),
		ContactId:     conID.String(),
		StatusFilter:  sf,
		SubtypeFilter: stfs,
		Limit:         int32(limit),
		Offset:        int32(offset),
	}

	lbaps, err := blas.lac.GetMany(context.Background(), gmr)
	if err != nil {
		return lbas, err
	}

	for _, lbap := range lbaps.LinkedAccounts {
		lba, err := transformProtoLinkedAccountToLinkedBankAccount(lbap)
		if err != nil {
			return lbas, err
		}

		lbas = append(lbas, *lba)
	}

	return lbas, nil
}

func (blas bankingLinkedAccountService) List(bID shared.BusinessID, stfs []grpcBanking.LinkedSubtype, limit, offset int) ([]LinkedBankAccount, error) {
	var lbas []LinkedBankAccount

	defer blas.conn.CloseAndCancel()

	sf := []grpcBankLinkedAccount.LinkedAccountStatus{
		grpcBankLinkedAccount.LinkedAccountStatus_LAS_ACTIVE,
	}

	gmr := &grpcBankLinkedAccount.GetManyRequest{
		BusinessId:    bID.ToPrefixString(),
		StatusFilter:  sf,
		SubtypeFilter: stfs,
		Limit:         int32(limit),
		Offset:        int32(offset),
	}

	lbaps, err := blas.lac.GetMany(context.Background(), gmr)
	if err != nil {
		return lbas, err
	}

	for _, lbap := range lbaps.LinkedAccounts {
		lba, err := transformProtoLinkedAccountToLinkedBankAccount(lbap)
		if err != nil {
			return lbas, err
		}

		lbas = append(lbas, *lba)
	}

	return lbas, nil
}

func (blas bankingLinkedAccountService) GetById(slaID string) (*LinkedBankAccount, error) {
	var lba *LinkedBankAccount

	defer blas.conn.CloseAndCancel()

	s := slaID
	if !strings.HasPrefix(slaID, id.IDPrefixLinkedBankAccount.String()) {
		s = fmt.Sprintf("%s%s", id.IDPrefixLinkedBankAccount, slaID)
	}

	laID, err := id.ParseLinkedBankAccountID(s)
	if err != nil {
		return lba, err
	}

	gr := &grpcBankLinkedAccount.GetRequest{
		LinkedAccountId: laID.String(),
	}

	lbap, err := blas.lac.Get(context.Background(), gr)
	if err != nil {
		//TODO change this once we return error types
		//TODO this is a hack because of the ids we pass around are ambiguous
		if strings.Contains(err.Error(), "Cannot find linked account to Get") {
			return blas.GetByAccountIDInternal(slaID)
		}
		return lba, err
	}

	return transformProtoLinkedAccountToLinkedBankAccount(lbap)
}

func (blas bankingLinkedAccountService) UnlinkBankAccount(slaID string) (*LinkedBankAccount, error) {
	var lba *LinkedBankAccount

	defer blas.conn.CloseAndCancel()

	s := slaID
	if !strings.HasPrefix(slaID, id.IDPrefixLinkedBankAccount.String()) {
		s = fmt.Sprintf("%s%s", id.IDPrefixLinkedBankAccount, slaID)
	}

	laID, err := id.ParseLinkedBankAccountID(s)
	if err != nil {
		return lba, err
	}

	dr := &grpcBankLinkedAccount.DeleteRequest{
		LinkedAccountId: laID.String(),
		ActorType:       grpcRoot.ActorType_AT_LEGACY_CORE,
	}

	_, err = blas.lac.Delete(context.Background(), dr)
	if err != nil {
		return lba, err
	}

	gr := &grpcBankLinkedAccount.GetRequest{
		LinkedAccountId: laID.String(),
	}

	lbap, err := blas.lac.Get(context.Background(), gr)
	if err != nil {
		return lba, err
	}

	return transformProtoLinkedAccountToLinkedBankAccount(lbap)
}

func (blas bankingLinkedAccountService) Update(ur *LinkedAccountUpdate) (*LinkedBankAccount, error) {
	var lba *LinkedBankAccount

	defer blas.conn.CloseAndCancel()

	if ur.UsageType == nil {
		return lba, nil
	}

	laID, err := id.ParseLinkedBankAccountID(fmt.Sprintf("%s%s", id.IDPrefixLinkedBankAccount, ur.ID))
	if err != nil {
		return lba, err
	}

	last := getProtoSubTypeFromUsageType(ur.UsageType)

	pur := &grpcBankLinkedAccount.UpdateRequest{
		LinkedAccountId:      laID.String(),
		LinkedAccountSubtype: last,
	}

	_, err = blas.lac.Update(context.Background(), pur)
	if err != nil {
		return lba, err
	}

	gr := &grpcBankLinkedAccount.GetRequest{
		LinkedAccountId: laID.String(),
	}

	lbap, err := blas.lac.Get(context.Background(), gr)
	if err != nil {
		return lba, err
	}

	return transformProtoLinkedAccountToLinkedBankAccount(lbap)
}

func transformProtoLinkedAccountToLinkedBankAccount(lbap *grpcBankLinkedAccount.LinkedAccount) (*LinkedBankAccount, error) {
	lba := new(LinkedBankAccount)

	//Let's check if it even exists?
	if lbap.Id == "" {
		return lba, services.ErrorNotFound{}.New("")
	}

	bID, err := shared.ParseBusinessID(lbap.BusinessId)
	if err != nil {
		return lba, err
	}

	created, err := ptypes.Timestamp(lbap.Created)
	if err != nil {
		return lba, err
	}

	modified, err := ptypes.Timestamp(lbap.Modified)
	if err != nil {
		return lba, err
	}

	var deactivated *time.Time
	if lbap.LinkedAccountStatus == grpcBankLinkedAccount.LinkedAccountStatus_LAS_INACTIVE {
		deactivated = &modified
	}

	sourceName := ""
	if lbap.ExternalLinkId != "" {
		sourceName = "plaid"
	}

	var bbaID *string
	var conID *string

	aID, err := id.ParseBankAccountID(lbap.AccountId)
	if err != nil {
		return lba, err
	}

	if !aID.IsZero() {
		s := aID.String()
		bbaID = &s
	}

	ctID, err := id.ParseContactID(lbap.ContactId)
	if err != nil {
		return lba, err
	}

	if !ctID.IsZero() {
		s := ctID.UUIDString()
		conID = &s
	}

	var sourceAccountID *string
	if lbap.ExternalLinkId != "" {
		sourceAccountID = &lbap.ExternalLinkId
	}

	var sourceID *string
	if lbap.ExternalLinkRequestId != "" {
		sourceID = &lbap.ExternalLinkRequestId
	}

	var accountName *string
	if lbap.IssuerAccountName != "" {
		accountName = &lbap.IssuerAccountName
	}

	perm := getPerm(lbap.UsageType)

	at := getAccountType(lbap.LinkedAccountType)

	ut := getUsageType(lbap.LinkedAccountSubtype)

	return &LinkedBankAccount{
		Id:                    lbap.Id,
		BusinessID:            bID,
		BusinessBankAccountId: bbaID,
		ContactId:             conID,
		RegisteredAccountId:   lbap.PartnerReferenceId,
		RegisteredBankName:    "bbva",
		AccountHolderName:     lbap.HolderName,
		AccountNumber:         AccountNumber(lbap.IssuerAccountNumber),
		AccountName:           accountName,
		BankName:              &lbap.IssuerBankName,
		Currency:              banking.CurrencyUSD,
		AccountType:           at,
		UsageType:             ut,
		RoutingNumber:         lbap.IssuerRoutingNumber,
		SourceId:              sourceID,
		SourceAccountId:       sourceAccountID,
		SourceName:            &sourceName,
		Permission:            perm,
		Alias:                 &lbap.Alias,
		Deactivated:           deactivated,
		Created:               created,
		Modified:              modified,
	}, nil
}

func getPerm(ut grpcBanking.UsageType) banking.LinkedAccountPermission {
	var perm banking.LinkedAccountPermission

	switch ut {
	case grpcBanking.UsageType_UT_RECEIVE_ONLY:
		perm = banking.LinkedAccountPermissionRecieveOnly
	case grpcBanking.UsageType_UT_SEND_AND_RECEIVE:
		perm = banking.LinkedAccountPermissionSendAndRecieve
	}

	return perm
}

func getAccountType(lat grpcBankLinkedAccount.LinkedAccountType) banking.AccountType {
	var at banking.AccountType

	switch lat {
	case grpcBankLinkedAccount.LinkedAccountType_LAT_BUSINESS_SAVINGS:
		at = banking.AccountTypeSavings
	case grpcBankLinkedAccount.LinkedAccountType_LAT_BUSINESS_CHECKING:
		at = banking.AccountTypeChecking
	}

	return at
}

func getUsageType(lst grpcBanking.LinkedSubtype) *UsageType {
	ut := UsageTypeExternal

	switch lst {
	case grpcBanking.LinkedSubtype_LST_CLEARING:
		ut = UsageTypeClearing
	case grpcBanking.LinkedSubtype_LST_PRIMARY:
		ut = UsageTypePrimary
	case grpcBanking.LinkedSubtype_LST_MERCHANT:
		ut = UsageTypeMerchant
	case grpcBanking.LinkedSubtype_LST_CONTACT:
		ut = UsageTypeContact
	case grpcBanking.LinkedSubtype_LST_EXTERNAL:
		ut = UsageTypeExternal
	}

	return &ut
}

func (blas bankingLinkedAccountService) UpdateLinkedAccountUsageType(slaID string, usageType *UsageType) (*LinkedBankAccount, error) {
	var lba *LinkedBankAccount

	defer blas.conn.CloseAndCancel()

	last := getProtoSubTypeFromUsageType(usageType)

	gr := &grpcBankLinkedAccount.UpdateRequest{
		LinkedAccountId:      slaID,
		LinkedAccountSubtype: last,
	}

	lbap, err := blas.lac.Update(context.Background(), gr)
	if err != nil {
		return lba, err
	}

	return transformProtoLinkedAccountToLinkedBankAccount(lbap)
}
