package handlers

import (
	"errors"
	"golist/auth"
	"golist/database"
	"golist/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Getter functions

func getUser(c *fiber.Ctx) (*models.User, error) {
	// Check JWT et get username
	username, err := auth.CheckJWT(c)
	if err != nil {
		return nil, errors.New("unauthorized: invalid token")
	}

	// Search user by its username
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// return pointer of its ID
	return &user, nil
}

func getList(c *fiber.Ctx) (*models.List, error) {

	//Retrieve the list ID from the URL and handle error
	listID, err := strconv.ParseUint(c.Params("listId"), 10, 32)
	if err != nil {
		return nil, errors.New("invalid list ID")
	}

	// Find the list by its ID
	var list models.List
	if err := database.DB.Preload("Tasks").First(&list, listID).Error; err != nil {
		return nil, errors.New("list not found")
	}

	return &list, nil
}

// Handler functions

func AllListsHandler(c *fiber.Ctx) error {
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
}

func ListHandler(c *fiber.Ctx) error {
	//Get the list from the URL
	list, err := getList(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get user from username and handle error
	user, err := getUser(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Find if the user is the owner of the list
	if list.UserID != user.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden access",
		})
	}
	return c.JSON(list)
}

func AddListHandler(c *fiber.Ctx) error {

	// Get user from username and handle error
	user, err := getUser(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
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
}

func UpdateListHandler(c *fiber.Ctx) error {
	//Get the list from the URL and handle error
	list, err := getList(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get user from username and handle error
	user, err := getUser(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
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
}

func DeleteListHandler(c *fiber.Ctx) error {
	//Get the list from the URL and handle error
	list, err := getList(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get user from username and handle error
	user, err := getUser(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Find if the user is the owner of the list
	if list.UserID != user.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden access",
		})
	}

	// Start to delete all tasks from the list
	if err := database.DB.Where("list_id = ?", list.ID).Delete(&models.Task{}).Error; err != nil {
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
}

func ListTasksHandler(c *fiber.Ctx) error {
	//Get the list from the URL and handle error
	list, err := getList(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get user from username and handle error
	user, err := getUser(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Find if the user is the owner of the list
	if list.UserID != user.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden access",
		})
	}

	return c.JSON(list.Tasks)
}
