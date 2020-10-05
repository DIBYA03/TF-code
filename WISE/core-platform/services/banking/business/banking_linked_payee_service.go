package business

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/golang/protobuf/ptypes"
	partnerbank "github.com/wiseco/core-platform/partner/bank"
	"github.com/wiseco/core-platform/services/address"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/go-lib/grpc"
	"github.com/wiseco/go-lib/id"
	grpcRoot "github.com/wiseco/protobuf/golang"
	grpcBanking "github.com/wiseco/protobuf/golang/banking"
	grpcBankLinkedPayee "github.com/wiseco/protobuf/golang/banking/linked_payee"
)

type bankingLinkedPayeeService struct {
	conn grpc.Client
	lpc  grpcBankLinkedPayee.LinkedPayeeServiceClient
}

func NewBankingLinkedPayeeService() (*bankingLinkedPayeeService, error) {
	var blps *bankingLinkedPayeeService

	bsn, err := grpc.GetConnectionStringForService(grpc.ServiceNameBanking)
	if err != nil {
		return blps, err
	}

	conn, err := grpc.NewInsecureClient(bsn)
	if err != nil {
		return blps, err
	}

	lac := grpcBankLinkedPayee.NewLinkedPayeeServiceClient(conn.GetConn())

	return &bankingLinkedPayeeService{
		conn: conn,
		lpc:  lac,
	}, nil
}

func (blps bankingLinkedPayeeService) GetByID(pID string) (*LinkedPayee, error) {
	var lp *LinkedPayee

	defer blps.conn.CloseAndCancel()

	s := pID
	if !strings.HasPrefix(pID, id.IDPrefixLinkedPayee.String()) {
		s = fmt.Sprintf("%s%s", id.IDPrefixLinkedPayee, pID)
	}

	lpID, err := id.ParseLinkedPayeeID(s)
	if err != nil {
		return lp, err
	}

	gr := &grpcBankLinkedPayee.GetRequest{
		LinkedPayeeId: lpID.String(),
	}

	lpp, err := blps.lpc.Get(context.Background(), gr)
	//TODO remove this in favor of new error proto
	if err != nil {
		if strings.Contains(err.Error(), "Cannot find linked payee to Get") {
			return lp, sql.ErrNoRows
		} else {
			return lp, err
		}
	}

	return transformProtoLinkedPayeeToLinkedPayee(lpp)
}

func (blps bankingLinkedPayeeService) GetByAddressID(aID shared.AddressID) (*LinkedPayee, error) {
	var lp *LinkedPayee

	defer blps.conn.CloseAndCancel()

	gr := &grpcBankLinkedPayee.GetRequest{
		AddressId: aID.ToPrefixString(),
	}

	lpp, err := blps.lpc.Get(context.Background(), gr)
	//TODO remove this in favor of new error proto
	if err != nil {
		if strings.Contains(err.Error(), "Cannot find linked payee to Get") {
			return lp, sql.ErrNoRows
		} else {
			return lp, err
		}
	}

	return transformProtoLinkedPayeeToLinkedPayee(lpp)
}

func (blps bankingLinkedPayeeService) List(cID string, bID shared.BusinessID, limit, offset int) ([]*LinkedPayee, error) {
	var lps []*LinkedPayee

	defer blps.conn.CloseAndCancel()

	s := string(cID)
	if !strings.HasPrefix(s, id.IDPrefixContact.String()) {
		s = fmt.Sprintf("%s%s", id.IDPrefixContact, s)
	}

	conID, err := id.ParseContactID(s)
	if err != nil {
		return lps, err
	}

	sf := []grpcBankLinkedPayee.LinkedPayeeStatus{
		grpcBankLinkedPayee.LinkedPayeeStatus_LPS_ACTIVE,
	}

	gmr := &grpcBankLinkedPayee.GetManyRequest{
		BusinessId:   bID.ToPrefixString(),
		ContactId:    conID.String(),
		StatusFilter: sf,
		Limit:        int32(limit),
		Offset:       int32(offset),
	}

	lpps, err := blps.lpc.GetMany(context.Background(), gmr)
	if err != nil {
		return lps, err
	}

	for _, lpp := range lpps.LinkedPayees {
		lp, err := transformProtoLinkedPayeeToLinkedPayee(lpp)
		if err != nil {
			return lps, err
		}

		lps = append(lps, lp)
	}

	return lps, nil
}

func (blps bankingLinkedPayeeService) Create(payeeCreate LinkedPayeeCreate, a *address.Address) (*LinkedPayee, error) {
	var lp *LinkedPayee

	defer blps.conn.CloseAndCancel()

	s := string(payeeCreate.ContactID)
	if !strings.HasPrefix(s, id.IDPrefixContact.String()) {
		s = fmt.Sprintf("%s%s", id.IDPrefixContact, s)
	}

	conID, err := id.ParseContactID(s)
	if err != nil {
		return lp, err
	}

	addy := &grpcRoot.Address{
		Line_1:     a.StreetAddress,
		Line_2:     a.AddressLine2,
		Locality:   a.Locality,
		AdminArea:  a.AdminArea,
		Country:    a.Country,
		PostalCode: a.PostalCode,
	}

	rr := &grpcBankLinkedPayee.RegisterRequest{
		BusinessId:  payeeCreate.BusinessID.ToPrefixString(),
		ContactId:   conID.String(),
		Address:     addy,
		PartnerName: grpcBanking.PartnerName_PN_BBVA,
		HolderName:  payeeCreate.AccountHolderName,
		PayeeName:   payeeCreate.PayeeName,
		AddressId:   a.ID.ToPrefixString(),
	}

	lpp, err := blps.lpc.Register(context.Background(), rr)
	if err != nil {
		return lp, err
	}

	return transformProtoLinkedPayeeToLinkedPayee(lpp)
}

func (blps bankingLinkedPayeeService) Deactivate(pID string) error {
	defer blps.conn.CloseAndCancel()

	s := pID
	if !strings.HasPrefix(pID, id.IDPrefixLinkedPayee.String()) {
		s = fmt.Sprintf("%s%s", id.IDPrefixLinkedPayee, pID)
	}

	lpID, err := id.ParseLinkedPayeeID(s)
	if err != nil {
		return err
	}

	dr := &grpcBankLinkedPayee.DeleteRequest{
		LinkedPayeeId: lpID.String(),
	}

	_, err = blps.lpc.Delete(context.Background(), dr)

	return err
}

func transformProtoLinkedPayeeToLinkedPayee(lpp *grpcBankLinkedPayee.LinkedPayee) (*LinkedPayee, error) {
	var lp *LinkedPayee

	var conID shared.ContactID

	bID, err := shared.ParseBusinessID(lpp.BusinessId)
	if err != nil {
		return lp, err
	}

	ctID, err := id.ParseContactID(lpp.ContactId)
	if err != nil {
		return lp, err
	}

	if !ctID.IsZero() {
		s := ctID.UUIDString()
		conID, err = shared.ParseContactID(s)
		if err != nil {
			return lp, err
		}
	}

	addID, err := shared.ParseAddressID(lpp.AddressId)
	if err != nil {
		return lp, err
	}

	created, err := ptypes.Timestamp(lpp.Created)
	if err != nil {
		return lp, err
	}

	modified, err := ptypes.Timestamp(lpp.Modified)
	if err != nil {
		return lp, err
	}

	status := getPayeeStatusFromProto(lpp.LinkedPayeeStatus)

	//TODO figure out address
	return &LinkedPayee{
		ID:                lpp.Id,
		BusinessID:        bID,
		ContactID:         conID,
		AddressID:         addID,
		BankPayeeID:       lpp.PartnerReferenceId,
		BankName:          partnerbank.ProviderNameBBVA,
		AccountHolderName: lpp.HolderName,
		PayeeName:         lpp.PayeeName,
		Status:            status,
		Created:           created,
		Modified:          modified,
	}, nil
}

func getPayeeStatusFromProto(lps grpcBankLinkedPayee.LinkedPayeeStatus) PayeeStatus {
	ps := PayeeStatusActive

	switch lps {
	case grpcBankLinkedPayee.LinkedPayeeStatus_LPS_INACTIVE:
		ps = PayeeStatusInactive
	}

	return ps
}
