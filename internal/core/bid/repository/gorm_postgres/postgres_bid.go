package gorm_postgres

import (
	"cleanstandarts/internal/domain"
	"database/sql"
	"gorm.io/gorm"
)

type pgBidRepository struct {
	Conn *gorm.DB
}

func NewPostgresBidRepository(conn *gorm.DB) domain.BidRepository {
	return &pgBidRepository{conn}
}

func (br *pgBidRepository) GetParticipantsByID(auctionID uint64) (result []uint64, err error) {
	err = br.Conn.Model(&domain.Bid{}).
		Where("auction_id = ?", auctionID).
		Distinct("user_id").
		Pluck("user_id", &result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (br *pgBidRepository) GetMaxBidPrice(auctionID uint64) (result uint64, err error) {
	var temp sql.NullInt64
	err = br.Conn.Model(&domain.Bid{}).Select("MAX(price)").Where("auction_id = ?", auctionID).Scan(&temp).Error
	if !temp.Valid {
		return 0, err
	}
	result = uint64(temp.Int64)
	return result, err
}

func (br *pgBidRepository) Fetch(userID uint64, auctionID uint64) (result []domain.Bid, err error) {
	query := br.Conn.Model(&domain.Bid{})

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	if auctionID > 0 {
		query = query.Where("auction_id = ?", auctionID)
	}

	err = query.Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (br *pgBidRepository) GetByID(id uint64) (result domain.Bid, err error) {
	err = br.Conn.First(&result, "id = ?", id).Error
	if err != nil {
		return domain.Bid{}, err
	}
	return result, nil
}

//func (br *pgBidRepository) GetPriceByID(id uint64) (result uint64, err error) {
//	err = br.Conn.Model(&domain.Bid{}).Select("price").Where("id = ?", id).Scan(&result).Error
//	if err!= nil {
//		return 0, err
//	}
//	return result, nil
//}

func (br *pgBidRepository) Create(c *domain.Bid) (err error) {
	err = br.Conn.Create(c).Error
	return err
}

func (br *pgBidRepository) Delete(id uint64) (err error) {
	err = br.Conn.Delete(&domain.Bid{}, "id = ?", id).Error
	return err
}
