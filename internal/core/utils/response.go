package response

import (
	"cleanstandarts/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func ResponseOk(c *gin.Context, val any, method string, path string, status int) {
	c.JSON(status, gin.H{
		"time":   time.Now().Format(time.UnixDate),
		"path":   path,
		"method": method,
		"value":  val,
	})
}

func ResponseError(c *gin.Context, err error, method string, path string) {
	c.JSON(getStatus(err), gin.H{
		"time":   time.Now().Format(time.UnixDate),
		"path":   path,
		"method": method,
		"error":  err.Error(),
	})
}

func getStatus(err error) int {
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
