package usecase

import (
	"cleanstandarts/internal/core/auction/repository/qstruct"
	"cleanstandarts/internal/core/auction/usecase/validation"
	"cleanstandarts/internal/domain"
	"cleanstandarts/internal/domain/dto"
	"cleanstandarts/internal/domain/dto/models"
	"cleanstandarts/pkg/validator"
	"cleanstandarts/pkg/yoopay"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"sync"
	"time"
)

type auctionUsecase struct {
	auctionRepo    domain.AuctionRepository
	bidRepo        domain.BidRepository
	carUsecase     domain.CarUsecase
	userUsecase    domain.UserUsecase
	telegramRepo   domain.TelegramRepository
	paymentService yoopay.Yoopay
	imagesRepo     domain.ImageRepository
	kafkaProducer  sarama.SyncProducer
}

func NewBaseCarUsecase(ar domain.AuctionRepository, br domain.BidRepository, cu domain.CarUsecase,
	uu domain.UserUsecase, tr domain.TelegramRepository, ps yoopay.Yoopay, ir domain.ImageRepository, kp *sarama.SyncProducer) domain.AuctionUsecase {
	return &auctionUsecase{
		auctionRepo:    ar,
		bidRepo:        br,
		carUsecase:     cu,
		userUsecase:    uu,
		telegramRepo:   tr,
		paymentService: ps,
		imagesRepo:     ir,
		kafkaProducer:  *kp,
	}
}

func (au *auctionUsecase) GetParticipantsByID(id uint64) (participants models.ParticipantOutputDTO, err error) {
	users, err := au.bidRepo.GetParticipantsByID(id)
	if err != nil {
		return
	}
	carName, err := au.auctionRepo.GetCarNameByID(id)
	if err != nil {
		return
	}
	seller, err := au.auctionRepo.GetSellerId(id)
	if err != nil {
		return
	}
	users = append(users, seller)
	chats, err := au.telegramRepo.GetChatsID(users)
	if err != nil {
		return
	}
	participants.ChatIDs = chats
	participants.AuctionName = carName
	return
}

func (au *auctionUsecase) FinishByID(id uint64) (err error) {
	auction, err := au.auctionRepo.GetByID(id)
	if err != nil {
		return err
	}
	bid, err := au.bidRepo.GetByID(auction.CurrentBidID)
	if err != nil {
		return err
	}

	if bid.Price < auction.Reserve {
		err = au.auctionRepo.FinishByID(id, 0)
	} else {

		err = au.auctionRepo.FinishByID(id, bid.UserID)
		pay_id, err := au.auctionRepo.GetWinnerPayIDByID(id)

		err = au.paymentService.CapturePayment(pay_id, "1g123")
		if err != nil {
			return err
		} //исправитб
	}
	if err != nil {
		return err
	}
	return
}

func (au *auctionUsecase) Fetch(q *qstruct.QueryParams) (result []models.AuctionOutputDTO, err error) {
	raw, err := au.auctionRepo.Fetch(q)
	if err != nil {
		return nil, err
	}
	for _, val := range raw {
		var (
			price  uint64
			car    models.CarOutputDTO
			seller models.UserOutputDTO
			winner models.UserOutputDTO
			images []domain.Image
		)
		wg := sync.WaitGroup{}
		wg.Add(4)
		go func() {
			defer wg.Done()
			price, err = au.bidRepo.GetMaxBidPrice(val.CurrentBidID)
		}()
		go func() {
			defer wg.Done()
			car, err = au.carUsecase.GetByID(val.CarID)
		}()
		go func() {
			defer wg.Done()
			seller, err = au.userUsecase.GetByID(val.SellerID)
		}()
		go func() {
			defer wg.Done()
			winner, err = au.userUsecase.GetByID(val.WinnerID)
		}()
		go func() {
			defer wg.Done()
			images, err = au.imagesRepo.GetByAucID(uint64(val.ID))
		}()
		wg.Wait()
		if err != nil {
			return nil, err
		}
		DTO := dto.AuctionToOutputDto(&val, &car, &seller, &winner, images, price)
		result = append(result, DTO)
	}
	return result, err
}

func (au *auctionUsecase) GetByID(id uint64) (result models.AuctionOutputDTO, err error) {
	auction, err := au.auctionRepo.GetByID(id)
	if err != nil {
		return result, err
	}
	var (
		price  uint64
		car    models.CarOutputDTO
		seller models.UserOutputDTO
		winner models.UserOutputDTO
		images []domain.Image
	)
	wg := sync.WaitGroup{}
	wg.Add(5)
	go func() {
		defer wg.Done()
		price, err = au.bidRepo.GetMaxBidPrice(id)
	}()
	go func() {
		defer wg.Done()
		car, err = au.carUsecase.GetByID(auction.CarID)
	}()
	go func() {
		defer wg.Done()
		seller, err = au.userUsecase.GetByID(auction.SellerID)
	}()
	go func() {
		defer wg.Done()
		winner, err = au.userUsecase.GetByID(auction.WinnerID)
	}()
	go func() {
		defer wg.Done()
		images, err = au.imagesRepo.GetByAucID(id)
	}()
	wg.Wait()
	if err != nil {
		return result, err
	}
	result = dto.AuctionToOutputDto(&auction, &car, &seller, &winner, images, price)
	return result, err
}

func (au *auctionUsecase) Delete(id uint64) (err error) {
	err = au.auctionRepo.Delete(id)
	return err
}

func (au *auctionUsecase) validateInputFormat(base *validator.BaseValidator, a *models.AuctionInputDTO) {
	v := validation.NewAuctionValidator(base)
	v.Validate(a)
}

func (au *auctionUsecase) validateAuction(a *models.AuctionInputDTO) (errors map[string]string, ok bool) {
	ok = true
	base := validator.NewBaseValidator()
	au.validateInputFormat(base, a)
	if !base.Valid() {
		errors = base.Errors
		ok = false
	}
	return errors, ok
}

func (au *auctionUsecase) Create(a *models.AuctionInputDTO, uID uint64) (errors map[string]string, ok bool) {
	errors, ok = au.validateAuction(a)
	if !ok {
		return errors, false
	}
	errors = map[string]string{}
	auction := dto.AuctionDtoToDomain(a, uID)
	err := au.auctionRepo.Create(&auction)
	if err != nil {
		errors["internal"] = err.Error()
		return errors, false
	}
	err = au.storeCron(uint64(auction.ID), auction.DateEnd)
	if err != nil {
		errors["internal"] = err.Error()
		return errors, false
	}
	a.Id = uint64(auction.ID)
	return errors, true
}

func (au *auctionUsecase) storeCron(aID uint64, timeEnd time.Time) (err error) {
	type dataStruct struct {
		AuctionId uint64    `json:"auction_id"`
		TimeEnd   time.Time `json:"time_end"`
	}
	data := dataStruct{
		AuctionId: aID,
		TimeEnd:   timeEnd,
	}
	dataJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error:", err.Error())
		return
	}
	msg := &sarama.ProducerMessage{
		Topic: "cron-new-auc",
		Value: sarama.ByteEncoder(dataJSON),
	}
	partition, offset, err := au.kafkaProducer.SendMessage(msg)
	if err != nil {
		fmt.Println("error:", err.Error())
		return
	}
	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", "cron-new-auc", partition, offset)
	return
}

//	HTTP storeCron
//func (au *auctionUsecase) storeCron(aID uint64, timeEnd time.Time) (err error) {
//	type bodyStruct struct {
//		AuctionId uint64    `json:"auction_id"`
//		TimeEnd   time.Time `json:"time_end"`
//	}
//
//	url := "_"
//
//	body := bodyStruct{
//		AuctionId: aID,
//		TimeEnd:   timeEnd,
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
