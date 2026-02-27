package service

import (
	"recommendation-system/internal/model/interfaces"
)

type recommendationService struct {
	recommendationRepository interfaces.RecommendationRepository
}

func NewRecommendationService(recommendationRepository interfaces.RecommendationRepository) *recommendationService {
	return &recommendationService{
		recommendationRepository,
	}
}
