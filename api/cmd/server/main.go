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

// Add these structs with other type definitions
type UserItemInteraction struct {
	UserID    string   `json:"user_id"`
	ItemIDs   []string `json:"item_ids"`
	Interactions int    `json:"interactions"`
}

type SimilarityScore struct {
	User1    string  `json:"user1"`
	User2    string  `json:"user2"`
	Score    float64 `json:"score"`
	SharedItems []string `json:"shared_items"`
}

type AlgorithmState struct {
	CurrentUser  string              `json:"current_user"`
	SimilarUsers []SimilarityScore   `json:"similar_users"`
	Interactions []UserItemInteraction `json:"interactions"`
	Recommendations []string         `json:"recommendations"`
	Explanation  string              `json:"explanation"`
}

// Add these global variables
var (
	userInteractions = make(map[string][]string)
	algorithmStates  = make(map[string]*AlgorithmState)
)

// Add this endpoint to main()
http.HandleFunc("/algorithm-visualization", algorithmVisualizationHandler)
http.HandleFunc("/algorithm-visualization/", algorithmVisualizationDetailHandler)

// NEW: Algorithm visualization endpoint
func algorithmVisualizationHandler(w http.ResponseWriter, r *http.Request) {
	// Generate mock algorithm state if empty
	if len(algorithmStates) == 0 {
		generateMockAlgorithmData()
	}

	// Get overview of all users and their algorithm states
	overview := make([]map[string]interface{}, 0)
	for userID, state := range algorithmStates {
		overview = append(overview, map[string]interface{}{
			"user_id": userID,
			"interaction_count": len(userInteractions[userID]),
			"similar_users_count": len(state.SimilarUsers),
			"recommendation_count": len(state.Recommendations),
			"last_updated": time.Now().Add(-time.Duration(rand.Intn(60)) * time.Second),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(overview)
}

// NEW: Algorithm visualization detail endpoint
func algorithmVisualizationDetailHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, `{"error": "User ID required"}`, http.StatusBadRequest)
		return
	}
	
	userID := pathParts[2]
	state, exists := algorithmStates[userID]
	if !exists {
		// Create new algorithm state for this user
		state = createAlgorithmState(userID)
		algorithmStates[userID] = state
	}

	// Update algorithm state with fresh data
	updateAlgorithmState(state)

	response := map[string]interface{}{
		"algorithm_state": state,
		"user_interactions": getUserInteractionsMatrix(),
		"similarity_matrix": getSimilarityMatrix(),
		"performance": map[string]interface{}{
			"calculation_time_ms": rand.Intn(50) + 10,
			"similarity_threshold": 0.3,
			"min_shared_items": 2,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func generateMockAlgorithmData() {
	users := []string{"alice", "bob", "charlie", "diana", "eve"}
	items := []string{"tech_ai_news", "science_space", "business_trends", "health_wellness", "entertainment_pop"}
	
	for _, user := range users {
		// Generate random interactions
		interactionCount := 3 + rand.Intn(5)
		for i := 0; i < interactionCount; i++ {
			item := items[rand.Intn(len(items))]
			if !contains(userInteractions[user], item) {
				userInteractions[user] = append(userInteractions[user], item)
			}
		}
		
		// Create algorithm state
		algorithmStates[user] = createAlgorithmState(user)
	}
}

func createAlgorithmState(userID string) *AlgorithmState {
	state := &AlgorithmState{
		CurrentUser: userID,
		Explanation: "Analyzing user behavior patterns...",
	}
	updateAlgorithmState(state)
	return state
}

func updateAlgorithmState(state *AlgorithmState) {
	// Calculate similar users
	state.SimilarUsers = calculateSimilarUsers(state.CurrentUser)
	
	// Generate recommendations based on similar users
	state.Recommendations = generateRecommendations(state.CurrentUser, state.SimilarUsers)
	
	// Update explanation
	state.Explanation = generateExplanation(state.CurrentUser, state.SimilarUsers, state.Recommendations)
	
	// Update interactions visualization
	state.Interactions = getUserInteractionsForState(state.CurrentUser)
}

func calculateSimilarUsers(userID string) []SimilarityScore {
	var similarities []SimilarityScore
	currentUserItems := userInteractions[userID]
	
	for otherUser, otherItems := range userInteractions {
		if otherUser == userID {
			continue
		}
		
		// Calculate Jaccard similarity
		shared := intersection(currentUserItems, otherItems)
		similarity := float64(len(shared)) / float64(len(union(currentUserItems, otherItems)))
		
		if similarity > 0.1 { // Only include meaningful similarities
			similarities = append(similarities, SimilarityScore{
				User1: userID,
				User2: otherUser,
				Score: similarity,
				SharedItems: shared,
			})
		}
	}
	
	// Sort by similarity score
	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i].Score > similarities[j].Score
	})
	
	// Return top 3
	if len(similarities) > 3 {
		return similarities[:3]
	}
	return similarities
}

func generateRecommendations(userID string, similarUsers []SimilarityScore) []string {
	userItems := userInteractions[userID]
	var recommendations []string
	
	for _, similarUser := range similarUsers {
		otherItems := userInteractions[similarUser.User2]
		for _, item := range otherItems {
			if !contains(userItems, item) && !contains(recommendations, item) {
				recommendations = append(recommendations, item)
			}
		}
	}
	
	return recommendations
}

func generateExplanation(userID string, similarUsers []SimilarityScore, recommendations []string) string {
	if len(similarUsers) == 0 {
		return "Not enough data for collaborative filtering. Using content-based recommendations."
	}
	
	explanation := fmt.Sprintf("Found %d similar users. ", len(similarUsers))
	
	for i, similarUser := range similarUsers {
		if i < 2 { // Limit to top 2 for brevity
			explanation += fmt.Sprintf("User %s (%.0f%% similar) shared %d interests. ", 
				similarUser.User2, similarUser.Score*100, len(similarUser.SharedItems))
		}
	}
	
	explanation += fmt.Sprintf("Based on these patterns, recommending %d new items.", len(recommendations))
	return explanation
}

func getUserInteractionsForState(userID string) []UserItemInteraction {
	var interactions []UserItemInteraction
	
	for user, items := range userInteractions {
		interactions = append(interactions, UserItemInteraction{
			UserID: user,
			ItemIDs: items,
			Interactions: len(items),
		})
	}
	
	return interactions
}

func getUserInteractionsMatrix() map[string][]string {
	return userInteractions
}

func getSimilarityMatrix() map[string]map[string]float64 {
	matrix := make(map[string]map[string]float64)
	
	for user1 := range userInteractions {
		matrix[user1] = make(map[string]float64)
		for user2 := range userInteractions {
			if user1 == user2 {
				matrix[user1][user2] = 1.0
			} else {
				shared := intersection(userInteractions[user1], userInteractions[user2])
				similarity := float64(len(shared)) / float64(len(union(userInteractions[user1], userInteractions[user2])))
				matrix[user1][user2] = similarity
			}
		}
	}
	
	return matrix
}

// Helper functions
func intersection(a, b []string) []string {
	set := make(map[string]bool)
	for _, item := range a {
		set[item] = true
	}
	
	var result []string
	for _, item := range b {
		if set[item] {
			result = append(result, item)
		}
	}
	return result
}

func union(a, b []string) []string {
	set := make(map[string]bool)
	for _, item := range a {
		set[item] = true
	}
	for _, item := range b {
		set[item] = true
	}
	
	var result []string
	for item := range set {
		result = append(result, item)
	}
	return result
}

// User Session structure
type UserSession struct {
	UserID      string    `json:"user_id"`
	SessionID   string    `json:"session_id"`
	StartTime   time.Time `json:"start_time"`
	LastActive  time.Time `json:"last_active"`
	PageViews   int       `json:"page_views"`
	Clicks      int       `json:"clicks"`
	SessionTime int       `json:"session_time"` // in seconds
	Categories  []string  `json:"categories"`
}

// Global variables for live metrics
var (
	totalImpressions = 10000
	totalClicks      = 2500
	activeUsers      = 342
	systemStartTime  = time.Now()
	engagementStats  = []map[string]interface{}{
		{"timestamp": time.Now().Add(-5 * time.Minute), "ctr": 0.23, "users": 45},
		{"timestamp": time.Now().Add(-4 * time.Minute), "ctr": 0.26, "users": 52},
		{"timestamp": time.Now().Add(-3 * time.Minute), "ctr": 0.28, "users": 61},
		{"timestamp": time.Now().Add(-2 * time.Minute), "ctr": 0.25, "users": 58},
		{"timestamp": time.Now().Add(-1 * time.Minute), "ctr": 0.27, "users": 55},
	}
	userSessions = make(map[string]*UserSession)
)

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
	http.HandleFunc("/ab-tests", abTestsHandler)
	http.HandleFunc("/ab-tests/", abTestDetailHandler)
	http.HandleFunc("/engagement-metrics", engagementMetricsHandler)
	http.HandleFunc("/user-sessions", userSessionsHandler) // NEW ENDPOINT
	http.HandleFunc("/user-sessions/", userSessionDetailHandler) // NEW ENDPOINT

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// NEW: User sessions list endpoint
func userSessionsHandler(w http.ResponseWriter, r *http.Request) {
	// Generate some mock user sessions if empty
	if len(userSessions) == 0 {
		generateMockSessions()
	}

	sessionsList := make([]map[string]interface{}, 0)
	for _, session := range userSessions {
		// Update session time for realism
		session.SessionTime = int(time.Since(session.StartTime).Seconds())
		
		sessionsList = append(sessionsList, map[string]interface{}{
			"user_id":      session.UserID,
			"session_id":   session.SessionID,
			"start_time":   session.StartTime,
			"session_time": session.SessionTime,
			"page_views":   session.PageViews,
			"clicks":       session.Clicks,
			"engagement":   float64(session.Clicks) / float64(session.PageViews),
			"categories":   session.Categories,
			"status":       "active",
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessionsList)
}

// NEW: User session detail endpoint
func userSessionDetailHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, `{"error": "User ID required"}`, http.StatusBadRequest)
		return
	}
	
	userID := pathParts[2]
	session, exists := userSessions[userID]
	if !exists {
		// Create a new session for this user
		session = &UserSession{
			UserID:     userID,
			SessionID:  "session_" + userID + "_" + strconv.FormatInt(time.Now().Unix(), 10),
			StartTime:  time.Now(),
			LastActive: time.Now(),
			PageViews:  1,
			Clicks:     0,
			Categories: []string{"general"},
		}
		userSessions[userID] = session
	}

	// Update session activity
	session.LastActive = time.Now()
	session.SessionTime = int(time.Since(session.StartTime).Seconds())

	// Generate mock click stream
	clickStream := generateClickStream(session)

	response := map[string]interface{}{
		"session": session,
		"analytics": map[string]interface{}{
			"click_stream":    clickStream,
			"avg_time_per_click": session.SessionTime / max(session.Clicks, 1),
			"category_distribution": getCategoryDistribution(session),
			"engagement_score": calculateEngagementScore(session),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func generateMockSessions() {
	users := []string{"alice", "bob", "charlie", "diana", "eve", "frank", "grace", "henry"}
	categories := [][]string{
		{"technology", "science"},
		{"business", "finance"},
		{"entertainment", "lifestyle"},
		{"sports", "health"},
		{"technology", "business"},
		{"science", "health"},
		{"entertainment", "technology"},
		{"finance", "lifestyle"},
	}

	for i, user := range users {
		userSessions[user] = &UserSession{
			UserID:     user,
			SessionID:  "session_" + user + "_" + strconv.FormatInt(time.Now().Unix(), 10),
			StartTime:  time.Now().Add(-time.Duration(rand.Intn(60)) * time.Minute),
			LastActive: time.Now().Add(-time.Duration(rand.Intn(10)) * time.Minute),
			PageViews:  5 + rand.Intn(20),
			Clicks:     2 + rand.Intn(10),
			Categories: categories[i],
		}
	}
}

func generateClickStream(session *UserSession) []map[string]interface{} {
	stream := []map[string]interface{}{
		{
			"timestamp": session.StartTime.Format(time.RFC3339),
			"action":    "session_start",
			"item_id":   "homepage",
			"duration":  0,
		},
	}

	// Generate mock click events
	currentTime := session.StartTime
	for i := 0; i < session.Clicks; i++ {
		currentTime = currentTime.Add(time.Duration(10+rand.Intn(30)) * time.Second)
		stream = append(stream, map[string]interface{}{
			"timestamp": currentTime.Format(time.RFC3339),
			"action":    "click",
			"item_id":   "item_" + session.Categories[rand.Intn(len(session.Categories))] + "_" + strconv.Itoa(i),
			"duration":  10 + rand.Intn(50),
		})
	}

	return stream
}

func getCategoryDistribution(session *UserSession) map[string]int {
	distribution := make(map[string]int)
	for _, category := range session.Categories {
		distribution[category] = 5 + rand.Intn(10)
	}
	return distribution
}

func calculateEngagementScore(session *UserSession) float64 {
	if session.PageViews == 0 {
		return 0.0
	}
	clickRatio := float64(session.Clicks) / float64(session.PageViews)
	timeRatio := float64(session.SessionTime) / float64(session.Clicks*10)
	return (clickRatio*0.6 + timeRatio*0.4) * 100
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Engagement metrics endpoint
func engagementMetricsHandler(w http.ResponseWriter, r *http.Request) {
	// Update metrics with some random variation
	updateLiveMetrics()
	
	currentCTR := float64(totalClicks) / float64(totalImpressions)
	
	metrics := map[string]interface{}{
		"timestamp":         time.Now().Format(time.RFC3339),
		"total_impressions": totalImpressions,
		"total_clicks":      totalClicks,
		"current_ctr":       currentCTR,
		"active_users":      activeUsers,
		"uptime_minutes":    time.Since(systemStartTime).Minutes(),
		"engagement_trend":  engagementStats,
		"performance": map[string]interface{}{
			"p95_latency_ms": 24,
			"error_rate":      0.0023,
			"throughput_rps":  1250,
		},
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

func updateLiveMetrics() {
	// Simulate live data changes
	totalImpressions += rand.Intn(50)
	totalClicks += rand.Intn(20)
	activeUsers = 300 + rand.Intn(100)
	
	// Update engagement trend
	newCTR := 0.22 + rand.Float64()*0.1
	engagementStats = append(engagementStats[1:], map[string]interface{}{
		"timestamp": time.Now(),
		"ctr":       newCTR,
		"users":     activeUsers,
	})
}

// A/B Tests list endpoint
func abTestsHandler(w http.ResponseWriter, r *http.Request) {
	// Update CTRs with some random variation to simulate live data
	updateExperimentData()
	
	experimentList := make([]map[string]interface{}, 0)
	for id, exp := range experiments {
		experimentList = append(experimentList, map[string]interface{}{
			"id":                id,
			"name":              exp.Name,
			"start_time":        exp.StartTime,
			"total_impressions": getTotalImpressions(exp),
			"winning_variant":   getWinningVariant(exp),
			"status":            "running",
		})
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(experimentList)
}

// A/B Test detail endpoint
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
			"duration":        time.Since(exp.StartTime).String(),
			"total_users":     getTotalImpressions(exp),
			"overall_ctr":     getOverallCTR(exp),
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
			variation := rand.Float64()*0.02 - 0.01 // ±1% variation
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
		"user_id":         userID,
		"recommendations": recommendations,
		"latency_ms":      time.Since(start).Milliseconds(),
		"strategy":        strategy,
		"timestamp":       time.Now().Format(time.RFC3339),
		"version":         "simple-v1",
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

	// Track user session activity
	trackUserSession(event.UserID, event.EventType, event.ItemID)

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

// NEW: Track user session activity
func trackUserSession(userID, eventType, itemID string) {
	session, exists := userSessions[userID]
	if !exists {
		session = &UserSession{
			UserID:     userID,
			SessionID:  "session_" + userID + "_" + strconv.FormatInt(time.Now().Unix(), 10),
			StartTime:  time.Now(),
			LastActive: time.Now(),
			PageViews:  0,
			Clicks:     0,
			Categories: []string{},
		}
		userSessions[userID] = session
	}

	session.LastActive = time.Now()

	switch eventType {
	case "view":
		session.PageViews++
		// Extract category from item ID
		if parts := strings.Split(itemID, "_"); len(parts) > 1 {
			category := parts[0]
			if !contains(session.Categories, category) {
				session.Categories = append(session.Categories, category)
			}
		}
	case "click":
		session.Clicks++
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "simple-v1",
		"features":  []string{"mock-recommendations", "event-logging", "health-check", "diversity-scoring", "ab-testing", "live-metrics", "user-sessions"},
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