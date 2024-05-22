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
type CarModsHandler struct {
	CarModsUsecase domain.CarModsUsecase
}

func NewCarModsHandler(e *gin.Engine, cmu domain.CarModsUsecase) {
	handler := &CarModsHandler{
		CarModsUsecase: cmu,
	}
	e.GET("/car_mods/:id", handler.GetByID)
	e.GET("/car_mods", handler.Fetch)
	e.POST("/car_mods", handler.Create)
	e.DELETE("/car_mods/:id", handler.Delete)
}

func (cmh *CarModsHandler) Fetch(c *gin.Context) {
	var output []models.CarModsOutputDTO

	output, err := cmh.CarModsUsecase.Fetch()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"got": output,
	})
}

func (cmh *CarModsHandler) GetByID(c *gin.Context) {
	var output models.CarModsOutputDTO

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	output, err = cmh.CarModsUsecase.GetByID(uint64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"got": output,
	})
}

func (cmh *CarModsHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	err = cmh.CarModsUsecase.Delete(uint64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (cmh *CarModsHandler) Create(c *gin.Context) {
	var input models.CarModsInputDTO

	if c.Bind(&input) != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	errors, ok := cmh.CarModsUsecase.Create(&input)
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
