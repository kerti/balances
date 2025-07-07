package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kerti/balances/backend/config"
	"github.com/kerti/balances/backend/database"
	"github.com/kerti/balances/backend/handler"
	"github.com/kerti/balances/backend/inject"
	"github.com/kerti/balances/backend/repository"
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
	container.RegisterService("bankAccountRepository", new(repository.BankAccountMySQLRepo))
	container.RegisterService("userRepository", new(repository.UserMySQLRepo))
	container.RegisterService("vehicleRepository", new(repository.VehicleMySQLRepo))
	container.RegisterService("propertyRepository", new(repository.PropertyMySQLRepo))

	// Prepare containers - services
	container.RegisterService("authService", new(service.AuthImpl))
	container.RegisterService("bankAccountService", new(service.BankAccountImpl))
	container.RegisterService("userService", new(service.UserImpl))
	container.RegisterService("vehicleService", new(service.VehicleImpl))
	container.RegisterService("propertyService", new(service.PropertyImpl))

	// Prepare containers - handlers
	container.RegisterService("authHandler", new(handler.AuthImpl))
	container.RegisterService("bankAccountHandler", new(handler.BankAccountImpl))
	container.RegisterService("healthHandler", new(handler.HealthImpl))
	container.RegisterService("userHandler", new(handler.UserImpl))
	container.RegisterService("vehicleHandler", new(handler.VehicleImpl))
	container.RegisterService("propertyHandler", new(handler.PropertyImpl))

	// Prepare containers - HTTP server
	var s server.Server
	container.RegisterService("server", &s)

	// call this after all dependencies are registered
	if err := container.Ready(); err != nil {
		logger.Fatal("Failed to populate services -- %v", err)
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
	ch := make(chan os.Signal, 1)
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
