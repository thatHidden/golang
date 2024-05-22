package validation

import (
	"cleanstandarts/internal/domain"
	"cleanstandarts/pkg/validator"
	"fmt"
)

type PaymentRulesValidator struct {
	Base        *validator.BaseValidator
	auctionRepo domain.AuctionRepository
	bidRepo     domain.BidRepository
}

func NewPaymentRulesValidator(base *validator.BaseValidator, ar domain.AuctionRepository, br domain.BidRepository) validator.InterfaceValidator[domain.Payment] {
	return &PaymentRulesValidator{
		Base:        base,
		auctionRepo: ar,
		bidRepo:     br,
	}
}

func (bv *PaymentRulesValidator) Validate(p *domain.Payment) {
	bv.CheckSeller(p.AuctionID, p.UserID)
	bv.CheckPrice(p.AuctionID, p.Price)
}

func (bv *PaymentRulesValidator) CheckPrice(aID, price uint64) {
	actualPrice, err := bv.bidRepo.GetMaxBidPrice(aID)
	if err != nil {
		bv.Base.Check(false, "internal", err.Error())
		return
	}
	bv.Base.Check(actualPrice < price, "price", "must not be lower then actual")
}

func (bv *PaymentRulesValidator) CheckSeller(aID, uID uint64) {
	sellerId, err := bv.auctionRepo.GetSellerId(aID)
	fmt.Println(aID)
	if err != nil {
		bv.Base.Check(false, "internal", err.Error())
		return
	}
	bv.Base.Check(sellerId != uID, "user_id", "must not be seller")
}
