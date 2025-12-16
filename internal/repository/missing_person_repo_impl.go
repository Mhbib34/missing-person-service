package repository

import (
	"context"

	"github.com/Mhbib34/missing-person-service/internal/exception"
	"github.com/Mhbib34/missing-person-service/internal/model"
	"gorm.io/gorm"
)

type MissingPersonRepositoryImpl struct{
	db *gorm.DB
}

func NewMissingPersonRepository(db *gorm.DB) MissingPersonRepository {
	return &MissingPersonRepositoryImpl{ db : db }
}

func (r *MissingPersonRepositoryImpl) Create(ctx context.Context, missingPerson *model.MissingPersons) (*model.MissingPersons, error) {
	err := r.db.WithContext(ctx).Create(missingPerson).Error
	exception.PanicIfError(err)
	return missingPerson, nil
}