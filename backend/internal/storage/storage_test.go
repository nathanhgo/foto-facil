package storage

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestStorageCRUD(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open memory db: %v", err)
	}
	defer db.Close()

	store := NewSQLiteStore(db)
	err = store.Init()
	if err != nil {
		t.Fatalf("Failed to init db schema: %v", err)
	}

	// Test Insert
	flowID := "flow_123"
	flowData := `{"nodes": [], "edges": []}`
	err = store.SaveFlow(flowID, "My First Flow", flowData)
	if err != nil {
		t.Fatalf("Failed to save flow: %v", err)
	}

	// Test Retrieve
	flow, err := store.GetFlow(flowID)
	if err != nil {
		t.Fatalf("Failed to get flow: %v", err)
	}
	if flow.Name != "My First Flow" || flow.Data != flowData {
		t.Errorf("Retrieved data mismatch: got %+v", flow)
	}

	// Test Update
	newData := `{"nodes": [{"id":"1"}], "edges": []}`
	err = store.SaveFlow(flowID, "Updated Flow", newData)
	if err != nil {
		t.Fatalf("Failed to update flow: %v", err)
	}

	flow, _ = store.GetFlow(flowID)
	if flow.Name != "Updated Flow" || flow.Data != newData {
		t.Errorf("Retrieved updated data mismatch: got %+v", flow)
	}
}
