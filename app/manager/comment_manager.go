package manager

import (
    "context"

    "github.com/google/uuid"

    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/api/client/request"
    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/api/client/response"
    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/database/entity"
    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/database/repository"
)

type CommentManager interface {
    CreateComment(ctx context.Context, req *request.CreateCommentRequest, userID uuid.UUID, postID uuid.UUID) (*response.CommentResponse, error)
    ListComments(ctx context.Context, postID uuid.UUID, page, pageSize int) (*response.CommentListResponse, error)
    DeleteComment(ctx context.Context, id uuid.UUID) error
}

type DefaultCommentManager struct{
    commentRepo repository.CommentRepository
}

func NewCommentManager(commentRepo repository.CommentRepository) CommentManager {
    return &DefaultCommentManager{commentRepo: commentRepo}
}

func (m *DefaultCommentManager) CreateComment(ctx context.Context, req *request.CreateCommentRequest, userID uuid.UUID, postID uuid.UUID) (*response.CommentResponse, error) {
    c := &entity.Comment{
        PostID:  postID,
        UserID:  userID,
        Content: req.Content,
    }
    if req.ParentID != nil {
        c.ParentID = req.ParentID
    }
    created, err := m.commentRepo.Insert(ctx, c)
    if err != nil { return nil, err }
    return &response.CommentResponse{ID: created.ID, Content: created.Content, UserID: created.UserID, PostID: created.PostID, CreatedAt: created.CreatedAt}, nil
}

func (m *DefaultCommentManager) ListComments(ctx context.Context, postID uuid.UUID, page, pageSize int) (*response.CommentListResponse, error) {
    comments, total, err := m.commentRepo.FindByPostID(ctx, postID, page, pageSize)
    if err != nil { return nil, err }
    out := make([]response.CommentResponse, len(comments))
    for i, c := range comments {
        out[i] = response.CommentResponse{ID: c.ID, Content: c.Content, UserID: c.UserID, PostID: c.PostID, CreatedAt: c.CreatedAt}
    }
    return &response.CommentListResponse{Comments: out, Total: total, Page: page, PageSize: pageSize}, nil
}

func (m *DefaultCommentManager) DeleteComment(ctx context.Context, id uuid.UUID) error {
    return m.commentRepo.DeleteByID(ctx, id)
}
