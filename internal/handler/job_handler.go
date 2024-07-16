package handler

import (
	"fliqt/internal/repository"
	"fliqt/internal/util"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type JobHandler struct {
	repo   *repository.JobRepository
	logger *zerolog.Logger
}

func NewJobHandler(
	repo *repository.JobRepository,
	logger *zerolog.Logger,
) *JobHandler {
	return &JobHandler{
		repo:   repo,
		logger: logger,
	}
}

// ListJobs is a handler for listing all jobs.
func (h *JobHandler) ListJobs(ctx *gin.Context) {
	var filterParams repository.JobFilterParams
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

	accounts, err := h.repo.ListJobs(tracerCtx, filterParams)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

// GetJob is a handler for getting job details.
func (h *JobHandler) GetJob(ctx *gin.Context) {
	tracerCtx, span := tracer.Start(ctx.Request.Context(), util.GetSpanNameFromCaller())
	defer span.End()

	id := ctx.Param("id")
	if id == "" {
		ctx.Error(ErrNotFound)
		return
	}
	account, err := h.repo.GetJobByID(tracerCtx, id)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (h *JobHandler) CreateJob(ctx *gin.Context) {
	var req repository.CreateJobDTO
	if err := ctx.BindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	tracerCtx, span := tracer.Start(
		ctx.Request.Context(),
		util.GetSpanNameFromCaller(),
		trace.WithAttributes(
			attribute.String("request", fmt.Sprintf("%+v", req)),
		),
	)
	defer span.End()

	if err := req.Validate(); err != nil {
		ctx.Error(err)
		return
	}

	job, err := h.repo.CreateJob(tracerCtx, req)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, job)
}
