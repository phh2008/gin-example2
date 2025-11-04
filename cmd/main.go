package main

import (
	"flag"
	"log/slog"

	"com.example/example/pkg/config"
	"com.example/example/pkg/logger"
	"go.uber.org/zap/exp/zapslog"
)

// @title						示例API
// @version					    1.0
// @description				    示例API文档.
// @host						localhost:8089
// @BasePath					/api/v1
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {
	// 参数
	flag.StringVar(&config.Path, "config", "./config", "指定配置文件目录，示例: -config ./config")
	flag.StringVar(&config.Active, "active", "dev", "指定当前运行环境,示例: -env dev")
	flag.Parse()
	// 初始化配置
	conf := config.NewConfig(config.Path)
	// zap 配置
	zapLog := logger.NewLogger(conf)
	// slog 配置
	slogger := slog.New(zapslog.NewHandler(zapLog.Core(), zapslog.WithCaller(true)))
	slog.SetDefault(slogger)
	// wire
	server := BuildServer(conf)
	// 启动服务
	server.Start()
}
