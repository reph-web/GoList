package main

import (
	"golist/auth"
	"golist/database"
	"golist/models"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Static("/static", "./static")
	database.ConnectDB()

	app.Get("/register", func(c *fiber.Ctx) error {
		return c.SendFile("static/register.html")
	})

	app.Get("/login", func(c *fiber.Ctx) error {
		return c.SendFile("static/login.html")
	})

	app.Post("/register", func(c *fiber.Ctx) error {
		return auth.CreateUser(c)
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		return auth.LoginUser(c)
	})

	app.Get("/logout", func(c *fiber.Ctx) error {
		// Remove the value and expire the cookie
		c.Cookie(&fiber.Cookie{
			Name:     "token",
			Value:    "",
			Expires:  time.Now().Add(-time.Hour),
			HTTPOnly: true,
		})

		return c.JSON(fiber.Map{"message": "Logged out successfully"})
	})

	app.Get("/todos", func(c *fiber.Ctx) error {
		//Check if the user is logged in in order to have acces to the webpage
		_, err := auth.CheckJWT(c)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendFile("static/todos.html")
	})

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

	app.Get("/list/:ListID/tasks", func(c *fiber.Ctx) error {

		//Retrieve the list ID from the URL
		listID, err := strconv.ParseUint(c.Params("listID"), 10, 32)
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
		if err := database.DB.Where("id = ?", listID).First(&list).Error; err != nil {
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

		// Retrieve all tasks for the list
		var tasks []models.Task
		if err := database.DB.Where("list_id = ?", listID).Find(&tasks).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Tasks not found",
			})
		}

		return c.JSON(tasks)
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

	app.Post("/list/:listID/addTask", func(c *fiber.Ctx) error {
		// Get the listID from the URL
		listID, err := strconv.ParseUint(c.Params("listID"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid list ID"})
		}

		// Get the username from the JWT
		username, err := auth.CheckJWT(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		// Verify if the user has access to the list
		var list models.List
		if err := database.DB.Joins("JOIN users ON users.id = lists.user_id").
			Where("lists.id = ? AND users.username = ?", listID, username).
			First(&list).Error; err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Fordidden access or list not found"})
		}

		// Get content from the body
		var request struct {
			Description string `json:"description"`
		}
		if err := c.BodyParser(&request); err != nil || request.Description == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Task description invalid"})
		}

		// Find the max task order for the list
		var maxOrder models.Task
		err = database.DB.Where("list_id = ?", list.ID).Order("task_order desc").First(&maxOrder).Error
		if err != nil && err.Error() != "record not found" {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		newOrder := maxOrder.TaskOrder + 1

		// Create the task
		task := models.Task{
			Description: request.Description,
			TaskOrder:   newOrder,
			Checked:     false,
			ListID:      list.ID,
			UserID:      list.UserID,
		}

		// Push to DB
		if err := database.DB.Create(&task).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Task not created"})
		}

		return c.Status(fiber.StatusCreated).JSON(task)
	})

	port := ":9000"
	log.Printf("Starting server on port%s\n", port)
	log.Fatal(app.Listen(port))
}
