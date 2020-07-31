package hub

import (
	"encoding/json"
	"goproject/streaming/conn"
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

type Subscription struct {
	Conn     *conn.Connection
	Room     string
	Playlist []string
}

type socketMessage struct {
	action string `json:"action"`
	data string `json:"data"`
}

type videoData struct {
	time string `json:"time"`
	start string `json:"start"`
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

		log.Println(objmap)

		var action string
		err = json.Unmarshal(objmap["action"], &action)
		if err != nil {
			log.Fatal(err)
		}
		var data string
		err = json.Unmarshal(objmap["data"], &data)
		if err != nil {
			log.Fatal(err)
		}

		// payload := socketMessage{action: action, data: data}

		if action == "REQUEST" {
			s.Playlist = append(s.Playlist, data)
		}

		// If there is something in the playlist send the video data to all the conns
		if len(s.Playlist) > 0 {
			videoData := getVideoData()

			jsonData, _ := json.Marshal(&videoData)
			var a interface{}
			json.Unmarshal(jsonData, &a)
			log.Println(a)
			m := Message{jsonData, s.Room}
			Instance.Broadcast <- m
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


func getVideoData() videoData {
	data := videoData{time: "2", start: "12"}
	return data
}