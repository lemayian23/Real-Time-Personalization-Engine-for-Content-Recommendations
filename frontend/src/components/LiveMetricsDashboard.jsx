import React, { useState, useEffect } from 'react'

const LiveMetricsDashboard = () => {
  const [metrics, setMetrics] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchMetrics()
    const interval = setInterval(fetchMetrics, 5000) // Update every 5 seconds
    return () => clearInterval(interval)
  }, [])

  const fetchMetrics = async () => {
    try {
      const response = await fetch('http://localhost:8080/engagement-metrics')
      const data = await response.json()
      setMetrics(data)
      setLoading(false)
    } catch (error) {
      console.error('Failed to fetch metrics:', error)
    }
  }

  if (loading || !metrics) {
    return <div className="loading">Loading live metrics...</div>
  }

  return (
    <div className="live-metrics-dashboard">
      <div className="metrics-header">
        <h2>📊 Live Engagement Metrics</h2>
        <div className="last-updated">
          Updated: {new Date(metrics.timestamp).toLocaleTimeString()}
        </div>
      </div>

      <div className="metrics-grid">
        <MetricCard 
          title="Total Impressions" 
          value={metrics.total_impressions?.toLocaleString()} 
          trend="up"
          icon="👁️"
        />
        <MetricCard 
          title="Total Clicks" 
          value={metrics.total_clicks?.toLocaleString()} 
          trend="up" 
          icon="👆"
        />
        <MetricCard 
          title="Current CTR" 
          value={`${(metrics.current_ctr * 100).toFixed(1)}%`} 
          trend="up"
          icon="🎯"
        />
        <MetricCard 
          title="Active Users" 
          value={metrics.active_users?.toLocaleString()} 
          trend="stable"
          icon="👥"
        />
      </div>

      <div className="performance-metrics">
        <h3>System Performance</h3>
        <div className="performance-grid">
          <div className="perf-metric">
            <span className="perf-label">P95 Latency</span>
            <span className="perf-value">{metrics.performance?.p95_latency_ms}ms</span>
          </div>
          <div className="perf-metric">
            <span className="perf-label">Error Rate</span>
            <span className="perf-value">{(metrics.performance?.error_rate * 100).toFixed(2)}%</span>
          </div>
          <div className="perf-metric">
            <span className="perf-label">Throughput</span>
            <span className="perf-value">{metrics.performance?.throughput_rps?.toLocaleString()}/s</span>
          </div>
          <div className="perf-metric">
            <span className="perf-label">Uptime</span>
            <span className="perf-value">{Math.round(metrics.uptime_minutes)}min</span>
          </div>
        </div>
      </div>

      <div className="engagement-trend">
        <h3>CTR Trend (Last 5 Minutes)</h3>
        <div className="trend-chart">
          {metrics.engagement_trend?.map((point, index) => (
            <TrendPoint 
              key={index}
              ctr={point.ctr}
              users={point.users}
              isCurrent={index === metrics.engagement_trend.length - 1}
            />
          ))}
        </div>
      </div>
    </div>
  )
}

const MetricCard = ({ title, value, trend, icon }) => {
  const getTrendColor = () => {
    switch(trend) {
      case 'up': return 'var(--forest-500)'
      case 'down': return '#ef4444'
      default: return 'var(--stone-500)'
    }
  }

  const getTrendIcon = () => {
    switch(trend) {
      case 'up': return '↗️'
      case 'down': return '↘️'
      default: return '→'
    }
  }

  return (
    <div className="metric-card">
      <div className="metric-icon">{icon}</div>
      <div className="metric-content">
        <div className="metric-value">{value}</div>
        <div className="metric-title">{title}</div>
      </div>
      <div className="metric-trend" style={{ color: getTrendColor() }}>
        {getTrendIcon()}
      </div>
    </div>
  )
}

const TrendPoint = ({ ctr, users, isCurrent }) => {
  const height = Math.max(20, ctr * 200) // Scale CTR to bar height
  
  return (
    <div className="trend-point">
      <div className="trend-bar-container">
        <div 
          className={`trend-bar ${isCurrent ? 'current' : ''}`}
          style={{ height: `${height}px` }}
        />
      </div>
      <div className="trend-label">
        <div className="ctr-value">{(ctr * 100).toFixed(0)}%</div>
        <div className="users-value">{users}u</div>
      </div>
    </div>
  )
}

export default LiveMetricsDashboard