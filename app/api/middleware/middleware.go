package middleware

import (
	"github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/internal/runtime"

	"github.com/labstack/echo/v4"
)

type Middleware struct {
	JwtAuthentication       JwtAuthentication
	ApiKeyAuthentication    ApiKeyAuthentication
	HttpBasicAuthentication HttpBasicAuthentication
	AdminAuth               *AdminAuthMiddleware
}

func NewMiddleware(res runtime.Resource) *Middleware {
	return &Middleware{
		JwtAuthentication:       NewJwtAuthentication(res),
		ApiKeyAuthentication:    NewApiKeyAuthentication(res),
		HttpBasicAuthentication: NewHttpBasicAuthentication(res),
		AdminAuth:               NewAdminAuthMiddleware(res),
	}
}

func (m *Middleware) RequireAuth() echo.MiddlewareFunc {
	return m.JwtAuthentication.RequireAuth()
}

func (m *Middleware) RequireRole(requiredRole string) echo.MiddlewareFunc {
	return m.JwtAuthentication.RequireRole(requiredRole)
}

func (m *Middleware) RequireAdmin() echo.MiddlewareFunc {
	return m.AdminAuth.RequireAdmin()
}
