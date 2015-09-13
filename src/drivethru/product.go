package main

import (
    "fmt"
    "os"
    "encoding/json"
    "io/ioutil"
	"log"
)

type Product struct {
	Code string
	Description string
	Amount int
}

type Products []Product

func (products Products) Find(code string) (product *Product) {

	for _, p := range products {
		if p.Code == code {
			product = &p
			break;
		}
	}

	return
}

/*
 * Our 'Db' of products is loaded from a json file.
 */
func load_products(filename string) (Products) {

	file, e := ioutil.ReadFile(filename)

    if e != nil {
        fmt.Printf("File error: %v\n", e)
        os.Exit(1)
    }

	var products Products
	err := json.Unmarshal(file, &products)
	
	if err != nil {
		log.Printf("Error reading json %v", err)
		os.Exit(1)
	}

	return products
}
