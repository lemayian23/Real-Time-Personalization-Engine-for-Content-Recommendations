import React from 'react'

export const RecommendationGrid = ({ recommendations, onInteraction, loading }) => {
  if (loading) {
    return <div className="loading-grid">Analyzing user preferences...</div>
  }

  return (
    <div className="recommendation-grid">
      {recommendations.map((rec, index) => (
        <RecommendationCard
          key={`${rec.item_id}-${index}`}
          recommendation={rec}
          rank={index + 1}
          onInteraction={onInteraction}
        />
      ))}
    </div>
  )
}

const RecommendationCard = ({ recommendation, rank, onInteraction }) => {
  const { item_id, score, explanation, strategy } = recommendation

  return (
    <div className="recommendation-card" data-strategy={strategy}>
      <div className="card-header">
        <span className="rank-badge">#{rank}</span>
        <span className="strategy-tag">{strategy}</span>
      </div>
      
      <div className="card-content">
        <h3 className="item-title">{item_id.replace(/_/g, ' ').toUpperCase()}</h3>
        <div className="confidence-score">
          <div 
            className="score-bar" 
            style={{ width: `${score * 100}%` }}
          />
          <span className="score-text">{Math.round(score * 100)}% match</span>
        </div>
        <p className="explanation">{explanation}</p>
      </div>

      <div className="card-actions">
        <button 
          onClick={() => onInteraction(item_id, 'view')}
          className="action-btn primary"
        >
          View Content
        </button>
        <button 
          onClick={() => onInteraction(item_id, 'click')}
          className="action-btn secondary"
        >
          Engage
        </button>
      </div>
    </div>
  )
}