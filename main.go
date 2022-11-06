package main

import (
	"ecart/src/database"
	"ecart/src/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	database.Connect()
	database.Automigrate()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	},
	))
	routes.Setup(app)

	log.Fatal(app.Listen(":8000"))
}
