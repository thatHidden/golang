package usecase

import (
	"cleanstandarts/internal/core/payment/usecase/validation"
	"cleanstandarts/internal/domain"
	"cleanstandarts/pkg/validator"
	"cleanstandarts/pkg/yoopay"
	"github.com/google/uuid"
)

type paymentUsecase struct {
	paymentRepo    domain.PaymentRepository
	bidRepo        domain.BidRepository
	auctionRepo    domain.AuctionRepository
	paymentService yoopay.Yoopay
}

func NewPaymentUsecase(pr domain.PaymentRepository, br domain.BidRepository, ps yoopay.Yoopay,
	ar domain.AuctionRepository) domain.PaymentUsecase {
	return &paymentUsecase{
		paymentRepo:    pr,
		bidRepo:        br,
		paymentService: ps,
		auctionRepo:    ar,
	}
}

func (pu *paymentUsecase) validateInputFormat(base *validator.BaseValidator, p *domain.Payment) {
	v := validation.NewPaymentFormatValidator(base)
	v.Validate(p)
}

func (pu *paymentUsecase) validateBusinessRules(base *validator.BaseValidator, p *domain.Payment) {
	v := validation.NewPaymentRulesValidator(base, pu.auctionRepo, pu.bidRepo)
	v.Validate(p)
}

func (pu *paymentUsecase) validatePayment(p *domain.Payment) (errors map[string]string, ok bool) {
	base := validator.NewBaseValidator()
	pu.validateInputFormat(base, p)
	if !base.Valid() {
		return base.Errors, false
	}
	pu.validateBusinessRules(base, p)
	if !base.Valid() {
		return base.Errors, false
	}
	return errors, true
}

func (pu *paymentUsecase) Fetch() (result []domain.Payment, err error) {
	result, err = pu.paymentRepo.Fetch()
	return result, err
}

func (pu *paymentUsecase) GetByID(id uint64) (result domain.Payment, err error) {
	result, err = pu.paymentRepo.GetByID(id)
	return result, err
}

func (pu *paymentUsecase) GetByPaymentID(id string) (result domain.Payment, err error) {
	result, err = pu.paymentRepo.GetByPaymentID(id)
	return result, err
}

func (pu *paymentUsecase) Delete(id uint64) (err error) {
	err = pu.paymentRepo.Delete(id)
	return err
}

func (pu *paymentUsecase) Create(p *domain.Payment, u *domain.User) (confUrl string, errors map[string]string, ok bool) {
	p.UserID = uint64(u.ID)
	p.IsListed = false

	errors, ok = pu.validatePayment(p)
	if !ok {
		return "", errors, false
	}
	errors = map[string]string{}

	idempotenceKey := uuid.New().String()

	p.IdempotenceKey = idempotenceKey

	data, err := pu.paymentService.CreatePayment(p.Price, "bid", idempotenceKey)
	if err != nil {
		errors["error"] = err.Error()
		return "", errors, false
	}

	p.Status = data["status"]
	p.PayID = data["payment_id"]

	err = pu.paymentRepo.Create(p)

	confirmationUrl := data["confirmation_url"]

	return confirmationUrl, nil, true
}
