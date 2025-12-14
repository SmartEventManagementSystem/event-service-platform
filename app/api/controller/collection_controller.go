package controller

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"

	"backend/event-service-platform/app/api/client/request"
	"backend/event-service-platform/app/api/client/response"
	collectionconst "backend/event-service-platform/app/database/constant/collection"
	"backend/event-service-platform/app/database/entity"
	"backend/event-service-platform/app/internal/runtime"
	"backend/event-service-platform/app/manager"
	"backend/event-service-platform/app/pkg/jwt"
)

type CollectionController struct {
	res      runtime.Resource
	managers *manager.Managers
	jwt      jwt.Jwt
}

func NewCollectionController(managers *manager.Managers, res runtime.Resource) *CollectionController {
	return &CollectionController{res: res, managers: managers, jwt: jwt.NewJwt(res.Config.JwtConfig)}
}

// ListCollections godoc
//
//	@Summary		List collections
//	@Description	List all collections with optional filtering
//	@Tags			collections
//	@Accept			json
//	@Produce		json
//	@Param			type		query	string	false	"Filter by collection type"	Enums(THEME,LOCATION,COUNTRY,GLOBAL)
//	@Param			enabled	query	boolean	false	"Filter by enabled status"
//	@Param			limit		query	int		false	"Limit results"			default(50)
//	@Param			offset		query	int		false	"Offset results"			default(0)
//	@Success		200			{array}		response.CollectionWithItemsResponse
//	@Failure		400
//	@Failure		401
//	@Failure		500
//	@Router			/api/v1/collections [get]
func (c *CollectionController) ListCollections(ec echo.Context) error {
	var filter manager.ListCollectionsFilter

	// Parse query parameters
	if typeParam := ec.QueryParam("type"); typeParam != "" {
		filter.Types = []collectionconst.Type{collectionconst.Type(typeParam)}
	}

	if enabledParam := ec.QueryParam("enabled"); enabledParam != "" {
		enabled := enabledParam == "true"
		filter.IsEnabled = &enabled
	}

	filter.Limit = 50
	filter.Offset = 0

	if limit := ec.QueryParam("limit"); limit != "" {
		if l, err := parseIntParam(limit); err == nil && l > 0 {
			filter.Limit = l
		}
	}

	if offset := ec.QueryParam("offset"); offset != "" {
		if o, err := parseIntParam(offset); err == nil && o >= 0 {
			filter.Offset = o
		}
	}

	collections, err := c.managers.CollectionManager.ListCollections(ec.Request().Context(), filter)
	if err != nil {
		return ec.JSON(http.StatusInternalServerError, response.ToErrorResponse(http.StatusInternalServerError, "Internal server error"))
	}

	// Convert to response DTOs
	result := make([]response.CollectionWithItemsResponse, len(collections))
	for i, col := range collections {
		result[i] = response.CollectionWithItemsResponse{
			Collection: response.CollectionResponse{
				ID:             col.Collection.ID,
				Name:           col.Collection.Name,
				Description:    col.Collection.Description,
				Type:           col.Collection.Type,
				RewardAmount:   col.Collection.RewardAmount,
				RewardCurrency: col.Collection.RewardCurrency,
				IsEnabled:      col.Collection.IsEnabled,
				CreatedAt:      col.Collection.CreatedAt,
				UpdatedAt:      col.Collection.UpdatedAt,
			},
			Items: convertCollectionItems(col.Items),
		}
	}

	return ec.JSON(http.StatusOK, response.ToSuccessResponse(result))
}

// GetCollection godoc
//
//	@Summary		Get collection by ID
//	@Description	Get a specific collection with its items
//	@Tags			collections
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Collection ID"
//	@Success		200		{object}	response.CollectionWithItemsResponse
//	@Failure		400
//	@Failure		401
//	@Failure		404
//	@Failure		500
//	@Router			/api/v1/collections/{id} [get]
func (c *CollectionController) GetCollection(ec echo.Context) error {
	idParam := ec.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Invalid collection ID"))
	}

	filter := manager.ListCollectionsFilter{
		IDs:    []uuid.UUID{id},
		Limit:  1,
		Offset: 0,
	}

	collections, err := c.managers.CollectionManager.ListCollections(ec.Request().Context(), filter)
	if err != nil {
		return ec.JSON(http.StatusInternalServerError, response.ToErrorResponse(http.StatusInternalServerError, "Internal server error"))
	}

	if len(collections) == 0 {
		return ec.JSON(http.StatusNotFound, response.ToErrorResponse(http.StatusNotFound, "Collection not found"))
	}

	col := collections[0]
	result := response.CollectionWithItemsResponse{
		Collection: response.CollectionResponse{
			ID:             col.Collection.ID,
			Name:           col.Collection.Name,
			Description:    col.Collection.Description,
			Type:           col.Collection.Type,
			RewardAmount:   col.Collection.RewardAmount,
			RewardCurrency: col.Collection.RewardCurrency,
			IsEnabled:      col.Collection.IsEnabled,
			CreatedAt:      col.Collection.CreatedAt,
			UpdatedAt:      col.Collection.UpdatedAt,
		},
		Items: convertCollectionItems(col.Items),
	}

	return ec.JSON(http.StatusOK, response.ToSuccessResponse(result))
}

// CreateCollection godoc
//
//	@Summary		Create collection
//	@Description	Create a new collection (admin only)
//	@Tags			collections
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.CreateCollectionRequest	true	"Collection data"
//	@Success		201		{object}	response.CollectionResponse
//	@Failure		400
//	@Failure		401
//	@Failure		403
//	@Failure		500
//	@Router			/api/v1/collections [post]
func (c *CollectionController) CreateCollection(ec echo.Context) error {
	// TODO: Add admin authorization check
	var req request.CreateCollectionRequest
	if err := ec.Bind(&req); err != nil {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Invalid request body"))
	}

	if err := req.Validate(); err != nil {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Validation failed: "+err.Error()))
	}

	// TODO: Implement collection creation logic
	return ec.JSON(http.StatusNotImplemented, response.ToErrorResponse(http.StatusNotImplemented, "Not implemented yet"))
}

// UpdateCollection godoc
//
//	@Summary		Update collection
//	@Description	Update an existing collection (admin only)
//	@Tags			collections
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"Collection ID"
//	@Param			request	body		request.UpdateCollectionRequest	true	"Collection data"
//	@Success		200		{object}	response.CollectionResponse
//	@Failure		400
//	@Failure		401
//	@Failure		403
//	@Failure		404
//	@Failure		500
//	@Router			/api/v1/collections/{id} [put]
func (c *CollectionController) UpdateCollection(ec echo.Context) error {
	// TODO: Add admin authorization check
	idParam := ec.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Invalid collection ID"))
	}

	var req request.UpdateCollectionRequest
	if err := ec.Bind(&req); err != nil {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Invalid request body"))
	}

	if err := req.Validate(); err != nil {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Validation failed: "+err.Error()))
	}

	// TODO: Implement collection update logic
	_ = id
	return ec.JSON(http.StatusNotImplemented, response.ToErrorResponse(http.StatusNotImplemented, "Not implemented yet"))
}

// DeleteCollection godoc
//
//	@Summary		Delete collection
//	@Description	Delete a collection (admin only)
//	@Tags			collections
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Collection ID"
//	@Success		204
//	@Failure		400
//	@Failure		401
//	@Failure		403
//	@Failure		404
//	@Failure		500
//	@Router			/api/v1/collections/{id} [delete]
func (c *CollectionController) DeleteCollection(ec echo.Context) error {
	// TODO: Add admin authorization check
	idParam := ec.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Invalid collection ID"))
	}

	err = c.managers.CollectionManager.DeleteCollection(ec.Request().Context(), id)
	if err != nil {
		return ec.JSON(http.StatusInternalServerError, response.ToErrorResponse(http.StatusInternalServerError, "Failed to delete collection"))
	}

	return ec.NoContent(http.StatusNoContent)
}

func convertCollectionItems(items []entity.CollectionItem) []response.CollectionItemResponse {
	result := make([]response.CollectionItemResponse, len(items))
	for i, item := range items {
		result[i] = response.CollectionItemResponse{
			ID:           item.ID,
			CollectionID: item.CollectionID,
			ItemID:       item.ItemID,
			ImageURL:     item.ImageURL,
			CreatedAt:    item.CreatedAt,
			UpdatedAt:    item.UpdatedAt,
		}
	}
	return result
}

func parseIntParam(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}
