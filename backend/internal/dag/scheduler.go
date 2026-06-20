package dag

import (
	"errors"
)

// Node represents a generic processing node in the DAG
type Node struct {
	ID           string
	Dependencies []string
}

// Scheduler handles the execution order of the DAG
type Scheduler struct{}

// NewScheduler creates a new DAG Scheduler
func NewScheduler() *Scheduler {
	return &Scheduler{}
}

// Sort performs a topological sort on the given nodes.
// Returns an ordered slice of Node IDs, or an error if a cycle is detected.
func (s *Scheduler) Sort(nodes map[string]*Node) ([]string, error) {
	var order []string
	visited := make(map[string]bool)
	tempMark := make(map[string]bool)
	var sortErr error

	var visit func(nodeID string)
	visit = func(nodeID string) {
		if sortErr != nil {
			return
		}
		if tempMark[nodeID] {
			sortErr = errors.New("cycle detected in DAG")
			return
		}
		if !visited[nodeID] {
			tempMark[nodeID] = true
			node := nodes[nodeID]
			if node != nil {
				for _, depID := range node.Dependencies {
					visit(depID)
				}
			}
			visited[nodeID] = true
			tempMark[nodeID] = false
			order = append(order, nodeID)
		}
	}

	for id := range nodes {
		if !visited[id] {
			visit(id)
		}
	}

	if sortErr != nil {
		return nil, sortErr
	}

	return order, nil
}
