package controller

import (
	"backend/event-service-platform/app/api/client/request"
	"backend/event-service-platform/app/api/client/response"
	"backend/event-service-platform/app/database/entity"
	"backend/event-service-platform/app/internal/runtime"
	"backend/event-service-platform/app/manager"
	"backend/event-service-platform/app/pkg/jwt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ItemController struct {
	res      runtime.Resource
	managers *manager.Managers
	jwt      jwt.Jwt
}

func NewItemController(managers *manager.Managers, res runtime.Resource) *ItemController {
	return &ItemController{res: res, managers: managers, jwt: jwt.NewJwt(res.Config.JwtConfig)}
}

// ListItems godoc
//
//	@Summary		List items
//	@Description	List all items with optional filtering
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			collection_id	query	string	false	"Filter by collection ID"
//	@Param			rarity_id		query	string	false	"Filter by rarity config ID"
//	@Param			limit			query	int		false	"Limit results"	default(50)
//	@Param			offset			query	int		false	"Offset results"	default(0)
//	@Success		200				{array}		response.ItemWithRarityResponse
//	@Failure		400
//	@Failure		401
//	@Failure		500
//	@Router			/api/v1/items [get]
func (c *ItemController) ListItems(ec echo.Context) error {
	var collectionID *uuid.UUID
	var rarityConfigID *uuid.UUID

	if collectionIDParam := ec.QueryParam("collection_id"); collectionIDParam != "" {
		if id, err := uuid.Parse(collectionIDParam); err == nil {
			collectionID = &id
		}
	}

	if rarityIDParam := ec.QueryParam("rarity_id"); rarityIDParam != "" {
		if id, err := uuid.Parse(rarityIDParam); err == nil {
			rarityConfigID = &id
		}
	}

	// TODO: Use collectionID and rarityConfigID for filtering when implemented
	_ = collectionID
	_ = rarityConfigID

	// TODO: Implement item listing logic with filters
	return ec.JSON(http.StatusNotImplemented, response.ToErrorResponse(http.StatusNotImplemented, "Not implemented yet"))
}

// GetItem godoc
//
//	@Summary		Get item by ID
//	@Description	Get a specific item with rarity information
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Item ID"
//	@Success		200	{object}	response.ItemWithRarityResponse
//	@Failure		400
//	@Failure		401
//	@Failure		404
//	@Failure		500
//	@Router			/api/v1/items/{id} [get]
func (c *ItemController) GetItem(ec echo.Context) error {
	idParam := ec.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Invalid item ID"))
	}

	item, err := c.managers.ItemManager.GetItem(ec.Request().Context(), id)
	if err != nil {
		return ec.JSON(http.StatusInternalServerError, response.ToErrorResponse(http.StatusInternalServerError, "Internal server error"))
	}

	if item == nil {
		return ec.JSON(http.StatusNotFound, response.ToErrorResponse(http.StatusNotFound, "Item not found"))
	}

	// TODO: Fetch rarity config and convert to response
	itemResponse := response.ItemResponse{
		ID:             item.ID,
		Name:           item.Name,
		Description:    item.Description,
		RarityConfigID: item.RarityConfigID,
		ImageURL:       item.ImageURL,
		CountryID:      item.CountryID,
		LocationID:     item.LocationID,
		CollectionID:   item.CollectionID,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
	}

	return ec.JSON(http.StatusOK, response.ToSuccessResponse(response.ItemWithRarityResponse{
		ItemResponse: itemResponse,
	}))
}

// CreateItem godoc
//
//	@Summary		Create item
//	@Description	Create a new item (admin only)
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.CreateItemRequest	true	"Item data"
//	@Success		201		{object}	response.ItemResponse
//	@Failure		400
//	@Failure		401
//	@Failure		403
//	@Failure		500
//	@Router			/api/v1/items [post]
func (c *ItemController) CreateItem(ec echo.Context) error {
	// TODO: Add admin authorization check
	var req request.CreateItemRequest
	if err := ec.Bind(&req); err != nil {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Invalid request body"))
	}

	if err := req.Validate(); err != nil {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Validation failed: "+err.Error()))
	}

	item := &entity.Item{
		Name:           req.Name,
		Description:    req.Description,
		RarityConfigID: &req.RarityConfigID,
		ImageURL:       req.ImageURL,
		CountryID:      req.CountryID,
		LocationID:     req.LocationID,
		CollectionID:   req.CollectionID,
	}

	createdItem, err := c.managers.ItemManager.CreateItem(ec.Request().Context(), item)
	if err != nil {
		return ec.JSON(http.StatusInternalServerError, response.ToErrorResponse(http.StatusInternalServerError, "Failed to create item"))
	}

	itemResponse := response.ItemResponse{
		ID:             createdItem.ID,
		Name:           createdItem.Name,
		Description:    createdItem.Description,
		RarityConfigID: createdItem.RarityConfigID,
		ImageURL:       createdItem.ImageURL,
		CountryID:      createdItem.CountryID,
		LocationID:     createdItem.LocationID,
		CollectionID:   createdItem.CollectionID,
		CreatedAt:      createdItem.CreatedAt,
		UpdatedAt:      createdItem.UpdatedAt,
	}

	return ec.JSON(http.StatusCreated, response.ToSuccessResponse(itemResponse))
}

// UpdateItem godoc
//
//	@Summary		Update item
//	@Description	Update an existing item (admin only)
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"Item ID"
//	@Param			request	body		request.UpdateItemRequest	true	"Item data"
//	@Success		200		{object}	response.ItemResponse
//	@Failure		400
//	@Failure		401
//	@Failure		403
//	@Failure		404
//	@Failure		500
//	@Router			/api/v1/items/{id} [put]
func (c *ItemController) UpdateItem(ec echo.Context) error {
	// TODO: Add admin authorization check
	idParam := ec.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Invalid item ID"))
	}

	var req request.UpdateItemRequest
	if err := ec.Bind(&req); err != nil {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Invalid request body"))
	}

	if err := req.Validate(); err != nil {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Validation failed: "+err.Error()))
	}

	// Get existing item
	item, err := c.managers.ItemManager.GetItem(ec.Request().Context(), id)
	if err != nil {
		return ec.JSON(http.StatusInternalServerError, response.ToErrorResponse(http.StatusInternalServerError, "Internal server error"))
	}

	if item == nil {
		return ec.JSON(http.StatusNotFound, response.ToErrorResponse(http.StatusNotFound, "Item not found"))
	}

	// Update fields
	if req.Name != nil {
		item.Name = *req.Name
	}
	if req.Description != nil {
		item.Description = req.Description
	}
	if req.RarityConfigID != nil {
		item.RarityConfigID = req.RarityConfigID
	}
	if req.ImageURL != nil {
		item.ImageURL = req.ImageURL
	}
	if req.CountryID != nil {
		item.CountryID = req.CountryID
	}
	if req.LocationID != nil {
		item.LocationID = req.LocationID
	}
	if req.CollectionID != nil {
		item.CollectionID = req.CollectionID
	}

	updatedItem, err := c.managers.ItemManager.UpdateItem(ec.Request().Context(), item)
	if err != nil {
		return ec.JSON(http.StatusInternalServerError, response.ToErrorResponse(http.StatusInternalServerError, "Failed to update item"))
	}

	itemResponse := response.ItemResponse{
		ID:             updatedItem.ID,
		Name:           updatedItem.Name,
		Description:    updatedItem.Description,
		RarityConfigID: updatedItem.RarityConfigID,
		ImageURL:       updatedItem.ImageURL,
		CountryID:      updatedItem.CountryID,
		LocationID:     updatedItem.LocationID,
		CollectionID:   updatedItem.CollectionID,
		CreatedAt:      updatedItem.CreatedAt,
		UpdatedAt:      updatedItem.UpdatedAt,
	}

	return ec.JSON(http.StatusOK, response.ToSuccessResponse(itemResponse))
}

// DeleteItem godoc
//
//	@Summary		Delete item
//	@Description	Delete an item (admin only)
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Item ID"
//	@Success		204
//	@Failure		400
//	@Failure		401
//	@Failure		403
//	@Failure		404
//	@Failure		500
//	@Router			/api/v1/items/{id} [delete]
func (c *ItemController) DeleteItem(ec echo.Context) error {
	// TODO: Add admin authorization check
	idParam := ec.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ec.JSON(http.StatusBadRequest, response.ToErrorResponse(http.StatusBadRequest, "Invalid item ID"))
	}

	err = c.managers.ItemManager.DeleteItems(ec.Request().Context(), []uuid.UUID{id})
	if err != nil {
		return ec.JSON(http.StatusInternalServerError, response.ToErrorResponse(http.StatusInternalServerError, "Failed to delete item"))
	}

	return ec.NoContent(http.StatusNoContent)
}
