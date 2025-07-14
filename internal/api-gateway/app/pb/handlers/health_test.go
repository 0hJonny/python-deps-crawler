package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/0hJonny/python-deps-crawler/internal/pkg/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupHealthTestRouter() (*gin.Engine, *mocks.MockLogger) {
	gin.SetMode(gin.TestMode)

	mockLogger := mocks.NewMockLogger()
	handler := NewHealthHandler(mockLogger)

	router := gin.New()
	router.GET("/health", handler.HealthCheck)
	router.GET("/health/live", handler.LiveCheck)

	return router, mockLogger
}

func TestHealthHandler_HealthCheck(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedFields map[string]interface{}
	}{
		{
			name:           "успешный health check",
			path:           "/health",
			expectedStatus: http.StatusOK,
			expectedFields: map[string]interface{}{
				"status":  "healthy",
				"service": "api-gateway",
				"version": "1.0.0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			router, mockLogger := setupHealthTestRouter()
			mockLogger.On("Debug", mock.AnythingOfType("string"), mock.Anything).Return()

			// Act
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			for key, expectedValue := range tt.expectedFields {
				assert.Equal(t, expectedValue, response[key])
			}

			// Проверяем, что timestamp существует и является числом
			timestamp, exists := response["timestamp"]
			assert.True(t, exists)
			assert.IsType(t, float64(0), timestamp)

			mockLogger.AssertExpectations(t)
		})
	}
}

func TestHealthHandler_LiveCheck(t *testing.T) {
	tests := []struct {
		name           string
		expectedStatus int
		expectedFields map[string]interface{}
	}{
		{
			name:           "успешный liveness check",
			expectedStatus: http.StatusOK,
			expectedFields: map[string]interface{}{
				"status":  "alive",
				"service": "api-gateway",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			router, _ := setupHealthTestRouter()

			// Act
			req := httptest.NewRequest(http.MethodGet, "/health/live", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			for key, expectedValue := range tt.expectedFields {
				assert.Equal(t, expectedValue, response[key])
			}

			// Проверяем timestamp
			timestamp, exists := response["timestamp"]
			assert.True(t, exists)
			assert.IsType(t, float64(0), timestamp)
		})
	}
}

func TestHealthHandler_HTTPMethods(t *testing.T) {
	router, mockLogger := setupHealthTestRouter()

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "GET health endpoint",
			method:         http.MethodGet,
			path:           "/health",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST health endpoint (не поддерживается)",
			method:         http.MethodPost,
			path:           "/health",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "PUT health endpoint (не поддерживается)",
			method:         http.MethodPut,
			path:           "/health",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.method == http.MethodGet {
				mockLogger.On("Debug", mock.AnythingOfType("string"), mock.Anything).Return()
			}

			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

// Benchmark тесты для health endpoints
func BenchmarkHealthHandler_HealthCheck(b *testing.B) {
	router, mockLogger := setupHealthTestRouter()
	mockLogger.On("Debug", mock.AnythingOfType("string"), mock.Anything).Return()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

func BenchmarkHealthHandler_LiveCheck(b *testing.B) {
	router, _ := setupHealthTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/health/live", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
