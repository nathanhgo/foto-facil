package main

import (
	"database/sql"
	"log"
	"net/http"

	"foto-facil-backend/internal/api"
	"foto-facil-backend/internal/storage"

	_ "modernc.org/sqlite"
)

func main() {
	// Initialize SQLite database
	db, err := sql.Open("sqlite", "./foto-facil.db")
	if err != nil {
		log.Fatal("Error opening SQLite database: ", err)
	}
	defer db.Close()

	store := storage.NewSQLiteStore(db)
	if err := store.Init(); err != nil {
		log.Fatal("Error initializing SQLite database tables: ", err)
	}

	api.Store = store
	log.Println("SQLite database initialized successfully at ./foto-facil.db")

	http.HandleFunc("/ws", api.HandleWebSocket)

	log.Println("Backend server starting on port :8080...")
	log.Println("WebSocket endpoint available at ws://localhost:8080/ws")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
