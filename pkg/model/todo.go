package model

import (
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Title  string
	Done   bool
	UserID uint
}

func AddTodo(db *gorm.DB, title string, done bool, userId uint) (*Todo, error) {
	todo := Todo{
		Title:  title,
		Done:   done,
		UserID: userId,
	}

	// insert new todo
	err := db.Create(&todo).Error
	return &todo, err
}

func DeleteTodo(db *gorm.DB, id uint, userId uint) error {
	return db.Delete(&Todo{UserID: userId}, id).Error
}

func ToggleTodo(db *gorm.DB, id uint, userId uint) (*Todo, error) {
	todo := Todo{UserID: userId}
	err := db.First(&todo, id).Error
	if err != nil {
		return nil, err
	}

	// toggle the todo
	todo.Done = !todo.Done

	err = db.Save(&todo).Error
	return &todo, err
}

func AllTodos(db *gorm.DB, userId uint) (*[]Todo, error) {
	var todos []Todo
	err := db.Find(&todos).Error
	return &todos, err
}
