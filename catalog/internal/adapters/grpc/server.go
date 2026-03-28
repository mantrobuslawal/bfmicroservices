package grpc

import (
   "fmt"
   "net"
   "google.golang.org/grpc"
   "google.golang.org/grpc/reflection"
   log "github.com/sirupsen/logrus"
   
   "github.com/mantrobuslawal/bfproto/golang/catalog"
   "github.com/mantrobuslawal/bfmicroservices/catalog.git/config"
   "github.com/mantrobuslawal/bfmicroservices/catalog.git/internal/ports"
)

type Adapter struct {
	api ports.APIPort
        port int
        catalog.UnimplementedCatalogServer
}

func NewAdapter(api ports.APIPort, port int) *Adapter {
	return &Adapter{api: api, port: port}
}

func (a Adapter) Run() {
     var err error

     listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
     if err != nil {
	log.Fatalf("failed to listen on port %d, error: %v" a.port, err)
     }

     grpcServer := grpc.NewServer()
     a.server = grpcServer
     catalog.RegisterCatalogServer(grpcServer, a)
     if config.GetEnv() == "development" {
         reflection.Register(grpcServer)
     }

     
     log.Printf("starting catalog service on port %d ...", a.port)
     if err := grpcServer.Serve(listen); err != nil {
	log.Fatalf("failed to serve grpc on port %d", a.port)
     }
}

func (a Adapter) Stop() {
	a.server.Stop()
}
