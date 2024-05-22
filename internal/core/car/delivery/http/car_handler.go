package http

import (
	"cleanstandarts/internal/domain"
	"cleanstandarts/internal/domain/dto/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ResponseError struct {
	Message string `json:"message"`
}

// ArticleHandler  represent the httphandler for article
type CarHandler struct {
	CarUsecase domain.CarUsecase
}

func NewCarHandler(e *gin.Engine, cu domain.CarUsecase) {
	handler := &CarHandler{
		CarUsecase: cu,
	}
	e.POST("/cars", handler.Create)
	e.GET("/cars", handler.Fetch)
	e.GET("/cars/:id", handler.GetByID)
	e.DELETE("/cars/:id", handler.Delete)
}

func (ch *CarHandler) Fetch(c *gin.Context) {
	var output []models.CarOutputDTO

	output, err := ch.CarUsecase.Fetch()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"got": output,
	})
}

func (ch *CarHandler) GetByID(c *gin.Context) {
	var output models.CarOutputDTO

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	output, err = ch.CarUsecase.GetByID(uint64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"got": output,
	})
}

func (ch *CarHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	err = ch.CarUsecase.Delete(uint64(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (ch *CarHandler) Create(c *gin.Context) {
	var input models.CarInputDTO

	err := c.Bind(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	errors, ok := ch.CarUsecase.Create(&input)
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
