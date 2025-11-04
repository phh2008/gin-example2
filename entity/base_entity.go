package entity

import (
	"com.example/example/pkg/common"
	"com.example/example/pkg/types"
	"com.example/example/pkg/xjwt"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

func init() {
	soft_delete.FlagDeleted = 2
	soft_delete.FlagActived = 1
}

type BaseEntity struct {
	Id        int64                 `gorm:"primaryKey" json:"id"`                                              // 主键id
	CreatedAt types.LocalDateTime   `gorm:"autoCreateTime" json:"createdAt"`                                   // 创建时间
	UpdatedAt types.LocalDateTime   `gorm:"autoUpdateTime" json:"updatedAt"`                                   // 更新时间
	CreatedBy string                `json:"createdBy"`                                                         // 创建人
	UpdatedBy string                `json:"updatedBy"`                                                         // 更新人
	Deleted   soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:UpdatedAt;default:1" json:"deleted"` // 是否删除 1-否，2-是
}

func (a *BaseEntity) BeforeCreate(tx *gorm.DB) (err error) {
	a.Deleted = 1
	ctx := tx.Statement.Context
	user, ok := ctx.Value(common.UserKey{}).(xjwt.UserClaims)
	if ok {
		a.CreatedBy = user.ID
		a.UpdatedBy = a.CreatedBy
	}
	return
}

func (a *BaseEntity) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	user, ok := ctx.Value(common.UserKey{}).(xjwt.UserClaims)
	if ok {
		a.UpdatedBy = user.ID
	}
	return
}

func (a *BaseEntity) BeforeDelete(tx *gorm.DB) error {
	ctx := tx.Statement.Context
	user, ok := ctx.Value(common.UserKey{}).(xjwt.UserClaims)
	if ok {
		a.UpdatedBy = user.ID
	}
	return nil
}
