package repository

import (
	"context"

	"github.com/Mhbib34/missing-person-service/internal/exception"
	"github.com/Mhbib34/missing-person-service/internal/model"
	"github.com/google/uuid"
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

func (r *MissingPersonRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*model.MissingPersons, error) {
	var missingPerson model.MissingPersons
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&missingPerson).Error
	exception.PanicIfError(err)
	return &missingPerson, nil
}

func (r *MissingPersonRepositoryImpl) GetAll(
	ctx context.Context,
	page int,
	limit int,
) ([]model.MissingPersons, int64, error) {

	var (
		missingPersons []model.MissingPersons
		total          int64
	)

	offset := (page - 1) * limit

	// hitung total data
	err := r.db.WithContext(ctx).
		Where("image_status = ?", "ready").
		Model(&model.MissingPersons{}).
		Count(&total).Error
		
	if err != nil {
		return nil, 0, err
	}

	// ambil data per page
	err = r.db.WithContext(ctx).
		Where("image_status = ?", "ready").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&missingPersons).Error
	if err != nil {
		return nil, 0, err
	}

	return missingPersons, total, nil
}
