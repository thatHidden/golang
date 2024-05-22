package usecase

import (
	"cleanstandarts/internal/core/basecar/usecase/validation"
	"cleanstandarts/internal/domain"
	"cleanstandarts/internal/domain/dto"
	"cleanstandarts/internal/domain/dto/models"
	"cleanstandarts/pkg/validator"
)

type baseCarUsecase struct {
	carRepo domain.BaseCarRepository
}

func NewBaseCarUsecase(cr domain.BaseCarRepository) domain.BaseCarUsecase {
	return &baseCarUsecase{
		carRepo: cr,
	}
}

func (bcu *baseCarUsecase) validateInputFormat(base *validator.BaseValidator, a *models.CarBaseInputDTO) {
	v := validation.NewBaseCarValidator(base)
	v.Validate(a)
}

//func (au *auctionUsecase) validateBusinessRules(a *domain.Auction) {
//
//}

func (bcu *baseCarUsecase) validateCar(a *models.CarBaseInputDTO) (errors map[string]string, ok bool) {
	ok = true
	base := validator.NewBaseValidator()
	bcu.validateInputFormat(base, a)
	if !base.Valid() {
		errors = base.Errors
		ok = false
	}
	return errors, ok
}

func (bcu *baseCarUsecase) Fetch() (result []models.CarBaseOutputDTO, err error) {
	raw, err := bcu.carRepo.Fetch()
	if err != nil {
		return nil, err
	}
	for _, car := range raw {
		DTO := dto.CarBaseToOutputDto(&car)
		result = append(result, DTO)
	}
	return result, err
}

func (bcu *baseCarUsecase) Create(c *models.CarBaseInputDTO) (errors map[string]string, ok bool) {
	errors, ok = bcu.validateCar(c)
	if !ok {
		return errors, false
	}
	errors = map[string]string{}
	carBase := dto.CarBaseDtoToDomain(c)
	err := bcu.carRepo.Create(&carBase)
	if err != nil {
		errors["error"] = err.Error()
		return errors, false
	}
	c.Id = uint64(carBase.ID) //удоли
	return errors, true
}

func (bcu *baseCarUsecase) GetByID(id uint64) (result models.CarBaseOutputDTO, err error) {
	raw, err := bcu.carRepo.GetByID(id)
	if err != nil {
		return result, err
	}
	result = dto.CarBaseToOutputDto(&raw)
	return result, err
}

func (bcu *baseCarUsecase) Delete(id uint64) (err error) {
	err = bcu.carRepo.Delete(id)
	return err
}
