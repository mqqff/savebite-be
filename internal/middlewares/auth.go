package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mqqff/savebite-be/pkg/jwt"
	"strings"
)

func (m *Middleware) RequireAuth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if len(authHeader) < 1 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error_code": "missing_authentication_token",
			"error":      "missing authentication token",
			"message":    "You are not authorized",
		})
	}

	bearer := strings.Split(authHeader, " ")
	if len(bearer) < 2 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error_code": "invalid_authentication_format",
			"error":      "invalid authentication token format",
			"message":    "You are not authorized",
		})
	}

	token := bearer[1]

	claims := jwt.Claims{}
	err := m.jwt.Decode(token, &claims)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error_code": "invalid_authentication",
			"error":      err.Error(),
			"message":    "You are not authorized",
		})
	}

	c.Locals("userID", claims.UserID)

	return c.Next()
}
