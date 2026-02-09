package models

import (
	"gorm.io/gorm"
)

type Job struct {
	gorm.Model
	CompanyID         uint    `gorm:"not null;index" json:"company_id"`
	Title             string  `gorm:"not null" json:"title"`
	Description       string  `gorm:"not null" json:"description"`
	Requirements      string  `json:"requirements"`
	YearsOfExperience string  `json:"years_of_experience"`
	Field             string  `json:"field"`
	Location          string  `json:"location"`
	Type              string  `json:"type"`
	Stipend           string  `json:"stipend"`
	SalaryRange       string  `json:"salary_range"`
	Status            string  `gorm:"default:'open'" json:"status"`
	Company           Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}
