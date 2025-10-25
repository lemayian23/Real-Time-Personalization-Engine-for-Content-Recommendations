import psycopg2
import json
from datetime import datetime

class ContentProcessor:
    def __init__(self, db_url):
        self.db_url = db_url
    
    def add_content_item(self, item_id, title, description, category=None, tags=None):
        """Add a new content item to the database"""
        conn = psycopg2.connect(self.db_url)
        cursor = conn.cursor()
        
        try:
            cursor.execute("""
                INSERT INTO content_items (id, title, description, category, tags)
                VALUES (%s, %s, %s, %s, %s)
                ON CONFLICT (id) DO UPDATE SET
                    title = EXCLUDED.title,
                    description = EXCLUDED.description,
                    category = EXCLUDED.category,
                    tags = EXCLUDED.tags
            """, (item_id, title, description, category, json.dumps(tags) if tags else None))
            
            conn.commit()
            print(f"Added/Updated content item: {item_id}")
            
        except Exception as e:
            conn.rollback()
            print(f"Error adding content item: {e}")
        finally:
            cursor.close()
            conn.close()
    
    def track_user_event(self, user_id, item_id, event_type, duration=None):
        """Track user interaction event"""
        conn = psycopg2.connect(self.db_url)
        cursor = conn.cursor()
        
        try:
            # Ensure user exists
            cursor.execute("""
                INSERT INTO users (id, last_active)
                VALUES (%s, %s)
                ON CONFLICT (id) DO UPDATE SET last_active = EXCLUDED.last_active
            """, (user_id, datetime.now()))
            
            # Record event
            cursor.execute("""
                INSERT INTO user_events (user_id, item_id, event_type, duration_seconds)
                VALUES (%s, %s, %s, %s)
            """, (user_id, item_id, event_type, duration))
            
            conn.commit()
            print(f"Tracked event: {user_id} {event_type} {item_id}")
            
        except Exception as e:
            conn.rollback()
            print(f"Error tracking event: {e}")
        finally:
            cursor.close()
            conn.close()