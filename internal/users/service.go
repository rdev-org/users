package users

import (
	"context"
	"fmt"

	pb "github.com/rdev-org/proto-contracts/pb/user"
	"github.com/rdev-org/users/internal/jwt"
	"github.com/rdev-org/users/internal/models"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	repo IUserRepository
}

func NewUserServiceServer(repo IUserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Login(ctx context.Context, req *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	if req.Username == "" {
		return nil, status.Error(codes.Unauthenticated, "username must not be empty")
	}

	user, err := s.repo.GetUserByUsernameOrEmail(ctx, req.Username)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid password")
	}

	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return nil, status.Error(codes.Unavailable, fmt.Sprintf("failed to generate token %v", err))
	}

	resp := &pb.UserLoginResponse{
		Token:    token,
		Username: user.Username,
	}

	return resp, nil
}

func (s *UserService) Register(ctx context.Context, req *pb.UserRegisterRequest) (*pb.UserRegisterResponse, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt password %w", err)
	}

	payload := &models.User{
		Username:   req.Username,
		Email:      req.Email,
		Mobile:     req.Mobile,
		Password:   string(passwordHash),
		FirstName:  req.FirstName,
		MiddleName: req.MiddleName,
		LastName:   req.LastName,
	}

	if err := s.repo.RegisterUser(ctx, payload); err != nil {
		return nil, fmt.Errorf("failed to register user %w", err)
	}

	resp := &pb.UserRegisterResponse{
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
	}

	return resp, nil
}
