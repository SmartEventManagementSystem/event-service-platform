package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPublished PostStatus = "published"
	PostStatusArchived  PostStatus = "archived"
)

type Post struct {
	bun.BaseModel `bun:"table:posts,alias:p"`

	ID            uuid.UUID  `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	UserID        uuid.UUID  `bun:"user_id,notnull"`
	Title         string     `bun:"title,notnull"`
	Content       string     `bun:"content,notnull"`
	Image         string     `bun:"image"`
	LikesCount    int        `bun:"likes_count,default:0"`
	CommentsCount int        `bun:"comments_count,default:0"`
	Status        PostStatus `bun:"status,default:'draft'"`
	Tags          []string   `bun:"tags,array"`
	CreatedAt     time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt     *time.Time `bun:"updated_at"`
	DeletedAt     *time.Time `bun:"deleted_at,soft_delete"`

	// Relations
	User *User `bun:"rel:belongs-to,join:user_id=id"`
}

func (p Post) Alias() string {
	return "p"
}
