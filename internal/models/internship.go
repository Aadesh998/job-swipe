package models

import (
	"time"

	"gorm.io/gorm"
)

type Internship struct {
	gorm.Model
	JobSeekerProfileID uint      `gorm:"not null;index" json:"job_seeker_profile_id"`
	Company            string    `json:"company" binding:"required"`
	Role               string    `json:"role" binding:"required"`
	Description        string    `json:"description"`
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
	IsCurrent          bool      `json:"is_current"`
	Location           string    `json:"location"`
	Technologies       string    `json:"technologies"`
}
