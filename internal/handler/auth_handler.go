package handler

import (
	"fliqt/internal/model"
	"fliqt/internal/service"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func AuthHandler(authService *service.AuthService, allowedRoles []model.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := authService.CurrentUser(c)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		if !slices.Contains(allowedRoles, user.Role) {
			c.AbortWithError(http.StatusForbidden, ErrForbidden)
			return
		}

		c.Next()
	}
}
