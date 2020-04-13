// Package app contains the router for status web app.
package app

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/sdslabs/status/pkg/api/app/oauth"
	"github.com/sdslabs/status/pkg/api/app/providers"
	"github.com/sdslabs/status/pkg/api/handlers"
	"github.com/sdslabs/status/pkg/config"
	"github.com/sdslabs/status/pkg/database"
)

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

func getProvidersFromConf(conf *config.AppConfig) ([]oauth.Provider, error) {
	oauthConf := conf.Oauth
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

func getRouter(conf *config.AppConfig) (*gin.Engine, error) {
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

	apiRouter.GET("/user/:id", handlers.GetUser)
	apiRouter.PATCH("/user/:id", handlers.UpdateUser)
	apiRouter.DELETE("/user", handlers.DeleteUser)

	return router, nil
}

// Serve starts the HTTP server on default port "8080".
func Serve(conf *config.AppConfig) error {
	if err := database.SetupDB(conf); err != nil {
		return fmt.Errorf("error setting up postgresql: %s", err.Error())
	}

	r, err := getRouter(conf)
	if err != nil {
		return err
	}

	addr := fmt.Sprintf("0.0.0.0:%d", conf.Port)
	return r.Run(addr)
}
