package main

import (
	"PromptGallery/internal/config"
	"PromptGallery/internal/database"
	"PromptGallery/internal/handlers"
	"PromptGallery/internal/repositories"
	"PromptGallery/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"log"
)

func main() {

	cfg := config.LoadConfig()

	err := database.ConnectDatabase(cfg.DatabaseURL, cfg.Environment)
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	defer database.CloseDatabase()

	app := fiber.New(fiber.Config{
		AppName: "PromptGallery API v1.0",
	})

	setUpMiddlewares(app)

	setupDependencies(app)

	log.Println("Server running on port %s ", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))

}

func setUpMiddlewares(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Use(logger.New(logger.Config{
		Format: "[${ip}] ${status} - ${method} ${path}\n",
	}))
}

func setupDependencies(app *fiber.App) {
	db := database.GetDb()

	promptRepo := repositories.NewPromptRepository(db)

	promptService := services.NewPromptService(promptRepo)

	promptHandler := handlers.NewPromptHandler(promptService)

	setupRoutes(app, promptHandler)
}

func setupRoutes(app *fiber.App, promptHandler *handlers.PromptHandler) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Prompt Gallery API is running",
		})
	})

	api := app.Group("/api/v1")

	// Prompt routes
	setupPromptRoutes(api, promptHandler)

	// 404 handler (catch-all)
	app.Use("*", func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Route not found",
		})
	})
}

func setupPromptRoutes(router fiber.Router, handler *handlers.PromptHandler) {
	prompts := router.Group("/prompts")

	// CRUD routes
	prompts.Get("/", handler.GetPrompts)
	prompts.Post("/", handler.CreatePrompt)
	prompts.Get("/:id", handler.GetPromptByID)
	prompts.Delete("/:id", handler.DeletePrompt)

}
