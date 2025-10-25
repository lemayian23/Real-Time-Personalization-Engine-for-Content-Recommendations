import hashlib
import json
from datetime import datetime
import psycopg2

class ExperimentConfig:
    def __init__(self, db_url):
        self.db_url = db_url
        self.experiments = {
            'homepage_recommendations_v1': {
                'control': {
                    'weight': 0.5,
                    'algorithm': 'hybrid_collaborative'
                },
                'treatment': {
                    'weight': 0.5,
                    'algorithm': 'hybrid_content_boosted'
                },
                'start_date': '2024-01-01',
                'metrics': ['ctr', 'session_duration', 'conversion_rate']
            }
        }
    
    def assign_variant(self, user_id, experiment_name):
        """Assign user to A/B test variant using consistent hashing"""
        # Create consistent hash based on user_id + experiment_name
        hash_input = f"{user_id}_{experiment_name}".encode()
        hash_value = int(hashlib.md5(hash_input).hexdigest()[:8], 16)
        
        experiment = self.experiments.get(experiment_name)
        if not experiment:
            return 'control'
        
        # Assign based on hash modulo
        hash_mod = hash_value % 100
        control_weight = experiment['control']['weight'] * 100
        
        variant = 'treatment' if hash_mod >= control_weight else 'control'
        
        # Store assignment in database
        self._store_assignment(user_id, experiment_name, variant)
        
        return variant
    
    def _store_assignment(self, user_id, experiment_name, variant):
        """Store A/B test assignment in database"""
        conn = psycopg2.connect(self.db_url)
        cursor = conn.cursor()
        
        try:
            cursor.execute("""
                INSERT INTO ab_test_assignments (user_id, experiment_name, variant)
                VALUES (%s, %s, %s)
                ON CONFLICT (user_id) DO UPDATE SET
                    variant = EXCLUDED.variant,
                    assigned_at = CURRENT_TIMESTAMP
            """, (user_id, experiment_name, variant))
            
            conn.commit()
        except Exception as e:
            print(f"Error storing assignment: {e}")
            conn.rollback()
        finally:
            cursor.close()
            conn.close()
    
    def get_variant_config(self, experiment_name, variant):
        """Get configuration for specific variant"""
        experiment = self.experiments.get(experiment_name, {})
        return experiment.get(variant, {})
    
    def track_recommendation_served(self, user_id, items, strategy, variant):
        """Track recommendations served for analytics"""
        conn = psycopg2.connect(self.db_url)
        cursor = conn.cursor()
        
        try:
            cursor.execute("""
                INSERT INTO recommendations_served 
                (user_id, recommended_items, strategy, ab_test_variant)
                VALUES (%s, %s, %s, %s)
            """, (user_id, json.dumps(items), strategy, variant))
            
            conn.commit()
        except Exception as e:
            print(f"Error tracking recommendation: {e}")
            conn.rollback()
        finally:
            cursor.close()
            conn.close()

if __name__ == "__main__":
    # Test the A/B testing system
    config = ExperimentConfig("postgresql://user:pass@localhost:5432/recommendations")
    
    # Test assignments
    test_users = ["user_123", "user_456", "user_789"]
    for user_id in test_users:
        variant = config.assign_variant(user_id, "homepage_recommendations_v1")
        print(f"User {user_id} assigned to: {variant}")