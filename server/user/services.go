package user

import "fmt"

type UserServer interface {
	CreateUser(string, string, string) (string, error)
}

type UserService struct{}

func (us UserService) CreateUser(email string, name string, password string) (string, error) {
	//if email is in database return error
	fmt.Print("")
	//else create the user.
	return "created", nil
}
