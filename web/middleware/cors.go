package middleware

import (
	"strings"
	"time"

	"com.example/example/pkg/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Cors 跨域请求
func Cors(conf config.Cors) gin.HandlerFunc {
	aware := cors.New(cors.Config{
		AllowOrigins:     conf.AllowedOriginPatterns,
		AllowMethods:     strings.Split(conf.AllowedMethods, ","),
		AllowHeaders:     strings.Split(conf.AllowedHeaders, ","),
		ExposeHeaders:    strings.Split(conf.ExposeHeaders, ","),
		AllowCredentials: conf.AllowCredentials,
		MaxAge:           time.Second * time.Duration(conf.MaxAge),
	})
	return aware
}

// DefaultCors 跨域请求
func DefaultCors(conf config.Cors) gin.HandlerFunc {
	return cors.Default()
}
