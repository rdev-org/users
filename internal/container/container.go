package container

import (
	"context"
	"log"
	"time"

	"github.com/rdev-org/users/internal/users"
	"github.com/uptrace/bun"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	pb "github.com/rdev-org/proto-contracts/pb/user"
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

	c.RegisterUserService(grpcServer)

	reflection.Register(grpcServer)

	return grpcServer
}

func (c *Container) RegisterUserService(grpcServer *grpc.Server) {
	userRepo := users.NewUserRepository(c.db)
	userService := users.NewUserServiceServer(userRepo)
	pb.RegisterUserServiceServer(grpcServer, userService)
}

func unaryLogInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	resp, err := handler(ctx, req)

	duration := time.Since(start)
	st, _ := status.FromError(err)

	log.Printf("[gRPC] Method: %s | Duration: %v | Status: %v", info.FullMethod, duration, st.Code())

	return resp, err
}
