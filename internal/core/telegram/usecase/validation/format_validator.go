package validation

import (
	"cleanstandarts/internal/domain/dto/models"
	"cleanstandarts/pkg/validator"
	"regexp"
)

var UserNameRX = regexp.MustCompile(`^[a-zA-Z]\w{4,31}$`)

type BidFormatValidator struct {
	Base *validator.BaseValidator
}

func NewBidFormatValidator(base *validator.BaseValidator) validator.InterfaceValidator[models.TelegramInputDTO] {
	return &BidFormatValidator{
		Base: base,
	}
}

func (bv *BidFormatValidator) Validate(t *models.TelegramInputDTO) {
	bv.CheckUsername(t.Username)
}

func (bv *BidFormatValidator) CheckUsername(u string) {
	bv.Base.Check(u != "", "username", "must be provided")
	bv.Base.Check(bv.Base.Matches(u, UserNameRX), "username", "must be a valid username")
}
