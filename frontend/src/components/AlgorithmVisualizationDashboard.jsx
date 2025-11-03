import React, { useState, useEffect } from 'react'

const AlgorithmVisualizationDashboard = () => {
  const [users, setUsers] = useState([])
  const [selectedUser, setSelectedUser] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchUsers()
  }, [])

  const fetchUsers = async () => {
    try {
      const response = await fetch('http://localhost:8080/algorithm-visualization')
      const data = await response.json()
      setUsers(data)
    } catch (error) {
      console.error('Failed to fetch users:', error)
    } finally {
      setLoading(false)
    }
  }

  const fetchUserDetail = async (userId) => {
    try {
      const response = await fetch(`http://localhost:8080/algorithm-visualization/${userId}`)
      const data = await response.json()
      setSelectedUser(data)
    } catch (error) {
      console.error('Failed to fetch user details:', error)
    }
  }

  if (loading) {
    return <div className="loading">Loading algorithm visualization...</div>
  }

  return (
    <div className="algorithm-visualization-dashboard">
      <div className="dashboard-header">
        <h2>🔍 Algorithm Visualization</h2>
        <button onClick={fetchUsers} className="refresh-btn">
          🔄 Refresh
        </button>
      </div>

      <div className="dashboard-content">
        <div className="users-list">
          <h3>User Analysis</h3>
          {users.map(user => (
            <div 
              key={user.user_id}
              className={`user-card ${selectedUser?.algorithm_state?.current_user === user.user_id ? 'selected' : ''}`}
              onClick={() => fetchUserDetail(user.user_id)}
            >
              <div className="user-header">
                <h4>User: {user.user_id}</h4>
                <span className="status-badge">Analyzed</span>
              </div>
              
              <div className="user-stats">
                <div className="user-stat">
                  <span className="stat-value">{user.interaction_count}</span>
                  <span className="stat-label">Interactions</span>
                </div>
                <div className="user-stat">
                  <span className="stat-value">{user.similar_users_count}</span>
                  <span className="stat-label">Similar Users</span>
                </div>
                <div className="user-stat">
                  <span className="stat-value">{user.recommendation_count}</span>
                  <span className="stat-label">Recommendations</span>
                </div>
              </div>

              <div className="last-updated">
                Updated: {new Date(user.last_updated).toLocaleTimeString()}
              </div>
            </div>
          ))}
        </div>

        <div className="visualization-detail">
          {selectedUser ? (
            <AlgorithmDetailView data={selectedUser} />
          ) : (
            <div className="empty-state">
              <h3>Select a user to view algorithm visualization</h3>
              <p>Click on a user card to see how collaborative filtering works for that user</p>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

const AlgorithmDetailView = ({ data }) => {
  const { algorithm_state, user_interactions, similarity_matrix, performance } = data

  return (
    <div className="algorithm-detail-view">
      <h3>Algorithm Analysis: {algorithm_state.current_user}</h3>
      
      <div className="explanation-box">
        <h4>🤔 How It Works</h4>
        <p className="explanation-text">{algorithm_state.explanation}</p>
      </div>

      <div className="algorithm-sections">
        <div className="section">
          <h4>👥 Similar Users</h4>
          <div className="similar-users-grid">
            {algorithm_state.similar_users.map((similar, index) => (
              <SimilarUserCard key={similar.user2} similar={similar} rank={index + 1} />
            ))}
          </div>
        </div>

        <div className="section">
          <h4>📊 User-Item Interactions</h4>
          <InteractionMatrix interactions={user_interactions} currentUser={algorithm_state.current_user} />
        </div>

        <div className="section">
          <h4>🎯 Generated Recommendations</h4>
          <RecommendationsList recommendations={algorithm_state.recommendations} />
        </div>

        <div className="section">
          <h4>⚡ Performance</h4>
          <PerformanceMetrics performance={performance} />
        </div>
      </div>
    </div>
  )
}

const SimilarUserCard = ({ similar, rank }) => {
  return (
    <div className="similar-user-card">
      <div className="similar-user-header">
        <span className="rank">#{rank}</span>
        <span className="username">{similar.user2}</span>
        <span className="similarity-score">{Math.round(similar.score * 100)}%</span>
      </div>
      
      <div className="shared-items">
        <span className="shared-label">Shared interests:</span>
        <div className="items-list">
          {similar.shared_items.map(item => (
            <span key={item} className="shared-item">{item}</span>
          ))}
        </div>
      </div>
    </div>
  )
}

const InteractionMatrix = ({ interactions, currentUser }) => {
  const allItems = [...new Set(Object.values(interactions).flat())]
  
  return (
    <div className="interaction-matrix">
      <div className="matrix-header">
        <div className="user-column">User</div>
        {allItems.map(item => (
          <div key={item} className="item-header">{item.split('_')[0]}</div>
        ))}
      </div>
      
      {Object.entries(interactions).map(([user, items]) => (
        <div key={user} className={`matrix-row ${user === currentUser ? 'current-user' : ''}`}>
          <div className="user-cell">{user}</div>
          {allItems.map(item => (
            <div key={item} className="interaction-cell">
              {items.includes(item) ? '✅' : '❌'}
            </div>
          ))}
        </div>
      ))}
    </div>
  )
}

const RecommendationsList = ({ recommendations }) => {
  return (
    <div className="recommendations-list">
      {recommendations.length === 0 ? (
        <div className="no-recommendations">No recommendations generated yet</div>
      ) : (
        recommendations.map((rec, index) => (
          <div key={rec} className="recommendation-item">
            <span className="rec-rank">#{index + 1}</span>
            <span className="rec-item">{rec}</span>
            <span className="rec-source">From similar users</span>
          </div>
        ))
      )}
    </div>
  )
}

const PerformanceMetrics = ({ performance }) => {
  return (
    <div className="performance-metrics">
      <div className="perf-metric">
        <span className="perf-label">Calculation Time</span>
        <span className="perf-value">{performance.calculation_time_ms}ms</span>
      </div>
      <div className="perf-metric">
        <span className="perf-label">Similarity Threshold</span>
        <span className="perf-value">{performance.similarity_threshold}</span>
      </div>
      <div className="perf-metric">
        <span className="perf-label">Min Shared Items</span>
        <span className="perf-value">{performance.min_shared_items}</span>
      </div>
    </div>
  )
}

export default AlgorithmVisualizationDashboard