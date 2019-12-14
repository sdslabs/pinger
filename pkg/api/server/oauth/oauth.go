package oauth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/status/pkg/database"
	"github.com/sdslabs/status/pkg/defaults"
	"github.com/sdslabs/status/pkg/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// SetupGoogleOAuth initialises the oAuth conf
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

func getGoogleLoginURL(state string) string {
	state = randToken()
	return config.AuthCodeURL(state)
}

// HandleGoogleLogin sends the response as login url using google oauth
func HandleGoogleLogin(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"login_url": getGoogleLoginURL(state),
	})
}

// HandleGoogleRedirect when user allows for oAuth access
func HandleGoogleRedirect(ctx *gin.Context) {
	token, err := config.Exchange(oauth2.NoContext, ctx.Query("code"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	client := config.Client(oauth2.NoContext, token)
	info, err := client.Get(googleUserInfoEndpoint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer info.Body.Close()

	data, err := ioutil.ReadAll(info.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	u := new(user)
	err = json.Unmarshal(data, u)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	createdUser, err := database.DBConn.CreateUser(u.Email, u.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	jwt, err := newToken(createdUser.ID, createdUser.Email, createdUser.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token":      jwt,
		"expires_in": defaults.JWTExpireInterval / time.Second,
		"user_id":    createdUser.ID,
		"user_email": createdUser.Email,
		"user_name":  createdUser.Name,
	})
}