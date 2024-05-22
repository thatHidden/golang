package domain

import (
	"cleanstandarts/internal/domain/dto/models"
	"gorm.io/gorm"
)

type Telegram struct {
	gorm.Model
	UserID   uint64
	Username string `gorm:"unique; not null" json:"username"`
	ChatID   uint64 `gorm:"unique; not null" json:"chat_id"`
	DoAlerts bool   `gorm:"not null" json:"do_alerts"`
}

type TelegramUsecase interface {
	SetChatID(username string, cID uint64) (err error)
	Fetch() (result []models.TelegramOutputDTO, err error)
	Update(id uint64, val bool) (err error)
	GetByID(id uint64) (result models.TelegramOutputDTO, err error)
	Create(dto *models.TelegramInputDTO) (errors map[string]string, ok bool)
	Delete(id uint64) (err error)
}

type TelegramRepository interface {
	GetChatsID(uID []uint64) (result []uint64, err error)
	GetByUserID(id uint64) (result Telegram, err error)
	GetByUsername(username string) (result Telegram, err error)
	Update(id uint64, val bool) (err error)
	SetChatID(username string, cID uint64) (err error)
	Fetch() (result []Telegram, err error)
	GetByID(id uint64) (result Telegram, err error)
	Create(u *Telegram) (err error)
	Delete(id uint64) (err error)
}
