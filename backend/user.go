package gorsk

import (
	"time"

	"github.com/satori/uuid"
)

// User represents user domain model
type User struct {
	ID uuid.UUID `json:"id" gorm:"type:char(36)"`

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Password  string `json:"-"`
	Email     string `json:"email"`

	Mobile  string `json:"mobile,omitempty"`
	Phone   string `json:"phone,omitempty"`
	Address string `json:"address,omitempty"`

	Active bool `json:"active"`

	LastLogin          *time.Time `json:"last_login,omitempty"`
	LastPasswordChange *time.Time `json:"last_password_change,omitempty"`

	Token string `json:"-"`

	Role *Role `json:"role,omitempty"`

	RoleID     uuid.UUID `json:"-"`
	CompanyID  uuid.UUID `json:"company_id"`
	LocationID uuid.UUID `json:"location_id"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" pg:",soft_delete"`
}

// AuthUser represents data stored in JWT token for user
type AuthUser struct {
	ID         uuid.UUID
	CompanyID  uuid.UUID
	LocationID uuid.UUID
	Username   string
	Email      string
	Role       AccessRole
}

// ChangePassword updates user's password related fields
func (u *User) ChangePassword(hash string) {
	now := time.Now()
	u.Password = hash
	u.LastPasswordChange = &now
}

// UpdateLastLogin updates last login field
func (u *User) UpdateLastLogin(token string) {
	now := time.Now()
	u.Token = token
	u.LastLogin = &now
}
