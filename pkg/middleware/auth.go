package middleware

import (
	"sachahjkl/htmx_go/pkg/common"
	"sachahjkl/htmx_go/pkg/model"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func (m *Middleware) Protected(c *fiber.Ctx) error {
	userJWT := c.Locals(common.USER_COOKIE_JWT_KEY).(*jwt.Token)

	if userJWT == nil {
		return c.RedirectToRoute("login", nil)
	}

	claims := userJWT.Claims.(jwt.MapClaims)
	userId := claims["userId"].(uint)

	user, err := model.GetUser(m.DB, userId)

	if err != nil {
		// clear cookie, it is invalid since it didn't
		// provide a key to any real user
		c.ClearCookie(common.USER_COOKIE_JWT_KEY)
		return c.RedirectToRoute("login", nil)
	}

	// make the user available for next request
	c.Locals(common.USER_LOCALS_KEY, user)
	return c.Next()

}
