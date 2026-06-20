package storage

import (
	"database/sql"
	"errors"
)

type Flow struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Data string `json:"data"`
}

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(db *sql.DB) *SQLiteStore {
	return &SQLiteStore{db: db}
}

// Init creates the necessary tables if they don't exist
func (s *SQLiteStore) Init() error {
	query := `
	CREATE TABLE IF NOT EXISTS flows (
		id TEXT PRIMARY KEY,
		name TEXT,
		data TEXT
	);
	`
	_, err := s.db.Exec(query)
	return err
}

// SaveFlow inserts or updates a flow
func (s *SQLiteStore) SaveFlow(id, name, data string) error {
	query := `
	INSERT INTO flows (id, name, data) VALUES (?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET name=excluded.name, data=excluded.data;
	`
	_, err := s.db.Exec(query, id, name, data)
	return err
}

// GetFlow retrieves a flow by ID
func (s *SQLiteStore) GetFlow(id string) (*Flow, error) {
	query := `SELECT id, name, data FROM flows WHERE id = ?`
	row := s.db.QueryRow(query, id)

	var flow Flow
	err := row.Scan(&flow.ID, &flow.Name, &flow.Data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &flow, nil
}

// GetAllFlows retrieves all flows from the database
func (s *SQLiteStore) GetAllFlows() ([]Flow, error) {
	query := `SELECT id, name, data FROM flows`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flows []Flow
	for rows.Next() {
		var f Flow
		if err := rows.Scan(&f.ID, &f.Name, &f.Data); err != nil {
			return nil, err
		}
		flows = append(flows, f)
	}
	return flows, nil
}

// DeleteFlow deletes a flow by ID
func (s *SQLiteStore) DeleteFlow(id string) error {
	query := `DELETE FROM flows WHERE id = ?`
	_, err := s.db.Exec(query, id)
	return err
}
