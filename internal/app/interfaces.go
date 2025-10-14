package app

import (
	"context"
	"nbf-user/internal/models"

	"github.com/google/uuid"
	userv1 "github.com/hesoyamTM/nbf-protos/gen/go/user"
)

type UserService interface {
	CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error)
	GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error)
	GetUsers(ctx context.Context, req *userv1.GetUsersRequest) (*userv1.GetUsersResponse, error)
	UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error)
	DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error)
}

type UserDatabase interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

//s3 interface
