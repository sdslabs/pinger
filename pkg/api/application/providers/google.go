package providers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/sdslabs/status/pkg/api/application/oauth"
	"github.com/sdslabs/status/pkg/config"
)

const (
	googleUserInfoEndpoint = "https://www.googleapis.com/oauth2/v3/userinfo"
	googleProviderType     = "google"
)

// Google oauth `Provider`.
var Google = &googleProvider{}

type googleProvider struct {
	config *oauth2.Config
	state  string
}

func (p *googleProvider) Type() oauth.ProviderType {
	return oauth.ProviderType(googleProviderType)
}

func (p *googleProvider) Setup(conf *config.OauthProviderConfig) error {
	p.config = &oauth2.Config{
		ClientID:     conf.ClientID,
		ClientSecret: conf.ClientSecret,
		RedirectURL:  conf.RedirectURL,
		Scopes:       conf.Scopes,
		Endpoint:     google.Endpoint,
	}
	return nil
}

func (p *googleProvider) GetLoginURL() string {
	p.state = randToken()
	return p.config.AuthCodeURL(p.state)
}

func (p *googleProvider) GetUser(ctx *gin.Context) (*oauth.User, int, error) {
	token, err := p.config.Exchange(context.TODO(), ctx.Query("code"))
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	client := p.config.Client(context.TODO(), token)
	info, err := client.Get(googleUserInfoEndpoint)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	defer info.Body.Close() //nolint:errcheck

	data, err := ioutil.ReadAll(info.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	u := new(oauth.User)
	if err := json.Unmarshal(data, u); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return u, http.StatusOK, nil
}
