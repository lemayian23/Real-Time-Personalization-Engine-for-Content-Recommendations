package cache

import (
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(redisURL string) (*RedisCache, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)
	ctx := context.Background()

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	log.Println("✅ Connected to Redis")
	return &RedisCache{client: client, ctx: ctx}, nil
}

func (r *RedisCache) SetUserRecommendations(userID string, recommendations interface{}, ttl time.Duration) error {
	jsonData, err := json.Marshal(recommendations)
	if err != nil {
		return err
	}

	return r.client.Set(r.ctx, "recs:"+userID, jsonData, ttl).Err()
}

func (r *RedisCache) GetUserRecommendations(userID string) ([]map[string]interface{}, error) {
	val, err := r.client.Get(r.ctx, "recs:"+userID).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	} else if err != nil {
		return nil, err
	}

	var recommendations []map[string]interface{}
	err = json.Unmarshal([]byte(val), &recommendations)
	return recommendations, err
}

func (r *RedisCache) IncrementUserActivity(userID string) error {
	key := "user_activity:" + userID
	return r.client.Incr(r.ctx, key).Err()
}