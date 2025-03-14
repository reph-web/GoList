package routes

import (
	"golist/handlers"

	"github.com/gofiber/fiber/v2"
)

func ListRoutes(app *fiber.App) {
	// Get all lists
	app.Get("/allLists", handlers.AllListsHandler)

	// Get a specific list from id
	app.Get("/list/:listId", handlers.ListHandler)

	app.Post("/addList", handlers.AddListHandler)

	app.Patch("/updateList/:listId", handlers.UpdateListHandler)

	app.Delete("/deleteList/:listId", handlers.DeleteListHandler)
}
