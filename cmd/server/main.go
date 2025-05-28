package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"packing-service/internal/packing"
)

// CalculateRequest represents the JSON request structure for item amount.
type CalculateRequest struct {
	Amount int `json:"amount"`
}

// CalculateResponse represents the JSON response structure with calculated packs.
type CalculateResponse struct {
	Packs map[int]int `json:"packs"`
}

// PackSizesRequest represents the JSON request to update available pack sizes.
type PackSizesRequest struct {
	Sizes []int `json:"sizes"`
}

func main() {
	// Serve static files from the ./ui directory as the frontend
	http.Handle("/", http.FileServer(http.Dir("./ui")))

	// Handle POST /calculate to compute optimal pack combination
	http.HandleFunc("/calculate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse the incoming JSON request
		var req CalculateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Load available pack sizes from config file
		f, err := os.Open("config/pack_sizes.json")
		if err != nil {
			http.Error(w, "Cannot read config", http.StatusInternalServerError)
			return
		}
		defer f.Close()

		var sizes []int
		if err := json.NewDecoder(f).Decode(&sizes); err != nil {
			http.Error(w, "Invalid config format", http.StatusInternalServerError)
			return
		}

		// Call the packing algorithm
		result := packing.CalculatePacks(sizes, req.Amount)

		// Send JSON response with the result
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CalculateResponse{Packs: result})
	})

	// Handle GET and POST for /pack-sizes endpoint
	http.HandleFunc("/pack-sizes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Serve current pack sizes from config
			f, err := os.Open("config/pack_sizes.json")
			if err != nil {
				http.Error(w, "Cannot read config", http.StatusInternalServerError)
				return
			}
			info, err := f.Stat()
			if err != nil {
				http.Error(w, "Unable to stat file", http.StatusInternalServerError)
				return
			}
			defer f.Close()
			http.ServeContent(w, r, "pack_sizes.json", info.ModTime(), f)

		case http.MethodPost:
			// Update the pack sizes configuration
			var req PackSizesRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}

			f, err := os.Create("config/pack_sizes.json")
			if err != nil {
				http.Error(w, "Cannot write config", http.StatusInternalServerError)
				return
			}
			defer f.Close()

			json.NewEncoder(f).Encode(req.Sizes)
			w.WriteHeader(http.StatusOK)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Start the HTTP server
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
