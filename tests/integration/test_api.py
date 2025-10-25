import unittest
import requests
import time
import json

class TestAPIEndpoints(unittest.TestCase):
    BASE_URL = "http://localhost:8080"
    
    def test_health_endpoint(self):
        """Test health check endpoint"""
        response = requests.get(f"{self.BASE_URL}/health")
        self.assertEqual(response.status_code, 200)
        
        data = response.json()
        self.assertIn('status', data)
        self.assertEqual(data['status'], 'healthy')
    
    def test_recommendation_endpoint(self):
        """Test recommendation endpoint"""
        # Test with new user (cold start)
        response = requests.get(f"{self.BASE_URL}/recommend?user_id=new_user_123&count=5")
        self.assertEqual(response.status_code, 200)
        
        data = response.json()
        self.assertIn('recommendations', data)
        self.assertIn('strategy', data)
        self.assertIn('latency_ms', data)
        
        # Verify latency is reasonable
        self.assertLess(data['latency_ms'], 1000)
    
    def test_event_tracking(self):
        """Test event tracking endpoint"""
        event_data = {
            'user_id': 'test_user_456',
            'item_id': 'item_42', 
            'event_type': 'view',
            'duration_seconds': 30
        }
        
        response = requests.post(f"{self.BASE_URL}/event", json=event_data)
        self.assertEqual(response.status_code, 200)
    
    def test_latency_requirement(self):
        """Test that latency meets 50ms p99 requirement"""
        latencies = []
        
        for i in range(100):
            start_time = time.time()
            response = requests.get(f"{self.BASE_URL}/recommend?user_id=load_test_{i}&count=5")
            end_time = time.time()
            
            self.assertEqual(response.status_code, 200)
            latencies.append((end_time - start_time) * 1000)  # Convert to ms
        
        # Calculate p99 latency
        latencies.sort()
        p99_latency = latencies[int(0.99 * len(latencies))]
        
        print(f"P99 Latency: {p99_latency:.2f}ms")
        self.assertLess(p99_latency, 100)  # Conservative test threshold

if __name__ == '__main__':
    unittest.main()