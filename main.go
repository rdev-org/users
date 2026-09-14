package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/rdev-org/users/internal/config"
	"github.com/rdev-org/users/internal/container"
)

func main() {
	appConfig, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config %v", err)
	}

	addr := fmt.Sprintf("%s:%s", appConfig.Server.Host, appConfig.Server.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", appConfig.Server.Port, err)
	}

	db := config.ConnectDB(appConfig.Database)
	containerInstance := container.NewContainer(db)
	grpcServer := containerInstance.Initialize()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("User gRPC server is running on port %s...", appConfig.Server.Port)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	sig := <-sigChan
	log.Printf("Received signal %v, initiating graceful shutdown...", sig)

	grpcServer.GracefulStop()
	db.Close()
	log.Println("Server stopped gracefully.")
}
