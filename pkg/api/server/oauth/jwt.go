package oauth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/status/pkg/defaults"
	"github.com/sdslabs/status/pkg/utils"
)

const (
	errExpiredTokenPrefix = "token is expired by"
	errExpiredToken       = "token expired"

	contextBucketKey = "currentUser"
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
	claims, ok := ctx.Get(contextBucketKey)
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
		return "", errors.New("Missing Authorization Header")
	}
	authTypeLen := len(defaults.JWTAuthType)
	if authHeader[:authTypeLen] != defaults.JWTAuthType {
		return "", fmt.Errorf("Missing '%s' Authorization type", defaults.JWTAuthType)
	}
	return authHeader[authTypeLen+1:], nil
}

// RefreshTokenHandler refreshes the token, given it claims it under the
// max refresh time.
func RefreshTokenHandler(ctx *gin.Context) {
	token, err := getTokenFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	claims, err := getClaimsFromToken(token)
	// Since we're refreshing token, we don't care if the previous token
	// claims are valid or invalid
	if err != nil && err.Error() != errExpiredToken {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// we need to check if the time since previous expired token is under
	// the max refresh time
	if time.Now().Unix() > claims.ExpiresAt+int64(defaults.JWTRefreshInterval) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Time exceeds max refresh time, login again",
		})
		return
	}

	refreshedToken, err := newToken(claims.ID, claims.Email, claims.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token":      refreshedToken,
		"expires_in": defaults.JWTExpireInterval / time.Second,
		"user_id":    claims.ID,
		"user_email": claims.Email,
		"user_name":  claims.Name,
	})
}

// VerifyJWTMiddleware is used to authenticate any request for user
// Sets the ctx 'currentUser' value to the user email
func VerifyJWTMiddleware(ctx *gin.Context) {
	token, err := getTokenFromCtx(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	claims, err := getClaimsFromToken(token)
	if err != nil {
		if err.Error() == errExpiredToken {
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
	ctx.Set(contextBucketKey, claims)
	ctx.Next()
}
