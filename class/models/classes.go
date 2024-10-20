package models

import (
	"time"

	"gorm.io/gorm"
)

type Class struct {
	ID              uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Title           string    `json:"title" gorm:"size:255;not null"`
	Description     string    `json:"description" gorm:"size:1024"`
	CompanyID       uint      `json:"company_id" gorm:"not null"`
	CourseID        uint      `json:"course_id" gorm:"not null"`
	InstructorID    uint      `json:"instructor_id" gorm:"not null"`
	ScheduledAt     time.Time `json:"scheduled_at" gorm:"not null"`
	Duration        uint      `json:"duration" gorm:"not null"`
	MaxParticipants uint      `json:"max_participants" gorm:"not null"`
	CurrentEnrolled uint      `json:"current_enrolled" gorm:"default:0"`
	WaitlistEnabled bool      `json:"waitlist_enabled" gorm:"default:false"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func CreateClass(db *gorm.DB, class *Class) error {
	return db.Create(class).Error
}
