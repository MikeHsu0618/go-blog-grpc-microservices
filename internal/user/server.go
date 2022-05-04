package user

import (
	"context"

	v1 "blog-grpc-microservices/api/protobuf/user/v1"
	"blog-grpc-microservices/internal/pkg/exception"
	"blog-grpc-microservices/internal/pkg/log"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewServer(logger *log.Logger, repo Repository) v1.UserServiceServer {
	return &Server{
		logger: logger,
		repo:   repo,
	}
}

type Server struct {
	v1.UnimplementedUserServiceServer
	logger *log.Logger
	repo   Repository
}

func (s *Server) ListUsersByIDs(ctx context.Context, in *v1.ListUsersByIDsRequest) (*v1.ListUsersByIDsResponse, error) {
	ids := in.GetIds()
	users, err := s.repo.ListUsersByIDs(ctx, ids)
	if err != nil {
		return nil, status.Errorf(codes.Internal, exception.Msg.User.GetUserListByIDsFail, err)
	}
	protoUsers := make([]*v1.User, len(users))
	for i, user := range users {
		protoUsers[i] = entityToProtobuf(user)
	}

	return &v1.ListUsersByIDsResponse{
		Users: protoUsers,
	}, nil
}

func (s *Server) GetUser(ctx context.Context, in *v1.GetUserRequest) (*v1.GetUserResponse, error) {
	id := in.GetId()
	user, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, exception.Msg.User.GetUserByIDFail, err)
	}

	return &v1.GetUserResponse{
		User: entityToProtobuf(user),
	}, nil
}
func (s *Server) GetUserByEmail(ctx context.Context, in *v1.GetUserByEmailRequest) (*v1.GetUserResponse, error) {
	email := in.GetEmail()
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, exception.Msg.User.GetUserByEmailFail, err)
	}

	// 如果有傳密碼，則需要驗證密碼
	if in.GetPassword() != "" {
		ok := isCorrectPassword(user.Password, in.GetPassword())
		if !ok {
			return nil, status.Errorf(codes.Internal, exception.Msg.User.IncorrectPassword, err)
		}
	}

	return &v1.GetUserResponse{
		User: entityToProtobuf(user),
	}, nil
}

func (s *Server) GetUserByUsername(ctx context.Context, in *v1.GetUserByUsernameRequest) (*v1.GetUserResponse, error) {
	username := in.GetUsername()
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, exception.Msg.User.GetUserByUsernameFail)
	}

	// 如果有傳密碼，則需要驗證密碼
	if in.GetPassword() != "" {
		ok := isCorrectPassword(user.Password, in.GetPassword())
		if !ok {
			return nil, status.Errorf(codes.Internal, exception.Msg.User.IncorrectPassword, err)
		}
	}

	return &v1.GetUserResponse{
		User: entityToProtobuf(user),
	}, nil
}
func (s *Server) CreateUser(ctx context.Context, in *v1.CreateUserRequest) (*v1.CreateUserResponse, error) {
	password, err := generateFromPassword(in.GetUser().GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, exception.Msg.User.GeneratePasswordFail, err)
	}
	user := &User{
		UUID:     uuid.New().String(),
		Username: in.GetUser().GetUsername(),
		Email:    in.GetUser().GetEmail(),
		Avatar:   in.GetUser().GetAvatar(),
		Password: password,
	}

	err = s.repo.Create(ctx, user)

	// TODO 判端 AlreadyExist

	if err != nil {
		return nil, status.Errorf(codes.Internal, exception.Msg.User.CreateUserFail, err)
	}
	return &v1.CreateUserResponse{
		User: entityToProtobuf(user),
	}, nil
}

func (s *Server) UpdateUser(ctx context.Context, in *v1.UpdateUserRequest) (*v1.UpdateUserResponse, error) {
	user, err := s.repo.Get(ctx, in.GetUser().GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, exception.Msg.User.GetUserByIDFail, err)
	}

	if in.GetUser().GetUsername() != "" {
		user.Username = in.GetUser().GetUsername()
	}
	if in.GetUser().GetEmail() != "" {
		user.Email = in.GetUser().GetEmail()
	}
	if in.GetUser().GetAvatar() != "" {
		user.Avatar = in.GetUser().GetAvatar()
	}

	if in.GetUser().GetPassword() != "" {
		password, err := generateFromPassword(in.GetUser().GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, exception.Msg.User.GeneratePasswordFail)
		}
		user.Password = password
	}

	err = s.repo.Update(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, exception.Msg.User.UpdateUserFail, err)
	}

	return &v1.UpdateUserResponse{
		Success: true,
	}, nil
}
func (s *Server) DeleteUser(ctx context.Context, in *v1.DeleteUserRequest) (*v1.DeleteUserResponse, error) {
	id := in.GetId()
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, exception.Msg.User.UpdateUserFail, err)
	}

	return &v1.DeleteUserResponse{
		Success: true,
	}, nil
}

func entityToProtobuf(user *User) *v1.User {
	return &v1.User{
		Id:        user.ID,
		Uuid:      user.UUID,
		Username:  user.Username,
		Email:     user.Email,
		Avatar:    user.Avatar,
		Password:  user.Password,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

func isCorrectPassword(password string, input string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(input))
	return err == nil
}

func generateFromPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}
