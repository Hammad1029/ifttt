package auth

import "context"

type Repository interface {
	StoreTokenPair(email string, tokens *TokenPair, ctx context.Context) error
	GetTokenPair(email string, ctx context.Context) (*TokenPair, error)
	DeleteTokenPair(email string, ctx context.Context) error
}
