import React from 'react'

export const A_BTestPanel = ({ userId }) => {
  return (
    <div className="sidebar-panel">
      <h3>A/B Testing</h3>
      <div className="ab-test-info">
        <div className="test-item">
          <span className="test-name">Algorithm Variant</span>
          <span className="test-value">Personalized v2</span>
        </div>
        <div className="test-item">
          <span className="test-name">User Bucket</span>
          <span className="test-value">B-{userId.slice(-3)}</span>
        </div>
      </div>
    </div>
  )
}