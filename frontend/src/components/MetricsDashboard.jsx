import React from 'react'

export const MetricsDashboard = ({ metrics, recommendations = [] }) => {
  const calculateDiversity = (recs) => {
    if (!recs.length) return 0
    const categories = new Set(recs.map(r => {
      const parts = r.item_id.split('_')
      return parts.length > 1 ? parts[0] : 'general'
    }))
    return (categories.size / recs.length).toFixed(2)
  }

  const diversityScore = calculateDiversity(recommendations)
  const getDiversityColor = (score) => {
    if (score >= 0.7) return 'var(--forest-500)'
    if (score >= 0.4) return 'var(--amber-500)'
    return '#ef4444'
  }

  return (
    <div className="sidebar-panel">
      <h3>System Metrics</h3>
      <div className="metrics-grid">
        <div className="metric-card">
          <div className="metric-value">24ms</div>
          <div className="metric-label">Avg Latency</div>
        </div>
        <div className="metric-card">
          <div className="metric-value">98.7%</div>
          <div className="metric-label">Uptime</div>
        </div>
        <div className="metric-card">
          <div className="metric-value">1.2k</div>
          <div className="metric-label">Req/Min</div>
        </div>
        {/* NEW: Diversity Score Card */}
        <div className="metric-card">
          <div 
            className="metric-value" 
            style={{ color: getDiversityColor(diversityScore) }}
          >
            {diversityScore}
          </div>
          <div className="metric-label">Diversity Score</div>
        </div>
      </div>
      
      {/* NEW: Diversity Explanation */}
      {recommendations.length > 0 && (
        <div className="diversity-info">
          <div className="diversity-bar">
            <div 
              className="diversity-fill"
              style={{ 
                width: `${diversityScore * 100}%`,
                backgroundColor: getDiversityColor(diversityScore)
              }}
            />
          </div>
          <div className="diversity-text">
            {diversityScore >= 0.7 ? '🎯 Excellent diversity' :
             diversityScore >= 0.4 ? '⚠️ Moderate diversity' :
             '🔴 Low diversity - risk of filter bubbles'}
          </div>
        </div>
      )}
    </div>
  )
}