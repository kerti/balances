package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kerti/balances/backend/repository"

	"github.com/kerti/balances/backend/config"
	"github.com/kerti/balances/backend/database"
	"github.com/kerti/balances/backend/handler"
	"github.com/kerti/balances/backend/inject"
	"github.com/kerti/balances/backend/server"
	"github.com/kerti/balances/backend/service"
	"github.com/kerti/balances/backend/util/logger"
)

func main() {
	// Register logger
	logger.SetupLoggerAuto("", "")

	// Initialize config
	config.Get()

	// Prepare containers
	container := inject.NewContainer()

	// Prepare containers - database
	var db database.MySQL
	container.RegisterService("mysql", &db)

	// Prepare containers - repositories
	container.RegisterService("userRepository", new(repository.User))

	// Prepare containers - services
	container.RegisterService("authService", new(service.Auth))
	container.RegisterService("userService", new(service.User))

	// Prepare containers - handlers
	container.RegisterService("authHandler", new(handler.Auth))
	container.RegisterService("healthHandler", new(handler.Health))
	container.RegisterService("userHandler", new(handler.User))

	// Prepare containers - HTTP server
	var s server.Server
	container.RegisterService("server", &s)

	// call this after all dependencies are registered
	if err := container.Ready(); err != nil {
		logger.Fatal("Failed to populate services %v", err)
	} else {
		logger.Info("Service registry started successfully.")
	}

	// Handle shutdown
	handleShutdown(container)

	// Run server
	s.Serve()
}

// handle graceful shutdown
func handleShutdown(container inject.ServiceContainer) {
	config := config.Get()
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func(ch chan os.Signal) {
		<-ch
		defer os.Exit(0)
		duration := config.Server.ShutdownPeriod
		logger.Info("SIGTERM received. Waiting %v seconds to shutdown", duration.Seconds())
		container.PrepareShutdown()
		time.Sleep(duration)
		logger.Info("Cleaning up resources...")
		container.Shutdown()
		logger.Info("Bye!")
	}(ch)
}
