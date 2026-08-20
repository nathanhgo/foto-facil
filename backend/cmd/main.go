package main

import (
	"database/sql"
	"log"
	"net/http"

	"foto-facil-backend/internal/api"
	"foto-facil-backend/internal/config"
	"foto-facil-backend/internal/storage"

	_ "modernc.org/sqlite"
)

func main() {
	cfg := config.Load()

	// Initialize SQLite database
	db, err := sql.Open("sqlite", cfg.DBPath)
	if err != nil {
		log.Fatal("Error opening SQLite database: ", err)
	}
	defer db.Close()

	store := storage.NewSQLiteStore(db)
	if err := store.Init(); err != nil {
		log.Fatal("Error initializing SQLite database tables: ", err)
	}

	api.Store = store
	log.Printf("SQLite database initialized successfully at %s\n", cfg.DBPath)

	http.HandleFunc("/ws", api.HandleWebSocket)

	addr := ":" + cfg.WSPort
	log.Printf("Backend server starting on port %s...\n", addr)
	log.Printf("WebSocket endpoint available at ws://localhost:%s/ws\n", cfg.WSPort)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
