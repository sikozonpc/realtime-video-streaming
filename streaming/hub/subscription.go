package hub

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Maximum message size allowed from peer.
	maxMessageSize = 512
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

type Message struct {
	Data []byte
	Room string
}

// TODO: NEED TO PUSH ROOM VIDEO DATA TO ROOM NOT SUB!!!

type Subscription struct {
	Conn *Connection
	Room string
}

type SocketMessage struct {
	Action string    `json:"action"`
	Data   VideoData `json:"data"`
}

type VideoData struct {
	Url     string  `json:"url"`
	Time    float32 `json:"time"`
	Playing bool    `json:"playing"`
}

// Read pumps messages from the conn connection to the hub.
func (s Subscription) Read() {
	c := s.Conn

	defer func() {
		Instance.Unregister <- s
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	timeAllowedToRead := time.Now().Add(pongWait)
	c.Conn.SetReadDeadline(timeAllowedToRead)
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		var objmap map[string]json.RawMessage
		err = json.Unmarshal(msg, &objmap)
		if err != nil {
			log.Fatalln(err.Error())
		}

		var action string
		err = json.Unmarshal(objmap["action"], &action)
		if err != nil {
			log.Fatal(err)
		}

		var data VideoData
		if objmap["data"] != nil {
			err = json.Unmarshal(objmap["data"], &data)
			if err != nil {
				log.Fatal(err)
			}
		}

		log.Println(data)

		if action == "REQUEST" {
			// TODO: Implement queue separatly

			Instance.RoomsVideoData[s.Room] = VideoData{
				Time:    0,
				Playing: false,
				Url:     data.Url,
			}
		}

		log.Printf("CUUR PLAYING %s \n", Instance.RoomsVideoData[s.Room].Url)

		if action == "PLAY_VIDEO" {
			// TODO: Abstract these methods
			// If there is something in the playlist send the video data to all the conns
			if Instance.RoomsVideoData[s.Room].Url != "" {
				Instance.RoomsVideoData[s.Room] = VideoData{
					Time:    data.Time,
					Url:     Instance.RoomsVideoData[s.Room].Url,
					Playing: true,
				}

				res := SocketMessage{
					Action: "PLAY_VIDEO",
					Data:   Instance.RoomsVideoData[s.Room],
				}

				jsData, _ := json.Marshal(res)

				m := Message{jsData, s.Room}
				Instance.Broadcast <- m
			}
		}

		if action == "PAUSE_VIDEO" {
			//TODO: check for playlist
			if Instance.RoomsVideoData[s.Room].Url != "" {
				Instance.RoomsVideoData[s.Room] = VideoData{
					Time:    data.Time,
					Url:     Instance.RoomsVideoData[s.Room].Url,
					Playing: false,
				}

				log.Println(Instance.RoomsVideoData[s.Room].Time)

				res := SocketMessage{
					Action: "PAUSE_VIDEO",
					Data:   Instance.RoomsVideoData[s.Room],
				}

				jsData, _ := json.Marshal(res)

				m := Message{jsData, s.Room}
				Instance.Broadcast <- m
			}
		}

		m := Message{msg, s.Room}
		Instance.Broadcast <- m
	}
}

// Write writes messages from the hub to the streaming connection
func (s Subscription) Write() {
	c := s.Conn
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.Write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// SyncToRoom sends current media playing data to connection
func (s Subscription) SyncToRoom(roomID string) {
	msg := SocketMessage{
		Action: "SYNC",
		Data:   Instance.RoomsVideoData[roomID],
	}
	log.Println(msg)
	s.Conn.Conn.WriteJSON(msg)
}
