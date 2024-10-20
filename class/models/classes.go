package models

import (
	"time"

	"gorm.io/gorm"
)

type ClassType struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"size:100;not null;unique"`
	Description string    `json:"description" gorm:"size:1024"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

var defaultClassTypes = []ClassType{
	{Name: "Webinar", Description: "Live online presentations."},
	{Name: "Demo", Description: "Product demonstrations."},
	{Name: "Tutorial", Description: "Step-by-step instructional sessions."},
	{Name: "Masterclass", Description: "In-depth sessions by experts."},
}

type Class struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Title        string    `json:"title" gorm:"size:255;not null"`
	Description  string    `json:"description" gorm:"size:1024"`
	CompanyID    uint      `json:"company_id" gorm:"not null"`
	CourseID     uint      `json:"course_id" gorm:"not null"`
	InstructorID uint      `json:"instructor_id" gorm:"not null"`
	ScheduledAt  time.Time `json:"scheduled_at" gorm:"not null"`

	Duration        uint       `json:"duration" gorm:"not null"`
	MaxParticipants uint       `json:"max_participants" gorm:"not null"`
	CurrentEnrolled uint       `json:"current_enrolled" gorm:"default:0"`
	WaitlistEnabled bool       `json:"waitlist_enabled" gorm:"default:false"`
	ClassTypeID     uint       `json:"class_type_id" gorm:"not null"`
	ClassType       *ClassType `gorm:"foreignKey:ClassTypeID"`
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

func CreateClass(db *gorm.DB, class *Class) error {
	return db.Create(class).Error
}

func MigrateDefaultClassTypes(db *gorm.DB) error {
	for _, classType := range defaultClassTypes {
		var count int64
		// Check if the class type already exists
		if err := db.Model(&ClassType{}).Where("name = ?", classType.Name).Count(&count).Error; err != nil {
			return err
		}
		// Insert if it does not exist
		if count == 0 {
			if err := db.Create(&classType).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
