package main

import (
	"log"
	"sachahjkl/htmx_go/db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{

		Prefork:       true,
		CaseSensitive: true,
		// StrictRouting: true,
		ServerHeader: "Fiber",
		AppName:      "HTMX + Go",

		Views: engine,
		// Views Layout is the global layout for all template render until override on Render function.
		ViewsLayout: "layouts/main",
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Todos": db.Todos,
		})
	})

	app.Get("/todos", func(c *fiber.Ctx) error {
		return c.Render("todos/index", fiber.Map{
			"Todos": db.Todos,
		})
	}).Name("todos")

	app.Post("/todos", func(c *fiber.Ctx) error {

		var title = c.FormValue("todo-title")
		db.AddTodo(title, false)

		return c.RedirectToRoute("todos", fiber.Map{})
	})

	app.Post("/todos/:id/toggle", func(c *fiber.Ctx) error {

		var id, err = c.ParamsInt("id")
		if err != nil {
			return fiber.ErrInternalServerError
		}

		db.ToggleTodo(id)
		return c.RedirectToRoute("todos", fiber.Map{})

	})

	app.Static("/", "./public")

	log.Fatal(app.Listen(":3000"))
}
