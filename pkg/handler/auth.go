package handler

import (
	"sachahjkl/htmx_go/pkg/common"
	"sachahjkl/htmx_go/pkg/config"
	"sachahjkl/htmx_go/pkg/middleware"
	"sachahjkl/htmx_go/pkg/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type LoginRequestBody struct {
	Username   string `form:"username"`
	Password   string `form:"password"`
	RememberMe bool   `form:"remember-me"`
}

type RegisterRequestBody struct {
	Username        string `form:"username"`
	Password        string `form:"password"`
	PasswordConfirm string `form:"password-confirm"`
}

func RegisterAuthRoutes(app *fiber.App, db *gorm.DB, c *config.Config) {
	h := &handler{
		DB:     db,
		Config: c,
	}

	middleware := middleware.New(db)

	routes := app.Group("/auth")

	app.Get("/login", h.RedirectToLogin)
	app.Get("/register", func(c *fiber.Ctx) error { return c.RedirectToRoute("register", nil) })
	app.Get("/logout", func(c *fiber.Ctx) error { return c.RedirectToRoute("logout", nil) })

	routes.Get("/", middleware.UnauthenticatedOnly, h.RedirectToLogin)
	routes.Get("/login", middleware.UnauthenticatedOnly, h.LoginPage).Name("login")
	routes.Post("/login", middleware.UnauthenticatedOnly, h.LoginSubmit)
	routes.Get("/register", middleware.UnauthenticatedOnly, h.RegisterPage)
	routes.Post("/register", middleware.UnauthenticatedOnly, h.RegisterSubmit)

	// Authenticated only

	routes.Get("/logout", middleware.AuthenticatedOnly, h.LogoutPage).Name("logout")
	routes.Post("/logout", middleware.AuthenticatedOnly, h.Logout)

	// TODO: Delete, Patch user
}

func (h *handler) RedirectToLogin(c *fiber.Ctx) error {
	return c.RedirectToRoute("login", nil)
}

func (h *handler) LoginPage(c *fiber.Ctx) error {
	return c.Render("auth/login", fiber.Map{
		"UsernameMinLength": model.MIN_USERNAME_LEN,
	}, "layouts/main")
}
func (h *handler) LogoutPage(c *fiber.Ctx) error {
	return c.Render("auth/logout", nil)
}

func (h *handler) Logout(c *fiber.Ctx) error {

	c.Cookie(&fiber.Cookie{
		Name: common.USER_COOKIE_KEY,
		// expires in the past
		Expires:  time.Now().Add(-(time.Hour * 2)),
		HTTPOnly: true,
		SameSite: "lax",
	})
	return c.RedirectToRoute("login", nil)
}

func (h *handler) LoginSubmit(c *fiber.Ctx) error {

	body := LoginRequestBody{}

	// parse body, unmarshall to AddTodoRequestBody struct
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user, err := model.LoginUser(h.DB, body.Username, body.Password)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Create the Claims
	claims := jwt.MapClaims{
		common.USER_CLAIM_KEY: user.ID,
	}

	expires := time.Now().Add(time.Hour * 72) // expires in 3 days

	if !body.RememberMe {
		// set expiry date if the guy doesn't want to have a long (unsafe) session
		claims["exp"] = expires.Unix()
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	encryptedToken, err := token.SignedString([]byte(h.Config.EncryptionKey))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Cookie(&fiber.Cookie{
		Name:     common.USER_COOKIE_KEY,
		Value:    encryptedToken,
		SameSite: "l²ax",
		HTTPOnly: true,
		Expires:  expires,
	})

	// redirect the guy to the main course
	// return c.SendStatus(fiber.StatusOK)
	return c.RedirectToRoute("todos", nil)
}
func (h *handler) RegisterPage(c *fiber.Ctx) error {
	return c.Render("auth/register", fiber.Map{
		"UsernameMinLength": model.MIN_USERNAME_LEN,
	}, "layouts/main")
}

func (h *handler) RegisterSubmit(c *fiber.Ctx) error {
	body := RegisterRequestBody{}

	// parse body, unmarshall to AddTodoRequestBody struct
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	_, err := model.CreateUser(h.DB, body.Username, body.Password, body.PasswordConfirm)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// redirect the guy to the login page
	return c.RedirectToRoute("login", nil)
}
