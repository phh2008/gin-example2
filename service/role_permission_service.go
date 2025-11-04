package service

import (
	"com.example/example/repository"
)

// RolePermissionService 角色权限关系
type RolePermissionService struct {
	rolePermissionRepository *repository.RolePermissionRepository
}

// NewRolePermissionService 创建服务
func NewRolePermissionService(rolePermissionRepository *repository.RolePermissionRepository) *RolePermissionService {
	return &RolePermissionService{
		rolePermissionRepository: rolePermissionRepository,
	}
}
