package dto

import (
	"cleanstandarts/internal/domain"
	"cleanstandarts/internal/domain/dto/models"
	"fmt"
	"strconv"
)

func CarBaseToOutputDto(c *domain.BaseCar) models.CarBaseOutputDTO {
	return models.CarBaseOutputDTO{
		Name:       c.Brand + " " + c.Model,
		Generation: c.Generation,
		BodyStyle:  c.BodyStyle,
		BuildYears: strconv.Itoa(int(c.BuildFrom)) + "-" + strconv.Itoa(int(c.BuildTo)),
		PowerTrain: struct {
			Transmission string `json:"transmission"`
			Drivetrain   string `json:"drivetrain"`
			Engine       struct {
				EngineSize string `json:"size"`
				EngineType string `json:"type"`
				HorsePower string `json:"power"`
				Torque     string `json:"torque"`
			} `json:"engine"`
		}{
			Transmission: c.Transmission,
			Drivetrain:   c.Drivetrain,
			Engine: struct {
				EngineSize string `json:"size"`
				EngineType string `json:"type"`
				HorsePower string `json:"power"`
				Torque     string `json:"torque"`
			}{
				EngineSize: fmt.Sprintf("%.1f л³", c.EngineSize),
				EngineType: c.EngineType,
				HorsePower: strconv.Itoa(int(c.HorsePower)) + " л.с.",
				Torque:     strconv.Itoa(int(c.Torque)) + " Нм",
			},
		},
	}
}

func CarBaseDtoToDomain(dto *models.CarBaseInputDTO) domain.BaseCar {
	return domain.BaseCar{
		Brand:        dto.Brand,
		Model:        dto.Model,
		Generation:   dto.Generation,
		BodyStyle:    dto.BodyStyle,
		Transmission: dto.Transmission,
		Drivetrain:   dto.Drivetrain,
		EngineSize:   dto.EngineSize,
		EngineType:   dto.EngineType,
		HorsePower:   dto.HorsePower,
		Torque:       dto.Torque,
		BuildFrom:    dto.BuildFrom,
		BuildTo:      dto.BuildTo,
	}
}
