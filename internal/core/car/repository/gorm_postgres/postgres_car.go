package gorm_postgres

import (
	"cleanstandarts/internal/domain"
	"gorm.io/gorm"
)

type pgCarRepository struct {
	Conn *gorm.DB
}

func NewPostgresCarRepository(conn *gorm.DB) domain.CarRepository {
	return &pgCarRepository{conn}
}

func (cr *pgCarRepository) Fetch() (result []domain.Car, err error) {
	err = cr.Conn.Find(&result).Error
	return result, err
}

func (cr *pgCarRepository) GetByID(id uint64) (result domain.Car, err error) {
	err = cr.Conn.First(&result, "id = ?", id).Error
	return result, err
}

func (cr *pgCarRepository) Create(c *domain.Car) (err error) {
	err = cr.Conn.Create(c).Error
	return err
}

func (cr *pgCarRepository) Delete(id uint64) (err error) {
	err = cr.Conn.Delete(&domain.Car{}, "id = ?", id).Error
	return err
}
