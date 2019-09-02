package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/status/pkg/api/server/oauth"
)

// NewRouter returns the router for the main API service
func NewRouter() *gin.Engine {
	router := gin.Default()
	oauth.SetupGoogleOAuth()

	oauthRouter := router.Group("/oauth")
	{
		oauthRouter.GET("/google", oauth.HandleGoogleLogin) // sends the url for login
		oauthRouter.GET("/google/redirect", oauth.HandleGoogleRedirect)
	}

	apiRouter := router.Group("/api")
	apiRouter.Use(oauth.VerifyJWTMiddleware)
	{
		apiRouter.GET("/test", func(ctx *gin.Context) {
			currentUser, ok := ctx.Get("currentUser")
			if !ok {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": "Cannot find user from token",
				})
			}
			ctx.JSON(http.StatusOK, gin.H{
				"user": currentUser.(string),
			})
		})
	}

	return router
}
