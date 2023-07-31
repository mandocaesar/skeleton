package jwt

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/config/secret"
	"github.com/s12v/go-jwks"
	"github.com/square/go-jose"
)

type (
	JwtClaims struct {
		jwt.StandardClaims
		Audience []string               `json:"aud,omitempty"` // overriden from StandardClaims because token can have multiple audience
		Email    string                 `json:"sub"`           // must include valid url as namespace for email, TBD to make this configurable
		Scope    []string               `json:"scp,omitempty"`
		Ext      map[string]interface{} `json:"ext,omitempty"`
	}
	JwtModuleAccess struct {
		ModuleName      string   `json:"module_name"`
		JwtModuleAction []string `json:"action,omitempty"`
	}
	IDContextKey         struct{}
	NameContextKey       struct{}
	EmailContextKey      struct{}
	ScopeContextKey      struct{}
	RoleContextKey       struct{}
	PermissionContextKey struct{}
)

// GenerateJwtToken is
func (j JwtClaims) GenerateJwtToken(method, secret string) (token string, err error) {
	if method == "" {
		method = "HS256"
	}
	jwtToken := jwt.NewWithClaims(
		jwt.GetSigningMethod(method),
		j,
	)
	return jwtToken.SignedString([]byte(secret))
}

// ValidateJwtToken is
func ValidateJwtToken(tokenString, method, secret string) (*jwt.Token, error) {
	jwtToken, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if method == "" {
			method = "HS256"
		}
		if jwt.GetSigningMethod(method) != token.Method {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		if jwt.GetSigningMethod(method).Alg() == "RS256" {
			return GetRSAPublicKey()
		}

		return []byte(secret), nil
	})

	return jwtToken, err
}

// GetTokenClaims is
func GetTokenClaims(token *jwt.Token) (claims *JwtClaims, err error) {
	if token.Claims.Valid() == nil {
		return token.Claims.(*JwtClaims), nil
	}

	return claims, token.Claims.Valid()
}

func GetRSAPublicKey() (*rsa.PublicKey, error) {
	jwksSource := jwks.NewWebSource(secret.JWKS_URL)
	jwksClient := jwks.NewDefaultClient(
		jwksSource,
		time.Duration(secret.JWKS_REFRESH)*time.Hour,
		time.Duration(secret.JWKS_TTL)*time.Hour,
	)

	var jwk *jose.JSONWebKey
	jwk, err := jwksClient.GetEncryptionKey("public:hydra.jwt.access-token")
	if err != nil {
		return nil, err
	}

	return jwk.Key.(*rsa.PublicKey), nil
}
