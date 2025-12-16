package model

import (
	"time"

	"github.com/google/uuid"
)

type ImageStatus string

const (
	Pending    ImageStatus = "pending"
	Processing ImageStatus = "processing"
	Ready      ImageStatus = "ready"
	Failed     ImageStatus = "failed"
)




type MissingPersons struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	// Personal Info
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Age         int    `gorm:"type:int" json:"age"`
	Description string `gorm:"type:text;not null" json:"description"`
	LastSeen    string `gorm:"type:varchar(255);not null" json:"last_seen"`
	Contact     string `gorm:"type:varchar(100);not null" json:"contact"`

	// Image Info
	PhotoID     string      `gorm:"type:varchar(255);not null" json:"photo_id"` // Cloudinary public_id
	ImageStatus ImageStatus `gorm:"type:varchar(20);default:'pending'" json:"image_status"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
}