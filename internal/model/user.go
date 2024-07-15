package model

type UserRole string

const (
	RoleHR         UserRole = "hr"
	RoleInteviewer UserRole = "interviewer"
	RoleCandidate  UserRole = "candidate"
)

type User struct {
	Base

	Role       UserRole `gorm:"type:enum('hr', 'interviewer', 'candidate');default:'candidate';index:idx_user_role"`
	TotpSecret string   `gorm:"not null"`
}
