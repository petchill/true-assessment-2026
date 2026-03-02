package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"recommendation-system/src/internal/model/aggregation"
	"recommendation-system/src/internal/model/entity"
	"time"

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

func (r *recommendationRepository) GetUserRecommendationCache(ctx context.Context, userID int, limit int) (aggregation.UserRecommendationResponse, bool, error) {
	// set the time out 10 sec
	ctx, _ = context.WithTimeout(ctx, time.Second*5)

	key := fmt.Sprintf("rec:user:%d:limit:%d", userID, limit)
	raw, err := r.redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return aggregation.UserRecommendationResponse{}, false, nil
		}
		return aggregation.UserRecommendationResponse{}, false, err
	}

	response := aggregation.UserRecommendationResponse{}
	if err := json.Unmarshal([]byte(raw), &response); err != nil {
		return aggregation.UserRecommendationResponse{}, false, err
	}

	response.Metadata.CacheHit = true

	return response, true, nil
}

func (r *recommendationRepository) SetUserRecommendationCache(ctx context.Context, userID int, limit int, response aggregation.UserRecommendationResponse) error {
	key := fmt.Sprintf("rec:user:%d:limit:%d", userID, limit)
	raw, err := json.Marshal(response)
	if err != nil {
		return err
	}

	return r.redisClient.Set(ctx, key, raw, time.Minute*10).Err()
}

func (r *recommendationRepository) TruncateUserRecommendationCache(ctx context.Context, userID int) error {
	pattern := fmt.Sprintf("rec:user:%d:limit:*", userID)
	keys, err := r.redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}

	return r.redisClient.Del(ctx, keys...).Err()
}

func (r *recommendationRepository) InsertUserWatchHistory(ctx context.Context, userID int, contentID int) error {
	return r.db.WithContext(ctx).
		Table("user_watch_history").
		Create(map[string]any{
			"user_id":    userID,
			"content_id": contentID,
			"watched_at": time.Now().UTC(),
		}).Error
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
		log.Println("Error => ", err)
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
		log.Println("Error => ", err)
		return nil, err
	}

	return historyContents, nil
}
