package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type Room struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedBy   *int      `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	OnlineCount int       `json:"online_count"`
}

type Message struct {
	ID        int       `json:"id"`
	RoomID    int       `json:"room_id"`
	UserID    int       `json:"user_id"`
	UserName  string    `json:"user_name"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// WebSocket message envelope
type WSMessage struct {
	Type    string   `json:"type"`    // "message" | "join" | "leave" | "online_users"
	Message *Message `json:"message,omitempty"`
	Users   []string `json:"users,omitempty"`
	UserName string  `json:"user_name,omitempty"`
	RoomID  int      `json:"room_id,omitempty"`
}
