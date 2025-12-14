package manager

import (
	"backend/event-service-platform/app/database/repository"
	"backend/event-service-platform/app/internal/runtime"
	"backend/event-service-platform/app/pkg/bcrypt"
	"backend/event-service-platform/app/pkg/jwt"
	"backend/event-service-platform/app/pkg/queue"
)

type Managers struct {
	AuthManager         AuthManager
	JobManager          JobManager
	UserBalanceManager  UserBalanceManager
	CollectionManager   CollectionManager
	ItemManager         ItemManager
	RarityConfigManager RarityConfigManager
	UserManager         UserManager
	PostManager         PostManager
	CommentManager      CommentManager
	SocialClient        interface{}
}

func NewManagers(
	res runtime.Resource,
	_ interface{},
	repositories *repository.Repositories,
) *Managers {
	// Create bcrypt hasher from configuration
	bcryptHasher := bcrypt.NewBcrypt(res.Config.BcryptConfig.Cost)
	hasher := &bcryptHasher

	// Create a JWT manager from configuration
	jwtManager := jwt.NewJwt(res.Config.JwtConfig)

	// Initialize job-related components
	redisQueue := queue.NewRedisQueue(res.Redis.GetUniversalClient(), res.Logger)
	jobManager := NewJobManager(repositories.JobRepository, redisQueue, res.Logger)

	// Initialize user manager
	userManager := NewUserManager(
		repositories.UserRepository,
		repositories.EventRepository,
		repositories.PostRepository,
		repositories.UserFollowRepository,
	)

	return &Managers{
		AuthManager:         NewAuthManager(res, hasher, jwtManager, repositories),
		JobManager:          jobManager,
		UserBalanceManager:  NewUserBalanceManager(repositories.UserBalanceRepository, repositories.UserBalanceTransactionRepository),
		CollectionManager:   NewCollectionManager(res.DB, repositories.CollectionRepository, repositories.CollectionItemRepository),
		ItemManager:         NewItemManager(res.DB, repositories.ItemRepository),
		RarityConfigManager: NewRarityConfigManager(repositories.RarityConfigRepository),
		UserManager:         userManager,
		PostManager:         NewPostManager(repositories.PostRepository),
		CommentManager:      NewCommentManager(repositories.CommentRepository),
		SocialClient:        nil,
	}
}
