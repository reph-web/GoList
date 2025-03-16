package handlers

import (
	"errors"
	"fmt"
	"golist/database"
	"golist/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func getTask(c *fiber.Ctx) (*models.Task, error) {
	//Retrieve the list ID from the URL and handle error
	taskID, err := strconv.ParseUint(c.Params("taskId"), 10, 32)
	fmt.Println(c.Params("taskId"))
	if err != nil {
		return nil, errors.New("invalid task ID")
	}

	// Find the list by its ID
	var task models.Task
	if err := database.DB.First(&task, taskID).Error; err != nil {
		return nil, errors.New("task not found")
	}

	// Get user from username and handle error
	user, err := getUser(c)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Find if the user is the owner of the task
	if task.UserID != user.ID {
		return nil, errors.New("user not authorized")
	}

	return &task, nil
}

// Handler functions
func TaskHandler(c *fiber.Ctx) error {
	//Get the task from the URL
	task, err := getTask(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(task)
}

func AddTaskHandler(c *fiber.Ctx) error {
	// Get user from username and handle error
	user, err := getUser(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	// Get List ID from the URL

	listId, err := getList(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid list Id"})
	}

	// Get the task description from the body request
	var request struct {
		Description string `json:"description"`
	}
	if err := c.BodyParser(&request); err != nil || request.Description == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task description"})
	}

	// Find the number of tasks in the list (for the order)
	var count int64
	if err := database.DB.Model(&models.Task{}).Where("list_id = ?", listId.ID).Count(&count).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot count tasks"})
	}

	// Create the task and push to DB
	task := models.Task{
		Description: request.Description,
		Checked:     false,
		TaskOrder:   count + 1,

		UserID: user.ID,
		ListID: listId.ID,
	}

	if err := database.DB.Create(&task).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot create list"})
	}

	return c.Status(fiber.StatusCreated).JSON(task)
}

func UpdateDescriptionTaskHandler(c *fiber.Ctx) error {
	//Get the task from the URL and handle error
	task, err := getTask(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get the new task description from the body request
	var request struct {
		Description string `json:"description"`
	}
	if err := c.BodyParser(&request); err != nil || request.Description == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid list name"})
	}

	// Update the task description and push to DB
	task.Description = request.Description
	if err := database.DB.Save(&task).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot update task"})
	}

	return c.Status(fiber.StatusCreated).JSON(task)
}

func UpdateCheckTaskHandler(c *fiber.Ctx) error {
	//Get the task from the URL and handle error
	task, err := getTask(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Update the task check state and push to DB
	if task.Checked {
		task.Checked = false
	} else {
		task.Checked = true
	}

	if err := database.DB.Save(&task).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot update task"})
	}

	return c.Status(fiber.StatusCreated).JSON(task)
}

func UpdateOrderTaskHandler(c *fiber.Ctx) error {
	//Get the task from the URL and handle error
	task, err := getTask(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get the task to swap id from the body request
	var request struct {
		TaskToSwapId string `json:"taskToSwapId"`
	}
	if err := c.BodyParser(&request); err != nil || request.TaskToSwapId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid swapped task id"})
	}

	// Find the swapped task by its ID
	var taskToSwap models.Task
	if err := database.DB.First(&taskToSwap, request.TaskToSwapId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task to swap not found"})
	}

	// Swap the task orders and push to DB
	var tempTaskOrder int64

	tempTaskOrder = task.TaskOrder
	task.TaskOrder = taskToSwap.TaskOrder
	taskToSwap.TaskOrder = tempTaskOrder

	if err := database.DB.Save(&task).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot update task"})
	}

	if err := database.DB.Save(&taskToSwap).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot update task to swap"})
	}

	return c.Status(fiber.StatusCreated).JSON(task)
}

func DeleteTaskHandler(c *fiber.Ctx) error {
	//Get the task from the URL and handle error
	task, err := getTask(c)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Delete the list
	if err := database.DB.Delete(&task).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete list",
		})
	}
	return c.JSON(fiber.Map{"message": "Task deleted"})
}
