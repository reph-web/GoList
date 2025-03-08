package auth

import (
	"golist/database"
	"golist/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).SendString("Error parsing request body")
	}

	// Check if username and password are not empty
	if user.Username == "" || user.Password == "" {
		return c.Status(400).SendString("Username and password are required")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).SendString("Error hashing password")
	}

	// Replace the password with the hashed password
	user.Password = string(hashedPassword)

	// Push data to DB
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(500).SendString("Error creating user")
	}

	// Return the user without the password displayed
	user.Password = ""

	return c.Status(201).JSON(user)
}
