package models

import "gorm.io/gorm"

type BlacklistEntry struct {
	gorm.Model
	UserID    string `gorm:"UniqueIndex:idx_channel_user"`
	UserLogin string
	ChannelID string `gorm:"UniqueIndex:idx_channel_user"`
}

func init() {
	All = append(All, BlacklistEntry{})
}
