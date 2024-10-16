package services

import "leecho/models"

type ClassService struct{}

func NewClassService() *ClassService {
	return &ClassService{}
}

func (s *ClassService) GetAllClasses() ([]models.Class, error) {
	classes := []models.Class{
		{ID: 1, Title: "Kubernetes 101", Description: "From Murrica Kube isnt Kube"},
		{ID: 2, Title: "ArgoCD", Description: "ArgoCD to Git up and go?"},
		{ID: 3, Title: "Hashicorp Vault", Description: "Do you know the crochetage?"},
		{ID: 4, Title: "Opentelemetry", Description: "It depends on how much technology scares you"},
	}
	return classes, nil
}
