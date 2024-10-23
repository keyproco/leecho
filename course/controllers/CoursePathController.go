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

type CoursePathController struct {
	coursePathService *services.CoursePathService
	rabbitMQConfig    *config.RabbitMQConfig
}

func NewCoursePathController(CoursePathService *services.CoursePathService, rabbitMQConfig *config.RabbitMQConfig) *CoursePathController {
	return &CoursePathController{
		coursePathService: CoursePathService,
		rabbitMQConfig:    rabbitMQConfig,
	}
}

// ListAllCoursePaths handles listing all course paths.
// @Summary List all course paths
// @Description Retrieve a list of all course paths
// @Produce json
// @Success 200 {array} models.CoursePath
// @Failure 500 {object} object
// @Router /coursepaths [get]
// @tags CoursePaths
func (c *CoursePathController) ListAllCoursePaths(ctx *fiber.Ctx) error {
	coursePaths, err := c.coursePathService.ListAllCoursePaths()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to list course paths"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": coursePaths})
}

// CreateCoursePath handles the creation of a course path.
// @Summary Create a course path
// @Description Create a new course path
// @Accept json
// @Produce json
// @Param coursePath body models.CoursePath true "CoursePath"
// @Success 201 {object} models.CoursePath
// @Failure 400 {object} object
// @Router /coursepath [post]
// @tags CoursePaths
func (c *CoursePathController) CreateCoursePath(ctx *fiber.Ctx) error {
	var coursePath models.CoursePath

	if err := ctx.BodyParser(&coursePath); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	courseEvent := map[string]interface{}{
		"event_type":   "course_path.created",
		"service_name": "course_path_service",
		"course_path":  coursePath,
	}

	courseJSON, err := json.Marshal(courseEvent)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not serialize course path"})
	}

	if err := c.rabbitMQConfig.PublishMessage("course_path_events", courseJSON); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create course path"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(coursePath)
}

// UpdateCoursePath updates a course path.
// @Summary Update a course path
// @Description Update an existing course path by ID
// @Accept json
// @Produce json
// @Param coursePath body models.CoursePath true "CoursePath"
// @Success 200 {object} models.CoursePath
// @Failure 400 {object} object
// @Failure 500 {object} object
// @Router /coursepath [put]
// @tags CoursePaths
func (c *CoursePathController) UpdateCoursePath(ctx *fiber.Ctx) error {
	var coursePath models.CoursePath

	if err := ctx.BodyParser(&coursePath); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if coursePath.ID == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Course Path ID is required for updating"})
	}

	courseEvent := map[string]interface{}{
		"event_type":   "course_path.updated",
		"service_name": "course_path_service",
		"course_path":  coursePath,
		"timestamp":    time.Now().Unix(),
	}

	courseJSON, err := json.Marshal(courseEvent)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not serialize course path data"})
	}

	if err := c.rabbitMQConfig.PublishMessage("course_path_events", courseJSON); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update course path"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Course path updated successfully", "course_path": coursePath})
}

// DeleteCoursePath deletes a course path.
// @Summary Delete a course path
// @Description Delete a course path by ID
// @Accept json
// @Produce json
// @Param id body uint true "Course Path ID"
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 500 {object} object
// @Router /coursepath [delete]
// @tags CoursePaths
func (c *CoursePathController) DeleteCoursePath(ctx *fiber.Ctx) error {
	var requestBody struct {
		ID uint `json:"id"`
	}

	if err := ctx.BodyParser(&requestBody); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	courseEvent := map[string]interface{}{
		"event_type":   "course_path.deleted",
		"service_name": "course_path_service",
		"id":           requestBody.ID,
		"timestamp":    time.Now().Unix(),
	}

	courseJSON, err := json.Marshal(courseEvent)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not serialize course path ID"})
	}

	if err := c.rabbitMQConfig.PublishMessage("course_path_events", courseJSON); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete course path"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Course path deleted successfully", "id": requestBody.ID})
}

// GetCoursePathByID retrieves a course path by ID.
// @Summary Get a course path by ID
// @Description Retrieve a course path by ID
// @Accept json
// @Produce json
// @Param id path uint true "Course Path ID"
// @Success 200 {object} models.CoursePath
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Router /coursepath/{id} [get]
// @tags CoursePaths
func (c *CoursePathController) GetCoursePathByID(ctx *fiber.Ctx) error {
	coursePathID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid course path ID"})
	}

	coursePath, err := c.coursePathService.GetCoursePathByID(uint(coursePathID))
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Course path not found"})
	}

	return ctx.Status(fiber.StatusOK).JSON(coursePath)
}
