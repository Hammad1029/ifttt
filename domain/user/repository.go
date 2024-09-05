package user

type Repository interface {
	GetUser(email string, decodeFunc func(input any) (*User, error)) (*User, error)
	CreateUser(user User) error
}
