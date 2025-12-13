package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type EventAttendeeStatus string

const (
	AttendeeStatusRegistered EventAttendeeStatus = "registered"
	AttendeeStatusCheckedIn  EventAttendeeStatus = "checked_in"
	AttendeeStatusCancelled  EventAttendeeStatus = "cancelled"
	AttendeeStatusWaitlisted EventAttendeeStatus = "waitlisted"
)

type EventAttendee struct {
	bun.BaseModel `bun:"table:event_attendees,alias:ea"`

	ID        uuid.UUID           `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	EventID   uuid.UUID           `bun:"event_id,notnull"`
	UserID    uuid.UUID           `bun:"user_id,notnull"`
	Status    EventAttendeeStatus `bun:"status,default:'registered'"`
	CheckInAt *time.Time          `bun:"check_in_at"`
	TicketID  *uuid.UUID          `bun:"ticket_id"`
	Notes     string              `bun:"notes"`
	CreatedAt time.Time           `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt *time.Time          `bun:"updated_at"`

	// Relations
	Event  *Event  `bun:"rel:belongs-to,join:event_id=id"`
	User   *User   `bun:"rel:belongs-to,join:user_id=id"`
	Ticket *Ticket `bun:"rel:belongs-to,join:ticket_id=id"`
}

func (ea EventAttendee) Alias() string {
	return "ea"
}
