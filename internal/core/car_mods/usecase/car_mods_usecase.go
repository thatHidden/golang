package usecase

import (
	"cleanstandarts/internal/core/car_mods/usecase/validation"
	"cleanstandarts/internal/domain"
	"cleanstandarts/internal/domain/dto"
	"cleanstandarts/internal/domain/dto/models"
	"cleanstandarts/pkg/validator"
)

type carModsUsecase struct {
	carModsRepo domain.CarModsRepository
}

func NewCarModsUsecase(cmr domain.CarModsRepository) domain.CarModsUsecase {
	return &carModsUsecase{
		carModsRepo: cmr,
	}
}

func (cmu *carModsUsecase) validateInputFormat(base *validator.BaseValidator, c *models.CarModsInputDTO) {
	v := validation.NewCarModsFormatValidator(base)
	v.Validate(c)
}

func (cmu *carModsUsecase) validateCarMods(c *models.CarModsInputDTO) (errors map[string]string, ok bool) {
	base := validator.NewBaseValidator()
	cmu.validateInputFormat(base, c)
	if !base.Valid() {
		return base.Errors, false
	}
	return errors, true
}

func (cmu *carModsUsecase) Fetch() (result []models.CarModsOutputDTO, err error) {
	raw, err := cmu.carModsRepo.Fetch()
	if err != nil {
		return nil, err
	}
	for _, mod := range raw {
		DTO := dto.CarModsToOutputDto(&mod)
		result = append(result, DTO)
	}
	return result, err
}

func (cmu *carModsUsecase) GetByID(id uint64) (result models.CarModsOutputDTO, err error) {
	raw, err := cmu.carModsRepo.GetByID(id)
	if err != nil {
		return result, err
	}
	result = dto.CarModsToOutputDto(raw)
	return result, err
}

func (cmu *carModsUsecase) Delete(id uint64) (err error) {
	err = cmu.carModsRepo.Delete(id)
	return err
}

func (cmu *carModsUsecase) Create(c *models.CarModsInputDTO) (errors map[string]string, ok bool) {
	errors, ok = cmu.validateCarMods(c)
	if !ok {
		return errors, false
	}
	errors = map[string]string{}
	carMods := dto.CarModsDtoToDomain(c)
	err := cmu.carModsRepo.Create(&carMods)
	if err != nil {
		errors["internal"] = err.Error()
	}
	c.Id = uint64(carMods.ID)
	return errors, true
}
