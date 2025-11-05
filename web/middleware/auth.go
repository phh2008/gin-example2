package middleware

import (
	"time"

	"com.example/example/manager"
	"com.example/example/model/result"
	"com.example/example/pkg/common"
	"com.example/example/pkg/exception"
	"com.example/example/pkg/xjwt"
	"github.com/gin-gonic/gin"
)

// Auth 权限中间件
type Auth struct {
	jwt              *xjwt.JwtHelper
	authorizeManager *manager.AuthorizeManager
}

// NewAuth 创建权限中间件
func NewAuth(jwt *xjwt.JwtHelper, authorizeManager *manager.AuthorizeManager) *Auth {
	return &Auth{jwt: jwt, authorizeManager: authorizeManager}
}

func (a *Auth) authValid(ctx *gin.Context) (xjwt.UserClaims, bool) {
	var user xjwt.UserClaims
	token := ctx.GetHeader(common.AuthTokenKey)
	if token == "" {
		result.Error[any](exception.NoLogin).Response(ctx)
		ctx.Abort()
		return user, false
	}
	jwtToken, err := a.jwt.VerifyToken(token)
	if err != nil {
		result.Error[any](exception.NoLogin).Response(ctx)
		ctx.Abort()
		return user, false
	}
	user, err = a.jwt.ParseToken(jwtToken)
	if err != nil {
		result.Error[any](exception.NoLogin).Response(ctx)
		ctx.Abort()
		return user, false
	}
	if !user.IsValidExpiresAt(time.Now()) {
		result.Error[any](exception.NoLogin).Response(ctx)
		ctx.Abort()
		return user, false
	}
	ctx.Set(common.UserKey{}, user)
	return user, true
}

// Authenticate 认证校验
func (a *Auth) Authenticate() gin.HandlerFunc {
	auth := func(ctx *gin.Context) {
		if _, ok := a.authValid(ctx); !ok {
			return
		}
		ctx.Next()
	}
	return auth
}

// Authorization 授权校验
func (a *Auth) Authorization(action string) gin.HandlerFunc {
	authorize := func(ctx *gin.Context) {
		// 是否已登录
		user, ok := a.authValid(ctx)
		if !ok {
			return
		}
		// 是否超级管理员
		ok = a.authorizeManager.HasAllRole(user.ID, []string{common.Admin})
		if ok {
			ctx.Next()
			return
		}
		// 是否有权限
		ok = a.authorizeManager.HasPermission(user.ID, action)
		if !ok {
			// 无权限
			result.Error[any](exception.Unauthorized).Response(ctx)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
	return authorize
}
