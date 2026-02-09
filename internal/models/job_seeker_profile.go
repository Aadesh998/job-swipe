package models

import (
	"gorm.io/gorm"
)

type JobSeekerProfile struct {
	gorm.Model
	UserID       uint         `gorm:"not null;uniqueIndex" json:"user_id"`
	FirstName    string       `json:"first_name"`
	LastName     string       `json:"last_name"`
	ResumeURL    string       `json:"resume_url"`
	Skills       string       `json:"skills"`
	Experience   string       `json:"experience"`
	Education      string         `json:"education"`
	Bio            string         `json:"bio"`
	JobPreferences []string       `gorm:"type:jsonb;serializer:json" json:"job_preferences"`
	IsOpenToWork   bool           `gorm:"default:true" json:"is_open_to_work"`
	Internships  []Internship `json:"internships,omitempty" gorm:"foreignKey:JobSeekerProfileID"`
	User         User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
