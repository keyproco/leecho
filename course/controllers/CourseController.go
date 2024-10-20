package controllers

import (
	"course/config"
	"course/models"
	"course/services"
	"encoding/json"
	"time"

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
	courseEvent := map[string]interface{}{
		"event_type":   "course.created",
		"service_name": "course_service",
		"course":       course,
	}
	courseJSON, err := json.Marshal(courseEvent)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not serialize course"})
	}
	if err := c.rabbitMQConfig.PublishMessage("course_events", courseJSON); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create course"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(course)
}
func (c *CourseController) DeleteCourse(ctx *fiber.Ctx) error {
	var requestBody struct {
		ID uint `json:"id"`
	}

	if err := ctx.BodyParser(&requestBody); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	courseEvent := map[string]interface{}{
		"event_type":   "course.deleted",
		"service_name": "course_service",
		"id":           requestBody.ID,
		"timestamp":    time.Now().Unix(),
	}

	courseJSON, err := json.Marshal(courseEvent)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not serialize course ID"})
	}

	if err := c.rabbitMQConfig.PublishMessage("course_events", courseJSON); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete course"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Course deleted successfully", "id": requestBody.ID})
}
