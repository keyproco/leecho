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

func (s *CourseService) UpdateCourse(courseID uint, updatedData *models.Course) error {
	return models.UpdateCourse(s.DB, courseID, updatedData)
}

func (s *CourseService) DeleteCourse(courseID uint) error {
	return models.DeleteCourse(s.DB, courseID)
}

func (s *CourseService) DeleteMultipleCourses(courseIDs []uint) error {
	return models.DeleteMultipleCourses(s.DB, courseIDs)
}

func (s *CourseService) GetCourseWithSubcourses(courseID uint) (*models.Course, error) {
	var course models.Course
	if err := s.DB.Preload("SubCourses").First(&course, courseID).Error; err != nil {
		return nil, err
	}
	return &course, nil
}

func (s *CourseService) ListAllCourses() ([]models.Course, error) {
	var courses []models.Course

	if err := s.DB.Preload("SubCourses").
		Where("parent_course_id IS NULL").
		Find(&courses).Error; err != nil {
		return nil, err
	}

	return courses, nil
}
