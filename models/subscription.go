package models

import "gorm.io/gorm"

type Subscription struct {
	gorm.Model
	GameID    string
	GameName  string
	Filter    string
	GuildID   string
	ChannelID string

	Messages []Message
}

func init() {
	AllModels = append(AllModels, Subscription{})
}
