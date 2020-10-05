package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/wiseco/go-lib/url"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const timeoutSec = 100

//Client encapsulates grpc client details
type Client interface {
	GetConn() *ggrpc.ClientConn
	GetContext() context.Context
	CloseAndCancel()
}

type client struct {
	conn   *ggrpc.ClientConn
	ctx    context.Context
	cancel context.CancelFunc
}

//GetConnectionStringForService returns a connection string for a serviceName
func GetConnectionStringForService(sn serviceName) (string, error) {
	switch sn {
	case ServiceNameVerification:
		return url.GetSrvVerConnectionString(), nil
	case ServiceNameTransaction:
		return url.GetSrvTxnConnectionString(), nil
	case ServiceNameInvoice:
		return url.GetSrvInvoiceConnectionString(), nil
	case ServiceNameEvent:
		return url.GetSrvEvntConnectionString(), nil
	case ServiceNameBanking:
		return url.GetSrvBankingConnectionString(), nil
	case ServiceNameShopping:
		return url.GetSrvShpConnectionString(), nil
	case ServiceNameAuth:
		return url.GetSrvAuthConnectionString(), nil

	// API
	case ServiceNameApiPartner:
		return url.GetSrvApiPartnerConnectionString(), nil
	case ServiceNameApiClient:
		return url.GetSrvApiClientConnectionString(), nil
	case ServiceNameApiCsp:
		return url.GetSrvApiCspConnectionString(), nil
	}
	return "", fmt.Errorf("Service name must be registered. Unknown:%s", sn)
}

//NewInsecureClient should only be used for testing
func NewInsecureClient(cs string) (Client, error) {
	config := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := ggrpc.Dial(cs, ggrpc.WithTransportCredentials(credentials.NewTLS(config)))

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSec*time.Second)

	return &client{
		conn,
		ctx,
		cancel,
	}, nil
}

//NewClient Returns a new client with safe defaults set
func NewClient(cs string) (Client, error) {
	creds := credentials.NewClientTLSFromCert(nil, "")

	conn, err := ggrpc.Dial(cs, ggrpc.WithTransportCredentials(creds))

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSec*time.Second)

	return &client{
		conn,
		ctx,
		cancel,
	}, nil
}

func (c client) GetConn() *ggrpc.ClientConn {
	return c.conn
}

func (c client) GetContext() context.Context {
	return c.ctx
}

func (c client) CloseAndCancel() {
	c.cancel()
	c.conn.Close()
}
