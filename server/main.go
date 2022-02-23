package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/transport/http"
)

//admin
//uD5zvmnN92QeTwLM

const (
	isAdmin = 1 << iota
	isSeller
	isCustomer
)

type product struct {
	productId int
	ownerId   int
	name      string
	price     int16
}

type user struct {
	userId int
	roles  byte
	email  string
	name   string
}

type UserService interface {
	CreateUser(string, string, string) (string, error)
}

type ProductService interface {
	AddProduct(int, string, int16) error
	DeleteProduct(int) error
	GetAllProducts() error
}

type userService struct{}

func (userService) CreateUser(email string, name string, password string) (string, error) {
	//if email is in database return error
	fmt.Println("inside Create user")
	//else create the user.
	return "created", nil
}

type productService struct{}

func (productService) AddProduct(userId int, name string, price int16) error {
	//add product
	fmt.Println("product added")
	return nil
}

func (productService) DeleteProduct(productId int) error {
	//product with that id
	fmt.Println("product with id % v deleted", productId)
	return nil
}

func (productService) GetAllProducts() error {
	fmt.Println("Returning All products")
	return nil
}

type createUserRequest struct {
	Name     string `json:"name"`
	Mail     string `json:"mail"`
	Password string `json:"password"`
}

type createUserResponse struct {
	msg string `json: "msg"`
	Err string `json:"err,omitempty"`
}

type addProductRequest struct {
	userId int    `json:"id"`
	name   string `json:"name"`
	price  int16  `json:"price"`
}

type addProductResponse struct {
	Err string `json:"err,omitempty"`
}

type deleteProductRequest struct {
	productId int `json:"id"`
}

type deleteProductResponse struct {
	Err string `json:"err,omitempty"`
}

func makeCreateUserEndpoint(svc UserService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(createUserRequest)
		_, err := svc.CreateUser(req.Name, req.Mail, req.Password)
		if err != nil {
			return createUserResponse{
				"some error happened",
				err.Error(),
			}, nil
		}
		return createUserResponse{"successful operation", ""}, nil
	}
}

func makeAddProductEndpoint(svc ProductService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Printf("Printing context: %v\n", ctx)
		req := request.(addProductRequest)
		err := svc.AddProduct(req.userId, req.name, req.price)
		if err != nil {
			return addProductResponse{"some error happened"}, nil
		}
		return addProductResponse{"no error"}, nil
	}
}

func makeDeleteProductEndpoint(svc ProductService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteProductRequest)
		err := svc.DeleteProduct(req.productId)
		if err != nil {
			return deleteProductResponse{"some error happened"}, nil
		}
		return deleteProductResponse{"no error"}, nil
	}
}

func main() {
	userSvc := userService{}
	productSvc := productService{}

	createUserHandler := httptransport.NewServer(
		makeCreateUserEndpoint(userSvc),
		decodeCreateUserRequest,
		encodeResponse,
	)

	addProductHandler := httptransport.NewServer(
		makeAddProductEndpoint(productSvc),
		decodeAddProductRequest,
		encodeResponse,
	)

	deleteProductHandler := httptransport.NewServer(
		makeDeleteProductEndpoint(productSvc),
		decodeDeleteProductRequest,
		encodeResponse,
	)

	http.Handle("/createUser", createUserHandler)
	http.Handle("/addProduct", addProductHandler)
	http.Handle("/deleteProduct", deleteProductHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func decodeCreateUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeAddProductRequest(_ context.Context, r *http.Request) interface{}

func decodeDeleteProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request deleteProductRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
