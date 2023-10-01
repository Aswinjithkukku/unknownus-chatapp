package ws

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Handler struct {
	hub *Hub
}

// Upgrader.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Factory function for Handler.
func NewHandler(h *Hub) *Handler {
	return &Handler{
		hub: h,
	}
}

// Join Room => /ws/:roomId?username=username.
func (h *Handler) HandleWS(c *gin.Context) {
	username := c.Query("username")

	if len(username) == 0 {
		c.JSON(401, gin.H{
			"error": "Provide valid username.",
		})
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	var client *Client
	var msg *Message

	var gotRoom bool

	for _, room := range h.hub.Rooms {
		if len(room.Clients) < 2 {
			fmt.Println(room.Clients)

			client = &Client{
				Conn:     conn,
				Message:  make(chan *Message, 10),
				RoomID:   room.ID,
				Username: username,
			}

			msg = &Message{
				Content:  "A user has joined!",
				RoomID:   room.ID,
				Username: username,
			}

			gotRoom = true
			break

		}
	}

	if !gotRoom {
		randomId := uuid.New()

		h.hub.Rooms[randomId.String()] = &Room{
			ID:      randomId.String(),
			Name:    randomId.String() + "-Room",
			Clients: make(map[string]*Client),
		}

		client = &Client{
			Conn:     conn,
			Message:  make(chan *Message, 10),
			RoomID:   randomId.String(),
			Username: username,
		}

		msg = &Message{
			Content:  "A user has joined!",
			RoomID:   randomId.String(),
			Username: username,
		}

	}

	// roomID := c.Param("roomId")

	// Register new cliwnt in register channel.
	h.hub.Register <- client
	// Broadcast the message
	h.hub.Broadcast <- msg
	// Read and Write message.

	go client.WriteMessage()
	client.ReadMessage(h.hub)
}
