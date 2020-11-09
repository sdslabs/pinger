package oauther

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

// Provider configures the oauther.
type Provider interface {
	GetProvider() string // Provider name

	GetClientID() string     // ClientID of provider
	GetClientSecret() string // ClientSecret of provider
	GetRedirectURL() string  // RedirectURL of provider
	GetScopes() []string     // Scopes of provider
}

// Opts configures the gin Router to enable authentication.
type Opts struct {
	// Router group to add the routes to.
	//
	// Routes are added to the router as:
	//
	// 	/:oauther/login
	// 	/:oauther/redirect
	//
	Router gin.IRouter

	// When the user is fetched from the provider, what action to do. Non-zero
	// status code should be returned in case of non-nil error. The interface
	// returned is what is set in session for the user.
	OnUser func(user json.RawMessage) (interface{}, int, error)

	// Response to send on getting the auth URL from the oauth provider.
	LoginResponse func(url string) interface{}

	// Response to send on getting the user information success and processing
	// it (`OnUser`ing) without any error.
	RedirectResponse func(user interface{}) interface{}

	// ErrorResponse when an error is received.
	ErrorResponse func(error) interface{}
}
