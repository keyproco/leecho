package routes

import (
	"leecho/controllers"

	"github.com/gofiber/fiber/v2"
)

func ClassRoutes(app *fiber.App) {
	classController := controllers.NewClassController()

	app.Get("/classes", classController.ListClasses)
}
