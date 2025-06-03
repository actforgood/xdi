package xdi_test

import (
	"fmt"

	"github.com/actforgood/xdi"
)

// some application logic

// ProductRepository ...
type ProductRepository interface {
	// GetBySKU returns a product data by SKU.
	GetBySKU(sku string) (map[string]any, error)
}

// DummyProductRepository is the default implementation for ProductRepository contract.
type DummyProductRepository struct{}

// NewDummyProductRepository returns new instance of
// a DummyProductRepository.
func NewDummyProductRepository() *DummyProductRepository {
	return &DummyProductRepository{}
}

func (DummyProductRepository) GetBySKU(sku string) (map[string]any, error) {
	return map[string]any{
		"sku":   sku,
		"price": 99.9,
		"stock": uint(100),
	}, nil
}

// ProductService ...
type ProductService interface {
	// CheckAvailability checks if a product has given quantity available.
	CheckAvailability(sku string, qty uint) (bool, error)
}

// DummyProductService is the default implementation for [ProductService] contract.
type DummyProductService struct {
	repo ProductRepository
}

// NewDummyProductService returns new instance of
// a DummyProductService.
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
// Note: instead of declaring a variable, you can also use the singleton provided by xdi.ManagerInstance().
var DiManager = xdi.NewManager()

func init() {
	DiManager.AddDefinition(xdi.Definition{
		ID: "app.repository.product",
		Initializer: func() any {
			return NewDummyProductRepository()
		},
		Shared: true,
	})
}

func init() {
	DiManager.AddDefinition(xdi.Definition{
		ID: "app.service.product",
		Initializer: func() any {
			return NewDummyProductService(
				DiManager.Get("app.repository.product").(ProductRepository),
			)
		},
		Shared: true,
	})
}

// ... so on and so forth, build entire application objects.

// end some application logic

func ExampleManager() {
	productService := DiManager.Get("app.service.product").(ProductService)
	isAvailable, _ := productService.CheckAvailability("some-sku", 2)
	fmt.Println("isAvailable:", isAvailable)

	// Output:
	// isAvailable: true
}
