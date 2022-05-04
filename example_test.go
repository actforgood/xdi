package xdi_test

import (
	"fmt"

	"github.com/actforgood/xdi"
)

// some application logic

// ProductRepository ...
type ProductRepository interface {
	// GetBySKU returns a product data by SKU.
	GetBySKU(sku string) (map[string]interface{}, error)
}

// DummyProductRepository is the default implementation for ProductRepository contract.
type DummyProductRepository struct{}

func (DummyProductRepository) GetBySKU(sku string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"sku":   sku,
		"price": 99.9,
		"stock": uint(100),
	}, nil
}

// NewDummyProductRepository returns new instance of
// a DummyProductRepository.
func NewDummyProductRepository() *DummyProductRepository {
	return &DummyProductRepository{}
}

// ProductService ...
type ProductService interface {
	// CheckAvailability checks if a product has given quantity available.
	CheckAvailability(sku string, qty uint) (bool, error)
}

// DummyProductService is the default implementation for ProductService contract.
type DummyProductService struct {
	repo ProductRepository
}

// NewDummyProductService returns new instance of
// a dummyProductService.
func NewDummyProductService(repo ProductRepository) *DummyProductService {
	return &DummyProductService{repo: repo}
}

func (srv DummyProductService) CheckAvailability(sku string, qty uint) (bool, error) {
	product, err := srv.repo.GetBySKU(sku)
	if err != nil {
		return false, err
	}
	isAvailable := product["stock"].(uint) >= qty

	return isAvailable, nil
}

// DiManager holds application's objects, dependencies.
// Do not inject it/use it directly, in your application's objects.
// It should be used only in the bootstrap process of your application and/or main.go,
// as a centralized container of dependencies.
var DiManager = xdi.NewDiManager()

const diProductRepoID = "app.repository.Product"

func init() {
	DiManager.AddDefinition(xdi.DiManagerDef{
		ID: diProductRepoID,
		Initializer: func() interface{} {
			return NewDummyProductRepository()
		},
		Shared: true,
	})
}

const diProductServiceID = "app.service.Product"

func init() {
	DiManager.AddDefinition(xdi.DiManagerDef{
		ID: diProductServiceID,
		Initializer: func() interface{} {
			return NewDummyProductService(
				DiManager.Get(diProductRepoID).(ProductRepository),
			)
		},
		Shared: true,
	})
}

// ... so on and so forth, build entire application objects.

// end some application logic

func ExampleDiManager() {
	productService := DiManager.Get(diProductServiceID).(ProductService)
	isAvailable, _ := productService.CheckAvailability("some-sku", 2)
	fmt.Println("isAvailable:", isAvailable)

	// Output:
	// isAvailable: true
}
