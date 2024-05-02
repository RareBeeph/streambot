package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	MessageID string `gorm:"primaryKey"`
	UserID    string
	PostOrder int

	SubscriptionID uint // foreign key
}

func init() {
	All = append(All, Message{})
}
