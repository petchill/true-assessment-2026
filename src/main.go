package main

import (
	"log"

	_handler "recommendation-system/src/internal/handler"
	_repository "recommendation-system/src/internal/repository"
	_service "recommendation-system/src/internal/service"
	"recommendation-system/src/utils/config"
	database "recommendation-system/src/utils/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewGormDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	redisClient, err := database.NewRedisClient(cfg.RedisURL)
	if err != nil {
		log.Fatal(err)
	}

	recommendationRepository := _repository.NewRecommendationRepository(db, redisClient)
	recommendationService := _service.NewRecommendationService(recommendationRepository)
	recommendationHandler := _handler.NewRecommendationHandler(recommendationService)
	recommendationHandler.RegisterRoutes()
}
