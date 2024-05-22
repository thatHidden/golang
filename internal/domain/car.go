package domain

import (
	"cleanstandarts/internal/domain/dto/models"
	"gorm.io/gorm"
)

type Car struct {
	gorm.Model
	BaseCarID     uint64 `gorm:"not null; foreignkey:BaseCarID" json:"base_car_id"`
	CarModsID     uint64 `gorm:"not null; foreignkey:CarModsID" json:"car_mods_id"`
	Mileage       uint32 `gorm:"not null" json:"mileage"`
	Vin           string `gorm:"not null" json:"vin"`
	Location      string `gorm:"not null" json:"location"`
	ExteriorColor string `gorm:"not null" json:"exterior_color"`
	InteriorColor string `gorm:"not null" json:"interior_color"`
}

type CarUsecase interface {
	Fetch() (result []models.CarOutputDTO, err error)
	GetByID(id uint64) (result models.CarOutputDTO, err error)
	Create(c *models.CarInputDTO) (errors map[string]string, ok bool)
	Delete(id uint64) (err error)
}

type CarRepository interface {
	Fetch() (result []Car, err error)
	GetByID(id uint64) (result Car, err error)
	Create(c *Car) (err error)
	Delete(id uint64) (err error)
}
