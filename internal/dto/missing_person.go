package dto

import "mime/multipart"

type CreateMissingPersonRequest struct {
	Name        string                `form:"name" validate:"required"`
	Age         int                   `form:"age" validate:"required,gt=0"`
	Description string                `form:"description" validate:"required"`
	LastSeen    string                `form:"last_seen" validate:"required"`
	Contact     string                `form:"contact" validate:"required"`
	Photo       *multipart.FileHeader `form:"photo" validate:"required"`
}

type MissingPersonResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name,omitempty"`
	Age         int    `json:"age,omitempty"`
	Description string `json:"description,omitempty"`
	LastSeen    string `json:"last_seen,omitempty"`
	Contact     string `json:"contact,omitempty"`
	PhotoID     string `json:"photo_id,omitempty"`
	ImageStatus string `json:"image_status,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
}