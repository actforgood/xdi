# Xdi

[![Build Status](https://github.com/actforgood/xdi/actions/workflows/build.yml/badge.svg)](https://github.com/actforgood/xdi/actions/workflows/build.yml)
[![License](https://img.shields.io/badge/license-MIT-blue)](https://raw.githubusercontent.com/actforgood/xdi/main/LICENSE)
[![Coverage Status](https://coveralls.io/repos/github/actforgood/xdi/badge.svg?branch=main)](https://coveralls.io/github/actforgood/xdi?branch=main)
[![Goreportcard](https://goreportcard.com/badge/github.com/actforgood/xdi)](https://goreportcard.com/report/github.com/actforgood/xdi)
[![Go Reference](https://pkg.go.dev/badge/github.com/actforgood/xdi.svg)](https://pkg.go.dev/github.com/actforgood/xdi)  

---

Package `xdi` provides a centralized dependency injection manager which holds definitions for an application's objects/dependencies.  


### Installation

```shell
$ go get -u github.com/actforgood/xdi
```


### Example
Basic example:  
```go
// DiManager holds application's objects, dependencies.
// Do not inject it/use it directly, in your application's objects.
// It should be used only in the bootstrap process of your application and/or main.go,
// as a centralized container of dependencies.
var DiManager = xdi.NewManager()

func init() {
	DiManager.AddDefinition(xdi.Definition{
		ID: "app.repository.product",
		Initializer: func() interface{} {
			return NewDummyProductRepository()
		},
		Shared: true,
	})
}

func init() {
	DiManager.AddDefinition(xdi.Definition{
		ID: "app.service.product",
		Initializer: func() interface{} {
			return NewDummyProductService(
				DiManager.Get("app.repository.product").(ProductRepository),
			)
		},
		Shared: true,
	})
}

func main() {
	productService := DiManager.Get("app.service.product").(ProductService)
	isAvailable, _ := productService.CheckAvailability("some-sku", 2)
	fmt.Println("isAvailable:", isAvailable)
}
```


### Misc 
Feel free to use this pkg if you like it and fits your needs.   
As it is a light/lite pkg, you can also just copy-paste the code, instead of importing it, keeping the license header.  


### License
This package is released under a MIT license. See [LICENSE](LICENSE).  
