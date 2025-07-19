package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/0hJonny/python-deps-crawler/internal/api-gateway/app/pb/handlers"
	"github.com/0hJonny/python-deps-crawler/internal/pkg/mocks"
	pbapi "github.com/0hJonny/python-deps-crawler/pkg/proto/api_gateway"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"
)

func getValidProtoRequest() []byte {
	req := &pbapi.AnalyzeRequest{
		UserId:        "user123",
		PythonVersion: "3.10",
		RepositoryUrl: "https://github.com/user/project",
		Packages: []*pbapi.AnalyzeRequest_RequiredPackage{
			{
				PackageName:    "requests",
				PackageVersion: "2.28.1",
			},
		},
	}
	data, _ := proto.Marshal(req)
	return data
}

func TestStartAnalysis_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockProducer := mocks.NewMockKafkaProducer()
	mockProducer.On("PublishEvent", mock.Anything, mock.AnythingOfType("*kafka_message.AnalysisStartedEvent")).Return(nil)

	mockLogger := mocks.NewMockLogger()
	mockLogger.On("WithRequestID", "test-id-123").Return(mockLogger)
	mockLogger.On("Info", mock.Anything, mock.AnythingOfType("zapcore.Field"))
	mockLogger.On("Info",
		mock.MatchedBy(func(s string) bool {
			return s == "Analysis request validated"
		}),
		mock.AnythingOfType("zapcore.Field"),
		mock.AnythingOfType("zapcore.Field"),
		mock.AnythingOfType("zapcore.Field"),
		mock.AnythingOfType("zapcore.Field"),
	).Return()

	handler := handlers.NewAnalysisHandler(mockProducer, mockLogger)

	router := gin.New()
	router.POST("/analyze", func(c *gin.Context) {
		c.Set("request_id", "test-id-123")
		c.Set("is_protobuf", true)
		c.Set("protobuf_body", getValidProtoRequest())

		handler.StartAnalysis(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/analyze", nil)
	req.Header.Set("Content-Type", "application/x-protobuf")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	respProto := &pbapi.AnalyzeResponse{}
	err := proto.Unmarshal(w.Body.Bytes(), respProto)
	assert.NoError(t, err)
	assert.NotEmpty(t, respProto.RequestId)
	assert.Equal(t, "pending", respProto.Status)
	assert.Contains(t, respProto.Message, "queued for processing")

	mockProducer.AssertExpectations(t)
}

func TestStartAnalysis_ValidationWarn(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockProducer := mocks.NewMockKafkaProducer()

	mockLogger := mocks.NewMockLogger()
	mockLogger.On("WithRequestID", "test-id-123").Return(mockLogger)
	mockLogger.On("Warn", mock.MatchedBy(func(msg string) bool {
		return strings.Contains(msg, "Non-protobuf request")
	})).Return()

	handler := handlers.NewAnalysisHandler(mockProducer, mockLogger)

	router := gin.New()
	router.POST("/analyze", func(c *gin.Context) {
		c.Set("request_id", "test-id-123")
		c.Set("is_protobuf", false)

		handler.StartAnalysis(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/analyze", nil)
	req.Header.Set("Content-Type", "application/x-protobuf")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	mockLogger.AssertCalled(t, "Warn", mock.Anything)
	assert.Contains(t, w.Body.String(), "INVALID_CONTENT_TYPE")
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStartAnalysis_ValidationBodyError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockProducer := mocks.NewMockKafkaProducer()

	mockLogger := mocks.NewMockLogger()
	mockLogger.On("WithRequestID", "test-id-123").Return(mockLogger)
	mockLogger.On("Error", mock.MatchedBy(func(msg string) bool {
		return strings.Contains(msg, "No protobuf data found")
	}),
		mock.Anything,
	).Return()

	handler := handlers.NewAnalysisHandler(mockProducer, mockLogger)

	router := gin.New()
	router.POST("/analyze", func(c *gin.Context) {
		c.Set("request_id", "test-id-123")
		c.Set("is_protobuf", true)

		handler.StartAnalysis(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/analyze", nil)
	req.Header.Set("Content-Type", "application/x-protobuf")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	mockLogger.AssertCalled(t, "Error", mock.Anything, mock.Anything)
	assert.Contains(t, w.Body.String(), "MISSING_PROTOBUF_DATA")
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStartAnalysis_ValidationProtobufError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockProducer := mocks.NewMockKafkaProducer()

	mockLogger := mocks.NewMockLogger()
	mockLogger.On("WithRequestID", "test-id-123").Return(mockLogger)
	mockLogger.On("Error", mock.MatchedBy(func(msg string) bool {
		return strings.Contains(msg, "Failed to unmarshal protobuf")
	}),
		mock.Anything,
	).Return()

	handler := handlers.NewAnalysisHandler(mockProducer, mockLogger)

	badData := []byte("Bad Protobuf!")

	router := gin.New()
	router.POST("/analyze", func(c *gin.Context) {
		c.Set("request_id", "test-id-123")
		c.Set("is_protobuf", true)
		c.Set("protobuf_body", badData)

		handler.StartAnalysis(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/analyze", nil)
	req.Header.Set("Content-Type", "application/x-protobuf")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	mockLogger.AssertCalled(t, "Error", mock.Anything, mock.Anything)
	assert.Contains(t, w.Body.String(), "PROTOBUF_UNMARSHAL_ERROR")
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStartAnalysis_ValidationRequestError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockProducer := mocks.NewMockKafkaProducer()

	mockLogger := mocks.NewMockLogger()
	mockLogger.On("WithRequestID", "test-id-123").Return(mockLogger)
	mockLogger.On("Warn", mock.MatchedBy(func(msg string) bool {
		return strings.Contains(msg, "Request validation failed")
	}),
		mock.Anything,
	).Return()

	handler := handlers.NewAnalysisHandler(mockProducer, mockLogger)

	var request pbapi.AnalyzeRequest

	data, _ := proto.Marshal(&request)

	router := gin.New()
	router.POST("/analyze", func(c *gin.Context) {
		c.Set("request_id", "test-id-123")
		c.Set("is_protobuf", true)
		c.Set("protobuf_body", data)

		handler.StartAnalysis(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/analyze", nil)
	req.Header.Set("Content-Type", "application/x-protobuf")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "VALIDATION_ERROR")
	mockLogger.AssertCalled(t, "Warn", mock.Anything, mock.Anything)
}

func TestStartAnalysis_KafkaError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockProducer := mocks.NewMockKafkaProducer()
	mockProducer.On("PublishEvent", mock.Anything, mock.AnythingOfType("*kafka_message.AnalysisStartedEvent")).Return(assert.AnError)

	mockLogger := mocks.NewMockLogger()
	mockLogger.On("WithRequestID", "test-id-123").Return(mockLogger)
	mockLogger.On("Info",
		mock.MatchedBy(func(s string) bool {
			return s == "Analysis request validated"
		}),
		mock.AnythingOfType("zapcore.Field"),
		mock.AnythingOfType("zapcore.Field"),
		mock.AnythingOfType("zapcore.Field"),
		mock.AnythingOfType("zapcore.Field"),
	).Return()
	mockLogger.On("Error", mock.MatchedBy(func(msg string) bool {
		return strings.Contains(msg, "Failed to publish event")
	}), mock.Anything, mock.Anything).Return()

	handler := handlers.NewAnalysisHandler(mockProducer, mockLogger)

	router := gin.New()
	router.POST("/analyze", func(c *gin.Context) {
		c.Set("request_id", "test-id-123")
		c.Set("is_protobuf", true)
		c.Set("protobuf_body", getValidProtoRequest())

		handler.StartAnalysis(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/analyze", nil)
	req.Header.Set("Content-Type", "application/x-protobuf")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "KAFKA_PUBLISH_ERROR")

	mockLogger.AssertCalled(t, "Error", mock.Anything, mock.Anything, mock.Anything)
	mockProducer.AssertExpectations(t)
}
