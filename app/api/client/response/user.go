package response

import (
	"time"

	"github.com/google/uuid"
)

type UserProfileResponse struct {
	ID             uuid.UUID         `json:"id"`
	Username       string            `json:"username"`
	Email          *string           `json:"email"`
	Firstname      string            `json:"firstname"`
	Lastname       string            `json:"lastname"`
	Bio            string            `json:"bio"`
	WebsiteURL     string            `json:"website_url"`
	Location       string            `json:"location"`
	Avatar         string            `json:"avatar"`
	Cover          string            `json:"cover"`
	FollowerCount  int               `json:"follower_count"`
	FollowingCount int               `json:"following_count"`
	PostCount      int               `json:"post_count"`
	EventsCount    int               `json:"events_count"`
	IsFollowing    bool              `json:"is_following"`
	IsOwnProfile   bool              `json:"is_own_profile"`
	SocialLinks    map[string]string `json:"social_links"`
	IsE2EEEnabled  bool              `json:"is_e2ee_enabled"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      *time.Time        `json:"updated_at"`
}

type UserEventResponse struct {
	ID               uuid.UUID `json:"id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Location         string    `json:"location"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	Status           string    `json:"status"`
	Avatar           string    `json:"avatar"`
	CurrentAttendees int       `json:"current_attendees"`
	MaxAttendees     int       `json:"max_attees"`
	IsPublic         bool      `json:"is_public"`
	Price            float64   `json:"price"`
	Currency         string    `json:"currency"`
	CreatedAt        time.Time `json:"created_at"`
}

type UserPostResponse struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Image         string    `json:"image"`
	LikesCount    int       `json:"likes_count"`
	CommentsCount int       `json:"comments_count"`
	CreatedAt     time.Time `json:"created_at"`
}

type UserListResponse struct {
	Users    []UserProfileResponse `json:"users"`
	Total    int64                 `json:"total"`
	Page     int                   `json:"page"`
	PageSize int                   `json:"page_size"`
}

type UserEventListResponse struct {
	Events   []UserEventResponse `json:"events"`
	Total    int64               `json:"total"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"page_size"`
}

type UserPostListResponse struct {
	Posts    []UserPostResponse `json:"posts"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
}
