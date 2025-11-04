package repository

import (
	"context"

	"com.example/example/entity"
	"gorm.io/gorm"
)

type RolePermissionRepository struct {
	BaseRepository[entity.RolePermissionEntity]
}

// NewRolePermissionRepository 创建 RolePermissionRepository
func NewRolePermissionRepository(db *gorm.DB) *RolePermissionRepository {
	return &RolePermissionRepository{
		NewBaseRepository[entity.RolePermissionEntity](db),
	}
}

func (a *RolePermissionRepository) DeleteByRoleId(ctx context.Context, roleId int64) error {
	db := a.GetDb(ctx).Where("role_id=?", roleId).Delete(&entity.RolePermissionEntity{})
	return db.Error
}

func (a *RolePermissionRepository) BatchAdd(ctx context.Context, list []*entity.RolePermissionEntity) error {
	if len(list) == 0 {
		return nil
	}
	db := a.GetDb(ctx).Create(list)
	return db.Error
}

func (a *RolePermissionRepository) ListRoleIdByPermId(ctx context.Context, permId int64) []int64 {
	var roleIds []int64
	a.GetDb(ctx).Model(&entity.RolePermissionEntity{}).
		Where("perm_id=?", permId).
		Pluck("role_id", &roleIds)
	return roleIds
}
