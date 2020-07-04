package gorsk

import (
	"time"

	"github.com/satori/uuid"
)

// Company represents company model
type Company struct {
	ID        uuid.UUID  `json:"id" gorm:"type:char(36)"`
	Name      string     `json:"name"`
	Active    bool       `json:"active"`
	Locations []Location `json:"locations,omitempty"`
	Owner     User       `json:"owner"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" pg:",soft_delete"`
}
