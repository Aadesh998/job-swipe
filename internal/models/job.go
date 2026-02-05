package models

import (
	"gorm.io/gorm"
)

type Job struct {
	gorm.Model
	CompanyID       uint   `gorm:"not null;index" json:"company_id"`
	Title           string `gorm:"not null" json:"title"`
	Description     string `gorm:"not null" json:"description"`
	Requirements    string `json:"requirements"` // Skills, etc.
	YearsOfExperience string `json:"years_of_experience"` // e.g., "0-1", "2-5"
	Field           string `json:"field"` // e.g., "Software Engineering", "Marketing"
	Location        string `json:"location"`
	Type            string `json:"type"` // e.g., "Full-time", "Part-time", "Internship"
	Stipend         string `json:"stipend"` // Used for internships
	SalaryRange     string `json:"salary_range"` // e.g. "$50k-$70k"
	Status          string `gorm:"default:'open'" json:"status"` // open, closed
	Company         Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}
