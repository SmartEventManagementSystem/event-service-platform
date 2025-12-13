package manager

import (
	"context"

	"github.com/google/uuid"

	"github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/api/client/request"
	"github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/api/client/response"
	"github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/database/repository"
)

type UserManager interface {
	GetUserProfile(ctx context.Context, req *request.GetUserProfileRequest, currentUserID *uuid.UUID) (*response.UserProfileResponse, error)
	GetUserEvents(ctx context.Context, req *request.GetUserEventsRequest) (*response.UserEventListResponse, error)
	GetUserPosts(ctx context.Context, req *request.GetUserPostsRequest) (*response.UserPostListResponse, error)
	FollowUser(ctx context.Context, req *request.FollowUserRequest, currentUserID uuid.UUID) error
	UnfollowUser(ctx context.Context, req *request.UnfollowUserRequest, currentUserID uuid.UUID) error
	UpdateUserProfile(ctx context.Context, req *request.UpdateUserProfileRequest, currentUserID uuid.UUID) (*response.UserProfileResponse, error)
	UpdateUserSettings(ctx context.Context, req *request.UpdateUserSettingsRequest, currentUserID uuid.UUID) error
}

type DefaultUserManager struct {
	userRepo   repository.UserRepository
	eventRepo  repository.EventRepository
	postRepo   repository.PostRepository
	followRepo repository.UserFollowRepository
}

func NewUserManager(
	userRepo repository.UserRepository,
	eventRepo repository.EventRepository,
	postRepo repository.PostRepository,
	followRepo repository.UserFollowRepository,
) UserManager {
	return &DefaultUserManager{
		userRepo:   userRepo,
		eventRepo:  eventRepo,
		postRepo:   postRepo,
		followRepo: followRepo,
	}
}

func (m *DefaultUserManager) GetUserProfile(ctx context.Context, req *request.GetUserProfileRequest, currentUserID *uuid.UUID) (*response.UserProfileResponse, error) {
	user, err := m.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	// Check if current user follows this user
	var isFollowing bool
	if currentUserID != nil {
		isFollowing, _ = m.followRepo.IsFollowing(ctx, *currentUserID, user.ID)
	}

	// Check if this is the user's own profile
	var isOwnProfile bool
	if currentUserID != nil {
		isOwnProfile = *currentUserID == user.ID
	}

	return &response.UserProfileResponse{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		Firstname:      user.Firstname,
		Lastname:       user.Lastname,
		Bio:            user.Bio,
		WebsiteURL:     user.WebsiteURL,
		Location:       "", // Would need to be added to user entity
		Avatar:         user.Avatar,
		Cover:          user.Cover,
		FollowerCount:  user.FollowerCount,
		FollowingCount: user.FollowingCount,
		PostCount:      user.PostCount,
		EventsCount:    user.EventsCount,
		IsFollowing:    isFollowing,
		IsOwnProfile:   isOwnProfile,
		SocialLinks:    make(map[string]string), // Would need to be added to user entity
		IsE2EEEnabled:  user.IsE2EEEnabled,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}, nil
}

func (m *DefaultUserManager) GetUserEvents(ctx context.Context, req *request.GetUserEventsRequest) (*response.UserEventListResponse, error) {
	user, err := m.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	events, total, err := m.eventRepo.FindByUserID(ctx, user.ID, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	eventResponses := make([]response.UserEventResponse, len(events))
	for i, event := range events {
		eventResponses[i] = response.UserEventResponse{
			ID:               event.ID,
			Title:            event.Title,
			Description:      event.Description,
			Location:         event.Location,
			StartTime:        event.StartTime,
			EndTime:          event.EndTime,
			Status:           string(event.Status),
			Avatar:           event.Avatar,
			CurrentAttendees: event.CurrentAttendees,
			MaxAttendees:     event.MaxAttendees,
			IsPublic:         event.IsPublic,
			Price:            event.Price,
			Currency:         event.Currency,
			CreatedAt:        event.CreatedAt,
		}
	}

	return &response.UserEventListResponse{
		Events:   eventResponses,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func (m *DefaultUserManager) GetUserPosts(ctx context.Context, req *request.GetUserPostsRequest) (*response.UserPostListResponse, error) {
	user, err := m.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	posts, total, err := m.postRepo.FindByUserID(ctx, user.ID, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	postResponses := make([]response.UserPostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = response.UserPostResponse{
			ID:            post.ID,
			Title:         post.Title,
			Content:       post.Content,
			Image:         post.Image,
			LikesCount:    post.LikesCount,
			CommentsCount: post.CommentsCount,
			CreatedAt:     post.CreatedAt,
		}
	}

	return &response.UserPostListResponse{
		Posts:    postResponses,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func (m *DefaultUserManager) FollowUser(ctx context.Context, req *request.FollowUserRequest, currentUserID uuid.UUID) error {
	user, err := m.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return err
	}

	return m.followRepo.Create(ctx, currentUserID, user.ID)
}

func (m *DefaultUserManager) UnfollowUser(ctx context.Context, req *request.UnfollowUserRequest, currentUserID uuid.UUID) error {
	user, err := m.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return err
	}

	return m.followRepo.Delete(ctx, currentUserID, user.ID)
}

func (m *DefaultUserManager) UpdateUserProfile(ctx context.Context, req *request.UpdateUserProfileRequest, currentUserID uuid.UUID) (*response.UserProfileResponse, error) {
	user, err := m.userRepo.FindByID(ctx, currentUserID)
	if err != nil {
		return nil, err
	}

	// Update user fields
	user.Firstname = req.Firstname
	user.Lastname = req.Lastname
	user.Bio = req.Bio
	user.WebsiteURL = req.WebsiteURL
	user.Avatar = req.Avatar
	user.Cover = req.Cover

	updatedUser, err := m.userRepo.Update(ctx, *user)
	if err != nil {
		return nil, err
	}

	return &response.UserProfileResponse{
		ID:             updatedUser.ID,
		Username:       updatedUser.Username,
		Email:          updatedUser.Email,
		Firstname:      updatedUser.Firstname,
		Lastname:       updatedUser.Lastname,
		Bio:            updatedUser.Bio,
		WebsiteURL:     updatedUser.WebsiteURL,
		Location:       "",
		Avatar:         updatedUser.Avatar,
		Cover:          updatedUser.Cover,
		FollowerCount:  updatedUser.FollowerCount,
		FollowingCount: updatedUser.FollowingCount,
		PostCount:      updatedUser.PostCount,
		EventsCount:    updatedUser.EventsCount,
		IsFollowing:    false,
		IsOwnProfile:   true,
		SocialLinks:    make(map[string]string),
		IsE2EEEnabled:  updatedUser.IsE2EEEnabled,
		CreatedAt:      updatedUser.CreatedAt,
		UpdatedAt:      updatedUser.UpdatedAt,
	}, nil
}

func (m *DefaultUserManager) UpdateUserSettings(ctx context.Context, req *request.UpdateUserSettingsRequest, currentUserID uuid.UUID) error {
	user, err := m.userRepo.FindByID(ctx, currentUserID)
	if err != nil {
		return err
	}

	user.IsE2EEEnabled = req.IsE2EEEnabled
	if req.EncryptionPasscode != "" {
		// Hash the encryption passcode before storing
		user.EncryptionPasscode = req.EncryptionPasscode // Should be hashed
	}

	_, err = m.userRepo.Update(ctx, *user)
	return err
}
