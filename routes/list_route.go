package routes

import (
	"golist/handlers"

	"github.com/gofiber/fiber/v2"
)

func ListRoutes(app *fiber.App) {
	// Get all lists
	app.Get("/lists", handlers.AllListsHandler)

	// Get a specific list from id
	app.Get("/list/:listId", handlers.ListHandler)

	app.Post("/list", handlers.AddListHandler)

	app.Patch("/list/:listId/name", handlers.UpdateListHandler)

	app.Delete("/deleteList/:listId", handlers.DeleteListHandler)
}
