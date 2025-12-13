package controller

import (
	"net/http"

	"errors"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/api/client/request"
	"github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/api/client/response"
	"github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/database/constant/role"
	"github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/internal/runtime"
	"github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/manager"
	"github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/pkg/jwt"
	utilcookie "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/pkg/util/cookie"
)

type AuthController struct {
	res      runtime.Resource
	managers *manager.Managers
	jwt      jwt.Jwt
}

func NewAuthController(managers *manager.Managers, res runtime.Resource) *AuthController {
	jwtService := jwt.NewJwt(res.Config.JwtConfig)
	return &AuthController{
		res:      res,
		managers: managers,
		jwt:      jwtService,
	}
}

// Register godoc
//
//	@Summary        Register user
//	@Description    Create a new account with email and password
//	@Tags           auth
//	@Accept         json
//	@Produce        json
//	@Param          request body        request.RegisterRequest true "Registration"
//	@Success        200
//	@Failure        400
//	@Failure        409
//	@Failure        500
//	@Router         /api/v1/auth/register [post]
func (c *AuthController) Register(ec echo.Context) error {
	ctx := ec.Request().Context()
	var req request.RegisterRequest
	if err := ec.Bind(&req); err != nil {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Invalid request"))
	}

	if err := c.managers.AuthManager.Register(ctx, req); err != nil {
		if errors.Is(err, manager.ErrEmailAlreadyExists) || errors.Is(err, manager.ErrUsernameAlreadyExisted) {
			return ec.JSON(http.StatusConflict, response.ToErrorResponse(http.StatusConflict, err.Error()))
		}
		return ec.JSON(http.StatusInternalServerError, response.ToErrorResponse(http.StatusInternalServerError, "Internal server error"))
	}
	return ec.JSON(http.StatusOK, response.ToSuccessResponse("registered"))
}

// Login godoc
//
//	@Summary		User login
//	@Description	Authenticate user with username and password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.AuthUserRequest	true	"Login credentials"
//	@Success		200		{object}	response.AuthResponse
//	@Failure		400
//	@Failure		401
//	@Failure		500
//	@Router			/api/v1/auth/login [post]
func (c *AuthController) Login(ec echo.Context) error {
	ctx := ec.Request().Context()
	var req request.AuthUserRequest
	if err := ec.Bind(&req); err != nil {
		c.res.Logger.Error("Failed to bind request", zap.Error(err))
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Invalid request format"))
	}

	res, err := c.managers.AuthManager.Login(ctx, req)
	if err != nil {
		c.res.Logger.Error("Login failed", zap.Error(err))
		if errors.Is(err, manager.ErrInvalidCredentials) {
			return ec.JSON(http.StatusUnauthorized, response.ToErrorResponse(http.StatusUnauthorized, "Invalid credentials"))
		}
		return ec.JSON(http.StatusInternalServerError, response.ToErrorResponse(http.StatusInternalServerError, "Internal server error"))
	}

	cookie := utilcookie.NewRefreshTokenCookie(ec.Request(), res.RefreshToken, c.res.Config.JwtConfig.RefreshExpiration)
	ec.SetCookie(cookie)
	res.RefreshToken = ""
	return ec.JSON(http.StatusOK, response.ToSuccessResponse(res))
}

// Wallet/SIWE endpoints are removed in the simplified auth flow.

// RefreshToken godoc
//
//	@Summary		Refresh access token
//	@Description	Get new access token using refresh token from cookie
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	response.AuthResponse
//	@Failure		401
//	@Failure		500
//	@Router			/api/v1/auth/refresh-token [post]
func (c *AuthController) RefreshToken(ec echo.Context) error {
	rtCookie, errCookie := ec.Cookie("refresh_token")
	if errCookie != nil || rtCookie == nil || rtCookie.Value == "" {
		return ec.JSON(http.StatusUnauthorized, response.ToErrorResponse(http.StatusUnauthorized, "Missing refresh token"))
	}
	authResp, err := c.managers.AuthManager.RefreshToken(ec.Request().Context(), request.RefreshTokenRequest{RefreshToken: rtCookie.Value})
	if err != nil {
		c.res.Logger.Error("Token refresh failed", zap.Error(err))
		if errors.Is(err, manager.ErrInvalidRefreshToken) || errors.Is(err, manager.ErrRefreshTokenRevoked) || errors.Is(err, manager.ErrRefreshTokenExpired) {
			return ec.JSON(http.StatusUnauthorized, response.ToErrorResponse(http.StatusUnauthorized, err.Error()))
		}
		return ec.JSON(http.StatusInternalServerError, response.ToErrorResponse(http.StatusInternalServerError, "Internal server error"))
	}
	// Rotate cookie with new refresh token
	cookie := utilcookie.NewRefreshTokenCookie(ec.Request(), authResp.RefreshToken, c.res.Config.JwtConfig.RefreshExpiration)
	ec.SetCookie(cookie)
	authResp.RefreshToken = ""
	return ec.JSON(http.StatusOK, response.ToSuccessResponse(authResp))
}

// Logout godoc
//
//	@Summary		User logout
//	@Description	Revoke refresh token and logout user using refresh token from cookie
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Failure		401
//	@Failure		500
//	@Router			/api/v1/auth/logout [post]
func (c *AuthController) Logout(ec echo.Context) error {
	rtCookie, errCookie := ec.Cookie("refresh_token")
	if errCookie != nil || rtCookie == nil || rtCookie.Value == "" {
		return ec.JSON(http.StatusUnauthorized, response.ToErrorResponse(http.StatusUnauthorized, "Missing refresh token"))
	}
	err := c.managers.AuthManager.Logout(ec.Request().Context(), request.LogoutRequest{RefreshToken: rtCookie.Value})
	if err != nil {
		c.res.Logger.Error("Logout failed", zap.Error(err))
		return ec.JSON(http.StatusInternalServerError, response.ToErrorResponse(http.StatusInternalServerError, "Internal server error"))
	}

	cookie := utilcookie.ExpireCookie("refresh_token")
	ec.SetCookie(cookie)
	return ec.JSON(http.StatusOK, response.ToSuccessResponse("Logged out successfully"))
}

// Me godoc
//
//	@Summary		Get token principal info
//	@Description	Return identity info from the provided access token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.MeResponse
//	@Failure		401
//	@Router			/api/v1/auth/me [get]
func (c *AuthController) Me(ec echo.Context) error {
	claims, err := c.jwt.GetClaims(ec)
	if err != nil {
		c.res.Logger.Error("Failed to get claims", zap.Error(err))
		return ec.JSON(http.StatusUnauthorized, response.ToErrorResponse(http.StatusUnauthorized, "Authentication required"))
	}
	if claims.UserID == nil || claims.Username == nil {
		return ec.JSON(http.StatusUnauthorized, response.ToErrorResponse(http.StatusUnauthorized, "Authentication required"))
	}

	// Convert role string to role.Role type
	var userRole role.Role
	if claims.Role != nil {
		userRole = role.Role(*claims.Role)
	} else {
		userRole = role.User // default role
	}

	// Convert email verified and phone verified to bool with defaults
	emailVerified := false
	if claims.EmailVerified != nil {
		emailVerified = *claims.EmailVerified
	}

	phoneVerified := false
	if claims.PhoneVerified != nil {
		phoneVerified = *claims.PhoneVerified
	}

	meResponse := response.MeResponse{
		ID:            *claims.UserID,
		Username:      *claims.Username,
		Email:         claims.Email,
		PhoneNumber:   claims.PhoneNumber,
		Role:          userRole,
		EmailVerified: emailVerified,
		PhoneVerified: phoneVerified,
		LastLoginAt:   claims.LastLoginAt,
	}

	return ec.JSON(http.StatusOK, response.ToSuccessResponse(meResponse))
}
