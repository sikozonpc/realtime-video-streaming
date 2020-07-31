package hub

import (
	"goproject/streaming/conn"
	"log"
)

// Hub maintains the set of active connections and broadcasts messages to the connections
type Hub struct {
	// Registered connections.
	Rooms map[string]map[*conn.Connection]bool
	// Inbound messages from the connections.
	Broadcast chan Message
	// Register requests from the connections.
	Register chan Subscription
	// Unregister requests from connections.
	Unregister chan Subscription
}

var Instance = Hub{
	Broadcast:  make(chan Message),
	Register:   make(chan Subscription),
	Unregister: make(chan Subscription),
	Rooms:      make(map[string]map[*conn.Connection]bool),
}

func (h *Hub) Run() {
	log.Println("Hub started")

	for {
		select {
		case s := <-h.Register:
			connections := h.Rooms[s.Room]
			if connections == nil {
				connections = make(map[*conn.Connection]bool)
				h.Rooms[s.Room] = connections
			}
			h.Rooms[s.Room][s.Conn] = true
			log.Printf("New client registered %v \n", s)
		case s := <-h.Unregister:
			connections := h.Rooms[s.Room]
			if connections != nil {
				if _, ok := connections[s.Conn]; ok {
					delete(connections, s.Conn)
					close(s.Conn.Send)
					if len(connections) == 0 {
						delete(h.Rooms, s.Room)
					}
				}
			}
		case m := <-h.Broadcast:
			connections := h.Rooms[m.Room]
			for c := range connections {
				select {
				case c.Send <- m.Data:
				default:
					close(c.Send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(h.Rooms, m.Room)
					}
				}
			}
		}
	}
}

// CheckRoomAvailability checks if a room exists
func CheckRoomAvailability(id string) bool {
	connections := Instance.Rooms[id]
	if connections == nil {
		return true
	}
	return false
}