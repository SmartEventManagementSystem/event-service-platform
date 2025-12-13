package response

import (
    "time"

    "github.com/google/uuid"
)

type CommentResponse struct {
    ID        uuid.UUID  `json:"id"`
    PostID    uuid.UUID  `json:"post_id"`
    UserID    uuid.UUID  `json:"user_id"`
    Content   string     `json:"content"`
    CreatedAt time.Time  `json:"created_at"`
}

type CommentListResponse struct {
    Comments []CommentResponse `json:"comments"`
    Total    int64             `json:"total"`
    Page     int               `json:"page"`
    PageSize int               `json:"page_size"`
}
