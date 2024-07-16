package repository

import "testing"

func TestCreateJobDTOValidate(t *testing.T) {
	dto := CreateJobDTO{
		Title:     "Software Engineer",
		Company:   "Google",
		JobType:   "full-time",
		SalaryMin: 1000,
		SalaryMax: 2000,
	}

	if err := dto.Validate(); err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	dto = CreateJobDTO{
		Title:     "Software Engineer",
		Company:   "Google",
		JobType:   "full-time",
		SalaryMin: 2000,
		SalaryMax: 1000,
	}

	if err := dto.Validate(); err != ErrJobSalaryRange {
		t.Errorf("expected ErrJobSalaryRange, got %v", err)
	}
}
