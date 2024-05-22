package usecase

import (
	"cleanstandarts/internal/core/car/usecase/validation"
	"cleanstandarts/internal/domain"
	"cleanstandarts/internal/domain/dto"
	"cleanstandarts/internal/domain/dto/models"
	"cleanstandarts/pkg/validator"
)

type carUsecase struct {
	carRepo        domain.CarRepository
	carModsUsecase domain.CarModsUsecase
	baseCarUsecase domain.BaseCarUsecase
}

func NewCarUsecase(cr domain.CarRepository, cmu domain.CarModsUsecase, bcu domain.BaseCarUsecase) domain.CarUsecase {
	return &carUsecase{
		carRepo:        cr,
		carModsUsecase: cmu,
		baseCarUsecase: bcu,
	}
}

func (cu *carUsecase) validateInputFormat(base *validator.BaseValidator, c *models.CarInputDTO) {
	v := validation.NewCarFormatValidator(base)
	v.Validate(c)
}

func (cu *carUsecase) validateCar(c *models.CarInputDTO) (errors map[string]string, ok bool) {
	base := validator.NewBaseValidator()
	cu.validateInputFormat(base, c)
	if !base.Valid() {
		return base.Errors, false
	}
	return errors, true
}

func (cu *carUsecase) Fetch() (result []models.CarOutputDTO, err error) {
	raw, err := cu.carRepo.Fetch()
	if err != nil {
		return nil, err
	}
	for _, car := range raw {
		carMods, err := cu.carModsUsecase.GetByID(car.CarModsID)
		if err != nil {
			return result, err
		}
		carBase, err := cu.baseCarUsecase.GetByID(car.BaseCarID)
		if err != nil {
			return result, err
		}
		DTO := dto.CarToOutputDto(&car, &carBase, &carMods)
		result = append(result, DTO)
	}
	return result, err
}

func (cu *carUsecase) GetByID(id uint64) (result models.CarOutputDTO, err error) {
	raw, err := cu.carRepo.GetByID(id)
	if err != nil {
		return result, err
	}
	carMods, err := cu.carModsUsecase.GetByID(raw.CarModsID)
	if err != nil {
		return result, err
	}
	carBase, err := cu.baseCarUsecase.GetByID(raw.BaseCarID)
	if err != nil {
		return result, err
	}
	result = dto.CarToOutputDto(&raw, &carBase, &carMods)
	return result, err
}

func (cu *carUsecase) Delete(id uint64) (err error) {
	err = cu.carRepo.Delete(id)
	return err
}

func (cu *carUsecase) Create(c *models.CarInputDTO) (errors map[string]string, ok bool) {
	errors, ok = cu.validateCar(c)
	if !ok {
		return errors, false
	}
	errors = map[string]string{}
	car := dto.CarDtoToDomain(c)
	err := cu.carRepo.Create(&car)
	if err != nil {
		errors["internal"] = err.Error()
		return errors, false
	}
	c.Id = uint64(car.ID)
	return errors, true
}
