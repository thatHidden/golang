package gorm_postgres

import (
	"cleanstandarts/internal/domain"
	"database/sql"
	"gorm.io/gorm"
)

type pgImageRepository struct {
	Conn *gorm.DB
}

func NewPostgresImageRepository(conn *gorm.DB) domain.ImageRepository {
	return &pgImageRepository{conn}
}

//func (ar *pgAuctionRepository) Fetch(q *qstruct.QueryParams) (result []domain.Image, err error) {
//
//}

func (ir *pgImageRepository) GetByID(id uint64) (result domain.Image, err error) {
	err = ir.Conn.First(&result, "id = ?", id).Error
	return result, err
}

func (ir *pgImageRepository) GetByAucID(id uint64) (result []domain.Image, err error) {
	err = ir.Conn.Find(&result, "auction_id = ?", id).Error
	return result, err
}

func (ir *pgImageRepository) Delete(id uint64) (err error) {
	err = ir.Conn.Delete(&domain.Image{}, "id = ?", id).Error
	return err
}

func (ir *pgImageRepository) Create(a *domain.Image) (err error) {
	err = ir.Conn.Set("current_bid_id", sql.NullInt64{}).Create(a).Error
	return err
}

func (ir *pgImageRepository) MultipleCreate(a []domain.Image) (err error) {
	err = ir.Conn.Create(&a).Error
	return err
}
