package entity

import (
    "time"

    "github.com/google/uuid"
    "github.com/uptrace/bun"
)

type Comment struct {
    bun.BaseModel `bun:"table:comments,alias:c"`

    ID         uuid.UUID  `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
    PostID     uuid.UUID  `bun:"post_id,notnull,type:uuid"`
    UserID     uuid.UUID  `bun:"user_id,notnull,type:uuid"`
    ParentID   *uuid.UUID `bun:"parent_id,type:uuid"`
    Content    string     `bun:"content,notnull"`
    ReplyCount int        `bun:"reply_count,default:0"`
    CreatedAt  time.Time  `bun:"created_at,notnull,default:current_timestamp"`
    UpdatedAt  *time.Time `bun:"updated_at"`
    DeletedAt  *time.Time `bun:"deleted_at,soft_delete"`

    User *User `bun:"rel:belongs-to,join:user_id=id"`
}

func (c Comment) Alias() string { return "c" }
