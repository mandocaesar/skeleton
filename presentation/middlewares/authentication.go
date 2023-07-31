package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/machtwatch/catalyst-go-skeleton/domain/common/response"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/config"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/config/secret"
	"github.com/machtwatch/catalyst-go-skeleton/utils/jwt"
	"github.com/machtwatch/catalystdk/go/log"
)

type AccountContextKey string

var accountContextKey AccountContextKey = "account"

func (a *ApiMiddleware) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.Contains(header, "Bearer") {
			responseError(w, response.Response[bool]{
				RequestID:  log.GetCtxRequestID(r.Context()),
				Code:       response.CodeUnauthorized,
				Message:    "Invalid token",
				ServerTime: time.Now().Unix(),
			}, response.HttpStatus[response.CodeUnauthorized])
			return
		}

		tokenString := strings.Replace(header, "Bearer ", "", -1)

		jwt := jwt.NewJWT()
		token, err := jwt.ValidateJwtToken(tokenString, config.JWT_METHOD, secret.JWT_SECRET)
		if err != nil {
			responseError(w, response.Response[bool]{
				RequestID:  log.GetCtxRequestID(r.Context()),
				Code:       response.CodeInternalServerError,
				Message:    err.Error(),
				ServerTime: time.Now().Unix(),
			}, response.HttpStatus[response.CodeInternalServerError])
			return
		}

		claims, err := jwt.GetTokenClaims(token)
		if err != nil {
			responseError(w, response.Response[bool]{
				RequestID:  log.GetCtxRequestID(r.Context()),
				Code:       response.CodeInternalServerError,
				Message:    err.Error(),
				ServerTime: time.Now().Unix(),
			}, response.HttpStatus[response.CodeInternalServerError])
			return
		}

		if !a.isAuthTokenInWhitelist(r.Context(), claims["username"].(string), claims["secret_key"].(string)) {
			responseError(w, response.Response[bool]{
				Code:       response.CodeUnauthorized,
				Message:    "Unauthorized",
				ServerTime: time.Now().Unix(),
			}, response.HttpStatus[response.CodeUnauthorized])
			return
		}

		account := Account{
			ID:         int64(claims["account_id"].(float64)),
			ResellerID: int64(claims["reseller_id"].(float64)),
			Username:   claims["username"].(string),
			Email:      claims["email"].(string),
			SecretKey:  claims["secret_key"].(string),
		}

		ctx := context.WithValue(r.Context(), accountContextKey, account)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func responseError(w http.ResponseWriter, err interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Status-Code", strconv.Itoa(code))
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(err)
}

func (a *ApiMiddleware) isAuthTokenInWhitelist(ctx context.Context, username, secretKey string) bool {
	key := fmt.Sprintf("%s%s%s%s%s", config.CACHE_WHITELIST_TOKENS, config.SEPARATOR_KEY, username, config.SEPARATOR_KEY, secretKey)
	val, err := a.cache.GET(ctx, key)
	if err != nil {
		log.StdError(ctx, map[string]interface{}{"key": key}, err, "a.cache.GET() got error - isAuthTokenInWhitelist()")
	}

	return val != ""
}

type Account struct {
	ID         int64
	ResellerID int64
	Username   string
	Email      string
	SecretKey  string
}

func GetAccount(ctx context.Context) (account Account) {
	acc := ctx.Value(accountContextKey)
	if a, ok := acc.(Account); ok {
		account.ID = a.ID
		account.Email = a.Email
		account.Username = a.Username
		account.ResellerID = a.ResellerID
		account.SecretKey = a.SecretKey
	}

	return
}
