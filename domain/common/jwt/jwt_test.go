package jwt

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWTToken(t *testing.T) {
	t.Run("it should gerenate a jwt token", func(t *testing.T) {
		claims := JwtClaims{
			StandardClaims: jwt.StandardClaims{
				Issuer:    "local",
				ExpiresAt: time.Now().Add(3 * time.Hour).Unix(),
			},
			Email: "fabs@mail.com",
		}
		authToken, err := claims.GenerateJwtToken("HS256", "")
		assert.Nil(t, err)
		assert.True(t, authToken != "")
	})

	t.Run("it should gerenate a jwt token with default method HS256", func(t *testing.T) {
		claims := JwtClaims{
			StandardClaims: jwt.StandardClaims{
				Issuer:    "local",
				ExpiresAt: time.Now().Add(3 * time.Hour).Unix(),
			},
			Email: "fabs@mail.com",
		}
		authToken, err := claims.GenerateJwtToken("", "")
		assert.Nil(t, err)
		assert.True(t, authToken != "")
	})
}

func TestValidateJWTToken(t *testing.T) {
	claims := JwtClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "local",
			ExpiresAt: time.Now().Add(3 * time.Hour).Unix(),
		},
		Email: "fabs@mail.com",
	}
	authToken, _ := claims.GenerateJwtToken("HS256", "FABS_SECRET")

	t.Run("it should success validate jwt token", func(t *testing.T) {
		token, err := ValidateJwtToken(authToken, "HS256", "FABS_SECRET")
		assert.Nil(t, err)
		assert.True(t, token.Valid)
		assert.Nil(t, token.Claims.Valid())
	})

	t.Run("it should success validate jwt token with default method HS256", func(t *testing.T) {
		token, err := ValidateJwtToken(authToken, "", "FABS_SECRET")
		assert.Nil(t, err)
		assert.True(t, token.Valid)
		assert.Nil(t, token.Claims.Valid())
	})

	t.Run("it should failed validate jwt token", func(t *testing.T) {
		_, err := ValidateJwtToken(authToken, "HS256", "wadsadwad")
		assert.NotNil(t, err)
	})

	t.Run("it should failed validate jwt token because method is different", func(t *testing.T) {
		_, err := ValidateJwtToken(authToken, "HS512", "FABS_SECRET")
		assert.NotNil(t, err)
	})
}

func TestGetClaims(t *testing.T) {
	claims := JwtClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "local",
			ExpiresAt: time.Now().Add(3 * time.Hour).Unix(),
		},
		Email: "fabs@mail.com",
	}
	authToken, _ := claims.GenerateJwtToken("HS256", "FABS_SECRET")
	token, _ := ValidateJwtToken(authToken, "HS256", "FABS_SECRET")

	t.Run("it should success get claims", func(t *testing.T) {
		respClaims, err := GetTokenClaims(token)
		assert.Nil(t, err)
		assert.Equal(t, claims.Email, fmt.Sprintf("%v", respClaims.Email))
	})
}
