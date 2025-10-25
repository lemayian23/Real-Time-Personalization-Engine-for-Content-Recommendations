import pandas as pd
import numpy as np
from scipy import stats
import psycopg2

class ABTestAnalyzer:
    def __init__(self, db_url):
        self.db_url = db_url
    
    def calculate_ctr(self, experiment_name, days_back=7):
        """Calculate CTR for control vs treatment groups"""
        conn = psycopg2.connect(self.db_url)
        
        query = """
        SELECT 
            aba.variant,
            COUNT(rs.id) as impressions,
            COUNT(ue.id) as clicks,
            COUNT(ue.id) * 1.0 / COUNT(rs.id) as ctr
        FROM ab_test_assignments aba
        LEFT JOIN recommendations_served rs ON aba.user_id = rs.user_id 
            AND rs.created_at >= NOW() - INTERVAL '%s days'
        LEFT JOIN user_events ue ON aba.user_id = ue.user_id 
            AND ue.event_type = 'click'
            AND ue.created_at >= NOW() - INTERVAL '%s days'
        WHERE aba.experiment_name = %s
        GROUP BY aba.variant
        """
        
        df = pd.read_sql_query(query, conn, params=[days_back, days_back, experiment_name])
        conn.close()
        
        return df
    
    def run_significance_test(self, experiment_name):
        """Run statistical significance test on CTR"""
        data = self.calculate_ctr(experiment_name)
        
        if len(data) < 2:
            return {"error": "Not enough data"}
        
        # Extract CTR values for t-test
        control_data = data[data['variant'] == 'control']
        treatment_data = data[data['variant'] == 'treatment']
        
        if len(control_data) == 0 or len(treatment_data) == 0:
            return {"error": "Missing variant data"}
        
        control_ctr = control_data['ctr'].iloc[0]
        treatment_ctr = treatment_data['ctr'].iloc[0]
        
        # Simulate individual user data for t-test (in production, use actual user-level data)
        control_clicks = int(control_data['clicks'].iloc[0])
        control_impressions = int(control_data['impressions'].iloc[0])
        treatment_clicks = int(treatment_data['clicks'].iloc[0])
        treatment_impressions = int(treatment_data['impressions'].iloc[0])
        
        # Perform proportion test
        from statsmodels.stats.proportion import proportions_ztest
        count = np.array([control_clicks, treatment_clicks])
        nobs = np.array([control_impressions, treatment_impressions])
        
        z_stat, p_value = proportions_ztest(count, nobs)
        
        result = {
            'control_ctr': control_ctr,
            'treatment_ctr': treatment_ctr,
            'lift': (treatment_ctr - control_ctr) / control_ctr,
            'p_value': p_value,
            'significant': p_value < 0.05,
            'confidence_interval': self._calculate_confidence_interval(control_ctr, treatment_ctr, control_impressions, treatment_impressions)
        }
        
        return result
    
    def _calculate_confidence_interval(self, control_ctr, treatment_ctr, control_n, treatment_n):
        """Calculate 95% confidence interval for lift"""
        # Simplified implementation
        lift = treatment_ctr - control_ctr
        se = np.sqrt((control_ctr * (1 - control_ctr) / control_n) + 
                    (treatment_ctr * (1 - treatment_ctr) / treatment_n))
        
        margin_of_error = 1.96 * se
        return (lift - margin_of_error, lift + margin_of_error)

if __name__ == "__main__":
    analyzer = ABTestAnalyzer("postgresql://user:pass@localhost:5432/recommendations")
    results = analyzer.run_significance_test("homepage_recommendations_v1")
    print("A/B Test Results:", json.dumps(results, indent=2, default=str))