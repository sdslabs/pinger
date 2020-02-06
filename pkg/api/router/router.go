// Package router contains the router for status web app.
package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sdslabs/status/pkg/api/router/oauth"
	"github.com/sdslabs/status/pkg/api/router/providers"
	"github.com/sdslabs/status/pkg/database"
)

func getRouter() (*gin.Engine, error) {
	router := gin.Default()

	oauthRouter := router.Group("/oauth")
	if err := oauth.Initialize(oauthRouter, providers.Google); err != nil {
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
