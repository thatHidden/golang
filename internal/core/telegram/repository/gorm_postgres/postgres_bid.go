package gorm_postgres

import (
	"cleanstandarts/internal/domain"
	"gorm.io/gorm"
)

type pgTelegramRepository struct {
	Conn *gorm.DB
}

func NewPostgresTelegramRepository(conn *gorm.DB) domain.TelegramRepository {
	return &pgTelegramRepository{conn}
}

func (br *pgTelegramRepository) GetChatsID(uID []uint64) (result []uint64, err error) {
	err = br.Conn.Model(domain.Telegram{}).
		Select("chat_Id").
		Where("user_id IN ?", uID).
		Scan(&result).Error
	return
}

func (br *pgTelegramRepository) GetByUserID(id uint64) (result domain.Telegram, err error) {
	err = br.Conn.First(&result, "user_id = ?", id).Error
	return
}

func (br *pgTelegramRepository) GetByUsername(username string) (result domain.Telegram, err error) {
	err = br.Conn.First(&result, "username = ?", username).Error
	return
}

func (br *pgTelegramRepository) Update(uID uint64, val bool) (err error) {
	err = br.Conn.Model(&domain.Telegram{}).
		Where("user_id = ?", uID).
		Update("do_alerts", val).
		Error
	return
}

func (br *pgTelegramRepository) SetChatID(username string, cID uint64) (err error) {
	err = br.Conn.Model(&domain.Telegram{}).
		Where("username = ?", username).
		Update("chat_id", cID).
		Error
	return
}

func (br *pgTelegramRepository) Fetch() (result []domain.Telegram, err error) {
	query := br.Conn.Model(&domain.Telegram{})

	err = query.Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (br *pgTelegramRepository) GetByID(id uint64) (result domain.Telegram, err error) {
	err = br.Conn.First(&result, "id = ?", id).Error
	if err != nil {
		return domain.Telegram{}, err
	}
	return result, nil
}

func (br *pgTelegramRepository) Create(c *domain.Telegram) (err error) {
	err = br.Conn.Create(c).Error
	return err
}

func (br *pgTelegramRepository) Delete(id uint64) (err error) {
	err = br.Conn.Delete(&domain.Telegram{}, "id = ?", id).Error
	return err
}
