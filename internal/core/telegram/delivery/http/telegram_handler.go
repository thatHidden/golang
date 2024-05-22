package http

import (
	"cleanstandarts/internal/core/user/delivery/http/middleware"
	"cleanstandarts/internal/domain"
	"cleanstandarts/internal/domain/dto/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ResponseError struct {
	Message string `json:"message"`
}

// ArticleHandler  represent the httphandler for article
type TelegramHandler struct {
	TelegramUsecase domain.TelegramUsecase
}

func NewTelegramHandler(e *gin.Engine, tu domain.TelegramUsecase, um *middleware.UserMiddleware) {
	handler := &TelegramHandler{
		TelegramUsecase: tu,
	}
	e.POST("/telegrams/set_chat", handler.SetChatID)
	e.PUT("/telegrams/:id", um.MustAuth, handler.Update)
	e.POST("/telegrams", um.MustAuth, handler.Create)
	e.GET("/telegrams", handler.Fetch)
	e.GET("/telegrams/:id", handler.GetByID)
	e.DELETE("/telegrams/:id", handler.Delete)
}

// вызывается только ботом
func (th *TelegramHandler) SetChatID(c *gin.Context) {
	var input struct {
		UserName string `json:"username"`
		ChatId   uint64 `json:"chat_id"`
	}

	if c.Bind(&input) != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	err := th.TelegramUsecase.SetChatID(input.UserName, input.ChatId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (th *TelegramHandler) Update(c *gin.Context) {
	var input struct {
		Alerts bool `json:"alerts"`
	}

	if c.Bind(&input) != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	_user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	user, ok := _user.(domain.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	fmt.Println(input.Alerts)
	err := th.TelegramUsecase.Update(uint64(user.ID), input.Alerts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (th *TelegramHandler) Fetch(c *gin.Context) {
	var output []models.TelegramOutputDTO

	output, err := th.TelegramUsecase.Fetch()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"got": output,
	})
}

func (th *TelegramHandler) GetByID(c *gin.Context) {
	var output models.TelegramOutputDTO

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	output, err = th.TelegramUsecase.GetByID(uint64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"got": output,
	})
}

func (th *TelegramHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	err = th.TelegramUsecase.Delete(uint64(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (th *TelegramHandler) Create(c *gin.Context) {
	var input models.TelegramInputDTO

	if c.Bind(&input) != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
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

	input.UserID = uint64(user.ID)

	errors, ok := th.TelegramUsecase.Create(&input)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}
