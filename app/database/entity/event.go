package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type EventStatus string

const (
	EventStatusDraft     EventStatus = "draft"
	EventStatusPublished EventStatus = "published"
	EventStatusCancelled EventStatus = "cancelled"
	EventStatusCompleted EventStatus = "completed"
	EventStatusOngoing   EventStatus = "ongoing"
)

type Event struct {
	bun.BaseModel `bun:"table:events,alias:e"`

	ID               uuid.UUID   `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	WorkspaceID      uuid.UUID   `bun:"workspace_id,notnull"`
	Title            string      `bun:"title,notnull"`
	Description      string      `bun:"description"`
	Location         string      `bun:"location,notnull"`
	StartTime        time.Time   `bun:"start_time,notnull"`
	EndTime          time.Time   `bun:"end_time,notnull"`
	Status           EventStatus `bun:"status,default:'draft'"`
	Avatar           string      `bun:"avatar"`
	Cover            string      `bun:"cover"`
	MaxAttendees     int         `bun:"max_attendees"`
	CurrentAttendees int         `bun:"current_attendees,default:0"`
	IsPublic         bool        `bun:"is_public,default:true"`
	Price            float64     `bun:"price,default:0"`
	Currency         string      `bun:"currency,default:'USD'"`
	Tags             []string    `bun:"tags,array"`
	CreatedAt        time.Time   `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt        *time.Time  `bun:"updated_at"`
	DeletedAt        *time.Time  `bun:"deleted_at,soft_delete"`

	// Relations
	CreatorID uuid.UUID       `bun:"creator_id,notnull"`
	Creator   *User           `bun:"rel:belongs-to,join:creator_id=id"`
	Attendees []EventAttendee `bun:"rel:has-many,join:event_id=id"`
	Tickets   []Ticket        `bun:"rel:has-many,join:event_id=id"`
}

func (e Event) Alias() string {
	return "e"
}
