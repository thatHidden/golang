package http

import (
	"cleanstandarts/internal/core/user/delivery/http/middleware"
	"cleanstandarts/internal/domain"
	"cleanstandarts/internal/domain/dto/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

// ArticleHandler  represent the httphandler for article
type UserHandler struct {
	UserUsecase domain.UserUsecase
}

func NewUserHandler(e *gin.Engine, uu domain.UserUsecase, um *middleware.UserMiddleware) {
	handler := &UserHandler{
		UserUsecase: uu,
	}
	e.GET("/users", handler.Fetch)
	e.POST("/users", handler.Create)
	e.GET("/users/:id", handler.GetByID)
	e.DELETE("/users/:id", handler.Delete)
	e.POST("/users/activate", um.MustAuth, handler.Activate)
	e.POST("/login", handler.Login)
}

func (uh *UserHandler) Fetch(c *gin.Context) {
	var output []models.UserOutputDTO

	output, err := uh.UserUsecase.Fetch()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"got": output,
	})
}

func (uh *UserHandler) Activate(c *gin.Context) {
	var input struct {
		Token string
	}

	if c.Bind(&input) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
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

	err := uh.UserUsecase.Activate(user, input.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}

func (uh *UserHandler) Login(c *gin.Context) {
	var input struct {
		Email    string
		Password string
	}

	if c.Bind(&input) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	token, err := uh.UserUsecase.Auth(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", token, 3600, "", "", false, true)
	c.JSON(http.StatusNoContent, gin.H{})
}

func (uh *UserHandler) Create(c *gin.Context) {
	var input models.UserInputDTO

	if c.Bind(&input) != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	errors, ok := uh.UserUsecase.Create(&input)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": errors,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": input,
	})
}

func (uh *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	err = uh.UserUsecase.Delete(uint64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func (uh *UserHandler) GetByID(c *gin.Context) {
	var output struct {
		Email    string `json:"email"`
		Username string `json:"username"`
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	user, err := uh.UserUsecase.GetByID(uint64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	//output.Email = user.Email
	output.Username = user.Username

	c.JSON(http.StatusOK, gin.H{
		"got": output,
	})
}
