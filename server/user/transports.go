package user

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

type createUserRequest struct {
	Name     string `json:"name"`
	Mail     string `json:"mail"`
	Password string `json:"password"`
}

type createUserResponse struct {
	Msg string `json:"msg"`
	Err string `json:"err,omitempty"`
}

func DecodeCreateUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func MakeCreateUserEndpoint(svc UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createUserRequest)
		_, err := svc.CreateUser(req.Name, req.Mail, req.Password)
		if err != nil {
			return createUserResponse{"Oh no", "some error happened"}, nil
		}
		return createUserResponse{"Oh yes", "no error"}, nil
	}
}
