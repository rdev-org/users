package container

import (
	"context"
	"log"
	"time"

	"github.com/rdev-org/users/internal/user_address"
	"github.com/rdev-org/users/internal/user_credential"
	"github.com/rdev-org/users/internal/user_preference"
	"github.com/rdev-org/users/internal/user_profile"
	"github.com/uptrace/bun"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	pbAddress "github.com/rdev-org/proto-contracts/pb/user_address"
	pbCredential "github.com/rdev-org/proto-contracts/pb/user_credential"
	pbPreference "github.com/rdev-org/proto-contracts/pb/user_preference"
	pbProfile "github.com/rdev-org/proto-contracts/pb/user_profile"
)

type Container struct {
	db *bun.DB
}

func NewContainer(db *bun.DB) *Container {
	return &Container{db: db}
}

func (c *Container) Initialize() *grpc.Server {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(unaryLogInterceptor),
	)

	c.RegisterUserCredentialService(grpcServer)
	c.RegisterUserProfileService(grpcServer)
	c.RegisterUserAddressService(grpcServer)
	c.RegisterUserPreferenceService(grpcServer)

	reflection.Register(grpcServer)

	return grpcServer
}

func (c *Container) RegisterUserCredentialService(grpcServer *grpc.Server) {
	r := user_credential.NewUserCredentialRepository(c.db)
	s := user_credential.NewUserCredentialService(r)
	pbCredential.RegisterUserCredentialServiceServer(grpcServer, s)
}

func (c *Container) RegisterUserProfileService(grpcServer *grpc.Server) {
	r := user_profile.NewUserProfileRepository(c.db)
	s := user_profile.NewUserProfileService(r)
	pbProfile.RegisterUserProfileServiceServer(grpcServer, s)
}

func (c *Container) RegisterUserAddressService(grpcServer *grpc.Server) {
	r := user_address.NewUserAddressRepository(c.db)
	s := user_address.NewUserAddressService(r)
	pbAddress.RegisterUserAddressServiceServer(grpcServer, s)
}

func (c *Container) RegisterUserPreferenceService(grpcServer *grpc.Server) {
	r := user_preference.NewUserPreferenceRepository(c.db)
	s := user_preference.NewUserPreferenceService(r)
	pbPreference.RegisterUserPreferenceServiceServer(grpcServer, s)
}

func unaryLogInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	resp, err := handler(ctx, req)

	duration := time.Since(start)
	st, _ := status.FromError(err)

	log.Printf("[gRPC] Method: %s | Duration: %v | Status: %v", info.FullMethod, duration, st.Code())

	return resp, err
}
