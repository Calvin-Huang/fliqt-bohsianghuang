package handler

import (
	"fliqt/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func NewRouter(
	app *gin.Engine,
	logger *zerolog.Logger,
	jobRepo *repository.JobRepository,
) {
	r := app.Group("/api")

	jobHandler := NewJobHandler(jobRepo, logger)
	r.GET("/jobs", jobHandler.ListJobs)
	r.GET("/jobs/:id", jobHandler.GetJob)
}
