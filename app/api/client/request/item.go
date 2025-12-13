package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CreateItemRequest struct {
	Name           string     `json:"name" validate:"required,min=1,max=100"`
	Description    *string    `json:"description" validate:"max=500"`
	RarityConfigID uuid.UUID  `json:"rarity_config_id" validate:"required"`
	ImageURL       *string    `json:"image_url" validate:"omitempty,url"`
	CountryID      *uuid.UUID `json:"country_id"`
	LocationID     *uuid.UUID `json:"location_id"`
	CollectionID   *uuid.UUID `json:"collection_id"`
}

type UpdateItemRequest struct {
	Name           *string    `json:"name" validate:"omitempty,min=1,max=100"`
	Description    *string    `json:"description" validate:"omitempty,max=500"`
	RarityConfigID *uuid.UUID `json:"rarity_config_id"`
	ImageURL       *string    `json:"image_url" validate:"omitempty,url"`
	CountryID      *uuid.UUID `json:"country_id"`
	LocationID     *uuid.UUID `json:"location_id"`
	CollectionID   *uuid.UUID `json:"collection_id"`
}

func (r *CreateItemRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *UpdateItemRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
