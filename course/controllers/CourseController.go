package controllers

import (
	"course/config"
	"course/models"
	"course/services"
	"encoding/json"
	"strconv"
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

// TODO Panigation

func (c *CourseController) ListAllCourses(ctx *fiber.Ctx) error {
	courses, err := c.courseService.ListAllCourses()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to list courses"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": courses})
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

func (c *CourseController) UpdateCourse(ctx *fiber.Ctx) error {
	var course models.Course
	if err := ctx.BodyParser(&course); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	courseEvent := map[string]interface{}{
		"event_type":   "course.updated",
		"service_name": "course_service",
		"course":       course,
		"timestamp":    time.Now().Unix(),
	}

	courseJSON, err := json.Marshal(courseEvent)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not serialize course data"})
	}

	if err := c.rabbitMQConfig.PublishMessage("course_events", courseJSON); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update course"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Course updated successfully", "course": course})
}

func (c *CourseController) DeleteAllCourses(ctx *fiber.Ctx) error {
	var requestBody struct {
		IDs []uint `json:"ids"`
	}

	if err := ctx.BodyParser(&requestBody); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	for _, id := range requestBody.IDs {
		courseEvent := map[string]interface{}{
			"event_type":   "course.deleted",
			"service_name": "course_service",
			"id":           id,
			"timestamp":    time.Now().Unix(),
		}

		courseJSON, err := json.Marshal(courseEvent)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not serialize course ID"})
		}

		if err := c.rabbitMQConfig.PublishMessage("course_events", courseJSON); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete course"})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Courses deleted successfully", "ids": requestBody.IDs})
}

func (c *CourseController) GetCourseWithSubcourses(ctx *fiber.Ctx) error {
	courseID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid course ID"})
	}

	course, err := c.courseService.GetCourseWithSubcourses(uint(courseID))
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Course not found"})
	}

	return ctx.Status(fiber.StatusOK).JSON(course)
}
