package grpc

import (
	
     "fmt"
     "context"

     pb "github.com/mantrobuslawal/bfproto/golang/catalog"
     "github.com/mantrobuslawal/bfmicroservices/catalog.git/internal/application/core/domain"
)


func (a Adapter) GetProducts(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	reqSearchType := req.GetSearchType()
	var opt domain.SearchOption
	var searchValue []string
	
	if _ , ok := reqSearchType.(*pb.GetProductRequest_Sku); ok {
		opt = domain.SKU
		searchValue = []string{req.GetSku()}
	} else if _, ok := reqSearchType.(*pb.GetProductRequest_ProductName); ok {
		opt = domain.ProductName
		searchValue = []string{req.GetProductName()}
	} else if _, ok := reqSearchType.(*pb.GetProductRequest_Brand); ok {
		opt = domain.Brand
		searchValue = []string{req.GetBrand()}
	} else if _, ok := reqSearchType.(*pb.GetProductRequest_CatSearch); ok {
		opt = domain.Category
		cat := req.GetCatSearch()
		searchValue = []string{cat.GetCategory(), cat.GetSubCategory()}
	}
	   
	st := domain.SearchType{opt, searchValue}
	result, err := a.api.GetProducts(ctx, st)
	if err != nil {
	   return nil, err
	}

	if result == nil {  // No matching products in database
	   return nil, fmt.Errorf("No products found with property %s matching %s", opt, searchValue[0])
	}

	var products []*pb.Product
	for _, product := range result {
		subcat := product.Subcategory
	    products = append(products, &pb.Product{
		Sku:          product.SKU,
                Name:         product.Name,
                Brand:        product.Brand,
                UnitPrice:    product.UnitPrice,
                Sizes:        product.Sizes,
                Description:  product.Description,
                Category:     product.Category,
                SubCategory:  &subcat, 
	    })
	}
	return &pb.GetProductResponse{Products: products}, nil 
}
