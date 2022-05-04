package main

import (
	"fmt"
	"github.com/eltaljohn/go-db/pkg/product"
	"github.com/eltaljohn/go-db/pkg/storage"
	"log"
)

func main() {
	driver := storage.MySQL
	storage.New(driver)

	myStorage, err := storage.DAOProduct(driver)
	if err != nil {
		log.Fatalf("DAOProduct:  %v", err)
	}

	serviceProduct := product.NewService(myStorage)

	ms, err := serviceProduct.GetAll()
	if err != nil {
		log.Fatalf("Product.GetAll %v", err)
	}

	fmt.Println(ms)
}
