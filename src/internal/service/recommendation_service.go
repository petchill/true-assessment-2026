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

func (s *recommendationService) getUserRecommendations(ctx context.Context, userID int, limit int) ([]aggregation.RecommendationContent, error) {

	watchHistory, err := s.recommendationRepository.GetUserContentWatchedHistory(ctx, userID)
	if err != nil {
		return []aggregation.RecommendationContent{}, err
	}

	candidates, err := s.recommendationRepository.GetUserNeverSeenContent(ctx, userID)
	if err != nil {
		return []aggregation.RecommendationContent{}, err
	}
	scoreCal := aggregation.ScoreCalculator{}
	recommendationGenerator := aggregation.NewRecommendationGenerator(&scoreCal)
	return recommendationGenerator.GenerateRecommendationContents(ctx, watchHistory, candidates, limit)
}

func (s *recommendationService) GetUserRecommendations(ctx context.Context, userID int, limit int) (aggregation.UserRecommendationResponse, error) {
	_, err := s.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return aggregation.UserRecommendationResponse{}, errors.New("user not found")
		}

		return aggregation.UserRecommendationResponse{}, err
	}
	recommendations, err := s.getUserRecommendations(ctx, userID, limit)

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

func (s *recommendationService) GetBatchRecommendation(ctx context.Context, page, limit int) (aggregation.BatchRecommendationResponse, error) {
	userIDs, err := s.userRepository.GetAllUserIDs(ctx)
	if err != nil {
		return aggregation.BatchRecommendationResponse{}, err
	}

	startTime := time.Now()

	filteredUserIDs := getPaginateUserIDs(userIDs, page, limit)

	results := s.getBatchResultFromUserIDs(ctx, filteredUserIDs)

	processingTime := time.Since(startTime)
	summary := getSummary(results, processingTime)

	return aggregation.BatchRecommendationResponse{
		Page:       page,
		Limit:      limit,
		TotalUsers: len(userIDs),
		Results:    results,
		Summary:    summary,
		Metadata: aggregation.BatchRecommendationResponseMeta{
			GeneratedAt: time.Now().Format(time.RFC3339),
		},
	}, nil
}

func getSummary(results []aggregation.BatchRecommendationResult, processingTime time.Duration) aggregation.BatchRecommendationSummary {
	successCount, failedCount := int(0), int(0)
	for _, r := range results {
		switch r.Status {
		case aggregation.BatchStatusSuccess:
			{
				successCount += 1

			}
		case aggregation.BatchStatusFailed:
			{
				failedCount += 1
			}
		}

	}
	return aggregation.BatchRecommendationSummary{
		ProcessingTimeMs: int(processingTime.Milliseconds()),
		SuccessCount:     successCount,
		FailedCount:      failedCount,
	}
}

func (s *recommendationService) getBatchResultFromUserIDs(ctx context.Context, userIDs []int) []aggregation.BatchRecommendationResult {
	result := []aggregation.BatchRecommendationResult{}
	for _, userID := range userIDs {
		recommendations, err := s.getUserRecommendations(ctx, userID, 50)
		if err != nil {
			batch := aggregation.BatchRecommendationResult{
				UserID:  userID,
				Status:  aggregation.BatchStatusFailed,
				Error:   err.Error(),
				Message: err.Error(), // TODO: change it later
			}
			result = append(result, batch)
			continue
		}
		batch := aggregation.BatchRecommendationResult{
			UserID:          userID,
			Status:          aggregation.BatchStatusSuccess,
			Recommendations: recommendations,
		}
		result = append(result, batch)
	}
	return result
}

func getPaginateUserIDs(userIDs []int, page, limit int) []int {
	if page <= 0 || limit <= 0 {
		return []int{}
	}

	offset := (page - 1) * limit
	if offset >= len(userIDs) {
		return []int{}
	}

	end := offset + limit
	if end > len(userIDs) {
		end = len(userIDs)
	}

	return userIDs[offset:end]
}
