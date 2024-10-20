package services

import (
	"course/config"
	"course/models"

	"gorm.io/gorm"
)

type CourseService struct {
	DB             *gorm.DB
	rabbitMQConfig *config.RabbitMQConfig
}

func NewCourseService(db *gorm.DB, rabbitMQConfig *config.RabbitMQConfig) *CourseService {
	return &CourseService{
		DB:             db,
		rabbitMQConfig: rabbitMQConfig,
	}
}

func (s *CourseService) CreateCourse(course *models.Course) error {
	return models.CreateCourse(s.DB, course)
}
