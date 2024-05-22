package usecase

import (
	"cleanstandarts/internal/core/bid/usecase/validation"
	"cleanstandarts/internal/domain"
	"cleanstandarts/pkg/validator"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
)

type bidUsecase struct {
	bidRepo       domain.BidRepository
	paymentRepo   domain.PaymentRepository
	auctionRepo   domain.AuctionRepository
	telegramRepo  domain.TelegramRepository
	kafkaProducer sarama.SyncProducer
}

func NewBidUsecase(br domain.BidRepository, pr domain.PaymentRepository, ar domain.AuctionRepository,
	tr domain.TelegramRepository, kp *sarama.SyncProducer) domain.BidUsecase {
	return &bidUsecase{
		bidRepo:       br,
		paymentRepo:   pr,
		auctionRepo:   ar,
		telegramRepo:  tr,
		kafkaProducer: *kp,
	}
}

func (cu *bidUsecase) validateInputFormat(base *validator.BaseValidator, id string) {
	v := validation.NewBidFormatValidator(base)
	v.Validate(&id)
}

func (cu *bidUsecase) validateBusinessRules(base *validator.BaseValidator, id string) {
	v := validation.NewBidRulesValidator(base)
	v.Validate(&id)
}

func (cu *bidUsecase) validateBid(id string) (errors map[string]string, ok bool) {
	base := validator.NewBaseValidator()
	cu.validateInputFormat(base, id)
	if !base.Valid() {
		return base.Errors, false
	}
	cu.validateBusinessRules(base, id)
	if !base.Valid() {
		return base.Errors, false
	}
	return errors, true
}

func (cu *bidUsecase) Fetch(userID uint64, auctionID uint64) (result []domain.Bid, err error) {
	result, err = cu.bidRepo.Fetch(userID, auctionID)
	return result, err
}

func (cu *bidUsecase) GetByID(id uint64) (result domain.Bid, err error) {
	result, err = cu.bidRepo.GetByID(id)
	return result, err
}

func (cu *bidUsecase) Delete(id uint64) (err error) {
	err = cu.bidRepo.Delete(id)
	return err
}

func (cu *bidUsecase) Create(payID string) (errors map[string]string, ok bool) {
	errors, ok = cu.validateBid(payID)
	if !ok {
		return errors, false
	}
	errors = map[string]string{}
	payment, err := cu.paymentRepo.GetByPaymentID(payID)
	if err != nil {
		errors["error"] = err.Error()
		return errors, false
	}

	bid := domain.Bid{
		AuctionID: payment.AuctionID,
		UserID:    payment.UserID,
		Price:     payment.Price,
		PaymentID: uint64(payment.ID),
	}

	oldLeadId, err := cu.auctionRepo.GetLeadBidID(bid.AuctionID)
	if err != nil {
		errors["error"] = err.Error()
		return errors, false
	}

	err = cu.bidRepo.Create(&bid)
	if err != nil {
		errors["error"] = err.Error()
		return errors, false
	}

	err = cu.paymentRepo.SetListedByPaymentID(payID)
	if err != nil {
		errors["error"] = err.Error()
		return errors, false
	}

	err = cu.auctionRepo.SetLeadBidID(bid.AuctionID, uint64(bid.ID))
	if err != nil {
		errors["error"] = err.Error()
		return errors, false
	}

	//?go
	if oldLeadId != 0 {
		go cu.overbidAlert(oldLeadId, bid.AuctionID)
	}

	return errors, true
}

func (cu *bidUsecase) overbidAlert(oldLeadId, aucId uint64) {
	type dataStruct struct {
		AuctionName string `json:"auction_name"`
		Chat        uint64 `json:"chat"`
	}

	oldLead, err := cu.bidRepo.GetByID(oldLeadId)
	aucName, err := cu.auctionRepo.GetCarNameByID(aucId)
	tg, err := cu.telegramRepo.GetByUserID(oldLead.UserID)

	data := dataStruct{
		AuctionName: aucName,
		Chat:        tg.ChatID,
	}
	dataJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error:", err.Error())
		return
	}
	msg := &sarama.ProducerMessage{
		Topic: "alert-overbid-auc",
		Value: sarama.ByteEncoder(dataJSON),
	}
	partition, offset, err := cu.kafkaProducer.SendMessage(msg)
	if err != nil {
		fmt.Println("error:", err.Error())
		return
	}
	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", "alert-overbid-auc", partition, offset)
}

// HTTP overbidAlert
//func (cu *bidUsecase) overbidAlert(oldLeadId, aucId uint64) (errors map[string]string, ok bool) {
//	oldLead, err := cu.bidRepo.GetByID(oldLeadId)
//	aucName, err := cu.auctionRepo.GetCarNameByID(aucId)
//	tg, err := cu.telegramRepo.GetByUserID(oldLead.UserID)
//
//	type bodyStruct struct {
//		AuctionName string `json:"auction_name"`
//		Chat        uint64 `json:"chat"`
//	}
//
//	url := "_"
//
//	body := bodyStruct{
//		AuctionName: aucName,
//		Chat:        tg.ChatID,
//	}
//
//	jsonData, err := json.Marshal(body)
//	if err != nil {
//		return
//	}
//
//	r, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
//	if err != nil {
//		return
//	}
//
//	r.Header.Add("Content-Type", "application/json")
//
//	client := &http.Client{}
//	res, err := client.Do(r)
//	if err != nil {
//		return
//	}
//
//	defer res.Body.Close()
//	return
//}
