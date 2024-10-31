package auth

import "ifttt/manager/domain/user"

type LoginResponse struct {
	Tokens *TokenPair `json:"tokens" mapstructure:"tokens"`
	User   *user.User `json:"user" mapstructure:"user"`
}
