package dto

import (
	"cleanstandarts/internal/domain"
	"cleanstandarts/internal/domain/dto/models"
	"time"
)

func AuctionToOutputDto(a *domain.Auction, c *models.CarOutputDTO, s *models.UserOutputDTO, w *models.UserOutputDTO,
	i []domain.Image, p uint64) models.AuctionOutputDTO {
	var interiorImages []string
	var exteriorImages []string
	for _, img := range i {
		if img.Mark == "interior" {
			interiorImages = append(interiorImages, img.Path)
		} else {
			exteriorImages = append(exteriorImages, img.Path)
		}
	}
	return models.AuctionOutputDTO{
		Id:  uint64(a.ID),
		Car: *c,
		Images: struct {
			Exterior []string `json:"exterior"`
			Interior []string `json:"interior"`
		}(struct {
			Exterior []string
			Interior []string
		}{Exterior: exteriorImages,
			Interior: interiorImages}),
		Seller:  *s,
		Winner:  *w,
		IsEnded: a.IsEnded,
		Reserve: a.Reserve != 0,
		Price:   p,
		DateEnd: a.DateEnd.Format(time.DateTime),
	}
}

func AuctionDtoToDomain(dto *models.AuctionInputDTO, sellerID uint64) domain.Auction {
	return domain.Auction{
		CarID:    dto.CarID,
		SellerID: sellerID,
		Reserve:  dto.Reserve,
		DateEnd:  dto.DateEnd,
	}
}
