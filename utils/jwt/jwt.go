package jwt

import (
	"fmt"

	gojwt "github.com/golang-jwt/jwt/v4"
)

type JWT interface {
	ValidateJwtToken(tokenString, method, secret string) (*gojwt.Token, error)
	GetTokenClaims(token *gojwt.Token) (claims map[string]interface{}, err error)
	GenerateJwtToken(claims gojwt.Claims, method, secret string) (token string, err error)
}

type jwt struct{}

func NewJWT() JWT {
	return &jwt{}
}

func (j *jwt) ValidateJwtToken(tokenString, method, secret string) (*gojwt.Token, error) {
	jwtToken, err := gojwt.Parse(tokenString, func(token *gojwt.Token) (interface{}, error) {
		if gojwt.GetSigningMethod(method) != token.Method {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
	return jwtToken, err
}

func (j *jwt) GetTokenClaims(token *gojwt.Token) (claims map[string]interface{}, err error) {
	if token.Claims.Valid() == nil {
		return token.Claims.(gojwt.MapClaims), nil
	}

	return claims, token.Claims.Valid()
}

func (j *jwt) GenerateJwtToken(claims gojwt.Claims, method, secret string) (token string, err error) {
	jwtToken := gojwt.NewWithClaims(
		gojwt.GetSigningMethod(method),
		claims,
	)
	return jwtToken.SignedString([]byte(secret))
}
