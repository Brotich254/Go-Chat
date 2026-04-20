package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gochat/db"
	"github.com/gochat/middleware"
	"github.com/gochat/models"
)

func GetRooms(hub interface{ onlineCount(int) int }) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.DB.Query(`SELECT id, name, COALESCE(description,''), created_at FROM rooms ORDER BY id`)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var rooms []models.Room
		for rows.Next() {
			var room models.Room
			rows.Scan(&room.ID, &room.Name, &room.Description, &room.CreatedAt)
			room.OnlineCount = hub.onlineCount(room.ID)
			rooms = append(rooms, room)
		}
		if rooms == nil {
			rooms = []models.Room{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rooms)
	}
}

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var room models.Room
	err := db.DB.QueryRow(
		`INSERT INTO rooms (name, description, created_by) VALUES ($1, $2, $3)
		 RETURNING id, name, COALESCE(description,''), created_at`,
		req.Name, req.Description, userID,
	).Scan(&room.ID, &room.Name, &room.Description, &room.CreatedAt)
	if err != nil {
		http.Error(w, "Room name already taken", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(room)
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room_id")
	if roomID == "" {
		http.Error(w, "room_id required", http.StatusBadRequest)
		return
	}

	rows, err := db.DB.Query(`
		SELECT m.id, m.room_id, m.user_id, u.name, m.content, m.created_at
		FROM messages m JOIN users u ON m.user_id = u.id
		WHERE m.room_id = $1
		ORDER BY m.created_at DESC LIMIT 50`, roomID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		rows.Scan(&msg.ID, &msg.RoomID, &msg.UserID, &msg.UserName, &msg.Content, &msg.CreatedAt)
		messages = append(messages, msg)
	}
	// Reverse so oldest first
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	if messages == nil {
		messages = []models.Message{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
