package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func NewRouter(
	app *gin.Engine,
	logger *zerolog.Logger,
) {
	r := app.Group("/api")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
