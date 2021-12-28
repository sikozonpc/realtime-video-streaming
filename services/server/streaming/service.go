package streaming

import (
	"database/sql"
)

// Initialize initializes streaming application service
func Initialize(db *sql.DB) Socket {
	return Socket{db}
}

// Service represents auth service interface
type Service interface {
	ServeRoom(id string)
	CreateRoom(id string) (*RoomData, error)
	GetRoomPlaylist(roomID string) []string
}

// Socket represents streaming application service
type Socket struct {
	DB *sql.DB
}

// RoomData struct
type RoomData struct {
	ID string
}