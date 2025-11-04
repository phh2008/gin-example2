package router

import (
	_ "com.example/example/docs"
	"com.example/example/pkg/config"
	"com.example/example/web/controller"
	"com.example/example/web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// ProviderSet is router provider set
var ProviderSet = wire.NewSet(wire.Struct(new(Router), "*"))

type Router struct {
	Engine        *gin.Engine
	HelloApi      *controller.HelloController
	Auth          *middleware.Auth
	UserApi       *controller.UserController
	PermissionApi *controller.PermissionController
	RoleApi       *controller.RoleController
}

// Register 注册路由
func (a *Router) Register() {
	v1 := a.Engine.Group("/api/v1")
	{
		//测试用的API，请删除
		v1.GET("/hello", a.HelloApi.Hello)
	}
	// API docs
	if config.Active == "dev" {
		v1.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	{
		//用户API
		v1.POST("/user/login", a.UserApi.Login)
		v1.GET("/user/list", a.Auth.Authorization("get"), a.UserApi.List)
		v1.POST("/user/createByEmail", a.UserApi.CreateByEmail)
		v1.POST("/user/assignRole", a.Auth.Authorization("post"), a.UserApi.AssignRole)
		v1.DELETE("/user/delete/:id", a.Auth.Authorization("delete"), a.UserApi.DeleteById)
		//角色API
		v1.GET("/role/list", a.Auth.Authorization("get"), a.RoleApi.List)
		v1.POST("/role/add", a.Auth.Authorization("post"), a.RoleApi.Add)
		v1.POST("/role/assignPermission", a.Auth.Authorization("post"), a.RoleApi.AssignPermission)
		v1.DELETE("/role/delete/:id", a.Auth.Authorization("delete"), a.RoleApi.Delete)
		//权限API
		v1.GET("/permission/list", a.Auth.Authorization("get"), a.PermissionApi.List)
		v1.POST("/permission/add", a.Auth.Authorization("post"), a.PermissionApi.Add)
		v1.POST("/permission/update", a.Auth.Authorization("post"), a.PermissionApi.Update)
	}
}
