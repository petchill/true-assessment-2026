package main

import (
	_handler "recommendation-system/internal/handler"
	_repository "recommendation-system/internal/repository"
	_service "recommendation-system/internal/service"
)

func main() {
	recommendationRepository := _repository.NewRecommendationRepository()
	recommendationService := _service.NewRecommendationService(recommendationRepository)
	recommendationHandler := _handler.NewRecommendationHandler(recommendationService)
	recommendationHandler.RegisterRoutes()
}
