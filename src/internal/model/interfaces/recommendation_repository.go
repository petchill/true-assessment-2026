package interfaces

import (
	"context"
	"recommendation-system/src/internal/model/entity"
)

type RecommendationRepository interface {
	GetUserNeverSeenContent(ctx context.Context, userID int) ([]entity.Content, error)
	GetUserContentWatchedHistory(ctx context.Context, userID int) ([]entity.HistoryContent, error)
}
