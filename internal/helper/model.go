package helper

import (
	"github.com/Mhbib34/missing-person-service/internal/dto"
	"github.com/Mhbib34/missing-person-service/internal/model"
)

func ToMissingPersonResponse(user model.MissingPersons) dto.MissingPersonResponse {
	return dto.MissingPersonResponse{
		ID:          user.ID.String(),
		Name:        user.Name,
		Age:         user.Age,
		Description: user.Description,
		LastSeen:    user.LastSeen,
		Contact:     user.Contact,
		PhotoID:     user.PhotoID,
		ImageStatus: string(user.ImageStatus),
		CreatedAt:   user.CreatedAt.String(),
	}
}