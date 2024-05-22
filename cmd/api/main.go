package main

import (
	_auctionHTTP "cleanstandarts/internal/core/auction/delivery/http"
	_auctionRepository "cleanstandarts/internal/core/auction/repository/gorm_postgres"
	_auctionUsecase "cleanstandarts/internal/core/auction/usecase"
	_baseCarHTTP "cleanstandarts/internal/core/basecar/delivery/http"
	_baseCarRepository "cleanstandarts/internal/core/basecar/repository/gorm_postgres"
	_baseCarUsecase "cleanstandarts/internal/core/basecar/usecase"
	_bidHTTP "cleanstandarts/internal/core/bid/delivery/http"
	_bidRepository "cleanstandarts/internal/core/bid/repository/gorm_postgres"
	_bidUsecase "cleanstandarts/internal/core/bid/usecase"
	_carHTTP "cleanstandarts/internal/core/car/delivery/http"
	_carRepository "cleanstandarts/internal/core/car/repository/gorm_postgres"
	_carUsecase "cleanstandarts/internal/core/car/usecase"
	_carModsHTTP "cleanstandarts/internal/core/car_mods/delivery/http"
	_carModsRepository "cleanstandarts/internal/core/car_mods/repository/gorm_postgres"
	_carModsUsecase "cleanstandarts/internal/core/car_mods/usecase"
	_commentHTTP "cleanstandarts/internal/core/comment/delivery/http"
	_commentRepository "cleanstandarts/internal/core/comment/repository/gorm_postgres"
	_commentUsecase "cleanstandarts/internal/core/comment/usecase"
	_imageRepository "cleanstandarts/internal/core/image/repository/gorm_postgres"
	_imageUsecase "cleanstandarts/internal/core/image/usecase"
	_paymentHTTP "cleanstandarts/internal/core/payment/delivery/http"
	_paymentRepository "cleanstandarts/internal/core/payment/repository/gorm_postgres"
	_paymentUsecase "cleanstandarts/internal/core/payment/usecase"
	_telegramHTTP "cleanstandarts/internal/core/telegram/delivery/http"
	_telegramRepository "cleanstandarts/internal/core/telegram/repository/gorm_postgres"
	_telegramUsecase "cleanstandarts/internal/core/telegram/usecase"
	_tokenRepository "cleanstandarts/internal/core/token/repository/gorm_postgres"
	_userHTTP "cleanstandarts/internal/core/user/delivery/http"
	_userHandlerMiddleware "cleanstandarts/internal/core/user/delivery/http/middleware"
	_userRepository "cleanstandarts/internal/core/user/repository/gorm_postgres"
	_userUsecase "cleanstandarts/internal/core/user/usecase"
	"fmt"
	"github.com/IBM/sarama"

	"cleanstandarts/internal/domain"
	"cleanstandarts/pkg/yoopay"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println(os.Getenv("DB_DSN"))
	fmt.Println(os.Getenv("YOOPAY_SECRET"))
}

func main() {
	db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{})

	if err != nil {
		log.Fatal("Error connecting to database")
	}

	err = db.AutoMigrate(&domain.User{}, &domain.BaseCar{}, &domain.Car{}, &domain.Bid{}, &domain.Auction{},
		&domain.CarMods{}, &domain.Comment{}, &domain.Payment{}, &domain.Telegram{}, &domain.Image{})
	if err != nil {
		log.Fatal("Error auto migrate database")
	}

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	brokers := []string{"172.30.73.87:9092"}

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatal("Failed to start Sarama producer:", err)
	}
	defer producer.Close()

	r := gin.New()

	r.Static("/static", "./static")

	userRepository := _userRepository.NewPostgresUserRepository(db)
	tokenRepository := _tokenRepository.NewPostgresTokenRepository(db)
	userUsecase := _userUsecase.NewUserUsecase(userRepository, tokenRepository)
	userMiddleware := _userHandlerMiddleware.NewUserMiddleware(userUsecase)

	baseCarRepository := _baseCarRepository.NewPostgresCarRepository(db)
	baseCarUsecase := _baseCarUsecase.NewBaseCarUsecase(baseCarRepository)

	carModsRepository := _carModsRepository.NewPostgresCarModsRepository(db)
	carModsUsecase := _carModsUsecase.NewCarModsUsecase(carModsRepository)

	carRepository := _carRepository.NewPostgresCarRepository(db)
	carUsecase := _carUsecase.NewCarUsecase(carRepository, carModsUsecase, baseCarUsecase)

	commentRepository := _commentRepository.NewPostgresCommentRepository(db)
	commentUsecase := _commentUsecase.NewCommentUsecase(commentRepository)

	bidRepository := _bidRepository.NewPostgresBidRepository(db)

	auctionRepository := _auctionRepository.NewPostgresAuctionRepository(db)

	paymentRepository := _paymentRepository.NewPostgresPaymentRepository(db)
	paymentService := yoopay.NewYoopayService(os.Getenv("YOOPAY_ID"), os.Getenv("YOOPAY_SECRET"))
	paymentUsecase := _paymentUsecase.NewPaymentUsecase(paymentRepository, bidRepository, paymentService, auctionRepository)

	telegramRepository := _telegramRepository.NewPostgresTelegramRepository(db)
	telegramUsecase := _telegramUsecase.NewTelegramUsecase(telegramRepository)

	bidUsecase := _bidUsecase.NewBidUsecase(bidRepository, paymentRepository, auctionRepository, telegramRepository,
		&producer)

	imageRepository := _imageRepository.NewPostgresImageRepository(db)
	imageUsecase := _imageUsecase.NewCommentUsecase(imageRepository)

	auctionUsecase := _auctionUsecase.NewBaseCarUsecase(auctionRepository, bidRepository, carUsecase, userUsecase,
		telegramRepository, paymentService, imageRepository, &producer)

	_userHTTP.NewUserHandler(r, userUsecase, userMiddleware)
	_baseCarHTTP.NewBaseCarHandler(r, baseCarUsecase)
	_carHTTP.NewCarHandler(r, carUsecase)
	_carModsHTTP.NewCarModsHandler(r, carModsUsecase)
	_auctionHTTP.NewAuctionHandler(r, auctionUsecase, imageUsecase, userMiddleware)
	_commentHTTP.NewCommentHandler(r, commentUsecase, userMiddleware)
	_bidHTTP.NewBidHandler(r, bidUsecase, paymentUsecase, userMiddleware)
	_paymentHTTP.NewPaymentHandler(r, paymentUsecase, userMiddleware)
	_telegramHTTP.NewTelegramHandler(r, telegramUsecase, userMiddleware)

	log.Fatal(r.Run(":8000"))
}
