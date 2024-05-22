package validation

import (
	"cleanstandarts/internal/domain/dto/models"
	"cleanstandarts/pkg/validator"
)

var BodyStyleList = [8]string{"Coupe", "Convertible", "Hatchback", "Sedan", "SUV/Crossover", "Truck", "Van/Minivan", "Wagon"}
var TransmissionList = [2]string{"Automatic", "Manual"}
var DrivetrainList = [3]string{"AWD", "RWD", "FWD"}
var EngineTypeList = [4]string{"F", "I", "V", "W"}
var EngineCylinderList = [...]string{"2", "3", "4", "5", "6", "8", "10", "12", "16"}

type BidFormatValidator struct {
	Base *validator.BaseValidator
}

func NewCarModsFormatValidator(base *validator.BaseValidator) validator.InterfaceValidator[models.CarModsInputDTO] {
	return &BidFormatValidator{
		Base: base,
	}
}

func (bv *BidFormatValidator) Validate(c *models.CarModsInputDTO) {
	bv.validateDrivetrain(c.Drivetrain)
	bv.validateBodystyle(c.BodyStyle)
	bv.validateTransmission(c.Transmission)
	bv.validateTorque(c.Torque)
	bv.validateEngineType(c.EngineType)
	bv.validateEngineSize(c.EngineSize)
	bv.validateHorsePower(c.HorsePower)
}

func (bv *BidFormatValidator) validateEngineSize(engineSize float64) {
	bv.Base.Check(engineSize != 0.0, "engineType", "must be provided")
	bv.Base.Check(engineSize <= 10, "engineType", "must not be more than 10")
	bv.Base.Check(engineSize >= 0.6, "engineType", "must be at least 0.6")
}

func (bv *BidFormatValidator) validateEngineType(engineType string) {
	if engineType == "Rotor" {
		return
	}
	var minBytes = len(engineType) >= 2
	var maxBytes = len(engineType) <= 3
	bv.Base.Check(engineType != "", "engineType", "must be provided")
	bv.Base.Check(maxBytes, "engineType", "must not be more than 3 bytes long")
	bv.Base.Check(minBytes, "engineType", "must be at least 2 bytes long")
	if !(minBytes && maxBytes) {
		return
	}
	bv.Base.Check(bv.Base.In(engineType[0:1], EngineTypeList[:]...), "engineType", "must be a valid engineType")
	bv.Base.Check(bv.Base.In(engineType[1:], EngineCylinderList[:]...), "engineType", "must be a valid engineType")

}

func (bv *BidFormatValidator) validateHorsePower(horsePower uint16) {
	bv.Base.Check(horsePower != 0, "horsepower", "must be provided")
	bv.Base.Check(horsePower <= 5000, "horsepower", "must not be more than 5000")
	bv.Base.Check(horsePower >= 50, "horsepower", "must be at least 50")
}

func (bv *BidFormatValidator) validateTorque(torque uint16) {
	bv.Base.Check(torque != 0, "torque", "must be provided")
	bv.Base.Check(torque <= 3000, "torque", "must not be more than 3000")
	bv.Base.Check(torque >= 33, "torque", "must be at least 50")
}

func (bv *BidFormatValidator) validateBodystyle(bodystyle string) {
	bv.Base.Check(bodystyle != "", "bodystyle", "must be provided")
	bv.Base.Check(bv.Base.In(bodystyle, BodyStyleList[:]...), "bodystyle", "must be a valid bodystyle")
}

func (bv *BidFormatValidator) validateTransmission(transmission string) {
	bv.Base.Check(transmission != "", "transmission", "must be provided")
	bv.Base.Check(bv.Base.In(transmission, TransmissionList[:]...), "transmission", "must be a valid transmission")
}

func (bv *BidFormatValidator) validateDrivetrain(drivetrain string) {
	bv.Base.Check(drivetrain != "", "drivetrain", "must be provided")
	bv.Base.Check(bv.Base.In(drivetrain, DrivetrainList[:]...), "drivetrain", "must be a valid drivetrain")
}
