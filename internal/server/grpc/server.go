package grpc_server

import (
	"context"
	"nbf-user/internal/app"

	userv1 "github.com/hesoyamTM/nbf-protos/gen/go/user"
)

type UserServer struct {
	userv1.UnimplementedUserServer
	userService app.UserService
}

func NewUserServer(userService app.UserService) *UserServer {
	return &UserServer{
		userService: userService,
	}
}

func (s *UserServer) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	return s.userService.CreateUser(ctx, req)
}

func (s *UserServer) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	return s.userService.GetUser(ctx, req)
}

func (s *UserServer) GetUsers(ctx context.Context, req *userv1.GetUsersRequest) (*userv1.GetUsersResponse, error) {
	return s.userService.GetUsers(ctx, req)
}

func (s *UserServer) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error) {
	return s.userService.UpdateUser(ctx, req)
}

func (s *UserServer) DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error) {
	return s.userService.DeleteUser(ctx, req)
}
