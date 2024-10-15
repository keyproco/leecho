package controllers

import (
	"github.com/gofiber/fiber/v2"
)

type ClassController struct{}

func NewClassController() *ClassController {
	return &ClassController{}
}

func (c *ClassController) ListClasses(ctx *fiber.Ctx) error {
	classes := []map[string]interface{}{
		{"id": 1, "title": "Kubernetes 101", "description": "From Murrica Kube isnt Kube"},
		{"id": 2, "title": "ArgoCD", "description": "Argocd for are you?"},
		{"id": 3, "title": "Hashicorp Vault", "description": "Do you know the crochetage?"},
		{"id": 4, "title": "Opentelemetry", "description": "It depends on how much technology scares you"},
	}

	return ctx.JSON(classes)
}
