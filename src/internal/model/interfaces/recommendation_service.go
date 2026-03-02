package interfaces

import (
	"context"
	"recommendation-system/src/internal/model/aggregation"
)

type RecommendationService interface {
	GetUserRecommendations(ctx context.Context, userID int, limit int) (aggregation.UserRecommendationResponse, error)
	GetBatchRecommendation(ctx context.Context, page, limit int) (aggregation.BatchRecommendationResponse, error)
	InsertUserWatchHistory(ctx context.Context, userID int, contentID int) error
}
