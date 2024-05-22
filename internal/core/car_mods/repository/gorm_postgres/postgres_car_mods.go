package gorm_postgres

import (
	"cleanstandarts/internal/domain"
	"gorm.io/gorm"
)

type pgCarModsRepository struct {
	Conn *gorm.DB
}

func NewPostgresCarModsRepository(conn *gorm.DB) domain.CarModsRepository {
	return &pgCarModsRepository{conn}
}

func (cmr *pgCarModsRepository) Fetch() (result []domain.CarMods, err error) {
	err = cmr.Conn.Find(&result).Error
	return result, err
}

func (cmr *pgCarModsRepository) GetByID(id uint64) (result *domain.CarMods, err error) {
	err = cmr.Conn.First(&result, "id = ?", id).Error
	if err != nil {
		return &domain.CarMods{}, err
	}
	return result, nil
}

func (cmr *pgCarModsRepository) Create(c *domain.CarMods) (err error) {
	err = cmr.Conn.Create(c).Error
	return err
}

func (cmr *pgCarModsRepository) Delete(id uint64) (err error) {
	err = cmr.Conn.Delete(&domain.CarMods{}, "id = ?", id).Error
	return err
}
