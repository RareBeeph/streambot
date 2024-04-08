package models

import "gorm.io/gorm"

type RegisteredCommand struct {
	gorm.Model
	ID      string `gorm:"primaryKey"`
	GuildID string
}

func init() {
	AllModels = append(AllModels, RegisteredCommand{})
}
