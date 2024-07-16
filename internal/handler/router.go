package handler

import (
	"fliqt/internal/model"
	"fliqt/internal/repository"
	"fliqt/internal/service"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func NewRouter(
	app *gin.Engine,
	s3Client *s3.Client,
	logger *zerolog.Logger,
	jobRepo *repository.JobRepository,
	authService service.AuthServiceInterface,
) {
	r := app.Group("/api")

	jobHandler := NewJobHandler(jobRepo, logger)
	r.GET("/jobs", jobHandler.ListJobs)
	r.GET("/jobs/:id", jobHandler.GetJob)
	r.POST("/jobs", AuthHandler(authService, []model.UserRole{model.RoleHR}), jobHandler.CreateJob)
	r.PUT("/jobs/:id", AuthHandler(authService, []model.UserRole{model.RoleHR}), jobHandler.UpdateJob)
	r.DELETE("/jobs/:id", AuthHandler(authService, []model.UserRole{model.RoleHR}), jobHandler.DeleteJob)
}
