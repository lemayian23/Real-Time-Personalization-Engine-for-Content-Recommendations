import React, { useState, useEffect } from 'react'

const ContentAnalyticsDashboard = () => {
  const [analytics, setAnalytics] = useState(null)
  const [selectedItem, setSelectedItem] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchAnalytics()
    const interval = setInterval(fetchAnalytics, 10000) // Update every 10 seconds
    return () => clearInterval(interval)
  }, [])

  const fetchAnalytics = async () => {
    try {
      const response = await fetch('http://localhost:8080/content-analytics')
      const data = await response.json()
      setAnalytics(data)
      setLoading(false)
    } catch (error) {
      console.error('Failed to fetch content analytics:', error)
    }
  }

  const fetchItemDetail = async (itemId) => {
    try {
      const response = await fetch(`http://localhost:8080/content-analytics/${itemId}`)
      const data = await response.json()
      setSelectedItem(data)
    } catch (error) {
      console.error('Failed to fetch item details:', error)
    }
  }

  if (loading || !analytics) {
    return <div className="loading">Loading content analytics...</div>
  }

  return (
    <div className="content-analytics-dashboard">
      <div className="dashboard-header">
        <h2>📊 Content Performance Analytics</h2>
        <div className="last-updated">
          Updated: {new Date(analytics.summary.last_updated).toLocaleTimeString()}
        </div>
      </div>

      <div className="summary-cards">
        <SummaryCard 
          title="Total Items" 
          value={analytics.summary.total_items} 
          icon="📄"
        />
        <SummaryCard 
          title="Total Impressions" 
          value={analytics.summary.total_impressions?.toLocaleString()} 
          icon="👁️"
        />
        <SummaryCard 
          title="Total Clicks" 
          value={analytics.summary.total_clicks?.toLocaleString()} 
          icon="👆"
        />
        <SummaryCard 
          title="Overall CTR" 
          value={`${(analytics.summary.overall_ctr * 100).toFixed(1)}%`} 
          icon="🎯"
        />
        <SummaryCard 
          title="Avg Duration" 
          value={`${Math.round(analytics.summary.avg_duration)}s`} 
          icon="⏱️"
        />
      </div>

      <div className="dashboard-content">
        <div className="content-lists">
          <div className="list-section">
            <h3>🏆 Top Performing Content</h3>
            <div className="content-list">
              {analytics.top_performing.map((item, index) => (
                <ContentItemCard 
                  key={item.item_id}
                  item={item}
                  rank={index + 1}
                  onSelect={() => fetchItemDetail(item.item_id)}
                  isSelected={selectedItem?.content_performance?.item_id === item.item_id}
                />
              ))}
            </div>
          </div>

          <div className="list-section">
            <h3>🚀 Trending Now</h3>
            <div className="content-list">
              {analytics.trending_items.map((item, index) => (
                <TrendingItemCard 
                  key={item.item_id}
                  item={item}
                  rank={index + 1}
                  onSelect={() => fetchItemDetail(item.item_id)}
                />
              ))}
            </div>
          </div>
        </div>

        <div className="detail-panel">
          {selectedItem ? (
            <ContentItemDetail itemData={selectedItem} />
          ) : (
            <div className="empty-state">
              <h3>Select a content item to view details</h3>
              <p>Click on any content card to see detailed performance analytics</p>
            </div>
          )}
        </div>
      </div>

      <div className="categories-section">
        <h3>📈 Category Performance</h3>
        <div className="categories-grid">
          {analytics.categories.map(category => (
            <CategoryCard key={category.category} category={category} />
          ))}
        </div>
      </div>
    </div>
  )
}

const SummaryCard = ({ title, value, icon }) => {
  return (
    <div className="summary-card">
      <div className="summary-icon">{icon}</div>
      <div className="summary-content">
        <div className="summary-value">{value}</div>
        <div className="summary-title">{title}</div>
      </div>
    </div>
  )
}

const ContentItemCard = ({ item, rank, onSelect, isSelected }) => {
  return (
    <div 
      className={`content-item-card ${isSelected ? 'selected' : ''}`}
      onClick={onSelect}
    >
      <div className="item-header">
        <span className="item-rank">#{rank}</span>
        <span className="item-category">{item.category}</span>
      </div>
      <h4 className="item-title">{item.title}</h4>
      <div className="item-stats">
        <div className="item-stat">
          <span className="stat-value">{item.impressions?.toLocaleString()}</span>
          <span className="stat-label">Impressions</span>
        </div>
        <div className="item-stat">
          <span className="stat-value">{item.clicks?.toLocaleString()}</span>
          <span className="stat-label">Clicks</span>
        </div>
        <div className="item-stat">
          <span className="stat-value highlight">{(item.ctr * 100).toFixed(1)}%</span>
          <span className="stat-label">CTR</span>
        </div>
      </div>
    </div>
  )
}

const TrendingItemCard = ({ item, rank, onSelect }) => {
  return (
    <div className="trending-item-card" onClick={onSelect}>
      <div className="trending-header">
        <span className="trending-rank">#{rank}</span>
        <span className="trending-badge">🔥 Trending</span>
      </div>
      <h4 className="item-title">{item.title}</h4>
      <div className="trending-stats">
        <span className="ctr">{(item.ctr * 100).toFixed(1)}% CTR</span>
        <span className="duration">{Math.round(item.avg_duration)}s avg</span>
      </div>
    </div>
  )
}

const CategoryCard = ({ category }) => {
  const getTrendIcon = () => {
    switch(category.trend) {
      case 'up': return '↗️'
      case 'down': return '↘️'
      default: return '→'
    }
  }

  const getTrendColor = () => {
    switch(category.trend) {
      case 'up': return 'var(--forest-500)'
      case 'down': return '#ef4444'
      default: return 'var(--stone-500)'
    }
  }

  return (
    <div className="category-card">
      <div className="category-header">
        <h4 className="category-name">{category.category}</h4>
        <span className="category-trend" style={{ color: getTrendColor() }}>
          {getTrendIcon()}
        </span>
      </div>
      <div className="category-stats">
        <div className="category-stat">
          <span className="stat-value">{category.total_items}</span>
          <span className="stat-label">Items</span>
        </div>
        <div className="category-stat">
          <span className="stat-value">{category.total_clicks?.toLocaleString()}</span>
          <span className="stat-label">Clicks</span>
        </div>
        <div className="category-stat">
          <span className="stat-value">{(category.avg_ctr * 100).toFixed(1)}%</span>
          <span className="stat-label">Avg CTR</span>
        </div>
      </div>
    </div>
  )
}

const ContentItemDetail = ({ itemData }) => {
  const { content_performance, performance_history, similar_items, recommendation_impact } = itemData

  return (
    <div className="content-item-detail">
      <h3>{content_performance.title}</h3>
      <div className="item-meta">
        <span className="item-id">ID: {content_performance.item_id}</span>
        <span className="item-category">{content_performance.category}</span>
      </div>

      <div className="performance-overview">
        <h4>Performance Overview</h4>
        <div className="overview-grid">
          <div className="overview-item">
            <span className="overview-label">Impressions</span>
            <span className="overview-value">{content_performance.impressions?.toLocaleString()}</span>
          </div>
          <div className="overview-item">
            <span className="overview-label">Clicks</span>
            <span className="overview-value">{content_performance.clicks?.toLocaleString()}</span>
          </div>
          <div className="overview-item">
            <span className="overview-label">CTR</span>
            <span className="overview-value highlight">{(content_performance.ctr * 100).toFixed(1)}%</span>
          </div>
          <div className="overview-item">
            <span className="overview-label">Avg Duration</span>
            <span className="overview-value">{Math.round(content_performance.avg_duration)}s</span>
          </div>
        </div>
      </div>

      <div className="performance-history">
        <h4>7-Day Performance History</h4>
        <div className="history-chart">
          {performance_history.map((day, index) => (
            <div key={index} className="history-day">
              <div className="day-label">{day.date.split('-')[2]}</div>
              <div className="day-bar">
                <div 
                  className="ctr-bar"
                  style={{ height: `${day.ctr * 500}px` }}
                />
              </div>
              <div className="day-value">{(day.ctr * 100).toFixed(1)}%</div>
            </div>
          ))}
        </div>
      </div>

      <div className="recommendation-impact">
        <h4>Recommendation Impact</h4>
        <div className="impact-stats">
          <div className="impact-stat">
            <span className="impact-value">{recommendation_impact.times_recommended}</span>
            <span className="impact-label">Times Recommended</span>
          </div>
          <div className="impact-stat">
            <span className="impact-value">{(recommendation_impact.conversion_rate * 100).toFixed(1)}%</span>
            <span className="impact-label">Conversion Rate</span>
          </div>
          <div className="impact-stat">
            <span className="impact-value">{Math.round(recommendation_impact.engagement_score)}</span>
            <span className="impact-label">Engagement Score</span>
          </div>
        </div>
      </div>

      <div className="similar-items">
        <h4>Similar Content</h4>
        <div className="similar-list">
          {similar_items.map(item => (
            <div key={item.item_id} className="similar-item">
              <span className="similar-title">{item.title}</span>
              <span className="similar-ctr">{(item.ctr * 100).toFixed(1)}% CTR</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

export default ContentAnalyticsDashboard