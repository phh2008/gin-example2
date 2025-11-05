package repository

import (
	"context"

	"com.example/example/entity"
	"com.example/example/model"
	"com.example/example/pkg/orm"
	"gorm.io/gorm"
)

type PermissionRepository struct {
	BaseRepository[entity.PermissionEntity]
}

// NewPermissionRepository 创建 PermissionRepository
func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{
		NewBaseRepository[entity.PermissionEntity](db),
	}
}

func (a *PermissionRepository) ListPage(ctx context.Context, req model.PermissionListReq) model.PageData[model.PermissionModel] {
	db := a.GetDb(ctx)
	db = db.Model(&entity.PermissionEntity{})
	if req.PermName != "" {
		db = db.Where("perm_name like ?", "%"+req.PermName+"%")
	}
	if req.Url != "" {
		db = db.Where("url=?", req.Url)
	}
	if req.Action != "" {
		db = db.Where("action=?", req.Action)
	}
	if req.PermType != 0 {
		db = db.Where("perm_type=?", req.PermType)
	}
	pageData, _ := orm.QueryPage[model.PermissionModel](db, req.GetPageNo(), req.GetPageSize())
	return pageData
}

func (a *PermissionRepository) Add(ctx context.Context, permission entity.PermissionEntity) (entity.PermissionEntity, error) {
	db := a.GetDb(ctx).Create(&permission)
	return permission, db.Error
}

func (a *PermissionRepository) Update(ctx context.Context, permission entity.PermissionEntity) (entity.PermissionEntity, error) {
	db := a.GetDb(ctx).Model(&permission).Updates(permission)
	return permission, db.Error
}

func (a *PermissionRepository) FindByIdList(idList []int64) []entity.PermissionEntity {
	var list []entity.PermissionEntity
	if len(idList) == 0 {
		return list
	}
	a.GetDb(context.TODO()).Find(&list, idList)
	return list
}

func (a *PermissionRepository) FindByUserId(ctx context.Context, userId int64) ([]entity.PermissionEntity, error) {
	var list []entity.PermissionEntity
	err := a.GetDb(ctx).Unscoped().Table("sys_user_role a").
		Select("c.*").
		Joins("JOIN sys_role_permission b ON a.role_id = b.role_id").
		Joins("JOIN sys_permission c ON b.perm_id = c.id AND c.deleted=1").
		Joins("JOIN sys_role d ON a.role_id=d.id AND d.deleted=1").
		Where("a.user_id=?", userId).
		Find(&list).Error
	return list, err
}
