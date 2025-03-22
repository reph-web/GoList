package handlers

import (
	"errors"
	"fmt"
	"golist/database"
	"golist/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Getter function
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

	// Check if the user is the owner of the list
	if err := isListOwner(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	// Get user from username and handle error
	user, err := getUser(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	// Get List ID from the URL
	listId, err := getList(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Get the task description from the body request
	var request struct {
		Description string `json:"description"`
	}
	if err := c.BodyParser(&request); err != nil || request.Description == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Find the number of tasks in the list (for the order)
	var count int64
	if err := database.DB.Model(&models.Task{}).Where("list_id = ?", listId.ID).Count(&count).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Increment the count
	count++
	// Create the task and push to DB
	task := models.Task{
		Description: request.Description,
		Checked:     false,
		TaskOrder:   count,

		UserID: user.ID,
		ListID: listId.ID,
	}

	if err := database.DB.Create(&task).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot create task in list"})
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

func SwapOrderTaskHandler(c *fiber.Ctx) error {
	//Get the task from the URL and handle error
	task1, err := getTask(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Get payload from request
	var request struct {
		OrderToSwap int `json:"OrderToSwap"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON payload"})
	}

	var task2 models.Task
	if err := database.DB.Where("task_order = ?", request.OrderToSwap).First(&task2).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task with OrderToSwap not found"})
	}
	// Find the max order to have unique value
	var maxOrder int
	database.DB.Table("tasks").Select("MAX(task_order)").Row().Scan(&maxOrder)
	tempOrder := maxOrder + 1
	TaskOneOrder := task1.TaskOrder
	TaskTwoOrder := task2.TaskOrder
	// Swap TaskOrder
	if err := database.DB.Model(&task1).Update("task_order", tempOrder).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := database.DB.Model(&task2).Update("task_order", TaskOneOrder).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := database.DB.Model(&task1).Update("task_order", TaskTwoOrder).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "task order swapped successfully",
		"task1":   task1,
		"task2":   task2,
	})
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
