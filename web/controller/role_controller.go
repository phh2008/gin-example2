package controller

import (
	"strconv"

	"com.example/example/model"
	"com.example/example/model/result"
	"com.example/example/pkg/exception"
	"com.example/example/service"
	"github.com/gin-gonic/gin"
)

// RoleController 角色 api
type RoleController struct {
	roleService *service.RoleService
	userService *service.UserService
}

// NewRoleController 创建角色 api
func NewRoleController(
	roleService *service.RoleService,
	userService *service.UserService,
) *RoleController {
	return &RoleController{
		roleService: roleService,
		userService: userService,
	}
}

// List 角色管理列表
//
//	@Summary		角色管理列表
//	@Description	角色管理列表
//	@Tags			角色
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			role	query		model.RoleListReq	false	"查询条件"
//	@Success		200		{object}	result.Result[model.PageData[model.RoleModel]]
//	@Router			/role/list [get]
func (a *RoleController) List(ctx *gin.Context) {
	var req model.RoleListReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		result.Error[any](err).Response(ctx)
		return
	}
	a.roleService.ListPage(ctx, req).Response(ctx)
}

// Add 添加角色
//
//	@Summary		添加角色
//	@Description	添加角色
//	@Tags			角色
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			role	body		model.RoleModel	false	"角色信息"
//	@Success		200		{object}	result.Result[entity.RoleEntity]
//	@Router			/role/add [post]
func (a *RoleController) Add(ctx *gin.Context) {
	var roleModel model.RoleModel
	if err := ctx.ShouldBindJSON(&roleModel); err != nil {
		result.Error[any](err).Response(ctx)
		return
	}
	a.roleService.Add(ctx, roleModel).Response(ctx)
}

// Delete 删除角色
//
//	@Summary		删除角色
//	@Description	删除角色
//	@Tags			角色
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int	true	"角色ID"
//	@Success		200	{object}	result.Result[any]
//	@Router			/role/delete/:id [delete]
func (a *RoleController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		result.Error[any](exception.ParamError)
		return
	}
	a.roleService.DeleteById(ctx, id).Response(ctx)
}

// AssignPermission 分配权限
//
//	@Summary		分配权限
//	@Description	分配权限
//	@Tags			角色
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			role	body		model.RoleAssignPermModel	true	"角色与权限信息"
//	@Success		200		{object}	result.Result[any]
//	@Router			/role/assignPermission [post]
func (a *RoleController) AssignPermission(ctx *gin.Context) {
	var assignModel model.RoleAssignPermModel
	if err := ctx.ShouldBindJSON(&assignModel); err != nil {
		result.Error[any](err).Response(ctx)
		return
	}
	a.roleService.AssignPermission(ctx, assignModel).Response(ctx)
}
