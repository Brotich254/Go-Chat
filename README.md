# GoChat — Real-time Chat App

Go (Gorilla WebSocket) + React + PostgreSQL.

## Stack
- Backend: Go, Gorilla Mux, Gorilla WebSocket, JWT, bcrypt, PostgreSQL
- Frontend: React, Vite, Tailwind
- Real-time: WebSocket with Go channels/goroutines

## Architecture
The hub is a single goroutine that owns all room state. Clients communicate
with it exclusively through channels — no mutexes needed for the core message flow.
Each client runs two goroutines: readPump and writePump.

## Setup

### Prerequisites
- Go 1.21+
- PostgreSQL
- Node.js 18+

### Database
```bash
psql -U postgres -c "CREATE DATABASE go_chat;"
```

### Backend
```bash
cd backend
cp .env.example .env    # fill in your values
go mod tidy
go run .                # runs on port 8080, auto-migrates DB + seeds 3 rooms
```

### Frontend
```bash
cd frontend
npm install
npm run dev             # runs on port 5173
```

## Features
- Real-time messaging with WebSocket
- Multiple rooms with online user count
- Create new rooms
- Message history (last 50 per room)
- Online users panel (live presence)
- Auto-reconnect on disconnect
- JWT auth
- Messages bubble left/right (yours vs others)
