package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"

	"github.com/ojash123/tasksync/internal/sync"
	"github.com/ojash123/tasksync/internal/task"
	"github.com/ojash123/tasksync/pkg/proto"
)

func main() {
	// --- Define our peer list ---
	// For this test, we are the only peer. It will try to sync with itself.
	peerAddresses := []string{"localhost:9090"}

	// --- Common setup: Create the shared task store ---
	taskStore := task.NewTaskStore()

	// --- Start gRPC Server (in a goroutine) ---
	go startGRPCServer(taskStore)

	// --- Start REST API Server (this will block) ---
	// Pass the peer list to the task handler.
	taskHandler := task.NewTaskHandler(taskStore, peerAddresses)
	router := http.NewServeMux()
	taskHandler.RegisterRoutes(router)

	fmt.Println("Starting REST API server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("could not start REST API server: %v\n", err)
	}
}

func startGRPCServer(ts *task.TaskStore) {
	// 1. Create a TCP listener on port 9090
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("failed to listen on port 9090: %v", err)
	}

	// 2. Create the gRPC server instance
	s := grpc.NewServer()

	// 3. Create our sync server implementation and register it
	syncServer := &sync.Server{Store: ts}
	proto.RegisterTaskSyncServiceServer(s, syncServer)

	fmt.Println("Starting gRPC server on :9090")
	// 4. Start the server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to start gRPC server: %v", err)
	}
}
