package models

import (
	"time"

	"gorm.io/gorm"
)

type Class struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string    `json:"title" gorm:"size:255;not null"`
	Description string    `json:"description" gorm:"size:1024"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func CreateClass(db *gorm.DB, class *Class) error {
	return db.Create(class).Error
}
