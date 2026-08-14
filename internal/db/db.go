package db

import (
	"fmt"
	"log"
	"sachahjkl/htmx_go/internal/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func Init(url string) *gorm.DB {

	dsn := fmt.Sprintf(
		"file:%s?cache=private&_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)",
		url,
	)
	ormDb, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Migrating db @ %v", url)
	ormDb.AutoMigrate(&model.User{}, &model.Todo{})
	log.Printf("Finished migration")

	return ormDb
}
