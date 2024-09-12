package user

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CreateUserRequest struct {
	Email    string `mapstructure:"email" json:"email"`
	Password string `mapstructure:"password" json:"password"`
}

func (c *CreateUserRequest) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Email, validation.Required, is.Email),
		validation.Field(&c.Password, validation.Required, validation.Length(5, 50)),
	)
}
