package users

import (
	"context"
	"fmt"

	pb "github.com/rdev-org/proto-contracts/pb/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	repo *UserRepository
}

func NewUserServiceServer(repo *UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Login(ctx context.Context, req *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	fmt.Println("Received login request for user:", req.Username)

	if req.Username == "" {
		return nil, status.Error(codes.Unauthenticated, "username must not be empty")
	}

	user, err := s.repo.GetUserByUsername(req.Username)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &pb.UserLoginResponse{Token: "dummy-tokens", Username: user.Username}, nil
}
