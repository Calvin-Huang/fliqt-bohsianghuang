package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// Logger is a handler that logs requests using zerolog
func Logger(logger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		logger.Info().
			Str("client_ip", c.ClientIP()).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Str("latency", time.Since(start).String()).
			Str("user_agent", c.Request.UserAgent()).
			Msg("request")
	}
}
