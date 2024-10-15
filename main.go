package main

import (
	"leecho/routes"

	"github.com/gofiber/fiber/v2"
)

// @title School Management API Leecho
// @version 0.1
// @description API for a school management system.
// @contact.name Akme
// @contact.url
// @contact.email
// @BasePath /api/v1
func main() {

	app := fiber.New()
	app.Static("/docs", "./public/")

	routes.ClassRoutes(app)

	app.Listen(":3000")
}
