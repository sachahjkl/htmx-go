package router

import (
	"fmt"
	"sachahjkl/htmx_go/pkg/common"
	"sachahjkl/htmx_go/pkg/config"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/gofiber/template/html/v2"
)

func New(c *config.Config) *fiber.App {

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		// Prefork:       true,
		CaseSensitive: true,
		// StrictRouting: true,
		ServerHeader: "Fiber",
		AppName:      "HTMX + Go",
		Views:        engine,
		// Views Layout is the global layout for all template render until override on Render function.
		// ViewsLayout: "layouts/main",
	})

	// Log in console
	app.Use(logger.New())

	// Compress request
	app.Use(compress.New())

	// Protect server
	app.Use(helmet.New())

	// use jwt middleware for authentication
	app.Use(jwtware.New(jwtware.Config{
		SigningKey:  jwtware.SigningKey{Key: []byte(c.EncryptionKey)},
		TokenLookup: fmt.Sprintf("cookie:%v", common.USER_COOKIE_JWT_KEY), // should be "cookie:userJWT"
		ContextKey:  common.USER_LOCALS_KEY,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.RedirectToRoute("login", nil)
		},
	}))

	return app
}

func RegisterDefaultRoutes(app *fiber.App) {
	app.Static("/", "./assets")
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", nil, "layouts/main")
	})
}
