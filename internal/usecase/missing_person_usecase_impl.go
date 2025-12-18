package usecase

import (
	"context"

	"github.com/Mhbib34/missing-person-service/internal/dto"
	"github.com/Mhbib34/missing-person-service/internal/exception"
	"github.com/Mhbib34/missing-person-service/internal/helper"
	"github.com/Mhbib34/missing-person-service/internal/model"
	"github.com/Mhbib34/missing-person-service/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type MissingPersonUsecaseImpl struct {
	repository repository.MissingPersonRepository
	Validate       *validator.Validate
}

func NewMissingPersonUsecase(repository repository.MissingPersonRepository, validate *validator.Validate) MissingPersonUsecase {
	return &MissingPersonUsecaseImpl{repository: repository, Validate: validate}
}

func (service *MissingPersonUsecaseImpl) Create(ctx context.Context, request dto.CreateMissingPersonRequest) (dto.MissingPersonResponse, error) {
	err := service.Validate.Struct(request)
	exception.PanicIfError(err)

	missingPerson := &model.MissingPersons{
		Name: request.Name, 
		Age: request.Age, 
		Description: request.Description,  
		LastSeen: request.LastSeen, 
		Contact: request.Contact, 
		PhotoID: request.Photo.Filename,
	}
	
	missingPerson, err = service.repository.Create(ctx, missingPerson)
	exception.PanicIfError(err)

	return helper.ToMissingPersonResponse(*missingPerson), err
}

func (service *MissingPersonUsecaseImpl) FindByID(ctx context.Context, id uuid.UUID) (*model.MissingPersons, error) {
	missingPerson, err := service.repository.FindByID(ctx, id)
	exception.PanicIfError(err)
	
	return missingPerson, nil
}

func (service *MissingPersonUsecaseImpl) GetAll(ctx context.Context, page int, limit int) ([]model.MissingPersons, int64, error) {
	missingPersons, total, err := service.repository.GetAll(ctx, page, limit)
	exception.PanicIfError(err)
	
	return missingPersons, total, nil
}