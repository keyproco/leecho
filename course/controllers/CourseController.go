package controllers

import (
	"course/config"
	"course/models"
	"course/services"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

type CourseController struct {
	courseService  *services.CourseService
	rabbitMQConfig *config.RabbitMQConfig
}

func NewCourseController(CourseService *services.CourseService, rabbitMQConfig *config.RabbitMQConfig) *CourseController {
	return &CourseController{
		courseService:  CourseService,
		rabbitMQConfig: rabbitMQConfig,
	}
}

func (c *CourseController) CreateCourse(ctx *fiber.Ctx) error {
	var course models.Course
	if err := ctx.BodyParser(&course); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	courseJSON, err := json.Marshal(course)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not serialize course"})
	}
	if err := c.rabbitMQConfig.PublishMessage("course_created", courseJSON); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create course"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(course)

}
