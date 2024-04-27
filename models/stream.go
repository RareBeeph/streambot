package models

import "gorm.io/gorm"

type Stream struct {
	gorm.Model
	UserID   string `copier:"must"`
	UserName string `copier:"must"`
	Title    string `copier:"must"`
}

func init() {
	All = append(All, Stream{})
}
