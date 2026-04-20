package main

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gochat/models"
)

// Hub manages all rooms and their connected clients
type Hub struct {
	// roomID -> set of clients
	rooms map[int]map[*Client]bool
	mu    sync.RWMutex

	register   chan *Client
	unregister chan *Client
	broadcast  chan *RoomMessage
}

type RoomMessage struct {
	RoomID  int
	Payload []byte
}

func newHub() *Hub {
	return &Hub{
		rooms:      make(map[int]map[*Client]bool),
		register:   make(chan *Client, 256),
		unregister: make(chan *Client, 256),
		broadcast:  make(chan *RoomMessage, 256),
	}
}

// Run is the single goroutine that owns all hub state — no locks needed here
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.rooms[client.roomID] == nil {
				h.rooms[client.roomID] = make(map[*Client]bool)
			}
			h.rooms[client.roomID][client] = true
			h.mu.Unlock()
			h.broadcastOnlineUsers(client.roomID)
			log.Printf("[hub] %s joined room %d", client.userName, client.roomID)

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.rooms[client.roomID]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.send)
					if len(clients) == 0 {
						delete(h.rooms, client.roomID)
					}
				}
			}
			h.mu.Unlock()
			h.broadcastOnlineUsers(client.roomID)
			log.Printf("[hub] %s left room %d", client.userName, client.roomID)

		case msg := <-h.broadcast:
			h.mu.RLock()
			clients := h.rooms[msg.RoomID]
			h.mu.RUnlock()
			for client := range clients {
				select {
				case client.send <- msg.Payload:
				default:
					// Slow client — drop and clean up
					h.unregister <- client
				}
			}
		}
	}
}

func (h *Hub) broadcastOnlineUsers(roomID int) {
	h.mu.RLock()
	clients := h.rooms[roomID]
	names := make([]string, 0, len(clients))
	for c := range clients {
		names = append(names, c.userName)
	}
	h.mu.RUnlock()

	msg := models.WSMessage{Type: "online_users", Users: names, RoomID: roomID}
	payload, _ := json.Marshal(msg)
	h.mu.RLock()
	for c := range h.rooms[roomID] {
		select {
		case c.send <- payload:
		default:
		}
	}
	h.mu.RUnlock()
}

func (h *Hub) onlineCount(roomID int) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.rooms[roomID])
}
