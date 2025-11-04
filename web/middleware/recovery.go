package middleware

import (
	"errors"
	"log/slog"
	"net"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"com.example/example/model/result"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// GinRecovery recover 掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					var se *os.SyscallError
					if errors.As(ne, &se) {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					slog.Error("recover", "path", c.Request.URL.Path, "error", err, "request", string(httpRequest))
					// If the connection is dead, we can't write a status to it.
					_ = c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}
				if stack {
					slog.Error("recover", "error", err, "request", string(httpRequest), "stack", string(debug.Stack()))
				} else {
					slog.Error("recover", "error", err, "request", string(httpRequest))
				}
				//c.AbortWithStatus(http.StatusInternalServerError)
				result.FailMsg[string](cast.ToString(err)).Response(c)
				c.Abort()
			}
		}()
		c.Next()
	}
}
