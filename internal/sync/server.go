package sync

import (
	"context"
	"log"

	// Import your internal task package and the generated proto package
	"github.com/ojash123/tasksync/internal/task"
	"github.com/ojash123/tasksync/pkg/proto"
)

// Server implements the gRPC TaskSyncServiceServer interface.
type Server struct {
	// Embed the UnimplementedServer for forward compatibility.
	proto.UnimplementedTaskSyncServiceServer
	// Store is the in-memory task store, shared with the REST API.
	Store *task.TaskStore
}

// SyncTask is the implementation of the RPC method defined in our .proto file.
// It's called when another node sends a task update.
func (s *Server) SyncTask(ctx context.Context, req *proto.SyncRequest) (*proto.SyncResponse, error) {
	// Get the task from the incoming request.
	incomingTaskProto := req.GetTask()
	if incomingTaskProto == nil {
		return nil, nil // Or return an error if a task is always expected
	}

	log.Printf("Received sync request for task ID: %s", incomingTaskProto.Id)

	// Convert the proto message to our internal task struct.
	incomingTask := fromProto(incomingTaskProto)

	// --- Conflict Resolution (Last-Write-Wins) ---
	existingTask, err := s.Store.GetTask(incomingTask.ID)
	// If the task exists and our existing version is newer, ignore the sync request.
	if err == nil && existingTask.LastUpdated.After(incomingTask.LastUpdated) {
		log.Printf("Ignoring sync for task %s; local version is newer.", incomingTask.ID)
		return &proto.SyncResponse{}, nil
	}

	// Otherwise, create or update the task in our local store.
	// Our CreateTask method works as an "upsert" (update or insert).
	s.Store.CreateTask(incomingTask)
	log.Printf("Synced task %s successfully.", incomingTask.ID)

	return &proto.SyncResponse{}, nil
}

// fromProto converts a protobuf TaskMessage to our internal task.Task struct.
func fromProto(p *proto.TaskMessage) task.Task {
	return task.Task{
		ID:             p.Id,
		Title:          p.Title,
		Description:    p.Description,
		Status:         p.Status,
		Priority:       p.Priority,
		AssignedUserID: p.AssignedUserId,
		DueDate:        p.DueDate.AsTime(),
		LastUpdated:    p.LastUpdated.AsTime(),
	}
}
