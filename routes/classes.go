package routes

import (
	"leecho/controllers"
	"leecho/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func ClassRoutes(app *fiber.App) {
	classService := services.NewClassService()
	classController := controllers.NewClassController(classService)

	app.Get("/classes", classController.ListClasses)
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "http://localhost:3000/docs/swagger.json",
	}))

}
