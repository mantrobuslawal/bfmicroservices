package grpc

import (
   "context"
   "log"
   "net"
   "testing"
   "time"
  // "fmt"

   "google.golang.org/grpc"
   "google.golang.org/grpc/reflection"
   "github.com/stretchr/testify/assert"
  
    pb "github.com/mantrobuslawal/bfproto/golang/catalog"
    "github.com/mantrobuslawal/bfmicroservices/catalog.git/internal/application/core/api"
    repo "github.com/mantrobuslawal/bfmicroservices/catalog.git/internal/adapters/repository"
    "github.com/mantrobuslawal/bfmicroservices/catalog.git/internal/application/core/domain"
)  

const (
    address = "localhost:50051"
    port = 50051
)

type testTable map[string]struct{req *pb.GetProductRequest
				products []domain.Product
				expectedErr error
			 }
// TODO: created sentinel errors & match errors returned by grpc
// to implement correct errors are returned to client!!!
func TestServer_GetProducts(t *testing.T) {
	goodSubcat := "wall decor" // Used in subcategory tests
	badSubcat := "foobarbazz" // Used in subcategory tests
	var tests testTable
	tests = testTable{
			"sku in repo": {
				req: &pb.GetProductRequest{SearchType: &pb.GetProductRequest_Sku{Sku: repo.SliceCatalog[0].SKU}},
				products: repo.SliceCatalog[0:1],
			},
			"sku not in repo": {
				req: &pb.GetProductRequest{SearchType: &pb.GetProductRequest_Sku{Sku:"0040505050"}},
				expectedErr: nil,	//TODO: incorrect update to return correct grpc status error			
			},
			"sku empty string": {
				req: &pb.GetProductRequest{SearchType: &pb.GetProductRequest_Sku{}},
				expectedErr: nil,      //TODO: incorrect update to return correct grpc status error	
			},
			"name in repo": {
				req: &pb.GetProductRequest{SearchType: &pb.GetProductRequest_ProductName{ProductName:"gopher desk"}},
				products: repo.SliceCatalog[0:1],
			},
			"name not in repo": {
				req: &pb.GetProductRequest{SearchType: &pb.GetProductRequest_ProductName{ProductName:"spam and eggs"}},
				expectedErr: nil, //TODO: incorrect update to return correct grpc status error	
			},
			"name as empty string": {
				req: &pb.GetProductRequest{SearchType: &pb.GetProductRequest_ProductName{}},
				expectedErr: nil, //TODO: incorrect update to return correct grpc status error	
			},
			"brand in repo": {
				req: &pb.GetProductRequest{SearchType: &pb.GetProductRequest_Brand{Brand: "rob pike tapestry"}},
				products: repo.SliceCatalog[1:2],
			},
			"brand not in repo": {
				req: &pb.GetProductRequest{SearchType: &pb.GetProductRequest_Brand{Brand: "foo furniture"}},
				expectedErr: nil, //TODO: incorrect update to return correct grpc status error	
			},
			"brand as empty string": {
				req: &pb.GetProductRequest{SearchType: &pb.GetProductRequest_Brand{}},
				expectedErr: nil, //TODO: incorrect update to return correct grpc status error	
			},
			"category in repo": {
				req: &pb.GetProductRequest{SearchType: &pb.GetProductRequest_CatSearch{CatSearch: 
									&pb.Category{Category: "office furniture" }}},
				products: []domain.Product{repo.SliceCatalog[0], repo.SliceCatalog[2]},
			},
			"category and subcategory in repo": {
				req: &pb.GetProductRequest{SearchType: &pb.GetProductRequest_CatSearch{CatSearch:
						&pb.Category{Category: "home decor", SubCategory: &goodSubcat}}},
				products: repo.SliceCatalog[1:2],
			},
			"category in repo, but subcategory not in repo": {
				req: &pb.GetProductRequest{SearchType: &pb.GetProductRequest_CatSearch{CatSearch:
					&pb.Category{Category: "home decor", SubCategory: &badSubcat}}},
				expectedErr: nil, //TODO: incorrect update to return correct grpc status error	
			},
			"category empty string and subcategory non-empty string": {
				req: &pb.GetProductRequest{SearchType: &pb.GetProductRequest_CatSearch{CatSearch:
					&pb.Category{SubCategory: &goodSubcat}}},
				expectedErr: nil, //TODO: incorrect update to return correct grpc status error	
			},
			"category and subcategory not in repo": {
				req: &pb.GetProductRequest{SearchType: &pb.GetProductRequest_CatSearch{CatSearch:
					&pb.Category{Category: badSubcat, SubCategory: &badSubcat}}},
				expectedErr: nil, //TODO: incorrect update to return correct grpc status error	
			},
			
		}	


	initGRPCServer(t)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	cc := pb.NewCatalogClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T){
			res, _ := cc.GetProducts(ctx, tc.req)
			/*
			if err != nil && tc.expectedErr != nil {
				assert.Equal(t, err.Error(), tc.expectedErr.Error())
			}
			assert.Equal(t, tc.expectedErr, err)
			*/
			got := res.GetProducts()
			var gotProducts []domain.Product
			for _, product := range got {
				subcat := *product.SubCategory
				gotProducts = append(gotProducts, domain.Product{
					SKU: product.Sku,
					Name: product.Name,
					Brand: product.Brand,
					UnitPrice: product.UnitPrice,
					Sizes: product.Sizes,
					Description: product.Description,
					Category: product.Category,
					Subcategory: subcat,	
				})
			}
		       assert.Equal(t, tc.products, gotProducts)
		})
	}
}

func initGRPCServer(t *testing.T) {
	t.Helper()
	
	lis, err := net.Listen("tcp", address)
	
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	app := NewAdapter(initApp(t), port)
	grpcServer := grpc.NewServer()
	app.server = grpcServer
	pb.RegisterCatalogServer(grpcServer, app)
	reflection.Register(grpcServer)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}

func initApp(t *testing.T) *api.Application {
	t.Helper()

	repo, _ := repo.NewAdapter(repo.SliceCatalog)
	app := api.NewApplication(repo)

	return app	
}
 
