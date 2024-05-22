package middleware

import (
	"cleanstandarts/internal/domain"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"time"
)

type UserMiddleware struct {
	UserUsecase domain.UserUsecase
}

func NewUserMiddleware(uu domain.UserUsecase) *UserMiddleware {
	middleware := &UserMiddleware{
		UserUsecase: uu,
	}
	return middleware
}

func (uhm *UserMiddleware) MustAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil || tokenString == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var user domain.User

	idFloat, ok := claims["sub"].(float64)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	id := uint64(idFloat)

	user, err = uhm.UserUsecase.GetByIDRaw(id)

	if user.ID == 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("user", user)
	c.Next()
}
