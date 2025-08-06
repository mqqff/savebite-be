package rest

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mqqff/savebite-be/internal/app/auth/usecase"
	"github.com/mqqff/savebite-be/internal/domain/dto"
	"github.com/mqqff/savebite-be/internal/domain/env"
	"time"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecaseItf
	validator   *validator.Validate
}

func NewAuthHandler(r fiber.Router, u usecase.AuthUsecaseItf, v *validator.Validate) {
	AuthHandler := AuthHandler{
		authUsecase: u,
		validator:   v,
	}

	r = r.Group("/auth/google/oauth")
	r.Get("/redirect", AuthHandler.HandleRedirect)
	r.Post("/callback", AuthHandler.HandleCallback)
}

func (h *AuthHandler) HandleRedirect(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "oauth2state",
		Value:    env.AppEnv.OAuthState,
		MaxAge:   3600,
		Expires:  time.Now().Add(5 * time.Minute),
		HTTPOnly: true,
		Secure:   env.AppEnv.AppEnv != "development",
	})

	url, _ := h.authUsecase.HandleRedirect(env.AppEnv.OAuthState)

	payload := fiber.Map{
		"redirect_url": url,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"payload": payload,
	})
}

func (h *AuthHandler) HandleCallback(c *fiber.Ctx) error {
	req := dto.GoogleCallbackRequest{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error_code": "bad_request",
			"error":      err.Error(),
			"message":    "Invalid body request",
		})
	}

	err := h.validator.Struct(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error_code": "validation_error",
			"error":      err.Error(),
			"message":    "Validation error",
		})
	}

	token, err := h.authUsecase.HandleCallback(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error_code": "internal_server_error",
			"error":      err.Error(),
			"message":    "Something happened with our end. Please try again later",
		})
	}

	payload := fiber.Map{
		"token": token,
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"payload": payload,
	})
}
