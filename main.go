package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gobot/agent"
)

// QueryRequest represents the structure of the incoming query request
type QueryRequest struct {
	Query string `json:"query"`
}

// QueryResponse represents the structure of the response
type QueryResponse struct {
	Answer string `json:"answer"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func main() {
	// Create AI agent
	aiAgent := agent.NewAgent("")
	defer aiAgent.Close()

	// Set up HTTP routes
	http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		// Enable CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Only allow POST requests
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse the request body
		var req QueryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Process the query
		result, err := aiAgent.ProcessQuery(req.Query)
		if err != nil {
			response := QueryResponse{
				Status: "error",
				Error:  err.Error(),
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		// Send the response
		response := QueryResponse{
			Answer: result.Answer,
			Status: result.Status,
		}
		if result.Error != nil {
			response.Error = result.Error.Error()
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Start the server
	fmt.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
