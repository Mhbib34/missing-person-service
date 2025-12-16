package entity

import (
	"time"

	"github.com/Mhbib34/missing-person-service/internal/model"
	"github.com/google/uuid"
)

type MissingPersons struct {
	ID  uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string
	Age         int
	Description string
	LastSeen    string
	Contact     string
	PhotoID     string
	ImageStatus model.ImageStatus
	CreatedAt time.Time 
}