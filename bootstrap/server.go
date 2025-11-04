package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"com.example/example/pkg/config"
	"com.example/example/web/middleware"
	"com.example/example/web/router"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/wire"
	"github.com/gookit/validate"
	"github.com/gookit/validate/locales/zhcn"
	"gorm.io/gorm"
)

var ServerSet = wire.NewSet(wire.Struct(new(Server), "*"), NewGin)

type Server struct {
	Config *config.Config
	DB     *gorm.DB
	Router *router.Router
	Engine *gin.Engine
}

func NewGin() *gin.Engine {
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		slog.Info("endpoint", "method", httpMethod, "path", absolutePath, "handlerName", handlerName, "nuHandlers", nuHandlers)
	}
	//gin.SetMode(gin.ReleaseMode)
	return gin.New()
}

func (a *Server) Start() {
	// 数据校验器
	zhcn.RegisterGlobal()
	validate.Config(func(opt *validate.GlobalOption) {
		opt.StopOnError = true // 如果为 true，则出现第一个错误时，将停止继续验证
		opt.SkipOnEmpty = true // 跳过对字段不存在或值为空的检查
	})
	binding.Validator = &customValidator{}
	// gin engine
	a.Engine.Use(middleware.GinLogger)
	a.Engine.Use(middleware.GinRecovery(true))
	a.Engine.Use(middleware.Cors(a.Config.Cors))
	a.Router.Register()
	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", a.Config.Server.Port),
		Handler: a.Engine.Handler(),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("启动服务出错", "err", err.Error())
			os.Exit(1)
		}
	}()
	addr := fmt.Sprintf("http://127.0.0.1:%s", a.Config.Server.Port)
	slog.Info("服务已启动", "addr", addr)
	slog.Info("接口文档", "docs", fmt.Sprintf("%s/api/v1/docs/index.html", addr))
	// 优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("服务关闭中 ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("服务关闭出错", "error", err)
	}
	slog.Info("服务已退出")
}

// implements the binding.StructValidator
type customValidator struct{}

func (c *customValidator) ValidateStruct(ptr any) error {
	v := validate.Struct(ptr)
	v.Validate() // 调用验证
	if v.Errors != nil && v.Errors.Error() == "" {
		return nil
	}
	return v.Errors
}

func (c *customValidator) Engine() any {
	return nil
}
