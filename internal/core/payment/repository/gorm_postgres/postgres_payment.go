package gorm_postgres

import (
	"cleanstandarts/internal/domain"
	"gorm.io/gorm"
)

type pgPaymentRepository struct {
	Conn *gorm.DB
}

func NewPostgresPaymentRepository(conn *gorm.DB) domain.PaymentRepository {
	return &pgPaymentRepository{conn}
}

func (pr *pgPaymentRepository) IsListedByPaymentID(id string) (result bool, err error) {
	err = pr.Conn.Model(&domain.Payment{}).Select("is_listed").Where("pay_id = ?", id).Scan(&result).Error
	return result, err
}

func (pr *pgPaymentRepository) SetListedByPaymentID(id string) (err error) {
	err = pr.Conn.Model(&domain.Payment{}).Where("pay_id = ?", id).Update("is_listed", true).Error
	return err
}

func (pr *pgPaymentRepository) Fetch() (result []domain.Payment, err error) {
	query := pr.Conn.Model(&domain.Payment{})

	err = query.Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (pr *pgPaymentRepository) IsExistsByPaymentID(id string) (result bool, err error) {
	err = pr.Conn.Model(&domain.Payment{}).Select("count(*) > 0").Where("pay_id = ?", id).Find(&result).Error
	if err != nil {
		return false, err
	}
	return result, nil
}

func (pr *pgPaymentRepository) GetByPaymentID(id string) (result domain.Payment, err error) {
	err = pr.Conn.First(&result, "pay_id = ?", id).Error
	if err != nil {
		return domain.Payment{}, err
	}
	return result, nil
}

func (pr *pgPaymentRepository) GetByID(id uint64) (result domain.Payment, err error) {
	err = pr.Conn.First(&result, "id = ?", id).Error
	if err != nil {
		return domain.Payment{}, err
	}
	return result, nil
}

func (pr *pgPaymentRepository) Create(c *domain.Payment) (err error) {
	err = pr.Conn.Create(c).Error
	return err
}

func (pr *pgPaymentRepository) Delete(id uint64) (err error) {
	err = pr.Conn.Delete(&domain.Payment{}, "id = ?", id).Error
	return err
}
