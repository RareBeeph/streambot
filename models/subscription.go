package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Subscription struct {
	gorm.Model
	GameID    string
	GameName  string
	Filter    string
	GuildID   string
	ChannelID string

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
