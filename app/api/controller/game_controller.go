package controller

import (
	"backend/event-service-platform/app/api/client/request"
	"backend/event-service-platform/app/api/client/response"
	"backend/event-service-platform/app/database/constant/currency"
	"backend/event-service-platform/app/database/entity"
	"backend/event-service-platform/app/internal/runtime"
	"backend/event-service-platform/app/manager"
	"backend/event-service-platform/app/pkg/jwt"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type GamePlayController struct {
	res      runtime.Resource
	managers *manager.Managers
	jwt      jwt.Jwt
}

func NewGamePlayController(managers *manager.Managers, res runtime.Resource) *GamePlayController {
	return &GamePlayController{res: res, managers: managers, jwt: jwt.NewJwt(res.Config.JwtConfig)}
}

// Spin godoc
//
//	@Summary		Spin for items
//	@Description	Use SPIN currency to get random items based on rarity weights
//	@Tags			game
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.SpinRequest	true	"Spin request"
//	@Success		200		{object}	response.SpinResultResponse
//	@Failure		400
//	@Failure		401
//	@Failure		402
//	@Failure		500
//	@Router			/api/v1/game/spin [post]
func (c *GamePlayController) Spin(ec echo.Context) error {
	claims, err := c.jwt.GetClaims(ec)
	if err != nil || claims.UserID == nil {
		return ec.JSON(http.StatusUnauthorized, response.ToErrorResponse(http.StatusUnauthorized, "Authentication required"))
	}

	var req request.SpinRequest
	if err := ec.Bind(&req); err != nil {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Invalid request body"))
	}

	if err := req.Validate(); err != nil {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Validation failed: "+err.Error()))
	}

	// Check if user has enough SPIN balance
	spinCost := int64(1) // Default spin cost
	balances, err := c.managers.UserBalanceManager.GetBalanceByCurrency(ec.Request().Context(), *claims.UserID, currency.SPIN)
	if err != nil {
		return ec.JSON(http.StatusInternalServerError, response.ToErrorResponse(http.StatusInternalServerError, "Failed to check balance"))
	}

	spinBalance := balances[string(currency.SPIN)]
	if spinBalance < spinCost {
		return ec.JSON(http.StatusPaymentRequired, response.ToErrorResponse(http.StatusPaymentRequired, "Insufficient SPIN balance"))
	}

	// TODO: Implement actual spinning logic with rarity weights
	// For now, return a mock response
	result := response.SpinResultResponse{
		Success:    true,
		Message:    "Spin successful!",
		SpinCost:   spinCost,
		NewBalance: spinBalance - spinCost,
	}

	return ec.JSON(http.StatusOK, response.ToSuccessResponse(result))
}

// ClaimDailyReward godoc
//
//	@Summary		Claim daily reward
//	@Description	Claim daily reward in COIN or SPIN currency
//	@Tags			game
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.ClaimDailyRewardRequest	true	"Daily reward request"
//	@Success		200		{object}	map[string]int64
//	@Failure		400
//	@Failure		401
//	@Failure		409
//	@Failure		500
//	@Router			/api/v1/game/daily-reward [post]
func (c *GamePlayController) ClaimDailyReward(ec echo.Context) error {
	claims, err := c.jwt.GetClaims(ec)
	if err != nil || claims.UserID == nil {
		return ec.JSON(http.StatusUnauthorized, response.ToErrorResponse(http.StatusUnauthorized, "Authentication required"))
	}

	var req request.ClaimDailyRewardRequest
	if err := ec.Bind(&req); err != nil {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Invalid request body"))
	}

	if err := req.Validate(); err != nil {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Validation failed: "+err.Error()))
	}

	// TODO: Check if user already claimed daily reward today
	// For now, just add the reward
	rewardAmount := int64(100) // Default daily reward amount

	// TODO: Implement actual daily reward logic with transaction
	_ = rewardAmount

	return ec.JSON(http.StatusNotImplemented, response.ToErrorResponse(http.StatusNotImplemented, "Not implemented yet"))
}

// GetGameStats godoc
//
//	@Summary		Get game statistics
//	@Description	Get user's game statistics
//	@Tags			game
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.GameStatsResponse
//	@Failure		401
//	@Failure		500
//	@Router			/api/v1/game/stats [get]
func (c *GamePlayController) GetGameStats(ec echo.Context) error {
	claims, err := c.jwt.GetClaims(ec)
	if err != nil || claims.UserID == nil {
		return ec.JSON(http.StatusUnauthorized, response.ToErrorResponse(http.StatusUnauthorized, "Authentication required"))
	}

	// TODO: Implement actual stats calculation
	stats := response.GameStatsResponse{
		TotalCollections: 0,
		TotalItems:       0,
		UniqueItems:      0,
		TotalSpins:       0,
	}

	return ec.JSON(http.StatusOK, response.ToSuccessResponse(stats))
}

// Helper function for weighted random selection
func (c *GamePlayController) selectRandomItem(items []entity.Item, rarityConfigs map[uuid.UUID]entity.RarityConfig) *entity.Item {
	if len(items) == 0 {
		return nil
	}

	// Calculate total weight
	totalWeight := 0
	for _, item := range items {
		if item.RarityConfigID != nil {
			if rarity, exists := rarityConfigs[*item.RarityConfigID]; exists {
				totalWeight += rarity.DropWeight
			}
		}
	}

	if totalWeight == 0 {
		// If no weights, return random item
		return &items[rand.Intn(len(items))]
	}

	// Weighted random selection
	randomWeight := rand.Intn(totalWeight)
	currentWeight := 0

	for _, item := range items {
		if item.RarityConfigID != nil {
			if rarity, exists := rarityConfigs[*item.RarityConfigID]; exists {
				currentWeight += rarity.DropWeight
				if randomWeight < currentWeight {
					return &item
				}
			}
		}
	}

	return &items[len(items)-1] // Fallback
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
