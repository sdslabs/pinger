package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/sdslabs/status/pkg/api/response"
	"github.com/sdslabs/status/pkg/database"
	"github.com/sdslabs/status/pkg/defaults"
	"github.com/sdslabs/status/pkg/utils"
)

const googleUserInfoEndpoint = "https://www.googleapis.com/oauth2/v3/userinfo"

var (
	config *oauth2.Config
	state  string
)

type user struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// SetupGoogleOAuth initializes the oAuth conf
func SetupGoogleOAuth() error {
	conf := utils.StatusConf.Oauth.Google

	config = &oauth2.Config{
		ClientID:     conf.ClientID,
		ClientSecret: conf.ClientSecret,
		RedirectURL:  conf.RedirectURL,
		Scopes:       conf.Scopes,
		Endpoint:     google.Endpoint,
	}
	return nil
}

// HandleGoogleLogin sends the response as login url using google oauth
func HandleGoogleLogin(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, response.HTTPLogin{
		LoginURL: getGoogleLoginURL(),
	})
}

// HandleGoogleRedirect when user allows for oAuth access
func HandleGoogleRedirect(ctx *gin.Context) {
	token, err := config.Exchange(context.TODO(), ctx.Query("code"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.HTTPError{
			Error: err.Error(),
		})
		return
	}

	client := config.Client(context.TODO(), token)
	info, err := client.Get(googleUserInfoEndpoint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.HTTPError{
			Error: err.Error(),
		})
		return
	}
	defer info.Body.Close() //nolint:errcheck

	data, err := ioutil.ReadAll(info.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.HTTPError{
			Error: err.Error(),
		})
		return
	}

	u := new(user)
	if err = json.Unmarshal(data, u); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.HTTPError{
			Error: err.Error(),
		})
		return
	}

	createdUser, err := database.DBConn.CreateUser(u.Email, u.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.HTTPError{
			Error: err.Error(),
		})
		return
	}

	jwt, err := newToken(createdUser.ID, createdUser.Email, createdUser.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.HTTPError{
			Error: err.Error(),
		})
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

func randToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}

func getGoogleLoginURL() string {
	state = randToken()
	return config.AuthCodeURL(state)
}
