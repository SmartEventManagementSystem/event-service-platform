package response

import (
	"time"

	"github.com/google/uuid"
)

type RarityConfigResponse struct {
	ID         uuid.UUID  `json:"id"`
	Code       string     `json:"code"`
	Label      string     `json:"label"`
	Rank       int16      `json:"rank"`
	ColorHex   *string    `json:"color_hex,omitempty"`
	DropWeight int        `json:"drop_weight"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
}
