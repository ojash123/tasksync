package main

import (
	"fmt"
	"log"
	"net/http"

	// Import the internal task package
	"github.com/ojash123/tasksync/internal/task"
)

func main() {
	// 1. Create the store
	taskStore := task.NewTaskStore()

	// 2. Create the handler and inject the store
	taskHandler := task.NewTaskHandler(taskStore)

	// 3. Create the router and register the routes
	router := http.NewServeMux()
	taskHandler.RegisterRoutes(router) // Use our new RegisterRoutes function

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("could not start server: %v\n", err)
	}
}
