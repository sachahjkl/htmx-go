package model

import "gorm.io/gorm"

type Settings struct {
	gorm.Model
	Title  string `gorm:"not null"`
	Done   bool
	UserID uint
}
