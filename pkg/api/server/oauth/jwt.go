package oauth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/status/pkg/defaults"
	"github.com/sdslabs/status/pkg/utils"
)

const (
	errInvalidToken = "UNAUTHORIZED_INVALID_TOKEN"
)

var (
	jwtSecret = []byte(utils.StatusConf.JWTSecret)
)

type Claims struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	jwt.StandardClaims
}

func CurrentUserFromCtx(ctx *gin.Context) (*Claims, bool) {
	claims, ok := ctx.Get(defaults.JWTContextKey)
	if !ok {
		return nil, false
	}
	return claims.(*Claims), true
}

func newToken(id uint, email, name string) (string, error) {
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
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func getClaimsFromToken(token string) (*Claims, error) {
	// when returning ErrInvalidToken, it means the status
	// of error is unauthorized rather than bad request
	c := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, c, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New(errInvalidToken)
		}
		return nil, err
	}
	if !tkn.Valid {
		return nil, errors.New(errInvalidToken)
	}
	return c, nil
}

// VerifyJWTMiddleware is used to authenticate any request for user
// Sets the ctx 'currentUser' value to the user email
func VerifyJWTMiddleware(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Missing Authorization Header",
		})
		return
	}
	authTypeLen := len(defaults.JWTAuthType)
	if authHeader[:authTypeLen] != defaults.JWTAuthType {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Missing '%s' Authorization type", defaults.JWTAuthType),
		})
		return
	}
	claims, err := getClaimsFromToken(authHeader[authTypeLen+1:])
	if err != nil {
		if err.Error() == errInvalidToken {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.Set(defaults.JWTContextKey, claims)
	ctx.Next()
}
