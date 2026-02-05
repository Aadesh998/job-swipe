package models

import (
	"gorm.io/gorm"
)

type JobSwipe struct {
	gorm.Model
	UserID uint `gorm:"not null;index" json:"user_id"`
	JobID  uint `gorm:"not null;index" json:"job_id"`
	Action string `gorm:"not null" json:"action"` // "like" (apply) or "pass"
}
