package validation

import (
	"cleanstandarts/internal/domain"
	"cleanstandarts/pkg/validator"
)

type PaymentFormatValidator struct {
	Base *validator.BaseValidator
}

func NewPaymentFormatValidator(base *validator.BaseValidator) validator.InterfaceValidator[domain.Payment] {
	return &PaymentFormatValidator{
		Base: base,
	}
}

func (bv *PaymentFormatValidator) Validate(p *domain.Payment) {
	bv.checkPrice(p.Price)
	bv.checkAuction(p.AuctionID)
}

func (bv *PaymentFormatValidator) checkAuction(aID uint64) {
	bv.Base.Check(aID != 0, "auction_id", "must be provided")
}

func (bv *PaymentFormatValidator) checkPrice(p uint64) {
	bv.Base.Check(p != 0, "price", "must be provided")
}
