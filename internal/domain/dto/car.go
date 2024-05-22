package dto

import (
	"cleanstandarts/internal/domain"
	"cleanstandarts/internal/domain/dto/models"
)

func CarToOutputDto(c *domain.Car, b *models.CarBaseOutputDTO, m *models.CarModsOutputDTO) models.CarOutputDTO {
	return models.CarOutputDTO{
		CarMods:  *m,
		CarBase:  *b,
		Mileage:  c.Mileage,
		Vin:      c.Vin,
		Location: c.Location,
		Color: struct {
			ExteriorColor string `json:"exterior"`
			InteriorColor string `json:"interior"`
		}{
			ExteriorColor: c.ExteriorColor,
			InteriorColor: c.InteriorColor,
		},
	}
}

func CarDtoToDomain(dto *models.CarInputDTO) domain.Car {
	return domain.Car{
		BaseCarID:     dto.BaseCarID,
		CarModsID:     dto.CarModsID,
		Mileage:       dto.Mileage,
		Vin:           dto.Vin,
		Location:      dto.Location,
		ExteriorColor: dto.ExteriorColor,
		InteriorColor: dto.InteriorColor,
	}
}
