package grpc

import (
     "context"

     "github.com/mantrobuslawal/bfproto/golang/catalog"
     "github.com/mantrobuslawal/bfmicroservices/catalog.git/internal/application/core/domain"
)


func (a Adapter) Get(ctx context.Context, req *catalog.GetProductRequest) (*catalog.GetProductResponse, error) {
	var opt domain.SearchOption
	var searchValue []string

	switch (req.GetSearchType()).(type) {
	case *GetProductRequest_Sku:
              opt = domain.SKU
              searchValue = []string{req.GetSku()}
        
	case *GetProductRequest_ProductName:
              opt = domain.ProductName
              searchValue = []string{req.GetProductName()}
  
	}
	case *GetProductRequest_Brand:
              opt = domain.Brand
              searchValue = []string{req.GetBrand()}
  
	}

	case *GetProductRequest_CatSearch:
              opt = domain.Category
              searchValue = []string{req.GetCatSearch().GetCategory(), req.GetCatSearch().GetSubCategory()}
  
	}
	st := domain.SearchType{opt, searchValue}
	result, err := a.api.GetProducts(ctx, st)
	if err != nil {
	   return nil, err
	}

	if result == nil {  // No matching products in database
	   return nil, fmt.Errorf("No products found with property %s matching %s", opt, searchValue[0])
	}

	var products []*catalog.Products
	for _, product := range result {
	    products = append(products, &catalog.Product{
		Sku:          product.SKU,
                Name:         product.Name,
                Brand:        product.Brand,
                UnitPrice:    product.UnitPrice,
                Sizes:        product.Sizes,
                Description:  product.Description,
                Category:     product.Category,
                SubCategory: *product.Subcategory, 
	    }
	}
	return &catalog.GetProductResponse{Products: products}, nil 
}
