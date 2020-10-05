/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package business

import (
	"errors"

	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/wiseco/core-platform/partner/bank/bbva"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/data"
	"github.com/wiseco/core-platform/shared"
	"github.com/wiseco/go-lib/grpc"
	"github.com/wiseco/protobuf/golang/shopping/shopify"
)

type partnerDatastore struct {
	sourceReq services.SourceRequest
	*sqlx.DB
}

type PartnerService interface {
	GetShopifyBusinessByID(shared.BusinessID) (*shopify.ShopifyBusiness, error)
	ActivatePartnerBusiness(*Partner) error
}

func NewPartnerService(r services.SourceRequest) PartnerService {
	return &partnerDatastore{r, data.DBWrite}
}

func (db *partnerDatastore) ActivatePartnerBusiness(p *Partner) error {
	if len(p.Name) == 0 {
		return errors.New("partner name cannot be empty")
	}

	if len(p.ActivationCode) == 0 {
		return errors.New("partner activation cannot be empty")
	}

	switch p.Name {
	case PartnerNameShopify:
		return p.activateShopifyBusiness()
	default:
		return errors.New("invalid partner name")
	}
}

func (p *Partner) activateShopifyBusiness() error {
	sn, err := grpc.GetConnectionStringForService(grpc.ServiceNameShopping)
	if err != nil {
		return err
	}

	client, err := grpc.NewInsecureClient(sn)
	if err != nil {
		return err
	}

	defer client.CloseAndCancel()
	shopifyServiceClient := shopify.NewShopifyBusinessServiceClient(client.GetConn())

	req := shopify.ActivateRequest{
		BusinessId:      p.BusinessID.ToPrefixString(),
		ActivationToken: p.ActivationCode,
	}

	_, err = shopifyServiceClient.ActivateShopifyBusiness(client.GetContext(), &req)

	return err
}

func (db *partnerDatastore) GetShopifyBusinessByID(bID shared.BusinessID) (*shopify.ShopifyBusiness, error) {
	sn, err := grpc.GetConnectionStringForService(grpc.ServiceNameShopping)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	client, err := grpc.NewInsecureClient(sn)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer client.CloseAndCancel()
	shopifyServiceClient := shopify.NewShopifyBusinessServiceClient(client.GetConn())

	req := shopify.ShopifyBusinessIDRequest{
		BusinessId: bID.ToPrefixString(),
	}

	b, err := shopifyServiceClient.GetShopifyBusinessByBusinessID(client.GetContext(), &req)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return b, nil
}
