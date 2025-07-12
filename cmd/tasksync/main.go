package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// This will eventually hold our API handlers
	router := http.NewServeMux()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Server is running!")
	})

	fmt.Println("Starting server on :8080")
	// Start the HTTP server
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("could not start server: %v\n", err)
	}
}
