package handler

import (
	"sachahjkl/htmx_go/internal/common"
	"sachahjkl/htmx_go/internal/config"
	"sachahjkl/htmx_go/internal/middleware"
	"sachahjkl/htmx_go/internal/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterUserRoutes(app *fiber.App, db *gorm.DB, c *config.Config) {
	h := &handler{
		DB:     db,
		Config: c,
	}

	middleware := middleware.New(db)

	routes := app.Group("/user")

	// Authenticated only

	routes.Get("/settings", middleware.AuthenticatedOnly, h.SettingsPage).Name("settings")
	routes.Delete("/delete-me", h.DeleteMe)
	// TODO: Delete, Patch user
}

func (h *handler) DeleteMe(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h handler) SettingsPage(c *fiber.Ctx) error {
	user := c.Locals(common.USER_KEY).(*model.User)

	if user == nil {
		return fiber.NewError(fiber.StatusUnauthorized)
	}

	bind := fiber.Map{
		"Username": user.Username,
		"Config":   h.Config,
	}

	layout := ""

	// Called from HTMX, no layout
	if c.Get("HX-Request") == "" {
		layout = "layouts/main"
	}

	return c.Render("user/settings", bind, layout)
}
