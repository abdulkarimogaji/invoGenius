package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/abdulkarimogaji/invoGenius/config"
	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(subject string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    config.C.App_Uri,
		Subject:   subject,
		Audience:  []string{config.C.App_Uri},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(config.C.Token_Expire))),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})
	return t.SignedString([]byte(config.C.Token_Secret))
}

func ValidateToken(tokenString string) (jwt.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.C.Token_Secret), nil
	}, jwt.WithAudience(config.C.App_Uri), jwt.WithExpirationRequired(), jwt.WithIssuer(config.C.App_Uri))

	switch {
	case token.Valid:
		return token.Claims, nil
	case errors.Is(err, jwt.ErrTokenMalformed):
		return nil, fmt.Errorf("malformed token")
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return nil, fmt.Errorf("invalid signature")
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		return nil, fmt.Errorf("token expired")
	}
	return token.Claims, err
}
