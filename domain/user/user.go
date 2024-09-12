package user

import (
	"fmt"
	"ifttt/manager/domain/roles"

	"github.com/gin-gonic/gin"
	"github.com/go-viper/mapstructure/v2"
)

type User struct {
	Email    string       `mapstructure:"email" json:"email"`
	Password string       `mapstructure:"password" json:"password"`
	Roles    []roles.Role `mapstructure:"roles" json:"roles"`
}

func GetUserFromContext(c *gin.Context) *User {
	userAny, ok := c.Get("user")
	if !ok {
		return nil
	}
	user, ok := userAny.(*User)
	if !ok {
		return nil
	}
	return user
}

func DecodeUser(input any) (*User, error) {
	var currUser User
	if err := mapstructure.Decode(input, &currUser); err != nil {
		return nil, fmt.Errorf("could not decode user: %s", err)
	}
	return &currUser, nil
}
