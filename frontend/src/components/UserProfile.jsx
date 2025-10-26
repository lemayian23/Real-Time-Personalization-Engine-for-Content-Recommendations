import React from 'react'

export const UserProfile = ({ userId, profile, onUserIdChange }) => {
  return (
    <div className="sidebar-panel">
      <h3>User Profile</h3>
      <div className="profile-section">
        <label className="input-label">User ID</label>
        <input
          type="text"
          value={userId}
          onChange={(e) => onUserIdChange(e.target.value)}
          className="text-input"
          placeholder="Enter user ID..."
        />
      </div>
      
      {profile && (
        <div className="profile-stats">
          <div className="stat-item">
            <span className="stat-label">Preference Strength</span>
            <span className="stat-value">Strong</span>
          </div>
          <div className="stat-item">
            <span className="stat-label">Content Categories</span>
            <span className="stat-value">4</span>
          </div>
        </div>
      )}
    </div>
  )
}