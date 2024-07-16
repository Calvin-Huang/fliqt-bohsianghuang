package repository

import (
	"context"
	"fliqt/internal/model"
	"fliqt/internal/util"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog"
)

func TestApplicationRepositroyListApplications(t *testing.T) {
	db, mock, cleanup := util.SetupMockDB(t)
	defer cleanup()

	logger := zerolog.New(zerolog.ConsoleWriter{Out: nil})

	repo := NewApplicationRepository(db, &logger)

	t.Run("KeywordSearch", func(t *testing.T) {
		mock.ExpectQuery("SELECT count\\(\\*\\) FROM `applications` JOIN jobs ON jobs.id =  applications.job_id WHERE MATCH\\(jobs\\.title, jobs\\.company\\) AGAINST \\(\\?\\) AND `applications`\\.`deleted_at` IS NULL").
			WithArgs("Hello World").
			WillReturnRows(
				sqlmock.NewRows([]string{"count(*)"}).AddRow(1),
			)

		mock.ExpectQuery("SELECT applications.id, applications.job_id, jobs.title as job_title, jobs.company as company, applications.status, applications.resume_object_key, applications.created_at, applications.updated_at FROM `applications` JOIN jobs ON jobs.id =  applications.job_id WHERE MATCH\\(jobs\\.title, jobs\\.company\\) AGAINST \\(\\?\\) AND `applications`\\.`deleted_at` IS NULL ORDER BY id DESC LIMIT \\?").
			WithArgs("Hello World", 10).
			WillReturnRows(
				sqlmock.
					NewRows([]string{"id", "job_id", "job_title", "company", "status", "resume_object_key", "created_at", "updated_at"}).
					AddRow("1", "1", "New Job", "Hello World", "applied", "resume_object_key", "2021-01-01", "2021-01-01"),
			)

		result, err := repo.ListApplications(context.TODO(), ApplicationFilterParams{
			Keyword: "Hello World",
			PaginationParams: model.PaginationParams{
				PageSize: 10,
			},
		})
		if err != nil {
			t.Fatalf("Error: %s", err)
		}

		if result.Items[0].JobTitle != "New Job" {
			t.Errorf("Expected job title to be 'Hello World', got %s", result.Items)
		}

		if result.Items[0].Company != "Hello World" {
			t.Errorf("Expected company to be 'Hello World', got %s", result.Items)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	})

	t.Run("StatusFilter", func(t *testing.T) {
		mock.ExpectQuery("SELECT count\\(\\*\\) FROM `applications` JOIN jobs ON jobs.id =  applications.job_id WHERE applications\\.status = \\? AND `applications`\\.`deleted_at` IS NULL").
			WithArgs("applied").
			WillReturnRows(
				sqlmock.NewRows([]string{"count(*)"}).AddRow(1),
			)

		mock.ExpectQuery("SELECT applications.id, applications.job_id, jobs.title as job_title, jobs.company as company, applications.status, applications.resume_object_key, applications.created_at, applications.updated_at FROM `applications` JOIN jobs ON jobs.id =  applications.job_id WHERE applications.status = \\? AND `applications`\\.`deleted_at` IS NULL ORDER BY id DESC LIMIT \\?").
			WithArgs("applied", 10).
			WillReturnRows(
				sqlmock.
					NewRows([]string{"id", "job_id", "job_title", "company", "status", "resume_object_key", "created_at", "updated_at"}).
					AddRow("1", "1", "New Job", "Hello World", "applied", "resume_object_key", "2021-01-01", "2021-01-01"),
			)

		result, err := repo.ListApplications(context.TODO(), ApplicationFilterParams{
			Status: "applied",
			PaginationParams: model.PaginationParams{
				PageSize: 10,
			},
		})
		if err != nil {
			t.Fatalf("Error: %s", err)
		}

		if result.Items[0].JobTitle != "New Job" {
			t.Errorf("Expected job title to be 'Hello World', got %s", result.Items)
		}

		if result.Items[0].Company != "Hello World" {
			t.Errorf("Expected company to be 'Hello World', got %s", result.Items)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	})
}
