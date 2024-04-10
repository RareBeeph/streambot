package models

import "gorm.io/gorm"

type Stream struct {
	gorm.Model
	UserID   string
	UserName string
	Title    string
}
