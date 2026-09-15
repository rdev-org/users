package container

import (
	"context"
	"log"
	"time"

	"github.com/rdev-org/users/internal/user_addresses"
	"github.com/rdev-org/users/internal/users"
	"github.com/uptrace/bun"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	pbUser "github.com/rdev-org/proto-contracts/pb/user"
	pbUserAddr "github.com/rdev-org/proto-contracts/pb/user_address"
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
	c.RegisterUserAddressService(grpcServer)

	reflection.Register(grpcServer)

	return grpcServer
}

func (c *Container) RegisterUserService(grpcServer *grpc.Server) {
	repo := users.NewUserRepository(c.db)
	svc := users.NewUserServiceServer(repo)
	pbUser.RegisterUserServiceServer(grpcServer, svc)
}

func (c *Container) RegisterUserAddressService(grpcServer *grpc.Server) {
	repo := user_addresses.NewUserAddressRepository(c.db)
	svc := user_addresses.NewUserAddressServiceServer(repo)
	pbUserAddr.RegisterUserAddressServiceServer(grpcServer, svc)
}

func unaryLogInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	resp, err := handler(ctx, req)

	duration := time.Since(start)
	st, _ := status.FromError(err)

	log.Printf("[gRPC] Method: %s | Duration: %v | Status: %v", info.FullMethod, duration, st.Code())

	return resp, err
}
