package grpc_server

import (
	"context"
	"nbf-user/internal/app"
	"nbf-user/internal/models"

	"github.com/google/uuid"
	userv1 "github.com/hesoyamTM/nbf-protos/gen/go/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	userInfo := req.GetUser()
	if userInfo == nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid UserInfo")
	}
	if err := validateUserInfo(userInfo); err != nil {
		return nil, err
	}

	NewUser := &models.User{
		ID:          uuid.New(),
		Name:        userInfo.Name,
		Surname:     userInfo.Surname,
		Contacts:    userInfo.Contacts,
		Description: userInfo.Description,
	}

	if err := s.userService.CreateUser(ctx, NewUser); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create user: %v", err)
	}

	return &userv1.CreateUserResponse{}, nil
}

func (s *UserServer) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {

	userID, err := parseUserID(req.GetId())
	if err != nil {
		return nil, err
	}

	userModel, err := s.userService.GetUser(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "User not found: %v", err)
	}

	return &userv1.GetUserResponse{
		User: convertToUserInfo(userModel),
	}, nil
}

func (s *UserServer) GetUsers(ctx context.Context, req *userv1.GetUsersRequest) (*userv1.GetUsersResponse, error) {
	var users []*models.User

	if len(req.GetIds()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid format")
	}

	ids := make([]uuid.UUID, 0, len(req.GetIds()))
	for _, string_id := range req.GetIds() {
		id, err := parseUserID(string_id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "Invalid user id: %s", string_id)
		}
		ids = append(ids, id)
	}

	users, err := s.userService.GetUsers(ctx, ids)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get users: %v", err)
	}

	response := &userv1.GetUsersResponse{}
	for _, userModel := range users {
		response.Users = append(response.Users, convertToUserInfo(userModel))
	}

	return response, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error) {

	userInfo := req.GetUser()
	if userInfo == nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid UserInfo")
	}

	userID, err := parseUserID(userInfo.GetId())
	if err != nil {
		return nil, err
	}

	userModel := &models.User{
		ID:          userID,
		Name:        userInfo.Name,
		Surname:     userInfo.Surname,
		Contacts:    userInfo.Contacts,
		Description: userInfo.Description,
	}

	if err := s.userService.UpdateUser(ctx, userModel); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update user: %v", err)
	}

	return &userv1.UpdateUserResponse{}, nil
}

func (s *UserServer) DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error) {

	userID, err := parseUserID(req.GetId())
	if err != nil {
		return nil, err
	}

	if err := s.userService.DeleteUser(ctx, userID); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete user: %v", err)
	}

	return &userv1.DeleteUserResponse{}, nil
}

// ВНУТРЯНКА

func validateUserInfo(userInfo *userv1.UserInfo) error {
	if userInfo.GetName() == "" {
		return status.Error(codes.InvalidArgument, "User name is empty")
	}
	if userInfo.GetSurname() == "" {
		return status.Error(codes.InvalidArgument, "User surname is empty")
	}
	return nil
}

func parseUserID(id string) (uuid.UUID, error) {
	if id == "" {
		return uuid.Nil, status.Error(codes.InvalidArgument, "User id is empty")
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, status.Error(codes.InvalidArgument, "Invalid user id format")
	}

	return userID, nil
}

func convertToUserInfo(userModel *models.User) *userv1.UserInfo {
	return &userv1.UserInfo{
		Id:          userModel.ID.String(),
		Name:        userModel.Name,
		Surname:     userModel.Surname,
		Contacts:    userModel.Contacts,
		Description: userModel.Description,
	}
}
