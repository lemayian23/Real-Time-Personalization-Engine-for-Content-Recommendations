package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/recommend", recommendHandler)
	http.HandleFunc("/event", eventHandler)
	http.HandleFunc("/health", healthHandler)
	
	log.Println("Starting recommendation API on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func recommendHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		userID = "anonymous"
	}
	
	// Simple mock recommendations
	recommendations := []map[string]interface{}{
		{
			"item_id": "1",
			"score":   0.95,
			"explanation": "Popular in your network",
		},
		{
			"item_id": "2", 
			"score":   0.87,
			"explanation": "Similar to your interests",
		},
		{
			"item_id": "3",
			"score":   0.76,
			"explanation": "Trending now",
		},
	}
	
	response := map[string]interface{}{
		"user_id":        userID,
		"recommendations": recommendations,
		"latency_ms":     time.Since(start).Milliseconds(),
		"strategy":       "hybrid",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	var event struct {
		UserID  string `json:"user_id"`
		ItemID  string `json:"item_id"`
		EventType string `json:"event_type"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	log.Printf("Event: %s %s %s", event.UserID, event.EventType, event.ItemID)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "recorded",
		"user_id": event.UserID,
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}