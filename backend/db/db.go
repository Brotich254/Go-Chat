package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("Database unreachable:", err)
	}
	migrate()
	log.Println("Database connected")
}

func migrate() {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS rooms (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) UNIQUE NOT NULL,
		description TEXT,
		created_by INT REFERENCES users(id),
		created_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		room_id INT REFERENCES rooms(id) ON DELETE CASCADE,
		user_id INT REFERENCES users(id),
		content TEXT NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);

	-- Seed default rooms
	INSERT INTO rooms (name, description, created_by)
	SELECT 'general', 'General discussion', NULL
	WHERE NOT EXISTS (SELECT 1 FROM rooms WHERE name = 'general');

	INSERT INTO rooms (name, description, created_by)
	SELECT 'random', 'Random chat', NULL
	WHERE NOT EXISTS (SELECT 1 FROM rooms WHERE name = 'random');

	INSERT INTO rooms (name, description, created_by)
	SELECT 'tech', 'Tech talk', NULL
	WHERE NOT EXISTS (SELECT 1 FROM rooms WHERE name = 'tech');
	`
	if _, err := DB.Exec(schema); err != nil {
		log.Fatal("Migration failed:", err)
	}
}
