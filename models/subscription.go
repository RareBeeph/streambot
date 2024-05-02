package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Subscription struct {
	gorm.Model
	GameID    string `gorm:"uniqueIndex:idx_channel_game_filter"`
	GameName  string
	Filter    string `gorm:"uniqueIndex:idx_channel_game_filter"`
	GuildID   string
	ChannelID string `gorm:"uniqueIndex:idx_channel_game_filter"`

	Messages []Message
}

func (s *Subscription) String() string {
	out := fmt.Sprintf("Game: `%s`", s.GameName)
	if s.Filter != "" {
		out += fmt.Sprintf(" | Filter: `%s`", s.Filter)
	}

	return out
}

func init() {
	All = append(All, Subscription{})
}
