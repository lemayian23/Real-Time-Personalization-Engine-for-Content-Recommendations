const API_BASE = 'http://localhost:8080'

class RecommendationService {
  async getRecommendations(userId, count = 10) {
    const response = await fetch(`${API_BASE}/recommend?user_id=${userId}&count=${count}`)
    const data = await response.json()
    return data.recommendations || []
  }

  async getUserProfile(userId) {
    // Mock user profile for now
    return {
      id: userId,
      preferenceStrength: 'strong',
      categories: ['tech', 'science', 'business']
    }
  }
}

class AnalyticsService {
  async trackInteraction(userId, itemId, interactionType) {
    const response = await fetch(`${API_BASE}/event`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        user_id: userId,
        item_id: itemId,
        event_type: interactionType
      })
    })
    return response.json()
  }

  async getMetrics() {
    const response = await fetch(`${API_BASE}/metrics`)
    return response.json()
  }
}

export const recommendationService = new RecommendationService()
export const analyticsService = new AnalyticsService()