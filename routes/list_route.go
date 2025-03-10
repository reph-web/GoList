package routes

import (
	"golist/auth"
	"golist/database"
	"golist/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func ListRoutes(app *fiber.App) {
	// Get all lists
	app.Get("/lists", func(c *fiber.Ctx) error {
		username, err := auth.CheckJWT(c)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": err.Error()})
		}

		// Retrieve all lists for the user
		var user models.User
		if err := database.DB.Preload("Lists").First(&user, "username = ?", username).Error; err != nil {
			return c.Status(500).SendString("Error fetching user from database")
		}

		return c.JSON(user.Lists)
	})

	// Get a specific list from id
	app.Get("/list/:id", func(c *fiber.Ctx) error {

		//Retrieve the list ID from the URL
		listID, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid list ID",
			})
		}

		// Get the username from the JWT
		username, err := auth.CheckJWT(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		// Find the list by its ID
		var list models.List
		if err := database.DB.Preload("Tasks").Where("id = ?", listID).First(&list).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "List not found",
			})
		}

		// Find the userID by its username
		var user models.User
		if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		// Find if the user is the owner of the list
		if list.UserID != user.ID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden access",
			})
		}
		return c.JSON(list)
	})

	app.Post("/addList", func(c *fiber.Ctx) error {
		// Get the username from the JWT
		username, err := auth.CheckJWT(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		// Find the userID by its username
		var user models.User
		if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}

		// Get the list name from the body request
		var request struct {
			Name string `json:"name"`
		}
		if err := c.BodyParser(&request); err != nil || request.Name == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid list name"})
		}

		// Create the list and push to DB
		list := models.List{Name: request.Name, UserID: user.ID}
		if err := database.DB.Create(&list).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot create list"})
		}

		return c.Status(fiber.StatusCreated).JSON(list)
	})

	app.Patch("/updateList/:id", func(c *fiber.Ctx) error {
		//Retrieve the list ID from the URL
		listID, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid list ID",
			})
		}

		// Get the username from the JWT
		username, err := auth.CheckJWT(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		// Find the list by its ID
		var list models.List
		if err := database.DB.Preload("Tasks").Where("id = ?", listID).First(&list).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "List not found",
			})
		}

		// Find the userID by its username
		var user models.User
		if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		// Find if the user is the owner of the list
		if list.UserID != user.ID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden access",
			})
		}

		// Get the new list name from the body request
		var request struct {
			Name string `json:"name"`
		}
		if err := c.BodyParser(&request); err != nil || request.Name == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid list name"})
		}

		// Update the list name and push to DB
		list.Name = request.Name
		if err := database.DB.Save(&list).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot update list"})
		}

		return c.Status(fiber.StatusCreated).JSON(list)
	})

	app.Delete("/deleteList/:id", func(c *fiber.Ctx) error {
		//Retrieve the list ID from the URL
		listID, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid list ID",
			})
		}

		// Get the username from the JWT
		username, err := auth.CheckJWT(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		// Find the list by its ID
		var list models.List
		if err := database.DB.Preload("Tasks").Where("id = ?", listID).First(&list).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "List not found",
			})
		}

		// Find the userID by its username
		var user models.User
		if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		// Find if the user is the owner of the list
		if list.UserID != user.ID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden access",
			})
		}

		// Start to delete all tasks from the list
		if err := database.DB.Where("list_id = ?", listID).Delete(&models.Task{}).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete tasks",
			})
		}

		// Delete the list
		if err := database.DB.Delete(&list).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete list",
			})
		}
		return c.JSON(user.Lists)
	})
}
