// Package router contains the router for status web app.
package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sdslabs/status/pkg/api/router/oauth"
	"github.com/sdslabs/status/pkg/api/router/providers"
)

// NewRouter returns the router for the main API service
func NewRouter() (*gin.Engine, error) {
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
