package container

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	pbEmail "github.com/rdev-org/proto-contracts/pb/notification_email"
	pbAddress "github.com/rdev-org/proto-contracts/pb/user_address"
	pbCredential "github.com/rdev-org/proto-contracts/pb/user_credential"
	pbPreference "github.com/rdev-org/proto-contracts/pb/user_preference"
	pbProfile "github.com/rdev-org/proto-contracts/pb/user_profile"

	"github.com/rdev-org/users/internal/user_address"
	"github.com/rdev-org/users/internal/user_credential"
	"github.com/rdev-org/users/internal/user_preference"
	"github.com/rdev-org/users/internal/user_profile"
	"github.com/uptrace/bun"
)

type Container struct {
	db               *bun.DB
	notificationConn *grpc.ClientConn
	email            pbEmail.NotificationEmailServiceClient
}

func NewContainer(db *bun.DB) *Container {
	return &Container{
		db: db,
	}
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

func (c *Container) Close() {
	if c.notificationConn != nil {
		if err := c.notificationConn.Close(); err != nil {
			log.Printf("error closing notification grpc conn: %v", err)
		}
	}
	if c.db != nil {
		if err := c.db.Close(); err != nil {
			log.Printf("error closing database conn: %v", err)
		}
	}
}

func (c *Container) SetNotificationService(url string) error {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to notification service: %w", err)
	}

	c.notificationConn = conn
	c.email = pbEmail.NewNotificationEmailServiceClient(conn)
	return nil
}

func (c *Container) RegisterUserCredentialService(grpcServer *grpc.Server) {
	r := user_credential.NewUserCredentialRepository(c.db)
	s := user_credential.NewUserCredentialService(c.email, r)
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
