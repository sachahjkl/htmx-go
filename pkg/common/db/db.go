package db

import (
	"log"
	"sachahjkl/htmx_go/pkg/common/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Init(url string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&models.Todo{})

	return db
}
