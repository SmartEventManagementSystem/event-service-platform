package request

import (
	collectionconst "backend/event-service-platform/app/database/constant/collection"
	"backend/event-service-platform/app/database/constant/currency"
	"github.com/go-playground/validator/v10"
)

type CreateCollectionRequest struct {
	Name           string               `json:"name" validate:"required,min=1,max=100"`
	Description    *string              `json:"description" validate:"max=500"`
	Type           collectionconst.Type `json:"type" validate:"required,oneof=THEME LOCATION COUNTRY GLOBAL"`
	RewardAmount   int64                `json:"reward_amount" validate:"min=0"`
	RewardCurrency currency.Currency    `json:"reward_currency" validate:"required,oneof=COIN SPIN"`
	IsEnabled      *bool                `json:"is_enabled"`
}

type UpdateCollectionRequest struct {
	Name           *string               `json:"name" validate:"omitempty,min=1,max=100"`
	Description    *string               `json:"description" validate:"omitempty,max=500"`
	Type           *collectionconst.Type `json:"type" validate:"omitempty,oneof=THEME LOCATION COUNTRY GLOBAL"`
	RewardAmount   *int64                `json:"reward_amount" validate:"omitempty,min=0"`
	RewardCurrency *currency.Currency    `json:"reward_currency" validate:"omitempty,oneof=COIN SPIN"`
	IsEnabled      *bool                 `json:"is_enabled"`
}

func (r *CreateCollectionRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *UpdateCollectionRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
