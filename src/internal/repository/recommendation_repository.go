package repository

import (
	"context"
	"recommendation-system/src/internal/model/entity"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type recommendationRepository struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func NewRecommendationRepository(db *gorm.DB, redisClient *redis.Client) *recommendationRepository {
	return &recommendationRepository{
		db:          db,
		redisClient: redisClient,
	}
}

func (r *recommendationRepository) GetUserNeverSeenContent(ctx context.Context, userID int) ([]entity.Content, error) {
	contents := []entity.Content{}
	subQuery := r.db.WithContext(ctx).
		Table("user_watch_history").
		Select("content_id").
		Where("user_id = ?", userID)

	err := r.db.WithContext(ctx).
		Table("content").
		Select("id, title, genre, popularity_score, created_at").
		Where("id NOT IN (?)", subQuery).
		Order("popularity_score DESC").
		Limit(100).
		Find(&contents).Error
	if err != nil {
		return []entity.Content{}, err
	}

	return contents, nil
}

func (r *recommendationRepository) GetUserContentWatchedHistory(ctx context.Context, userID int) ([]entity.HistoryContent, error) {
	historyContents := []entity.HistoryContent{}

	err := r.db.WithContext(ctx).
		Table("user_watch_history AS uwh").
		Select("c.id, c.genre, uwh.watched_at").
		Joins("JOIN content AS c ON c.id = uwh.content_id").
		Where("uwh.user_id = ?", userID).
		Order("uwh.watched_at DESC").
		Limit(50).
		Find(&historyContents).Error
	if err != nil {
		return nil, err
	}

	return historyContents, nil
}
