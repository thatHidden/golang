package models

import "time"

type AuctionInputDTO struct {
	Id      uint64    `json:"id"` //удоли
	CarID   uint64    `json:"car_id"`
	Reserve uint64    `json:"reserve"`
	DateEnd time.Time `json:"date_end"`
}

type AuctionOutputDTO struct {
	Id     uint64       `json:"id"`
	Car    CarOutputDTO `json:"car"`
	Images struct {
		Exterior []string `json:"exterior"`
		Interior []string `json:"interior"`
	} `json:"images"`
	Seller  UserOutputDTO `json:"seller"`
	Winner  UserOutputDTO `json:"winner"`
	Reserve bool          `json:"reserve"`
	Price   uint64        `json:"price"`
	IsEnded bool          `json:"is_ended"`
	DateEnd string        `json:"date_end"`
}

type ParticipantOutputDTO struct {
	AuctionName string   `json:"auction_name"`
	ChatIDs     []uint64 `json:"telegram_id"`
}
