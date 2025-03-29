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

	// api routes
	routes.ListRoutes(app)
	routes.TaskRoutes(app)

	// webpage routes
	routes.AuthRoutes(app)
	app.Get("*", func(c *fiber.Ctx) error {
		//Check if the user is logged in in order to have acces to the webpage
		_, err := auth.CheckJWT(c)
		if err != nil {
			return c.SendFile("static/views/login.html")
		}
		return c.SendFile("static/views/index.html")
	})

	port := ":9000"
	log.Printf("Starting server on port%s\n", port)
	log.Fatal(app.Listen(port))
}
