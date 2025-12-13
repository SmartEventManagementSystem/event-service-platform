package controller_test

import (
    "bytes"
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/google/uuid"
    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/assert"

    ctrl "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/api/controller"
    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/api/client/request"
    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/api/client/response"
    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/internal/runtime"
    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/manager"
)

type fakePostManager struct{}

func (f *fakePostManager) CreatePost(ctx context.Context, req *request.CreatePostRequest, userID uuid.UUID) (*response.UserPostResponse, error) {
    return &response.UserPostResponse{ID: uuid.New(), Title: req.Title, Content: req.Content}, nil
}
func (f *fakePostManager) GetPost(ctx context.Context, id uuid.UUID) (*response.UserPostResponse, error) { return &response.UserPostResponse{ID: id, Title: "t", Content: "c"}, nil }
func (f *fakePostManager) ListPostsByUser(ctx context.Context, userID uuid.UUID, page, pageSize int) (*response.UserPostListResponse, error) {
    return &response.UserPostListResponse{Posts: []response.UserPostResponse{{ID: uuid.New(), Title: "t"}}}, nil
}
func (f *fakePostManager) UpdatePost(ctx context.Context, id uuid.UUID, req *request.UpdatePostRequest) (*response.UserPostResponse, error) { return &response.UserPostResponse{ID: id, Title: "t"}, nil }
func (f *fakePostManager) DeletePost(ctx context.Context, id uuid.UUID) error { return nil }

func makeController() *ctrl.PostController {
    managers := &manager.Managers{PostManager: &fakePostManager{}}
    res := runtime.Resource{}
    return ctrl.NewPostController(managers, res)
}

func TestCreatePost_Unauthorized(t *testing.T) {
    e := echo.New()
    reqBody := map[string]string{"title": "hi", "content": "hello"}
    b, _ := json.Marshal(reqBody)
    req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    pc := makeController()
    if assert.NoError(t, pc.CreatePost(c)) == false {
        return
    }
    assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestCreatePost_Success(t *testing.T) {
    e := echo.New()
    reqBody := map[string]string{"title": "hi", "content": "hello"}
    b, _ := json.Marshal(reqBody)
    req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    // set user id in context
    c.Set("user_id", uuid.New().String())

    pc := makeController()
    err := pc.CreatePost(c)
    assert.NoError(t, err)
    assert.Equal(t, http.StatusCreated, rec.Code)
    var resp response.UserPostResponse
    _ = json.Unmarshal(rec.Body.Bytes(), &resp)
    assert.NotEmpty(t, resp.ID)
}
