package request

type GetUserProfileRequest struct {
	Username string `param:"username" validate:"required,min=3,max=50"`
}

type GetUserEventsRequest struct {
	Username string `param:"username" validate:"required,min=3,max=50"`
	Page     int    `query:"page" validate:"min=1"`
	PageSize int    `query:"page_size" validate:"min=1,max=100"`
}

type GetUserPostsRequest struct {
	Username string `param:"username" validate:"required,min=3,max=50"`
	Page     int    `query:"page" validate:"min=1"`
	PageSize int    `query:"page_size" validate:"min=1,max=100"`
}

type FollowUserRequest struct {
	Username string `param:"username" validate:"required,min=3,max=50"`
}

type UnfollowUserRequest struct {
	Username string `param:"username" validate:"required,min=3,max=50"`
}

type UpdateUserProfileRequest struct {
	Firstname   string            `json:"firstname" validate:"required,min=1,max=100"`
	Lastname    string            `json:"lastname" validate:"required,min=1,max=100"`
	Bio         string            `json:"bio" validate:"max=500"`
	WebsiteURL  string            `json:"website_url" validate:"omitempty,url"`
	Location    string            `json:"location" validate:"max=100"`
	Avatar      string            `json:"avatar"`
	Cover       string            `json:"cover"`
	SocialLinks map[string]string `json:"social_links"`
}

type UpdateUserSettingsRequest struct {
	IsE2EEEnabled      bool   `json:"is_e2ee_enabled"`
	EncryptionPasscode string `json:"encryption_passcode" validate:"omitempty,min=6,max=50"`
}
