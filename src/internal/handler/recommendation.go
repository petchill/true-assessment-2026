package handler

import (
	"recommendation-system/src/internal/model/interfaces"

	"github.com/labstack/echo/v5"
)

type recommendationHandler struct {
	recommendationService interfaces.RecommendationService
}

func NewRecommendationHandler(recommendationService interfaces.RecommendationService) *recommendationHandler {
	return &recommendationHandler{
		recommendationService,
	}
}

func (h *recommendationHandler) RegisterRoutes(e *echo.Group) {

}
