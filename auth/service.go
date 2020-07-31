package auth

import (
	"database/sql"
	"goproject/user"
)

// Initialize initializes auth application service
func Initialize(db *sql.DB) Auth {
	return *New(db)
}

// New creates new service
func New(db *sql.DB) *Auth {
	return &Auth{db}
}

// Service represents auth service interface
type Service interface {
	Register(creds Credentials) error
	GetUsers() ([]user.User, error)
}

// Auth represents auth application service
type Auth struct {
	DB *sql.DB
}

// Credentials for registering a user
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
