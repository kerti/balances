package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kerti/balances/backend/util/logger"
)

// InitRoutes initializes the routes
func (s *Server) InitRoutes() {
	logger.Trace("Initializing routes...")
	s.router = mux.NewRouter()

	// Health
	s.router.HandleFunc("/health", s.HealthHandler.HandleHealthCheck).Methods("GET")

	// Authentication/Authorization
	s.router.HandleFunc("/auth/login", s.AuthHandler.HandlePreflight).Methods("OPTIONS")
	s.router.HandleFunc("/auth/login", s.AuthHandler.HandleAuthLogin).Methods("POST")
	s.router.HandleFunc("/auth/token", s.AuthHandler.HandleGetToken).Methods("GET")

	// Users
	s.router.HandleFunc("/users/{id}", s.UserHandler.HandleGetUserByID).Methods("GET")
	s.router.HandleFunc("/users", s.UserHandler.HandleCreateUser).Methods("POST")
	s.router.HandleFunc("/users/{id}", s.UserHandler.HandleUpdateUser).Methods("PATCH")

	http.Handle("/", s.router)
}
