package oauth

import (
	"context"
	"github.com/mqqff/savebite-be/internal/domain/env"
	"github.com/mqqff/savebite-be/pkg/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauth2api "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type UserInfo struct {
	Email string
	Name  string
}

type OAuthItf interface {
	GenerateAuthLink(state string) string
	ExchangeToken(code string) (*oauth2.Token, error)
	GetUserInfo(token *oauth2.Token) (UserInfo, error)
}

type OAuthStruct struct {
	config *oauth2.Config
}

var OAuth = getOAuth()

func getOAuth() OAuthItf {
	config := &oauth2.Config{
		ClientID:     env.AppEnv.GoogleClientID,
		ClientSecret: env.AppEnv.GoogleClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  env.AppEnv.GoogleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid",
		},
	}

	return &OAuthStruct{config}
}

func (o *OAuthStruct) GenerateAuthLink(state string) string {
	return o.config.AuthCodeURL(state)
}

func (o *OAuthStruct) ExchangeToken(code string) (*oauth2.Token, error) {
	return o.config.Exchange(context.Background(), code)
}

func (o *OAuthStruct) GetUserInfo(token *oauth2.Token) (UserInfo, error) {
	client := o.config.Client(context.Background(), token)
	oauth2Service, err := oauth2api.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[OAuth][GetUserInfo] failed to create new oauth2 service")
		return UserInfo{}, err
	}

	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[OAuth][GetUserInfo] failed to get user info")
		return UserInfo{}, err
	}

	return UserInfo{
		Email: userInfo.Email,
		Name:  userInfo.Name,
	}, nil
}
