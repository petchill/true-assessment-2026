package service

import (
	"context"
	"errors"
	"recommendation-system/src/internal/model/aggregation"
	"recommendation-system/src/internal/model/interfaces"
	"time"

	"gorm.io/gorm"
)

type recommendationService struct {
	recommendationRepository interfaces.RecommendationRepository
	userRepository           interfaces.UserRepository
}

func NewRecommendationService(
	recommendationRepository interfaces.RecommendationRepository,
	userRepository interfaces.UserRepository,
) *recommendationService {
	return &recommendationService{
		recommendationRepository: recommendationRepository,
		userRepository:           userRepository,
	}
}

func (s *recommendationService) GetUserRecommendations(ctx context.Context, userID int, limit int) (aggregation.UserRecommendationResponse, error) {
	_, err := s.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return aggregation.UserRecommendationResponse{}, errors.New("user not found")
		}

		return aggregation.UserRecommendationResponse{}, err
	}

	watchHistory, err := s.recommendationRepository.GetUserContentWatchedHistory(ctx, userID)
	if err != nil {
		return aggregation.UserRecommendationResponse{}, err
	}

	candidates, err := s.recommendationRepository.GetUserNeverSeenContent(ctx, userID)
	if err != nil {
		return aggregation.UserRecommendationResponse{}, err
	}
	scoreCal := aggregation.ScoreCalculator{}
	recommendationGenerator := aggregation.NewRecommendationGenerator(&scoreCal)
	recommendations, err := recommendationGenerator.GenerateRecommendationContents(ctx, watchHistory, candidates, limit)

	if err != nil {
		return aggregation.UserRecommendationResponse{}, err
	}

	return aggregation.UserRecommendationResponse{
		UserID:          userID,
		Recommendations: recommendations,
		Metadata: aggregation.RecommendationResponseMeta{
			CacheHit:    false,
			GeneratedAt: time.Now().UTC().Format(time.RFC3339),
			TotalCount:  len(recommendations),
		},
	}, nil
}
