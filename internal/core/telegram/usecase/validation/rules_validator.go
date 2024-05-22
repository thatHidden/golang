package validation

import (
	"cleanstandarts/internal/domain"
	"cleanstandarts/internal/domain/dto/models"
	"cleanstandarts/pkg/validator"
	"errors"
	"gorm.io/gorm"
)

type BidValidator struct {
	Base         *validator.BaseValidator
	telegramRepo domain.TelegramRepository
}

func NewBidRulesValidator(base *validator.BaseValidator, tr domain.TelegramRepository) validator.InterfaceValidator[models.TelegramInputDTO] {
	return &BidValidator{
		Base:         base,
		telegramRepo: tr,
	}
}

func (bv *BidValidator) Validate(t *models.TelegramInputDTO) {
	bv.CheckUsernameUnique(t.Username)
	bv.CheckOnlyUsername(t.UserID)
}

func (bv *BidValidator) CheckUsernameUnique(username string) {
	_, err := bv.telegramRepo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return
		} else {
			bv.Base.Check(false, "internal", err.Error())
		}
	}
	bv.Base.Check(false, "username", "must be unique")
}

func (bv *BidValidator) CheckOnlyUsername(id uint64) {
	_, err := bv.telegramRepo.GetByUserID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return
		} else {
			bv.Base.Check(false, "internal", err.Error())
		}
	}
	bv.Base.Check(false, "username", "must be only for user")
}
