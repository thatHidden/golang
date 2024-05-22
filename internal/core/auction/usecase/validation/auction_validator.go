package validation

import (
	"cleanstandarts/internal/domain/dto/models"
	"cleanstandarts/pkg/validator"
	"fmt"
	"time"
)

type AuctionValidator struct {
	Base *validator.BaseValidator
}

func NewAuctionValidator(base *validator.BaseValidator) validator.InterfaceValidator[models.AuctionInputDTO] {
	return &AuctionValidator{
		Base: base,
	}
}

func (au *AuctionValidator) Validate(a *models.AuctionInputDTO) {
	au.ValidateDate(a.DateEnd)
	au.ValidateReserve(a.Reserve)
}

func (au *AuctionValidator) ValidateReserve(r uint64) {
	fmt.Println(r)
	au.Base.Check(r >= 10000, "reserve", "must not be less then 10000")
}

func (au *AuctionValidator) ValidateDate(d time.Time) {
	au.Base.Check(d.After(time.Now()), "DateEnd", "must be a future")
}
