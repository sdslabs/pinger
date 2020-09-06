// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package oauther

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sdslabs/pinger/pkg/appcontext"
)

// This map stores all the oauthers. The only way to add a new oauther in
// this map is to use the `Register` method.
var oauthers = map[string]newFunc{}

// newFunc is an alias for the function that can create a new oauther.
type newFunc = func() Oauther

// Register adds a new oauther to the package. This does not throw an
// error, rather panics if the oauther with the same name is already
// registered, hence an oauther should be registered inside the init method
// of the package.
func Register(name string, fn newFunc) {
	if _, ok := oauthers[name]; ok {
		panic(fmt.Errorf("oauther with same name already exists: %s", name))
	}

	oauthers[name] = fn
}

// Oauther represents any service that provides third-party authentication
// through the Oauth2 protocol.
type Oauther interface {
	// Provision enables any initialization of variables or any other config
	// required through the oauth options.
	Provision(*appcontext.Context, Provider) error

	// AuthURL returns the URL where the user authenticates the application.
	// It takes the state parameter and generates the authorization URL with
	// other required parameters.
	AuthURL(state string) string

	// FetchUser fetches the user information from the provider API. The info
	// should be in valid JSON format. The method should return a non-zero status
	// code in case of an error to specify what kind of error it was during the
	// API call, example, 400 is bad request and 500 is an internal server error.
	FetchUser(ctx context.Context, code string) (json.RawMessage, int, error)
}

// Initialize adds the required routes like the login and redirect routes to
// the router for all the given oauthers.
func Initialize(ctx *appcontext.Context, provider Provider, opts *Opts) error {
	newOauther, ok := oauthers[provider.GetProvider()]
	if !ok {
		return fmt.Errorf("provider %q does not exist", provider.GetProvider())
	}

	oauther := newOauther()

	if err := oauther.Provision(ctx, provider); err != nil {
		return err
	}

	groupRoute := fmt.Sprintf("/%s", provider.GetProvider())
	oautherGroup := opts.Router.Group(groupRoute)
	oautherGroup.GET("/login", loginRoute(oauther, opts))
	oautherGroup.GET("/redirect", redirectRoute(oauther, opts))

	return nil
}

// loginRoute returns the login route for the provider.
func loginRoute(oauther Oauther, opts *Opts) func(*gin.Context) {
	return func(ctx *gin.Context) {
		state := randomToken()
		authURL := oauther.AuthURL(state)
		ctx.PureJSON(http.StatusOK, opts.LoginResponse(authURL))
	}
}

// redirectRoute returns the redirect route for the provider.
func redirectRoute(oauther Oauther, opts *Opts) func(*gin.Context) {
	return func(ctx *gin.Context) {
		code := ctx.Query("code")
		user, sc, err := oauther.FetchUser(ctx.Request.Context(), code)
		if err != nil {
			ctx.PureJSON(sc, opts.ErrorResponse(err))
			return
		}

		v, sc, err := opts.OnUser(user)
		if err != nil {
			ctx.PureJSON(sc, opts.ErrorResponse(err))
			return
		}

		ctx.PureJSON(http.StatusOK, opts.RedirectResponse(v))
	}
}

// randomToken generates a random string.
func randomToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}
