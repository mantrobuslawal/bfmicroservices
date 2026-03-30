package grpc

import (
   "context"
   "log"
   "net"
   "testing"
   "time"

   "google.golang.org/grpc"
   "google.golang.org/grpc/reflection"
   "google.golang.org/grpc/test/bufconn"
  
    "github.com/mantrobuslawal/bfproto/golang/catalog"
    "github.com/mantrobuslawal/bfmicroservices/catalog.git/internal/application/core/api"
    "github.com/mantrobuslawal/bfmicroservices/catalog.git/internal/adapters/repository"
    
)  

const (
    address = "localhost:50051"
    bufSize = 1024 * 1024
    port = 50051
)

var listener *bufconn.Listener

func TestServer_GetProducts(t *testing.T) {
	initGRPCServerBuffCon(t)
	cc, _ := NewClientAdapter()
	
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	
	// search by SKU
	res, err := cc.GetProducts(ctx, &catalog.GetProductRequest{SearchType: &catalog.GetProductRequest_Sku{Sku: "abdcdegh12345"}})
	if err != nil {
		log.Fatalf("Could not retrieve product: %v", err)
	}	
	log.Printf("%v", res)
}

type clientAdapter struct {
	catalog catalog.CatalogClient
}

func NewClientAdapter() (*clientAdapter, error) {
	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(getBufDialer(listener)), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := catalog.NewCatalogClient(conn)
	return &clientAdapter{catalog: client}, nil
}

func (ca *clientAdapter) GetProducts(ctx context.Context, in *catalog.GetProductRequest) (*catalog.GetProductResponse, error) {
	res, err := ca.catalog.GetProducts(ctx, in)
	return res, err
} 

func getBufDialer(listener *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, url string) (net.Conn, error) {
		return listener.Dial()
	}
}

// Initialization of BufConn and catalog service business
// logic. Ignoring server Run() method and using this
// to so bufConn can be utilised in place of actutal tcp listner
func initGRPCServerBuffCon(t *testing.T){
	t.Helper()

	listener := bufconn.Listen(bufSize)
	s := NewAdapter(initApp(t), port)

	grpcServer := grpc.NewServer()
	s.server = grpcServer
	catalog.RegisterCatalogServer(grpcServer, s)
        // Register reflection service on gRPC server
	reflection.Register(grpcServer)
        
	go func() {
		log.Printf("starting catalog service on port %d...", s.port)
		if err := s.server.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}


func initApp(t *testing.T) *api.Application {
	t.Helper()

	repo, _ := repository.NewAdapter(repository.SliceCatalog)
	app := api.NewApplication(repo)

	return app	
}
 
