package business

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/wiseco/core-platform/partner/bank/data"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking"
	"github.com/wiseco/core-platform/services/business"
	"github.com/wiseco/core-platform/services/user"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/go-lib/grpc"
	"github.com/wiseco/go-lib/id"
	"github.com/wiseco/go-lib/num"
	grpcRoot "github.com/wiseco/protobuf/golang"
	grpcBanking "github.com/wiseco/protobuf/golang/banking"
	grpcBankAccount "github.com/wiseco/protobuf/golang/banking/account"
)

type bankingAccountService struct {
	conn grpc.Client
	ac   grpcBankAccount.AccountServiceClient
}

func NewBankingAccountService() (*bankingAccountService, error) {
	var bas *bankingAccountService

	bsn, err := grpc.GetConnectionStringForService(grpc.ServiceNameBanking)
	if err != nil {
		return bas, err
	}

	conn, err := grpc.NewInsecureClient(bsn)
	if err != nil {
		return bas, err
	}

	ac := grpcBankAccount.NewAccountServiceClient(conn.GetConn())

	return &bankingAccountService{
		conn: conn,
		ac:   ac,
	}, nil
}

func (bas bankingAccountService) DeactivateAccountBlock(blID string) error {
	defer bas.conn.CloseAndCancel()

	s := blID
	if !strings.HasPrefix(blID, id.IDPrefixBankAccountBlock.String()) {
		s = fmt.Sprintf("%s%s", id.IDPrefixBankAccountBlock, blID)
	}

	pblID, err := id.ParseBankAccountBlockID(s)
	if err != nil {
		return err
	}

	cbr := &grpcBankAccount.ClearBlockRequest{
		BlockId:   pblID.String(),
		ActorType: grpcRoot.ActorType_AT_LEGACY_CORE,
	}

	_, err = bas.ac.ClearBlock(context.Background(), cbr)

	return err
}

func (bas bankingAccountService) GetByBlockID(blID string) (*banking.AccountBlock, error) {
	var ab *banking.AccountBlock

	defer bas.conn.CloseAndCancel()

	pblID, err := id.ParseBankAccountBlockID(fmt.Sprintf("%s%s", id.IDPrefixBankAccountBlock, blID))
	if err != nil {
		return ab, err
	}

	gbr := &grpcBankAccount.GetBlockRequest{
		PartnerBlockId: pblID.String(),
	}

	bp, err := bas.ac.GetBlock(context.Background(), gbr)
	if err != nil {
		return ab, err
	}

	return transformProtoAccountBlockToAccountBlock(bp)
}

func (bas bankingAccountService) GetManyAccountBlocks(aID string) ([]banking.AccountBlock, error) {
	var abs []banking.AccountBlock

	defer bas.conn.CloseAndCancel()

	s := aID
	if !strings.HasPrefix(aID, id.IDPrefixBankAccount.String()) {
		s = fmt.Sprintf("%s%s", id.IDPrefixBankAccount, aID)
	}

	paID, err := id.ParseBankAccountID(s)
	if err != nil {
		return abs, err
	}

	gbr := &grpcBankAccount.GetManyBlockRequest{
		AccountId: paID.String(),
	}

	bps, err := bas.ac.GetManyBlocks(context.Background(), gbr)
	if err != nil {
		return abs, err
	}

	for _, bp := range bps.Blocks {
		ab, err := transformProtoAccountBlockToAccountBlock(bp)
		if err != nil {
			return abs, err
		}

		abs = append(abs, *ab)
	}

	return abs, nil
}

func (bas bankingAccountService) CreateAccountBlock(abc banking.AccountBlockCreate) (*banking.AccountBlock, error) {
	var ab *banking.AccountBlock

	defer bas.conn.CloseAndCancel()

	blockID := ""
	if abc.BlockID != "" {
		blID, err := id.ParseBankAccountBlockID(fmt.Sprintf("%s%s", id.IDPrefixBankAccountBlock, abc.BlockID))
		if err != nil {
			return ab, err
		}
		blockID = blID.String()
	}

	s := abc.AccountID
	if !strings.HasPrefix(abc.AccountID, id.IDPrefixBankAccount.String()) {
		s = fmt.Sprintf("%s%s", id.IDPrefixBankAccount, abc.AccountID)
	}

	baID, err := id.ParseBankAccountID(s)
	if err != nil {
		return ab, err
	}

	of := getProtoOriginatedFromString(abc.OriginatedFrom)

	bt := getProtoBlockTypeFromAccountBlockType(abc.BlockType)

	cbr := &grpcBankAccount.CreateBlockRequest{
		AccountId:      baID.String(),
		PartnerBlockId: blockID,
		BlockType:      bt,
		Reason:         abc.Reason,
		OriginatedFrom: of,
	}

	bp, err := bas.ac.CreateBlock(context.Background(), cbr)
	if err != nil {
		return ab, err
	}

	return transformProtoAccountBlockToAccountBlock(bp)
}

func (bas bankingAccountService) Create(bac *BankAccountCreate, bus business.Business, usr user.User, bankBus *data.Business, bankCon *data.Consumer) (*BankAccount, error) {
	var ba *BankAccount

	defer bas.conn.CloseAndCancel()

	pat := grpcBanking.PartnerAccountType_PAT_DEPOSITORY_CHECKING

	pbt := getPartnerBusinessTypeFromEntityType(*bus.EntityType)

	prs := &grpcBankAccount.ParticipantRequests{
		ParticipantRequests: []*grpcBankAccount.ParticipantRequest{
			&grpcBankAccount.ParticipantRequest{
				ConsumerId:        usr.ConsumerID.ToPrefixString(),
				PartnerConsumerId: string(bankCon.BankID),
				BusinessId:        bus.ID.ToPrefixString(),
				PartnerBusinessId: string(bankBus.BankID),
				Role:              grpcBankAccount.Role_R_HOLDER,
			},
		},
	}

	alias := ""
	if bac.Alias != nil {
		alias = *bac.Alias
	}

	cr := &grpcBankAccount.CreateRequest{
		PartnerName:         grpcBanking.PartnerName_PN_BBVA,
		PartnerAccountType:  pat,
		PartnerBusinessType: pbt,
		Tier:                grpcBankAccount.Tier_T_NORMAL,
		Alias:               alias,
		PrimaryBusinessId:   bus.ID.ToPrefixString(),
		PrimaryConsumerId:   usr.ConsumerID.ToPrefixString(),
		PrimaryHolderName:   bus.Name(),
		AccountType:         grpcBanking.AccountType_AT_BUSINESS,
		AccountSubtype:      grpcBanking.AccountSubtype_AST_PRIMARY,
		ParticipantRequests: prs,
	}

	ap, err := bas.ac.Create(context.Background(), cr)
	if err != nil {
		return ba, err
	}

	return transformProtoAccountToBankAccount(ap)
}

func (bas bankingAccountService) GetByID(aID string) (*BankAccount, error) {
	var ba *BankAccount

	defer bas.conn.CloseAndCancel()

	s := aID
	if !strings.HasPrefix(aID, id.IDPrefixBankAccount.String()) {
		s = fmt.Sprintf("%s%s", id.IDPrefixBankAccount, aID)
	}

	baID, err := id.ParseBankAccountID(s)
	if err != nil {
		return ba, err
	}

	gr := &grpcBankAccount.GetRequest{
		AccountId: baID.String(),
	}

	ap, err := bas.ac.Get(context.Background(), gr)
	if err != nil {
		return ba, err
	}

	return transformProtoAccountToBankAccount(ap)
}

func (bas bankingAccountService) GetByBankAccountId(paID string) (*BankAccount, error) {
	var ba *BankAccount

	defer bas.conn.CloseAndCancel()

	gr := &grpcBankAccount.GetRequest{
		PartnerAccountId: paID,
	}

	ap, err := bas.ac.Get(context.Background(), gr)
	if err != nil {
		return ba, err
	}

	return transformProtoAccountToBankAccount(ap)
}

func (bas bankingAccountService) GetByBusinessID(bID shared.BusinessID, limit, offset int) ([]BankAccount, error) {
	var ba []BankAccount

	defer bas.conn.CloseAndCancel()

	sf := []grpcBanking.AccountStatus{
		grpcBanking.AccountStatus_AS_ACTIVE,
		grpcBanking.AccountStatus_AS_BLOCKED,
	}

	gr := &grpcBankAccount.GetManyRequest{
		BusinessId:    bID.ToPrefixString(),
		Limit:         int32(limit),
		Offset:        int32(offset),
		SubtypeFilter: []grpcBanking.AccountSubtype{grpcBanking.AccountSubtype_AST_PRIMARY},
		StatusFilter:  sf,
	}

	aps, err := bas.ac.GetMany(context.Background(), gr)
	if err != nil {
		return ba, err
	}

	for _, a := range aps.Accounts {
		b, err := transformProtoAccountToBankAccount(a)
		if err != nil {
			return ba, err
		}

		ba = append(ba, *b)
	}

	return ba, nil
}

func (bas bankingAccountService) GetAllPrimaryByBusinessID(bID shared.BusinessID, limit, offset int) ([]BankAccount, error) {
	var ba []BankAccount

	defer bas.conn.CloseAndCancel()

	gr := &grpcBankAccount.GetManyRequest{
		BusinessId:    bID.ToPrefixString(),
		Limit:         int32(limit),
		Offset:        int32(offset),
		SubtypeFilter: []grpcBanking.AccountSubtype{grpcBanking.AccountSubtype_AST_PRIMARY},
	}

	aps, err := bas.ac.GetMany(context.Background(), gr)
	if err != nil {
		return ba, err
	}

	for _, a := range aps.Accounts {
		b, err := transformProtoAccountToBankAccount(a)
		if err != nil {
			return ba, err
		}

		ba = append(ba, *b)
	}

	return ba, nil
}
func (bas bankingAccountService) GetByUsageType(bID shared.BusinessID, ut UsageType, limit, offset int) ([]BankAccount, error) {
	var ba []BankAccount

	defer bas.conn.CloseAndCancel()

	stf := grpcBanking.AccountSubtype_AST_PRIMARY
	if ut == UsageTypeClearing {
		stf = grpcBanking.AccountSubtype_AST_CLEARING
	}

	sf := []grpcBanking.AccountStatus{
		grpcBanking.AccountStatus_AS_ACTIVE,
		grpcBanking.AccountStatus_AS_BLOCKED,
	}

	gr := &grpcBankAccount.GetManyRequest{
		BusinessId:    bID.ToPrefixString(),
		SubtypeFilter: []grpcBanking.AccountSubtype{stf},
		StatusFilter:  sf,
		Limit:         int32(limit),
		Offset:        int32(offset),
	}

	aps, err := bas.ac.GetMany(context.Background(), gr)
	if err != nil {
		return ba, err
	}

	for _, a := range aps.Accounts {
		b, err := transformProtoAccountToBankAccount(a)
		if err != nil {
			return ba, err
		}

		ba = append(ba, *b)
	}

	return ba, nil
}

func (bas bankingAccountService) GetByConsumerIDAndUsageType(cID shared.ConsumerID, ut UsageType, limit, offset int) ([]BankAccount, error) {
	var ba []BankAccount

	defer bas.conn.CloseAndCancel()

	stf := grpcBanking.AccountSubtype_AST_PRIMARY
	if ut == UsageTypeClearing {
		stf = grpcBanking.AccountSubtype_AST_CLEARING
	}

	sf := []grpcBanking.AccountStatus{
		grpcBanking.AccountStatus_AS_ACTIVE,
		grpcBanking.AccountStatus_AS_BLOCKED,
	}

	gr := &grpcBankAccount.GetManyRequest{
		ConsumerId:    cID.ToPrefixString(),
		SubtypeFilter: []grpcBanking.AccountSubtype{stf},
		StatusFilter:  sf,
		Limit:         int32(limit),
		Offset:        int32(offset),
	}

	aps, err := bas.ac.GetMany(context.Background(), gr)
	if err != nil {
		return ba, err
	}

	for _, a := range aps.Accounts {
		b, err := transformProtoAccountToBankAccount(a)
		if err != nil {
			return ba, err
		}

		ba = append(ba, *b)
	}

	return ba, nil
}

func (bas bankingAccountService) ListInternal(limit, offset int) ([]*BankAccount, error) {
	var ba []*BankAccount

	defer bas.conn.CloseAndCancel()

	sf := []grpcBanking.AccountStatus{
		grpcBanking.AccountStatus_AS_ACTIVE,
		grpcBanking.AccountStatus_AS_BLOCKED,
	}

	gr := &grpcBankAccount.GetManyRequest{
		StatusFilter: sf,
		Limit:        int32(limit),
		Offset:       int32(offset),
	}

	aps, err := bas.ac.GetMany(context.Background(), gr)
	if err != nil {
		return ba, err
	}

	for _, a := range aps.Accounts {
		b, err := transformProtoAccountToBankAccount(a)
		if err != nil {
			return ba, err
		}

		ba = append(ba, b)
	}

	return ba, nil
}

func (bas bankingAccountService) Update(u *BankAccountUpdate) (*BankAccount, error) {
	var ba *BankAccount

	defer bas.conn.CloseAndCancel()

	status := getProtoStatus(u.Status)

	s := u.Id
	if !strings.HasPrefix(u.Id, id.IDPrefixBankAccount.String()) {
		s = fmt.Sprintf("%s%s", id.IDPrefixBankAccount, u.Id)
	}

	baID, err := id.ParseBankAccountID(s)
	if err != nil {
		return ba, err
	}

	ur := &grpcBankAccount.UpdateRequest{
		AccountId:     baID.String(),
		Alias:         u.Alias,
		AccountStatus: status,
	}

	ap, err := bas.ac.Update(context.Background(), ur)
	if err != nil {
		return ba, err
	}

	return transformProtoAccountToBankAccount(ap)
}

func (bas bankingAccountService) GetBalanceByID(aID string, preferCache bool) (*BankAccount, error) {
	var ba *BankAccount

	defer bas.conn.CloseAndCancel()

	s := aID
	if !strings.HasPrefix(aID, id.IDPrefixBankAccount.String()) {
		s = fmt.Sprintf("%s%s", id.IDPrefixBankAccount, aID)
	}

	baID, err := id.ParseBankAccountID(s)
	if err != nil {
		return ba, err
	}

	bur := &grpcBankAccount.BalanceUpdateRequest{
		AccountId:   baID.String(),
		PreferCache: preferCache,
	}

	ap, err := bas.ac.BalanceUpdate(context.Background(), bur)
	if err != nil {
		return ba, err
	}

	return transformProtoAccountToBankAccount(ap)
}

func (bas bankingAccountService) DeactivateAccount(aID string, reason grpcBanking.AccountStatusReason) (*BankAccount, error) {
	var ba *BankAccount

	defer bas.conn.CloseAndCancel()

	s := aID
	if !strings.HasPrefix(aID, id.IDPrefixBankAccount.String()) {
		s = fmt.Sprintf("%s%s", id.IDPrefixBankAccount, aID)
	}

	baID, err := id.ParseBankAccountID(s)
	if err != nil {
		return ba, err
	}

	dr := &grpcBankAccount.DeleteRequest{
		AccountId:           baID.String(),
		ActorType:           grpcRoot.ActorType_AT_LEGACY_CORE,
		AccountStatusReason: reason,
	}

	ap, err := bas.ac.Delete(context.Background(), dr)
	if err != nil {
		return ba, err
	}

	return transformProtoAccountToBankAccount(ap)
}

func transformProtoAccountToBankAccount(ap *grpcBankAccount.Account) (*BankAccount, error) {
	ba := new(BankAccount)

	if ap.Id == "" {
		return ba, services.ErrorNotFound{}.New("")
	}

	as := getAccountStatusFromProto(ap.AccountStatus)

	opened, err := ptypes.Timestamp(ap.Opened)
	if err != nil {
		return ba, err
	}

	created, err := ptypes.Timestamp(ap.Created)
	if err != nil {
		return ba, err
	}

	modified, err := ptypes.Timestamp(ap.Modified)
	if err != nil {
		return ba, err
	}

	if ap.Participants == nil {
		return ba, fmt.Errorf("No participants returned for AccountID:%s", ap.Id)
	}

	ut := UsageTypePrimary
	if ap.AccountSubtype == grpcBanking.AccountSubtype_AST_CLEARING {
		ut = UsageTypeClearing
	}

	pb, err := num.ParseDecimal(ap.PostedBalance)
	if err != nil {
		return ba, err
	}

	pdb, err := num.ParseDecimal(ap.PendingDebitBalance)
	if err != nil {
		return ba, err
	}

	acb, err := num.ParseDecimal(ap.ActualBalance)
	if err != nil {
		return ba, err
	}

	rfb, err := num.ParseDecimal(ap.RemainingFundingBalance)
	if err != nil {
		return ba, err
	}

	fl, err := num.ParseDecimal(ap.FundingLimit)
	if err != nil {
		return ba, err
	}

	pbf, _ := pb.Float64()
	pbdf, _ := pdb.Float64()
	acbf, _ := acb.Float64()
	rfbf, _ := rfb.Float64()
	flf, _ := fl.Float64()

	bID, err := id.ParseBusinessID(ap.Participants.Participants[0].BusinessId)
	if err != nil {
		return ba, err
	}

	ba.Id = ap.Id
	ba.BusinessID = shared.BusinessID(bID.UUIDString())
	ba.UsageType = ut
	ba.BankName = banking.BankNameBBVA
	ba.BankAccountId = ap.PartnerAccountId
	ba.AccountType = "checking"
	ba.AccountStatus = as
	ba.AccountNumber = ap.PartnerAccountNumber
	ba.RoutingNumber = ap.PartnerRoutingNumber
	ba.Alias = &ap.Alias
	ba.AvailableBalance = acbf
	ba.PostedBalance = pbf
	ba.PendingDebitBalance = pbdf
	ba.ActualBalance = acbf
	ba.Currency = banking.CurrencyUSD
	ba.RemainingFundingBalance = rfbf
	ba.FundingLimit = flf
	ba.Opened = opened
	ba.Created = created
	ba.Modified = modified

	return ba, nil
}

func (bas bankingAccountService) FindCount() (int64, error) {
	var count int64

	defer bas.conn.CloseAndCancel()

	ac, err := bas.ac.FindCount(context.Background(), &grpcBankAccount.AccountCountRequest{})
	if err != nil {
		return count, err
	}

	return int64(ac.Count), nil
}

func transformProtoAccountBlockToAccountBlock(bp *grpcBankAccount.Block) (*banking.AccountBlock, error) {
	var ab *banking.AccountBlock

	created, err := ptypes.Timestamp(bp.Created)
	if err != nil {
		return ab, err
	}

	var deactivated *time.Time
	if bp.BlockStatus == grpcBankAccount.BlockStatus_BS_INACTIVE {
		modified, err := ptypes.Timestamp(bp.Created)
		if err != nil {
			return ab, err
		}

		deactivated = &modified
	}

	bt := getBlockTypeFromProto(bp.BlockType)

	return &banking.AccountBlock{
		ID:             bp.Id,
		AccountID:      bp.AccountId,
		BlockID:        bp.PartnerBlockId,
		BlockType:      bt,
		Reason:         bp.Reason,
		Created:        created,
		Deactivated:    deactivated,
		OriginatedFrom: bp.OriginatedFrom.String(),
	}, nil
}

func getBlockTypeFromProto(bt grpcBankAccount.BlockType) banking.AccountBlockType {
	var abt banking.AccountBlockType

	switch bt {
	case grpcBankAccount.BlockType_BT_DEBITS:
		abt = banking.AccountBlockTypeDebit
	case grpcBankAccount.BlockType_BT_CREDITS:
		abt = banking.AccountBlockTypeCredit
	case grpcBankAccount.BlockType_BT_CHECKS:
		abt = banking.AccountBlockTypeCheck
	case grpcBankAccount.BlockType_BT_ALL:
		abt = banking.AccountBlockTypeAll
	}

	return abt
}

func getProtoBlockTypeFromAccountBlockType(abt banking.AccountBlockType) grpcBankAccount.BlockType {
	var bt grpcBankAccount.BlockType

	switch abt {
	case banking.AccountBlockTypeDebit:
		bt = grpcBankAccount.BlockType_BT_DEBITS
	case banking.AccountBlockTypeCredit:
		bt = grpcBankAccount.BlockType_BT_CREDITS
	case banking.AccountBlockTypeCheck:
		bt = grpcBankAccount.BlockType_BT_CHECKS
	case banking.AccountBlockTypeAll:
		bt = grpcBankAccount.BlockType_BT_ALL
	}

	return bt
}

func getProtoOriginatedFromString(s string) grpcBankAccount.OriginatedFromType {
	var of grpcBankAccount.OriginatedFromType

	switch s {
	case "Notification service":
		of = grpcBankAccount.OriginatedFromType_OFT_NOTIFICATIONS
	case "Customer Success Portal":
		of = grpcBankAccount.OriginatedFromType_OFT_AGENT
	default:
		of = grpcBankAccount.OriginatedFromType_OFT_UNSPECIFIED
	}

	return of
}

func getPartnerBusinessTypeFromEntityType(et string) grpcBankAccount.PartnerBusinessType {
	pbt := grpcBankAccount.PartnerBusinessType_PBT_CORPORATE

	switch et {
	case "limitedLiabilityCompany", "singleMemberLLC":
		pbt = grpcBankAccount.PartnerBusinessType_PBT_LLC
	}

	return pbt
}

func getProtoStatus(s string) grpcBanking.AccountStatus {
	var as grpcBanking.AccountStatus

	switch s {
	case banking.BankAccountStatusActive:
		as = grpcBanking.AccountStatus_AS_ACTIVE
	case banking.BankAccountStatusBlocked:
		as = grpcBanking.AccountStatus_AS_BLOCKED
	case banking.BankAccountStatusLocked:
		as = grpcBanking.AccountStatus_AS_LOCKED
	case banking.BankAccountStatusClosePending:
		as = grpcBanking.AccountStatus_AS_CLOSE_PENDING
	case banking.BankAccountStatusClosed:
		as = grpcBanking.AccountStatus_AS_CLOSED
	case banking.BankAccountStatusDormant:
		as = grpcBanking.AccountStatus_AS_DORMANT
	case banking.BankAccountStatusAbandoned:
		as = grpcBanking.AccountStatus_AS_ABANDONED
	case banking.BankAccountStatusEscheated:
		as = grpcBanking.AccountStatus_AS_ESCHEATED
	case banking.BankAccountStatusChargeOff:
		as = grpcBanking.AccountStatus_AS_CHARGE_OFF
	default:
		as = grpcBanking.AccountStatus_AS_UNSPECIFIED
	}

	return as
}

func getAccountStatusFromProto(as grpcBanking.AccountStatus) string {
	var s string

	switch as {
	case grpcBanking.AccountStatus_AS_ACTIVE:
		s = banking.BankAccountStatusActive
	case grpcBanking.AccountStatus_AS_BLOCKED:
		s = banking.BankAccountStatusBlocked
	case grpcBanking.AccountStatus_AS_LOCKED:
		s = banking.BankAccountStatusLocked
	case grpcBanking.AccountStatus_AS_CLOSE_PENDING:
		s = banking.BankAccountStatusClosePending
	case grpcBanking.AccountStatus_AS_CLOSED:
		s = banking.BankAccountStatusClosed
	case grpcBanking.AccountStatus_AS_DORMANT:
		s = banking.BankAccountStatusDormant
	case grpcBanking.AccountStatus_AS_ABANDONED:
		s = banking.BankAccountStatusAbandoned
	case grpcBanking.AccountStatus_AS_ESCHEATED:
		s = banking.BankAccountStatusEscheated
	case grpcBanking.AccountStatus_AS_CHARGE_OFF:
		s = banking.BankAccountStatusChargeOff
	}

	return s
}
