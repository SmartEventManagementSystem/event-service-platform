package router

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	"backend/event-service-platform/app/api/controller"
	"backend/event-service-platform/app/api/middleware"
	"backend/event-service-platform/app/database/repository"
	"backend/event-service-platform/app/internal/runtime"
	"backend/event-service-platform/app/internal/validator"
	ctxutil "backend/event-service-platform/app/pkg/util/context"
	echoUtil "backend/event-service-platform/app/pkg/util/echo"
	//_ "backend/event-service-platform/docs"/
)

const (
	// Base paths
	apiV1BasePath = "/api/v1"
	swaggerPath   = "/api/v1/swagger/*"
	healthPath    = "/health"

	// Route prefixes
	authPrefix        = "/auth"
	usersPrefix       = "/users"
	collectionsPrefix = "/collections"
	itemsPrefix       = "/items"
	rarityPrefix      = "/rarity-configs"
	gamePrefix        = "/game"
	inventoryPrefix   = "/inventory"
)

type Router struct {
	*echo.Echo
	res          runtime.Resource
	vals         *validator.Validators
	middleware   *middleware.Middleware
	controllers  *controller.Controllers
	repositories *repository.Repositories
}

// NewRouter creates a new router instance
func NewRouter(
	res runtime.Resource,
	vals *validator.Validators,
	middleware *middleware.Middleware,
	controllers *controller.Controllers,
	repositories *repository.Repositories,
) *Router {
	if controllers == nil {
		panic("controllers cannot be nil")
	}
	if vals == nil {
		panic("validators cannot be nil")
	}

	r := &Router{
		Echo:         echo.New(),
		res:          res,
		vals:         vals,
		middleware:   middleware,
		controllers:  controllers,
		repositories: repositories,
	}

	r.setupEcho()
	r.setupMiddlewares()
	r.setupSwagger()
	r.setupHealthRoutes()
	r.setupRoutes()

	return r
}

func (r *Router) setupEcho() {
	r.Echo.HidePort = true
	r.Echo.HideBanner = true
	r.Echo.Validator = r.vals
}

func (r *Router) setupMiddlewares() {
	r.Echo.Use(echoMiddleware.RequestID())
	r.Echo.Use(echoUtil.SetupCORSMiddleware(r.res))
	r.Echo.Use(echoUtil.SetupLoggerMiddleware(r.res))
}

func (r *Router) setupSwagger() {
	env := ctxutil.GetAppModeFromEnv()
	if env == ctxutil.AppModeDev || env == ctxutil.AppModeLocal {
		r.Echo.Debug = true
		r.Echo.GET(swaggerPath, echoSwagger.WrapHandler)
	}
}

func (r *Router) setupHealthRoutes() {
	r.Echo.GET(healthPath, r.controllers.HealthController.HealthCheck)
}

func (r *Router) setupRoutes() {
	apiGroup := r.Echo.Group(apiV1BasePath)

	r.setupAuthRoutes(apiGroup)
	r.setupUserRoutes(apiGroup)
	r.setupCollectionRoutes(apiGroup)
	r.setupItemRoutes(apiGroup)
	r.setupRarityConfigRoutes(apiGroup)
	r.setupGameRoutes(apiGroup)
	r.setupInventoryRoutes(apiGroup)
	r.setupRestoredRoutes(apiGroup)
}

func (r *Router) setupRestoredRoutes(apiGroup *echo.Group) {
	// Posts
	apiGroup.POST("/posts", r.controllers.PostController.CreatePost)
	apiGroup.GET("/posts", r.controllers.PostController.ListPosts)
	apiGroup.GET("/posts/:id", r.controllers.PostController.GetPost)

	// Comments
	apiGroup.POST("/posts/:id/comments", r.controllers.CommentController.CreateComment)
	apiGroup.GET("/posts/:id/comments", r.controllers.CommentController.ListComments)
	apiGroup.POST("/comments/:id/react", r.controllers.CommentController.ReactToComment)

	// Chat
	apiGroup.POST("/chat/rooms", r.controllers.ChatController.CreateRoom, r.middleware.RequireAuth())
	apiGroup.GET("/chat/rooms", r.controllers.ChatController.ListRooms, r.middleware.RequireAuth())
	apiGroup.POST("/chat/messages", r.controllers.ChatController.SendMessage, r.middleware.RequireAuth())

	// Tickets
	apiGroup.POST("/ticket_sales", r.controllers.TicketController.CreateTicketSale)
	apiGroup.POST("/ticket", r.controllers.TicketController.ListTickets)
	apiGroup.POST("/tickets/validate", r.controllers.TicketController.ValidateTicket)

	// Sessions & Speakers
	apiGroup.POST("/sessions", r.controllers.SessionController.CreateSession)
	apiGroup.GET("/sessions", r.controllers.SessionController.ListSessions)
	apiGroup.POST("/speakers", r.controllers.SpeakerController.CreateSpeaker)
	apiGroup.GET("/speakers", r.controllers.SpeakerController.ListSpeakers)

	// Groups
	apiGroup.POST("/groups", r.controllers.GroupController.CreateGroup)
	apiGroup.GET("/groups", r.controllers.GroupController.GetGroups)

	// Notifications
	apiGroup.GET("/notifications/:user_id", r.controllers.NotificationController.GetNotifications)
	apiGroup.PUT("/notifications/:user_id/:notification_id/read", r.controllers.NotificationController.MarkAsRead)

	// Media
	apiGroup.POST("/chat/upload", r.controllers.MediaController.UploadChatFile)
	apiGroup.GET("/chat/messages/:messageId/attachments", r.controllers.MediaController.GetChatAttachments)

	// Tags
	apiGroup.GET("/tags", r.controllers.TagController.GetAllTags)
	apiGroup.POST("/tags", r.controllers.TagController.CreateTag)

	// Workspace
	apiGroup.POST("/create_workspace", r.controllers.WorkspaceController.CreateWorkspace)
	apiGroup.POST("/list_workspace", r.controllers.WorkspaceController.ListWorkspaces)

	// Subscriptions
	apiGroup.POST("/subscription/create", r.controllers.SubscriptionController.CreateSubscription)
	apiGroup.GET("/subscription/user/:user_id", r.controllers.SubscriptionController.GetUserSubscription)

	// Payment methods
	apiGroup.GET("/payment-methods/:eventId", r.controllers.PaymentController.GetPaymentMethods)
	apiGroup.POST("/payment-methods", r.controllers.PaymentController.CreatePaymentMethod)

	// Newsfeed
	apiGroup.GET("/newsfeed", r.controllers.NewsfeedController.GetUserNewsfeed)
	apiGroup.POST("/newsfeed/invalidate", r.controllers.NewsfeedController.InvalidateNewsfeed)
}

func (r *Router) setupAuthRoutes(apiGroup *echo.Group) {
	authGroup := apiGroup.Group(authPrefix)
	authGroup.POST("/register", r.controllers.AuthController.Register)
	authGroup.POST("/login", r.controllers.AuthController.Login)
	authGroup.POST("/logout", r.controllers.AuthController.Logout)
	authGroup.POST("/refresh-token", r.controllers.AuthController.RefreshToken)
	authGroup.GET("/me", r.controllers.AuthController.Me, r.middleware.RequireAuth())
}

func (r *Router) setupUserRoutes(apiGroup *echo.Group) {
	usersGroup := apiGroup.Group(usersPrefix)
	usersGroup.GET("/balances", r.controllers.UserController.GetBalances, r.middleware.RequireAuth())
	usersGroup.GET("/balances/history", r.controllers.UserController.GetBalanceHistory, r.middleware.RequireAuth())
}

func (r *Router) setupCollectionRoutes(apiGroup *echo.Group) {
	collectionsGroup := apiGroup.Group(collectionsPrefix)
	collectionsGroup.GET("", r.controllers.CollectionController.ListCollections)
	collectionsGroup.GET("/:id", r.controllers.CollectionController.GetCollection)
	collectionsGroup.POST("", r.controllers.CollectionController.CreateCollection, r.middleware.RequireAdmin())
	collectionsGroup.PUT("/:id", r.controllers.CollectionController.UpdateCollection, r.middleware.RequireAdmin())
	collectionsGroup.DELETE("/:id", r.controllers.CollectionController.DeleteCollection, r.middleware.RequireAdmin())
}

func (r *Router) setupItemRoutes(apiGroup *echo.Group) {
	itemsGroup := apiGroup.Group(itemsPrefix)
	itemsGroup.GET("", r.controllers.ItemController.ListItems)
	itemsGroup.GET("/:id", r.controllers.ItemController.GetItem)
	itemsGroup.POST("", r.controllers.ItemController.CreateItem, r.middleware.RequireAdmin())
	itemsGroup.PUT("/:id", r.controllers.ItemController.UpdateItem, r.middleware.RequireAdmin())
	itemsGroup.DELETE("/:id", r.controllers.ItemController.DeleteItem, r.middleware.RequireAdmin())
}

func (r *Router) setupRarityConfigRoutes(apiGroup *echo.Group) {
	rarityGroup := apiGroup.Group(rarityPrefix)
	rarityGroup.GET("", r.controllers.RarityConfigController.ListRarityConfigs)
	rarityGroup.GET("/:id", r.controllers.RarityConfigController.GetRarityConfig)
	rarityGroup.POST("", r.controllers.RarityConfigController.CreateRarityConfig, r.middleware.RequireAdmin())
	rarityGroup.PUT("/:id", r.controllers.RarityConfigController.UpdateRarityConfig, r.middleware.RequireAdmin())
	rarityGroup.DELETE("/:id", r.controllers.RarityConfigController.DeleteRarityConfig, r.middleware.RequireAdmin())
}

func (r *Router) setupGameRoutes(apiGroup *echo.Group) {
	gameGroup := apiGroup.Group(gamePrefix)
	gameGroup.POST("/spin", r.controllers.GamePlayController.Spin, r.middleware.RequireAuth())
	gameGroup.POST("/daily-reward", r.controllers.GamePlayController.ClaimDailyReward, r.middleware.RequireAuth())
	gameGroup.GET("/stats", r.controllers.GamePlayController.GetGameStats, r.middleware.RequireAuth())
}

func (r *Router) setupInventoryRoutes(apiGroup *echo.Group) {
	inventoryGroup := apiGroup.Group(inventoryPrefix)
	inventoryGroup.GET("", r.controllers.InventoryController.GetInventory, r.middleware.RequireAuth())
	inventoryGroup.GET("/:id", r.controllers.InventoryController.GetInventoryItem, r.middleware.RequireAuth())
	inventoryGroup.GET("/collection-progress", r.controllers.InventoryController.GetCollectionProgress, r.middleware.RequireAuth())
}
