package usecase

import (
	"context"

	"github.com/Mhbib34/missing-person-service/internal/dto"
	"github.com/Mhbib34/missing-person-service/internal/model"
	"github.com/google/uuid"
)

type MissingPersonUsecase interface {
	Create(ctx context.Context, request dto.CreateMissingPersonRequest)(dto.MissingPersonResponse, error)
	FindByID(ctx context.Context, id uuid.UUID)(*model.MissingPersons, error)
	GetAll(ctx context.Context, page int, limit int)([]model.MissingPersons, int64, error)
}