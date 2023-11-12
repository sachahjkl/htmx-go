package handler

import (
	"sachahjkl/htmx_go/pkg/common"
	"sachahjkl/htmx_go/pkg/config"
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
		DB: db,
	}

	routes := app.Group("/auth")
	routes.Get("/", h.RedirectToLogin)
	routes.Get("/login", h.LoginPage).Name("login")
	routes.Post("/login", h.LoginSubmit)
	routes.Get("/register", h.RegisterPage)
	routes.Post("/register", h.RegisterSubmit)

	// TODO: Delete, Patch user
}

func (h *handler) RedirectToLogin(c *fiber.Ctx) error {
	return c.RedirectToRoute("login", fiber.Map{})
}

func (h *handler) LoginPage(c *fiber.Ctx) error {

	user := c.Locals(common.USER_LOCALS_KEY).(*model.User)

	if user != nil {
		return c.RedirectToRoute("todos", nil)
	}

	return c.Render("login/index", nil)
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
		"username": user.Username,
	}

	if !body.RememberMe {
		// set expiry date if the guy doesn't want to have a long (unsafe) session
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // expires in 24 days
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	encryptedToken, err := token.SignedString([]byte(h.Config.EncryptionKey))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Cookie(&fiber.Cookie{
		Name:  common.USER_COOKIE_JWT_KEY,
		Value: encryptedToken,
	})

	// redirect the guy to the main course
	return c.RedirectToRoute("todos", nil)
}
func (h *handler) RegisterPage(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotFound)
}

func (h *handler) RegisterSubmit(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotFound)

}
