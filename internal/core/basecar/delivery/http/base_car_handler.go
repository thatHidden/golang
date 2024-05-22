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
type BaseCarHandler struct {
	BaseCarUsecase domain.BaseCarUsecase
}

func NewBaseCarHandler(e *gin.Engine, bcu domain.BaseCarUsecase) {
	handler := &BaseCarHandler{
		BaseCarUsecase: bcu,
	}
	e.POST("/base_cars", handler.Create)
	e.GET("/base_cars", handler.Fetch)
	e.GET("/base_cars/:id", handler.GetByID)
	e.DELETE("/base_cars/:id", handler.Delete)
}

func (bch *BaseCarHandler) Fetch(c *gin.Context) {
	var output []models.CarBaseOutputDTO

	output, err := bch.BaseCarUsecase.Fetch()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"got": output,
	})
}

func (bch *BaseCarHandler) Create(c *gin.Context) {
	var input models.CarBaseInputDTO

	if c.Bind(&input) != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	errors, ok := bch.BaseCarUsecase.Create(&input)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"created": input,
	})
}

func (bch *BaseCarHandler) GetByID(c *gin.Context) {
	var output models.CarBaseOutputDTO

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	output, err = bch.BaseCarUsecase.GetByID(uint64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"got": output,
	})
}

func (bch *BaseCarHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	err = bch.BaseCarUsecase.Delete(uint64(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
	}

	c.JSON(http.StatusOK, gin.H{})
}
