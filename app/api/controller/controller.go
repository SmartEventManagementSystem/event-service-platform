package controller

import (
	"github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/internal/runtime"
	"github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/manager"
)

type Controllers struct {
	AuthController         *AuthController
	HealthController       *HealthController
	UserController         *UserController
	CollectionController   *CollectionController
	ItemController         *ItemController
	RarityConfigController *RarityConfigController
	GamePlayController     *GamePlayController
	InventoryController    *UserInventoryController
	// Restored controllers from legacy /events
	PostController         *PostController
	CommentController      *CommentController
	ChatController         *ChatController
	TicketController       *TicketController
	SessionController      *SessionController
	SpeakerController      *SpeakerController
	GroupController        *GroupController
	NotificationController *NotificationController
	MediaController        *MediaController
	TagController          *TagController
	WorkspaceController    *WorkspaceController
	SubscriptionController *SubscriptionController
	PaymentController      *PaymentController
	NewsfeedController     *NewsfeedController
}

func NewControllers(managers *manager.Managers, res runtime.Resource) *Controllers {
	return &Controllers{
		AuthController:         NewAuthController(managers, res),
		HealthController:       NewHealthController(managers, res),
		UserController:         NewUserController(managers, res),
		CollectionController:   NewCollectionController(managers, res),
		ItemController:         NewItemController(managers, res),
		RarityConfigController: NewRarityConfigController(managers, res),
		GamePlayController:     NewGamePlayController(managers, res),
		InventoryController:    NewUserInventoryController(managers, res),
		// instantiate restored controllers as skeletons for Plan A
		PostController:         NewPostController(managers, res),
		CommentController:      NewCommentController(managers, res), // Already exists
		ChatController:         NewChatController(managers, res),
		TicketController:       NewTicketController(managers, res),
		SessionController:      NewSessionController(managers, res),
		SpeakerController:      NewSpeakerController(managers, res),
		GroupController:        NewGroupController(managers, res),
		NotificationController: NewNotificationController(managers, res),
		MediaController:        NewMediaController(managers, res),
		TagController:          NewTagController(managers, res),
		WorkspaceController:    NewWorkspaceController(managers, res),
		SubscriptionController: NewSubscriptionController(managers, res),
		PaymentController:      NewPaymentController(managers, res),
		NewsfeedController:     NewNewsfeedController(managers, res),
	}
}
