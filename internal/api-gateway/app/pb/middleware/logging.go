package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/0hJonny/python-deps-crawler/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ZapLoggingMiddleware(logger *logger.Logger) gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(params gin.LogFormatterParams) string {
			fields := []zap.Field{
				zap.Int("status_code", params.StatusCode),
				zap.Duration("latency", params.Latency),
				zap.String("client_ip", params.ClientIP),
				zap.String("method", params.Method),
				zap.String("path", params.Path),
				zap.String("user_agent", params.Request.UserAgent()),
				zap.Int("body_size", params.BodySize),
			}

			if requestID := params.Request.Header.Get("X-Request-ID"); requestID != "" {
				fields = append(fields, zap.String("request_id", requestID))
			}

			if params.ErrorMessage != "" {
				fields = append(fields, zap.String("error", params.ErrorMessage))
			}

			switch {
			case params.StatusCode >= 500:
				logger.Error("HTTP Request", fields...)
			case params.StatusCode >= 400:
				logger.Warn("HTTP Request", fields...)
			default:
				logger.Info("HTTP Request", fields...)
			}

			return ""
		},
		Output: gin.DefaultWriter,
	})
}

func ZapRecoveryMiddleware(logger *logger.Logger) gin.HandlerFunc {
	return gin.RecoveryWithWriter(gin.DefaultErrorWriter, func(c *gin.Context, rec any) {
		if err, ok := rec.(string); ok {
			logger.Error("Panic recovered",
				zap.String("error", err),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.String("client_ip", c.ClientIP()),
			)
		}
		c.AbortWithStatus(500)
	})
}

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Header("X-Request-ID", requestID)

		c.Set("request_id", requestID)

		c.Next()
	}
}

func generateRequestID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return fmt.Sprintf("%d-%s", time.Now().UnixNano(), hex.EncodeToString(bytes))
}
