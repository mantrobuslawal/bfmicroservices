package grpc

import (
	
     "fmt"
     "context"

     "google.golang.org/protobuf/reflect/protoreflect"

     "github.com/mantrobuslawal/bfproto/golang/catalog"
     "github.com/mantrobuslawal/bfmicroservices/catalog.git/internal/application/core/domain"
)


func (a Adapter) GetProducts(ctx context.Context, req *catalog.GetProductRequest) (*catalog.GetProductResponse, error) {
/* 
 	1. Get protobuf msg
	2. Get msg descriptor
	3. Get field descriptors (fd) from msg descriptor
	4. iterate over fds and see which is not nil for ContainingOneof()
	5. Check all non-nil fds with msg.Has(fd) to see if it's been set
	6. When you have set fd get name using method TextName()
	7. Get value using msg.Get(fd) and type cast to string or use  request get method	
	
*/	

	// Retrieve protobuf message
	msg := req.ProtoReflect()
	
	// Get msg descriptor
	msgDsp := msg.Descriptor()

	// Get fds
	fds := msgDsp.Fields()

	// Iterate over fds to find set Oneof fd
	var fd protoreflect.FieldDescriptor
	
	for i := range fds.Len() {
		fd := fds.Get(i+1)
		if chk := fd.ContainingOneof(); chk == nil {
			continue
		}
		if ok := msg.Has(fd); !ok {
			continue
		}
	}
	
	// Get field name
	searchField := fd.TextName()
	
	var opt domain.SearchOption
	var searchValue []string

	
	switch searchField {
	case "sku":
              opt = domain.SKU
	      searchValue = []string{req.GetSku()}
        
	case "name":
              opt = domain.ProductName
	      searchValue = []string{req.GetProductName()}
  
	case "brand":
              opt = domain.Brand
	      searchValue = []string{req.GetBrand()}
  
	case "category":
              opt = domain.Category
	      _, ok := req.SearchType.(*catalog.GetProductRequest_CatSearch)
	      if ok { cat := req.GetCatSearch() 
		      searchValue = []string{cat.GetCategory(), cat.GetSubCategory()}
	      } 	
        }
 
	st := domain.SearchType{opt, searchValue}
	result, err := a.api.GetProducts(ctx, st)
	if err != nil {
	   return nil, err
	}

	if result == nil {  // No matching products in database
	   return nil, fmt.Errorf("No products found with property %s matching %s", opt, searchValue[0])
	}

	var products []*catalog.Product
	for _, product := range result {
		subcat := product.Subcategory
	    products = append(products, &catalog.Product{
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
	return &catalog.GetProductResponse{Products: products}, nil 
}
