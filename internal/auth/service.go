package auth

import (
	"context"
	"fmt"
	"log"

	"github.com/rdev-org/users/internal/constant"
	"github.com/rdev-org/users/internal/jwt"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/rdev-org/proto-contracts/pb/auth"
	pbEmail "github.com/rdev-org/proto-contracts/pb/notification"
)

type CredentialService struct {
	pb.UnimplementedAuthenticationServiceServer
	email pbEmail.NotificationServiceClient
	r     CredentialRepository
}

func NewAuthenticationService(email pbEmail.NotificationServiceClient, r CredentialRepository) *CredentialService {
	return &CredentialService{
		email: email,
		r:     r,
	}
}

func (s *CredentialService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := s.r.Login(ctx, req.GetUsername())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.GetPassword()))
	if err != nil {
		log.Printf("[Credential]service - failed invalid password %v", err)
		return nil, status.Error(codes.Unauthenticated, "invalid password")
	}

	token, err := jwt.GenerateToken(user.Meta.Uuid)
	if err != nil {
		log.Printf("[Credential]service - failed to generate token %v", err)
		return nil, status.Error(codes.Unavailable, err.Error())
	}

	name := user.Profile.FirstName
	if user.Profile.MiddleName != "" {
		name += " " + user.Profile.MiddleName
	}
	name += " " + user.Profile.LastName

	_, err = s.email.SendEmail(ctx, &pbEmail.SendEmailRequest{
		Module:         "LOGIN_EMAIL",
		RequestContext: req.RequestContext,
		Email: &pbEmail.Email{
			Recipient: user.Email,
			Name:      name,
			Subject:   fmt.Sprintf("New sign in from %s", user.Email),
		},
	})
	if err != nil {
		log.Printf("[Credential]service - failed to send email login %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.LoginResponse{
		Token:    token,
		Username: user.Username,
		Email:    user.Email,
		Uuid:     user.Meta.Uuid,
	}, nil
}

func (s *CredentialService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	passwordHash := jwt.HashPassword(req.GetPassword())
	payload := &RegisterRequest{
		Email:    req.GetEmail(),
		Password: passwordHash,
	}

	if _, err := s.r.Register(ctx, payload); err != nil {
		log.Printf("[Credential]service - failed to register %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RegisterResponse{
		Message: constant.Registration,
	}, nil
}

func (s *CredentialService) DeleteCredential(ctx context.Context, req *pb.DeleteCredentialRequest) (*pb.DeleteCredentialResponse, error) {
	err := s.r.DeleteCredential(ctx, req.GetUuid())
	if err != nil {
		log.Printf("[Credential]service - failed to delete credential %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeleteCredentialResponse{
		Message: "Credential successfully deleted",
	}, nil
}

func (s *CredentialService) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*pb.UpdatePasswordResponse, error) {
	passwordHash := jwt.HashPassword(req.GetPassword())
	payload := &UpdatePasswordRequest{
		Uuid:     req.GetUuid(),
		Password: passwordHash,
	}

	err := s.r.UpdatePassword(ctx, payload)
	if err != nil {
		log.Printf("[Credential]service - failed to update password %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdatePasswordResponse{
		Message: constant.PasswordUpdate,
	}, nil
}

func (s *CredentialService) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	passwordHash := jwt.HashPassword(req.GetPassword())
	payload := &UpdatePasswordRequest{
		Uuid:     req.GetUuid(),
		Password: passwordHash,
	}

	// check token

	err := s.r.UpdatePassword(ctx, payload)
	if err != nil {
		log.Printf("[Credential]service - failed to reset password %v", err)
		return nil, status.Error(codes.Unavailable, err.Error())
	}

	return &pb.ResetPasswordResponse{
		Message: constant.PasswordReset,
	}, nil
}

func (s *CredentialService) ForgetPassword(ctx context.Context, req *pb.ForgetPasswordRequest) (*pb.ForgetPasswordResponse, error) {
	_, err := s.r.CheckEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	grpcResp, err := s.email.ResetPasswordEmail(ctx, &pbEmail.ResetPasswordEmailRequest{
		Email: &pbEmail.Email{
			Recipient: req.GetEmail(),
			Subject:   fmt.Sprintf("Reset Password for %s", req.GetEmail()),
		},
	})
	if err != nil {
		log.Printf("[Credential]service - failed to forget password %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ForgetPasswordResponse{
		Message: grpcResp.GetMessage(),
		SentAt:  grpcResp.GetSentAt(),
	}, nil
}
