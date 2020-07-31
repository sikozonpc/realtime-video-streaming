package streaming

import (
	"database/sql"
)

// Initialize initializes streaming application service
func Initialize(db *sql.DB) Socket {
	return *New(db)
}

// New creates new service
func New(db *sql.DB) *Socket {
	return &Socket{db}
}

// Service represents auth service interface
type Service interface {
	ServeRoom(id string)
	ValidateRoom(id string) RoomData
}

// Socket represents streaming application service
type Socket struct {
	DB *sql.DB
}

type RoomData struct {
	ID string
}