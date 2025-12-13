package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/database/entity"
	"github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/internal/runtime"
)

type PostRepository interface {
	FindByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]entity.Post, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Post, error)
	Insert(ctx context.Context, post *entity.Post) (*entity.Post, error)
	Update(ctx context.Context, post *entity.Post) (*entity.Post, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

type DefaultPostRepository struct {
	res runtime.Resource
}

func NewPostRepository(res runtime.Resource) PostRepository {
	return &DefaultPostRepository{res: res}
}

func (r *DefaultPostRepository) FindByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]entity.Post, int64, error) {
	var posts []entity.Post

	offset := (page - 1) * pageSize

	count, err := r.res.DB.
		ReplicaNewSelect().
		Model((*entity.Post)(nil)).
		Where("user_id = ?", userID).
		Where("deleted_at IS NULL").
		Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	err = r.res.DB.
		ReplicaNewSelect().
		Model(&posts).
		Where("user_id = ?", userID).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Scan(ctx)
	if err != nil {
		return nil, 0, err
	}

	return posts, int64(count), nil
}

func (r *DefaultPostRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Post, error) {
	post := new(entity.Post)
	err := r.res.DB.
		ReplicaNewSelect().
		Model(post).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (r *DefaultPostRepository) Insert(ctx context.Context, post *entity.Post) (*entity.Post, error) {
	err := r.res.DB.
		NewInsert().
		Model(post).
		Returning("*").
		Scan(ctx, post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (r *DefaultPostRepository) Update(ctx context.Context, post *entity.Post) (*entity.Post, error) {
	var p entity.Post
	err := r.res.DB.
		NewUpdate().
		Model(post).
		WherePK().
		Where("deleted_at IS NULL").
		Returning("*").
		Scan(ctx, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *DefaultPostRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	_, err := r.res.DB.
		NewUpdate().
		Model((*entity.Post)(nil)).
		Set("deleted_at", "NOW()").
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}
