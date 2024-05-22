package dto

import (
	"cleanstandarts/internal/domain"
	"cleanstandarts/internal/domain/dto/models"
	"fmt"
	"strconv"
)

func CarModsToOutputDto(m *domain.CarMods) models.CarModsOutputDTO {
	return models.CarModsOutputDTO{
		Id:        uint64(m.ID),
		BodyStyle: m.BodyStyle,
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
			Transmission: m.Transmission,
			Drivetrain:   m.Drivetrain,
			Engine: struct {
				EngineSize string `json:"size"`
				EngineType string `json:"type"`
				HorsePower string `json:"power"`
				Torque     string `json:"torque"`
			}{
				EngineSize: fmt.Sprintf("%.1f л³", m.EngineSize),
				EngineType: m.EngineType,
				HorsePower: strconv.Itoa(int(m.HorsePower)) + " л.с.",
				Torque:     strconv.Itoa(int(m.Torque)) + " Нм",
			},
		},
	}
}

func CarModsDtoToDomain(dto *models.CarModsInputDTO) domain.CarMods {
	return domain.CarMods{
		BodyStyle:    dto.BodyStyle,
		Transmission: dto.Transmission,
		Drivetrain:   dto.Drivetrain,
		EngineSize:   dto.EngineSize,
		EngineType:   dto.EngineType,
		HorsePower:   dto.HorsePower,
		Torque:       dto.Torque,
	}
}
