package models

import (
	"gorm.io/gorm"
)

type JobProviderProfile struct {
	gorm.Model
	UserID        uint   `gorm:"not null;uniqueIndex" json:"user_id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Title         string `json:"title"`
	ContactNumber string `json:"contact_number"`
	Bio           string `json:"bio"`
	User          User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
