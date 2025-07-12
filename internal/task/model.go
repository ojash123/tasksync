package task

import "time"

// Task represents a single task in the system.
type Task struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Status         string    `json:"status"`   // e.g., "pending", "in-progress", "completed"
	Priority       string    `json:"priority"` // e.g., "low", "medium", "high"
	AssignedUserID string    `json:"assigned_user_id"`
	DueDate        time.Time `json:"due_date"`
	LastUpdated    time.Time `json:"last_updated"` // Used for conflict resolution
}
