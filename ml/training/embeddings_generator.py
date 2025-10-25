from sentence_transformers import SentenceTransformer
import pandas as pd
import numpy as np
import json

class EmbeddingsGenerator:
    def __init__(self, model_name='all-MiniLM-L6-v2'):
        self.model = SentenceTransformer(model_name)
    
    def generate_content_embeddings(self, content_data):
        """Generate embeddings for content items"""
        texts = []
        item_ids = []
        
        for item in content_data:
            # Combine title and description for embedding
            text = f"{item['title']} {item.get('description', '')}"
            texts.append(text)
            item_ids.append(item['id'])
        
        embeddings = self.model.encode(texts, show_progress_bar=True)
        
        # Return as dictionary for easy lookup
        return {
            item_id: embedding 
            for item_id, embedding in zip(item_ids, embeddings)
        }
    
    def save_embeddings(self, embeddings, path):
        """Save embeddings to file"""
        # Convert numpy arrays to lists for JSON serialization
        serializable = {
            item_id: embedding.tolist() 
            for item_id, embedding in embeddings.items()
        }
        
        with open(path, 'w') as f:
            json.dump(serializable, f)
    
    def load_embeddings(self, path):
        """Load embeddings from file"""
        with open(path, 'r') as f:
            serializable = json.load(f)
        
        return {
            item_id: np.array(embedding) 
            for item_id, embedding in serializable.items()
        }

if __name__ == "__main__":
    # Example usage
    generator = EmbeddingsGenerator()
    
    # Mock content data
    mock_content = [
        {'id': '1', 'title': 'Machine Learning Basics', 'description': 'Introduction to ML algorithms'},
        {'id': '2', 'title': 'Deep Learning Advanced', 'description': 'Neural networks and deep learning'},
        {'id': '3', 'title': 'Data Science Fundamentals', 'description': 'Data analysis and visualization'},
    ]
    
    embeddings = generator.generate_content_embeddings(mock_content)
    generator.save_embeddings(embeddings, 'ml/models/content_embeddings.json')