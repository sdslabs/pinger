package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/status/api/oauth"
)

// NewRouter returns the router for the main API service
func NewRouter() *gin.Engine {
	router := gin.Default()
	oauth.SetupGoogleOAuth()

	apiRouter := router.Group("/api")
	{
		// Oauth routes
		oauthRouter := apiRouter.Group("/oauth")
		{
			oauthRouter.GET("/google", oauth.HandleGoogleLogin) // sends the url for login
			oauthRouter.GET("/google/redirect", oauth.HandleGoogleRedirect)
		}
	}

	return router
}
