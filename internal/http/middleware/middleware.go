package middleware

import (
	"log/slog"
	"time"

	"go-boilerplate/internal/http/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIDHeader = "X-Request-ID"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(requestIDHeader)
		if requestID == "" {
			requestID = newRequestID()
		}
		c.Header(requestIDHeader, requestID)
		c.Next()
	}
}

func Logger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startedAt := time.Now()
		c.Next()
		logger.Info("http request",
			slog.String("requestId", getRequestID(c)),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", c.Writer.Status()),
			slog.String("latency", time.Since(startedAt).String()),
			slog.String("ip", c.ClientIP()),
		)
	}
}

func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				logger.Error("panic recovered",
					slog.String("requestId", getRequestID(c)),
					slog.Any("error", recovered),
				)
				response.Error(c, response.InternalServerError("internal server error"))
				c.Abort()
			}
		}()
		c.Next()
	}
}

func newRequestID() string {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.NewString()
	}
	return id.String()
}

func getRequestID(c *gin.Context) string {
	return c.Writer.Header().Get(requestIDHeader)
}
