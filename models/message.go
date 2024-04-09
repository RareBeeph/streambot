package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	MessageID string `gorm:"primaryKey"`

	SubscriptionID uint // foreign key
}

func init() {
	AllModels = append(AllModels, Message{})
}
