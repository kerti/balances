package service

import (
	"github.com/google/uuid"
	"github.com/kerti/balances/backend/model"
	"github.com/kerti/balances/backend/repository"
	"github.com/kerti/balances/backend/util/failure"
	"github.com/kerti/balances/backend/util/logger"
)

// UserImpl is the service provider implementation
type UserImpl struct {
	Repository repository.User `inject:"userRepository"`
}

// Startup perform startup functions
func (s *UserImpl) Startup() {
	logger.Trace("User Service starting up...")
}

// Shutdown cleans up everything and shuts down
func (s *UserImpl) Shutdown() {
	logger.Trace("User Service shutting down...")
}

// GetByID fetches a User by its ID
func (s *UserImpl) GetByID(id uuid.UUID) (*model.User, error) {
	users, err := s.Repository.ResolveByIDs([]uuid.UUID{id})
	if err != nil {
		return nil, err
	}

	if len(users) != 1 {
		return nil, failure.EntityNotFound("User")
	}

	return &users[0], nil
}

// GetByFilter fetches a set of Users by its filter
func (s *UserImpl) GetByFilter(input model.UserFilterInput) ([]model.User, model.PageInfoOutput, error) {
	return s.Repository.ResolveByFilter(input.ToFilter())
}

// Create creates a new User
func (s *UserImpl) Create(input model.UserInput, userID uuid.UUID) (model.User, error) {
	user := model.NewUserFromInput(input, userID)
	err := s.Repository.Create(user)
	return user, err
}

// Update updates an existing User
func (s *UserImpl) Update(input model.UserInput, userID uuid.UUID) (model.User, error) {
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
