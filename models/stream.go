package models

import "gorm.io/gorm"

type Stream struct {
	gorm.Model
	UserID      string `copier:"must,nopanic" gorm:"UniqueIndex:idx_user"`
	UserName    string `copier:"must,nopanic"`
	Title       string `copier:"must,nopanic"`
	GameID      string `copier:"must,nopanic"`
	Language    string `copier:"must,nopanic"`
	ViewerCount string `copier:"must,nopanic"`
}

func init() {
	All = append(All, Stream{})
}
