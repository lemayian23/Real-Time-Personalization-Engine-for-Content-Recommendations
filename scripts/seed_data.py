import psycopg2
import json
import random
from datetime import datetime, timedelta

def seed_data():
    """Seed database with sample data"""
    conn = psycopg2.connect(
        host="localhost",
        database="recommendations", 
        user="user",
        password="pass",
        port="5432"
    )
    
    cursor = conn.cursor()
    
    # Sample content categories
    categories = ['technology', 'science', 'business', 'health', 'entertainment', 'sports']
    
    try:
        # Seed content items
        for i in range(1, 101):
            cursor.execute("""
                INSERT INTO content_items (id, title, description, category, tags)
                VALUES (%s, %s, %s, %s, %s)
                ON CONFLICT (id) DO NOTHING
            """, (
                f"item_{i}",
                f"Sample Article {i}",
                f"This is a sample description for article {i}",
                random.choice(categories),
                json.dumps([f"tag_{j}" for j in range(1, random.randint(2, 5))])
            ))
        
        # Seed some users
        for i in range(1, 51):
            cursor.execute("""
                INSERT INTO users (id)
                VALUES (%s)
                ON CONFLICT (id) DO NOTHING
            """, (f"user_{i}",))
        
        # Seed some user events
        for i in range(1, 1001):
            user_id = f"user_{random.randint(1, 50)}"
            item_id = f"item_{random.randint(1, 100)}"
            event_type = random.choice(['view', 'click', 'like'])
            
            cursor.execute("""
                INSERT INTO user_events (user_id, item_id, event_type, duration_seconds)
                VALUES (%s, %s, %s, %s)
            """, (
                user_id,
                item_id,
                event_type,
                random.randint(10, 300) if event_type == 'view' else None
            ))
        
        conn.commit()
        print("✅ Sample data seeded successfully!")
        print("   - 100 content items")
        print("   - 50 users") 
        print("   - 1000 user events")
        
    except Exception as e:
        print(f"❌ Error seeding data: {e}")
        conn.rollback()
    finally:
        cursor.close()
        conn.close()

if __name__ == "__main__":
    seed_data()