import React, { useState, useEffect } from 'react'

const ABTestDashboard = () => {
  const [experiments, setExperiments] = useState([])
  const [selectedExperiment, setSelectedExperiment] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchExperiments()
  }, [])

  const fetchExperiments = async () => {
    try {
      const response = await fetch('http://localhost:8080/ab-tests')
      const data = await response.json()
      setExperiments(data)
    } catch (error) {
      console.error('Failed to fetch experiments:', error)
    } finally {
      setLoading(false)
    }
  }

  const fetchExperimentDetail = async (experimentId) => {
    try {
      const response = await fetch(`http://localhost:8080/ab-tests/${experimentId}`)
      const data = await response.json()
      setSelectedExperiment(data)
    } catch (error) {
      console.error('Failed to fetch experiment details:', error)
    }
  }

  if (loading) {
    return <div className="loading">Loading experiments...</div>
  }

  return (
    <div className="ab-test-dashboard">
      <div className="dashboard-header">
        <h2>A/B Test Results</h2>
        <button onClick={fetchExperiments} className="refresh-btn">
          🔄 Refresh
        </button>
      </div>

      <div className="dashboard-content">
        <div className="experiments-list">
          <h3>Active Experiments</h3>
          {experiments.map(exp => (
            <div 
              key={exp.id}
              className={`experiment-card ${selectedExperiment?.experiment?.name === exp.name ? 'selected' : ''}`}
              onClick={() => fetchExperimentDetail(exp.id)}
            >
              <div className="experiment-header">
                <h4>{exp.name}</h4>
                <span className={`status-badge ${exp.status}`}>{exp.status}</span>
              </div>
              <div className="experiment-stats">
                <div className="stat">
                  <span className="stat-value">{exp.total_impressions?.toLocaleString()}</span>
                  <span className="stat-label">Impressions</span>
                </div>
                <div className="stat">
                  <span className="stat-value">{exp.winning_variant}</span>
                  <span className="stat-label">Leading</span>
                </div>
              </div>
            </div>
          ))}
        </div>

        <div className="experiment-detail">
          {selectedExperiment ? (
            <ExperimentDetailView experiment={selectedExperiment} />
          ) : (
            <div className="empty-state">
              <h3>Select an experiment to view details</h3>
              <p>Click on an experiment card to see detailed results and statistical significance</p>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

const ExperimentDetailView = ({ experiment }) => {
  const { experiment: exp, summary } = experiment

  return (
    <div className="experiment-detail-view">
      <h3>{exp.Name}</h3>
      
      <div className="experiment-summary">
        <div className="summary-grid">
          <div className="summary-item">
            <span className="summary-label">Duration</span>
            <span className="summary-value">{summary.duration}</span>
          </div>
          <div className="summary-item">
            <span className="summary-label">Total Users</span>
            <span className="summary-value">{summary.total_users?.toLocaleString()}</span>
          </div>
          <div className="summary-item">
            <span className="summary-label">Overall CTR</span>
            <span className="summary-value">{(summary.overall_ctr * 100).toFixed(1)}%</span>
          </div>
          <div className="summary-item">
            <span className="summary-label">Detected Effect</span>
            <span className="summary-value positive">
              +{(summary.detected_effect * 100).toFixed(1)}%
            </span>
          </div>
        </div>
      </div>

      <div className="variants-comparison">
        <h4>Variant Performance</h4>
        {exp.Variants.map((variant, index) => (
          <VariantCard key={variant.Name} variant={variant} isControl={index === 0} />
        ))}
      </div>

      <div className="statistical-significance">
        <h4>Statistical Significance</h4>
        <div className="confidence-level">
          <div className="confidence-bar">
            <div 
              className="confidence-fill"
              style={{ width: `${exp.Variants[1]?.Confidence * 100}%` }}
            />
          </div>
          <span className="confidence-text">
            {exp.Variants[1]?.IsSignificant ? '✅ Statistically Significant' : '⚠️ Needs More Data'}
          </span>
        </div>
      </div>
    </div>
  )
}

const VariantCard = ({ variant, isControl }) => {
  return (
    <div className={`variant-card ${isControl ? 'control' : 'treatment'}`}>
      <div className="variant-header">
        <h5>{variant.Name}</h5>
        {variant.IsSignificant && <span className="significant-badge">★ Significant</span>}
      </div>
      
      <div className="variant-stats">
        <div className="variant-stat">
          <span className="stat-value">{variant.Impressions.toLocaleString()}</span>
          <span className="stat-label">Impressions</span>
        </div>
        <div className="variant-stat">
          <span className="stat-value">{variant.Clicks.toLocaleString()}</span>
          <span className="stat-label">Clicks</span>
        </div>
        <div className="variant-stat">
          <span className="stat-value highlight">{(variant.CTR * 100).toFixed(1)}%</span>
          <span className="stat-label">CTR</span>
        </div>
      </div>

      <div className="confidence">
        <span className="confidence-label">Confidence:</span>
        <span className="confidence-value">{(variant.Confidence * 100).toFixed(0)}%</span>
      </div>
    </div>
  )
}

export default ABTestDashboard