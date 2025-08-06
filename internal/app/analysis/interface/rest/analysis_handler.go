package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mqqff/savebite-be/internal/app/analysis/usecase"
	"github.com/mqqff/savebite-be/internal/domain/dto"
	"github.com/mqqff/savebite-be/internal/middlewares"
	"strings"
)

type AnalysisHandler struct {
	analysisUsecase usecase.AnalysisUsecaseItf
}

func NewAnalysisHandler(r fiber.Router, m middlewares.MiddlewareItf, u usecase.AnalysisUsecaseItf) {
	AnalysisHandler := AnalysisHandler{
		analysisUsecase: u,
	}

	r = r.Group("/", m.RequireAuth)
	r.Get("/me/analyses", AnalysisHandler.GetHistory)

	r = r.Group("/analyses")
	r.Post("/", AnalysisHandler.Analyze)
}

func (h *AnalysisHandler) Analyze(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error_code": "bad_request",
			"error":      err.Error(),
			"message":    "Invalid body request",
		})
	}

	if file == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error_code": "validation_error",
			"error":      "missing image",
			"message":    "Please provide an image",
		})
	}

	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error_code": "validation_error",
			"error":      "invalid image content type",
			"message":    "Please provide valid image",
		})
	}

	userID := c.Locals("userID").(uuid.UUID)

	res, err := h.analysisUsecase.Analyze(file, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error_code": "internal_server_error",
			"error":      err.Error(),
			"message":    "Something happened with our end. Please try again later",
		})
	}

	payload := fiber.Map{
		"analysis": res,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"payload": payload,
	})
}

func (h *AnalysisHandler) GetHistory(c *fiber.Ctx) error {
	req := dto.PaginationRequest{}
	err := c.QueryParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error_code": "bad_request",
			"error":      err.Error(),
			"message":    "Invalid query request",
		})
	}

	userID := c.Locals("userID").(uuid.UUID)

	res, meta, err := h.analysisUsecase.GetHistory(userID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error_code": "internal_server_error",
			"error":      err.Error(),
			"message":    "Something happened with our end. Please try again later",
		})
	}

	payload := fiber.Map{
		"meta":     meta,
		"analyses": res,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"payload": payload,
	})
}
