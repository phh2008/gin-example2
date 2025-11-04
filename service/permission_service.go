package service

import (
	"context"

	"com.example/example/entity"
	"com.example/example/model"
	"com.example/example/model/result"
	"com.example/example/pkg/exception"
	"com.example/example/pkg/logger"
	"com.example/example/repository"
	"github.com/casbin/casbin/v2"
	"github.com/jinzhu/copier"
)

// PermissionService 权限服务
type PermissionService struct {
	permissionRepository *repository.PermissionRepository
	enforcer             *casbin.Enforcer
}

// NewPermissionService 创建服务
func NewPermissionService(
	permissionRepository *repository.PermissionRepository,
	enforcer *casbin.Enforcer,
) *PermissionService {
	return &PermissionService{
		permissionRepository: permissionRepository,
		enforcer:             enforcer,
	}
}

// ListPage 权限资源列表
func (a *PermissionService) ListPage(ctx context.Context, req model.PermissionListReq) *result.Result[model.PageData[model.PermissionModel]] {
	data := a.permissionRepository.ListPage(ctx, req)
	return result.Ok[model.PageData[model.PermissionModel]](data)
}

// Add 添加权限资源
func (a *PermissionService) Add(ctx context.Context, perm model.PermissionModel) *result.Result[entity.PermissionEntity] {
	var permission entity.PermissionEntity
	copier.Copy(&permission, &perm)
	res, err := a.permissionRepository.Add(ctx, permission)
	if err != nil {
		logger.Errorf("添加权限资源失败，%s", err.Error())
		return result.FailMsg[entity.PermissionEntity]("添加权限资源失败")
	}
	return result.Ok(res)
}

// Update 更新权限资源
func (a *PermissionService) Update(ctx context.Context, perm model.PermissionModel) *result.Result[*entity.PermissionEntity] {
	oldPerm, _ := a.permissionRepository.GetById(ctx, perm.Id)
	if oldPerm.Id == 0 {
		return result.Error[*entity.PermissionEntity](exception.NotFound)
	}
	var permission entity.PermissionEntity
	copier.Copy(&permission, &perm)
	// 更新权限资源表
	res, err := a.permissionRepository.Update(ctx, permission)
	if err != nil {
		logger.Errorf("更新权限资源失败，%s", err.Error())
		return result.FailMsg[*entity.PermissionEntity]("更新权限资源失败")
	}
	// 获取角色与资源列表,比如：[[systemAdmin /api/v1/user/list get] [guest /api/v1/user/list get]]
	perms, _ := a.enforcer.GetFilteredPolicy(1, oldPerm.Url, oldPerm.Action)
	// 更新casbin中的数据
	if len(perms) > 0 {
		for i, v := range perms {
			item := v
			item[1] = res.Url
			item[2] = res.Action
			perms[i] = item
		}
		_, err = a.enforcer.UpdateFilteredPolicies(perms, 1, oldPerm.Url, oldPerm.Action)
		if err != nil {
			logger.Errorf("更新casbin中的权限错误: %s", err)
		}
	}
	return result.Ok(&res)
}
