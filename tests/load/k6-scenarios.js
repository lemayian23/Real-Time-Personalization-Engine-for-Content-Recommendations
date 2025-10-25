import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    stages: [
        { duration: '2m', target: 1000 },  // Ramp up to 1000 users
        { duration: '5m', target: 1000 },  // Stay at 1000 users
        { duration: '2m', target: 0 },     // Ramp down to 0 users
    ],
    thresholds: {
        http_req_duration: ['p(99)<50'],  // 99% of requests must complete under 50ms
    },
};

export default function () {
    const userId = `user_${Math.floor(Math.random() * 10000)}`;
    
    // Test recommendation endpoint
    const recResponse = http.get(`http://localhost:8080/recommend?user_id=${userId}&count=10`);
    
    check(recResponse, {
        'recommendation status is 200': (r) => r.status === 200,
        'recommendation latency < 50ms': (r) => r.timings.duration < 50,
        'response has recommendations': (r) => r.json('recommendations').length > 0,
    });
    
    // Test event tracking
    const eventPayload = JSON.stringify({
        user_id: userId,
        item_id: `item_${Math.floor(Math.random() * 100)}`,
        event_type: 'view',
        duration_seconds: 30
    });
    
    const eventResponse = http.post('http://localhost:8080/event', eventPayload, {
        headers: { 'Content-Type': 'application/json' },
    });
    
    check(eventResponse, {
        'event tracking status is 200': (r) => r.status === 200,
    });
    
    sleep(0.1); // 100ms between requests
}