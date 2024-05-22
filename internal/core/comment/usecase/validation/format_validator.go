package validation

import (
	"cleanstandarts/internal/domain"
	"cleanstandarts/pkg/validator"
)

type BidFormatValidator struct {
	Base *validator.BaseValidator
}

func NewCommentFormatValidator(base *validator.BaseValidator) validator.InterfaceValidator[domain.Comment] {
	return &BidFormatValidator{
		Base: base,
	}
}

func (bv *BidFormatValidator) Validate(c *domain.Comment) {
	bv.validateAuctionID(c.AuctionID)
	bv.validateText(c.Text)
}

func (bv *BidFormatValidator) validateAuctionID(id uint64) {
	bv.Base.Check(id != 0, "auction_id", "must be provided")
}
func (bv *BidFormatValidator) validateText(t string) {
	bv.Base.Check(len(t) != 0, "text", "must be provided")
	bv.Base.Check(len(t) <= 200, "text", "must not be more than 200 bytes long")
}
