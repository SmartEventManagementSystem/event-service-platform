package repository

import (
	"context"

	"github.com/google/uuid"

	"backend/event-service-platform/app/database/entity"
	"backend/event-service-platform/app/internal/runtime"
)

type EventRepository interface {
	Insert(ctx context.Context, event *entity.Event) (*entity.Event, error)
	Update(ctx context.Context, event *entity.Event) (*entity.Event, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Event, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]entity.Event, int64, error)
	FindByCreatorID(ctx context.Context, creatorID uuid.UUID, page, pageSize int) ([]entity.Event, int64, error)
	FindUpcomingEvents(ctx context.Context, page, pageSize int) ([]entity.Event, int64, error)
	FindEventsByDateRange(ctx context.Context, startTime, endTime string, page, pageSize int) ([]entity.Event, int64, error)
	UpdateStatus(ctx context.Context, eventID uuid.UUID, status entity.EventStatus) error
	UpdateAttendeeCount(ctx context.Context, eventID uuid.UUID, count int) error
}

type DefaultEventRepository struct {
	res runtime.Resource
}

func NewEventRepository(res runtime.Resource) EventRepository {
	return &DefaultEventRepository{res: res}
}

func (r *DefaultEventRepository) Insert(ctx context.Context, event *entity.Event) (*entity.Event, error) {
	err := r.res.DB.
		NewInsert().
		Model(event).
		Returning("*").
		Scan(ctx, event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (r *DefaultEventRepository) Update(ctx context.Context, event *entity.Event) (*entity.Event, error) {
	var e entity.Event
	err := r.res.DB.
		NewUpdate().
		Model(event).
		WherePK().
		Where("deleted_at IS NULL").
		Returning("*").
		Scan(ctx, &e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *DefaultEventRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	_, err := r.res.DB.
		NewUpdate().
		Model((*entity.Event)(nil)).
		Set("deleted_at", "NOW()").
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}

func (r *DefaultEventRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Event, error) {
	event := new(entity.Event)
	err := r.res.DB.
		ReplicaNewSelect().
		Model(event).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (r *DefaultEventRepository) FindByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]entity.Event, int64, error) {
	var events []entity.Event

	offset := (page - 1) * pageSize

	count, err := r.res.DB.
		ReplicaNewSelect().
		Model((*entity.Event)(nil)).
		Where("creator_id = ?", userID).
		Where("deleted_at IS NULL").
		Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	err = r.res.DB.
		ReplicaNewSelect().
		Model(&events).
		Where("creator_id = ?", userID).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Scan(ctx)
	if err != nil {
		return nil, 0, err
	}

	return events, int64(count), nil
}

func (r *DefaultEventRepository) FindByCreatorID(ctx context.Context, creatorID uuid.UUID, page, pageSize int) ([]entity.Event, int64, error) {
	return r.FindByUserID(ctx, creatorID, page, pageSize)
}

func (r *DefaultEventRepository) FindUpcomingEvents(ctx context.Context, page, pageSize int) ([]entity.Event, int64, error) {
	var events []entity.Event

	offset := (page - 1) * pageSize

	count, err := r.res.DB.
		ReplicaNewSelect().
		Model((*entity.Event)(nil)).
		Where("start_time > NOW()").
		Where("deleted_at IS NULL").
		Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	err = r.res.DB.
		ReplicaNewSelect().
		Model(&events).
		Where("start_time > NOW()").
		Where("deleted_at IS NULL").
		Order("start_time ASC").
		Limit(pageSize).
		Offset(offset).
		Scan(ctx)
	if err != nil {
		return nil, 0, err
	}

	return events, int64(count), nil
}

func (r *DefaultEventRepository) FindEventsByDateRange(ctx context.Context, startTime, endTime string, page, pageSize int) ([]entity.Event, int64, error) {
	var events []entity.Event

	offset := (page - 1) * pageSize

	count, err := r.res.DB.
		ReplicaNewSelect().
		Model((*entity.Event)(nil)).
		Where("start_time BETWEEN ? AND ?", startTime, endTime).
		Where("deleted_at IS NULL").
		Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	err = r.res.DB.
		ReplicaNewSelect().
		Model(&events).
		Where("start_time BETWEEN ? AND ?", startTime, endTime).
		Where("deleted_at IS NULL").
		Order("start_time ASC").
		Limit(pageSize).
		Offset(offset).
		Scan(ctx)
	if err != nil {
		return nil, 0, err
	}

	return events, int64(count), nil
}

func (r *DefaultEventRepository) UpdateStatus(ctx context.Context, eventID uuid.UUID, status entity.EventStatus) error {
	_, err := r.res.DB.
		NewUpdate().
		Model((*entity.Event)(nil)).
		Set("status = ?", status).
		Where("id = ?", eventID).
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}

func (r *DefaultEventRepository) UpdateAttendeeCount(ctx context.Context, eventID uuid.UUID, count int) error {
	_, err := r.res.DB.
		NewUpdate().
		Model((*entity.Event)(nil)).
		Set("current_attendees = ?", count).
		Where("id = ?", eventID).
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}
