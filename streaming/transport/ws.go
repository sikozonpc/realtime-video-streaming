package transport

import (
	"fmt"
	"goproject/responses"
	"goproject/streaming"
	"goproject/streaming/hub"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/gorilla/websocket"
)

// WS represents the web streaming connection
type WS struct {
	svc streaming.Service
}

// Upgrader specifies parameters for upgrading an HTTP connection to a WebSocket connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// NewWS creates a new websocket service
func NewWS(svc streaming.Service, r *mux.Router) {
	h := WS{svc}

	r.HandleFunc("/ws/{roomID}", h.handleRoomConn).Methods("GET")
	r.HandleFunc("/room", h.handleValidateRoom).Methods("GET")
}

func (h *WS) handleRoomConn(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID := vars["roomID"]

	if len(roomID) == 0 {
		fmt.Println("No room ID")
		return
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	/// Upgrade the connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}

	c := &hub.Connection{Send: make(chan []byte, 256), Conn: ws}
	sub := hub.Subscription{Conn: c, Room: roomID}
	hub.Instance.Register <- sub

	// Sync to current room
	sub.SyncToRoom()

	go sub.Write()
	go sub.Read()
}

func (h *WS) handleValidateRoom(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.NewUUID()
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	rd := h.svc.ValidateRoom(id.String())

	if rd.ID != "" {
		responses.JSON(w, http.StatusOK, rd)
		return
	}
	responses.JSON(w, http.StatusConflict, rd)
}
