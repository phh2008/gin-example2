package middleware

import (
	"log/slog"

	"com.example/example/pkg/common"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

func GinLogger(c *gin.Context) {
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery
	requestID := c.GetHeader(common.RequestIDHeaderKey)
	if requestID == "" {
		requestID = xid.New().String()
		c.Request.Header.Set(common.RequestIDHeaderKey, requestID)
	}
	c.Set(common.RequestIdStringKey, requestID)
	// 设置响应头
	c.Header(common.RequestIDHeaderKey, requestID)
	slog.Info("incoming request",
		slog.String(common.RequestIdStringKey, requestID),
		"status", c.Writer.Status(),
		"method", c.Request.Method,
		"path", path,
		"query", query,
		"ip", c.ClientIP(),
		"userAgent", c.Request.UserAgent(),
	)
	c.Next()
}
