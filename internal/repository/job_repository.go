package repository

import (
	"context"
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

type JobFilterParams struct {
	model.PaginationParams

	Keyword   string `form:"keyword,omitempty"`
	SalaryMin int    `form:"salary_min,omitempty"`
	SalaryMax int    `form:"salary_max,omitempty"`
	JobType   string `form:"job_type,omitempty"`
}

// ListJobs returns a list of jobs
func (r *JobRepository) ListJobs(ctx context.Context, filterParams JobFilterParams) (model.PaginationResponse[model.Job], error) {
	var jobs []model.Job
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
	var result model.PaginationResponse[model.Job]

	if err := query.Count(&total).Error; err != nil {
		return result, err
	}

	if filterParams.NextToken != "" {
		query = query.Where("id < ?", filterParams.NextToken)
	}

	if filterParams.PageSize > 0 {
		query = query.Limit(filterParams.PageSize)
	}

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
