package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/TorekhanUssembay/subscription_service/docs"

	"github.com/TorekhanUssembay/subscription_service/internal/config"
	"github.com/TorekhanUssembay/subscription_service/internal/repository"
	"github.com/TorekhanUssembay/subscription_service/internal/service"
	"github.com/TorekhanUssembay/subscription_service/internal/handler"
)

// @title Subscription Service API
// @version 1.0
// @description API for managing subscriptions
// @host localhost:8080
// @BasePath /
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbPool, err := repository.NewPostgres(ctx, cfg.DatabaseDSN())
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %v", err)
	}
	defer dbPool.Close()

	log.Printf("Connected to DB %s", cfg.DBName)

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	repo := repository.NewSubscriptionRepo(dbPool)
	svc := service.NewSubscriptionService(repo)
	h := handler.NewSubscriptionHandler(svc)

	r.POST("/subscriptions", h.CreateSubscription)
	r.GET("/subscriptions/:id", h.GetSubscription)
	r.PUT("/subscriptions/:id", h.UpdateSubscription)
	r.DELETE("/subscriptions/:id", h.DeleteSubscription)
	r.GET("/subscriptions", h.ListSubscriptions)
	r.GET("/subscriptions/sum", h.SumSubscriptions)

	log.Printf("Server started on port %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}