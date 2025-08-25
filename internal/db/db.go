package db

import (
	"log"
	"sachahjkl/htmx_go/internal/model"

	"github.com/jacob2161/sqlitebp"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Init(url string) *gorm.DB {

	// Creates or opens the database with best practice defaults
	// (WAL, foreign keys, busy timeout, NORMAL synchronous, private cache, etc.)
	db, err := sqlitebp.OpenReadWriteCreate(url)

	if err != nil {
		log.Fatalln(err)
	}

	dialector := sqlite.New(sqlite.Config{
		Conn: db,
	})

	ormDb, err := gorm.Open(dialector, &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Migrating db @ %v", url)
	ormDb.AutoMigrate(&model.User{}, &model.Todo{})
	log.Printf("Finished migration")

	return ormDb
}
