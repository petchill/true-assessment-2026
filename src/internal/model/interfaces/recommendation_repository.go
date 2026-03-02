package interfaces

import (
	"context"
	"recommendation-system/src/internal/model/aggregation"
	"recommendation-system/src/internal/model/entity"
)

type RecommendationRepository interface {
	GetUserRecommendationCache(ctx context.Context, userID int, limit int) (aggregation.UserRecommendationResponse, bool, error)
	SetUserRecommendationCache(ctx context.Context, userID int, limit int, response aggregation.UserRecommendationResponse) error
	TruncateUserRecommendationCache(ctx context.Context, userID int) error
	InsertUserWatchHistory(ctx context.Context, userID int, contentID int) error
	GetUserNeverSeenContent(ctx context.Context, userID int) ([]entity.Content, error)
	GetUserContentWatchedHistory(ctx context.Context, userID int) ([]entity.HistoryContent, error)
}
