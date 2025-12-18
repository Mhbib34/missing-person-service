package repository

import (
	"context"

	"github.com/Mhbib34/missing-person-service/internal/model"
	"github.com/google/uuid"
)

type MissingPersonRepository interface {
	Create(ctx context.Context, missingPerson *model.MissingPersons)(*model.MissingPersons, error)
	FindByID(ctx context.Context, id uuid.UUID)(*model.MissingPersons, error)
	GetAll(ctx context.Context,page int, limit int) ([]model.MissingPersons, int64, error)
}