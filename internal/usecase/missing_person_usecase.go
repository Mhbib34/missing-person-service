package usecase

import (
	"context"

	"github.com/Mhbib34/missing-person-service/internal/dto"
)

type MissingPersonUsecase interface {
	Create(ctx context.Context, request dto.CreateMissingPersonRequest)(dto.MissingPersonResponse, error)
}