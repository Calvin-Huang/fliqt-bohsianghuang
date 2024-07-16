package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"fliqt/internal/model"
)

type mockedAuthServiceForHR struct {
}

func (m *mockedAuthServiceForHR) CurrentUser(ctx *gin.Context) (*model.User, error) {
	return &model.User{
		Base: model.Base{
			ID: "1",
		},
		Role: model.RoleHR,
	}, nil
}
func (m *mockedAuthServiceForHR) VerifyTOTP(ctx *gin.Context, secret string, passcode string) error {
	return nil
}

func TestAuthHandler(t *testing.T) {
	authService := &mockedAuthServiceForHR{}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	AuthHandler(authService, []model.UserRole{model.RoleHR})(ctx)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code 200, got %d", w.Code)
	}

	AuthHandler(authService, []model.UserRole{model.RoleCandidate})(ctx)
	if w.Code != http.StatusForbidden {
		t.Errorf("expected status code 401, got %d", w.Code)
	}
}
