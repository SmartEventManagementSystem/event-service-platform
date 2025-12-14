package response

import (
	collectionconst "backend/event-service-platform/app/database/constant/collection"
	"backend/event-service-platform/app/database/constant/currency"
	"time"

	"github.com/google/uuid"
)

type CollectionResponse struct {
	ID             uuid.UUID            `json:"id"`
	Name           string               `json:"name"`
	Description    *string              `json:"description,omitempty"`
	Type           collectionconst.Type `json:"type"`
	RewardAmount   int64                `json:"reward_amount"`
	RewardCurrency currency.Currency    `json:"reward_currency"`
	IsEnabled      bool                 `json:"is_enabled"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      *time.Time           `json:"updated_at,omitempty"`
}

type CollectionWithItemsResponse struct {
	Collection CollectionResponse       `json:"collection"`
	Items      []CollectionItemResponse `json:"items"`
}

type CollectionItemResponse struct {
	ID           uuid.UUID  `json:"id"`
	CollectionID uuid.UUID  `json:"collection_id"`
	ItemID       uuid.UUID  `json:"item_id"`
	ImageURL     *string    `json:"image_url,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}
