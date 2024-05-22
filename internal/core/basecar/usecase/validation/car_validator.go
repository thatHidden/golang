package validation

import (
	"cleanstandarts/internal/domain/dto/models"
	"cleanstandarts/pkg/validator"
	"regexp"
	"time"
)

var CarNameRX = regexp.MustCompile("^[A-Za-z0-9\\s-]+$")

var BodyStyleList = [8]string{"Coupe", "Convertible", "Hatchback", "Sedan", "SUV/Crossover", "Truck", "Van/Minivan", "Wagon"}
var TransmissionList = [2]string{"Automatic", "Manual"}
var DrivetrainList = [3]string{"AWD", "RWD", "FWD"}
var EngineTypeList = [4]string{"F", "I", "V", "W"}
var EngineCylinderList = [...]string{"2", "3", "4", "5", "6", "8", "10", "12", "16"}

type BaseCarValidator struct {
	Base *validator.BaseValidator
}

func NewBaseCarValidator(base *validator.BaseValidator) validator.InterfaceValidator[models.CarBaseInputDTO] {
	return &BaseCarValidator{
		Base: base,
	}
}

func (cv *BaseCarValidator) Validate(car *models.CarBaseInputDTO) {
	cv.validateBrand(car.Brand)
	cv.validateModel(car.Model)
	cv.validateGeneration(car.Generation)
	cv.validateEngineSize(car.EngineSize)
	cv.validateEngineType(car.EngineType)
	cv.validateHorsePower(car.HorsePower)
	cv.validateTorque(car.Torque)
	cv.validateBuildYears(car.BuildFrom, car.BuildTo)
	cv.validateTransmission(car.Transmission)
	cv.validateBodystyle(car.BodyStyle)
	cv.validateDrivetrain(car.Drivetrain)
}

func In(value string, list ...string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
}

func (cv *BaseCarValidator) validateBrand(brand string) {
	cv.Base.Check(brand != "", "brand", "must be provided")
	cv.Base.Check(len(brand) <= 12, "brand", "must not be more than 10 bytes long")
	cv.Base.Check(len(brand) >= 2, "brand", "must be at least 2 bytes long")
	cv.Base.Check(cv.Base.Matches(brand, CarNameRX), "brand", "must be a valid brand")
}

func (cv *BaseCarValidator) validateModel(model string) {
	cv.Base.Check(model != "", "model", "must be provided")
	cv.Base.Check(len(model) <= 50, "model", "must not be more than 10 bytes long")
	cv.Base.Check(len(model) >= 2, "model", "must be at least 2 bytes long")
	cv.Base.Check(cv.Base.Matches(model, CarNameRX), "model", "must be a valid model")
}

func (cv *BaseCarValidator) validateGeneration(generation string) {
	cv.Base.Check(generation != "", "generation", "must be provided")
	cv.Base.Check(len(generation) <= 10, "generation", "must not be more than 10 bytes long")
	cv.Base.Check(len(generation) >= 1, "generation", "must be at least 2 bytes long")
}

func (cv *BaseCarValidator) validateEngineSize(engineSize float64) {
	cv.Base.Check(engineSize != 0.0, "engineType", "must be provided")
	cv.Base.Check(engineSize <= 10, "engineType", "must not be more than 10")
	cv.Base.Check(engineSize >= 0.6, "engineType", "must be at least 0.6")
}

func (cv *BaseCarValidator) validateEngineType(engineType string) {
	if engineType == "Rotor" {
		return
	}
	var minBytes = len(engineType) >= 2
	var maxBytes = len(engineType) <= 3
	cv.Base.Check(engineType != "", "engineType", "must be provided")
	cv.Base.Check(maxBytes, "engineType", "must not be more than 3 bytes long")
	cv.Base.Check(minBytes, "engineType", "must be at least 2 bytes long")
	if !(minBytes && maxBytes) {
		return
	}
	cv.Base.Check(In(engineType[0:1], EngineTypeList[:]...), "engineType", "must be a valid engineType")
	cv.Base.Check(In(engineType[1:], EngineCylinderList[:]...), "engineType", "must be a valid engineType")

}

func (cv *BaseCarValidator) validateHorsePower(horsePower uint16) {
	cv.Base.Check(horsePower != 0, "horsepower", "must be provided")
	cv.Base.Check(horsePower <= 5000, "horsepower", "must not be more than 5000")
	cv.Base.Check(horsePower >= 50, "horsepower", "must be at least 50")
}

func (cv *BaseCarValidator) validateTorque(torque uint16) {
	cv.Base.Check(torque != 0, "torque", "must be provided")
	cv.Base.Check(torque <= 3000, "torque", "must not be more than 3000")
	cv.Base.Check(torque >= 33, "torque", "must be at least 50")
}

func (cv *BaseCarValidator) validateBuildYears(buildFrom, buildTo uint16) {
	currentYear := time.Now().Year()
	cv.Base.Check(buildFrom != 0, "build_from", "must be provided")
	cv.Base.Check(buildTo != 0, "build_to", "must be provided")
	cv.Base.Check(buildFrom <= buildTo, "build_years", "must be a valid range")
	cv.Base.Check(buildTo <= uint16(currentYear)+uint16(10), "build_to", "last year must not be more then current + 10")
	cv.Base.Check(buildFrom >= uint16(1866), "build_from", "first year must be at least 1866")
	cv.Base.Check(buildTo-buildFrom <= 11, "build_years", "difference not be more than 11")
}

func (cv *BaseCarValidator) validateBodystyle(bodystyle string) {
	cv.Base.Check(bodystyle != "", "bodystyle", "must be provided")
	cv.Base.Check(In(bodystyle, BodyStyleList[:]...), "bodystyle", "must be a valid bodystyle")
}

func (cv *BaseCarValidator) validateTransmission(transmission string) {
	cv.Base.Check(transmission != "", "transmission", "must be provided")
	cv.Base.Check(In(transmission, TransmissionList[:]...), "transmission", "must be a valid transmission")
}

func (cv *BaseCarValidator) validateDrivetrain(drivetrain string) {
	cv.Base.Check(drivetrain != "", "drivetrain", "must be provided")
	cv.Base.Check(In(drivetrain, DrivetrainList[:]...), "drivetrain", "must be a valid drivetrain")
}
