package grpc

import (
	"fmt"
	"net"
	"os"
	"time"

	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

//Server is an interface describing the grpc server
type Server interface {
	GetGRPCServer() *ggrpc.Server
	Start() error
}

type server struct {
	server *ggrpc.Server
}

//NewServer returns a new server with middelware interceptors registered
func NewServer(sn serviceName) Server {
	certLocation, err := getCertLocation(sn)
	if err != nil {
		panic(err)
	}

	keyLocation, err := getKeyLocation(sn)
	if err != nil {
		panic(err)
	}

	creds, err := credentials.NewServerTLSFromFile(certLocation, keyLocation)
	if err != nil {
		panic(err)
	}

	s := ggrpc.NewServer(
		unaryInterceptor(),
		ggrpc.Creds(creds),
		ggrpc.KeepaliveParams(keepalive.ServerParameters{
			Timeout: 100 * time.Second,
		}),
	)

	return &server{
		s,
	}
}

func (s server) GetGRPCServer() *ggrpc.Server {
	return s.server
}

func (s server) Start() error {
	cp := os.Getenv("GRPC_SERVICE_PORT")

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", cp))
	if err != nil {
		return err
	}

	return s.server.Serve(l)
}
