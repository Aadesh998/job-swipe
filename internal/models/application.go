package models

import (
	"gorm.io/gorm"
)

type Application struct {
	gorm.Model
	JobID  uint   `gorm:"not null;index" json:"job_id"`
	UserID uint   `gorm:"not null;index" json:"user_id"`
	Status string `gorm:"default:'applied'" json:"status"` // applied, reviewing, interviewed, rejected, hired
	Job    Job    `gorm:"foreignKey:JobID" json:"job,omitempty"`
	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
