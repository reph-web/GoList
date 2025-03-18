package main

import (
	"golist/auth"
	"golist/database"
	"golist/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Static("/static", "./static")
	database.ConnectDB()

	routes.AuthRoutes(app)
	routes.ListRoutes(app)
	routes.TaskRoutes(app)

	app.Get("/", func(c *fiber.Ctx) error {
		//Check if the user is logged in in order to have acces to the webpage
		_, err := auth.CheckJWT(c)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendFile("static/index.html")
	})

	port := ":9000"
	log.Printf("Starting server on port%s\n", port)
	log.Fatal(app.Listen(port))
}
