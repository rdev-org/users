package main

import (
	"fmt"
	"log"
	"net"

	"github.com/rdev-org/users/internal/container"
)

func main() {
	port := "50051"
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}

	container := container.NewContainer()
	grpcServer := container.Initialize()

	log.Printf("User gRPC server is running on port %s...", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
