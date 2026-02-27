package handler

import (
	"recommendation-system/internal/model/interfaces"
)

type recommendationHandler struct {
	recommendationService interfaces.RecommendationService
}

func NewRecommendationHandler(recommendationService interfaces.RecommendationService) *recommendationHandler {
	return &recommendationHandler{
		recommendationService,
	}
}

func (h *recommendationHandler) RegisterRoutes() {

}
