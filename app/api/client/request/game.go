package request

import (
	"backend/event-service-platform/app/database/constant/currency"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type SpinRequest struct {
	CollectionID *uuid.UUID `json:"collection_id" validate:"omitempty,uuid4"`
}

type ClaimDailyRewardRequest struct {
	Currency currency.Currency `json:"currency" validate:"required,oneof=COIN SPIN"`
}

func (r *SpinRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *ClaimDailyRewardRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
