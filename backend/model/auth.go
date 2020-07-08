package model

import "time"

// AuthenticationInfo is the wrapper object for authentication information
type AuthenticationInfo struct {
	Expiration *time.Time
	Token      *string
	User       *User
}
