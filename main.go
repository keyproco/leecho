package main

import (
	"leecho/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	routes.ClassRoutes(app)

	app.Listen(":3000")
}
