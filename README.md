# 💬 GoChat — Real-time Chat Application

> Instant messaging with multiple rooms, live presence indicators, and persistent history — built on Go WebSockets.

## 🛠 Tech Stack

- Go 1.21 + Gorilla WebSocket + Gorilla Mux
- React 18 + Vite + Tailwind CSS
- PostgreSQL
- JWT authentication

## 🏗 Architecture

The core is a single **hub goroutine** that owns all room state and communicates exclusively via channels (`register`, `unregister`, `broadcast`). Each connected client runs two goroutines: `readPump` (WebSocket → hub) and `writePump` (hub → WebSocket). No shared state, no race conditions.

## ✨ Features

- Real-time messaging via WebSocket — zero polling
- Multiple chat rooms with online user count
- Create new rooms from the sidebar
- Message history (last 50 messages loaded on join)
- Live online presence — see who's in the room right now
- Auto-reconnect on dropped connection with 2s delay
- Messages styled left/right (yours vs. others)
- JWT auth with 7-day tokens

## 🚀 Getting Started

### Prerequisites

- Go 1.21+
- Node.js 18+
- PostgreSQL 14+

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/go-chat.git
   cd go-chat
   ```

2. Set up the backend:
   ```bash
   cd backend
   cp .env.example .env
   # Fill in DATABASE_URL and JWT_SECRET
   go mod tidy
   go run .
   # Auto-migrates DB and seeds 3 default rooms
   ```

3. Set up the frontend:
   ```bash
   cd ../frontend
   npm install
   npm run dev
   ```

Open two browser tabs at `http://localhost:5173`, register two users, and see real-time messaging in action.

## 📸 Screenshots

> Run the app to see it in action. Screenshots coming soon.

## 📄 License

MIT
