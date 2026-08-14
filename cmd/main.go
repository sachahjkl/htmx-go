package main

import (
	"log"
	"sachahjkl/htmx_go/internal/config"
	"sachahjkl/htmx_go/internal/db"
	"sachahjkl/htmx_go/internal/handler"
	"sachahjkl/htmx_go/internal/router"

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

	log.Printf("Version : %v\n", c.Version)

	db := db.Init(c.DBUrl)

	app := router.New(c)

	router.RegisterDefaultRoutes(app, db)
	handler.RegisterTodoRoutes(app, db, c)
	handler.RegisterAuthRoutes(app, db, c)
	handler.RegisterUserRoutes(app, db, c)

	log.Fatal(app.Listen("0.0.0.0:" + c.Port))

}
