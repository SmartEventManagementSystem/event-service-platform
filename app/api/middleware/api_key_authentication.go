package middleware

import (
	"fmt"
	"net/http"

	"backend/event-service-platform/app/api/client/response"
	"backend/event-service-platform/app/internal/runtime"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type ApiKeyAuthentication struct {
	res runtime.Resource
	// In production, API keys should be stored securely (database, vault, etc.)
	// For now, we'll use a configurable map
	validAPIKeys map[string]string // key -> service name
}

func NewApiKeyAuthentication(res runtime.Resource) ApiKeyAuthentication {
	// Initialize with some example API keys (in production, load from secure storage)
	validKeys := map[string]string{
		"service-key-1": "financial-service",
		"service-key-2": "kyc-service",
		"service-key-3": "notification-service",
	}

	return ApiKeyAuthentication{
		res:          res,
		validAPIKeys: validKeys,
	}
}

func (ak *ApiKeyAuthentication) GetName() string {
	return authMethodAPIKey
}

func (ak *ApiKeyAuthentication) CanHandle(c echo.Context) bool {
	apiKey := c.Request().Header.Get(apiKeyHeader)
	return apiKey != ""
}

func (ak *ApiKeyAuthentication) GetAuthenticationMethod(c echo.Context) (string, error) {
	method := c.Get(contextAuthMethod)
	if method == nil {
		return "", fmt.Errorf("authentication method not found in context")
	}

	methodStr, ok := method.(string)
	if !ok {
		return "", fmt.Errorf("invalid authentication method type in context")
	}

	return methodStr, nil
}

func (ak *ApiKeyAuthentication) RequireAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !ak.CanHandle(c) {
				return c.JSON(http.StatusUnauthorized, response.ToErrorResponse(http.StatusUnauthorized, errMsgAuthRequired))
			}

			result, err := ak.Authenticate(c)
			if err != nil {
				ak.res.Logger.Debug("API Key Authentication failed",
					zap.String("handler", ak.GetName()),
					zap.String("error", err.Error()))
				return c.JSON(http.StatusUnauthorized, response.ToErrorResponse(http.StatusUnauthorized, errMsgInvalidAPIKey))
			}

			if !result.Success {
				return c.JSON(http.StatusUnauthorized, response.ToErrorResponse(http.StatusUnauthorized, errMsgInvalidAPIKey))
			}

			ak.SetUserContext(c, result)
			return next(c)
		}
	}
}

func (ak *ApiKeyAuthentication) RequireRole(requiredRole string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Service accounts typically have specific roles
			// For API key authentication, we might want different logic
			// Since services often have elevated permissions, we can either:
			// 1. Always allow service accounts (bypass role check)
			// 2. Check against service-specific roles

			// Option 1: Allow all authenticated service accounts
			if ak.IsServiceAccount(c) {
				ak.res.Logger.Debug("Service account authenticated, bypassing role check")
				return next(c)
			}

			// Option 2: Check role normally (uncomment if needed)
			/*
			   roleInterface, exists := c.Get(contextRole)
			   if !exists {
			       c.JSON(http.StatusForbidden, response.ToErrorResponse(http.StatusForbidden, "Access denied: no role found"))
			       c.Abort()
			       return
			   }

			   userRole, ok := roleInterface.(string)
			   if !ok {
			       c.JSON(http.StatusForbidden, response.ToErrorResponse(http.StatusForbidden, "Access denied: invalid role"))
			       c.Abort()
			       return
			   }

			   if !ak.HasRequiredRole(userRole, requiredRole) {
			       c.JSON(http.StatusForbidden, response.ToErrorResponse(http.StatusForbidden, "Access denied: insufficient permissions"))
			       c.Abort()
			       return
			   }
			*/

			return next(c)
		}
	}
}

func (ak *ApiKeyAuthentication) Authenticate(c echo.Context) (*AuthenticationResult, error) {
	apiKey := c.Request().Header.Get(apiKeyHeader)
	if apiKey == "" {
		return nil, fmt.Errorf(errMsgInvalidAPIKey)
	}

	serviceName, exists := ak.validAPIKeys[apiKey]
	if !exists {
		return nil, fmt.Errorf(errMsgInvalidAPIKey)
	}

	// Parse service UUID
	serviceUUID, err := uuid.Parse(serviceUserID)
	if err != nil {
		return nil, fmt.Errorf("invalid service UUID")
	}

	username := serviceUsername
	serviceNamePtr := &serviceName

	serviceRole := "SERVICE"
	return &AuthenticationResult{
		Success:     true,
		UserID:      &serviceUUID,
		Username:    &username,
		Role:        &serviceRole, // Service accounts get a Service role
		Method:      authMethodAPIKey,
		ServiceName: serviceNamePtr,
	}, nil
}

func (ak *ApiKeyAuthentication) SetUserContext(c echo.Context, result *AuthenticationResult) {
	if result.UserID != nil {
		c.Set(contextUserID, *result.UserID)
		c.Set(contextUserUUID, *result.UserID)
	}
	if result.Username != nil {
		c.Set(contextUsername, *result.Username)
	}
	if result.Email != nil {
		c.Set(contextEmail, *result.Email)
	}
	if result.PhoneNumber != nil {
		c.Set(contextPhoneNumber, *result.PhoneNumber)
	}
	if result.Role != nil {
		c.Set(contextRole, *result.Role)
	}
	if result.EmailVerified != nil {
		c.Set(contextEmailVerified, *result.EmailVerified)
	}
	if result.PhoneVerified != nil {
		c.Set(contextPhoneVerified, *result.PhoneVerified)
	}
	if result.LastLoginAt != nil {
		c.Set(contextLastLoginAt, *result.LastLoginAt)
	}

	// Set an authentication method and service name for tracking
	c.Set(contextAuthMethod, result.Method)
	if result.ServiceName != nil {
		c.Set(contextServiceName, *result.ServiceName)
	}
}

func (ak *ApiKeyAuthentication) CreateErrorResponse(statusCode int, message string) error {
	return echo.NewHTTPError(statusCode, response.ToErrorResponse(statusCode, message))
}

func (ak *ApiKeyAuthentication) HasRequiredRole(userRole string, requiredRole string) bool {
	return userRole == requiredRole
}

// IsServiceAccount checks if the current request is from a service account
func (ak *ApiKeyAuthentication) IsServiceAccount(c echo.Context) bool {
	method, err := ak.GetAuthenticationMethod(c)
	return err == nil && method == authMethodAPIKey
}
