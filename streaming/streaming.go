package streaming

import (
	"goproject/streaming/hub"
	"log"
)

// ServeRoom implements the logic to serve a room
func (s Socket) ServeRoom(id string) {
	log.Print("SERVING ROOM")
}

func (s Socket) ValidateRoom(id string) RoomData {
	v := hub.CheckRoomAvailability(id)
	if v == true {
		return RoomData{ID: id}
	}
	return RoomData{ID: ""}
}

