package model

import (
	"time"
)

type RoleModel struct {
	Id       int64     `json:"id"`                           // 主键id
	RoleCode string    `json:"roleCode" validate:"required"` // 角色编号
	RoleName string    `json:"roleName" validate:"required"` // 角色名称
	CreateAt time.Time `json:"createAt"`                     // 创建时间
	UpdateAt time.Time `json:"updateAt"`                     // 更新时间
	CreateBy int64     `json:"createBy"`                     // 创建人
	UpdateBy int64     `json:"updateBy"`                     // 更新人
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
