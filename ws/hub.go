package ws

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

type Room struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"clients"`
}

// Factory function for the Hub.
func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 10),
	}
}

// Infinite Running loop.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if _, ok := h.Rooms[client.RoomID]; ok {
				room := h.Rooms[client.RoomID]

				if _, ok := room.Clients[client.Username]; !ok {
					room.Clients[client.Username] = client
				}

			}
		case client := <-h.Unregister:
			if _, ok := h.Rooms[client.RoomID]; ok {
				room := h.Rooms[client.RoomID]
				// Broadcast the message that client left the room.
				if len(h.Rooms[client.RoomID].Clients) > 0 {
					h.Broadcast <- &Message{
						Content:  "A user left",
						RoomID:   client.RoomID,
						Username: client.Username,
					}
				}

				if _, ok := room.Clients[client.Username]; ok {
					delete(h.Rooms[client.RoomID].Clients, client.Username)
					close(client.Message)
				}
			}
		case msg := <-h.Broadcast:

			if _, ok := h.Rooms[msg.RoomID]; ok {
				for _, client := range h.Rooms[msg.RoomID].Clients {
					client.Message <- msg
				}
			}
		}
	}
}
