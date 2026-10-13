package main

import (
	"log"

	"github.com/rdev-org/helpers/app"
	"github.com/rdev-org/users/internal/config"
	"github.com/rdev-org/users/internal/container"
)

func main() {
	appConfig, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db := config.ConnectDB(appConfig.Database)
	c := container.NewContainer(db)

	if err := c.SetNotificationService(appConfig.Services.Notification); err != nil {
		log.Fatalf("failed to set notification service: %v", err)
	}

	grpcServer := c.Initialize(app.NewGRPCServer())

	srvCfg := app.ServerConfig{
		Name: "User",
		Host: appConfig.Server.Host,
		Port: appConfig.Server.Port,
	}

	app.RunGRPC(srvCfg, grpcServer, func() {
		c.Close()
	})
}
