package repository

import (
	"context"
	"errors"
	"fliqt/internal/model"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type JobRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

func NewJobRepository(
	db *gorm.DB,
	logger *zerolog.Logger,
) *JobRepository {
	return &JobRepository{
		db:     db,
		logger: logger,
	}
}

var (
	ErrJobSalaryRange = errors.New("salary_min must be less than or equal to salary_max")
)

type JobFilterParams struct {
	model.PaginationParams

	Keyword   string `form:"keyword,omitempty"`
	SalaryMin int    `form:"salary_min,omitempty"`
	SalaryMax int    `form:"salary_max,omitempty"`
	JobType   string `form:"job_type,omitempty"`
}

type JobResponseDTO struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Company   string `json:"company"`
	JobType   string `json:"job_type"`
	SalaryMin int    `json:"salary_min"`
	SalaryMax int    `json:"salary_max"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ListJobs returns a list of jobs
func (r *JobRepository) ListJobs(ctx context.Context, filterParams JobFilterParams) (model.PaginationResponse[JobResponseDTO], error) {
	var jobs []JobResponseDTO
	query := r.db.WithContext(ctx).Model(&model.Job{}).Order("id DESC")

	if filterParams.Keyword != "" {
		query = query.Where("MATCH(title, company) AGAINST (?)", filterParams.Keyword)
	}
	if filterParams.SalaryMin != 0 {
		query = query.Where("salary_min >= ?", filterParams.SalaryMin)
	}
	if filterParams.SalaryMax != 0 {
		query = query.Where("salary_max <= ?", filterParams.SalaryMax)
	}
	if filterParams.JobType != "" {
		query = query.Where("job_type = ?", filterParams.JobType)
	}

	var total int64
	var result model.PaginationResponse[JobResponseDTO]

	if err := query.Count(&total).Error; err != nil {
		return result, err
	}

	if filterParams.NextToken != "" {
		query = query.Where("id < ?", filterParams.NextToken)
	}

	query = query.Limit(filterParams.PageSize)

	if err := query.Find(&jobs).Error; err != nil {
		return result, err
	}

	result.Total = total
	result.Items = jobs

	if len(jobs) > 0 && len(jobs) == filterParams.PageSize {
		result.NextToken = jobs[len(jobs)-1].ID
	}

	return result, nil
}

// GetJobByID returns a job by its ID
func (r *JobRepository) GetJobByID(ctx context.Context, ID string) (*model.Job, error) {
	var job model.Job
	if err := r.db.WithContext(ctx).Where("id = ?", ID).First(&job).Error; err != nil {
		return nil, err
	}

	return &job, nil
}

type JobValidator interface {
	Validate() error
}

type CreateJobDTO struct {
	Title     string `json:"title" binding:"required"`
	Company   string `json:"company" binding:"required"`
	JobType   string `json:"job_type" binding:"required,oneof=full-time part-time contract"`
	SalaryMin int    `json:"salary_min" binding:"required"`
	SalaryMax int    `json:"salary_max" binding:"required"`
}

func (dto CreateJobDTO) Validate() error {
	if dto.SalaryMin > dto.SalaryMax {
		return ErrJobSalaryRange
	}

	return nil
}

// CreateJob creates a new job
func (r *JobRepository) CreateJob(ctx context.Context, dto CreateJobDTO) (*model.Job, error) {
	job := model.Job{
		Title:     dto.Title,
		Company:   dto.Company,
		JobType:   model.JobType(dto.JobType),
		SalaryMin: dto.SalaryMin,
		SalaryMax: dto.SalaryMax,
	}

	if err := r.db.WithContext(ctx).Create(&job).Error; err != nil {
		return nil, err
	}

	return &job, nil
}

type UpdateJobDTO struct {
	Title     string `json:"title" binding:"required"`
	Company   string `json:"company" binding:"required"`
	JobType   string `json:"job_type" binding:"required,oneof=full-time part-time contract"`
	SalaryMin int    `json:"salary_min" binding:"required"`
	SalaryMax int    `json:"salary_max" binding:"required"`
}

func (dto UpdateJobDTO) Validate() error {
	if dto.SalaryMin > dto.SalaryMax {
		return ErrJobSalaryRange
	}

	return nil
}

// UpdateJob updates a job
func (r *JobRepository) UpdateJob(ctx context.Context, ID string, dto UpdateJobDTO) (*model.Job, error) {
	if err := r.db.WithContext(ctx).Model(&model.Job{}).Where("id = ?", ID).UpdateColumns(dto).Error; err != nil {
		return nil, err
	}

	return r.GetJobByID(ctx, ID)
}

// DeleteJob deletes a job
func (r *JobRepository) DeleteJob(ctx context.Context, ID string) error {
	if err := r.db.WithContext(ctx).Where("id = ?", ID).Delete(&model.Job{}).Error; err != nil {
		return err
	}

	return nil
}
