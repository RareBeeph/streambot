package twitch

import (
	"streambot/config"

	"github.com/nicklaw5/helix/v2"
)

var Client *helix.Client

func LoadConfig(conf *config.Config) {
	var err error
	Client, err = helix.NewClient(&helix.Options{
		ClientID:     conf.Twitch.ClientID,
		ClientSecret: conf.Twitch.ClientSecret,
	})
	if err != nil {
		return
	}

	resp, err := Client.RequestAppAccessToken([]string{"user:read:email"})
	if err != nil {
		return
	}

	Client.SetAppAccessToken(resp.Data.AccessToken)
}
