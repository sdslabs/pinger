// Package application contains the router for status web app.
package application

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sdslabs/status/pkg/api/application/oauth"
	"github.com/sdslabs/status/pkg/api/application/providers"
	"github.com/sdslabs/status/pkg/config"
	"github.com/sdslabs/status/pkg/database"
)

// ErrInvalidConfigDeploy is returned from `Serve` when application.deploy is set to false in `config.yml`.
var ErrInvalidConfigDeploy = fmt.Errorf("application.deploy is set false in config, cannot start server without postgres")

func getProvider(providerType string) (oauth.Provider, error) {
	switch oauth.ProviderType(providerType) {
	case providers.Google.Type():
		return providers.Google, nil
	case providers.Github.Type():
		return providers.Github, nil
	default:
		return nil, fmt.Errorf("invalid oauth provider '%s'", providerType)
	}
}

func getProvidersFromConf(conf *config.StatusConfig) ([]oauth.Provider, error) {
	oauthConf := conf.Application.Oauth
	p := []oauth.Provider{}
	for key := range oauthConf {
		provider, err := getProvider(key)
		if err != nil {
			return nil, err
		}
		p = append(p, provider)
	}
	return p, nil
}

func getRouter(conf *config.StatusConfig) (*gin.Engine, error) {
	router := gin.Default()

	oauthProviders, err := getProvidersFromConf(conf)
	if err != nil {
		return nil, err
	}
	oauthRouter := router.Group("/oauth")
	if err := oauth.Initialize(oauthRouter, conf, oauthProviders...); err != nil {
		return nil, err
	}

	apiRouter := router.Group("/api")
	apiRouter.Use(oauth.GetJWTVerficationMiddleware(conf.Secret()))
	apiRouter.GET("/test", func(ctx *gin.Context) {
		currentUser, ok := oauth.CurrentUserFromCtx(ctx)

		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Cannot find user from token",
			})
		}
		ctx.JSON(http.StatusOK, gin.H{
			"name":  currentUser.Name,
			"email": currentUser.Email,
		})
	})

	return router, nil
}

// Serve starts the HTTP server on default port "8080".
func Serve(conf *config.StatusConfig, port int) error {
	if err := database.SetupDB(conf); err != nil {
		return fmt.Errorf("error setting up postgresql: %s", err.Error())
	}

	r, err := getRouter(conf)
	if err != nil {
		return err
	}

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	return r.Run(addr)
}
