package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/api/client/request"
	"github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/api/client/response"
	"github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/service/social"
    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/database/constant/currency"
    txconst "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/database/constant/transaction"
	"github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/internal/runtime"
	"github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/manager"
	"github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/pkg/jwt"
)

type UserController struct {
	res      runtime.Resource
	managers *manager.Managers
	jwt      jwt.Jwt
}

func NewUserController(managers *manager.Managers, res runtime.Resource) *UserController {
	return &UserController{res: res, managers: managers, jwt: jwt.NewJwt(res.Config.JwtConfig)}
}

// GetUserProfile gets a user profile by username
//
//	@Summary		Get user profile by username
//	@Description	Get user profile information including followers, following count, etc.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			username	path	string	true	"Username"
//	@Success		200		{object}	response.UserProfileResponse
//	@Failure		400
//	@Failure		401
//	@Failure		404
//	@Failure		500
//	@Router			/api/v1/users/profiles/{username} [get]
func (c *UserController) GetUserProfile(ec echo.Context) error {
	var req request.GetUserProfileRequest
	req.Username = ec.Param("username")
	if req.Username == "" {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Invalid request"))
	}

	// Get current user ID from JWT claims
	claims, err := c.jwt.GetClaims(ec)
	if err != nil || claims.UserID == nil {
		return ec.JSON(http.StatusUnauthorized, response.ToErrorResponse(http.StatusUnauthorized, "Authentication required"))
	}

	profile, err := c.managers.UserManager.GetUserProfile(ec.Request().Context(), &req, claims.UserID)
	if err != nil {
		return ec.JSON(http.StatusInternalServerError, response.ToErrorResponse(http.StatusInternalServerError, "Failed to get user profile"))
	}

	return ec.JSON(http.StatusOK, response.ToSuccessResponse(profile))
}

// GetBalances godoc
//
//  @Summary        Get user balances
//  @Description    Return current balances for the authenticated user; optionally filter by currency
//  @Tags           users
//  @Accept          json
//  @Produce        json
//  @Param          currency  query   string  false  "Filter by currency"   Enums(COIN,SPIN)
//  @Success        200     {object}    map[string]int64
//  @Failure        400
//  @Failure        401
//  @Failure        500
//  @Router         /api/v1/users/balances [get]
func (c *UserController) GetBalances(ec echo.Context) error {
	claims, err := c.jwt.GetClaims(ec)
	if err != nil || claims.UserID == nil {
		return ec.JSON(http.StatusUnauthorized, response.ToErrorResponse(http.StatusUnauthorized, "Authentication required"))
	}

	currencyParam := ec.QueryParam("currency")

	if currencyParam != "" {
		cur := currency.Currency(currencyParam)
		if cur != currency.COIN && cur != currency.SPIN {
			return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "invalid currency"))
		}
		items, err := c.managers.UserBalanceManager.GetBalanceByCurrency(ec.Request().Context(), *claims.UserID, cur)
		if err != nil {
			return ec.JSON(http.StatusInternalServerError, response.ToErrorResponse(http.StatusInternalServerError, "Internal server error"))
		}
		return ec.JSON(http.StatusOK, response.ToSuccessResponse(items))
	}

	items, err := c.managers.UserBalanceManager.GetAllBalances(ec.Request().Context(), *claims.UserID)
	if err != nil {
		return ec.JSON(http.StatusInternalServerError, response.ToErrorResponse(http.StatusInternalServerError, "Internal server error"))
	}
	return ec.JSON(http.StatusOK, response.ToSuccessResponse(items))
}

// GetBalanceHistory godoc
//
//  @Summary        Get user balance history
//  @Description    List balance change transactions for the authenticated user
//  @Tags           users
//  @Accept         json
//  @Produce        json
//  @Param          currency    query   string  true    "Currency"    Enums(COIN,SPIN)
//  @Param          type        query   string  false   "Transaction type"  Enums(DAILY_REWARD,AD_WATCH,PURCHASE,SPIN_USE)
//  @Success        200     {array}     response.BalanceTransactionResponse
//  @Failure        400
//  @Failure        401
//  @Failure        500
//  @Router         /api/v1/users/balances/history [get]
func (c *UserController) GetBalanceHistory(ec echo.Context) error {
	claims, err := c.jwt.GetClaims(ec)
	if err != nil || claims.UserID == nil {
		return ec.JSON(http.StatusUnauthorized, response.ToErrorResponse(http.StatusUnauthorized, "Authentication required"))
	}

	currencyParam := ec.QueryParam("currency")
	typeParam := ec.QueryParam("type")

	var curPtr *currency.Currency
	if currencyParam == "" {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "currency is required"))
	}
	cur := currency.Currency(currencyParam)
	if cur != currency.COIN && cur != currency.SPIN {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "invalid currency"))
	}
	curPtr = &cur

	var typePtr *string
	if typeParam != "" {
		t := typeParam
		typePtr = &t
	}

	// Convert string type to typed enum if provided
	var typedTypePtr *txconst.TransactionType
	if typePtr != nil {
		tt := txconst.TransactionType(*typePtr)
		typedTypePtr = &tt
	}

	items, err := c.managers.UserBalanceManager.GetHistory(ec.Request().Context(), *claims.UserID, curPtr, typedTypePtr)
	if err != nil {
		return ec.JSON(http.StatusInternalServerError, response.ToErrorResponse(http.StatusInternalServerError, "Internal server error"))
	}
	dtos := make([]response.BalanceTransactionResponse, 0, len(items))
	for _, it := range items {
		dtos = append(dtos, response.BalanceTransactionResponse{
			ID:        it.ID,
			Amount:    it.Amount,
			Currency:  string(it.Currency),
			Type:      string(it.Type),
			Source:    string(it.Source),
			Status:    string(it.Status),
			CreatedAt: it.CreatedAt,
		})
	}
	return ec.JSON(http.StatusOK, response.ToSuccessResponse(dtos))
}

// GetUserEvents gets events created by a user
//
//	@Summary		Get user events
//	@Description	Get events created by a specific user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			username	path	string	true	"Username"
//	@Param			page		query	int		false	"Page number"	default(1)
//	@Param			page_size	query	int		false	"Page size"	default(20)
//	@Success		200		{object}	response.UserEventListResponse
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/api/v1/users/{username}/events [get]
func (c *UserController) GetUserEvents(ec echo.Context) error {
	var req request.GetUserEventsRequest
	req.Username = ec.Param("username")
	if req.Username == "" {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Invalid request"))
	}

	// Parse pagination query params
	if p := ec.QueryParam("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			req.Page = v
		}
	}
	if ps := ec.QueryParam("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil {
			req.PageSize = v
		}
	}

	// Set default pagination values
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	events, err := c.managers.UserManager.GetUserEvents(ec.Request().Context(), &req)
	if err != nil {
		return ec.JSON(http.StatusInternalServerError, response.ToErrorResponse(http.StatusInternalServerError, "Failed to get user events"))
	}

	return ec.JSON(http.StatusOK, response.ToSuccessResponse(events))
}

// GetUserPosts gets posts created by a user
//
//	@Summary		Get user posts
//	@Description	Get posts created by a specific user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			username	path	string	true	"Username"
//	@Param			page		query	int		false	"Page number"	default(1)
//	@Param			page_size	query	int		false	"Page size"	default(20)
//	@Success		200		{object}	response.UserPostListResponse
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/api/v1/users/{username}/posts [get]
func (c *UserController) GetUserPosts(ec echo.Context) error {
	var req request.GetUserPostsRequest
	req.Username = ec.Param("username")
	if req.Username == "" {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Invalid request"))
	}

	// Parse pagination query params
	if p := ec.QueryParam("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			req.Page = v
		}
	}
	if ps := ec.QueryParam("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil {
			req.PageSize = v
		}
	}

	// Set default pagination values
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	posts, err := c.managers.UserManager.GetUserPosts(ec.Request().Context(), &req)
	if err != nil {
		return ec.JSON(http.StatusInternalServerError, response.ToErrorResponse(http.StatusInternalServerError, "Failed to get user posts"))
	}

	return ec.JSON(http.StatusOK, response.ToSuccessResponse(posts))
}

// FollowUser follows a user
//
//	@Summary		Follow a user
//	@Description	Follow a user by username
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			username	path	string	true	"Username to follow"
//	@Success		200
//	@Failure		400
//	@Failure		401
//	@Failure		404
//	@Failure		500
//	@Router			/api/v1/users/{username}/follow [post]
func (c *UserController) FollowUser(ec echo.Context) error {
	var req request.FollowUserRequest
	req.Username = ec.Param("username")
	if req.Username == "" {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Invalid request"))
	}

	claims, err := c.jwt.GetClaims(ec)
	if err != nil || claims.UserID == nil {
		return ec.JSON(http.StatusUnauthorized, response.ToErrorResponse(http.StatusUnauthorized, "Authentication required"))
	}

	// delegate to social-service-platform if available
	if c.managers.SocialClient != nil {
		if sc, ok := c.managers.SocialClient.(social.SocialClient); ok {
			if err := sc.FollowUser(ec.Request().Context(), req.Username, *claims.UserID); err != nil {
				return ec.JSON(http.StatusInternalServerError, response.ToErrorResponse(http.StatusInternalServerError, "Failed to follow user"))
			}
			return ec.JSON(http.StatusOK, response.ToSuccessResponse(map[string]string{"message": "User followed successfully"}))
		}
	}
	// fallback to local manager if social client not configured
	err = c.managers.UserManager.FollowUser(ec.Request().Context(), &req, *claims.UserID)
	if err != nil {
		return ec.JSON(http.StatusInternalServerError, response.ToErrorResponse(http.StatusInternalServerError, "Failed to follow user"))
	}

	return ec.JSON(http.StatusOK, response.ToSuccessResponse(map[string]string{"message": "User followed successfully"}))
}

// UnfollowUser unfollows a user
//
//	@Summary		Unfollow a user
//	@Description	Unfollow a user by username
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			username	path	string	true	"Username to unfollow"
//	@Success		200
//	@Failure		400
//	@Failure		401
//	@Failure		404
//	@Failure		500
//	@Router			/api/v1/users/{username}/unfollow [post]
func (c *UserController) UnfollowUser(ec echo.Context) error {
	var req request.UnfollowUserRequest
	req.Username = ec.Param("username")
	if req.Username == "" {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Invalid request"))
	}

	claims, err := c.jwt.GetClaims(ec)
	if err != nil || claims.UserID == nil {
		return ec.JSON(http.StatusUnauthorized, response.ToErrorResponse(http.StatusUnauthorized, "Authentication required"))
	}

	if c.managers.SocialClient != nil {
		if sc, ok := c.managers.SocialClient.(social.SocialClient); ok {
			if err := sc.UnfollowUser(ec.Request().Context(), req.Username, *claims.UserID); err != nil {
				return ec.JSON(http.StatusInternalServerError, response.ToErrorResponse(http.StatusInternalServerError, "Failed to unfollow user"))
			}
			return ec.JSON(http.StatusOK, response.ToSuccessResponse(map[string]string{"message": "User unfollowed successfully"}))
		}
	}
	// fallback to local manager
	err = c.managers.UserManager.UnfollowUser(ec.Request().Context(), &req, *claims.UserID)
	if err != nil {
		return ec.JSON(http.StatusInternalServerError, response.ToErrorResponse(http.StatusInternalServerError, "Failed to unfollow user"))
	}

	return ec.JSON(http.StatusOK, response.ToSuccessResponse(map[string]string{"message": "User unfollowed successfully"}))
}

// UpdateUserProfile updates user profile
//
//	@Summary		Update user profile
//	@Description	Update authenticated user's profile information
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body	request.UpdateUserProfileRequest	true	"Profile update request"
//	@Success		200		{object}	response.UserProfileResponse
//	@Failure		400
//	@Failure		401
//	@Failure		500
//	@Router			/api/v1/users/profile [put]
func (c *UserController) UpdateUserProfile(ec echo.Context) error {
	var req request.UpdateUserProfileRequest
	if err := ec.Bind(&req); err != nil {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Invalid request"))
	}

	claims, err := c.jwt.GetClaims(ec)
	if err != nil || claims.UserID == nil {
		return ec.JSON(http.StatusUnauthorized, response.ToErrorResponse(http.StatusUnauthorized, "Authentication required"))
	}

	profile, err := c.managers.UserManager.UpdateUserProfile(ec.Request().Context(), &req, *claims.UserID)
	if err != nil {
		return ec.JSON(http.StatusInternalServerError, response.ToErrorResponse(http.StatusInternalServerError, "Failed to update profile"))
	}

	return ec.JSON(http.StatusOK, response.ToSuccessResponse(profile))
}

// UpdateUserSettings updates user settings
//
//	@Summary		Update user settings
//	@Description	Update authenticated user's settings (E2EE, etc.)
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body	request.UpdateUserSettingsRequest	true	"Settings update request"
//	@Success		200
//	@Failure		400
//	@Failure		401
//	@Failure		500
//	@Router			/api/v1/users/settings [put]
func (c *UserController) UpdateUserSettings(ec echo.Context) error {
	var req request.UpdateUserSettingsRequest
	if err := ec.Bind(&req); err != nil {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Invalid request"))
	}

	claims, err := c.jwt.GetClaims(ec)
	if err != nil || claims.UserID == nil {
		return ec.JSON(http.StatusUnauthorized, response.ToErrorResponse(http.StatusUnauthorized, "Authentication required"))
	}

	err = c.managers.UserManager.UpdateUserSettings(ec.Request().Context(), &req, *claims.UserID)
	if err != nil {
		return ec.JSON(http.StatusInternalServerError, response.ToErrorResponse(http.StatusInternalServerError, "Failed to update settings"))
	}

	return ec.JSON(http.StatusOK, response.ToSuccessResponse(map[string]string{"message": "Settings updated successfully"}))
}
