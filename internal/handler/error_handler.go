package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrBadRequest = errors.New("bad request")
)

var errStatusMap = map[error]int{
	gorm.ErrRecordNotFound: http.StatusNotFound,
	ErrNotFound:            http.StatusNotFound,
	ErrBadRequest:          http.StatusBadRequest,
}

func ErrorHandler(logger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logger.Error().Err(err).Msg("error occurred")

				status := http.StatusInternalServerError

				var ve validator.ValidationErrors
				if errors.As(err, &ve) {
					status = http.StatusBadRequest
				} else if s, ok := errStatusMap[err.Err]; ok {
					status = s
				}

				c.JSON(status, gin.H{"error": err.Error()})
			}
		}
	}
}

func NotFoundHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
	}
}