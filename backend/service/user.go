package service

import (
	"errors"

	"github.com/kerti/balances/backend/model"
	"github.com/kerti/balances/backend/repository"
	"github.com/kerti/balances/backend/util/logger"
	"github.com/satori/uuid"
)

// User is the service provider
type User struct {
	Repository *repository.User `inject:"userRepository"`
}

// Startup perform startup functions
func (s *User) Startup() {
	logger.Trace("User Service starting up...")
}

// Shutdown cleans up everything and shuts down
func (s *User) Shutdown() {
	logger.Trace("User Service shutting down...")
}

// GetByIDs fetches Users by their IDs
func (s *User) GetByIDs(ids []uuid.UUID) ([]model.User, error) {
	return s.Repository.ResolveByIDs(ids)
}

// Create creates a new User
func (s *User) Create(input model.UserInput, userID uuid.UUID) (model.User, error) {
	user := model.NewUserFromInput(input, userID)
	err := s.Repository.Create(user)
	return user, err
}

// Update updates an existing User
func (s *User) Update(input model.UserInput, userID uuid.UUID) (model.User, error) {
	users, err := s.Repository.ResolveByIDs([]uuid.UUID{input.ID})
	if err != nil {
		return model.User{}, err
	}

	if len(users) != 1 {
		return model.User{}, errors.New("failed resolving user for update")
	}

	user := users[0]
	err = user.Update(input, userID)
	if err != nil {
		return model.User{}, err
	}

	err = s.Repository.Update(user)
	return user, err
}
