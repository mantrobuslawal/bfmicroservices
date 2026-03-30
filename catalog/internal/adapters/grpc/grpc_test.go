package grpc

import (
   "context"
   "log"
   "net"
   "testing"

   "google.golang.org/grpc"
   "google.golang.org/grpc/reflection"
   "google.golang.org/grpc/test/bufconn"
  
    "github.com/mantrobuslawal/bfproto/golang/catalog"
    "github.com/mantrobuslawal/bfmircoservices/catalog.git/internal/application/core/api"
    "github.com/mantrobuslawal/bfmircoservices/catalog.git/internal/adapters/repository"
    
)  

const (
    address = "localhost:50051"
    bufSize = 1024 * 1024
    port = 50051
)

var listener *bufconn.Listener

func TestServer_GetProducts(t *testing.T) {
	ctx := context.Background()
	initGRPCServerBuffConn(t)
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(getBufDialer(listner)), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	cc := catalog.NewCatalogClient(conn)
	
	skuQuery := catalog.*GetProductRequest_Sku("abdcdegh12345")
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	res, err := cc.GetProducts(ctx, &catalog.GetProductRequest{SearchType: skuQuery})
	if err != nil {
		log.Fatalf("Could not retrieve product: %v", err)
	}	
	log.Printf(res.Value)
}

func getBufDialer(listener *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, url string) (net.Conn, error) {
		return listner.Dial()
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
	reflection.Register(s)
        
	go func() {
		log.Printf("starting catalog service on port %d...", s.port)
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}


func initApp(t *testing.T) *api.Application {
	t.Helper()

	repo := repository.NewAdapter(repository.SliceCatalog)
	app := api.NewApplication(repo)

	return app	
}
 
