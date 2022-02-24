package product

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

type addProductRequest struct {
	UserId int    `json:"id"`
	Name   string `json:"name"`
	Price  int16  `json:"price"`
}

type addProductResponse struct {
	Err string `json:"err,omitempty"`
}

type deleteProductRequest struct {
	ProductId int `json:"id"`
}

type deleteProductResponse struct {
	Err string `json:"err,omitempty"`
}

func MakeAddProductEndpoint(svc ProductServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addProductRequest)
		err := svc.AddProduct(req.UserId, req.Name, req.Price)
		if err != nil {
			return addProductResponse{"some error happened"}, nil
		}
		return addProductResponse{"no error"}, nil
	}
}

func MakeDeleteProductEndpoint(svc ProductServer) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteProductRequest)
		err := svc.DeleteProduct(req.ProductId)
		if err != nil {
			return deleteProductResponse{"some error happened"}, nil
		}
		return deleteProductResponse{"no error"}, nil
	}
}

func DecodeAddProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request addProductRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
func DecodeDeleteProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request deleteProductRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
