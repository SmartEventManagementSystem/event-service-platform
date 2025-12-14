package manager

import (
	"context"
	"github.com/google/uuid"

	"backend/event-service-platform/app/api/client/request"
	"backend/event-service-platform/app/api/client/response"
	"backend/event-service-platform/app/database/entity"
	"backend/event-service-platform/app/database/repository"
)

type PostManager interface {
	CreatePost(ctx context.Context, req *request.CreatePostRequest, userID uuid.UUID) (*response.UserPostResponse, error)
	GetPost(ctx context.Context, id uuid.UUID) (*response.UserPostResponse, error)
	ListPostsByUser(ctx context.Context, userID uuid.UUID, page, pageSize int) (*response.UserPostListResponse, error)
	UpdatePost(ctx context.Context, id uuid.UUID, req *request.UpdatePostRequest) (*response.UserPostResponse, error)
	DeletePost(ctx context.Context, id uuid.UUID) error
}

type DefaultPostManager struct {
	postRepo repository.PostRepository
}

func NewPostManager(postRepo repository.PostRepository) PostManager {
	return &DefaultPostManager{postRepo: postRepo}
}

func (m *DefaultPostManager) CreatePost(ctx context.Context, req *request.CreatePostRequest, userID uuid.UUID) (*response.UserPostResponse, error) {
	p := &entity.Post{
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
		Image:   req.Image,
		Tags:    req.Tags,
		Status:  entity.PostStatusPublished,
	}

	created, err := m.postRepo.Insert(ctx, p)
	if err != nil {
		return nil, err
	}

	return &response.UserPostResponse{
		ID:            created.ID,
		Title:         created.Title,
		Content:       created.Content,
		Image:         created.Image,
		LikesCount:    created.LikesCount,
		CommentsCount: created.CommentsCount,
		CreatedAt:     created.CreatedAt,
	}, nil
}

func (m *DefaultPostManager) GetPost(ctx context.Context, id uuid.UUID) (*response.UserPostResponse, error) {
	p, err := m.postRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &response.UserPostResponse{
		ID:            p.ID,
		Title:         p.Title,
		Content:       p.Content,
		Image:         p.Image,
		LikesCount:    p.LikesCount,
		CommentsCount: p.CommentsCount,
		CreatedAt:     p.CreatedAt,
	}, nil
}

func (m *DefaultPostManager) ListPostsByUser(ctx context.Context, userID uuid.UUID, page, pageSize int) (*response.UserPostListResponse, error) {
	posts, total, err := m.postRepo.FindByUserID(ctx, userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	res := make([]response.UserPostResponse, len(posts))
	for i, p := range posts {
		res[i] = response.UserPostResponse{
			ID:            p.ID,
			Title:         p.Title,
			Content:       p.Content,
			Image:         p.Image,
			LikesCount:    p.LikesCount,
			CommentsCount: p.CommentsCount,
			CreatedAt:     p.CreatedAt,
		}
	}

	return &response.UserPostListResponse{Posts: res, Total: total, Page: page, PageSize: pageSize}, nil
}

func (m *DefaultPostManager) UpdatePost(ctx context.Context, id uuid.UUID, req *request.UpdatePostRequest) (*response.UserPostResponse, error) {
	post, err := m.postRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Title != nil {
		post.Title = *req.Title
	}
	if req.Content != nil {
		post.Content = *req.Content
	}
	if req.Image != nil {
		post.Image = *req.Image
	}
	updated, err := m.postRepo.Update(ctx, post)
	if err != nil {
		return nil, err
	}
	return &response.UserPostResponse{
		ID:            updated.ID,
		Title:         updated.Title,
		Content:       updated.Content,
		Image:         updated.Image,
		LikesCount:    updated.LikesCount,
		CommentsCount: updated.CommentsCount,
		CreatedAt:     updated.CreatedAt,
	}, nil
}

func (m *DefaultPostManager) DeletePost(ctx context.Context, id uuid.UUID) error {
	return m.postRepo.DeleteByID(ctx, id)
}
