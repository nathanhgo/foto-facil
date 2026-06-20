package dag

import (
	"testing"
)

func TestSchedulerTopologicalSort(t *testing.T) {
	// Simple graph: Node1 -> Node2 -> Node3
	nodes := map[string]*Node{
		"node1": {ID: "node1", Dependencies: []string{}},
		"node2": {ID: "node2", Dependencies: []string{"node1"}},
		"node3": {ID: "node3", Dependencies: []string{"node2"}},
	}

	scheduler := NewScheduler()
	order, err := scheduler.Sort(nodes)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(order) != 3 {
		t.Fatalf("Expected order length 3, got %d", len(order))
	}

	if order[0] != "node1" || order[1] != "node2" || order[2] != "node3" {
		t.Errorf("Expected order [node1, node2, node3], got %v", order)
	}
}

func TestSchedulerDetectCycle(t *testing.T) {
	// Cyclic graph: Node1 -> Node2 -> Node3 -> Node1
	nodes := map[string]*Node{
		"node1": {ID: "node1", Dependencies: []string{"node3"}},
		"node2": {ID: "node2", Dependencies: []string{"node1"}},
		"node3": {ID: "node3", Dependencies: []string{"node2"}},
	}

	scheduler := NewScheduler()
	_, err := scheduler.Sort(nodes)

	if err == nil {
		t.Fatalf("Expected cycle detection error, got nil")
	}
}
