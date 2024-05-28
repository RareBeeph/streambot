package models

import (
	"fmt"

	"gorm.io/gorm"
)

type subHealthRegistry struct {
	Healthy  int
	Stale    int
	Orphaned int
}

// TODO: Load this from config
var SubHealths = subHealthRegistry{
	Healthy:  0,
	Stale:    5,
	Orphaned: 20,
}

type Subscription struct {
	gorm.Model
	GameID    string `gorm:"uniqueIndex:idx_channel_game_filter"`
	GameName  string
	Filter    string `gorm:"uniqueIndex:idx_channel_game_filter"`
	GuildID   string
	ChannelID string `gorm:"uniqueIndex:idx_channel_game_filter"`

	TimesFailed int

	Messages []Message
}

func (s *Subscription) String() string {
	out := fmt.Sprintf("Game: `%s`", s.GameName)
	if s.Filter != "" {
		out += fmt.Sprintf(" | Filter: `%s`", s.Filter)
	}
	if s.TimesFailed >= SubHealths.Stale {
		out += " (Deactivated)"
	}

	return out
}

func init() {
	All = append(All, Subscription{})
}
