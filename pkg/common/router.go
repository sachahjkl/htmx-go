package router

import (
	"sachahjkl/htmx_go/pkg/common/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

func New(c config.Config) *fiber.App {

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		// StrictRouting: true,
		ServerHeader: "Fiber",
		AppName:      "HTMX + Go",
		Views:        engine,
		// Views Layout is the global layout for all template render until override on Render function.
		ViewsLayout: "layouts/main",
	})
	app.Use(logger.New())
	return app
}

func RegisterDefaultRoutes(app *fiber.App) {
	app.Static("/", "./public")
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})
}
