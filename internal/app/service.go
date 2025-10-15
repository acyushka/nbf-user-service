package app

import (
	"context"
	"fmt"
	"nbf-user/internal/models"

	"github.com/google/uuid"
)

type userService struct {
	userRepo UserDatabase
	//s3 avatarStorage interface
}

func NewUserService(userRepo UserDatabase) *userService { //+s3
	return &userService{
		userRepo: userRepo,
		//s3
	}
}

func (s *userService) CreateUser(ctx context.Context, NewUser *models.User) error {

	if err := s.userRepo.Create(ctx, NewUser); err != nil {
		return fmt.Errorf("Failed to create user: %w", err)
	}

	return nil
}

func (s *userService) GetUser(ctx context.Context, userID uuid.UUID) (*models.User, error) {

	userModel, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("User not found: %w", err)
	}

	return userModel, nil
}

func (s *userService) GetUsers(ctx context.Context, ids []uuid.UUID) ([]*models.User, error) {
	var users []*models.User

	users, err := s.userRepo.GetByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("Failed to get users: %w", err)
	}

	return users, nil
}

func (s *userService) UpdateUser(ctx context.Context, UserModel *models.User) error {

	userForUpdate, err := s.userRepo.GetByID(ctx, UserModel.ID)
	if err != nil {
		return fmt.Errorf("User not found: %w", err)
	}

	updateUserModel(userForUpdate, UserModel)

	if err := s.userRepo.Update(ctx, userForUpdate); err != nil {
		return fmt.Errorf("Failed to update user: %w", err)
	}

	return nil
}

func (s *userService) DeleteUser(ctx context.Context, UserID uuid.UUID) error {

	if _, err := s.userRepo.GetByID(ctx, UserID); err != nil {
		return fmt.Errorf("User not found: %w", err)
	}

	if err := s.userRepo.Delete(ctx, UserID); err != nil {
		return fmt.Errorf("Failed to delete user: %w", err)
	}

	return nil
}

// ВНУТРЯНКА:

func updateUserModel(user *models.User, userModel *models.User) {
	if userModel.Name != "" {
		user.Name = userModel.Name
	}
	if userModel.Surname != "" {
		user.Surname = userModel.Surname
	}
	if len(userModel.Contacts) > 0 {
		user.Contacts = userModel.Contacts
	}
	if userModel.Description != "" {
		user.Description = userModel.Description
	}
}
