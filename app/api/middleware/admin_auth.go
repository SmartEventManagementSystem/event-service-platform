package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/internal/runtime"
	"github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/pkg/jwt"
)

type AdminAuthMiddleware struct {
	res runtime.Resource
	jwt jwt.Jwt
}

func NewAdminAuthMiddleware(res runtime.Resource) *AdminAuthMiddleware {
	return &AdminAuthMiddleware{
		res: res,
		jwt: jwt.NewJwt(res.Config.JwtConfig),
	}
}

// RequireAdmin checks if the authenticated user has admin privileges
func (m *AdminAuthMiddleware) RequireAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, err := m.jwt.GetClaims(c)
			if err != nil || claims.UserID == nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
			}

			if !m.isAdminUser(*claims.UserID) {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Admin access required"})
			}

			return next(c)
		}
	}
}

// isAdminUser checks if the user ID belongs to an admin
// This is a placeholder implementation - in production, you'd check user roles
func (m *AdminAuthMiddleware) isAdminUser(userID uuid.UUID) bool {
	// TODO: Replace with actual admin role checking from database
	// For now, return true to allow all authenticated users for development
	return true
}
