package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/database/constant/role"
	"github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/database/constant/user"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID            uuid.UUID   `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	Username      string      `bun:"username,notnull,unique"`
	Email         *string     `bun:"email,unique"`
	PhoneNumber   *string     `bun:"phone_number,unique"`
	Password      string      `bun:"password,notnull"`
	Status        user.Status `bun:"status,notnull,default:'UNVERIFIED'"`
	Role          role.Role   `bun:"role,notnull,default:'USER'"`
	EmailVerified bool        `bun:"email_verified,default:false"`
	PhoneVerified bool        `bun:"phone_verified,default:false"`
	LastLoginAt   *time.Time  `bun:"last_login_at"`
	CreatedAt     time.Time   `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt     *time.Time  `bun:"updated_at"`
	DeletedAt     *time.Time  `bun:"deleted_at,soft_delete"`
	DeactivatedAt *time.Time  `bun:"deactivated_at"`
	// Events-specific fields
	Firstname      string `bun:"firstname"`
	Lastname       string `bun:"lastname"`
	Cover          string `bun:"cover"`
	Avatar         string `bun:"avatar"`
	Bio            string `bun:"bio"`
	WebsiteURL     string `bun:"website_url"`
	FollowerCount  int    `bun:"follower_count,default:0"`
	FollowingCount int    `bun:"following_count,default:0"`
	PostCount      int    `bun:"post_count,default:0"`
	EventsCount    int    `bun:"events_count,default:0"`
	// E2EE fields for chat encryption
	EncryptionPasscode string `bun:"encryption_passcode"`
	IsE2EEEnabled      bool   `bun:"is_e2ee_enabled,default:false"`
}

func (u User) Alias() string {
	return "u"
}
