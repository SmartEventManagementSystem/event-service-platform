package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/internal/runtime"
)

type UserFollowRepository interface {
	Create(ctx context.Context, followerID, followedID uuid.UUID) error
	Delete(ctx context.Context, followerID, followedID uuid.UUID) error
	IsFollowing(ctx context.Context, followerID, followedID uuid.UUID) (bool, error)
	GetFollowers(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]uuid.UUID, int64, error)
	GetFollowing(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]uuid.UUID, int64, error)
}

type DefaultUserFollowRepository struct {
	res runtime.Resource
}

func NewUserFollowRepository(res runtime.Resource) UserFollowRepository {
	return &DefaultUserFollowRepository{res: res}
}

func (r *DefaultUserFollowRepository) Create(ctx context.Context, followerID, followedID uuid.UUID) error {
	_, err := r.res.DB.
		NewInsert().
		Model(&struct {
			FollowerID uuid.UUID `bun:"follower_id"`
			FollowedID uuid.UUID `bun:"followed_id"`
		}{
			FollowerID: followerID,
			FollowedID: followedID,
		}).
		Exec(ctx)
	return err
}

func (r *DefaultUserFollowRepository) Delete(ctx context.Context, followerID, followedID uuid.UUID) error {
	_, err := r.res.DB.
		NewDelete().
		Model((*struct{})(nil)).
		Table("user_follows").
		Where("follower_id = ?", followerID).
		Where("followed_id = ?", followedID).
		Exec(ctx)
	return err
}

func (r *DefaultUserFollowRepository) IsFollowing(ctx context.Context, followerID, followedID uuid.UUID) (bool, error) {
	count, err := r.res.DB.
		NewSelect().
		Model((*struct{})(nil)).
		Table("user_follows").
		Where("follower_id = ?", followerID).
		Where("followed_id = ?", followedID).
		Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *DefaultUserFollowRepository) GetFollowers(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]uuid.UUID, int64, error) {
	var followers []uuid.UUID

	offset := (page - 1) * pageSize

	count, err := r.res.DB.
		NewSelect().
		Model((*struct{})(nil)).
		Table("user_follows").
		Where("followed_id = ?", userID).
		Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	err = r.res.DB.
		NewSelect().
		Model((*struct{})(nil)).
		Table("user_follows").
		Column("follower_id").
		Where("followed_id = ?", userID).
		Limit(pageSize).
		Offset(offset).
		Scan(ctx, &followers)
	if err != nil {
		return nil, 0, err
	}

	return followers, int64(count), nil
}

func (r *DefaultUserFollowRepository) GetFollowing(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]uuid.UUID, int64, error) {
	var following []uuid.UUID

	offset := (page - 1) * pageSize

	count, err := r.res.DB.
		NewSelect().
		Model((*struct{})(nil)).
		Table("user_follows").
		Where("follower_id = ?", userID).
		Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	err = r.res.DB.
		NewSelect().
		Model((*struct{})(nil)).
		Table("user_follows").
		Column("followed_id").
		Where("follower_id = ?", userID).
		Limit(pageSize).
		Offset(offset).
		Scan(ctx, &following)
	if err != nil {
		return nil, 0, err
	}

	return following, int64(count), nil
}
