package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TicketStatus string

const (
	TicketStatusAvailable TicketStatus = "available"
	TicketStatusReserved  TicketStatus = "reserved"
	TicketStatusSold      TicketStatus = "sold"
	TicketStatusCancelled TicketStatus = "cancelled"
	TicketStatusUsed      TicketStatus = "used"
)

type Ticket struct {
	bun.BaseModel `bun:"table:tickets,alias:t"`

	ID          uuid.UUID    `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	EventID     uuid.UUID    `bun:"event_id,notnull"`
	Name        string       `bun:"name,notnull"`
	Description string       `bun:"description"`
	Price       float64      `bun:"price,notnull"`
	Currency    string       `bun:"currency,default:'USD'"`
	Quantity    int          `bun:"quantity,notnull"`
	Sold        int          `bun:"sold,default:0"`
	Status      TicketStatus `bun:"status,default:'available'"`
	SaleStartAt *time.Time   `bun:"sale_start_at"`
	SaleEndAt   *time.Time   `bun:"sale_end_at"`
	MinPurchase int          `bun:"min_purchase,default:1"`
	MaxPurchase int          `bun:"max_purchase,default:10"`
	CreatedAt   time.Time    `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt   *time.Time   `bun:"updated_at"`

	// Relations
	Event     *Event          `bun:"rel:belongs-to,join:event_id=id"`
	Attendees []EventAttendee `bun:"rel:has-many,join:ticket_id=id"`
}

func (t Ticket) Alias() string {
	return "t"
}
