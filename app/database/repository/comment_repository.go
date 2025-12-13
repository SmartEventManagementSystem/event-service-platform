package repository

import (
    "context"

    "github.com/google/uuid"

    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/database/entity"
    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/internal/runtime"
)

type CommentRepository interface {
    FindByPostID(ctx context.Context, postID uuid.UUID, page, pageSize int) ([]entity.Comment, int64, error)
    Insert(ctx context.Context, comment *entity.Comment) (*entity.Comment, error)
    DeleteByID(ctx context.Context, id uuid.UUID) error
    CountByPostID(ctx context.Context, postID uuid.UUID) (int64, error)
}

type DefaultCommentRepository struct{
    res runtime.Resource
}

func NewCommentRepository(res runtime.Resource) CommentRepository {
    return &DefaultCommentRepository{res: res}
}

func (r *DefaultCommentRepository) FindByPostID(ctx context.Context, postID uuid.UUID, page, pageSize int) ([]entity.Comment, int64, error) {
    var comments []entity.Comment
    offset := (page - 1) * pageSize

    count, err := r.res.DB.
        ReplicaNewSelect().
        Model((*entity.Comment)(nil)).
        Where("post_id = ?", postID).
        Where("parent_id IS NULL").
        Count(ctx)
    if err != nil { return nil, 0, err }

    err = r.res.DB.
        ReplicaNewSelect().
        Model(&comments).
        Where("post_id = ?", postID).
        Where("parent_id IS NULL").
        Order("created_at DESC").
        Limit(pageSize).
        Offset(offset).
        Scan(ctx)
    if err != nil { return nil, 0, err }
    return comments, int64(count), nil
}

func (r *DefaultCommentRepository) Insert(ctx context.Context, comment *entity.Comment) (*entity.Comment, error) {
    if err := r.res.DB.NewInsert().Model(comment).Returning("*").Scan(ctx, comment); err != nil { return nil, err }
    return comment, nil
}

func (r *DefaultCommentRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
    _, err := r.res.DB.NewUpdate().Model((*entity.Comment)(nil)).Set("deleted_at", "NOW()").Where("id = ?", id).Where("deleted_at IS NULL").Exec(ctx)
    return err
}

func (r *DefaultCommentRepository) CountByPostID(ctx context.Context, postID uuid.UUID) (int64, error) {
    count, err := r.res.DB.ReplicaNewSelect().Model((*entity.Comment)(nil)).Where("post_id = ?", postID).Count(ctx)
    return int64(count), err
}
