package oauth

import (
	"errors"
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

type claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func newToken(email string) (string, error) {
	expirationTime := time.Now().Add(defaults.JWTExpireInterval)
	c := &claims{
		Email: email,
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

func getEmailFromToken(token string) (string, error) {
	// when returning ErrInvalidToken, it means the status
	// of error is unauthorized rather than bad request
	c := &claims{}
	tkn, err := jwt.ParseWithClaims(token, c, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", errors.New(errInvalidToken)
		}
		return "", err
	}
	if !tkn.Valid {
		return "", errors.New(errInvalidToken)
	}
	return c.Email, nil
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
	if authHeader[:5] != "Basic" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Missing 'Basic' Authorization type",
		})
		return
	}
	email, err := getEmailFromToken(authHeader[6:])
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
	ctx.Set("currentUser", email)
	ctx.Next()
}
