package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/rs/cors"

	"github.com/gochat/db"
	"github.com/gochat/handlers"
	"github.com/gochat/middleware"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func serveWS(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomIDStr := r.URL.Query().Get("room_id")
		roomID, err := strconv.Atoi(roomIDStr)
		if err != nil || roomID == 0 {
			http.Error(w, "invalid room_id", http.StatusBadRequest)
			return
		}

		userID := r.Context().Value(middleware.UserIDKey).(int)
		userName := r.Context().Value(middleware.UserNameKey).(string)

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("ws upgrade error: %v", err)
			return
		}

		client := &Client{
			hub:      hub,
			conn:     conn,
			send:     make(chan []byte, 256),
			roomID:   roomID,
			userID:   userID,
			userName: userName,
		}

		hub.register <- client

		// Send join notification
		joinMsg, _ := json.Marshal(map[string]interface{}{
			"type":      "join",
			"user_name": userName,
			"room_id":   roomID,
		})
		hub.broadcast <- &RoomMessage{RoomID: roomID, Payload: joinMsg}

		go client.writePump()
		go client.readPump()
	}
}

func main() {
	godotenv.Load()
	db.Init()

	hub := newHub()
	go hub.run()

	r := mux.NewRouter()

	// Public routes
	r.HandleFunc("/api/auth/register", handlers.Register).Methods("POST")
	r.HandleFunc("/api/auth/login", handlers.Login).Methods("POST")
	r.Handle("/api/rooms", middleware.Auth(http.HandlerFunc(handlers.GetRooms(hub)))).Methods("GET")
	r.Handle("/api/rooms", middleware.Auth(http.HandlerFunc(handlers.CreateRoom))).Methods("POST")
	r.Handle("/api/messages", middleware.Auth(http.HandlerFunc(handlers.GetMessages))).Methods("GET")
	r.Handle("/ws", middleware.Auth(http.HandlerFunc(serveWS(hub))))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, c.Handler(r)))
}
