package controller

import (
	"github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/api/client/request"
	"github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/api/client/response"
	"github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/internal/runtime"
	"github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/manager"
	"github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/pkg/jwt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type RarityConfigController struct {
	res      runtime.Resource
	managers *manager.Managers
	jwt      jwt.Jwt
}

func NewRarityConfigController(managers *manager.Managers, res runtime.Resource) *RarityConfigController {
	return &RarityConfigController{res: res, managers: managers, jwt: jwt.NewJwt(res.Config.JwtConfig)}
}

// ListRarityConfigs godoc
//
//	@Summary		List rarity configurations
//	@Description	List all rarity configurations
//	@Tags			rarity-configs
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	response.RarityConfigResponse
//	@Failure		400
//	@Failure		401
//	@Failure		500
//	@Router			/api/v1/rarity-configs [get]
func (c *RarityConfigController) ListRarityConfigs(ec echo.Context) error {
	// TODO: Implement listing logic when manager has the method
	return ec.JSON(http.StatusNotImplemented, response.ToErrorResponse(http.StatusNotImplemented, "Not implemented yet"))
}

// GetRarityConfig godoc
//
//	@Summary		Get rarity configuration by ID
//	@Description	Get a specific rarity configuration
//	@Tags			rarity-configs
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Rarity Config ID"
//	@Success		200	{object}	response.RarityConfigResponse
//	@Failure		400
//	@Failure		401
//	@Failure		404
//	@Failure		500
//	@Router			/api/v1/rarity-configs/{id} [get]
func (c *RarityConfigController) GetRarityConfig(ec echo.Context) error {
	idParam := ec.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Invalid rarity config ID"))
	}

	// TODO: Implement get logic when manager has the method
	_ = id
	return ec.JSON(http.StatusNotImplemented, response.ToErrorResponse(http.StatusNotImplemented, "Not implemented yet"))
}

// CreateRarityConfig godoc
//
//	@Summary		Create rarity configuration
//	@Description	Create a new rarity configuration (admin only)
//	@Tags			rarity-configs
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.CreateRarityConfigRequest	true	"Rarity config data"
//	@Success		201		{object}	response.RarityConfigResponse
//	@Failure		400
//	@Failure		401
//	@Failure		403
//	@Failure		500
//	@Router			/api/v1/rarity-configs [post]
func (c *RarityConfigController) CreateRarityConfig(ec echo.Context) error {
	// TODO: Add admin authorization check
	var req request.CreateRarityConfigRequest
	if err := ec.Bind(&req); err != nil {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Invalid request body"))
	}

	if err := req.Validate(); err != nil {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Validation failed: "+err.Error()))
	}

	// TODO: Implement creation logic when manager has the method
	return ec.JSON(http.StatusNotImplemented, response.ToErrorResponse(http.StatusNotImplemented, "Not implemented yet"))
}

// UpdateRarityConfig godoc
//
//	@Summary		Update rarity configuration
//	@Description	Update an existing rarity configuration (admin only)
//	@Tags			rarity-configs
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"Rarity Config ID"
//	@Param			request	body		request.UpdateRarityConfigRequest	true	"Rarity config data"
//	@Success		200		{object}	response.RarityConfigResponse
//	@Failure		400
//	@Failure		401
//	@Failure		403
//	@Failure		404
//	@Failure		500
//	@Router			/api/v1/rarity-configs/{id} [put]
func (c *RarityConfigController) UpdateRarityConfig(ec echo.Context) error {
	// TODO: Add admin authorization check
	idParam := ec.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Invalid rarity config ID"))
	}

	var req request.UpdateRarityConfigRequest
	if err := ec.Bind(&req); err != nil {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Invalid request body"))
	}

	if err := req.Validate(); err != nil {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Validation failed: "+err.Error()))
	}

	// TODO: Implement update logic when manager has the method
	_ = id
	return ec.JSON(http.StatusNotImplemented, response.ToErrorResponse(http.StatusNotImplemented, "Not implemented yet"))
}

// DeleteRarityConfig godoc
//
//	@Summary		Delete rarity configuration
//	@Description	Delete a rarity configuration (admin only)
//	@Tags			rarity-configs
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Rarity Config ID"
//	@Success		204
//	@Failure		400
//	@Failure		401
//	@Failure		403
//	@Failure		404
//	@Failure		500
//	@Router			/api/v1/rarity-configs/{id} [delete]
func (c *RarityConfigController) DeleteRarityConfig(ec echo.Context) error {
	// TODO: Add admin authorization check
	idParam := ec.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Invalid rarity config ID"))
	}

	// TODO: Implement delete logic when manager has the method
	_ = id
	return ec.JSON(http.StatusNotImplemented, response.ToErrorResponse(http.StatusNotImplemented, "Not implemented yet"))
}
