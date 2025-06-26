package handlers

import (
	"PromptGallery/internal/models"
	"PromptGallery/internal/services"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
)

type PromptHandler struct {
	promptService *services.PromptService
}

func NewPromptHandler(promptService *services.PromptService) *PromptHandler {
	return &PromptHandler{
		promptService: promptService,
	}
}

type APIResponse struct {
	Status  string
	Message string
	Data    interface{}
	Error   string
}

func (h *PromptHandler) GetPrompts(c *fiber.Ctx) error {

	filter, page, limit, err := h.parsePromptQuery(c)
	if err != nil {
		return c.Status(400).JSON(APIResponse{
			Status:  "error",
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
	}

	result, err := h.promptService.GetAllPrompts(filter, page, limit)
	if err != nil {
		return c.Status(500).JSON(APIResponse{
			Status:  "error",
			Message: "Internal server error",
		})
	}

	return c.Status(200).JSON(APIResponse{
		Status:  "success",
		Message: "Prompts fetched successfully",
		Data:    result,
	})
}

func (h *PromptHandler) GetPromptByID(c *fiber.Ctx) error {
	id, err := h.parseUintParam(c, "id")

	if err != nil {
		return c.Status(400).JSON(APIResponse{
			Status: "error",
			Error:  "Invalid prompt ID",
		})
	}

	prompt, err := h.promptService.GetPromptByID(id)
	if err != nil {
		return c.Status(404).JSON(APIResponse{
			Status:  "error",
			Message: "Prompt not found",
		})
	}

	return c.Status(200).JSON(APIResponse{
		Status:  "success",
		Message: "Prompt fetched successfully",
		Data:    prompt,
	})

}

func (h *PromptHandler) CreatePrompt(c *fiber.Ctx) error {

	var createReq models.PromptCreateRequest

	if err := c.BodyParser(&createReq); err != nil {
		return c.Status(400).JSON(APIResponse{
			Status:  "error",
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	prompt, err := h.promptService.CreatePrompt(&createReq)

	if err != nil {
		// Handle validation errors
		if strings.Contains(err.Error(), "required") ||
			strings.Contains(err.Error(), "invalid") {
			return c.Status(400).JSON(APIResponse{
				Status: "error",
				Error:  err.Error(),
			})
		}
		return c.Status(500).JSON(APIResponse{
			Status: "error",
			Error:  "Failed to create prompt",
		})
	}

	return c.Status(201).JSON(APIResponse{
		Status:  "success",
		Message: "Prompt created successfully",
		Data:    prompt,
	})

}

func (h *PromptHandler) DeletePrompt(c *fiber.Ctx) error {
	// Parse path parameter
	id, err := h.parseUintParam(c, "id")
	if err != nil {
		return c.Status(400).JSON(APIResponse{
			Status: "error",
			Error:  "Invalid prompt ID",
		})
	}

	// Call service
	err = h.promptService.DeletePrompt(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(404).JSON(APIResponse{
				Status: "error",
				Error:  "Prompt not found",
			})
		}
		return c.Status(500).JSON(APIResponse{
			Status: "error",
			Error:  "Failed to delete prompt",
		})
	}

	return c.Status(200).JSON(APIResponse{
		Status:  "success",
		Message: "Prompt deleted successfully",
	})
}

func (h *PromptHandler) parsePromptQuery(c *fiber.Ctx) (models.PromptFilter, int, int, error) {
	var filter models.PromptFilter

	filter.Language = c.Query("language")
	filter.Category = c.Query("category")
	filter.Search = c.Query("search")

	if verifiedStr := c.Query("is_verified"); verifiedStr != "" {
		verified, err := strconv.ParseBool(verifiedStr)
		if err != nil {
			return filter, 0, 0, err
		}
		filter.IsVerified = &verified
	}

	page := h.parseIntQuery(c, "page", 1)
	limit := h.parseIntQuery(c, "limit", 10)

	return filter, page, limit, nil

}

func (h *PromptHandler) parseUintParam(c *fiber.Ctx, param string) (uint, error) {
	paramStr := c.Params(param)
	if paramStr == "" {
		return 0, fiber.NewError(400, "Parameter is required")
	}

	value, err := strconv.ParseUint(paramStr, 10, 32)
	if err != nil {
		return 0, fiber.NewError(400, "Invalid parameter format")
	}

	return uint(value), nil
}

func (h *PromptHandler) parseIntQuery(c *fiber.Ctx, key string, defaultValue int) int {
	valueStr := c.Query(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil || value < 1 {
		return defaultValue
	}

	return value
}
