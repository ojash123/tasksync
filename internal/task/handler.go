package task

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// TaskHandler holds the dependencies for task-related HTTP handlers.
type TaskHandler struct {
	Store *TaskStore
}

// NewTaskHandler creates a new TaskHandler with a given TaskStore.
func NewTaskHandler(store *TaskStore) *TaskHandler {
	return &TaskHandler{Store: store}
}

// RegisterRoutes sets up the routing for the task endpoints.
func (h *TaskHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/tasks", h.handleTasks)
	mux.HandleFunc("/tasks/", h.handleTaskByID) // Note the trailing slash
}

// handleTasks routes requests for /tasks based on the HTTP method.
func (h *TaskHandler) handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listTasks(w, r)
	case http.MethodPost:
		h.createTask(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleTaskByID routes requests for /tasks/{id} based on the HTTP method.
func (h *TaskHandler) handleTaskByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/tasks/"):] // Extract ID from URL
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getTask(w, r, id)
	case http.MethodPut:
		h.updateTask(w, r, id)
	case http.MethodDelete:
		h.deleteTask(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// --- Handler Functions ---

func (h *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Assign new values for ID, status, and timestamps
	newTask.ID = uuid.NewString()
	newTask.Status = "pending" // Default status
	newTask.LastUpdated = time.Now().UTC()

	h.Store.CreateTask(newTask)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func (h *TaskHandler) listTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.Store.GetAllTasks()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) getTask(w http.ResponseWriter, r *http.Request, id string) {
	task, err := h.Store.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) updateTask(w http.ResponseWriter, r *http.Request, id string) {
	var updatedTask Task
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Ensure the ID and timestamp are correctly set
	updatedTask.ID = id
	updatedTask.LastUpdated = time.Now().UTC()

	if err := h.Store.UpdateTask(id, updatedTask); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTask)
}

func (h *TaskHandler) deleteTask(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.Store.DeleteTask(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
