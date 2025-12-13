package request

import "github.com/google/uuid"

type CreateCommentRequest struct {
    Content  string     `json:"content" validate:"required"`
    ParentID *uuid.UUID `json:"parent_id"`
}
