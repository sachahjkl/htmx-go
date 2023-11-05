package main

import (
	"log"
	router "sachahjkl/htmx_go/pkg/common"
	"sachahjkl/htmx_go/pkg/common/config"
	"sachahjkl/htmx_go/pkg/common/db"
	"sachahjkl/htmx_go/pkg/todos"

	"github.com/joho/godotenv"
)

func main() {

	// Try to load a potentially existing ".env" file
	// ignore errors
	godotenv.Load()

	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	app := router.New(c)

	db := db.Init(c.DBUrl)

	router.RegisterDefaultRoutes(app)
	todos.RegisterRoutes(app, db)

	log.Fatal(app.Listen(":" + c.Port))

}
