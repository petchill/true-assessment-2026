package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_handler "recommendation-system/src/internal/handler"
	_repository "recommendation-system/src/internal/repository"
	_service "recommendation-system/src/internal/service"
	"recommendation-system/src/utils"
	"recommendation-system/src/utils/config"
	database "recommendation-system/src/utils/database"

	"github.com/labstack/echo/v5"
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

	e := utils.InitEchoApp()

	recommendationRepository := _repository.NewRecommendationRepository(db, redisClient)
	recommendationService := _service.NewRecommendationService(recommendationRepository)
	recommendationHandler := _handler.NewRecommendationHandler(recommendationService)
	recommendationHandler.RegisterRoutes(e.Group(""))

	e.GET("/healthz", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	sc := echo.StartConfig{
		Address:         ":8080",
		GracefulTimeout: 5 * time.Second,
	}
	if err := sc.Start(ctx, e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
