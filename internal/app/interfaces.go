package app

import (
	"context"
	"nbf-user/internal/models"

	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(ctx context.Context, NewUser *models.User) error
	GetUser(ctx context.Context, UserID uuid.UUID) (*models.User, error)
	GetUsers(ctx context.Context, ids []uuid.UUID) ([]*models.User, error)
	UpdateUser(ctx context.Context, UserModel *models.User) error
	DeleteUser(ctx context.Context, UserID uuid.UUID) error
}

type UserDatabase interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
