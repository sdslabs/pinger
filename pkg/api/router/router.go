// Package router contains the router for status web app.
package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sdslabs/status/pkg/api/router/oauth"
	"github.com/sdslabs/status/pkg/api/router/providers"
	"github.com/sdslabs/status/pkg/database"
	"github.com/sdslabs/status/pkg/utils"
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

func getProvidersFromConf() ([]oauth.Provider, error) {
	oauthConf := utils.StatusConf.Oauth
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

func getRouter() (*gin.Engine, error) {
	router := gin.Default()

	oauthProviders, err := getProvidersFromConf()
	if err != nil {
		return nil, err
	}
	oauthRouter := router.Group("/oauth")
	if err := oauth.Initialize(oauthRouter, oauthProviders...); err != nil {
		return nil, err
	}

	apiRouter := router.Group("/api")
	apiRouter.Use(oauth.JWTVerficationMiddleware)
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
func Serve() error {
	if err := database.SetupDB(); err != nil {
		return err
	}

	r, err := getRouter()
	if err != nil {
		return err
	}

	return r.Run()
}
