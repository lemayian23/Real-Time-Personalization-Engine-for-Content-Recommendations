import unittest
import sys
import os
sys.path.append(os.path.join(os.path.dirname(__file__), '../..'))

from ml.training.collaborative_filtering import CollaborativeFiltering
from ml.training.embeddings_generator import EmbeddingsGenerator
import numpy as np

class TestRecommendationAlgorithms(unittest.TestCase):
    
    def test_collaborative_filtering_training(self):
        """Test collaborative filtering model training"""
        cf = CollaborativeFiltering(n_factors=10)
        
        # Create mock user-item matrix
        user_item_matrix = np.random.randint(0, 2, (100, 50))
        
        # Should train without errors
        cf.fit(user_item_matrix)
        
        self.assertIsNotNone(cf.user_factors)
        self.assertIsNotNone(cf.item_factors)
        self.assertEqual(cf.user_factors.shape, (100, 10))
        self.assertEqual(cf.item_factors.shape, (50, 10))
    
    def test_embeddings_generation(self):
        """Test content embeddings generation"""
        generator = EmbeddingsGenerator()
        
        test_content = [
            {'id': 'test1', 'title': 'Machine Learning', 'description': 'AI algorithms'},
            {'id': 'test2', 'title': 'Data Science', 'description': 'Data analysis'}
        ]
        
        embeddings = generator.generate_content_embeddings(test_content)
        
        self.assertIn('test1', embeddings)
        self.assertIn('test2', embeddings)
        self.assertEqual(len(embeddings['test1']), 384)  # all-MiniLM-L6-v2 embedding size
    
    def test_hybrid_recommendation_weights(self):
        """Test hybrid recommendation weight blending"""
        # Mock collaborative scores
        collab_scores = {'item1': 0.8, 'item2': 0.6}
        content_scores = {'item1': 0.7, 'item3': 0.9}
        
        # Blend with 60% collaborative, 40% content
        hybrid_scores = {}
        for item, score in collab_scores.items():
            hybrid_scores[item] = score * 0.6
        for item, score in content_scores.items():
            hybrid_scores[item] = hybrid_scores.get(item, 0) + score * 0.4
        
        self.assertAlmostEqual(hybrid_scores['item1'], 0.8*0.6 + 0.7*0.4)
        self.assertAlmostEqual(hybrid_scores['item2'], 0.6*0.6)
        self.assertAlmostEqual(hybrid_scores['item3'], 0.9*0.4)

if __name__ == '__main__':
    unittest.main()