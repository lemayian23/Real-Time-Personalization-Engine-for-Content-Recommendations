package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	// For now, use simple version - we'll add database later
	log.Println("🚀 Starting SIMPLE recommendation API on :8080")
	log.Println("📝 Note: Using mock data - database integration pending")
	
	http.HandleFunc("/recommend", recommendHandler)
	http.HandleFunc("/event", eventHandler) 
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/metrics", metricsHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func calculateDiversityScore(recommendations []map[string]interface{}) float64 {
    if len(recommendations) == 0 {
        return 0.0
    }
    
    // Extract categories from item IDs (simple heuristic)
    categories := make(map[string]bool)
    for _, rec := range recommendations {
        itemID := rec["item_id"].(string)
        // Simple category extraction from item_id format: "category_content"
        if parts := strings.Split(itemID, "_"); len(parts) > 1 {
            categories[parts[0]] = true
        }
    }
    
    // Calculate diversity as ratio of unique categories to total recommendations
    uniqueCategories := len(categories)
    totalItems := len(recommendations)
    
    return float64(uniqueCategories) / float64(totalItems)
}

// Update the recommendHandler response to include diversity score:
response := map[string]interface{}{
    "user_id":        userID,
    "recommendations": recommendations,
    "latency_ms":     time.Since(start).Milliseconds(),
    "strategy":       strategy,
    "timestamp":      time.Now().Format(time.RFC3339),
    "version":       "simple-v1",
    "diversity_score": calculateDiversityScore(recommendations), // NEW
}

func recommendHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		userID = "anonymous"
	}

	countStr := r.URL.Query().Get("count")
	count := 10
	if countStr != "" {
		if parsed, err := strconv.Atoi(countStr); err == nil && parsed > 0 {
			count = parsed
		}
	}

	// SIMPLE VERSION - Mock recommendations
	recommendations, strategy := getMockRecommendations(userID, count)

	response := map[string]interface{}{
		"user_id":        userID,
		"recommendations": recommendations,
		"latency_ms":     time.Since(start).Milliseconds(),
		"strategy":       strategy,
		"timestamp":      time.Now().Format(time.RFC3339),
		"version":       "simple-v1",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getMockRecommendations(userID string, count int) ([]map[string]interface{}, string) {
	// Simple logic: if user has specific pattern, return personalized, else trending
	var recommendations []map[string]interface{}
	strategy := "trending"
	
	if len(userID) > 5 {
		strategy = "personalized"
		// Mock personalized recommendations
		personalizedItems := []map[string]interface{}{
			{"item_id": "tech_ai_news", "score": 0.95, "explanation": "Based on your tech interests"},
			{"item_id": "science_space", "score": 0.88, "explanation": "Similar to content you viewed"},
			{"item_id": "business_trends", "score": 0.82, "explanation": "Popular in your network"},
			{"item_id": "health_wellness", "score": 0.78, "explanation": "Complementary content"},
			{"item_id": "entertainment_pop", "score": 0.75, "explanation": "Trending now"},
		}
		recommendations = personalizedItems[:min(count, len(personalizedItems))]
	} else {
		// Mock trending recommendations
		trendingItems := []map[string]interface{}{
			{"item_id": "breaking_news_1", "score": 0.92, "explanation": "🔥 Trending worldwide"},
			{"item_id": "viral_video_1", "score": 0.89, "explanation": "📈 Going viral"},
			{"item_id": "popular_tutorial", "score": 0.85, "explanation": "⭐ Most watched today"},
			{"item_id": "celebrity_news", "score": 0.81, "explanation": "🌟 Top story"},
			{"item_id": "sports_highlight", "score": 0.79, "explanation": "🏆 Match of the day"},
		}
		recommendations = trendingItems[:min(count, len(trendingItems))]
	}
	
	return recommendations, strategy
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	var event struct {
		UserID    string `json:"user_id"`
		ItemID    string `json:"item_id"`
		EventType string `json:"event_type"`
		Duration  *int   `json:"duration_seconds"`
	}

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, `{"error": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if event.UserID == "" || event.ItemID == "" || event.EventType == "" {
		http.Error(w, `{"error": "user_id, item_id, and event_type are required"}`, http.StatusBadRequest)
		return
	}

	// SIMPLE VERSION - Just log the event
	log.Printf("📊 EVENT: user=%s item=%s type=%s duration=%v", 
		event.UserID, event.ItemID, event.EventType, event.Duration)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "recorded",
		"user_id": event.UserID,
		"note":    "mock-storage",
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "simple-v1",
		"features":  []string{"mock-recommendations", "event-logging", "health-check"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	metrics := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"system":    "recommendation-engine",
		"version":   "simple-v1",
		"status":    "operational",
		"uptime":    "since-last-deploy",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}