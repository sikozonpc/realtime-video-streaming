package hub

import (
	"log"
)

// Hub maintains the set of active connections and broadcasts messages to the connections
type Hub struct {
	// Registered connections.
	Rooms map[string]map[*Connection]bool
	// Inbound messages from the connections.
	Broadcast chan Message
	// Register requests from the connections.
	Register chan Subscription
	// Unregister requests from connections.
	Unregister chan Subscription
	// RoomsData registred rooms data
	RoomsPlaylist map[string]Playlist
}

// Instance is the global Hub instance that manages the connected subscriptions
var Instance = Hub{
	Broadcast:     make(chan Message),
	Register:      make(chan Subscription),
	Unregister:    make(chan Subscription),
	Rooms:         make(map[string]map[*Connection]bool),
	RoomsPlaylist: make(map[string]Playlist),
}

// Run the Hub instance
func (h *Hub) Run() {
	log.Println("Hub started")

	for {
		select {
		case s := <-h.Register:
			connections := h.Rooms[s.Room]
			if connections == nil {
				connections = make(map[*Connection]bool)
				h.Rooms[s.Room] = connections
			}
			h.Rooms[s.Room][s.Conn] = true

			log.Printf("New client registered %v, to room %s, with %v connections \n",
				s.Conn.Conn.RemoteAddr(),
				s.Room,
				len(h.Rooms[s.Room]),
			)
		case s := <-h.Unregister:
			connections := h.Rooms[s.Room]
			if connections != nil {
				if _, ok := connections[s.Conn]; ok {
					delete(connections, s.Conn)

					// No more users in the room
					if len(h.Rooms[s.Room]) <= 0 {
						h.deleteRoom(s)
					}

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

func (h *Hub) deleteRoom(s Subscription) {
	log.Printf("[server] Room %s deleted because not active users found: %v users \n",
		s.Room,
		len(h.Rooms[s.Room]),
	)

	delete(h.RoomsPlaylist, s.Room)
}
