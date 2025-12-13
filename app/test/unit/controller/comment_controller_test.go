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

type fakeCommentManager struct{}

func (f *fakeCommentManager) CreateComment(ctx context.Context, req *request.CreateCommentRequest, userID uuid.UUID, postID uuid.UUID) (*response.CommentResponse, error) {
    return &response.CommentResponse{ID: uuid.New(), Content: req.Content, UserID: userID, PostID: postID}, nil
}
func (f *fakeCommentManager) ListComments(ctx context.Context, postID uuid.UUID, page, pageSize int) (*response.CommentListResponse, error) {
    return &response.CommentListResponse{Comments: []response.CommentResponse{{ID: uuid.New(), Content: "ok", UserID: uuid.New(), PostID: postID}}}, nil
}
func (f *fakeCommentManager) DeleteComment(ctx context.Context, id uuid.UUID) error { return nil }

func makeCommentController() *ctrl.CommentController {
    managers := &manager.Managers{CommentManager: &fakeCommentManager{}}
    res := runtime.Resource{}
    return ctrl.NewCommentController(managers, res)
}

func TestCreateComment_Unauthorized(t *testing.T) {
    e := echo.New()
    reqBody := map[string]string{"content": "hello"}
    b, _ := json.Marshal(reqBody)
    req := httptest.NewRequest(http.MethodPost, "/api/v1/posts/111/comments", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    cc := makeCommentController()
    err := cc.CreateComment(c)
    assert.NoError(t, err)
    assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestCreateComment_Success(t *testing.T) {
    e := echo.New()
    reqBody := map[string]string{"content": "hello"}
    b, _ := json.Marshal(reqBody)
    postID := uuid.New()
    req := httptest.NewRequest(http.MethodPost, "/api/v1/posts/"+postID.String()+"/comments", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.Set("user_id", uuid.New().String())

    cc := makeCommentController()
    err := cc.CreateComment(c)
    assert.NoError(t, err)
    assert.Equal(t, http.StatusCreated, rec.Code)
}
