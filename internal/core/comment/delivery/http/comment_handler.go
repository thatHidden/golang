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

// ArticleHandler  represent the httphandler for article
type CommentHandler struct {
	CommentUsecase domain.CommentUsecase
}

func NewCommentHandler(e *gin.Engine, cu domain.CommentUsecase, um *middleware.UserMiddleware) {
	handler := &CommentHandler{
		CommentUsecase: cu,
	}
	e.POST("/comments", um.MustAuth, handler.Create)
	e.GET("/comments", handler.Fetch)
	e.GET("/comments/:id", handler.GetByID)
	e.DELETE("/comments/:id", um.MustAuth, handler.Delete)
}

func (ch *CommentHandler) Fetch(c *gin.Context) {
	var output []domain.Comment

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

	output, err = ch.CommentUsecase.Fetch(userID, auctionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"got": output,
	})
}

func (ch *CommentHandler) GetByID(c *gin.Context) {
	var output domain.Comment

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	output, err = ch.CommentUsecase.GetByID(uint64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"got": output,
	})
}

func (ch *CommentHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	err = ch.CommentUsecase.Delete(uint64(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (ch *CommentHandler) Create(c *gin.Context) {
	var input struct {
		AuctionID uint64 `json:"auction_id"`
		Text      string `json:"text"`
		ReplyID   uint64 `json:"reply_id"`
	}

	if c.Bind(&input) != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	comment := domain.Comment{
		AuctionID: input.AuctionID,
		Text:      input.Text,
		ReplyID:   input.ReplyID,
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

	errors, ok := ch.CommentUsecase.Create(&comment, &user)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"created": input,
	})
}
