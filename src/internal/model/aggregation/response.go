package aggregation

type UserRecommendationResponse struct {
	UserID          int                        `json:"user_id"`
	Recommendations []RecommendationItem       `json:"recommendations"`
	Metadata        RecommendationResponseMeta `json:"metadata"`
}

type RecommendationItem struct {
	ContentID       uint64  `json:"content_id"`
	Title           string  `json:"title"`
	Genre           string  `json:"genre"`
	PopularityScore float64 `json:"popularity_score"`
	Score           float64 `json:"score"`
}

type RecommendationResponseMeta struct {
	CacheHit    bool   `json:"cache_hit"`
	GeneratedAt string `json:"generated_at"`
	TotalCount  int    `json:"total_count"`
}
