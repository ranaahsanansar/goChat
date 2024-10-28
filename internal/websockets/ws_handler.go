package websockets

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ranaahsanansar/gochat/initializers"
	"github.com/ranaahsanansar/gochat/middleware"
	"github.com/ranaahsanansar/gochat/models"
)

type Handler struct {
	hub *Hub
}

func NewHandler(h *Hub) *Handler {
	return &Handler{
		hub: h,
	}
}

func (h *Handler) CreateRoomIfNotExists(roomId string, roomName string) {

	if _, ok := h.hub.Rooms[roomId]; !ok {

		r := &Room{
			ID:      roomId,
			Name:    roomName,
			Clients: make(map[string]*Client),
		}
		h.hub.Rooms[roomId] = r
	}
}

func (h *Handler) DeleteRoomIfLastUserLeft(roomId string) {
	if room, ok := h.hub.Rooms[roomId]; ok {
		if len(room.Clients) == 0 {
			delete(h.hub.Rooms, roomId)
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) JoinRoom(c *gin.Context) {

	isValid := middleware.JoinRoomGuard(c)

	if !isValid {
		return
	}
	username, _ := c.Get("username")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomID := c.Param("roomId")
	var userModel models.Users
	result := initializers.DB.First(&userModel, "username = ?", username)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"message": result.Error,
		})
		return
	}
	clientID := userModel.ID
	cl := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		ID:       strconv.FormatUint(uint64(clientID), 10),
		RoomID:   roomID,
		Username: username.(string),
	}

	m := &Message{
		Content:  "A new user has joined the room",
		RoomID:   roomID,
		Username: username.(string),
	}

	h.CreateRoomIfNotExists(roomID, roomID)

	h.hub.Register <- cl
	h.hub.Broadcast <- m

	go cl.listenMessage()
	cl.sendMessage(h.hub)
}
