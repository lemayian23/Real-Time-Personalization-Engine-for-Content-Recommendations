import time
import psycopg2
from datetime import datetime

class MetricsCollector:
    def __init__(self, db_url):
        self.db_url = db_url
    
    def track_recommendation_latency(self, user_id, latency_ms, strategy):
        """Track recommendation latency metrics"""
        conn = psycopg2.connect(self.db_url)
        cursor = conn.cursor()
        
        try:
            cursor.execute("""
                INSERT INTO recommendation_metrics 
                (user_id, metric_type, metric_value, strategy, created_at)
                VALUES (%s, %s, %s, %s, %s)
            """, (user_id, 'latency_ms', latency_ms, strategy, datetime.now()))
            
            conn.commit()
        except Exception as e:
            print(f"Error tracking metrics: {e}")
            conn.rollback()
        finally:
            cursor.close()
            conn.close()
    
    def track_ctr(self, user_id, recommended_items, clicked_items):
        """Track click-through rate metrics"""
        conn = psycopg2.connect(self.db_url)
        cursor = conn.cursor()
        
        try:
            ctr = len(clicked_items) / len(recommended_items) if recommended_items else 0
            
            cursor.execute("""
                INSERT INTO recommendation_metrics 
                (user_id, metric_type, metric_value, strategy, created_at)
                VALUES (%s, %s, %s, %s, %s)
            """, (user_id, 'ctr', ctr, 'hybrid', datetime.now()))
            
            conn.commit()
        except Exception as e:
            print(f"Error tracking CTR: {e}")
            conn.rollback()
        finally:
            cursor.close()
            conn.close()
    
    def get_system_health(self):
        """Get system health metrics"""
        conn = psycopg2.connect(self.db_url)
        cursor = conn.cursor()
        
        try:
            # Get recent latency
            cursor.execute("""
                SELECT AVG(metric_value) as avg_latency,
                       PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY metric_value) as p95_latency
                FROM recommendation_metrics
                WHERE metric_type = 'latency_ms' 
                AND created_at >= NOW() - INTERVAL '1 hour'
            """)
            latency_stats = cursor.fetchone()
            
            # Get current CTR
            cursor.execute("""
                SELECT AVG(metric_value) as avg_ctr
                FROM recommendation_metrics
                WHERE metric_type = 'ctr'
                AND created_at >= NOW() - INTERVAL '1 hour'
            """)
            ctr_stats = cursor.fetchone()
            
            return {
                'avg_latency_ms': latency_stats[0] or 0,
                'p95_latency_ms': latency_stats[1] or 0,
                'avg_ctr': ctr_stats[0] or 0,
                'timestamp': datetime.now().isoformat()
            }
            
        except Exception as e:
            print(f"Error getting health metrics: {e}")
            return {}
        finally:
            cursor.close()
            conn.close()

# Create metrics table
def create_metrics_table(db_url):
    conn = psycopg2.connect(db_url)
    cursor = conn.cursor()
    
    cursor.execute("""
        CREATE TABLE IF NOT EXISTS recommendation_metrics (
            id SERIAL PRIMARY KEY,
            user_id VARCHAR(255),
            metric_type VARCHAR(50) NOT NULL,
            metric_value FLOAT NOT NULL,
            strategy VARCHAR(100),
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
        
        CREATE INDEX IF NOT EXISTS idx_metrics_type_created ON recommendation_metrics(metric_type, created_at);
    """)
    
    conn.commit()
    cursor.close()
    conn.close()

if __name__ == "__main__":
    create_metrics_table("postgresql://user:pass@localhost:5432/recommendations")
    print("✅ Metrics table created")