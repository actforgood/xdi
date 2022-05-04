# Xdi

[![Build Status](https://github.com/actforgood/xdi/actions/workflows/build.yml/badge.svg)](https://github.com/actforgood/xdi/actions/workflows/build.yml)
[![License](https://img.shields.io/badge/license-MIT-blue)](https://raw.githubusercontent.com/actforgood/xdi/main/LICENSE)
[![Coverage Status](https://coveralls.io/repos/github/actforgood/xdi/badge.svg?branch=main)](https://coveralls.io/github/actforgood/xdi?branch=main)
[![Go Reference](https://pkg.go.dev/badge/github.com/actforgood/xdi.svg)](https://pkg.go.dev/github.com/actforgood/xdi)  

---

Package `xdi` provides a centralized dependency injection manager which holds definitions for an application's objects/dependencies.  


### Example
Basic example:  
```golang
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

func main() {
	productService := DiManager.Get(diProductServiceID).(ProductService)
	isAvailable, _ := productService.CheckAvailability("some-sku", 2)
	fmt.Println("isAvailable:", isAvailable)
}
```


### Misc 
Feel free to use this pkg if you like it and fits your needs.   
As it is a light/lite pkg, you can also just copy-paste the code, instead of importing it, keeping the license header.  


### License
This package is released under a MIT license. See [LICENSE](LICENSE).  
