package model

type Application struct {
	Base

	JobID           string `gorm:"not null;index:idx_application_job_id"`
	Job             Job    `gorm:"foreignKey:JobID"`
	UserID          string `gorm:"not null;index:idx_application_user_id"`
	User            User   `gorm:"foreignKey:UserID"`
	Status          string `gorm:"type:enum('accepted', 'pending', 'rejected');default:'accepted';index:idx_application_status"`
	ResumeObjectKey string `gorm:"not null"`
}
