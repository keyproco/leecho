package controllers

import (
	"course/config"
	"course/models"
	"course/requests"
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

// TODO Pagination

func (c *CourseController) ListAllCourses(ctx *fiber.Ctx) error {
	courses, err := c.courseService.ListAllCourses()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to list courses"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": courses})
}

// CreateCourse handles the creation of a course.
// @Summary Create a course
// @Description Create a new course
// @Accept json
// @Produce json
// @Param course body requests.CourseCreateRequest true "CourseCreateRequest"
// @Success 201 {object} requests.CourseCreateRequest
// @Failure 400 {object} object
// @Router /course [post]
// @tags Courses
func (c *CourseController) CreateCourse(ctx *fiber.Ctx) error {
	var courseRequest requests.CourseCreateRequest
	if err := ctx.BodyParser(&courseRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	course := models.Course{
		Title:       courseRequest.Title,
		Description: courseRequest.Description,
		Category:    courseRequest.Category,
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

// DeleteCourse deletes a course.
// @Summary Delete a course
// @Description Delete a course by ID
// @Accept json
// @Produce json
// @Param id body uint true "Course ID"
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 500 {object} object
// @Router /course [delete]
// @tags Courses
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

// UpdateCourse updates a course.
// @Summary Update a course
// @Description Update an existing course by ID
// @Accept json
// @Produce json
// @Param course body models.Course true "Course"
// @Success 200 {object} models.Course
// @Failure 400 {object} object
// @Failure 500 {object} object
// @Router /course [put]
// @tags Courses
func (c *CourseController) UpdateCourse(ctx *fiber.Ctx) error {
	var course models.Course

	if err := ctx.BodyParser(&course); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if course.ID == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Course ID is required for updating"})
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

// DeleteAllCourses deletes multiple courses.
// @Summary Delete multiple courses
// @Description Delete multiple courses by their IDs
// @Accept json
// @Produce json
// @Param ids body []uint true "Course IDs"
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 500 {object} object
// @Router /courses/delete [delete]
// @tags Courses
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

// GetCourseWithSubcourses gets a course with its subcourses.
// @Summary Get a course with subcourses
// @Description Retrieve a course and its subcourses by ID
// @Accept json
// @Produce json
// @Param id path uint true "Course ID"
// @Success 200 {object} models.Course
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Router /course/{id} [get]
// @tags Courses
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
