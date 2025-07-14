package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/0hJonny/python-deps-crawler/internal/api-gateway/app/pb/handlers"
	"github.com/0hJonny/python-deps-crawler/internal/api-gateway/app/pb/routes"
	"github.com/0hJonny/python-deps-crawler/internal/api-gateway/kafka"
	"github.com/0hJonny/python-deps-crawler/internal/pkg/config"
	"github.com/0hJonny/python-deps-crawler/internal/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	tLogg, _ := zap.NewDevelopment()
	defer tLogg.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		tLogg.Fatal("Failed to load config", zap.Error(err))
	}

	logger, err := logger.NewLogger(&cfg.Logger)
	if err != nil {
		tLogg.Fatal("Failed to initialize logger", zap.Error(err))
	}
	defer logger.Sync()

	logger.Info("Starting API Gateway",
		zap.String("version", "1.0.0"),
		zap.String("env", cfg.Server.Mode),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kafkaProducer, err := initKafkaProducer(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to initialize Kafka producer", zap.Error(err))
	}
	defer func() {
		logger.Info("Closing Kafka producer")
		if err := kafkaProducer.Close(); err != nil {
			logger.Error("Error closing Kafka producer", zap.Error(err))
		}
	}()

	analysisHandler := handlers.NewAnalysisHandler(kafkaProducer, logger)
	healthHandler := handlers.NewHealthHandler(logger)

	router := routes.SetupRoutes(analysisHandler, healthHandler, cfg, logger)

	server := &http.Server{
		Addr:         cfg.Server.GetConfig(),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		logger.Info("HTTP server starting",
			zap.String("address", server.Addr),
			zap.Duration("read_timeout", cfg.Server.ReadTimeout),
			zap.Duration("write_timeout", cfg.Server.WriteTimeout),
		)

		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	logger.Info("API Gateway started successfully",
		zap.Strings("kafka_brokers", cfg.Kafka.Brokers),
		zap.String("kafka_topic", cfg.Kafka.Topic),
		zap.String("server_mode", cfg.Server.Mode),
	)

	gracefulShutdown(ctx, server, cfg, logger)
}

func gracefulShutdown(ctx context.Context, server *http.Server, cfg *config.Config, logger *logger.Logger) {
	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	sig := <-term
	logger.Info("Shutdown signal received", zap.String("signal", sig.String()))

	shutDownCtx, cancel := context.WithTimeout(ctx, cfg.Server.ShutdownTimeout)
	defer cancel()

	logger.Info("Shutting down HTTP server...")

	if err := server.Shutdown(shutDownCtx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
		return
	}
	logger.Info("Server exited gracefully")
}

func initKafkaProducer(cfg *config.Config, logger *logger.Logger) (*kafka.APIGatewayProducer, error) {
	logger.Info("Connecting to Kafka",
		zap.Strings("brokers", cfg.Kafka.Brokers),
		zap.String("topic", cfg.Kafka.Topic),
	)

	producer, err := kafka.NewAPIGatewayProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	logger.Info("Kafka producer initialized successfully")
	return producer, nil
}
