package domain

import (
	"cleanstandarts/internal/domain/dto/models"
	"database/sql"
	"time"
)

type BaseCar struct {
	ID        uint         `gorm:"primarykey"`
	CreatedAt time.Time    `json:"-"`
	UpdatedAt time.Time    `json:"-"`
	DeletedAt sql.NullTime `gorm:"index" json:"-"`

	Brand      string `gorm:"not null" json:"brand"`
	Model      string `gorm:"not null" json:"model"`
	Generation string `gorm:"not null" json:"generation"`

	BodyStyle    string  `gorm:"not null" json:"body_style"`
	Transmission string  `gorm:"not null" json:"transmission"`
	Drivetrain   string  `gorm:"not null" json:"drivetrain"`
	EngineSize   float64 `gorm:"not null" json:"engine_size"`
	EngineType   string  `gorm:"not null" json:"engine_type"`
	HorsePower   uint16  `gorm:"not null" json:"horse_power"`
	Torque       uint16  `gorm:"not null" json:"torque"`

	BuildFrom uint16 `gorm:"not null" json:"build_from"`
	BuildTo   uint16 `gorm:"not null" json:"build_to"`
}

type BaseCarUsecase interface {
	Fetch() (result []models.CarBaseOutputDTO, err error)
	GetByID(id uint64) (result models.CarBaseOutputDTO, err error)
	Create(c *models.CarBaseInputDTO) (errors map[string]string, ok bool)
	Delete(id uint64) (err error)
}

type BaseCarRepository interface {
	Fetch() (result []BaseCar, err error)
	GetByID(id uint64) (result BaseCar, err error)
	Create(c *BaseCar) (err error)
	Delete(id uint64) (err error)
}
