package service

import (
	"context"
	"recommendation-system/src/internal/model/aggregation"
	"recommendation-system/src/internal/model/interfaces"
)

type recommendationService struct {
	recommendationRepository interfaces.RecommendationRepository
}

func NewRecommendationService(recommendationRepository interfaces.RecommendationRepository) *recommendationService {
	return &recommendationService{
		recommendationRepository,
	}
}

func (s *recommendationService) GetUserRecommendations(ctx context.Context, userID int) (aggregation.UserRecommendationResponse, error) {
	return aggregation.UserRecommendationResponse{}, nil
}
