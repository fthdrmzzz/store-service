package product

import "fmt"

type ProductServer interface {
	AddProduct(int, string, int16) error
	DeleteProduct(int) error
	GetAllProducts() error
}

type ProductService struct{}

func (ps ProductService) AddProduct(userId int, name string, price int16) error {
	//add product
	fmt.Print("")
	return nil
}

func (ps ProductService) DeleteProduct(productId int) error {
	//product with that id
	fmt.Print("")
	return nil
}

func (ps ProductService) GetAllProducts() error {
	fmt.Print("")
	return nil
}
