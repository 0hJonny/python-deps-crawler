package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/0hJonny/python-deps-crawler/internal/api-gateway/app/pb/middleware"
	"github.com/0hJonny/python-deps-crawler/internal/api-gateway/kafka"
	"github.com/0hJonny/python-deps-crawler/internal/pkg/logger"
	pbapi "github.com/0hJonny/python-deps-crawler/pkg/proto/api_gateway"
	eventspb "github.com/0hJonny/python-deps-crawler/pkg/proto/api_gateway_kafka_events"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-uuid"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AnalysisHandler struct {
	kafkaProducer *kafka.APIGatewayProducer
	logger        *logger.Logger
}

func NewAnalysisHandler(kafkaProducer *kafka.APIGatewayProducer, logger *logger.Logger) *AnalysisHandler {
	return &AnalysisHandler{
		kafkaProducer: kafkaProducer,
		logger:        logger,
	}
}

func (h *AnalysisHandler) StartAnalysis(c *gin.Context) {
	requestID := c.GetString("request_id")

	contextLogger := h.logger.WithRequestID(requestID)

	isPb, ok := c.Get("is_protobuf")
	if !ok || !isPb.(bool) {
		contextLogger.Warn("Non-protobuf request received")
		middleware.SendProtobufError(c, http.StatusBadRequest,
			"Expected protobuf content-type", "INVALID_CONTENT_TYPE")
		return
	}

	pbBody, ok := c.Get("protobuf_body")
	if !ok {
		contextLogger.Error("No protobuf data found")
		middleware.SendProtobufError(c, http.StatusBadRequest,
			"No protobuf data found", "MISSING_PROTOBUF_DATA")
		return
	}

	var request pbapi.AnalyzeRequest
	if err := proto.Unmarshal(pbBody.([]byte), &request); err != nil {
		contextLogger.Error("Failed to unmarshal protobuf", zap.Error(err))
		middleware.SendProtobufError(c, http.StatusBadRequest,
			"Invalid protobuf message", "PROTOBUF_UNMARSHAL_ERROR")
		return
	}

	if err := h.validateRequest(&request); err != nil {
		contextLogger.Warn("Request validation failed", zap.Error(err))
		middleware.SendProtobufError(c, http.StatusBadRequest,
			err.Error(), "VALIDATION_ERROR")
		return
	}

	analysisID, err := uuid.GenerateUUID()
	if err != nil {
		contextLogger.Warn("UUID generate failed", zap.Error(err))
		middleware.SendProtobufError(c, http.StatusBadRequest,
			err.Error(), "UUID_GENERATE_ERROR")
		return
	}

	contextLogger.Info("Analysis request validated",
		zap.String("analysis_id", analysisID),
		zap.String("user_id", request.UserId),
		zap.String("python_version", request.PythonVersion),
		zap.Int("packages_count", len(request.Packages)),
	)

	event := &eventspb.AnalysisStartedEvent{
		RequestId:     analysisID,
		UserId:        request.UserId,
		PythonVersion: request.PythonVersion,
		RepositoryUrl: request.RepositoryUrl,
		Packages:      h.convertPackages(request.Packages),
		Timestamp:     timestamppb.Now(),
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	if err := h.kafkaProducer.PublishEvent(ctx, event); err != nil {
		contextLogger.Error("Failed to publish event to Kafka",
			zap.String("analysis_id", analysisID),
			zap.Error(err),
		)
		middleware.SendProtobufError(c, http.StatusInternalServerError,
			"Failed to publish event", "KAFKA_PUBLISH_ERROR")
		return
	}

	contextLogger.Info("Event published to Kafka successfully",
		zap.String("analysis_id", analysisID),
	)

	response := &pbapi.AnalyzeResponse{
		RequestId: analysisID,
		Status:    "pending",
		Message:   "Analysis request received and queued for processing",
		CreatedAt: timestamppb.Now(),
	}

	middleware.SendProtobufResponse(c, response)
}

func (h *AnalysisHandler) convertPackages(apiPackages []*pbapi.AnalyzeRequest_RequiredPackage) []*eventspb.AnalysisStartedEvent_RequiredPackage {
	eventPackages := make([]*eventspb.AnalysisStartedEvent_RequiredPackage, len(apiPackages))
	for i, pkg := range apiPackages {
		eventPackages[i] = &eventspb.AnalysisStartedEvent_RequiredPackage{
			PackageName:    pkg.PackageName,
			PackageVersion: pkg.PackageVersion,
			Extras:         pkg.Extras,
		}
	}
	return eventPackages
}

func (h *AnalysisHandler) validateRequest(req *pbapi.AnalyzeRequest) error {
	if req.UserId == "" {
		return fmt.Errorf("user_id is required")
	}
	if req.PythonVersion == "" {
		return fmt.Errorf("python_version is required")
	}
	if len(req.Packages) == 0 {
		return fmt.Errorf("at least one package is required")
	}

	for i, pkg := range req.Packages {
		if pkg.PackageName == "" {
			return fmt.Errorf("package_name is required for package %d", i)
		}
	}

	return nil
}
