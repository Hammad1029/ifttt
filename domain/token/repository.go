package token

type Repository interface {
	StoreTokenPair(email string, tokens *TokenPair) error
	GetTokenPair(email string) (*TokenPair, error)
	DeleteTokenPair(email string) error
}
