package http

import (
	"cleanstandarts/internal/core/user/delivery/http/middleware"
	"cleanstandarts/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ResponseError struct {
	Message string `json:"message"`
}

type BidHandler struct {
	BidUsecase     domain.BidUsecase
	PaymentUsecase domain.PaymentUsecase
}

func NewBidHandler(e *gin.Engine, bu domain.BidUsecase, pu domain.PaymentUsecase, um *middleware.UserMiddleware) {
	handler := &BidHandler{
		BidUsecase:     bu,
		PaymentUsecase: pu,
	}
	e.POST("/bids", handler.Create)
	e.GET("/bids", handler.Fetch)
	e.GET("/bids/:id", handler.GetByID)
	e.DELETE("/bids/:id", um.MustAuth, handler.Delete)
}

func (bh *BidHandler) Fetch(c *gin.Context) {
	var output []domain.Bid

	var userID uint64
	var auctionID uint64

	var err error

	if c.Query("user_id") != "" {
		userID, err = strconv.ParseUint(c.Query("user_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}
	}

	if c.Query("auction_id") != "" {
		auctionID, err = strconv.ParseUint(c.Query("auction_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}
	}

	output, err = bh.BidUsecase.Fetch(userID, auctionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"got": output,
	})
}

func (bh *BidHandler) GetByID(c *gin.Context) {
	var output domain.Bid

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	output, err = bh.BidUsecase.GetByID(uint64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"got": output,
	})
}

func (bh *BidHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	err = bh.BidUsecase.Delete(uint64(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
	}

	c.JSON(http.StatusOK, gin.H{})
}

// Create обработчик может быть вызвана исключительно сервером уведомлений
func (bh *BidHandler) Create(c *gin.Context) {
	var input struct {
		PayID string `json:"payment_id"`
	}

	//десереализация
	if c.Bind(&input) != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	errors, ok := bh.BidUsecase.Create(input.PayID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}
