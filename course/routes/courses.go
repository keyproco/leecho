package routes

import (
	"course/config"
	"course/controllers"

	"course/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"gorm.io/gorm"
)

func ClassRoutes(app *fiber.App, rabbitMQConfig *config.RabbitMQConfig, db *gorm.DB) {

	courseService := services.NewCourseService(db, rabbitMQConfig)
	classController := controllers.NewCourseController(courseService, rabbitMQConfig)

	app.Post("/course", classController.CreateCourse)
	app.Put("/course", classController.UpdateCourse)
	app.Delete("/course", classController.DeleteCourse)
	app.Delete("/courses", classController.DeleteAllCourses)

	app.Get("/courses", classController.ListAllCourses)
	app.Get("/course/:id", classController.GetCourseWithSubcourses)

	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "http://localhost:3000/docs/swagger.json",
	}))
}
