package streaming

import (
	"fmt"
	"streamserver/streaming/hub"
	"log"
)

// ServeRoom implements the logic to serve a room
func (s Socket) ServeRoom(id string) {
	log.Print("SERVING ROOM")
}

func (s Socket) CreateRoom(id string) (*RoomData, error) {
	roomExists := hub.CheckRoomAvailability(id)
	if roomExists {
		return nil, fmt.Errorf("room already exists")
	}

	return &RoomData{ID: id}, nil
}

func (s Socket) GetRoomPlaylist(roomID string) []string {
	roomPlaylist := hub.Instance.RoomsPlaylist[roomID]
	var playlist []string

	for _, video := range roomPlaylist {
		playlist = append(playlist, video.Url)
	}

	return playlist
}
