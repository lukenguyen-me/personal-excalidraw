package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Personal Excalidraw Backend")
	fmt.Println("Server would start here")
	// TODO: Implement backend server
}

// healthCheckHandler handles health check requests
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok"}`)
}
