package gorm_postgres

import (
	"cleanstandarts/internal/domain"
	"gorm.io/gorm"
)

type pgTokenRepository struct {
	Conn *gorm.DB
}

func NewPostgresTokenRepository(conn *gorm.DB) domain.TokenRepository {
	return &pgTokenRepository{conn}
}

func (tr *pgTokenRepository) GetByPlaintext(plaintext string) (result domain.Token, err error) {
	err = tr.Conn.First(&result, "plaintext = ?", plaintext).Error
	if err != nil {
		return domain.Token{}, err
	}
	return result, nil
}

func (tr *pgTokenRepository) GetByUserID(id uint64) (result domain.Token, err error) {
	err = tr.Conn.First(&result, "user_id = ?", id).Error
	if err != nil {
		return domain.Token{}, err
	}
	return result, nil
}

func (tr *pgTokenRepository) Create(t *domain.Token) (err error) {
	err = tr.Conn.Create(t).Error
	//+отправка на почту (в usecase)
	return err
}

func (tr *pgTokenRepository) Delete(id uint64) (err error) {
	err = tr.Conn.Delete(&domain.Token{}, "id = ?", id).Error
	return err
}
