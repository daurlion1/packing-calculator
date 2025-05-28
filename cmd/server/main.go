package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"packing-service/internal/packing"
)

type CalculateRequest struct {
	Amount int `json:"amount"`
}

type CalculateResponse struct {
	Packs map[int]int `json:"packs"`
}

type PackSizesRequest struct {
	Sizes []int `json:"sizes"`
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./ui")))

	http.HandleFunc("/calculate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req CalculateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

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

		result := packing.CalculatePacks(sizes, req.Amount)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CalculateResponse{Packs: result})
	})

	http.HandleFunc("/pack-sizes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
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

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
