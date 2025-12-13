package response

import (
	"time"

	"github.com/google/uuid"
)

type UserInventoryResponse struct {
	ID        uuid.UUID              `json:"id"`
	UserID    uuid.UUID              `json:"user_id"`
	ItemID    uuid.UUID              `json:"item_id"`
	Quantity  int                    `json:"quantity"`
	IsNew     bool                   `json:"is_new"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt *time.Time             `json:"updated_at,omitempty"`
	Item      ItemWithRarityResponse `json:"item"`
}

type SpinResultResponse struct {
	Success      bool                    `json:"success"`
	Item         *ItemWithRarityResponse `json:"item,omitempty"`
	Message      string                  `json:"message"`
	SpinCost     int64                   `json:"spin_cost"`
	NewBalance   int64                   `json:"new_balance,omitempty"`
	CollectionID *uuid.UUID              `json:"collection_id,omitempty"`
}

type GameStatsResponse struct {
	TotalCollections int `json:"total_collections"`
	TotalItems       int `json:"total_items"`
	UniqueItems      int `json:"unique_items"`
	TotalSpins       int `json:"total_spins"`
}
