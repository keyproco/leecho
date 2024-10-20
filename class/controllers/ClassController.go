package controllers

import (
	"class/config"
	"class/models"
	"class/services"
	"encoding/json"

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

func (c *ClassController) ListClasses(ctx *fiber.Ctx) error {
	classes, err := c.classService.GetAllClasses()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to fetch classes"})
	}

	return ctx.JSON(classes)
}

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
