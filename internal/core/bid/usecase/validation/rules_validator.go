package validation

import (
	"cleanstandarts/internal/domain"
	"cleanstandarts/pkg/validator"
)

type BidValidator struct {
	Base        *validator.BaseValidator
	paymentRepo domain.PaymentRepository
}

func NewBidRulesValidator(base *validator.BaseValidator) validator.InterfaceValidator[string] {
	return &BidValidator{
		Base: base,
	}
}

func (bv *BidValidator) Validate(id *string) {
	bv.CheckListed(*id)
	bv.CheckExists(*id)
}

func (bv *BidValidator) CheckListed(id string) {
	//isListed, err := bv.paymentRepo.IsListedByPaymentID(id)
	//if err != nil {
	//	bv.Base.Check(false, "internal", err.Error())
	//	return
	//}
	//bv.Base.Check(isListed != true, "bid", "must be not listed already")
}

func (bv *BidValidator) CheckExists(id string) {
	//isExists, err := bv.paymentRepo.IsExistsByPaymentID(id)
	//if err != nil {
	//	bv.Base.Check(false, "internal", err.Error())
	//	return
	//}
	//bv.Base.Check(isExists == true, "bid", "must have exists payment")
}
