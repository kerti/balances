package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kerti/balances/backend/handler/response"
	"github.com/kerti/balances/backend/util/logger"
)

// InitRoutes initializes the routes
func (s *Server) InitRoutes() {
	logger.Trace("Initializing routes...")
	s.router = mux.NewRouter()

	// Preflight
	s.router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response.RespondWithNoContent(w)
	})

	// Health
	s.router.HandleFunc("/health", s.HealthHandler.HandleHealthCheck).Methods("GET")

	// Authentication/Authorization
	s.router.HandleFunc("/auth/login", s.AuthHandler.HandleAuthLogin).Methods("POST")
	s.router.HandleFunc("/auth/token", s.AuthHandler.HandleGetToken).Methods("GET")

	// Users
	s.router.HandleFunc("/users/{id}", s.UserHandler.HandleGetUserByID).Methods("GET")
	s.router.HandleFunc("/users/search", s.UserHandler.HandleGetUserByFilter).Methods("POST")
	s.router.HandleFunc("/users", s.UserHandler.HandleCreateUser).Methods("POST")
	s.router.HandleFunc("/users/{id}", s.UserHandler.HandleUpdateUser).Methods("PATCH")

	// Banks
	s.router.HandleFunc("/bankAccounts", s.BankAccountHandler.HandleCreateBankAccount).Methods("POST")
	s.router.HandleFunc("/bankAccounts/{id}", s.BankAccountHandler.HandleGetBankAccountByID).Methods("GET")
	s.router.HandleFunc("/bankAccounts/search", s.BankAccountHandler.HandleGetBankAccountByFilter).Methods("POST")
	s.router.HandleFunc("/bankAccounts/{id}", s.BankAccountHandler.HandleUpdateBankAccount).Methods("PATCH")
	s.router.HandleFunc("/bankAccounts/{id}", s.BankAccountHandler.HandleDeleteBankAccount).Methods("DELETE")
	s.router.HandleFunc("/bankAccounts/balances", s.BankAccountHandler.HandleCreateBankAccountBalance).Methods("POST")
	s.router.HandleFunc("/bankAccounts/balances/{id}", s.BankAccountHandler.HandleGetBankAccountBalanceByID).Methods("GET")
	s.router.HandleFunc("/bankAccounts/balances/search", s.BankAccountHandler.HandleGetBankAccountBalanceByFilter).Methods("POST")
	s.router.HandleFunc("/bankAccounts/balances/{id}", s.BankAccountHandler.HandleUpdateBankAccountBalance).Methods("PATCH")
	s.router.HandleFunc("/bankAccounts/balances/{id}", s.BankAccountHandler.HandleDeleteBankAccountBalance).Methods("DELETE")

	// Vehicles
	s.router.HandleFunc("/vehicles", s.VehicleHandler.HandleCreateVehicle).Methods("POST")
	s.router.HandleFunc("/vehicles/{id}", s.VehicleHandler.HandleGetVehicleByID).Methods("GET")
	s.router.HandleFunc("/vehicles/search", s.VehicleHandler.HandleGetVehicleByFilter).Methods("POST")
	s.router.HandleFunc("/vehicles/{id}", s.VehicleHandler.HandleUpdateVehicle).Methods("PATCH")
	s.router.HandleFunc("/vehicles/{id}", s.VehicleHandler.HandleDeleteVehicle).Methods("DELETE")
	s.router.HandleFunc("/vehicles/values", s.VehicleHandler.HandleCreateVehicleValue).Methods("POST")
	s.router.HandleFunc("/vehicles/values/{id}", s.VehicleHandler.HandleGetVehicleValueByID).Methods("GET")
	s.router.HandleFunc("/vehicles/values/search", s.VehicleHandler.HandleGetVehicleValueByFilter).Methods("POST")
	s.router.HandleFunc("/vehicles/values/{id}", s.VehicleHandler.HandleUpdateVehicleValue).Methods("PATCH")
	s.router.HandleFunc("/vehicles/values/{id}", s.VehicleHandler.HandleDeleteVehicleValue).Methods("DELETE")

	// Properties
	s.router.HandleFunc("/properties", s.PropertyHandler.HandleCreateProperty).Methods("POST")
	s.router.HandleFunc("/properties/{id}", s.PropertyHandler.HandleGetPropertyByID).Methods("GET")
	s.router.HandleFunc("/properties/search", s.PropertyHandler.HandleGetPropertyByFilter).Methods("POST")
	s.router.HandleFunc("/properties/{id}", s.PropertyHandler.HandleUpdateProperty).Methods("PATCH")
	s.router.HandleFunc("/properties/{id}", s.PropertyHandler.HandleDeleteProperty).Methods("DELETE")
	s.router.HandleFunc("/properties/values", s.PropertyHandler.HandleCreatePropertyValue).Methods("POST")
	s.router.HandleFunc("/properties/values/{id}", s.PropertyHandler.HandleGetPropertyValueByID).Methods("GET")
	s.router.HandleFunc("/properties/values/search", s.PropertyHandler.HandleGetPropertyValueByFilter).Methods("POST")
	s.router.HandleFunc("/properties/values/{id}", s.PropertyHandler.HandleUpdatePropertyValue).Methods("PATCH")
	s.router.HandleFunc("/properties/values/{id}", s.PropertyHandler.HandleDeletePropertyValue).Methods("DELETE")

	http.Handle("/", s.router)
}
