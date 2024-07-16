package handler

import (
	"fliqt/internal/model"
	"fliqt/internal/repository"
	"fliqt/internal/service"
	"fliqt/internal/util"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type ApplicationHandler struct {
	applicationRepo *repository.ApplicationRepository
	logger          *zerolog.Logger
	authService     service.AuthServiceInterface
}

func NewApplicationHandler(
	applicationRepo *repository.ApplicationRepository,
	logger *zerolog.Logger,
	authService service.AuthServiceInterface,
) *ApplicationHandler {
	return &ApplicationHandler{
		applicationRepo,
		logger,
		authService,
	}
}

// ListApplications returns a list of applications
func (h *ApplicationHandler) ListApplications(ctx *gin.Context) {
	var filterParams repository.ApplicationFilterParams
	if err := ctx.ShouldBindQuery(&filterParams); err != nil {
		ctx.Error(err)
		return
	}

	tracerCtx, span := tracer.Start(
		ctx.Request.Context(),
		util.GetSpanNameFromCaller(),
		trace.WithAttributes(
			attribute.String("query", fmt.Sprintf("%+v", filterParams)),
		),
	)
	defer span.End()

	filterParams.Normalize()

	user, err := h.authService.CurrentUser(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	// Only HR and Interviewer can list all applications
	if user.Role == model.RoleCandidate {
		filterParams.UserID = &user.ID
	}

	applications, err := h.applicationRepo.ListApplications(tracerCtx, filterParams)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, applications)
}
