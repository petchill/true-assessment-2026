package aggregation

type BatchStatus string

const (
	BatchStatusSuccess BatchStatus = "success"
	BatchStatusFailed  BatchStatus = "failed"
)

type BatchRecommendationResponse struct {
	Page       int                             `json:"page"`
	Limit      int                             `json:"limit"`
	TotalUsers int                             `json:"total_users"`
	Results    []BatchRecommendationResult     `json:"results"`
	Summary    BatchRecommendationSummary      `json:"summary"`
	Metadata   BatchRecommendationResponseMeta `json:"metadata"`
}

type BatchRecommendationResult struct {
	UserID          int                     `json:"user_id"`
	Recommendations []RecommendationContent `json:"recommendations,omitempty"`
	Status          BatchStatus             `json:"status"`
	Error           string                  `json:"error,omitempty"`
	Message         string                  `json:"message,omitempty"`
}

type BatchRecommendationSummary struct {
	SuccessCount     int `json:"success_count"`
	FailedCount      int `json:"failed_count"`
	ProcessingTimeMs int `json:"processing_time_ms"`
}

type BatchRecommendationResponseMeta struct {
	GeneratedAt string `json:"generated_at"`
}
