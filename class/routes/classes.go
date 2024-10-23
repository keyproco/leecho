package routes

import (
	"class/config"
	"class/controllers"

	"class/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"gorm.io/gorm"
)

func ClassRoutes(app *fiber.App, rabbitMQConfig *config.RabbitMQConfig, db *gorm.DB) {
	classService := services.NewClassService(db, rabbitMQConfig)
	classController := controllers.NewClassController(classService, rabbitMQConfig)

	app.Get("/classes", classController.ListClasses)

	app.Post("/class", classController.CreateClass)

	app.Put("/class/:id", classController.UpdateClass)

	app.Delete("/class/:id", classController.DeleteClass)

	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "http://localhost:3000/docs/swagger.json",
	}))
}
