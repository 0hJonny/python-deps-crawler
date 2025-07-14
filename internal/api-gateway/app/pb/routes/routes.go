package routes

import (
	"github.com/0hJonny/python-deps-crawler/internal/api-gateway/app/pb/handlers"
	"github.com/0hJonny/python-deps-crawler/internal/api-gateway/app/pb/middleware"
	"github.com/0hJonny/python-deps-crawler/internal/pkg/config"
	"github.com/0hJonny/python-deps-crawler/internal/pkg/logger"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	analysisHandler *handlers.AnalysisHandler,
	healthHandler *handlers.HealthHandler,
	cfg *config.Config,
	logger *logger.Logger,
) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)

	router := gin.New()

	router.Use(middleware.ZapRecoveryMiddleware(logger))
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.ZapLoggingMiddleware(logger))
	router.Use(middleware.ProtobufMiddleware())

	setupHealthRoutes(router, healthHandler)

	v1 := router.Group("/api/v1")
	{
		setupAnalysisRoutes(v1, analysisHandler)
	}

	return router
}

func setupAnalysisRoutes(group *gin.RouterGroup, analysisHandler *handlers.AnalysisHandler) {
	analysis := group.Group("/analysis")
	{
		analysis.POST("/start", analysisHandler.StartAnalysis)
		analysis.POST("", analysisHandler.StartAnalysis)
	}

	group.POST("/analyze", analysisHandler.StartAnalysis)
}

func setupHealthRoutes(router *gin.Engine, healthHandler *handlers.HealthHandler) {
	health := router.Group("/health")
	{
		health.GET("", healthHandler.HealthCheck)
		health.GET("/", healthHandler.HealthCheck)
		health.GET("/live", healthHandler.LiveCheck)
	}

	router.GET("/", healthHandler.HealthCheck)
}
