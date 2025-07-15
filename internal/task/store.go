package task

import (
	"fmt"
	"sync"
)

// TaskStore is an in-memory, thread-safe store for tasks.
type TaskStore struct {
	mu    sync.RWMutex
	tasks map[string]Task
}

// NewTaskStore creates and returns a new TaskStore.
func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks: make(map[string]Task),
	}
}

// CreateTask adds a new task to the store.
func (ts *TaskStore) CreateTask(task Task) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.tasks[task.ID] = task
}

// GetTask retrieves a task from the store by its ID.
func (ts *TaskStore) GetTask(id string) (Task, error) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	task, ok := ts.tasks[id]
	if !ok {
		return Task{}, fmt.Errorf("task with id %s not found", id)
	}
	return task, nil
}

// UpdateTask updates an existing task in the store.
func (ts *TaskStore) UpdateTask(id string, task Task) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	_, ok := ts.tasks[id]
	if !ok {
		return fmt.Errorf("task with id %s not found", id)
	}
	ts.tasks[id] = task
	return nil
}

// DeleteTask removes a task from the store by its ID.
func (ts *TaskStore) DeleteTask(id string) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	_, ok := ts.tasks[id]
	if !ok {
		return fmt.Errorf("task with id %s not found", id)
	}
	delete(ts.tasks, id)
	return nil
}

// GetAllTasks returns all tasks from the store.
func (ts *TaskStore) GetAllTasks() []Task {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	allTasks := make([]Task, 0, len(ts.tasks))
	for _, task := range ts.tasks {
		allTasks = append(allTasks, task)
	}
	return allTasks
}
