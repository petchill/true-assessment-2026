package service

import (
	"context"
	"errors"
	"log"
	"recommendation-system/src/internal/model/aggregation"
	"recommendation-system/src/internal/model/interfaces"
	_appErr "recommendation-system/src/utils/error"
	"sync"
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
	chacheReponse, found, err := s.recommendationRepository.GetUserRecommendationCache(ctx, userID, limit)
	if err != nil {
		log.Println("Warning: Get User Recommendation from cache error due to ", err.Error())
	}
	if found {
		return chacheReponse, nil
	}

	// if no cache
	_, err = s.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return aggregation.UserRecommendationResponse{}, errors.New("user not found")
		}

		return aggregation.UserRecommendationResponse{}, err
	}
	recommendations, err := s.getUserRecommendations(ctx, userID, limit)
	if err != nil {
		return aggregation.UserRecommendationResponse{}, err
	}

	response := aggregation.UserRecommendationResponse{
		UserID:          userID,
		Recommendations: recommendations,
		Metadata: aggregation.RecommendationResponseMeta{
			CacheHit:    false,
			GeneratedAt: time.Now().UTC().Format(time.RFC3339),
			TotalCount:  len(recommendations),
		},
	}

	err = s.recommendationRepository.SetUserRecommendationCache(ctx, userID, limit, response)
	if err != nil {
		log.Println("Warning: Set User Recommendation to cache error due to ", err.Error())
	}

	return response, nil
}

func (s *recommendationService) InsertUserWatchHistory(ctx context.Context, userID int, contentID int) error {
	if err := s.recommendationRepository.TruncateUserRecommendationCache(ctx, userID); err != nil {
		return err
	}

	return s.recommendationRepository.InsertUserWatchHistory(ctx, userID, contentID)
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
	ch := make(chan aggregation.BatchRecommendationResult)
	var wg sync.WaitGroup
	for _, userID := range userIDs {
		wg.Add(1)
		go func() {
			finish := make(chan aggregation.BatchRecommendationResult, 1)
			// timeout 2 sec
			ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
			go func() {
				recommendations, err := s.getUserRecommendations(ctx, userID, 50)
				if err != nil {
					batch := aggregation.BatchRecommendationResult{
						UserID:  userID,
						Status:  aggregation.BatchStatusFailed,
						Error:   err.Error(),
						Message: err.Error(), // TODO: change it later
					}
					finish <- batch
				} else {
					batch := aggregation.BatchRecommendationResult{
						UserID:          userID,
						Status:          aggregation.BatchStatusSuccess,
						Recommendations: recommendations,
					}
					finish <- batch
				}
			}()
			select {
			case <-ctx.Done():
				{
					batch := aggregation.BatchRecommendationResult{
						UserID:  userID,
						Status:  aggregation.BatchStatusFailed,
						Error:   "model_inference_timeout",
						Message: _appErr.ErrModelInferenceTimeout.Message(), // TODO: change it later
					}
					ch <- batch
				}
			case batch := <-finish:
				{
					ch <- batch
					wg.Done()
				}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	for v := range ch {
		result = append(result, v)
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
