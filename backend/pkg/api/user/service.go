package user

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo"
	"github.com/satori/uuid"

	gorsk "github.com/kerti/balances/backend"
	"github.com/kerti/balances/backend/pkg/api/user/platform/pgsql"
)

// Service represents user application interface
type Service interface {
	Create(echo.Context, gorsk.User) (gorsk.User, error)
	List(echo.Context, gorsk.Pagination) ([]gorsk.User, error)
	View(echo.Context, uuid.UUID) (gorsk.User, error)
	Delete(echo.Context, uuid.UUID) error
	Update(echo.Context, Update) (gorsk.User, error)
}

// New creates new user application service
func New(db *pg.DB, udb UDB, rbac RBAC, sec Securer) *User {
	return &User{db: db, udb: udb, rbac: rbac, sec: sec}
}

// Initialize initalizes User application service with defaults
func Initialize(db *pg.DB, rbac RBAC, sec Securer) *User {
	return New(db, pgsql.User{}, rbac, sec)
}

// User represents user application service
type User struct {
	db   *pg.DB
	udb  UDB
	rbac RBAC
	sec  Securer
}

// Securer represents security interface
type Securer interface {
	Hash(string) string
}

// UDB represents user repository interface
type UDB interface {
	Create(orm.DB, gorsk.User) (gorsk.User, error)
	View(orm.DB, uuid.UUID) (gorsk.User, error)
	List(orm.DB, *gorsk.ListQuery, gorsk.Pagination) ([]gorsk.User, error)
	Update(orm.DB, gorsk.User) error
	Delete(orm.DB, gorsk.User) error
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	User(echo.Context) gorsk.AuthUser
	EnforceUser(echo.Context, uuid.UUID) error
	AccountCreate(echo.Context, uuid.UUID, uuid.UUID, uuid.UUID) error
	IsLowerRole(echo.Context, gorsk.AccessRole) error
}
