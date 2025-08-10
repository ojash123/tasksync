package sync

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ojash123/tasksync/internal/task"
	"github.com/ojash123/tasksync/pkg/proto"
)

// BroadcastTask sends a task update to a list of peer nodes.
func BroadcastTask(peers []string, t task.Task) {
	// Convert our internal task to the protobuf message format.
	protoTask := toProto(t)

	// Iterate over each peer and send the update in a separate goroutine.
	for _, addr := range peers {
		go sendUpdateToPeer(addr, protoTask)
	}
}

// sendUpdateToPeer connects to a single peer and sends the task.
func sendUpdateToPeer(addr string, protoTask *proto.TaskMessage) {
	log.Printf("Attempting to sync task %s to peer %s", protoTask.Id, addr)

	// Set up a connection to the server.
	// We use WithInsecure() because we are not setting up TLS for this MVP.
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to peer %s: %v", addr, err)
		return
	}
	defer conn.Close()

	// Create a new client.
	client := proto.NewTaskSyncServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Call the remote SyncTask method.
	_, err = client.SyncTask(ctx, &proto.SyncRequest{Task: protoTask})
	if err != nil {
		log.Printf("Failed to sync task with peer %s: %v", addr, err)
		return
	}

	log.Printf("Successfully synced task %s with peer %s", protoTask.Id, addr)
}

// toProto converts our internal task.Task struct to a protobuf TaskMessage.
func toProto(t task.Task) *proto.TaskMessage {
	return &proto.TaskMessage{
		Id:             t.ID,
		Title:          t.Title,
		Description:    t.Description,
		Status:         t.Status,
		Priority:       t.Priority,
		AssignedUserId: t.AssignedUserID,
		DueDate:        timestamppb.New(t.DueDate),
		LastUpdated:    timestamppb.New(t.LastUpdated),
	}
}
