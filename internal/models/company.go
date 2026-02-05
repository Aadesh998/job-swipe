package models

import (
	"gorm.io/gorm"
)

type Company struct {
	gorm.Model
	UserID            uint   `gorm:"not null;index" json:"user_id"`
	CompanyName       string `gorm:"not null" json:"company_name"`
	CompanySize       string `json:"company_size"` // e.g., "1-10", "11-50", "50+"
	Location          string `json:"location"`
	Website           string `json:"website"`
	Description       string `json:"description"`
	Industry          string `json:"industry"` // e.g., "IT", "Healthcare"
	Services          string `json:"services"` // Could be a comma-separated list or JSON
	Products          []Product `json:"products,omitempty" gorm:"foreignKey:CompanyID"`
	LogoURL           string `json:"logo_url"`
	IsProfileComplete bool   `gorm:"default:false" json:"is_profile_complete"`
}
