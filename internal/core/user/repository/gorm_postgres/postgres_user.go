package gorm_postgres

import (
	"cleanstandarts/internal/domain"
	"gorm.io/gorm"
)

type pgUserRepository struct {
	Conn *gorm.DB
}

func NewPostgresUserRepository(conn *gorm.DB) domain.UserRepository {
	return &pgUserRepository{conn}
}

func (p *pgUserRepository) Fetch() (result []domain.User, err error) {
	err = p.Conn.Find(&result).Error
	return result, err
}

func (p *pgUserRepository) Activate(id uint64) (err error) {
	err = p.Conn.
		Model(&domain.User{}).
		Where("id = ?", id).
		Update("activated", true).Error
	return err
}

func (p *pgUserRepository) GetByEmail(email string) (result domain.User, err error) {
	err = p.Conn.First(&result, "email = ?", email).Error
	if err != nil {
		return domain.User{}, err
	}
	return result, nil
}

func (p *pgUserRepository) GetByID(id uint64) (result domain.User, err error) {
	err = p.Conn.First(&result, "id = ?", id).Error
	if err != nil {
		return domain.User{}, err
	}
	return result, nil
}

func (p *pgUserRepository) Create(u *domain.User) (err error) {
	err = p.Conn.Create(u).Error
	return err
}

func (p *pgUserRepository) Delete(id uint64) (err error) {
	err = p.Conn.Delete(&domain.User{}, "id = ?", id).Error
	return err
}

//func (p *pgUserRepository) Update(u *domain.User) (err error) {
//	err = p.Conn.Save(u).Error
//	return nil
//}
