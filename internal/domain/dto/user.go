package dto

import (
	"cleanstandarts/internal/domain"
	"cleanstandarts/internal/domain/dto/models"
)

func UserToOutputDto(u *domain.User) models.UserOutputDTO {
	return models.UserOutputDTO{
		Id:       uint64(u.ID),
		Username: u.Username,
		Photo:    u.Photo,
		//Email:       u.Email,
		//Phone:       u.Phone,
		Name: u.Name,
		//IsActivated: u.IsActivated,
	}
}

func UserDtoToDomain(dto *models.UserInputDTO) domain.User {
	return domain.User{
		Email:    dto.Email,
		Password: dto.Password,
		Username: dto.Username,
	}
}
