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
	r.HandleFunc("/room", h.handleCreateRoom).Methods("GET")
	r.HandleFunc("/room/{roomID}/playlist", h.handleGetRoomPlaylist).Methods("GET")
}

func (h *WS) handleRoomConn(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	roomID := params["roomID"]

	if len(roomID) == 0 {
		responses.ERROR(w, http.StatusBadRequest, fmt.Errorf("missing roomID param"))
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

func (h *WS) handleCreateRoom(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.NewUUID()
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	rd, err := h.svc.CreateRoom(id.String())
	if err != nil {
		responses.ERROR(w, http.StatusConflict, err)
		return
	}

	log.Println((rd))

	if rd.ID == "" {
		responses.JSON(w, http.StatusConflict, rd)
		return
	}

	responses.JSON(w, http.StatusOK, rd)
}

func (h *WS) handleGetRoomPlaylist(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	roomID := params["roomID"]

	if len(roomID) == 0 {
		responses.ERROR(w, http.StatusBadRequest, fmt.Errorf("missing roomID param"))
		return
	}

	roomExists := hub.CheckRoomAvailability(roomID)
	if !roomExists {
		responses.ERROR(w, http.StatusBadRequest, fmt.Errorf("room does not exist"))
		return
	}

	playlist := h.svc.GetRoomPlaylist(roomID)

	responses.JSON(w, http.StatusOK, playlist)
}
