package oauth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/sdslabs/status/pkg/api/response"
	"github.com/sdslabs/status/pkg/config"
	"github.com/sdslabs/status/pkg/database"
	"github.com/sdslabs/status/pkg/defaults"
)

var (
	errNilProvider = fmt.Errorf("provider cannot be nil")

	// Contains all the providers.
	providerRegister = map[ProviderType]Provider{}
)

// User represents the user details required from the oauth provider.
type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ProviderType for different OAuth Providers.
type ProviderType string

// Provider is an interface that represents any provider
// through which we can authenticate a user on the application.
type Provider interface {

	// Type returns the type of provider.
	Type() ProviderType

	// Setup enables us to initialize any variables or setup requirements.
	Setup(*config.OauthProviderConfig) error

	// GetLoginURL returns the URL which redirects user to the providers login page.
	GetLoginURL() string

	// GetUser gets the user after requesting the OAuth provider.
	// Returns the user, status code and error if any.
	GetUser(*gin.Context) (*User, int, error)
}

// Setup all the providers. This adds the refresh token route for JWT as well as
// login and redirect routes for all the providers.
func Setup(oauthRouter gin.IRoutes, conf *config.AppConfig) error {
	jwtSecret := conf.Secret()

	// Add refresh route.
	oauthRouter.GET("/refresh", GetRefreshTokenHandler(jwtSecret))

	for typ, provider := range providerRegister {
		if provider == nil {
			continue
		}

		typStr := string(typ)

		oauthConfig, ok := conf.Oauth[typStr]
		if !ok {
			return fmt.Errorf("could not find conf for %s OAuth provider", typStr)
		}
		if err := provider.Setup(&oauthConfig); err != nil {
			return fmt.Errorf(
				"error while setting up %s OAuth provider: %s",
				typStr,
				err.Error())
		}

		loginRoute := fmt.Sprintf("/%s", typStr)
		redirectRoute := fmt.Sprintf("/%s/redirect", typStr)

		// Add login route.
		oauthRouter.GET(loginRoute, loginHandler(provider))
		// Add redirect route.
		oauthRouter.GET(redirectRoute, redirectHandler(provider, jwtSecret))
	}

	return nil
}

// AddProvider adds a new provider to the register.
// This should be run before `Setup`.
//
// Returns an error when provider exists. Use `UpdateProvider` instead.
func AddProvider(provider Provider) error {
	if provider == nil {
		return errNilProvider
	}
	typ := provider.Type()
	if _, ok := providerRegister[typ]; ok {
		return fmt.Errorf("provider %s already exists", string(typ))
	}
	providerRegister[typ] = provider
	return nil
}

// UpdateProvider updates the provider for given type.
// This should be run before `Setup`.
//
// Returns error if provider doesn't exist. Use `AddProvider` instead.
func UpdateProvider(provider Provider) error {
	if provider == nil {
		return errNilProvider
	}
	typ := provider.Type()
	if _, ok := providerRegister[typ]; !ok {
		return fmt.Errorf("provider %s doesn't exist", string(typ))
	}
	providerRegister[typ] = provider
	return nil
}

// Initialize is a shorthand for adding multiple routers to the group and setting them up.
func Initialize(oauthRouter gin.IRoutes, conf *config.AppConfig, providers ...Provider) error {
	for _, provider := range providers {
		if err := AddProvider(provider); err != nil {
			return fmt.Errorf("error while adding %s provider: %s", string(provider.Type()), err.Error())
		}
	}
	if err := Setup(oauthRouter, conf); err != nil {
		return err
	}
	return nil
}

type ginHandler = func(*gin.Context)

func loginHandler(provider Provider) ginHandler {
	return func(ctx *gin.Context) {
		ctx.JSON(200, response.HTTPLogin{
			LoginURL: provider.GetLoginURL(),
		})
	}
}

func redirectHandler(provider Provider, jwtSecret []byte) ginHandler {
	return func(ctx *gin.Context) {
		u, code, err := provider.GetUser(ctx)
		if err != nil {
			ctx.JSON(code, response.HTTPError{
				Error: err.Error(),
			})
			return
		}

		createdUser, err := database.CreateUser(&database.User{
			Email: u.Email,
			Name:  u.Name,
		})
		if err != nil {
			logrus.WithError(err).Errorln("cannot create user in database")
			ctx.JSON(http.StatusInternalServerError, response.HTTPInternalServerError)
			return
		}

		jwt, err := newToken(createdUser.ID, createdUser.Email, jwtSecret)
		if err != nil {
			logrus.WithError(err).Errorln("cannot create jwt")
			ctx.JSON(http.StatusInternalServerError, response.HTTPInternalServerError)
			return
		}

		ctx.JSON(http.StatusOK, response.HTTPAuthorization{
			Token:     jwt,
			ExpiresIn: defaults.JWTExpireInterval / time.Second,
			UserID:    createdUser.ID,
			UserEmail: createdUser.Email,
			UserName:  createdUser.Name,
		})
	}
}
