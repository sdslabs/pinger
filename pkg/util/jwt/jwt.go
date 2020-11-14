package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// errExpiredTokenPrefix is the prefix of the error when token is expired.
const errExpiredTokenPrefix = "token is expired by"

// errExpiredToken is the error returned when token is expired.
var errExpiredToken = errors.New("token expired")

// claims define the JWT claims.
type claims struct {
	jwt.StandardClaims
	Values interface{} `json:"val,omitempty"`
}

// JWT defines and creates the token and it's properties.
type JWT struct {
	// ExpirationInterval of the token.
	ExpirationInterval time.Duration

	// RefreshInterval of the token.
	RefreshInterval time.Duration

	// Secret is used to sign the token.
	Secret []byte

	// AuthType defines the authentication type, i.e., the prefix before the
	// token in authorization header.
	AuthType string
}

// NewToken creates a new JWT.
func (t *JWT) NewToken(values interface{}) (string, error) {
	expirationTime := time.Now().Add(t.ExpirationInterval)
	cl := claims{
		Values: values,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &cl)
	return token.SignedString(t.Secret)
}

// RefreshToken returns the refresh token for an expired token.
func (t *JWT) RefreshToken(token string) (refreshToken string, statusCode int, _ error) {
	cl, err := t.getClaimsFromToken(token)
	if err != nil && !errors.Is(err, errExpiredToken) {
		return "", http.StatusBadRequest, err
	}

	if time.Now().Unix() > cl.ExpiresAt+int64(t.RefreshInterval) {
		return "", http.StatusUnauthorized, errors.New("time exceeds max refresh time, login again")
	}

	refreshToken, err = t.NewToken(cl.Values)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	return refreshToken, http.StatusOK, nil
}

// VerifyToken verifies if the token is valid and returns the values for
// its claims.
func (t *JWT) VerifyToken(token string) (values interface{}, statusCode int, _ error) {
	cl, err := t.getClaimsFromToken(token)
	if err != nil {
		if errors.Is(err, errExpiredToken) {
			return nil, http.StatusUnauthorized, err
		}

		return nil, http.StatusBadRequest, err
	}

	return cl.Values, http.StatusOK, nil
}

// getClaimsFromToken returns the values after verifying if the token
// is valid or not.
func (t *JWT) getClaimsFromToken(token string) (claims, error) {
	// when returning ErrInvalidToken, it means the status of error is
	// unauthorized rather than bad request
	c := claims{}
	tkn, err := jwt.ParseWithClaims(token, &c, func(token *jwt.Token) (interface{}, error) {
		return t.Secret, nil
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), errExpiredTokenPrefix) {
			return c, errExpiredToken
		}

		return c, err
	}

	if !tkn.Valid {
		return c, errors.New("invalid token")
	}

	return c, nil
}

// GetTokenFromHeader returns the token from gin Context.
func (t *JWT) GetTokenFromHeader(ctx *gin.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("missing Authorization Header")
	}

	authTypeLen := len(t.AuthType)
	if authHeader[:authTypeLen] != t.AuthType {
		return "", fmt.Errorf("missing '%s' authorization type", t.AuthType)
	}

	return authHeader[authTypeLen+1:], nil
}
