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

	coursePathService := services.NewCoursePathService(db, rabbitMQConfig)
	coursePathController := controllers.NewCoursePathController(coursePathService, rabbitMQConfig)

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})
	app.Post("/course", classController.CreateCourse)
	app.Put("/course", classController.UpdateCourse)
	app.Delete("/course", classController.DeleteCourse)
	app.Delete("/courses", classController.DeleteAllCourses)

	app.Get("/courses", classController.ListAllCourses)
	app.Get("/course/:id", classController.GetCourseWithSubcourses)

	app.Post("/coursepath", coursePathController.CreateCoursePath)
	app.Put("/coursepath/:id", coursePathController.UpdateCoursePath)
	app.Delete("/coursepath/:id", coursePathController.DeleteCoursePath)
	app.Get("/coursepaths", coursePathController.ListAllCoursePaths)
	app.Get("/coursepath/:id", coursePathController.GetCoursePathByID)

	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "http://localhost:3000/docs/swagger.json",
	}))
}
