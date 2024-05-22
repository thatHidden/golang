package http

import (
	"cleanstandarts/internal/core/user/delivery/http/middleware"
	"cleanstandarts/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResponseError struct {
	Message string `json:"message"`
}

// ArticleHandler  represent the httphandler for article
type PaymentHandler struct {
	PaymentUsecase domain.PaymentUsecase
}

func NewPaymentHandler(e *gin.Engine, pu domain.PaymentUsecase, um *middleware.UserMiddleware) {
	handler := &PaymentHandler{
		PaymentUsecase: pu,
	}
	e.POST("/payments", um.MustAuth, handler.Create)
	//e.GET("/payments", handler.Fetch)
	//e.GET("/payments/:id", handler.GetByID)
	//e.DELETE("/payments/:id", handler.Delete)
}

func (ph *PaymentHandler) Create(c *gin.Context) {
	var input struct {
		AuctionID uint64 `json:"auction_id"`
		Price     uint64 `json:"price"`
	}

	if c.Bind(&input) != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	payment := domain.Payment{
		AuctionID: input.AuctionID,
		Price:     input.Price,
	}

	_user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	user, ok := _user.(domain.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"lox": "lox",
		})
		return
	}

	url, err, ok := ph.PaymentUsecase.Create(&payment, &user)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"created":          input,
		"confirmation_url": url,
	})
}

//func (ph *PaymentHandler) Fetch(c *gin.Context) {
//
//}
//
//func (ph *PaymentHandler) GetByID(c *gin.Context) {
//
//}
//
//func (ph *PaymentHandler) Delete(c *gin.Context) {
//
//}
