package models

import "gorm.io/gorm"

type Stream struct {
	gorm.Model
	UserID   string `copier:"must,nopanic"`
	UserName string `copier:"must,nopanic"`
	Title    string `copier:"must,nopanic"`
	GameID   string `copier:"must,nopanic"`
}

func init() {
	All = append(All, Stream{})
}
