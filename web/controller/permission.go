package controller

import (
	"errors"

	"com.example/example/model"
	"com.example/example/model/result"
	"com.example/example/service"
	"github.com/gin-gonic/gin"
)

// PermissionController 权限资源 api
type PermissionController struct {
	permissionService *service.PermissionService
}

// NewPermissionController 创建权限 api
func NewPermissionController(permissionService *service.PermissionService) *PermissionController {
	return &PermissionController{permissionService: permissionService}
}

// List 权限管理列表
//
//	@Summary		权限管理列表
//	@Description	权限管理列表
//	@Tags			权限
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			perm	query		model.PermissionListReq	false	"查询条件"
//	@Success		200		{object}	result.Result[model.PageData[model.PermissionModel]]
//	@Router			/permission/list [get]
func (a *PermissionController) List(ctx *gin.Context) {
	var req model.PermissionListReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		result.Error[any](err).Response(ctx)
		return
	}
	a.permissionService.ListPage(ctx, req).Response(ctx)
}

// Add 添加权限
//
//	@Summary		添加权限
//	@Description	添加权限
//	@Tags			权限
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			perm	body		model.PermissionModel	true	"权限信息"
//	@Success		200		{object}	result.Result[entity.PermissionEntity]
//	@Router			/permission/add [post]
func (a *PermissionController) Add(ctx *gin.Context) {
	var perm model.PermissionModel
	if err := ctx.ShouldBindJSON(&perm); err != nil {
		result.Error[any](err).Response(ctx)
		return
	}
	a.permissionService.Add(ctx, perm).Response(ctx)
}

// Update 更新权限
//
//	@Summary		更新权限
//	@Description	更新权限
//	@Tags			权限
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			perm	body		model.PermissionModel	true	"权限信息"
//	@Success		200		{object}	result.Result[entity.PermissionEntity]
//	@Router			/permission/add [post]
func (a *PermissionController) Update(ctx *gin.Context) {
	var perm model.PermissionModel
	if err := ctx.ShouldBindJSON(&perm); err != nil {
		result.Error[any](err).Response(ctx)
		return
	}
	if perm.Id == 0 {
		result.Error[any](errors.New("id不能为空")).Response(ctx)
		return
	}
	a.permissionService.Update(ctx, perm).Response(ctx)
}
