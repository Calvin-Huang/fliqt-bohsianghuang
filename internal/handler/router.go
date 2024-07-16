package handler

import (
	"fliqt/internal/model"
	"fliqt/internal/repository"
	"fliqt/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"fliqt/config"
)

func NewRouter(
	cfg *config.Config,
	app *gin.Engine,
	logger *zerolog.Logger,
	jobRepo *repository.JobRepository,
	authService service.AuthServiceInterface,
	s3Service service.S3ServiceInterface,
) {
	r := app.Group("/api")

	jobHandler := NewJobHandler(jobRepo, logger)
	r.GET("/jobs", jobHandler.ListJobs)
	r.GET("/jobs/:id", jobHandler.GetJob)
	r.POST("/jobs", AuthHandler(authService, []model.UserRole{model.RoleHR}), jobHandler.CreateJob)
	r.PUT("/jobs/:id", AuthHandler(authService, []model.UserRole{model.RoleHR}), jobHandler.UpdateJob)
	r.DELETE("/jobs/:id", AuthHandler(authService, []model.UserRole{model.RoleHR}), jobHandler.DeleteJob)

	fileHandler := NewFileHandler(cfg, authService, s3Service)
	r.POST("/files", AuthHandler(authService, []model.UserRole{model.RoleCandidate}), fileHandler.GetUploadInfo)
	r.GET("/files/*object_key", AuthHandler(authService, []model.UserRole{model.RoleHR, model.RoleInteviewer, model.RoleCandidate}), fileHandler.GetDownloadInfo)
}
