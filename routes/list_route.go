package routes

import (
	"golist/handlers"

	"github.com/gofiber/fiber/v2"
)

func ListRoutes(app *fiber.App) {

	api := app.Group("/api")

	// Get all lists
	api.Get("/lists", handlers.AllListsHandler)

	// Get a specific list from id
	api.Get("/list/:listId", handlers.ListHandler)

	api.Post("/list", handlers.AddListHandler)

	api.Patch("/list/:listId/name", handlers.UpdateListHandler)

	api.Delete("/deleteList/:listId", handlers.DeleteListHandler)
}
