-- Users table
CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_active TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Content items table
CREATE TABLE content_items (
    id VARCHAR(255) PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    category VARCHAR(100),
    tags JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    embedding_vector FLOAT[] -- For PostgreSQL vector operations
);

-- User events table
CREATE TABLE user_events (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) REFERENCES users(id),
    item_id VARCHAR(255) REFERENCES content_items(id),
    event_type VARCHAR(50) NOT NULL, -- 'view', 'click', 'like', 'share'
    duration_seconds INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- A/B test assignments
CREATE TABLE ab_test_assignments (
    user_id VARCHAR(255) PRIMARY KEY,
    experiment_name VARCHAR(255) NOT NULL,
    variant VARCHAR(50) NOT NULL, -- 'control', 'treatment'
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Recommendations served (for analytics)
CREATE TABLE recommendations_served (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) REFERENCES users(id),
    recommended_items JSONB NOT NULL, -- Array of item IDs
    strategy VARCHAR(100) NOT NULL,
    ab_test_variant VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX idx_user_events_user_id ON user_events(user_id);
CREATE INDEX idx_user_events_item_id ON user_events(item_id);
CREATE INDEX idx_user_events_created_at ON user_events(created_at);
CREATE INDEX idx_ab_test_assignments_experiment ON ab_test_assignments(experiment_name);