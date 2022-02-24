package main

import (
	l "log"
	"net/http"
	"os"

	pkProduct "github.com/fthdrmzzz/store-service/server/product"
	pkUser "github.com/fthdrmzzz/store-service/server/user"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
)

//admin
//uD5zvmnN92QeTwLM

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)
	//userSvc := pkUser.UserService{}
	//productSvc := pkProduct.ProductService{}

	var userSvc pkUser.UserServer
	userSvc = pkUser.UserService{}
	userSvc = pkUser.LoggingMiddleware{logger, userSvc}

	createUserHandler := httptransport.NewServer(
		pkUser.MakeCreateUserEndpoint(userSvc),
		pkUser.DecodeCreateUserRequest,
		pkUser.EncodeResponse,
	)

	var productSvc pkProduct.ProductServer
	productSvc = pkProduct.ProductService{}
	productSvc = pkProduct.LoggingMiddleware{logger, productSvc}

	addProductHandler := httptransport.NewServer(
		pkProduct.MakeAddProductEndpoint(productSvc),
		pkProduct.DecodeAddProductRequest,
		pkUser.EncodeResponse,
	)

	deleteProductHandler := httptransport.NewServer(
		pkProduct.MakeDeleteProductEndpoint(productSvc),
		pkProduct.DecodeDeleteProductRequest,
		pkUser.EncodeResponse,
	)

	http.Handle("/createUser", createUserHandler)
	http.Handle("/addProduct", addProductHandler)
	http.Handle("/deleteProduct", deleteProductHandler)

	l.Fatal(http.ListenAndServe(":8080", nil))
}
