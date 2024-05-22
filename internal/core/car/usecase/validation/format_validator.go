package validation

import (
	"cleanstandarts/internal/domain/dto/models"
	"cleanstandarts/pkg/validator"
	"regexp"
)

var LocationRX = regexp.MustCompile("^[а-яА-ЯёЁ]{0,23}(?:-[^-]*){0,2},[ ]?[а-яА-ЯёЁ]{0,23}[ ]?район$")

type BidFormatValidator struct {
	Base *validator.BaseValidator
}

func NewCarFormatValidator(base *validator.BaseValidator) validator.InterfaceValidator[models.CarInputDTO] {
	return &BidFormatValidator{
		Base: base,
	}
}

func (bv *BidFormatValidator) Validate(c *models.CarInputDTO) {
	bv.CheckBaseCarID(c.BaseCarID)
	bv.CheckMileage(c.Mileage)
	bv.CheckVin(c.Vin)
	bv.CheckLocation(c.Location)
	bv.CheckExteriorColor(c.ExteriorColor)
	bv.CheckInteriorColor(c.InteriorColor)
}

func (bv *BidFormatValidator) CheckBaseCarID(id uint64) {
	bv.Base.Check(id != 0, "base_car_id", "must be provided")
	bv.Base.Check(id > 0, "base_car_id", "must be valid id")
}

func (bv *BidFormatValidator) CheckMileage(m uint32) {
	bv.Base.Check(m != 0, "mileage", "must be provided")
	bv.Base.Check(m < 10000000, "mileage", "must be less then 10mil")
}

func (bv *BidFormatValidator) CheckVin(v string) {
	bv.Base.Check(len(v) != 0, "vin", "must be provided")
	bv.Base.Check(len(v) == 17, "vin", "must be 17 characters")
}

func (bv *BidFormatValidator) CheckLocation(l string) {
	bv.Base.Check(len(l) != 0, "location", "must be provided")
	bv.Base.Check(bv.Base.Matches(l, LocationRX), "vin", "must be valid location: '<city>, <rayon> район")
}

func (bv *BidFormatValidator) CheckExteriorColor(e string) {
	bv.Base.Check(len(e) != 0, "exterior_color", "must be provided")
}

func (bv *BidFormatValidator) CheckInteriorColor(i string) {
	bv.Base.Check(len(i) != 0, "interior_color", "must be provided")
}
