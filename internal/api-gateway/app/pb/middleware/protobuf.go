package middleware

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

// ProtobufMiddleware обрабатывает protobuf запросы
func ProtobufMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверяем только POST/PUT/PATCH запросы
		if c.Request.Method != http.MethodPost &&
			c.Request.Method != http.MethodPut &&
			c.Request.Method != http.MethodPatch {
			c.Next()
			return
		}

		contentType := c.GetHeader("Content-Type")

		// Проверяем protobuf content-type
		if strings.Contains(contentType, "application/x-protobuf") ||
			strings.Contains(contentType, "application/protobuf") {

			// Читаем тело запроса
			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Failed to read request body",
					"code":  "INVALID_BODY",
				})
				c.Abort()
				return
			}

			// Проверяем, что тело не пустое
			if len(body) == 0 {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Empty request body",
					"code":  "EMPTY_BODY",
				})
				c.Abort()
				return
			}

			// Восстанавливаем тело для повторного чтения
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

			// Сохраняем protobuf данные в контекст
			c.Set("protobuf_body", body)
			c.Set("is_protobuf", true)
		}

		c.Next()
	}
}

// SendProtobufResponse отправляет protobuf ответ
func SendProtobufResponse(c *gin.Context, message proto.Message) {
	data, err := proto.Marshal(message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to serialize response",
			"code":  "SERIALIZATION_ERROR",
		})
		return
	}

	c.Header("Content-Type", "application/x-protobuf")
	c.Data(http.StatusOK, "application/x-protobuf", data)
}

// SendProtobufError отправляет protobuf ошибку
func SendProtobufError(c *gin.Context, statusCode int, message string, code string) {
	c.Header("Content-Type", "application/json")
	c.JSON(statusCode, gin.H{
		"error": message,
		"code":  code,
	})
}
