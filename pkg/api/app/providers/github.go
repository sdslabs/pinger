package providers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"

	"github.com/sdslabs/status/pkg/api/app/oauth"
	"github.com/sdslabs/status/pkg/config"
)

const (
	githubUserInfoEndpoint = "https://api.github.com/user"
	githubProviderType     = "github"
)

// Github oauth `Provider`.
var Github = &githubProvider{}

type githubProvider struct {
	config *oauth2.Config
	state  string
}

func (p *githubProvider) Type() oauth.ProviderType {
	return oauth.ProviderType(githubProviderType)
}

func (p *githubProvider) Setup(conf *config.OauthProviderConfig) error {
	p.config = &oauth2.Config{
		ClientID:     conf.ClientID,
		ClientSecret: conf.ClientSecret,
		RedirectURL:  conf.RedirectURL,
		Scopes:       conf.Scopes,
		Endpoint:     github.Endpoint,
	}
	return nil
}

func (p *githubProvider) GetLoginURL() string {
	p.state = randToken()
	return p.config.AuthCodeURL(p.state)
}

func (p *githubProvider) GetUser(ctx *gin.Context) (*oauth.User, int, error) {
	token, err := p.config.Exchange(context.TODO(), ctx.Query("code"))
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	client := p.config.Client(context.TODO(), token)
	info, err := client.Get(githubUserInfoEndpoint)
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
