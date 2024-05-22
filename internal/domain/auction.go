package domain

import (
	"cleanstandarts/internal/core/auction/repository/qstruct"
	"cleanstandarts/internal/domain/dto/models"
	"gorm.io/gorm"
	"time"
)

type Auction struct {
	gorm.Model
	CarID        uint64    `gorm:"not null; foreignkey:CarID" json:"car_id"`
	SellerID     uint64    `gorm:"not null; foreignkey:UserID" json:"seller_id"`
	WinnerID     uint64    `gorm:"default:0; foreignkey:UserID" json:"winner_id"`
	BidsAmount   uint64    `gorm:"default:0" json:"bids_amount"`
	CurrentBidID uint64    `gorm:"default:0; foreignkey:BidID" json:"current_bid_id"`
	Reserve      uint64    `gorm:"default:0" json:"reserve"`
	DateEnd      time.Time `gorm:"not null" json:"date_end"`
	IsEnded      bool      `gorm:"default:0" json:"is_ended"`
	Slug         string    `gorm:"not null" json:"slug"`
	//IsActive     bool      `gorm:"not null" json:"is_active"`
}

type AuctionUsecase interface {
	GetParticipantsByID(id uint64) (participants models.ParticipantOutputDTO, err error)
	FinishByID(id uint64) (err error)
	Fetch(q *qstruct.QueryParams) (result []models.AuctionOutputDTO, err error)
	GetByID(id uint64) (result models.AuctionOutputDTO, err error)
	Create(a *models.AuctionInputDTO, uID uint64) (errors map[string]string, ok bool)
	Delete(id uint64) (err error)
}

type AuctionRepository interface {
	GetWinnerPayIDByID(id uint64) (res string, err error)
	GetLeadBidID(id uint64) (result uint64, err error)
	GetCarNameByID(id uint64) (result string, err error)
	GetSellerId(aID uint64) (id uint64, err error)
	SetLeadBidID(auctionID, bidID uint64) (err error)
	FinishByID(id uint64, wId uint64) (err error)
	Fetch(q *qstruct.QueryParams) (result []Auction, err error)
	GetByID(id uint64) (result Auction, err error)
	Create(a *Auction) (err error)
	Delete(id uint64) (err error)
}

/*
ALTER TABLE auctions
ADD CONSTRAINT fk_current_bid_id
FOREIGN KEY (current_bid_id)
REFERENCES bids(id)
ON DELETE RESTRICT
ON UPDATE RESTRICT;

ALTER TABLE auctions
ADD CONSTRAINT fk_winner_id
FOREIGN KEY (winner_id)
REFERENCES users(id)
ON DELETE RESTRICT
ON UPDATE CASCADE;

ALTER TABLE auctions
ADD CONSTRAINT fk_seller_id
FOREIGN KEY (seller_id)
REFERENCES users(id)
ON DELETE RESTRICT
ON UPDATE CASCADE;

ALTER TABLE auctions
ADD CONSTRAINT fk_car_id
FOREIGN KEY (car_id)
REFERENCES cars(id)
ON DELETE RESTRICT
ON UPDATE CASCADE;
*/
