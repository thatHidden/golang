package gorm_postgres

import (
	"cleanstandarts/internal/domain"
	"gorm.io/gorm"
)

type pgCommentRepository struct {
	Conn *gorm.DB
}

func NewPostgresCommentRepository(conn *gorm.DB) domain.CommentRepository {
	return &pgCommentRepository{conn}
}

func (cr *pgCommentRepository) Fetch(userID uint64, auctionID uint64) (result []domain.Comment, err error) {
	query := cr.Conn.Model(&domain.Comment{})

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

func (cr *pgCommentRepository) GetByID(id uint64) (result domain.Comment, err error) {
	err = cr.Conn.First(&result, "id = ?", id).Error
	if err != nil {
		return domain.Comment{}, err
	}
	return result, nil
}

func (cr *pgCommentRepository) Create(c *domain.Comment) (err error) {
	err = cr.Conn.Create(c).Error
	return err
}

func (cr *pgCommentRepository) Delete(id uint64) (err error) {
	err = cr.Conn.Delete(&domain.Comment{}, "id = ?", id).Error
	return err
}
