package user_credential

import (
	"context"
	"log"

	pb "github.com/rdev-org/proto-contracts/pb/user_credential"
	"github.com/rdev-org/users/internal/constant"
	"github.com/rdev-org/users/internal/jwt"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserCredentialService struct {
	pb.UnimplementedUserCredentialServiceServer
	r UserCredentialRepository
}

func NewUserCredentialService(r UserCredentialRepository) *UserCredentialService {
	return &UserCredentialService{
		r: r,
	}
}

func (s *UserCredentialService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := s.r.Login(ctx, req.GetUsername())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "user not foud")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.GetPassword()))
	if err != nil {
		log.Printf("[UserCredential]service - failed invalid password %v", err)
		return nil, status.Error(codes.Unauthenticated, "invalid password")
	}

	token, err := jwt.GenerateToken(user.Meta.Uuid)
	if err != nil {
		log.Printf("[UserCredential]service - failed to generate token %v", err)
		return nil, status.Error(codes.Unavailable, err.Error())
	}

	return &pb.LoginResponse{
		Token:    token,
		Username: user.Username,
		Email:    user.Email,
		Uuid:     user.Meta.Uuid,
	}, nil
}

func (s *UserCredentialService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	passwordHash := jwt.HashPassword(req.GetPassword())
	payload := &RegisterRequest{
		Email:    req.GetEmail(),
		Password: passwordHash,
	}

	if _, err := s.r.Register(ctx, payload); err != nil {
		log.Printf("[UserCredential]service - failed to register %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RegisterResponse{
		Message: constant.Registration,
	}, nil
}

func (s *UserCredentialService) DeleteCredential(ctx context.Context, req *pb.DeleteCredentialRequest) (*pb.DeleteCredentialResponse, error) {
	err := s.r.DeleteCredential(ctx, req.GetUuid())
	if err != nil {
		log.Printf("[UserCredential]service - failed to delete profile %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeleteCredentialResponse{
		Message: "Credential successfully deleted",
	}, nil
}

func (s *UserCredentialService) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*pb.UpdatePasswordResponse, error) {
	passwordHash := jwt.HashPassword(req.GetPassword())
	payload := &UpdatePasswordRequest{
		Uuid:     req.GetUuid(),
		Password: passwordHash,
	}

	err := s.r.UpdatePassword(ctx, payload)
	if err != nil {
		log.Printf("[UserCredential]service - failed to update password %v", err)
		return nil, status.Error(codes.Unavailable, err.Error())
	}

	return &pb.UpdatePasswordResponse{
		Message: constant.PasswordUpdate,
	}, nil
}

func (s *UserCredentialService) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	passwordHash := jwt.HashPassword(req.GetPassword())
	payload := &UpdatePasswordRequest{
		Uuid:     req.GetUuid(),
		Password: passwordHash,
	}

	// check token

	err := s.r.UpdatePassword(ctx, payload)
	if err != nil {
		log.Printf("[UserCredential]service - failed to reset password %v", err)
		return nil, status.Error(codes.Unavailable, err.Error())
	}

	return &pb.ResetPasswordResponse{
		Message: constant.PasswordReset,
	}, nil
}
