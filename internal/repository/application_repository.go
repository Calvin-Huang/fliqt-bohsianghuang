package repository

import (
	"context"
	"fliqt/internal/model"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type ApplicationRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

func NewApplicationRepository(
	db *gorm.DB,
	logger *zerolog.Logger,
) *ApplicationRepository {
	return &ApplicationRepository{
		db:     db,
		logger: logger,
	}
}

type ApplicationFilterParams struct {
	model.PaginationParams

	Status  string `form:"status,omitempty"`
	Keyword string `form:"keyword,omitempty"`

	UserID *string `form:",omitempty"`
	JobID  *string `form:",omitempty"`
}

type ApplicationResponseDTO struct {
	ID              string `json:"id"`
	JobID           string `json:"job_id"`
	JobTitle        string `json:"job_title"`
	Company         string `json:"company"`
	UserID          string `json:"user_id"`
	Status          string `json:"status"`
	ResumeObjectKey string `json:"resume_object_key"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

// ListApplications returns a list of applications
func (r *ApplicationRepository) ListApplications(ctx context.Context, filterParams ApplicationFilterParams) (model.PaginationResponse[ApplicationResponseDTO], error) {
	var applications []ApplicationResponseDTO
	query := r.db.WithContext(ctx).Model(&model.Application{}).
		Select(`
			applications.id,
			applications.job_id,
			jobs.title as job_title,
			jobs.company as company,
			applications.status,
			applications.resume_object_key,
			applications.created_at,
			applications.updated_at
		`).
		Joins("JOIN jobs ON jobs.id =  applications.job_id").
		Order("id DESC")

	if filterParams.Status != "" {
		query = query.Where("applications.status = ?", filterParams.Status)
	}

	if filterParams.Keyword != "" {
		// Search Job's title or company
		query = query.Where("MATCH(jobs.title, jobs.company) AGAINST (?)", filterParams.Keyword)
	}

	if filterParams.UserID != nil {
		query = query.Where("applications.user_id = ?", *filterParams.UserID)
	}

	if filterParams.JobID != nil {
		query = query.Where("applications.job_id = ?", *filterParams.JobID)
	}

	var total int64
	var result model.PaginationResponse[ApplicationResponseDTO]

	if err := query.Count(&total).Error; err != nil {
		return model.PaginationResponse[ApplicationResponseDTO]{}, err
	}

	if filterParams.NextToken != "" {
		query = query.Where("id < ?", filterParams.NextToken)
	}

	query = query.Limit(filterParams.PageSize)

	if err := query.Find(&applications).Error; err != nil {
		return model.PaginationResponse[ApplicationResponseDTO]{}, err
	}

	result.Total = total
	result.Items = applications

	return result, nil
}

type CreateApplicationDTO struct {
	JobID           string `json:"job_id" binding:"required"`
	UserID          string `json:"user_id" binding:"required"`
	ResumeObjectKey string `json:"resume_object_key" binding:"required"`
}

func (r *ApplicationRepository) CreateApplication(ctx context.Context, dto CreateApplicationDTO) (model.Application, error) {
	application := model.Application{
		JobID:           dto.JobID,
		UserID:          dto.UserID,
		ResumeObjectKey: dto.ResumeObjectKey,
	}
	if err := r.db.WithContext(ctx).Create(&application).Error; err != nil {
		return model.Application{}, err
	}

	return application, nil
}
