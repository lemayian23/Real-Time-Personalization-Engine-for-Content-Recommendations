import React from 'react'

export const MetricsDashboard = ({ metrics }) => {
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
      </div>
    </div>
  )
}