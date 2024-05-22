package domain

import (
	"cleanstandarts/internal/domain/dto/models"
	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	AuctionID uint64 `gorm:"not null; foreignkey:AuctionID" json:"auction_id"`
	Mark      string `gorm:"not null" json:"mark"`
	Path      string `gorm:"not null" json:"path"`
}

type ImageUsecase interface {
	Fetch(userID uint64, auctionID uint64) (result []Image, err error)
	Create(c *models.ImagesInputDTO) (errors map[string]string, ok bool)
	Delete(id uint64) (err error)
}

type ImageRepository interface {
	//Fetch(userID uint64, auctionID uint64) (result []Image, err error)
	MultipleCreate(a []Image) (err error)
	GetByID(id uint64) (result Image, err error)
	GetByAucID(id uint64) (result []Image, err error)
	Create(c *Image) (err error)
	Delete(id uint64) (err error)
}
