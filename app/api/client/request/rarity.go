package request

import (
	"github.com/go-playground/validator/v10"
)

type CreateRarityConfigRequest struct {
	Code       string  `json:"code" validate:"required,min=1,max=10"`
	Label      string  `json:"label" validate:"required,min=1,max=50"`
	Rank       int16   `json:"rank" validate:"required,min=1"`
	ColorHex   *string `json:"color_hex" validate:"omitempty,hexcolor,len=7"`
	DropWeight int     `json:"drop_weight" validate:"required,min=0"`
}

type UpdateRarityConfigRequest struct {
	Code       *string `json:"code" validate:"omitempty,min=1,max=10"`
	Label      *string `json:"label" validate:"omitempty,min=1,max=50"`
	Rank       *int16  `json:"rank" validate:"omitempty,min=1"`
	ColorHex   *string `json:"color_hex" validate:"omitempty,hexcolor,len=7"`
	DropWeight *int    `json:"drop_weight" validate:"omitempty,min=0"`
}

func (r *CreateRarityConfigRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *UpdateRarityConfigRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
