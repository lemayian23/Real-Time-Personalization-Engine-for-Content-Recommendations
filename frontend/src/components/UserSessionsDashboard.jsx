import React, { useState, useEffect } from 'react'

const UserSessionsDashboard = () => {
  const [sessions, setSessions] = useState([])
  const [selectedSession, setSelectedSession] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchSessions()
  }, [])

  const fetchSessions = async () => {
    try {
      const response = await fetch('http://localhost:8080/user-sessions')
      const data = await response.json()
      setSessions(data)
    } catch (error) {
      console.error('Failed to fetch sessions:', error)
    } finally {
      setLoading(false)
    }
  }

  const fetchSessionDetail = async (userId) => {
    try {
      const response = await fetch(`http://localhost:8080/user-sessions/${userId}`)
      const data = await response.json()
      setSelectedSession(data)
    } catch (error) {
      console.error('Failed to fetch session details:', error)
    }
  }

  if (loading) {
    return <div className="loading">Loading user sessions...</div>
  }

  return (
    <div className="user-sessions-dashboard">
      <div className="dashboard-header">
        <h2>👤 User Session Analytics</h2>
        <button onClick={fetchSessions} className="refresh-btn">
          🔄 Refresh
        </button>
      </div>

      <div className="dashboard-content">
        <div className="sessions-list">
          <h3>Active User Sessions</h3>
          {sessions.map(session => (
            <div 
              key={session.session_id}
              className={`session-card ${selectedSession?.session?.user_id === session.user_id ? 'selected' : ''}`}
              onClick={() => fetchSessionDetail(session.user_id)}
            >
              <div className="session-header">
                <h4>User: {session.user_id}</h4>
                <span className={`status-badge ${session.status}`}>{session.status}</span>
              </div>
              
              <div className="session-stats">
                <div className="session-stat">
                  <span className="stat-value">{session.page_views}</span>
                  <span className="stat-label">Views</span>
                </div>
                <div className="session-stat">
                  <span className="stat-value">{session.clicks}</span>
                  <span className="stat-label">Clicks</span>
                </div>
                <div className="session-stat">
                  <span className="stat-value">{Math.round(session.session_time / 60)}m</span>
                  <span className="stat-label">Duration</span>
                </div>
                <div className="session-stat">
                  <span className="stat-value">{(session.engagement * 100).toFixed(0)}%</span>
                  <span className="stat-label">Engagement</span>
                </div>
              </div>

              <div className="session-categories">
                {session.categories?.map(category => (
                  <span key={category} className="category-tag">{category}</span>
                ))}
              </div>
            </div>
          ))}
        </div>

        <div className="session-detail">
          {selectedSession ? (
            <SessionDetailView sessionData={selectedSession} />
          ) : (
            <div className="empty-state">
              <h3>Select a user session to view details</h3>
              <p>Click on a user session card to see detailed analytics and click stream</p>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

const SessionDetailView = ({ sessionData }) => {
  const { session, analytics } = sessionData

  return (
    <div className="session-detail-view">
      <h3>User: {session.user_id}</h3>
      
      <div className="session-overview">
        <div className="overview-grid">
          <div className="overview-item">
            <span className="overview-label">Session Duration</span>
            <span className="overview-value">{Math.round(session.session_time / 60)} minutes</span>
          </div>
          <div className="overview-item">
            <span className="overview-label">Page Views</span>
            <span className="overview-value">{session.page_views}</span>
          </div>
          <div className="overview-item">
            <span className="overview-label">Clicks</span>
            <span className="overview-value">{session.clicks}</span>
          </div>
          <div className="overview-item">
            <span className="overview-label">Engagement Score</span>
            <span className="overview-value highlight">
              {Math.round(analytics.engagement_score)}/100
            </span>
          </div>
        </div>
      </div>

      <div className="category-distribution">
        <h4>Interest Categories</h4>
        <div className="categories-grid">
          {Object.entries(analytics.category_distribution || {}).map(([category, count]) => (
            <div key={category} className="category-item">
              <span className="category-name">{category}</span>
              <div className="category-bar">
                <div 
                  className="category-fill"
                  style={{ width: `${(count / 15) * 100}%` }}
                />
              </div>
              <span className="category-count">{count} views</span>
            </div>
          ))}
        </div>
      </div>

      <div className="click-stream">
        <h4>Click Stream Timeline</h4>
        <div className="timeline">
          {analytics.click_stream?.map((event, index) => (
            <div key={index} className="timeline-event">
              <div className="event-time">
                {new Date(event.timestamp).toLocaleTimeString()}
              </div>
              <div className="event-content">
                <div className="event-action">{event.action}</div>
                <div className="event-item">{event.item_id}</div>
                {event.duration > 0 && (
                  <div className="event-duration">{event.duration}s</div>
                )}
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

export default UserSessionsDashboard