package service

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/kerti/balances/backend/config"
	"github.com/kerti/balances/backend/model"
	"github.com/kerti/balances/backend/repository"
	"github.com/kerti/balances/backend/util/cachetime"
	"github.com/kerti/balances/backend/util/failure"
	"github.com/kerti/balances/backend/util/logger"
)

// Auth is the service provider interface
type Auth interface {
	Startup()
	Shutdown()
	Authenticate(basic string) (authInfo *model.AuthenticationInfo, err error)
	Authorize(bearer string) (userID *uuid.UUID, err error)
	GetToken(user model.User) (token *string, expiration *time.Time, err error)
}

// AuthImpl is the service provider implementation
type AuthImpl struct {
	UserRepository repository.User `inject:"userRepository"`
}

// Startup performs startup functions
func (s *AuthImpl) Startup() {
	logger.Trace("Auth service starting up...")
}

// Shutdown cleans up everything and shuts down
func (s *AuthImpl) Shutdown() {
	logger.Trace("Auth service shutting down...")
}

// Authenticate performs authentication
func (s *AuthImpl) Authenticate(basic string) (authInfo *model.AuthenticationInfo, err error) {
	identity, password, err := s.validateBasicAuthHeader(basic)
	if err != nil {
		return nil, err
	}

	user, err := s.UserRepository.ResolveByIdentity(identity)
	if err != nil {
		return nil, err
	}

	matched := user.ComparePassword(password)
	if !matched {
		return nil, failure.Unauthorized("authentication failed")
	}

	token, expiration, err := s.signJWT(&user)
	return &model.AuthenticationInfo{
		Expiration: expiration,
		Token:      token,
		User:       &user,
	}, err
}

// Authorize authorizes a request based on its Bearer token
func (s *AuthImpl) Authorize(bearer string) (userID *uuid.UUID, err error) {
	config := config.Get()
	jwtToken, err := s.validateBearerAuthHeader(bearer)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, err := s.getUserID(claims)
		if err != nil {
			return nil, err
		}
		err = s.checkExpiration(claims)
		if err != nil {
			return nil, err
		}
		return userID, nil
	}
	return nil, failure.Unauthorized("invalid JWT token")
}

// GetToken signs a new token for a specified user
func (s *AuthImpl) GetToken(user model.User) (token *string, expiration *time.Time, err error) {
	return s.signJWT(&user)
}

// ValidateBasicAuthHeader validates basic authentication header
func (s *AuthImpl) validateBasicAuthHeader(basic string) (string, string, error) {
	auth := strings.SplitN(basic, " ", 2)
	if len(auth) != 2 || auth[0] != "Basic" {
		return "", "", failure.BadRequestFromString("invalid authentication request")
	}

	payload, _ := base64.StdEncoding.DecodeString(auth[1])

	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 {
		return "", "", failure.BadRequestFromString("invalid authentication request")
	}

	return pair[0], pair[1], nil
}

func (s *AuthImpl) validateBearerAuthHeader(bearer string) (string, error) {
	auth := strings.SplitN(bearer, " ", 2)
	if len(auth) != 2 || auth[0] != "Bearer" {
		return "", failure.BadRequestFromString("invalid authentication header")
	}

	return auth[1], nil
}

func (s *AuthImpl) signJWT(user *model.User) (*string, *time.Time, error) {
	config := config.Get()
	now := time.Now()
	expTime := now.Add(config.JWT.Expiration)
	expiration := cachetime.CacheTime(expTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         base64.StdEncoding.EncodeToString([]byte(user.ID.String())),
		"created":    user.ToOutput().Created,
		"expiration": expiration,
		"iss":        "balances",
	})
	tokenString, err := token.SignedString([]byte(config.JWT.Secret))
	return &tokenString, &expTime, err
}

func (s *AuthImpl) getUserID(claims jwt.MapClaims) (*uuid.UUID, error) {
	userIDBase64 := claims["id"].(string)
	userIDBytes, err := base64.StdEncoding.DecodeString(userIDBase64)
	if err != nil {
		return nil, failure.Unauthorized("failed decoding ID from JWT token")
	}
	userID, err := uuid.Parse(string(userIDBytes))
	if err != nil {
		return nil, failure.Unauthorized("failed decoding ID from JWT token")
	}
	return &userID, nil
}

func (s *AuthImpl) checkExpiration(claims jwt.MapClaims) error {
	expiration := int64(claims["expiration"].(float64))
	expTime := time.Unix(expiration/1000, 0)
	if expTime.Before(time.Now()) {
		return failure.Unauthorized("token has expired")
	}
	return nil
}
