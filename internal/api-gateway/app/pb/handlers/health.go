package handlers

import (
	"net/http"
	"time"

	"github.com/0hJonny/python-deps-crawler/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HealthHandler struct {
	logger logger.LoggerInterface
}

func NewHealthHandler(logger logger.LoggerInterface) *HealthHandler {
	return &HealthHandler{
		logger: logger,
	}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	resourse := gin.H{
		"status":    "healthy",
		"service":   "api-gateway",
		"timestamp": time.Now().Unix(),
		"version":   "1.0.0",
	}

	h.logger.Debug("Health check requested",
		zap.String("client_ip", c.ClientIP()),
	)

	c.JSON(http.StatusOK, resourse)
}

func (h *HealthHandler) LiveCheck(c *gin.Context) {
	resourse := gin.H{
		"status":    "alive",
		"service":   "api-gateway",
		"timestamp": time.Now().Unix(),
	}

	c.JSON(http.StatusOK, resourse)
}
