package usecase

import (
	"cleanstandarts/internal/core/telegram/usecase/validation"
	"cleanstandarts/internal/domain"
	"cleanstandarts/internal/domain/dto"
	"cleanstandarts/internal/domain/dto/models"
	"cleanstandarts/pkg/validator"
)

type telegramUsecase struct {
	telegramRepo domain.TelegramRepository
}

func NewTelegramUsecase(tr domain.TelegramRepository) domain.TelegramUsecase {
	return &telegramUsecase{
		telegramRepo: tr,
	}
}

func (tu *telegramUsecase) validateInputFormat(base *validator.BaseValidator, t *models.TelegramInputDTO) {
	v := validation.NewBidFormatValidator(base)
	v.Validate(t)
}

func (tu *telegramUsecase) validateBusinessRules(base *validator.BaseValidator, t *models.TelegramInputDTO) {
	v := validation.NewBidRulesValidator(base, tu.telegramRepo)
	v.Validate(t)
}

func (tu *telegramUsecase) SetChatID(username string, cID uint64) (err error) {
	err = tu.telegramRepo.SetChatID(username, cID)
	return
}

func (tu *telegramUsecase) Update(uID uint64, val bool) (err error) {
	err = tu.telegramRepo.Update(uID, val)
	return
}

func (tu *telegramUsecase) validateTelegram(t *models.TelegramInputDTO) (errors map[string]string, ok bool) {
	base := validator.NewBaseValidator()
	tu.validateInputFormat(base, t)
	if !base.Valid() {
		return base.Errors, false
	}
	tu.validateBusinessRules(base, t)
	if !base.Valid() {
		return base.Errors, false
	}
	return errors, true
}

func (tu *telegramUsecase) Fetch() (result []models.TelegramOutputDTO, err error) {
	raw, err := tu.telegramRepo.Fetch()
	for _, val := range raw {
		DTO := dto.TelegramToOutputDto(&val)
		result = append(result, DTO)
	}
	return result, err
}

func (tu *telegramUsecase) GetByID(id uint64) (result models.TelegramOutputDTO, err error) {
	raw, err := tu.telegramRepo.GetByID(id)
	result = dto.TelegramToOutputDto(&raw)
	return result, err
}

func (tu *telegramUsecase) Delete(id uint64) (err error) {
	err = tu.telegramRepo.Delete(id)
	return err
}

func (tu *telegramUsecase) Create(dto *models.TelegramInputDTO) (errors map[string]string, ok bool) {
	errors, ok = tu.validateTelegram(dto)
	if !ok {
		return errors, false
	}
	errors = map[string]string{}
	tg := domain.Telegram{
		Username: dto.Username,
		UserID:   dto.UserID,
	}
	err := tu.telegramRepo.Create(&tg)
	if err != nil {
		errors["internal"] = err.Error()
		return errors, false
	}

	return errors, true
}
