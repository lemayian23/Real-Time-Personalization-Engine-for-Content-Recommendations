package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Experiment data structure
type Experiment struct {
	Name      string    `json:"name"`
	Variants  []Variant `json:"variants"`
	StartTime time.Time `json:"start_time"`
}

type Variant struct {
	Name          string  `json:"name"`
	Impressions   int     `json:"impressions"`
	Clicks        int     `json:"clicks"`
	CTR           float64 `json:"ctr"`
	Confidence    float64 `json:"confidence"`
	IsSignificant bool    `json:"is_significant"`
}

var experiments = map[string]*Experiment{
	"homepage_algo_v1": {
		Name:      "Homepage Algorithm v1",
		StartTime: time.Now().Add(-24 * time.Hour),
		Variants: []Variant{
			{Name: "Control (Hybrid)", Impressions: 15432, Clicks: 2314, CTR: 0.15},
			{Name: "Treatment (Content-Boosted)", Impressions: 15289, Clicks: 3822, CTR: 0.25},
		},
	},
	"cold_start_v2": {
		Name:      "Cold Start Strategy v2",
		StartTime: time.Now().Add(-12 * time.Hour),
		Variants: []Variant{
			{Name: "Control (Trending)", Impressions: 8231, Clicks: 987, CTR: 0.12},
			{Name: "Treatment (LLM-Powered)", Impressions: 8456, Clicks: 1856, CTR: 0.22},
		},
	},
}

func main() {
	// For now, use simple version - we'll add database later
	log.Println("🚀 Starting SIMPLE recommendation API on :8080")
	log.Println("📝 Note: Using mock data - database integration pending")
	
	http.HandleFunc("/recommend", recommendHandler)
	http.HandleFunc("/event", eventHandler) 
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/metrics", metricsHandler)
	http.HandleFunc("/ab-tests", abTestsHandler) // NEW ENDPOINT
	http.HandleFunc("/ab-tests/", abTestDetailHandler) // NEW ENDPOINT

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// NEW: A/B Tests list endpoint
func abTestsHandler(w http.ResponseWriter, r *http.Request) {
	// Update CTRs with some random variation to simulate live data
	updateExperimentData()
	
	experimentList := make([]map[string]interface{}, 0)
	for id, exp := range experiments {
		experimentList = append(experimentList, map[string]interface{}{
			"id": id,
			"name": exp.Name,
			"start_time": exp.StartTime,
			"total_impressions": getTotalImpressions(exp),
			"winning_variant": getWinningVariant(exp),
			"status": "running",
		})
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(experimentList)
}

// NEW: A/B Test detail endpoint
func abTestDetailHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, `{"error": "Experiment ID required"}`, http.StatusBadRequest)
		return
	}
	
	experimentID := pathParts[2]
	exp, exists := experiments[experimentID]
	if !exists {
		http.Error(w, `{"error": "Experiment not found"}`, http.StatusNotFound)
		return
	}
	
	// Calculate confidence intervals and significance
	calculateStatistics(exp)
	
	response := map[string]interface{}{
		"experiment": exp,
		"summary": map[string]interface{}{
			"duration": time.Since(exp.StartTime).String(),
			"total_users": getTotalImpressions(exp),
			"overall_ctr": getOverallCTR(exp),
			"detected_effect": exp.Variants[1].CTR - exp.Variants[0].CTR,
		},
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Helper functions for A/B testing
func updateExperimentData() {
	for _, exp := range experiments {
		for i := range exp.Variants {
			// Add some random variation to simulate live data
			variation := rand.Float64() * 0.02 - 0.01 // ±1% variation
			exp.Variants[i].CTR += variation
			if exp.Variants[i].CTR < 0.05 {
				exp.Variants[i].CTR = 0.05
			}
			
			// Update impressions and clicks based on CTR
			exp.Variants[i].Impressions += rand.Intn(100)
			exp.Variants[i].Clicks = int(float64(exp.Variants[i].Impressions) * exp.Variants[i].CTR)
		}
	}
}

func getTotalImpressions(exp *Experiment) int {
	total := 0
	for _, v := range exp.Variants {
		total += v.Impressions
	}
	return total
}

func getWinningVariant(exp *Experiment) string {
	if len(exp.Variants) == 0 {
		return ""
	}
	winner := exp.Variants[0]
	for _, v := range exp.Variants {
		if v.CTR > winner.CTR {
			winner = v
		}
	}
	return winner.Name
}

func getOverallCTR(exp *Experiment) float64 {
	totalClicks := 0
	totalImpressions := 0
	for _, v := range exp.Variants {
		totalClicks += v.Clicks
		totalImpressions += v.Impressions
	}
	if totalImpressions == 0 {
		return 0.0
	}
	return float64(totalClicks) / float64(totalImpressions)
}

func calculateStatistics(exp *Experiment) {
	for i := range exp.Variants {
		v := &exp.Variants[i]
		// Simple confidence calculation (in real system, use proper statistical test)
		v.Confidence = 0.85 + rand.Float64()*0.14 // 85-99% confidence
		v.IsSignificant = v.Confidence > 0.95
	}
}

// Existing functions remain the same...
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
	diversityScore := calculateDiversityScore(recommendations)

	response := map[string]interface{}{
		"user_id":        userID,
		"recommendations": recommendations,
		"latency_ms":     time.Since(start).Milliseconds(),
		"strategy":       strategy,
		"timestamp":      time.Now().Format(time.RFC3339),
		"version":       "simple-v1",
		"diversity_score": diversityScore,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
		"features":  []string{"mock-recommendations", "event-logging", "health-check", "diversity-scoring", "ab-testing"},
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