package routes

import (
	"golist/handlers"

	"github.com/gofiber/fiber/v2"
)

func TaskRoutes(app *fiber.App) {
	api := app.Group("/api")
	// Get all tasks of a specific list from id
	api.Get("/list/:listId/tasks", handlers.ListTasksHandler)

	// Get a specific task from id
	api.Get("/task/:taskId", handlers.TaskHandler)

	// Add task routes
	api.Post("/list/:listId/task", handlers.AddTaskHandler)

	// Update task routes
	api.Patch("/task/:taskId/description", handlers.UpdateDescriptionTaskHandler)
	api.Patch("/task/:taskId/check", handlers.UpdateCheckTaskHandler)
	api.Patch("/task/:taskId/swap", handlers.SwapOrderTaskHandler)

	api.Delete("/task/:taskId", handlers.DeleteTaskHandler)
}
