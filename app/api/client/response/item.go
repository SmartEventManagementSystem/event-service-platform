package response

import (
	"time"

	"github.com/google/uuid"
)

type ItemResponse struct {
	ID             uuid.UUID  `json:"id"`
	Name           string     `json:"name"`
	Description    *string    `json:"description,omitempty"`
	RarityConfigID *uuid.UUID `json:"rarity_config_id,omitempty"`
	ImageURL       *string    `json:"image_url,omitempty"`
	CountryID      *uuid.UUID `json:"country_id,omitempty"`
	LocationID     *uuid.UUID `json:"location_id,omitempty"`
	CollectionID   *uuid.UUID `json:"collection_id,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}

type ItemWithRarityResponse struct {
	ItemResponse
	Rarity *RarityConfigResponse `json:"rarity,omitempty"`
}
