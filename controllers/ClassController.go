package controllers

import (
	"leecho/services"

	"github.com/gofiber/fiber/v2"
)

type ClassController struct {
	classService *services.ClassService
}

func NewClassController(classService *services.ClassService) *ClassController {
	return &ClassController{classService: classService}
}

func (c *ClassController) ListClasses(ctx *fiber.Ctx) error {
	classes, err := c.classService.GetAllClasses()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to fetch classes"})
	}

	return ctx.JSON(classes)
}
