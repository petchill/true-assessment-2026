package repository

import (
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
