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

	pbAuth "github.com/rdev-org/proto-contracts/pb/auth"
	pbNotification "github.com/rdev-org/proto-contracts/pb/notification"
	pb "github.com/rdev-org/proto-contracts/pb/user"

	"github.com/rdev-org/users/internal/auth"
	"github.com/rdev-org/users/internal/repository"
	"github.com/rdev-org/users/internal/service"

	"github.com/uptrace/bun"
)

type Container struct {
	db                 *bun.DB
	notificationConn   *grpc.ClientConn
	notificationClient pbNotification.NotificationServiceClient
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

	c.RegisterCredentialService(grpcServer)

	c.RegisterUserService(grpcServer)

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
	c.notificationClient = pbNotification.NewNotificationServiceClient(conn)
	return nil
}

func (c *Container) RegisterUserService(grpcServer *grpc.Server) {
	profileRepo := repository.NewProfileRepository(c.db)
	addressRepo := repository.NewAddressRepository(c.db)
	preferenceRepo := repository.NewPreferenceRepository(c.db)

	s := service.NewUserService(c.notificationClient, profileRepo, addressRepo, preferenceRepo)
	pb.RegisterUserServiceServer(grpcServer, s)
}

func (c *Container) RegisterCredentialService(grpcServer *grpc.Server) {
	r := auth.NewAuthenticationRepository(c.db)
	s := auth.NewAuthenticationService(c.notificationClient, r)

	pbAuth.RegisterAuthenticationServiceServer(grpcServer, s)
}

func unaryLogInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	resp, err := handler(ctx, req)

	duration := time.Since(start)
	st, _ := status.FromError(err)

	log.Printf("[gRPC] Method: %s | Duration: %v | Status: %v", info.FullMethod, duration, st.Code())

	return resp, err
}
