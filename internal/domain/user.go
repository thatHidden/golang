package domain

import (
	"cleanstandarts/internal/domain/dto/models"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username    string    `gorm:"unique; not null" json:"username"`
	Password    string    `gorm:"not null" json:"password"`
	Email       string    `gorm:"unique; not null" json:"email"`
	Phone       string    `json:"phone"`
	Photo       string    `json:"photo"`
	Name        string    `json:"name"`
	LastLogin   time.Time `json:"last_login"`
	TelegramID  uint64    `json:"telegram_id"` //gorm:"not null"
	IsStaff     bool      `gorm:"default:false" json:"is_staff"`
	IsSuperuser bool      `gorm:"default:false" json:"is_superuser"`
	IsActivated bool      `json:"is_activated"`
}

type UserUsecase interface {
	Fetch() (result []models.UserOutputDTO, err error)
	GetByID(id uint64) (result models.UserOutputDTO, err error)
	GetByIDRaw(id uint64) (result User, err error)
	Auth(email string, password string) (result string, err error)
	Activate(authUser User, tokenRequest string) (err error) //DTO
	Create(u *models.UserInputDTO) (errors map[string]string, ok bool)
	Delete(id uint64) (err error)
}

type UserRepository interface {
	Fetch() (result []User, err error)
	GetByID(id uint64) (result User, err error)
	GetByEmail(email string) (result User, err error)
	Create(u *User) (err error)
	Delete(id uint64) (err error)
	Activate(id uint64) (err error)
}
