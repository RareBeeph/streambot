package config

import "github.com/spf13/viper"

type Config struct {
	Token       string
	GuildID     string
	DatabaseURL string

	Twitch struct {
		ClientID     string
		ClientSecret string
	}
}

var current *Config

func Load() *Config {
	if current != nil {
		return current
	}

	current = &Config{}
	viper.Unmarshal(current)
	return current
}
