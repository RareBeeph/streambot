package models

import "gorm.io/gorm"

type RegisteredCommand struct {
	gorm.Model
	ID string `gorm:"primaryKey"`
}

func init() {
	AllModels = append(AllModels, RegisteredCommand{})
}
