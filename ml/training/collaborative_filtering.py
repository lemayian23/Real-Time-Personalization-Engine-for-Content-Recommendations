import numpy as np
from scipy.sparse import csr_matrix
import pickle
import os

class CollaborativeFiltering:
    def __init__(self, n_factors=100):
        self.n_factors = n_factors
        self.user_factors = None
        self.item_factors = None
        
    def fit(self, user_item_matrix):
        """Train matrix factorization using ALS"""
        n_users, n_items = user_item_matrix.shape
        
        # Initialize factors
        self.user_factors = np.random.normal(0, 0.1, (n_users, self.n_factors))
        self.item_factors = np.random.normal(0, 0.1, (n_items, self.n_factors))
        
        # Simple ALS implementation
        for iteration in range(10):
            # Fix items, solve for users
            for u in range(n_users):
                self.user_factors[u] = self._solve_least_squares(
                    self.item_factors, user_item_matrix[u]
                )
            
            # Fix users, solve for items
            for i in range(n_items):
                self.item_factors[i] = self._solve_least_squares(
                    self.user_factors, user_item_matrix[:, i]
                )
            
            if iteration % 2 == 0:
                loss = self._calculate_loss(user_item_matrix)
                print(f"Iteration {iteration}, Loss: {loss:.4f}")
    
    def _solve_least_squares(self, factors, ratings):
        """Solve least squares problem for ALS"""
        # Simplified implementation
        return np.linalg.lstsq(factors.T @ factors, factors.T @ ratings, rcond=None)[0]
    
    def _calculate_loss(self, user_item_matrix):
        """Calculate reconstruction error"""
        predictions = self.user_factors @ self.item_factors.T
        return np.mean((user_item_matrix - predictions) ** 2)
    
    def save_model(self, path):
        """Save trained model"""
        with open(path, 'wb') as f:
            pickle.dump({
                'user_factors': self.user_factors,
                'item_factors': self.item_factors
            }, f)
    
    def load_model(self, path):
        """Load trained model"""
        with open(path, 'rb') as f:
            data = pickle.load(f)
            self.user_factors = data['user_factors']
            self.item_factors = data['item_factors']

if __name__ == "__main__":
    # Example usage
    cf = CollaborativeFiltering(n_factors=50)
    # Mock data: 1000 users, 500 items
    mock_matrix = np.random.randint(0, 2, (1000, 500))
    cf.fit(mock_matrix)
    cf.save_model('ml/models/collaborative_filtering.pkl')