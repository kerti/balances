package gorsk

import (
	"time"

	"github.com/satori/uuid"
)

// Location represents company location model
type Location struct {
	ID        uuid.UUID  `json:"id" gorm:"type:char(36)"`
	Name      string     `json:"name"`
	Active    bool       `json:"active"`
	Address   string     `json:"address"`
	CompanyID uuid.UUID  `json:"company_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" pg:",soft_delete"`
}
