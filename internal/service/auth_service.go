package service

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"fliqt/internal/model"
)

type AuthServiceInterface interface {
	CurrentUser(ctx *gin.Context) (*model.User, error)
}

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{
		db,
	}
}

func (s *AuthService) CurrentUser(ctx *gin.Context) (*model.User, error) {
	userID := ctx.GetHeader("X-FLIQT-USER")

	if userID == "" {
		return nil, gorm.ErrRecordNotFound
	}

	var user model.User

	if err := s.db.WithContext(ctx).Where("id", userID).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
