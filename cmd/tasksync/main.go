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
	peerAddresses := []string{"localhost:9090", "localhost:9999"}
	taskStore := task.NewTaskStore()

	// 1. Create the channel for broadcasting tasks.
	// A buffered channel is used so that sending doesn't block if the listener is busy.
	taskBroadcastChan := make(chan task.Task, 10)

	// Start gRPC Server (in a goroutine)
	go startGRPCServer(taskStore)

	// 2. Start a new goroutine to listen on the channel and broadcast tasks.
	go func() {
		for t := range taskBroadcastChan {
			sync.BroadcastTask(peerAddresses, t)
		}
	}()

	// 3. Update the NewTaskHandler call to pass the channel.
	taskHandler := task.NewTaskHandler(taskStore, taskBroadcastChan)
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
