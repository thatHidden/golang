package gorm_postgres

import (
	"cleanstandarts/internal/domain"
	"gorm.io/gorm"
)

type pgBaseCarRepository struct {
	Conn *gorm.DB
}

func NewPostgresCarRepository(conn *gorm.DB) domain.BaseCarRepository {
	return &pgBaseCarRepository{conn}
}

func (bcr *pgBaseCarRepository) Fetch() (result []domain.BaseCar, err error) {
	err = bcr.Conn.Find(&result).Error
	return result, err
}

func (bcr *pgBaseCarRepository) Create(c *domain.BaseCar) (err error) {
	err = bcr.Conn.Create(c).Error
	return err
}

func (bcr *pgBaseCarRepository) GetByID(id uint64) (result domain.BaseCar, err error) {
	err = bcr.Conn.First(&result, "id = ?", id).Error
	return result, err
}

func (bcr *pgBaseCarRepository) Delete(id uint64) (err error) {
	err = bcr.Conn.Delete(&domain.BaseCar{}, "id = ?", id).Error
	return err
}
