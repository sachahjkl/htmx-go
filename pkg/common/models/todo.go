package models

import "gorm.io/gorm"

type Todo struct {
	gorm.Model
	Title string `json:"title" form:"todo-title"`
	Done  bool   `json:"done" form:"todo-done"`
}