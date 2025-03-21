package routes

import (
	"golist/handlers"

	"github.com/gofiber/fiber/v2"
)

func ListRoutes(app *fiber.App) {

	// Get all lists
	app.Get("/api/lists", handlers.AllListsHandler)

	// Get a specific list from id
	app.Get("api/list/:listId", handlers.ListHandler)

	app.Post("api/list", handlers.AddListHandler)

	app.Patch("api/list/:listId/name", handlers.UpdateListHandler)

	app.Delete("api/deleteList/:listId", handlers.DeleteListHandler)
}
