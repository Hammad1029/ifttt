package auth

import (
	"fmt"
	"ifttt/manager/application/config"
	"ifttt/manager/common"
	"strconv"

	"github.com/golang-jwt/jwt"
)

type TokenService struct {
	Secret        string
	AccessExpiry  int
	RefreshExpiry int
}

type TokenPair struct {
	Access  *TokenDetails `json:"access"`
	Refresh *TokenDetails `json:"refresh"`
}

type TokenDetails struct {
	Expiry int64         `mapstructure:"expiry" json:"expiry"`
	Token  string        `mapstructure:"token" json:"token"`
	Claims jwt.MapClaims `mapstructure:"claims" json:"claims"`
}

func NewTokenService() (*TokenService, error) {
	accessExpiry, err := strconv.Atoi(config.GetConfigProp("jwt.expiry.access"))
	if err != nil {
		return nil, fmt.Errorf("method *NewTokenService: could not convert access expiry to int: %s", err)
	}

	refreshExpiry, err := strconv.Atoi(config.GetConfigProp("jwt.expiry.access"))
	if err != nil {
		return nil, fmt.Errorf("method *NewTokenService: could not convert refresh expiry to int: %s", err)
	}

	jwtSecret := config.GetConfigProp("jwt.secret")
	if jwtSecret == "" {
		return nil, fmt.Errorf("method *NewTokenService: could not get JWT secret: %s", err)
	}

	return &TokenService{
		Secret:        jwtSecret,
		AccessExpiry:  accessExpiry,
		RefreshExpiry: refreshExpiry,
	}, nil
}

func (t *TokenService) NewTokenPair(email string) (*TokenPair, error) {
	accessToken := TokenDetails{}
	if err := accessToken.createToken(t.AccessExpiry, email, t.Secret, common.AccessTokenKey); err != nil {
		return nil, fmt.Errorf("method *TokenService.NewTokenPair: could not create access token: %s", err)
	}
	refreshToken := TokenDetails{}
	if err := refreshToken.createToken(t.RefreshExpiry, email, t.Secret, common.RefreshTokenKey); err != nil {
		return nil, fmt.Errorf("method *TokenService.NewTokenPair: could not create refresh token: %s", err)
	}

	return &TokenPair{
		Access:  &accessToken,
		Refresh: &refreshToken,
	}, nil
}

func (t *TokenService) VerifyToken(header string) (*TokenDetails, error) {
	bearerToken := extractToken(header)
	if bearerToken == "" {
		return nil, nil
	}

	token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("method VerifyToken: unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(t.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("method VerifyToken: could not cast claims")
	}
	expiry, ok := claims["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("method VerifyToken: could not cast expiry")
	}

	return &TokenDetails{Expiry: int64(expiry), Token: bearerToken, Claims: claims}, nil
}
