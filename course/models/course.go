package models

import (
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"size:255;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Course struct {
	ID              uint         `json:"id" gorm:"primaryKey;autoIncrement"`
	Title           string       `json:"title" gorm:"size:255;not null"`
	Description     string       `json:"description" gorm:"size:1024"`
	Category        string       `json:"category" gorm:"size:100;not null"`
	CreatedAt       time.Time    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time    `json:"updated_at" gorm:"autoUpdateTime"`
	Instructors     []Instructor `json:"instructors" gorm:"many2many:course_instructors;constraint:OnDelete:CASCADE;"`
	EnrollmentLimit int          `json:"enrollment_limit"`
	Tags            []Tag        `json:"tags" gorm:"many2many:course_tags;constraint:OnDelete:CASCADE;"`
	SubCourses      []Course     `json:"sub_courses" gorm:"foreignKey:ParentCourseID;constraint:OnDelete:CASCADE;"`
	ParentCourseID  *uint        `json:"parent_course_id"`
}

type CoursePath struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string    `json:"title" gorm:"size:255;not null"`
	Description string    `json:"description" gorm:"size:1024"`
	Courses     []Course  `json:"courses" gorm:"many2many:path_courses;"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Instructor struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"size:255;not null"`
	Email     string    `json:"email" gorm:"size:255;unique;not null"`
	Biography string    `json:"biography" gorm:"size:1024"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Courses   []Course  `json:"courses" gorm:"many2many:course_instructors;constraint:OnDelete:CASCADE;"`
}

func CreateCourse(db *gorm.DB, course *Course) error {
	return db.Create(course).Error
}

func DeleteCourse(db *gorm.DB, courseID uint) error {
	return db.Where("id = ?", courseID).Delete(&Course{}).Error
}
