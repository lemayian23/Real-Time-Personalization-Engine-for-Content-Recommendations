package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func NewPostgresDB(connectionString string) (*DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("✅ Connected to PostgreSQL")
	return &DB{db}, nil
}

func (db *DB) LogUserEvent(userID, itemID, eventType string, duration *int) error {
	query := `
		INSERT INTO user_events (user_id, item_id, event_type, duration_seconds) 
		VALUES ($1, $2, $3, $4)
	`
	_, err := db.Exec(query, userID, itemID, eventType, duration)
	return err
}

func (db *DB) GetUserRecentViews(userID string, limit int) ([]string, error) {
	query := `
		SELECT item_id FROM user_events 
		WHERE user_id = $1 AND event_type = 'view' 
		ORDER BY created_at DESC 
		LIMIT $2
	`
	rows, err := db.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []string
	for rows.Next() {
		var itemID string
		if err := rows.Scan(&itemID); err != nil {
			return nil, err
		}
		items = append(items, itemID)
	}

	return items, nil
}