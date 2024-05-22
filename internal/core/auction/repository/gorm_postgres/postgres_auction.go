package gorm_postgres

import (
	"cleanstandarts/internal/core/auction/repository/qstruct"
	"cleanstandarts/internal/domain"
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type pgAuctionRepository struct {
	Conn *gorm.DB
}

func NewPostgresAuctionRepository(conn *gorm.DB) domain.AuctionRepository {
	return &pgAuctionRepository{conn}
}

func (ar *pgAuctionRepository) GetWinnerPayIDByID(id uint64) (res string, err error) {
	err = ar.Conn.Model(&domain.Auction{}).
		Joins("JOIN bids ON auctions.current_bid_id = bids.id").
		Joins("JOIN payments ON bids.payment_id = payments.id").
		Select("pay_id").
		Where("auctions.id = ?", id).
		Scan(&res).
		Error
	if err != nil {
		return
	}
	fmt.Println("PAY_ID", res)
	return
}

// ?сразу забирать домен Bid
func (ar *pgAuctionRepository) GetLeadBidID(id uint64) (result uint64, err error) {
	err = ar.Conn.Model(domain.Auction{}).
		Select("current_bid_id").
		Where("id = ?", id).
		Scan(&result).
		Error
	return
}

func (ar *pgAuctionRepository) GetCarNameByID(id uint64) (result string, err error) {
	var ans struct {
		Brand      string
		Model      string
		Generation string
	}
	err = ar.Conn.Model(&domain.Auction{}).
		Joins("JOIN cars ON auctions.car_id = cars.id").
		Joins("JOIN base_cars ON cars.base_car_id = base_cars.id").
		Select("brand, model, generation").
		Where("auctions.id = ?", id).
		Scan(&ans).
		Error
	if err != nil {
		return "", err
	}
	//result = ans.Brand + " " + ans.Model + " (" + ans.Generation + ")"
	result = ans.Brand + " " + ans.Model
	if ans.Generation[0] != '_' {
		result += " (" + ans.Generation + ")"
	}
	return
}

func (ar *pgAuctionRepository) GetSellerId(aID uint64) (id uint64, err error) {
	err = ar.Conn.Model(&domain.Auction{}).Select("seller_id").Where("id = ?", aID).Scan(&id).Error
	return
}

func (ar *pgAuctionRepository) SetLeadBidID(aID, bID uint64) (err error) {
	err = ar.Conn.Model(&domain.Auction{}).
		Where("id = ?", aID).
		Update("current_bid_id", bID).Error
	return
}

func (ar *pgAuctionRepository) FinishByID(id uint64, wId uint64) (err error) {
	err = ar.Conn.Model(&domain.Auction{}).
		Where("id = ?", id).
		Update("winner_id", wId).
		Update("is_ended", true).Error
	return
}

func (ar *pgAuctionRepository) Fetch(q *qstruct.QueryParams) (result []domain.Auction, err error) {
	query := ar.Conn.Model(&domain.Auction{}).
		Joins("JOIN cars ON auctions.car_id = cars.id").
		Joins("JOIN base_cars ON cars.base_car_id = base_cars.id")

	if q.UserID > 0 {
		query = query.Where("user_id = ?", q.UserID)
	}

	if q.IsNoReserve == true {
		query = query.Where("reserve = ?", 0)
	}

	if q.IsEnded == true {
		query = query.Where("is_ended = ?", true)
	}

	if q.IsEndingSoon == true {
		query = query.Where("date_end <= ?", time.Now().Add(24*time.Hour))
	}

	if q.MinYear > 0 {
		query = query.Where("build_from = ?", q.MinYear)
	}

	if q.MaxYear > 0 {
		query = query.Where("build_to = ?", q.MaxYear)
	}

	if q.Brand != "" {
		query = query.Where("brand = ?", q.Brand)
	}

	if q.Model != "" {
		query = query.Where("model = ?", q.Model)
	}

	if q.Generation != "" {
		query = query.Where("generation = ?", q.Generation)
	}

	//?mods
	if q.BodyStyle != "" {
		query = query.Where("body_style = ?", q.BodyStyle)
	}

	err = query.Find(&result).Error
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, domain.ErrNotFound
	}

	return result, nil
}

func (ar *pgAuctionRepository) GetByID(id uint64) (result domain.Auction, err error) {
	err = ar.Conn.First(&result, "id = ?", id).Error
	return result, err
}

func (ar *pgAuctionRepository) Delete(id uint64) (err error) {
	err = ar.Conn.Delete(&domain.Auction{}, "id = ?", id).Error
	return err
}

func (ar *pgAuctionRepository) Create(a *domain.Auction) (err error) {
	err = ar.Conn.Set("current_bid_id", sql.NullInt64{}).Create(a).Error
	return err
}
