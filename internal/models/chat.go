package models

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	SenderID   uint   `json:"sender_id"`
	ReceiverID uint   `json:"receiver_id"`
	Content    string `json:"content"`
	IsRead     bool   `gorm:"default:false" json:"is_read"`
}
