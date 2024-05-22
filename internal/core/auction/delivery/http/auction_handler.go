package http

import (
	"cleanstandarts/internal/core/auction/repository/qstruct"
	"cleanstandarts/internal/core/user/delivery/http/middleware"
	response "cleanstandarts/internal/core/utils"
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
type AuctionHandler struct {
	AuctionUsecase domain.AuctionUsecase
	ImageUsecase   domain.ImageUsecase
}

func NewAuctionHandler(e *gin.Engine, au domain.AuctionUsecase, iu domain.ImageUsecase, um *middleware.UserMiddleware) {
	handler := &AuctionHandler{
		AuctionUsecase: au,
		ImageUsecase:   iu,
	}
	e.GET("/auctions", handler.Fetch)
	e.GET("/auctions/:id/participants", handler.GetParticipantsByID)
	e.POST("/auctions", um.MustAuth, handler.Create)
	e.GET("/auctions/:id", handler.GetByID)
	e.POST("/auctions/:id/finish", handler.FinishByID)
	e.DELETE("/auctions/:id", um.MustAuth, handler.Delete)
	e.POST("/auctions/:id/images", handler.CreateImages)
}

func (ah *AuctionHandler) CreateImages(c *gin.Context) {
	aucId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	form, _ := c.MultipartForm()
	interiorImages := form.File["interior"]
	exteriorImages := form.File["exterior"]

	dto := models.ImagesInputDTO{
		Exterior: exteriorImages,
		Interior: interiorImages,
		AucID:    aucId,
	}

	errors, ok := ah.ImageUsecase.Create(&dto)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": errors,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func (ah *AuctionHandler) GetParticipantsByID(c *gin.Context) {
	var output models.ParticipantOutputDTO

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	output, err = ah.AuctionUsecase.GetParticipantsByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"got": output,
	})
}

func (ah *AuctionHandler) FinishByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	fmt.Println("finish:", id)

	err = ah.AuctionUsecase.FinishByID(uint64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (ah *AuctionHandler) Fetch(c *gin.Context) {
	var output []models.AuctionOutputDTO

	var query qstruct.QueryParams

	var err error

	if c.Query("user_id") != "" {
		query.UserID, err = strconv.ParseUint(c.Query("user_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}
	}

	if c.Query("no_reserve") == "true" {
		query.IsNoReserve = true
	}

	if c.Query("ended") == "true" {
		query.IsEnded = true
	}

	if c.Query("ending_soon") == "true" {
		query.IsEndingSoon = true
	}

	if c.Query("min_year") != "" {
		minYear, err := strconv.ParseUint(c.Query("min_year"), 10, 16)
		query.MinYear = uint16(minYear)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}
	}

	if c.Query("max_year") != "" {
		maxYear, err := strconv.ParseUint(c.Query("max_year"), 10, 16)
		query.MaxYear = uint16(maxYear)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}
	}

	if c.Query("brand") != "" {
		query.Brand = c.Query("brand")
	}

	if c.Query("model") != "" {
		query.Model = c.Query("model")
	}

	if c.Query("gen") != "" {
		query.Generation = c.Query("gen")
	}

	if c.Query("body_style") != "" {
		query.BodyStyle = c.Query("body_style")
	}

	output, err = ah.AuctionUsecase.Fetch(&query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"got": output,
	})
}

func (ah *AuctionHandler) GetByID(c *gin.Context) {
	var output models.AuctionOutputDTO

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	output, err = ah.AuctionUsecase.GetByID(uint64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	response.ResponseOk(c, output, c.Request.Method, c.FullPath(), http.StatusOK)
	//c.JSON(http.StatusOK, gin.H{
	//	"got": output,
	//})
}

func (ah *AuctionHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	err = ah.AuctionUsecase.Delete(uint64(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (ah *AuctionHandler) Create(c *gin.Context) {
	var input models.AuctionInputDTO

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

	errors, ok := ah.AuctionUsecase.Create(&input, uint64(user.ID))
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": errors,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"created": input,
	})
}
