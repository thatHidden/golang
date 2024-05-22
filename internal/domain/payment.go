package domain

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	UserID         uint64 `gorm:"not null" json:"user_id"`
	AuctionID      uint64 `gorm:"not null" json:"auction_id"`
	Price          uint64 `gorm:"not null" json:"price"`
	IdempotenceKey string `gorm:"not null" json:"idempotence_key"`
	PayID          string `gorm:"not null" json:"pay_id"`
	Status         string `gorm:"not null" json:"status"`
	IsListed       bool   `gorm:"not null" json:"is_listed"`
}

type PaymentUsecase interface {
	Fetch() (result []Payment, err error)
	GetByID(id uint64) (result Payment, err error)
	GetByPaymentID(id string) (result Payment, err error)
	Create(c *Payment, u *User) (confUrl string, errors map[string]string, ok bool)
	Delete(id uint64) (err error)
}

type PaymentRepository interface {
	IsListedByPaymentID(id string) (result bool, err error)
	SetListedByPaymentID(id string) (err error)
	Fetch() (result []Payment, err error)
	IsExistsByPaymentID(id string) (result bool, err error)
	GetByPaymentID(id string) (result Payment, err error)
	GetByID(id uint64) (result Payment, err error)
	Create(c *Payment) (err error)
	Delete(id uint64) (err error)
}
