package service

import (
	"github.com/kerti/balances/backend/model"
	"github.com/kerti/balances/backend/repository"
	"github.com/kerti/balances/backend/util/failure"
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

// GetByID fetches a User by its ID
func (s *User) GetByID(id uuid.UUID) (*model.User, error) {
	users, err := s.Repository.ResolveByIDs([]uuid.UUID{id})
	if err != nil {
		return nil, err
	}

	if len(users) != 1 {
		return nil, failure.EntityNotFound("User")
	}

	return &users[0], nil
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
		return model.User{}, failure.EntityNotFound("User")
	}

	user := users[0]
	err = user.Update(input, userID)
	if err != nil {
		return model.User{}, err
	}

	err = s.Repository.Update(user)
	return user, err
}
