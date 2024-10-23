package controllers

import (
	"class/config"
	"class/models"
	"class/services"
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ClassController struct {
	classService   *services.ClassService
	rabbitMQConfig *config.RabbitMQConfig
}

func NewClassController(classService *services.ClassService, rabbitMQConfig *config.RabbitMQConfig) *ClassController {
	return &ClassController{
		classService:   classService,
		rabbitMQConfig: rabbitMQConfig,
	}
}

// ListClasses handles fetching all classes.
// @Summary List all classes
// @Description Retrieve a list of all classes
// @Produce json
// @Success 200 {array} models.Class
// @Failure 500 {object} object
// @Tags Classes
// @Router /classes [get]
func (c *ClassController) ListClasses(ctx *fiber.Ctx) error {
	classes, err := c.classService.GetAllClasses()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to fetch classes"})
	}

	return ctx.JSON(classes)
}

// CreateClass handles the creation of a class.
// @Summary Create a class
// @Description Create a new class
// @Accept json
// @Produce json
// @Param class body models.Class true "Class"
// @Success 201 {object} models.Class
// @Failure 400 {object} object
// @Failure 500 {object} object
// @Tags Classes
// @Router /classes [post]
func (c *ClassController) CreateClass(ctx *fiber.Ctx) error {
	var class models.Class
	if err := ctx.BodyParser(&class); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	classJSON, err := json.Marshal(class)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not serialize class"})
	}
	if err := c.rabbitMQConfig.PublishMessage("class_created", classJSON); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create class"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(class)
}

// UpdateClass handles the update of an existing class.
// @Summary Update a class
// @Description Update an existing class
// @Accept json
// @Produce json
// @Param class body models.Class true "Class"
// @Success 200 {object} models.Class
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Failure 500 {object} object
// @Tags Classes
// @Router /classes [put]
func (c *ClassController) UpdateClass(ctx *fiber.Ctx) error {
	var class models.Class
	if err := ctx.BodyParser(&class); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	if class.ID == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Class ID is required"})
	}

	classJSON, err := json.Marshal(class)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not serialize class"})
	}
	if err := c.rabbitMQConfig.PublishMessage("class_updated", classJSON); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update class"})
	}

	return ctx.Status(fiber.StatusOK).JSON(class)
}

// DeleteClass handles the deletion of a class.
// @Summary Delete a class
// @Description Delete a class by ID
// @Accept json
// @Produce json
// @Param id path uint true "Class ID"
// @Success 200 {object} object
// @Failure 404 {object} object
// @Failure 500 {object} object
// @Tags Classes
// @Router /classes/{id} [delete]
func (c *ClassController) DeleteClass(ctx *fiber.Ctx) error {
	classID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid class ID"})
	}

	classEvent := map[string]interface{}{
		"event_type":   "class.deleted",
		"service_name": "class_service",
		"id":           classID,
	}

	classJSON, err := json.Marshal(classEvent)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not serialize class ID"})
	}

	if err := c.rabbitMQConfig.PublishMessage("class_deleted", classJSON); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete class"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Class deleted successfully", "id": classID})
}
