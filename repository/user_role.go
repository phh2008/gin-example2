package repository

import (
	"context"

	"com.example/example/entity"
	"gorm.io/gorm"
)

type UserRoleRepository struct {
	BaseRepository[entity.UserRoleEntity]
}

func NewUserRoleRepository(db *gorm.DB) *UserRoleRepository {
	return &UserRoleRepository{
		NewBaseRepository[entity.UserRoleEntity](db),
	}
}

func (a *UserRoleRepository) CreateBatch(ctx context.Context, userRoles []entity.UserRoleEntity) error {
	return a.GetDb(ctx).Create(userRoles).Error
}

func (a *UserRoleRepository) FindByRoleId(ctx context.Context, roleId int64) ([]entity.UserRoleEntity, error) {
	var list []entity.UserRoleEntity
	err := a.GetDb(ctx).Model(&entity.UserRoleEntity{}).Where("role_id=?", roleId).Find(&list).Error
	return list, err
}

func (a *UserRoleRepository) FindRoleCodeByUserId(ctx context.Context, userId int64) ([]string, error) {
	var roleCodes []string
	err := a.GetDb(ctx).Model(&entity.UserRoleEntity{}).
		Joins("JOIN sys_role ON sys_user_role.role_id=sys_role.id").
		Where("sys_user_role.user_id=?", userId).Pluck("sys_role.role_code", &roleCodes).Error
	return roleCodes, err
}

func (a *UserRoleRepository) DeleteByUserId(ctx context.Context, userId int64) error {
	return a.GetDb(ctx).Unscoped().Where("user_id=?", userId).Delete(&entity.UserRoleEntity{}).Error
}
