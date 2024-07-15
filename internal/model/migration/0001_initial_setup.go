package migration

import (
	"strconv"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"fliqt/internal/model"
)

// ID, Secret
var userSeedData = [][]string{
	// HR
	{"cqan84gcvavjif3csp4g", "UGLOBAFSYEIDW52JGKUEFEQFEB3RZFYL", "hr"},

	// Interviewer
	{"cqanb5gcvavjneudu13g", "7KIHH3TKGHNS67UHG4JLS5QPYN4SKTQC", "interviewer"},

	// Candidate
	{"cqanbg8cvavjpljmh7pg", "TXMJIAOMR42PQP2A5JWC7SPOIHEKI3X2", "candidate"},
}

// ID, Title, Company, JobType, SalaryMin, SalaryMax
var jobSeedData = [][]string{
	{"cqanjjocvavk14kcpc9g", "Software Engineer", "Google", "full-time", "100000", "200000"},
	{"cqankhocvavk4lkbb63g", "Sr. Software Engineer", "Google", "full-time", "200000", "300000"},
	{"cqank1gcvavk2ohtbi90", "Software Engineer", "Facebook", "full-time", "100000", "200000"},
	{"cqanl4gcvavk5aka0am0", "Infrastructure Engineer", "Facebook", "full-time", "150000", "200000"},
	{"cqank6ocvavk3bku7rhg", "Software Engineer", "Amazon", "full-time", "100000", "200000"},
	{"cqanlbocvavk60s8l840", "Designer Manager", "Amazon", "full-time", "100000", "200000"},
	{"cqank8gcvavk3o6tsqb0", "Software Engineer", "Apple", "full-time", "100000", "200000"},
}

func Migration0001() *gormigrate.Migration {
	type Job struct {
		model.Base

		Title     string `gorm:"not null"`
		Company   string `gorm:"not null"`
		JobType   string `gorm:"type:enum('full-time', 'part-time', 'contract');default:'full-time';index:idx_job_job_type"`
		SalaryMin int    `gorm:"not null;index:idx_job_salary_min"`
		SalaryMax int    `gorm:"not null;index:idx_job_salary_max"`
	}

	type User struct {
		model.Base

		Role       string `gorm:"type:enum('hr', 'interviewer', 'candidate');default:'candidate';index:idx_user_role"`
		TotpSecret string `gorm:"not null"`
	}

	type Application struct {
		model.Base

		JobID           string `gorm:"not null;index:idx_application_job_id"`
		Job             Job    `gorm:"foreignKey:JobID"`
		UserID          string `gorm:"not null;index:idx_application_user_id"`
		User            User   `gorm:"foreignKey:UserID"`
		Status          string `gorm:"type:enum('accepted', 'pending', 'rejected');default:'accepted';index:idx_application_status"`
		ResumeObjectKey string `gorm:"not null"`
	}

	return &gormigrate.Migration{
		ID: "0001",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Migrator().CreateTable(&Job{}); err != nil {
				return err
			}
			if err := tx.Exec("CREATE FULLTEXT INDEX idx_job_title_company ON jobs(title, company)").Error; err != nil {
				return err
			}
			if err := tx.Migrator().CreateTable(&User{}); err != nil {
				return err
			}
			if err := tx.Migrator().CreateTable(&Application{}); err != nil {
				return err
			}

			// Add seed data for users.
			users := []User{}
			for _, data := range userSeedData {
				users = append(users, User{
					Base: model.Base{
						ID: data[0],
					},
					TotpSecret: data[1],
					Role:       data[2],
				})
			}

			if err := tx.Create(&users).Error; err != nil {
				return err
			}

			// Add seed data for jobs.
			jobs := []Job{}
			for _, data := range jobSeedData {
				min, _ := strconv.Atoi(data[4])
				max, _ := strconv.Atoi(data[5])

				jobs = append(jobs, Job{
					Base: model.Base{
						ID: data[0],
					},
					Title:     data[1],
					Company:   data[2],
					JobType:   data[3],
					SalaryMin: min,
					SalaryMax: max,
				})
			}
			if err := tx.Create(&jobs).Error; err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropTable(&Application{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&User{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&Job{}); err != nil {
				return err
			}

			return nil
		},
	}
}
