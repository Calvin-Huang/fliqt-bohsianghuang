package model

type JobType string

const (
	JobTypeFullTime JobType = "full-time"
	JobTypePartTime JobType = "part-time"
	JobTypeContract JobType = "contract"
)

type Job struct {
	Base

	// Title and Company have a full-text index
	Title     string  `gorm:"not null"`
	Company   string  `gorm:"not null"`
	JobType   JobType `gorm:"type:enum('full-time', 'part-time', 'contract');default:'full-time';index:idx_job_job_type"`
	SalaryMin int     `gorm:"not null;index:idx_job_salary_min"`
	SalaryMax int     `gorm:"not null;index:idx_job_salary_max"`
}
