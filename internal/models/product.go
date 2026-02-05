package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	CompanyID   uint    `gorm:"not null;index" json:"company_id"`
	Name        string  `gorm:"not null" json:"name"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	Price       float64 `json:"price,omitempty"`
}
