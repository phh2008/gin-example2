package entity

type RolePermissionEntity struct {
	Id     int64 // 主键id
	RoleId int64 // 角色id
	PermId int64 // 权限id
}

func (RolePermissionEntity) TableName() string {
	return "sys_role_permission"
}
