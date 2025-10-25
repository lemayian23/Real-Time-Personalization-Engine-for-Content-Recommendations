import psycopg2
import os

def setup_database():
    """Initialize database with schema"""
    conn = psycopg2.connect(
        host="localhost",
        database="recommendations",
        user="user",
        password="pass",
        port="5432"
    )
    
    cursor = conn.cursor()
    
    # Read and execute schema file
    with open('data/migrations/001_initial_schema.sql', 'r') as f:
        schema_sql = f.read()
    
    try:
        cursor.execute(schema_sql)
        conn.commit()
        print("✅ Database schema created successfully!")
    except Exception as e:
        print(f"❌ Error creating schema: {e}")
        conn.rollback()
    finally:
        cursor.close()
        conn.close()

if __name__ == "__main__":
    setup_database()