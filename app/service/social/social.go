package social

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

// SocialClient defines operations provided by the external social-service-platform
type SocialClient interface {
	FollowUser(ctx context.Context, username string, followerID uuid.UUID) error
	UnfollowUser(ctx context.Context, username string, followerID uuid.UUID) error
	// other social operations can be added here (react, like, etc.)
}

// HTTPClient is a simple REST client implementation of SocialClient
type HTTPClient struct {
	client  *resty.Client
	baseURL string
}

func NewHTTPClient(client *resty.Client, baseURL string) *HTTPClient {
	return &HTTPClient{client: client, baseURL: baseURL}
}

func (h *HTTPClient) FollowUser(ctx context.Context, username string, followerID uuid.UUID) error {
	url := fmt.Sprintf("%s/api/v1/social/users/%s/follow", h.baseURL, username)
	resp, err := h.client.R().SetContext(ctx).SetBody(map[string]string{"user_id": followerID.String()}).Post(url)
	if err != nil {
		return err
	}
	if resp.StatusCode() >= 400 {
		return fmt.Errorf("social service error: %s", resp.Status())
	}
	return nil
}

func (h *HTTPClient) UnfollowUser(ctx context.Context, username string, followerID uuid.UUID) error {
	url := fmt.Sprintf("%s/api/v1/social/users/%s/unfollow", h.baseURL, username)
	resp, err := h.client.R().SetContext(ctx).SetBody(map[string]string{"user_id": followerID.String()}).Post(url)
	if err != nil {
		return err
	}
	if resp.StatusCode() >= 400 {
		return fmt.Errorf("social service error: %s", resp.Status())
	}
	return nil
}
