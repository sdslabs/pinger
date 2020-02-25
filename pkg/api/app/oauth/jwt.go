package oauth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/sdslabs/status/pkg/api/response"
	"github.com/sdslabs/status/pkg/defaults"
)

const (
	errExpiredTokenPrefix = "token is expired by"
	errExpiredToken       = "token expired"

	contextBucketKey = "currentUser"
)

// Claims represent the claims of an authorization JSON Web Token.
// Includes the user ID (Primary key in DB), Email and Name.
type Claims struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	jwt.StandardClaims
}

// CurrentUserFromCtx gets the current user for `currentUser` context bucket
// in the form of JWT Claims.
// Returns false when user is not found in the context.
func CurrentUserFromCtx(ctx *gin.Context) (*Claims, bool) {
	claims, ok := ctx.Get(contextBucketKey)
	if !ok {
		return nil, false
	}
	return claims.(*Claims), true
}

// GetRefreshTokenHandler returns the handler which refreshes the token,
// given it claims it under the max refresh time.
func GetRefreshTokenHandler(jwtSecret []byte) ginHandler {
	return func(ctx *gin.Context) {
		token, err := getTokenFromCtx(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.HTTPError{
				Error: err.Error(),
			})
			return
		}
		claims, err := getClaimsFromToken(token, jwtSecret)
		// Since we're refreshing token, we don't care if the previous token
		// claims are valid or invalid
		if err != nil && err.Error() != errExpiredToken {
			ctx.JSON(http.StatusBadRequest, response.HTTPError{
				Error: err.Error(),
			})
			return
		}
		// we need to check if the time since previous expired token is under
		// the max refresh time
		if time.Now().Unix() > claims.ExpiresAt+int64(defaults.JWTRefreshInterval) {
			ctx.JSON(http.StatusUnauthorized, response.HTTPError{
				Error: "Time exceeds max refresh time, login again",
			})
			return
		}

		refreshedToken, err := newToken(claims.ID, claims.Email, claims.Name, jwtSecret)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.HTTPError{
				Error: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, response.HTTPRefreshToken{
			Token:     refreshedToken,
			ExpiresIn: defaults.JWTExpireInterval / time.Second,
		})
	}
}

// GetJWTVerficationMiddleware returns a middleware that is used to authenticate
// any request for user. Sets the ctx 'currentUser' value to the user email.
func GetJWTVerficationMiddleware(jwtSecret []byte) ginHandler {
	return func(ctx *gin.Context) {
		token, err := getTokenFromCtx(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response.HTTPError{
				Error: err.Error(),
			})
			return
		}
		claims, err := getClaimsFromToken(token, jwtSecret)
		if err != nil {
			if err.Error() == errExpiredToken {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.HTTPError{
					Error: err.Error(),
				})
				return
			}
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response.HTTPError{
				Error: err.Error(),
			})
			return
		}
		ctx.Set(contextBucketKey, claims)
		ctx.Next()
	}
}

func newToken(id uint, email, name string, jwtSecret []byte) (string, error) {
	expirationTime := time.Now().Add(defaults.JWTExpireInterval)
	c := &Claims{
		ID:    id,
		Email: email,
		Name:  name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(jwtSecret)
}

func getClaimsFromToken(token string, jwtSecret []byte) (*Claims, error) {
	// when returning ErrInvalidToken, it means the status
	// of error is unauthorized rather than bad request
	c := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, c, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), errExpiredTokenPrefix) {
			return c, errors.New(errExpiredToken)
		}
		return c, err
	}
	if !tkn.Valid {
		return c, errors.New("invalid token")
	}
	return c, nil
}

func getTokenFromCtx(ctx *gin.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("missing Authorization Header")
	}
	authTypeLen := len(defaults.JWTAuthType)
	if authHeader[:authTypeLen] != defaults.JWTAuthType {
		return "", fmt.Errorf("missing '%s' Authorization type", defaults.JWTAuthType)
	}
	return authHeader[authTypeLen+1:], nil
}
