package user

const (
	isAdmin = 1 << iota
	isSeller
	isCustomer
)

type user struct {
	userId int
	roles  byte
	email  string
	name   string
}
