package controller_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"backend/event-service-platform/app/api/client/request"
	ctrl "backend/event-service-platform/app/api/controller"
	"backend/event-service-platform/app/internal/runtime"
	"backend/event-service-platform/app/manager"
	"github.com/google/uuid"
)

type fakeSocialClient struct {
	followErr   error
	unfollowErr error
}

func (f *fakeSocialClient) FollowUser(ctx context.Context, username string, followerID uuid.UUID) error {
	return f.followErr
}
func (f *fakeSocialClient) UnfollowUser(ctx context.Context, username string, followerID uuid.UUID) error {
	return f.unfollowErr
}

func makeUserControllerWithSocial(sc interface{}) *ctrl.UserController {
	mgrs := &manager.Managers{SocialClient: sc}
	res := runtime.Resource{}
	return ctrl.NewUserController(mgrs, res)
}

func TestFollowUser_DelegatesToSocialClient(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/alice/follow", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// set JWT claims fallback by setting user_id in context used by controller
	c.Set("user_id", uuid.New().String())

	sc := &fakeSocialClient{}
	uc := makeUserControllerWithSocial(sc)

	// Prepare request param (username path param)
	c.SetParamNames("username")
	c.SetParamValues("alice")

	err := uc.FollowUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestUnfollowUser_DelegatesToSocialClient(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/alice/unfollow", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("user_id", uuid.New().String())

	sc := &fakeSocialClient{}
	uc := makeUserControllerWithSocial(sc)

	c.SetParamNames("username")
	c.SetParamValues("alice")

	err := uc.UnfollowUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
