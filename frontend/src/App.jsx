import React, { useState, useEffect, useCallback } from 'react'
import { RecommendationGrid, UserProfile, MetricsDashboard, A_BTestPanel } from './components'
import ABTestDashboard from './components/ABTestDashboard'
import LiveMetricsDashboard from './components/LiveMetricsDashboard'
import UserSessionsDashboard from './components/UserSessionsDashboard'
import { recommendationService, analyticsService } from './services'
import AlgorithmVisualizationDashboard from './components/AlgorithmVisualizationDashboard'
import './App.css'

function App() {
  const [state, setState] = useState({
    userId: `user_${Date.now()}`,
    recommendations: [],
    userProfile: null,
    metrics: {},
    loading: false,
    error: null,
    currentView: 'recommendations'
  })

  const { userId, recommendations, userProfile, metrics, loading, error, currentView } = state

  const updateState = useCallback((updates) => {
    setState(prev => ({ ...prev, ...updates }))
  }, [])

  const fetchRecommendations = useCallback(async () => {
    updateState({ loading: true, error: null })
    try {
      const recs = await recommendationService.getRecommendations(userId)
      const profile = await recommendationService.getUserProfile(userId)
      const mets = await analyticsService.getMetrics()
      updateState({ 
        recommendations: recs, 
        userProfile: profile, 
        metrics: mets,
        loading: false 
      })
    } catch (err) {
      updateState({ error: err.message, loading: false })
    }
  }, [userId, updateState])

  const trackInteraction = useCallback(async (itemId, interactionType) => {
    try {
      await analyticsService.trackInteraction(userId, itemId, interactionType)
      setTimeout(fetchRecommendations, 500)
    } catch (err) {
      console.error('Tracking failed:', err)
    }
  }, [userId, fetchRecommendations])

  useEffect(() => {
    if (currentView === 'recommendations') {
      fetchRecommendations()
    }
  }, [fetchRecommendations, currentView])

  const renderCurrentView = () => {
    switch(currentView) {
      case 'recommendations':
        return (
          <div className="app-layout">
            <aside className="sidebar">
              <UserProfile 
                userId={userId}
                profile={userProfile}
                onUserIdChange={(newId) => updateState({ userId: newId })}
              />
              <MetricsDashboard metrics={metrics} recommendations={recommendations} />
              <A_BTestPanel userId={userId} />
            </aside>

            <main className="main-content">
              <div className="content-header">
                <h2>Personalized Recommendations</h2>
                <button 
                  onClick={fetchRecommendations}
                  disabled={loading}
                  className="refresh-btn"
                >
                  {loading ? '🔄 Updating...' : '🔄 Refresh'}
                </button>
              </div>

              {error && (
                <div className="error-banner">
                  <strong>Error:</strong> {error}
                </div>
              )}

              <RecommendationGrid
                recommendations={recommendations}
                onInteraction={trackInteraction}
                loading={loading}
              />
            </main>
          </div>
        )
      case 'ab-tests':
        return <ABTestDashboard />
      case 'metrics':
        return <LiveMetricsDashboard />
      case 'sessions':
        return <UserSessionsDashboard />
      default:
        return <ABTestDashboard />
    }
  }

  return (
    <div className="app">
      <header className="app-header">
        <div className="header-content">
          <h1>Recommendation Engine v2.0</h1>
          <div className="environment-tag">PRODUCTION</div>
        </div>
        
        <nav className="app-nav">
          <button 
            className={`nav-btn ${currentView === 'recommendations' ? 'active' : ''}`}
            onClick={() => updateState({ currentView: 'recommendations' })}
          >
            🎯 Recommendations
          </button>
          <button 
            className={`nav-btn ${currentView === 'ab-tests' ? 'active' : ''}`}
            onClick={() => updateState({ currentView: 'ab-tests' })}
          >
            📊 A/B Tests
          </button>
          <button 
            className={`nav-btn ${currentView === 'metrics' ? 'active' : ''}`}
            onClick={() => updateState({ currentView: 'metrics' })}
          >
            📈 Live Metrics
          </button>
          <button 
            className={`nav-btn ${currentView === 'sessions' ? 'active' : ''}`}
            onClick={() => updateState({ currentView: 'sessions' })}
          >
            👤 User Sessions
          </button>
            <button 
    className={`nav-btn ${currentView === 'recommendations' ? 'active' : ''}`}
    onClick={() => updateState({ currentView: 'recommendations' })}
  >
    🎯 Recommendations
  </button>
  <button 
    className={`nav-btn ${currentView === 'ab-tests' ? 'active' : ''}`}
    onClick={() => updateState({ currentView: 'ab-tests' })}
  >
    📊 A/B Tests
  </button>
  <button 
    className={`nav-btn ${currentView === 'metrics' ? 'active' : ''}`}
    onClick={() => updateState({ currentView: 'metrics' })}
  >
    📈 Live Metrics
  </button>
  <button 
    className={`nav-btn ${currentView === 'sessions' ? 'active' : ''}`}
    onClick={() => updateState({ currentView: 'sessions' })}
  >
    👤 User Sessions
  </button>
  <button 
    className={`nav-btn ${currentView === 'algorithm' ? 'active' : ''}`}
    onClick={() => updateState({ currentView: 'algorithm' })}
  >
    🔍 Algorithm
  </button>
        </nav>
      </header>
case 'algorithm':
  return <AlgorithmVisualizationDashboard />
      {renderCurrentView()}
    </div>
  )
}

export default App