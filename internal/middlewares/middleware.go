package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mqqff/savebite-be/pkg/jwt"
)

type MiddlewareItf interface {
	RequireAuth(c *fiber.Ctx) error
}

type Middleware struct {
	jwt jwt.JWTIf
}

func NewMiddleware(jwt jwt.JWTIf) MiddlewareItf {
	return &Middleware{
		jwt: jwt,
	}
}
