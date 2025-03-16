package routes

import (
	"golist/handlers"

	"github.com/gofiber/fiber/v2"
)

func TaskRoutes(app *fiber.App) {

	// Get all tasks of a specific list from id
	app.Get("/list/:listId/tasks", handlers.ListTasksHandler)

	// Get a specific task from id
	app.Get("/task/:taskId", handlers.TaskHandler)

	// Add task routes
	app.Post("/list/:listId/task", handlers.AddTaskHandler)

	// Update task routes
	app.Patch("/task/:taskId/description", handlers.UpdateDescriptionTaskHandler)
	app.Patch("/task/:taskId/check", handlers.UpdateCheckTaskHandler)
	app.Patch("/task/:taskId/swapOrder", handlers.UpdateOrderTaskHandler)

	app.Delete("/deleteTask/:taskId", handlers.DeleteTaskHandler)
}
