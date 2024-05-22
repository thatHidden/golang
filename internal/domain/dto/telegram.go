package dto

import (
	"cleanstandarts/internal/domain"
	"cleanstandarts/internal/domain/dto/models"
)

func TelegramToOutputDto(t *domain.Telegram) models.TelegramOutputDTO {
	return models.TelegramOutputDTO{
		ID:       uint64(t.ID),
		UserID:   t.UserID,
		ChatID:   t.ChatID,
		Username: t.Username,
		DoAlerts: t.DoAlerts,
	}
}

func TelegramDtoToDomain(dto *models.CarModsInputDTO) domain.Telegram {
	return domain.Telegram{}
}
