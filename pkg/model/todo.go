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
	if result := db.Create(&todo); result.Error != nil {
		return nil, result.Error
	}
	return &todo, nil
}

func DeleteTodo(db *gorm.DB, id uint, userId uint) error {
	result := db.Delete(&Todo{UserID: userId}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func ToggleTodo(db *gorm.DB, id uint, userId uint) (*Todo, error) {
	todo := Todo{UserID: userId}
	result := db.First(&todo, id)
	if result.Error != nil {
		return nil, result.Error
	}

	// toggle the todo
	todo.Done = !todo.Done

	result = db.Save(&todo)
	if result.Error != nil {
		return nil, result.Error
	}
	return &todo, nil
}

func AllTodos(db *gorm.DB, userId uint) (*[]Todo, error) {
	var todos []Todo
	result := db.Find(&todos)
	if result.Error != nil {
		return nil, result.Error
	}
	return &todos, nil
}
