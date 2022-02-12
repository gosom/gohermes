package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gosom/gohermes/pkg/container"
	"github.com/gosom/gohermes/pkg/utils"
)

const (
	headerAuthorization = "Authorization"
	headerPrefixBearer  = "BEARER"
	headerApiKey        = "X-API-KEY"
)

type JwtClaims struct {
	User  AuthenticatedUser      `json:"user`
	Extra map[string]interface{} `json:"extra,omitempty"`
	jwt.StandardClaims
}

func CreateJwtAccessToken(signingKey, issuer string, duration time.Duration, u IUser) (string, error) {
	claims := JwtClaims{
		User:  AuthenticatedUser{u.GetID(), u.GetRoles()},
		Extra: nil,
		StandardClaims: jwt.StandardClaims{
			Issuer:    issuer,
			Subject:   strconv.Itoa(u.GetID()),
			ExpiresAt: time.Now().UTC().Add(duration).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(signingKey))
}

func CreateJwtRefreshToken(signingKey, issuer, access string, duration time.Duration) (string, error) {
	claims := JwtClaims{
		Extra: map[string]interface{}{"token": access},
		StandardClaims: jwt.StandardClaims{
			Issuer:    issuer,
			ExpiresAt: time.Now().UTC().Add(duration).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(signingKey))
}

func AuthenticationJWT(di *container.ServiceContainer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if IsAuthenticatedUser(r.Context()) {
				next.ServeHTTP(w, r)
				return
			}
			tokenType, token, err := getAuthTokenFromHeader(r, false)
			if ae := utils.ApiErrorFromErr(err); ae != nil {
				utils.RenderJson(r, w, ae.StatusCode, ae)
				return
			}

			var user AuthenticatedUser
			switch tokenType {
			case headerPrefixBearer:
				user, err = ValidateAccessToken(di.Cfg.TokenSecret, token)
				if err != nil {
					di.Logger.Error().Msg(err.Error())
					ae := utils.NewAuthenticationError("")
					utils.RenderJson(r, w, ae.StatusCode, ae)
					return
				}
			default:
				ae := utils.NewAuthenticationError("")
				utils.RenderJson(r, w, ae.StatusCode, ae)
				return
			}
			ctx := context.WithValue(r.Context(), utils.Authenticated, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AuthenticationXAPIKey(di *container.ServiceContainer) func(next http.Handler) http.Handler {
	srv := IUserSrvFromDi(di)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if IsAuthenticatedUser(r.Context()) {
				next.ServeHTTP(w, r)
				return
			}
			tokenType, token, err := getAuthTokenFromHeader(r, true)
			if ae := utils.ApiErrorFromErr(err); ae != nil {
				utils.RenderJson(r, w, ae.StatusCode, ae)
				return
			}

			var user AuthenticatedUser
			switch tokenType {
			case headerApiKey:
				u, err := srv.GetUserFromApiKey(r.Context(), token)
				if err != nil {
					ae := utils.NewAuthenticationError("")
					utils.RenderJson(r, w, ae.StatusCode, ae)
					return
				}
				user = AuthenticatedUser{u.GetID(), u.GetRoles()}
			default:
				ae := utils.NewAuthenticationError("")
				utils.RenderJson(r, w, ae.StatusCode, ae)
				return
			}
			ctx := context.WithValue(r.Context(), utils.Authenticated, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ValidateAccessToken(signingKey string, accessToken string) (AuthenticatedUser, error) {
	token, err := jwt.ParseWithClaims(accessToken, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return AuthenticatedUser{}, err
	}
	payload, ok := token.Claims.(*JwtClaims)
	if ok && token.Valid {
		return payload.User, nil
	}

	return payload.User, errors.New("invalid token")
}

func getAuthTokenFromHeader(r *http.Request, withApiKey bool) (string, string, error) {
	if withApiKey {
		xapikey := r.Header.Get(headerApiKey)
		if len(xapikey) > 0 {
			return headerApiKey, xapikey, nil
		}
	} else {
		bearer := r.Header.Get(headerAuthorization)
		size := len(headerPrefixBearer) + 1
		if len(bearer) > size && strings.ToUpper(bearer[0:size-1]) == headerPrefixBearer {
			return headerPrefixBearer, bearer[size:], nil
		}
	}
	ae := utils.NewAuthorizationError(http.StatusText(http.StatusUnauthorized))
	return "", "", &ae
}
