package model

import (
	"time"
)

type RoleModel struct {
	Id        int64     `json:"id"`                           // 主键id
	RoleCode  string    `json:"roleCode" validate:"required"` // 角色编号
	RoleName  string    `json:"roleName" validate:"required"` // 角色名称
	CreatedAt time.Time `json:"createdAt"`                    // 创建时间
	UpdatedAt time.Time `json:"updatedAt"`                    // 更新时间
	CreatedBy string    `json:"createdBy"`                    // 创建人
	UpdatedBy string    `json:"updatedBy"`                    // 更新人
}

type RoleListReq struct {
	QueryPage
	RoleCode string `json:"roleCode" form:"roleCode"` // 角色编号
	RoleName string `json:"roleName" form:"roleName"` // 角色名称
}

type RoleAssignPermModel struct {
	RoleId     int64   `json:"roleId" validate:"required"`     // 角色ID
	PermIdList []int64 `json:"permIdList" validate:"required"` // 权限ID列表
}
