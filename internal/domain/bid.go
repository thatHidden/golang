package domain

import "gorm.io/gorm"

type Bid struct {
	gorm.Model
	AuctionID uint64 `gorm:"not null;foreignkey:AuctionID" json:"auction_id"`
	UserID    uint64 `gorm:"not null;foreignkey:UserID" json:"user_id"`
	Price     uint64 `gorm:"not null" json:"price"`
	PaymentID uint64 `gorm:"not null" json:"pay_id"`
}

type BidUsecase interface {
	Fetch(userID uint64, auctionID uint64) (result []Bid, err error)
	GetByID(id uint64) (result Bid, err error)
	Create(payID string) (errors map[string]string, ok bool)
	Delete(id uint64) (err error)
}

type BidRepository interface {
	GetParticipantsByID(auctionID uint64) (result []uint64, err error)
	GetMaxBidPrice(auctionID uint64) (result uint64, err error)
	Fetch(userID uint64, auctionID uint64) (result []Bid, err error)
	GetByID(id uint64) (result Bid, err error)
	Create(c *Bid) (err error)
	Delete(id uint64) (err error)
}
