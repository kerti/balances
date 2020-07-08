package response

import (
	"github.com/kerti/balances/backend/model"
	"github.com/kerti/balances/backend/util/cachetime"
)

// LoginResponse is the success response to a login request
type LoginResponse struct {
	Expiration cachetime.CacheTime `json:"expiration"`
	Token      string              `json:"token"`
	User       model.UserOutput    `json:"user"`
}

// TokenResponse is the response to a request for a new token
type TokenResponse struct {
	Expiration cachetime.CacheTime `json:"expiration"`
	Token      string              `json:"token"`
}
