package repository

import (
	"context"

	"com.example/example/entity"
	"com.example/example/model"
	"com.example/example/pkg/exception"
	"com.example/example/pkg/orm"
	"gorm.io/gorm"
)

type RoleRepository struct {
	BaseRepository[entity.RoleEntity]
}

// NewRoleRepository 创建 RoleRepository
func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{
		NewBaseRepository[entity.RoleEntity](db),
	}
}

func (a *RoleRepository) ListPage(ctx context.Context, req model.RoleListReq) model.PageData[model.RoleModel] {
	db := a.GetDb(ctx)
	db = db.Model(&entity.RoleEntity{})
	if req.RoleCode != "" {
		db = db.Where("role_code like ?", "%"+req.RoleCode+"%")
	}
	if req.RoleName != "" {
		db = db.Where("role_name like ?", "%"+req.RoleName+"%")
	}
	pageData, _ := orm.QueryPage[model.RoleModel](db, req.GetPageNo(), req.GetPageSize())
	return pageData
}

// Add 添加角色
func (a *RoleRepository) Add(ctx context.Context, entity entity.RoleEntity) (entity.RoleEntity, error) {
	// 检查角色是否存在
	role := a.GetByCode(ctx, entity.RoleCode)
	if role.Id > 0 {
		return entity, exception.NewBizError("500", "角色已存在")
	}
	db := a.GetDb(ctx).Create(&entity)
	return entity, db.Error
}

// GetByCode 根据角色编号获取角色
func (a *RoleRepository) GetByCode(ctx context.Context, code string) entity.RoleEntity {
	var role entity.RoleEntity
	a.GetDb(ctx).Where("role_code=?", code).First(&role)
	return role
}

// DeleteById 删除角色
func (a *RoleRepository) DeleteById(ctx context.Context, id int64) error {
	ret := a.GetDb(ctx).Delete(&entity.RoleEntity{}, id)
	return ret.Error
}

// ListByIds 根据角色ID集合查询角色列表
func (a *RoleRepository) ListByIds(ctx context.Context, ids []int64) []entity.RoleEntity {
	var list []entity.RoleEntity
	a.GetDb(ctx).Find(&list, ids)
	return list
}
