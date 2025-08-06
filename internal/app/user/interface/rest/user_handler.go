package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mqqff/savebite-be/internal/app/user/usecase"
	"github.com/mqqff/savebite-be/internal/middlewares"
)

type UserHandlerItf interface {
}

type UserHandler struct {
	userUsecase usecase.UserUsecaseItf
}

func NewUserHandler(r fiber.Router, m middlewares.MiddlewareItf, u usecase.UserUsecaseItf) {
	UserHandler := &UserHandler{userUsecase: u}

	r = r.Get("/me", m.RequireAuth, UserHandler.GetProfile)
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	user, err := h.userUsecase.GetProfile(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error_code": "internal_server_error",
			"error":      err.Error(),
			"message":    "Something happened with our end. Please try again later",
		})
	}

	payload := fiber.Map{
		"user": user,
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"payload": payload,
	})
}
