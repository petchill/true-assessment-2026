package main

import (
	_handler "recommendation-system/src/internal/handler"
	_repository "recommendation-system/src/internal/repository"
	_service "recommendation-system/src/internal/service"
)

func main() {
	recommendationRepository := _repository.NewRecommendationRepository()
	recommendationService := _service.NewRecommendationService(recommendationRepository)
	recommendationHandler := _handler.NewRecommendationHandler(recommendationService)
	recommendationHandler.RegisterRoutes()
}
