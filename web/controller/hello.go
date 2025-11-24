package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HelloController hello api
type HelloController struct {
}

// NewHelloController new hello api
func NewHelloController() *HelloController {
	return &HelloController{}
}

// Hello 测试 api
func (a *HelloController) Hello(ctx *gin.Context) {
	ctx.String(http.StatusOK, "请求成功：success")
}
