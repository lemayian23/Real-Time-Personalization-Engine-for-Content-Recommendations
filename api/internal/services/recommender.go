package services

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"

	"recommendation-engine/api/internal/cache"
	"recommendation-engine/api/internal/database"
)

type Recommender struct {
	db    *database.DB
	cache *cache.RedisCache
}

func NewRecommender(db *database.DB, cache *cache.RedisCache) *Recommender {
	return &Recommender{
		db:    db,
		cache: cache,
	}
}

type Recommendation struct {
	ItemID      string  `json:"item_id"`
	Score       float64 `json:"score"`
	Explanation string  `json:"explanation"`
	Strategy    string  `json:"strategy"`
}

func (r *Recommender) GetRecommendations(userID string, count int) ([]Recommendation, string, error) {
	// Try cache first
	if cached, err := r.cache.GetUserRecommendations(userID); err == nil && cached != nil {
		log.Printf("Cache hit for user: %s", userID)
		var recs []Recommendation
		for _, item := range cached {
			recs = append(recs, Recommendation{
				ItemID:      item["item_id"].(string),
				Score:       item["score"].(float64),
				Explanation: item["explanation"].(string),
				Strategy:    "cached",
			})
		}
		return recs, "cached", nil
	}

	// Check if new user (cold start)
	recentViews, err := r.db.GetUserRecentViews(userID, 5)
	if err != nil {
		return nil, "", err
	}

	var recommendations []Recommendation
	var strategy string

	if len(recentViews) < 3 {
		// Cold start - show trending/popular items
		recommendations = r.getTrendingRecommendations(count)
		strategy = "trending"
	} else {
		// Personalized recommendations
		recommendations = r.getPersonalizedRecommendations(userID, recentViews, count)
		strategy = "personalized"
	}

	// Cache the recommendations
	if err := r.cache.SetUserRecommendations(userID, recommendations, 5*time.Minute); err != nil {
		log.Printf("Warning: failed to cache recommendations: %v", err)
	}

	return recommendations, strategy, nil
}

func (r *Recommender) getTrendingRecommendations(count int) []Recommendation {
	// In production, this would query a trending items table
	// For now, return mock trending items
	trendingItems := []struct {
		ID    string
		Title string
	}{
		{"item_tech_1", "Latest AI Breakthroughs"},
		{"item_science_1", "Space Exploration Updates"},
		{"item_business_1", "Market Trends Analysis"},
		{"item_health_1", "Health & Wellness Tips"},
		{"item_entertainment_1", "Popular Entertainment News"},
	}

	var recs []Recommendation
	for i, item := range trendingItems {
		if i >= count {
			break
		}
		recs = append(recs, Recommendation{
			ItemID:      item.ID,
			Score:       0.9 - (float64(i) * 0.1),
			Explanation: "Trending in your network",
			Strategy:    "trending",
		})
	}

	return recs
}

func (r *Recommender) getPersonalizedRecommendations(userID string, recentViews []string, count int) []Recommendation {
	// In production, this would use real ML models
	// For now, return mock personalized recommendations
	personalizedItems := []struct {
		ID          string
		Score       float64
		Explanation string
	}{
		{"item_tech_2", 0.95, "Based on your interest in technology"},
		{"item_science_2", 0.88, "Similar to content you viewed"},
		{"item_business_2", 0.82, "Popular among users like you"},
		{"item_health_2", 0.78, "Complementary to your interests"},
		{"item_entertainment_2", 0.75, "Diversifying your content mix"},
	}

	var recs []Recommendation
	for i, item := range personalizedItems {
		if i >= count {
			break
		}
		recs = append(recs, Recommendation{
			ItemID:      item.ID,
			Score:       item.Score,
			Explanation: item.Explanation,
			Strategy:    "personalized",
		})
	}

	return recs
}

func (r *Recommender) TrackUserEvent(userID, itemID, eventType string, duration *int) error {
	// Log to database
	if err := r.db.LogUserEvent(userID, itemID, eventType, duration); err != nil {
		return err
	}

	// Update cache counters
	if eventType == "view" || eventType == "click" {
		if err := r.cache.IncrementUserActivity(userID); err != nil {
			log.Printf("Warning: failed to update user activity: %v", err)
		}
	}

	// Invalidate cached recommendations
	r.cache.SetUserRecommendations(userID, nil, 0)

	log.Printf("Tracked event: %s %s %s", userID, eventType, itemID)
	return nil
}