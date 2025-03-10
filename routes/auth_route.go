package routes

import (
	"golist/auth"
	"time"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {
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
}
