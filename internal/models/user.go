package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email              string              `gorm:"uniqueIndex;not null" json:"email"`
	Password           string              `json:"-"`
	IsVerified         bool                `gorm:"default:false" json:"is_verified"`
	Role               string              `gorm:"default:'job_seeker'" json:"role"` // job_seeker, job_provider, admin
	VerificationToken  string              `json:"-"`
	PasswordResetToken string              `json:"-"`
	ResetTokenExpiry   *time.Time          `json:"-"`
	GoogleID           string              `gorm:"uniqueIndex" json:"google_id,omitempty"`
	AvatarURL          string              `json:"avatar_url,omitempty"`
	JobSeekerProfile   *JobSeekerProfile   `json:"job_seeker_profile,omitempty" gorm:"foreignKey:UserID"`
	JobProviderProfile *JobProviderProfile `json:"job_provider_profile,omitempty" gorm:"foreignKey:UserID"`
}
