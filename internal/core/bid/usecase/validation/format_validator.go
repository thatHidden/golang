package validation

import (
	"cleanstandarts/pkg/validator"
)

type BidFormatValidator struct {
	Base *validator.BaseValidator
}

func NewBidFormatValidator(base *validator.BaseValidator) validator.InterfaceValidator[string] {
	return &BidFormatValidator{
		Base: base,
	}
}

func (bv *BidFormatValidator) Validate(id *string) {
	bv.CheckPaymentID(*id)
}

func (bv *BidFormatValidator) CheckPaymentID(id string) {
	bv.Base.Check(len(id) != 0, "payment_id", "must be provided")
}
