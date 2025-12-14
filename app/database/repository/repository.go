package repository

import (
	"backend/event-service-platform/app/internal/runtime"
)

type Repositories struct {
	UserRepository                   UserRepository
	SessionRepository                SessionRepository
	JobRepository                    JobRepository
	UserBalanceRepository            UserBalanceRepository
	UserBalanceTransactionRepository UserBalanceTransactionRepository
	CollectionRepository             CollectionRepository
	CollectionItemRepository         CollectionItemRepository
	ItemRepository                   ItemRepository
	RarityConfigRepository           RarityConfigRepository
	EventRepository                  EventRepository
	PostRepository                   PostRepository
	CommentRepository                CommentRepository
	UserFollowRepository             UserFollowRepository
}

func NewRepositories(res runtime.Resource) *Repositories {
	return &Repositories{
		UserRepository:                   NewUserRepository(res),
		SessionRepository:                NewSessionRepository(res),
		JobRepository:                    NewJobRepository(res),
		UserBalanceRepository:            NewUserBalanceRepository(res),
		UserBalanceTransactionRepository: NewUserBalanceTransactionRepository(res),
		CollectionRepository:             NewCollectionRepository(res),
		CollectionItemRepository:         NewCollectionItemRepository(res),
		ItemRepository:                   NewItemRepository(res),
		RarityConfigRepository:           NewRarityConfigRepository(res),
		EventRepository:                  NewEventRepository(res),
		PostRepository:                   NewPostRepository(res),
		CommentRepository:                NewCommentRepository(res),
		UserFollowRepository:             NewUserFollowRepository(res),
	}
}
