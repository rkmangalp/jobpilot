package dashboard

type Metrics struct {
	JobsDiscovered        int     `json:"jobs_discovered"`
	JobsAnalyzed          int     `json:"jobs_analyzed"`
	ApplicationsPrepared  int     `json:"applications_prepared"`
	ApplicationsSubmitted int     `json:"applications_submitted"`
	Interviews            int     `json:"interviews"`
	Offers                int     `json:"offers"`
	AverageMatchScore     float64 `json:"average_match_score"`
	ApplicationsPerWeek   int     `json:"applications_per_week"`
}

func Empty() Metrics { return Metrics{} }
