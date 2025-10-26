import React, { useState, useEffect, useCallback } from 'react'
import { RecommendationGrid, UserProfile, MetricsDashboard, A_BTestPanel } from './components'
import { recommendationService, analyticsService } from './services'
import './App.css'

function App() {
  const [state, setState] = useState({
    userId: `user_${Date.now()}`,
    recommendations: [],
    userProfile: null,
    metrics: {},
    loading: false,
    error: null
  })

  const { userId, recommendations, userProfile, metrics, loading, error } = state

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
      // Real-time update - refetch to show system learning
      setTimeout(fetchRecommendations, 500)
    } catch (err) {
      console.error('Tracking failed:', err)
    }
  }, [userId, fetchRecommendations])

  useEffect(() => {
    fetchRecommendations()
  }, [fetchRecommendations])

  return (
    <div className="app">
      <header className="app-header">
        <div className="header-content">
          <h1>Recommendation Engine v2.0</h1>
          <div className="environment-tag">PRODUCTION</div>
        </div>
      </header>

      <div className="app-layout">
        <aside className="sidebar">
          <UserProfile 
            userId={userId}
            profile={userProfile}
            onUserIdChange={(newId) => updateState({ userId: newId })}
          />
          <MetricsDashboard metrics={metrics} />
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
    </div>
  )
}

export default App