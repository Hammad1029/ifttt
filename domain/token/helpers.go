package token

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

func (t *TokenDetails) createToken(expiry int, email string, secret string) error {
	t.Expiry = time.Now().Add(time.Minute * time.Duration(expiry)).Unix()
	t.Claims = jwt.MapClaims{
		"email": email,
		"exp":   t.Expiry,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, t.Claims)
	tokenSigned, err := token.SignedString(secret)
	if err != nil {
		return fmt.Errorf("method *tokenDetails.createToken: could not sign access token: %s", err)
	}
	t.Token = tokenSigned
	return nil
}

func extractToken(header string) string {
	authHeaderSplit := strings.Split(header, " ")
	if len(authHeaderSplit) != 2 {
		return ""
	}
	return authHeaderSplit[1]
}
