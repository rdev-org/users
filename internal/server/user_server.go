package server

import (
	"context"
	"fmt"

	pb "github.com/rdev-org/proto-contracts/pb/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
}

func NewUserServiceServer() *UserServiceServer {
	return &UserServiceServer{}
}

func (s *UserServiceServer) Login(ctx context.Context, req *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	fmt.Println("Received login request for user:", req.Username)

	if req.Username == "" {
		return nil, status.Error(codes.Unauthenticated, "username must not be empty")
	}

	return &pb.UserLoginResponse{Token: "dummy-token-" + req.Username}, nil
}
