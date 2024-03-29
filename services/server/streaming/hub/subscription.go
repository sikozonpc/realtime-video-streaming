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

type actionType string

const (
	REQUEST     actionType = "REQUEST"
	PLAY_VIDEO  actionType = "PLAY_VIDEO"
	PAUSE_VIDEO actionType = "PAUSE_VIDEO"
	END_VIDEO   actionType = "END_VIDEO"
	SYNC        actionType = "SYNC"
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

type RequestSocketMessage struct {
	Action string   `json:"action"`
	Data   Playlist `json:"data"`
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

		action, data := unmarshalSocketMessage(msg)

		log.Println(data)

		itemsInPlaylist := len(Instance.RoomsPlaylist[s.Room]) > 0

		switch action {
		case REQUEST:
			// TODO: Proper validate if it's a valid youtube video
			Instance.RoomsPlaylist[s.Room] = Instance.RoomsPlaylist[s.Room].Enqueue(VideoData{
				Time:    0,
				Playing: true,
				Url:     data.Url,
			})

			log.Println("ARR: ", Instance.RoomsPlaylist[s.Room])

			res := RequestSocketMessage{
				Action: "REQUEST",
				Data:   Instance.RoomsPlaylist[s.Room],
			}

			jsData, _ := json.Marshal(res)

			m := Message{jsData, s.Room}
			Instance.Broadcast <- m

		case END_VIDEO:
			if len(Instance.RoomsPlaylist[s.Room]) > 0 {
				Instance.RoomsPlaylist[s.Room] = Instance.RoomsPlaylist[s.Room].Unqueue()
				log.Println("NEW ARR: ", Instance.RoomsPlaylist[s.Room])

				if len(Instance.RoomsPlaylist[s.Room]) > 0 {
					//TODO: Abstract this method
					res := SocketMessage{
						Action: "PLAY_VIDEO",
						Data:   Instance.RoomsPlaylist[s.Room].GetNext(),
					}

					jsData, _ := json.Marshal(res)

					m := Message{jsData, s.Room}
					Instance.Broadcast <- m

				}

				//TODO: Abstract this method
				res := SocketMessage{
					Action: "PLAY_VIDEO",
					Data:   Instance.RoomsPlaylist[s.Room].GetNext(),
				}

				jsData, _ := json.Marshal(res)

				m := Message{jsData, s.Room}
				Instance.Broadcast <- m

			}
		case PLAY_VIDEO:
			if itemsInPlaylist {
				Instance.RoomsPlaylist[s.Room].UpdateCurrent(
					VideoData{
						Time:    data.Time,
						Url:     Instance.RoomsPlaylist[s.Room].GetNext().Url,
						Playing: true,
					})

				res := SocketMessage{
					Action: "PLAY_VIDEO",
					Data:   Instance.RoomsPlaylist[s.Room].GetNext(),
				}

				jsData, _ := json.Marshal(res)

				m := Message{jsData, s.Room}
				Instance.Broadcast <- m
			}
		case PAUSE_VIDEO:
			if itemsInPlaylist {
				Instance.RoomsPlaylist[s.Room].UpdateCurrent(
					VideoData{
						Time:    data.Time,
						Url:     Instance.RoomsPlaylist[s.Room].GetNext().Url,
						Playing: false,
					})

				res := SocketMessage{
					Action: "PAUSE_VIDEO",
					Data:   Instance.RoomsPlaylist[s.Room].GetNext(),
				}

				jsData, _ := json.Marshal(res)

				m := Message{jsData, s.Room}
				Instance.Broadcast <- m
			}
		case SYNC:
			s.SyncToRoom()

		default:
			log.Printf("No valid action sent from Client, ACTION: %v \n", action)
		}

		if itemsInPlaylist {
			log.Printf("CURRENTLY PLAYING %s \n", Instance.RoomsPlaylist[s.Room].GetNext().Url)
		}

		// m := Message{msg, s.Room}
		//Instance.Broadcast <- m
	}
}

// Unpacks the marsheled json data by the socket message
func unmarshalSocketMessage(msg []byte) (actionType, VideoData) {
	var objmap map[string]json.RawMessage
	err := json.Unmarshal(msg, &objmap)
	if err != nil {
		log.Fatalln(err.Error())
	}

	var action actionType
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

	return action, data
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
func (s Subscription) SyncToRoom() {
	if len(Instance.RoomsPlaylist[s.Room]) > 0 {

		msg := SocketMessage{
			Action: "SYNC",
			Data:   Instance.RoomsPlaylist[s.Room].GetNext(),
		}

		log.Printf("Syncing... MSG: %v", msg)

		s.Conn.Conn.WriteJSON(msg)

	}
}
