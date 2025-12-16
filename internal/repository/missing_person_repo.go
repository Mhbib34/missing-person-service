package repository

import (
	"context"

	"github.com/Mhbib34/missing-person-service/internal/model"
)

type MissingPersonRepository interface {
	Create(ctx context.Context, missingPerson *model.MissingPersons)(*model.MissingPersons, error)
}