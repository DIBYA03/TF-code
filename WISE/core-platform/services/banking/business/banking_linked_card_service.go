package business

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/xtgo/uuid"

	"github.com/golang/protobuf/ptypes"
	"github.com/wiseco/core-platform/partner/bank/bbva"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/banking"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/go-lib/grpc"
	"github.com/wiseco/go-lib/id"
	grpcRoot "github.com/wiseco/protobuf/golang"
	grpcBanking "github.com/wiseco/protobuf/golang/banking"
	grpcBankLinkedCard "github.com/wiseco/protobuf/golang/banking/linked_card"
)

type bankingLinkedCardService struct {
	conn grpc.Client
	lcc  grpcBankLinkedCard.LinkedCardServiceClient
}

func NewBankingLinkedCardService() (*bankingLinkedCardService, error) {
	var lcs *bankingLinkedCardService

	bsn, err := grpc.GetConnectionStringForService(grpc.ServiceNameBanking)
	if err != nil {
		return lcs, err
	}

	conn, err := grpc.NewInsecureClient(bsn)
	if err != nil {
		return lcs, err
	}

	lcc := grpcBankLinkedCard.NewLinkedCardServiceClient(conn.GetConn())

	return &bankingLinkedCardService{
		conn: conn,
		lcc:  lcc,
	}, nil
}

func (lcs bankingLinkedCardService) Create(cc *LinkedCardCreate, cID shared.ConsumerID) (*LinkedCard, error) {
	var lc *LinkedCard
	var addy *grpcRoot.Address

	defer lcs.conn.CloseAndCancel()

	conID := ""
	if cc.ContactId != nil {
		cUUID, err := uuid.Parse(*cc.ContactId)
		if err != nil {
			return lc, err
		}

		conID = id.ContactID(cUUID).String()
	}

	if cc.BillingAddress != nil {
		addy = &grpcRoot.Address{
			Line_1:     cc.BillingAddress.StreetAddress,
			Line_2:     cc.BillingAddress.AddressLine2,
			Locality:   cc.BillingAddress.City,
			AdminArea:  cc.BillingAddress.State,
			Country:    cc.BillingAddress.Country,
			PostalCode: cc.BillingAddress.PostalCode,
		}
	}

	st := getProtoSubTypeFromUsageType(cc.UsageType)
	ut := getProtoUsageTypeFromPermission(cc.Permission)

	rr := &grpcBankLinkedCard.RegisterRequest{
		ConsumerId:          cID.ToPrefixString(),
		BusinessId:          cc.BusinessID.ToPrefixString(),
		ContactId:           conID,
		PartnerName:         grpcBanking.PartnerName_PN_BBVA,
		Address:             addy,
		Cvv2CvcCode:         cc.CVVCode,
		IssuerAccountNumber: string(cc.CardNumber),
		ExpirationDate:      cc.ExpirationDate.String(),
		HolderName:          cc.CardHolderName,
		Alias:               cc.Alias,
		UsageType:           ut,
		LinkedCardSubtype:   st,
	}

	if cc.ValidateCard {
		rr.ValidateCard = grpcRoot.Boolean_B_TRUE
	} else {
		rr.ValidateCard = grpcRoot.Boolean_B_FALSE
	}

	lcp, err := lcs.lcc.RegisterAndValidate(context.Background(), rr)
	if err != nil {
		return lc, err
	}

	return transformProtoLinkedCardToLinkedCard(lcp)
}

func (lcs bankingLinkedCardService) Delete(lcID string) error {
	defer lcs.conn.CloseAndCancel()

	s := lcID
	if !strings.HasPrefix(lcID, id.IDPrefixLinkedCard.String()) {
		s = fmt.Sprintf("%s%s", id.IDPrefixLinkedCard, lcID)
	}

	plcID, err := id.ParseLinkedCardID(s)
	if err != nil {
		return err
	}

	dr := &grpcBankLinkedCard.DeleteRequest{
		LinkedCardId: plcID.String(),
	}

	_, err = lcs.lcc.Delete(context.Background(), dr)

	return err
}

func (lcs bankingLinkedCardService) RegisterExistingCard(businessID shared.BusinessID, contactID string, hash string) (*LinkedCard, error) {
	var lc *LinkedCard

	defer lcs.conn.CloseAndCancel()

	s := contactID
	if !strings.HasPrefix(contactID, id.IDPrefixContact.String()) {
		s = fmt.Sprintf("%s%s", id.IDPrefixContact, contactID)
	}

	conID, err := id.ParseContactID(s)
	if err != nil {
		return lc, err
	}

	req := &grpcBankLinkedCard.RegisterExistingCardToNewContactRequest{
		BusinessId:     businessID.ToPrefixString(),
		ContactId:      conID.String(),
		CardInfoHashed: hash,
	}

	lcp, err := lcs.lcc.RegisterExistingCardToNewContact(context.Background(), req)
	if err != nil {
		return lc, err
	}

	return transformProtoLinkedCardToLinkedCard(lcp)
}

func (lcs bankingLinkedCardService) GetByID(lcID string) (*LinkedCard, error) {
	var lc *LinkedCard

	defer lcs.conn.CloseAndCancel()

	s := lcID
	if !strings.HasPrefix(lcID, id.IDPrefixLinkedCard.String()) {
		s = fmt.Sprintf("%s%s", id.IDPrefixLinkedCard, lcID)
	}

	plcID, err := id.ParseLinkedCardID(s)
	if err != nil {
		return lc, err
	}

	gr := &grpcBankLinkedCard.GetRequest{
		LinkedCardId: plcID.String(),
	}

	lcp, err := lcs.lcc.Get(context.Background(), gr)

	//TODO remove this in favor of new error proto
	if err != nil {
		if strings.Contains(err.Error(), "Cannot find linked card to Get") {
			return lc, sql.ErrNoRows
		} else {
			return lc, err
		}
	}

	return transformProtoLinkedCardToLinkedCard(lcp)
}

func (lcs bankingLinkedCardService) ListWithContact(bID shared.BusinessID, cID string, stfs []grpcBanking.LinkedSubtype, limit, offset int) ([]*LinkedCard, error) {
	var ls []*LinkedCard

	defer lcs.conn.CloseAndCancel()

	s := cID
	if !strings.HasPrefix(cID, id.IDPrefixContact.String()) {
		s = fmt.Sprintf("%s%s", id.IDPrefixContact, cID)
	}

	conID, err := id.ParseContactID(s)
	if err != nil {
		return ls, err
	}

	sf := []grpcBankLinkedCard.LinkedCardStatus{
		grpcBankLinkedCard.LinkedCardStatus_LCS_ACTIVE,
	}

	gmr := &grpcBankLinkedCard.GetManyRequest{
		BusinessId:    bID.ToPrefixString(),
		ContactId:     conID.String(),
		StatusFilter:  sf,
		SubtypeFilter: stfs,
		Limit:         int32(limit),
		Offset:        int32(offset),
	}

	lcps, err := lcs.lcc.GetMany(context.Background(), gmr)
	if err != nil {
		return ls, err
	}

	for _, lcp := range lcps.LinkedCards {
		lc, err := transformProtoLinkedCardToLinkedCard(lcp)
		if err != nil {
			return ls, err
		}

		ls = append(ls, lc)
	}

	return ls, nil
}

func (lcs bankingLinkedCardService) List(bID shared.BusinessID, limit, offset int) ([]*LinkedCard, error) {
	var ls []*LinkedCard

	defer lcs.conn.CloseAndCancel()

	stfs := []grpcBanking.LinkedSubtype{
		grpcBanking.LinkedSubtype_LST_EXTERNAL,
	}

	sf := []grpcBankLinkedCard.LinkedCardStatus{
		grpcBankLinkedCard.LinkedCardStatus_LCS_ACTIVE,
	}

	gmr := &grpcBankLinkedCard.GetManyRequest{
		BusinessId:    bID.ToPrefixString(),
		StatusFilter:  sf,
		SubtypeFilter: stfs,
		Limit:         int32(limit),
		Offset:        int32(offset),
	}

	lcps, err := lcs.lcc.GetMany(context.Background(), gmr)
	if err != nil {
		return ls, err
	}

	for _, lcp := range lcps.LinkedCards {
		lc, err := transformProtoLinkedCardToLinkedCard(lcp)
		if err != nil {
			return ls, err
		}

		ls = append(ls, lc)
	}

	return ls, nil
}

func (lcs bankingLinkedCardService) GetByLinkedCardHashAndContactID(cID string, hash string) (*LinkedCard, error) {
	var lc *LinkedCard

	defer lcs.conn.CloseAndCancel()

	s := cID
	if !strings.HasPrefix(cID, id.IDPrefixContact.String()) {
		s = fmt.Sprintf("%s%s", id.IDPrefixContact, cID)
	}

	conID, err := id.ParseContactID(s)
	if err != nil {
		return lc, err
	}

	gr := &grpcBankLinkedCard.GetRequest{
		CardInfoHashed: hash,
		ContactId:      conID.String(),
	}

	lcp, err := lcs.lcc.Get(context.Background(), gr)

	//TODO remove this in favor of new error proto
	if err != nil {
		if strings.Contains(err.Error(), "Cannot find linked card to Get") {
			return lc, sql.ErrNoRows
		} else {
			return lc, err
		}
	}

	return transformProtoLinkedCardToLinkedCard(lcp)
}

func (lcs bankingLinkedCardService) GetByLinkedCardHash(bID shared.BusinessID, hash string) ([]LinkedCard, error) {
	var ls []LinkedCard

	defer lcs.conn.CloseAndCancel()

	sf := []grpcBankLinkedCard.LinkedCardStatus{
		grpcBankLinkedCard.LinkedCardStatus_LCS_ACTIVE,
	}

	gmr := &grpcBankLinkedCard.GetManyRequest{
		BusinessId:     bID.ToPrefixString(),
		CardInfoHashed: hash,
		StatusFilter:   sf,
		Limit:          int32(10),
	}

	lcps, err := lcs.lcc.GetMany(context.Background(), gmr)
	if err != nil {
		return ls, err
	}

	for _, lcp := range lcps.LinkedCards {
		lc, err := transformProtoLinkedCardToLinkedCard(lcp)
		if err != nil {
			return ls, err
		}

		ls = append(ls, *lc)
	}

	return ls, nil
}

func (lcs bankingLinkedCardService) UpdateLinkedCardUsageType(lcID string, usageType *UsageType) (*LinkedCard, error) {
	var lc *LinkedCard

	defer lcs.conn.CloseAndCancel()

	lcst := getProtoSubTypeFromUsageType(usageType)

	gr := &grpcBankLinkedCard.UpdateRequest{
		LinkedCardId:      lcID,
		LinkedCardSubtype: lcst,
	}

	lbap, err := lcs.lcc.Update(context.Background(), gr)
	if err != nil {
		return lc, err
	}

	return transformProtoLinkedCardToLinkedCard(lbap)
}

func transformProtoLinkedCardToLinkedCard(lcp *grpcBankLinkedCard.LinkedCard) (*LinkedCard, error) {
	var addy *services.Address

	lc := new(LinkedCard)

	if lcp.Id == "" {
		return lc, services.ErrorNotFound{}.New("")
	}

	cnm := CardNumber(lcp.CardNumberLast_4)

	perm := getPermFromProtoUsageType(lcp.UsageType)

	if lcp.Address != nil {
		addy = &services.Address{
			Type:          services.AddressTypeMailing,
			StreetAddress: lcp.Address.Line_1,
			AddressLine2:  lcp.Address.Line_2,
			City:          lcp.Address.Locality,
			State:         lcp.Address.AdminArea,
			Country:       lcp.Address.Country,
			PostalCode:    lcp.Address.PostalCode,
		}
	}

	ut := getUsageTypeFromProtoSubtype(lcp.LinkedCardSubtype)

	created, err := ptypes.Timestamp(lcp.Created)
	if err != nil {
		return lc, err
	}

	modified, err := ptypes.Timestamp(lcp.Modified)
	if err != nil {
		return lc, err
	}

	var deactivated *time.Time
	if lcp.LinkedCardStatus == grpcBankLinkedCard.LinkedCardStatus_LCS_INACTIVE {
		deactivated = &modified
	}

	ffe := false
	if lcp.FastFundsEnabled == "yes" {
		ffe = true
	}

	var conID *string

	if lcp.ContactId != "" {
		id, err := id.ParseContactID(lcp.ContactId)
		if err != nil {
			return lc, err
		}

		if !id.IsZero() {
			s := id.UUIDString()
			conID = &s
		}
	}

	bID, err := id.ParseBusinessID(lcp.BusinessId)
	if err != nil {
		return lc, err
	}

	ct := bbva.GetWiseLinkedCardType(bbva.RegisterCardType(lcp.CardType))

	return &LinkedCard{
		Id:                 lcp.Id,
		BusinessID:         shared.BusinessID(bID.UUIDString()),
		ContactId:          conID,
		RegisteredCardId:   lcp.PartnerReferenceId,
		RegisteredBankName: "bbva",
		CardNumberMasked:   cnm,
		CardBrand:          lcp.CardBrand,
		CardType:           string(ct),
		CardIssuer:         lcp.IssuerBankName,
		FastFundsEnabled:   ffe,
		CardHolderName:     lcp.HolderName,
		Alias:              &lcp.Alias,
		Permission:         perm,
		BillingAddress:     addy,
		Verified:           false,
		UsageType:          ut,
		Deactivated:        deactivated,
		Created:            created,
		Modified:           modified,
	}, nil
}

func getProtoUsageTypeFromPermission(p banking.LinkedAccountPermission) grpcBanking.UsageType {
	var ut grpcBanking.UsageType

	switch p {
	case banking.LinkedAccountPermissionRecieveOnly:
		ut = grpcBanking.UsageType_UT_RECEIVE_ONLY
	case banking.LinkedAccountPermissionSendOnly:
		ut = grpcBanking.UsageType_UT_SEND_ONLY
	case banking.LinkedAccountPermissionSendAndRecieve:
		ut = grpcBanking.UsageType_UT_SEND_AND_RECEIVE
	}

	return ut
}

func getPermFromProtoUsageType(ut grpcBanking.UsageType) banking.LinkedAccountPermission {
	var perm banking.LinkedAccountPermission

	switch ut {
	case grpcBanking.UsageType_UT_RECEIVE_ONLY:
		perm = banking.LinkedAccountPermissionRecieveOnly
	case grpcBanking.UsageType_UT_SEND_ONLY:
		perm = banking.LinkedAccountPermissionSendOnly
	case grpcBanking.UsageType_UT_SEND_AND_RECEIVE:
		perm = banking.LinkedAccountPermissionSendAndRecieve
	}

	return perm
}

func getUsageTypeFromProtoSubtype(st grpcBanking.LinkedSubtype) *UsageType {
	var ut UsageType

	switch st {
	case grpcBanking.LinkedSubtype_LST_CONTACT:
		ut = UsageTypeContact
	case grpcBanking.LinkedSubtype_LST_CONTACT_INVISIBLE:
		ut = UsageTypeContactInvisible
	case grpcBanking.LinkedSubtype_LST_EXTERNAL:
		ut = UsageTypeExternal
	}

	return &ut
}

func getProtoSubTypeFromUsageType(u *UsageType) grpcBanking.LinkedSubtype {
	var st grpcBanking.LinkedSubtype

	ut := UsageTypeNone
	if u != nil {
		ut = *u
	}

	switch ut {
	case UsageTypeContact:
		st = grpcBanking.LinkedSubtype_LST_CONTACT
	case UsageTypeContactInvisible:
		st = grpcBanking.LinkedSubtype_LST_CONTACT_INVISIBLE
	default:
		st = grpcBanking.LinkedSubtype_LST_EXTERNAL
	}

	return st
}
