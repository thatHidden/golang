package domain

import (
	"cleanstandarts/internal/domain/dto/models"
	"gorm.io/gorm"
)

type CarMods struct {
	gorm.Model
	//CarID uint64 `json:"car_id"`
	BodyStyle    string  `json:"body_style"`
	Transmission string  `json:"transmission"`
	Drivetrain   string  `json:"drivetrain"`
	EngineSize   float64 `json:"engine_size"`
	EngineType   string  `json:"engine_type"`
	HorsePower   uint16  `json:"horse_power"`
	Torque       uint16  `json:"torque"`
}

type CarModsUsecase interface {
	Fetch() (result []models.CarModsOutputDTO, err error)
	GetByID(id uint64) (result models.CarModsOutputDTO, err error)
	Create(c *models.CarModsInputDTO) (errors map[string]string, ok bool)
	Delete(id uint64) (err error)
}

type CarModsRepository interface {
	Fetch() (result []CarMods, err error)
	GetByID(id uint64) (result *CarMods, err error)
	Create(c *CarMods) (err error)
	Delete(id uint64) (err error)
}
