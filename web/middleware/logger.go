package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func GinLogger(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery
	c.Next()
	cost := time.Since(start)
	slog.Info("request",
		"status", c.Writer.Status(),
		"method", c.Request.Method,
		"path", path,
		"query", query,
		"ip", c.ClientIP(),
		"user-agent", c.Request.UserAgent(),
		"errors", c.Errors.ByType(gin.ErrorTypePrivate).String(),
		"cost", cost.String(),
	)
}
