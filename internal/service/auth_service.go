package service

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"

	"fliqt/internal/model"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrFailedTOTP   = errors.New("failed to verify TOTP")
)

type AuthServiceInterface interface {
	CurrentUser(ctx *gin.Context) (*model.User, error)
	VerifyTOTP(ctx *gin.Context, secret string, passcode string) error
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
		return nil, ErrUnauthorized
	}

	var user model.User

	if err := s.db.WithContext(ctx).Where("id", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUnauthorized
		}
		return nil, err
	}

	return &user, nil
}

func (s *AuthService) VerifyTOTP(ctx *gin.Context, secret string, passcode string) error {
	if !totp.Validate(passcode, secret) {
		return ErrFailedTOTP
	}

	return nil
}
