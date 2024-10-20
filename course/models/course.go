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
	Instructors     []Instructor `json:"instructors" gorm:"many2many:course_instructors;"`
	EnrollmentLimit int          `json:"enrollment_limit"`
	Tags            []Tag        `json:"tags" gorm:"many2many:course_tags;"`
	SubCourses      []Course     `json:"sub_courses" gorm:"foreignKey:ParentCourseID"`
	ParentCourseID  *uint        `json:"parent_course_id"`
}

type Instructor struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"size:255;not null"`
	Email     string    `json:"email" gorm:"size:255;unique;not null"`
	Biography string    `json:"biography" gorm:"size:1024"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Courses   []Course  `json:"courses" gorm:"many2many:course_instructors;"`
}

func CreateCourse(db *gorm.DB, course *Course) error {
	return db.Create(course).Error
}
